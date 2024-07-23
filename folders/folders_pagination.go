package folders

import (
	"fmt"
	"github.com/gofrs/uuid"
)

// PaginatedFetchRequest adds pagination params to the original request
type PaginatedFetchRequest struct {
	OrgID uuid.UUID
	Limit int    // number of items per page
	Token string // pagination token for fetching next page
}

// PaginatedFetchResponse includes a token for the next page
type PaginatedFetchResponse struct {
	Folders   []*Folder
	NextToken string
}

// GetPaginatedFolders fetches folders with pagination
func GetPaginatedFolders(req *PaginatedFetchRequest) (*PaginatedFetchResponse, error) {
	if req.OrgID == uuid.Nil {
		return nil, fmt.Errorf("invalid org ID")
	}

	// use a sensible default if limit is not set
	if req.Limit <= 0 {
		req.Limit = 10
	}

	allFolders, err := FetchAllFoldersByOrgID(req.OrgID)
	if err != nil {
		return nil, err
	}

	// decode the starting index from the token
	startIndex := 0
	if req.Token != "" {
		startIndex = decodeToken(req.Token)
	}

	// make sure we don't go out of bounds
	if startIndex >= len(allFolders) {
		return &PaginatedFetchResponse{
			Folders:   []*Folder{},
			NextToken: "",
		}, nil
	}

	// calculate the end index, ensuring we don't exceed the slice bounds
	endIndex := startIndex + req.Limit
	if endIndex > len(allFolders) {
		endIndex = len(allFolders)
	}

	paginatedFolders := allFolders[startIndex:endIndex]

	// set the next token if there are more folders to fetch
	var nextToken string
	if endIndex < len(allFolders) {
		nextToken = encodeToken(endIndex)
	}

	return &PaginatedFetchResponse{
		Folders:   paginatedFolders,
		NextToken: nextToken,
	}, nil
}

// encodeToken converts an index to a string token
func encodeToken(index int) string {
	return fmt.Sprintf("%d", index)
}

// decodeToken converts a string token back to an index
func decodeToken(token string) int {
	var index int
	fmt.Sscanf(token, "%d", &index)
	return index
}
