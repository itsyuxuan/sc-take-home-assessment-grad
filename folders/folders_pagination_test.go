package folders

import (
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
		assert.NotEmpty(t, res.NextToken)
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
		assert.NotNil(t, res)
		assert.Empty(t, res.Folders)
		assert.Empty(t, res.NextToken)
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

	t.Run("fetch_with_invalid_token", func(t *testing.T) {
		req := &PaginatedFetchRequest{
			OrgID: testOrgID,
			Limit: 5,
			Token: "invalid_token",
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.Len(t, res.Folders, 5) // should start from beginning
	})

	t.Run("fetch_with_out_of_bounds_token", func(t *testing.T) {
		req := &PaginatedFetchRequest{
			OrgID: testOrgID,
			Limit: 5,
			Token: encodeToken(99999),
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.Empty(t, res.Folders)
		assert.Empty(t, res.NextToken)
	})
}

func TestTokenEncodingDecoding(t *testing.T) {
	t.Run("encode_and_decode", func(t *testing.T) {
		original := 42
		encoded := encodeToken(original)
		decoded, err := decodeToken(encoded)

		assert.NoError(t, err)
		assert.Equal(t, original, decoded)
	})

	t.Run("decode_empty_token", func(t *testing.T) {
		decoded, err := decodeToken("")

		assert.NoError(t, err)
		assert.Equal(t, 0, decoded)
	})

	t.Run("decode_invalid_token", func(t *testing.T) {
		_, err := decodeToken("invalid_token")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token format")
	})
}
