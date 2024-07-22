package folders

import (
	"github.com/gofrs/uuid"
)

// GetAllFolders retrieves folders for a given organisation
// Input: FetchFolderRequest with OrgID
// Output: FetchFolderResponse with matching Folders, or error
func GetAllFolders(req *FetchFolderRequest) (*FetchFolderResponse, error) {
	// TODO: declared but unused vars
	var (
		err error
		f1  Folder
		fs  []*Folder
	)

	// create an empty slice of folders
	f := []Folder{}

	// fetch folders by org id, ignoring potential error (seems risky?)
	r, _ := FetchAllFoldersByOrgID(req.OrgID)

	// loop through fetched folders and append each folder to f slice
	// k is unused
	for k, v := range r {
		f = append(f, *v)
	}

	// create another slice of folder pointers
	var fp []*Folder

	// k1 is unused
	for k1, v1 := range f {
		fp = append(fp, &v1)
	}

	// create response with folder pointers
	var ffr *FetchFolderResponse
	ffr = &FetchFolderResponse{Folders: fp}

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
	return resFolder, nil
}
