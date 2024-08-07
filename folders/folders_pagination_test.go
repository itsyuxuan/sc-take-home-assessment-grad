package folders

import (
	"encoding/base64"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPaginatedFolders(t *testing.T) {
	testOrgID := uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a")
	anotherOrgID := uuid.FromStringOrNil("52214b35-f4da-461a-9f93-fbd3590e700f")
	nonExistentOrgID := uuid.Must(uuid.NewV4())

	t.Run("fetch_first_page", func(t *testing.T) {
		req := &PaginatedFetchRequest{
			OrgID: testOrgID,
			Limit: 5,
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.Len(t, res.Folders, 5)
		assert.NotEmpty(t, res.NextCursor)
	})

	t.Run("fetch_next_page", func(t *testing.T) {
		firstReq := &PaginatedFetchRequest{
			OrgID: testOrgID,
			Limit: 5,
		}
		firstRes, err := GetPaginatedFolders(firstReq)
		assert.NoError(t, err)

		secondReq := &PaginatedFetchRequest{
			OrgID:  testOrgID,
			Limit:  5,
			Cursor: firstRes.NextCursor,
		}
		secondRes, err := GetPaginatedFolders(secondReq)

		assert.NoError(t, err)
		assert.Len(t, secondRes.Folders, 5)
		assert.NotEqual(t, secondRes.Folders[0].Id, firstRes.Folders[0].Id)
	})

	t.Run("fetch_for_different_org", func(t *testing.T) {
		req := &PaginatedFetchRequest{
			OrgID: anotherOrgID,
			Limit: 5,
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.NotEmpty(t, res.Folders)
		for _, folder := range res.Folders {
			assert.Equal(t, anotherOrgID, folder.OrgId)
		}
	})

	t.Run("fetch_for_non_existent_org", func(t *testing.T) {
		req := &PaginatedFetchRequest{
			OrgID: nonExistentOrgID,
			Limit: 5,
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.Empty(t, res.Folders)
		assert.Empty(t, res.NextCursor)
	})

	t.Run("fetch_with_invalid_org_id", func(t *testing.T) {
		req := &PaginatedFetchRequest{
			OrgID: uuid.Nil,
			Limit: 5,
		}
		res, err := GetPaginatedFolders(req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "invalid org ID")
	})

	t.Run("fetch_with_negative_limit", func(t *testing.T) {
		req := &PaginatedFetchRequest{
			OrgID: testOrgID,
			Limit: -5,
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.Len(t, res.Folders, 10) // should use default limit
	})

	t.Run("fetch_with_zero_limit", func(t *testing.T) {
		req := &PaginatedFetchRequest{
			OrgID: testOrgID,
			Limit: 0,
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.Len(t, res.Folders, 10) // should use default limit
	})

	t.Run("fetch_with_small_limit", func(t *testing.T) {
		req := &PaginatedFetchRequest{
			OrgID: testOrgID,
			Limit: 1,
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.Len(t, res.Folders, 1)
		assert.NotEmpty(t, res.NextCursor)
	})

	t.Run("fetch_with_limit_exceeding_maximum", func(t *testing.T) {
		req := &PaginatedFetchRequest{
			OrgID: testOrgID,
			Limit: 150,
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.Len(t, res.Folders, 100) // should use maximum limit
		assert.NotEmpty(t, res.NextCursor)
	})

	t.Run("fetch_with_invalid_cursor", func(t *testing.T) {
		req := &PaginatedFetchRequest{
			OrgID:  testOrgID,
			Limit:  5,
			Cursor: "invalid_cursor",
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.Folders, 5) // should start from beginning
	})

	t.Run("fetch_with_non_existent_folder_cursor", func(t *testing.T) {
		nonExistentID := uuid.Must(uuid.NewV4())
		cursor := encodeCursor(nonExistentID)

		req := &PaginatedFetchRequest{
			OrgID:  testOrgID,
			Limit:  10,
			Cursor: cursor,
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.Empty(t, res.Folders)    // Expecting empty folders
		assert.Empty(t, res.NextCursor) // Expecting empty next cursor
	})

	t.Run("fetch_all_pages", func(t *testing.T) {
		allFolders, _ := FetchAllFoldersByOrgID(testOrgID)
		totalFolders := len(allFolders)

		var fetchedFolders []*Folder
		var nextCursor string

		for {
			req := &PaginatedFetchRequest{
				OrgID:  testOrgID,
				Limit:  10,
				Cursor: nextCursor,
			}
			res, err := GetPaginatedFolders(req)

			assert.NoError(t, err)
			fetchedFolders = append(fetchedFolders, res.Folders...)

			if res.NextCursor == "" {
				break
			}
			nextCursor = res.NextCursor
		}

		assert.Len(t, fetchedFolders, totalFolders)
	})

	t.Run("fetch_with_last_folder_cursor", func(t *testing.T) {
		allFolders, _ := FetchAllFoldersByOrgID(testOrgID)
		lastFolder := allFolders[len(allFolders)-1]
		lastCursor := encodeCursor(lastFolder.Id)

		req := &PaginatedFetchRequest{
			OrgID:  testOrgID,
			Limit:  10,
			Cursor: lastCursor,
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.Empty(t, res.Folders)
		assert.Empty(t, res.NextCursor)
	})

	t.Run("fetch_empty_page", func(t *testing.T) {
		emptyOrgID := uuid.Must(uuid.NewV4())
		req := &PaginatedFetchRequest{
			OrgID: emptyOrgID,
			Limit: 10,
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.Empty(t, res.Folders)
		assert.Empty(t, res.NextCursor)
	})

	t.Run("fetch_last_page_with_fewer_items", func(t *testing.T) {
		allFolders, _ := FetchAllFoldersByOrgID(testOrgID)
		lastPageStart := len(allFolders) - 3 // Assuming there are at least 3 folders
		lastPageCursor := encodeCursor(allFolders[lastPageStart].Id)

		req := &PaginatedFetchRequest{
			OrgID:  testOrgID,
			Limit:  10,
			Cursor: lastPageCursor,
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.Len(t, res.Folders, 2)      // Expecting 2 folders instead of 3
		assert.NotEmpty(t, res.NextCursor) // Expecting a non-empty next cursor
	})

	t.Run("default_limit", func(t *testing.T) {
		req := &PaginatedFetchRequest{
			OrgID: testOrgID,
			Limit: -1,
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 10, req.Limit)
	})

	t.Run("max_limit", func(t *testing.T) {
		req := &PaginatedFetchRequest{
			OrgID: testOrgID,
			Limit: 200,
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 100, req.Limit)
	})
}

func TestCursorEncodingDecoding(t *testing.T) {
	t.Run("encode_and_decode", func(t *testing.T) {
		original := uuid.Must(uuid.NewV4())
		encoded := encodeCursor(original)
		decoded, err := decodeCursor(encoded)

		assert.NoError(t, err)
		assert.Equal(t, original, decoded)
	})

	t.Run("decode_empty_cursor", func(t *testing.T) {
		decoded, err := decodeCursor("")

		assert.NoError(t, err)
		assert.Equal(t, uuid.Nil, decoded)
	})

	t.Run("decode_invalid_cursor", func(t *testing.T) {
		_, err := decodeCursor("invalid_cursor")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid cursor format")
	})

	t.Run("decode_invalid_base64", func(t *testing.T) {
		_, err := decodeCursor("ThisIsNotBase64!")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid cursor format")
	})

	t.Run("decode_invalid_uuid", func(t *testing.T) {
		invalidUUID := base64.StdEncoding.EncodeToString([]byte("0123456789abcdef"))
		_, err := decodeCursor(invalidUUID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid cursor content")
	})
}
