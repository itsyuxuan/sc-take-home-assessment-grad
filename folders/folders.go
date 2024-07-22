package folders

import (
	"fmt"
	"github.com/gofrs/uuid"
)

// GetAllFolders retrieves folders for a given organisation
// Input: FetchFolderRequest with OrgID
// Output: FetchFolderResponse with matching Folders, or error
func GetAllFolders(req *FetchFolderRequest) (*FetchFolderResponse, error) {
	// fetch folders by org id
	folders, err := FetchAllFoldersByOrgID(req.OrgID)

	// error handling
	if err != nil {
		return nil, fmt.Errorf("failed to fetch folders: %w", err)
	}

	// create response directly with fetched folders
	return &FetchFolderResponse{Folders: folders}, nil
}

// FetchAllFoldersByOrgID gets folders for a specific org from sample data
// Input: orgID (UUID)
// Output: slice of matching Folder pointers, or error
func FetchAllFoldersByOrgID(orgID uuid.UUID) ([]*Folder, error) {
	// get all folders from sample data
	allFolders := GetSampleData()

	// estimate initial capacity to reduce allocations
	estimatedCapacity := len(allFolders) / 2 // assume about half will match
	resFolder := make([]*Folder, 0, estimatedCapacity)

	// loop through folders, append matches to resFolder
	for _, folder := range allFolders {
		if folder.OrgId == orgID {
			resFolder = append(resFolder, folder)
		}
	}

	// error handling
	if len(resFolder) == 0 {
		return nil, fmt.Errorf("no folders found for organization ID %s", orgID)
	}

	return resFolder, nil
}
