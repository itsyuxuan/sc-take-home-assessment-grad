package folders

import (
	"fmt"
	"github.com/gofrs/uuid"
)

// PaginatedFetchRequest paginated fetch request
type PaginatedFetchRequest struct {
	OrgID uuid.UUID
	Limit int    // items per page
	Token string // pagination token
}

// PaginatedFetchResponse paginated fetch response
type PaginatedFetchResponse struct {
	Folders   []*Folder
	NextToken string
}

// GetPaginatedFolders get paginated folders
func GetPaginatedFolders(req *PaginatedFetchRequest) (*PaginatedFetchResponse, error) {
	if req.OrgID == uuid.Nil {
		return nil, fmt.Errorf("invalid org ID")
	}

	// default limit if not set
	if req.Limit <= 0 {
		req.Limit = 10
	}

	// fetch all folders (in real world, we'd query the db with limit and offset)
	allFolders, err := FetchAllFoldersByOrgID(req.OrgID)
	if err != nil {
		return nil, err
	}

	// decode start index from token
	startIndex := 0
	if req.Token != "" {
		startIndex = decodeToken(req.Token)
	}

	// handle out of bounds
	if startIndex >= len(allFolders) {
		return &PaginatedFetchResponse{
			Folders:   []*Folder{},
			NextToken: "",
		}, nil
	}

	// calculate end index
	endIndex := startIndex + req.Limit
	if endIndex > len(allFolders) {
		endIndex = len(allFolders)
	}

	paginatedFolders := allFolders[startIndex:endIndex]

	// set next token if more folders exist
	var nextToken string
	if endIndex < len(allFolders) {
		nextToken = encodeToken(endIndex)
	}

	return &PaginatedFetchResponse{
		Folders:   paginatedFolders,
		NextToken: nextToken,
	}, nil
}

// encode token (simple for now)
func encodeToken(index int) string {
	return fmt.Sprintf("%d", index)
}

// decode token (simple for now)
func decodeToken(token string) int {
	var index int
	fmt.Sscanf(token, "%d", &index)
	return index
}
