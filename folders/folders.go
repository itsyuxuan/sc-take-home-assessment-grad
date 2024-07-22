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
	r, err := FetchAllFoldersByOrgID(req.OrgID)

	// error handling
	if err != nil {
		return nil, fmt.Errorf("failed to fetch folders: %w", err)
	}

	// create a slice of folder pointers
	fp := make([]*Folder, len(r))
	for i, v := range r {
		fp[i] = v
	}

	// create response with folder pointers
	ffr := &FetchFolderResponse{Folders: fp}

	return ffr, nil
}

// FetchAllFoldersByOrgID gets folders for a specific org from sample data
// Input: orgID (UUID)
// Output: slice of matching Folder pointers, or error
func FetchAllFoldersByOrgID(orgID uuid.UUID) ([]*Folder, error) {
	// get all folders from sample data
	folders := GetSampleData()

	// create slice for matching folders
	resFolder := []*Folder{}

	// loop through folders, append matches to resFolder
	for _, folder := range folders {
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
