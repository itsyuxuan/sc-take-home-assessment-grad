package folders

import (
	"fmt"
	"github.com/gofrs/uuid"
)

// GetAllFolders retrieves folders for a given organisation
// Input: FetchFolderRequest with OrgID
// Output: FetchFolderResponse with matching Folders, or error
func GetAllFolders(req *FetchFolderRequest) (*FetchFolderResponse, error) {
	// check for nil UUID
	if req.OrgID == uuid.Nil {
		return nil, fmt.Errorf("invalid organisation ID: nil UUID")
	}

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
	if orgID == uuid.Nil {
		return nil, fmt.Errorf("invalid organisation ID: nil UUID")
	}

	// get all folders from sample data
	allFolders := GetSampleData()

	// create slice for matching folders
	var resFolder []*Folder

	// loop through folders, append matches to resFolder
	for _, folder := range allFolders {
		if folder.OrgId == orgID {
			resFolder = append(resFolder, folder)
		}
	}

	// error handling
	if len(resFolder) == 0 {
		return nil, fmt.Errorf("no folders found for organisation ID %s", orgID)
	}

	return resFolder, nil
}
