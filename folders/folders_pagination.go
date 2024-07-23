package folders

import (
	"encoding/base64"
	"fmt"
	"github.com/gofrs/uuid"
	"log"
	"strconv"
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
	log.Printf("Received request for org: %s, limit: %d, token: %s", req.OrgID, req.Limit, req.Token)

	if req.OrgID == uuid.Nil {
		return nil, fmt.Errorf("invalid org ID: cannot be nil")
	}

	// default limit if not set or negative
	if req.Limit <= 0 {
		req.Limit = 10
		log.Printf("Using default limit: %d", req.Limit)
	}

	// fetch all folders (in real world, we'd query the db with limit and offset)
	allFolders, err := FetchAllFoldersByOrgID(req.OrgID)
	if err != nil {
		log.Printf("Error fetching folders: %v", err)
		// Return an empty response instead of an error for non-existent org ID
		if err.Error() == fmt.Sprintf("no folders found for organisation ID %s", req.OrgID) {
			return &PaginatedFetchResponse{
				Folders:   []*Folder{},
				NextToken: "",
			}, nil
		}
		return nil, fmt.Errorf("failed to fetch folders: %w", err)
	}

	// decode start index from token
	startIndex, err := decodeToken(req.Token)
	if err != nil {
		log.Printf("Error decoding token: %v", err)
		startIndex = 0 // fallback to start if token is invalid
	}

	// handle out of bounds
	if startIndex >= len(allFolders) {
		log.Printf("Start index out of bounds: %d >= %d", startIndex, len(allFolders))
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

	log.Printf("Returning %d folders, next token: %s", len(paginatedFolders), nextToken)

	return &PaginatedFetchResponse{
		Folders:   paginatedFolders,
		NextToken: nextToken,
	}, nil
}

// encodeToken encodes the index into a base64 string
func encodeToken(index int) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", index)))
}

// decodeToken decodes the token back into an index
func decodeToken(token string) (int, error) {
	if token == "" {
		return 0, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return 0, fmt.Errorf("invalid token format: %w", err)
	}

	index, err := strconv.Atoi(string(decoded))
	if err != nil {
		return 0, fmt.Errorf("invalid token content: %w", err)
	}

	return index, nil
}
