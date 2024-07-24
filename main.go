package main

import (
	"fmt"
	"github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
)

func main() {
	req := &folders.FetchFolderRequest{
		OrgID: uuid.FromStringOrNil(folders.DefaultOrgID),
	}

	res, err := folders.GetAllFolders(req)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	folders.PrettyPrint(res)

	//req := &folders.PaginatedFetchRequest{
	//	OrgID:  uuid.FromStringOrNil(folders.DefaultOrgID),
	//	Limit:  10, // specify the limit for pagination
	//	Cursor: "", // empty cursor for the first page
	//}
	//
	//for {
	//	res, err := folders.GetPaginatedFolders(req)
	//	if err != nil {
	//		fmt.Printf("%v", err)
	//		return
	//	}
	//
	//	folders.PrettyPrint(res.Folders)
	//
	//	// if the NextCursor is empty, the last page is reached
	//	if res.NextCursor == "" {
	//		break
	//	}
	//
	//	// set the cursor for the next request to fetch the next page
	//	req.Cursor = res.NextCursor
	//}
}
