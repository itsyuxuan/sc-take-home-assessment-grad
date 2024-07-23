package folders

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPaginatedFolders(t *testing.T) {
	testOrgID := uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a")

	// fetch first page
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

	// fetch next page
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

	// fetch last page
	t.Run("fetch_last_page", func(t *testing.T) {
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

	// fetch beyond last page
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

	// edge cases
	t.Run("fetch_with_limit_one", func(t *testing.T) {
		req := &PaginatedFetchRequest{
			OrgID: testOrgID,
			Limit: 1,
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.Len(t, res.Folders, 1)
		assert.NotEmpty(t, res.NextToken)
	})

	t.Run("fetch_with_large_limit", func(t *testing.T) {
		req := &PaginatedFetchRequest{
			OrgID: testOrgID,
			Limit: 1000,
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.NotEmpty(t, res.Folders)
		assert.Empty(t, res.NextToken)
	})

	t.Run("fetch_with_invalid_token", func(t *testing.T) {
		req := &PaginatedFetchRequest{
			OrgID: testOrgID,
			Limit: 5,
			Token: "invalid",
		}
		res, err := GetPaginatedFolders(req)

		assert.NoError(t, err)
		assert.Len(t, res.Folders, 5)
		assert.NotEmpty(t, res.NextToken)
	})
}
