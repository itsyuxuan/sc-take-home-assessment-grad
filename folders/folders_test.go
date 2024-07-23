package folders

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// assertValidFolder is a helper function to check if a folder has valid properties
func assertValidFolder(t *testing.T, folder *Folder, expectedOrgID uuid.UUID) {
	assert.NotEqual(t, uuid.Nil, folder.Id, "folder ID should be valid")
	assert.Equal(t, expectedOrgID, folder.OrgId, "folder should belong to the requested org")
	assert.NotEmpty(t, folder.Name, "folder name should not be empty")
}

func TestGetAllFolders(t *testing.T) {
	testOrgID := uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a")

	// Test case: Fetch existing folders for a valid organisation
	t.Run("fetch_existing_folders", func(t *testing.T) {
		req := &FetchFolderRequest{OrgID: testOrgID}
		res, err := GetAllFolders(req)

		assert.NoError(t, err, "should not return an error for existing folders")
		assert.NotNil(t, res, "response should not be nil")
		assert.NotEmpty(t, res.Folders, "should return non-empty folder list")

		for _, folder := range res.Folders {
			assertValidFolder(t, folder, testOrgID)
		}
	})

	// Test case: Attempt to fetch folders for a non-existent organisation
	t.Run("fetch_nonexistent_folders", func(t *testing.T) {
		nonexistentOrgID := uuid.Must(uuid.NewV4())
		req := &FetchFolderRequest{OrgID: nonexistentOrgID}
		res, err := GetAllFolders(req)

		assert.Error(t, err, "should return an error for non-existent folders")
		assert.Nil(t, res, "response should be nil when no folders found")
		assert.Contains(t, err.Error(), "no folders found", "error message should indicate no folders were found")
	})

	// Test case: Attempt to fetch folders with a nil UUID
	t.Run("fetch_with_nil_uuid", func(t *testing.T) {
		req := &FetchFolderRequest{OrgID: uuid.Nil}
		res, err := GetAllFolders(req)

		assert.Error(t, err, "should return an error for nil UUID")
		assert.Nil(t, res, "response should be nil for invalid input")
		assert.Equal(t, "invalid organisation ID: nil UUID", err.Error(), "error message should indicate invalid input")
	})

	// Test case: Fetch a large number of folders to check performance
	t.Run("fetch_large_number_of_folders", func(t *testing.T) {
		req := &FetchFolderRequest{OrgID: testOrgID}
		res, err := GetAllFolders(req)

		assert.NoError(t, err, "should handle large number of folders without error")
		assert.NotNil(t, res, "response should not be nil for large number of folders")
		assert.True(t, len(res.Folders) > 100, "should return a large number of folders")
	})
}

func TestFetchAllFoldersByOrgID(t *testing.T) {
	testOrgID := uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a")

	// Test case: Fetch existing folders for a valid organisation
	t.Run("fetch_existing_folders", func(t *testing.T) {
		folders, err := FetchAllFoldersByOrgID(testOrgID)

		assert.NoError(t, err, "should not return an error for existing folders")
		assert.NotEmpty(t, folders, "should return non-empty folder list")

		for _, folder := range folders {
			assertValidFolder(t, folder, testOrgID)
		}
	})

	// Test case: Attempt to fetch folders for a non-existent organisation
	t.Run("fetch_nonexistent_folders", func(t *testing.T) {
		nonexistentOrgID := uuid.Must(uuid.NewV4())
		folders, err := FetchAllFoldersByOrgID(nonexistentOrgID)

		assert.Error(t, err, "should return an error for non-existent folders")
		assert.Nil(t, folders, "should return nil when no folders found")
		assert.Contains(t, err.Error(), "no folders found", "error message should indicate no folders were found")
	})

	// Test case: Attempt to fetch folders with a nil UUID
	t.Run("fetch_with_nil_uuid", func(t *testing.T) {
		folders, err := FetchAllFoldersByOrgID(uuid.Nil)

		assert.Error(t, err, "should return an error for nil UUID")
		assert.Nil(t, folders, "should return nil for invalid input")
		assert.Contains(t, err.Error(), "invalid organisation ID", "error message should indicate invalid input")
	})
}

func TestGetPaginatedFolders(t *testing.T) {
	testOrgID := uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a")

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

	t.Run("fetch_next_page", func(t *testing.T) {
		req := &PaginatedFetchRequest{
			OrgID: testOrgID,
			Limit: 5,
			Token: "5", // token from previous page
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.Len(t, res.Folders, 5)
		assert.NotEqual(t, res.Folders[0].Id, uuid.Nil)
	})

	t.Run("fetch_last_page", func(t *testing.T) {
		// first, get total number of folders
		allFolders, _ := FetchAllFoldersByOrgID(testOrgID)
		lastPageStart := len(allFolders) - 3 // assume last page has 3 items

		req := &PaginatedFetchRequest{
			OrgID: testOrgID,
			Limit: 5, // more than remaining items
			Token: fmt.Sprintf("%d", lastPageStart),
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.NotEmpty(t, res.Folders)
		assert.Len(t, res.Folders, 3)  // should get remaining 3 items
		assert.Empty(t, res.NextToken) // should be empty for last page
	})

	t.Run("fetch_beyond_last_page", func(t *testing.T) {
		req := &PaginatedFetchRequest{
			OrgID: testOrgID,
			Limit: 5,
			Token: "99999", // way beyond last page
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.Empty(t, res.Folders)
		assert.Empty(t, res.NextToken)
	})
}
