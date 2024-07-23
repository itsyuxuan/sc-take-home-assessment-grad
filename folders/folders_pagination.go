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

	// decode start index from token
	startIndex, err := decodeToken(req.Token)
	if err != nil {
		log.Printf("Error decoding token: %v", err)
		startIndex = 0 // fallback to start if token is invalid
	}

	// Simulate efficient database query
	paginatedFolders, totalCount, err := fetchFoldersPage(req.OrgID, startIndex, req.Limit)
	if err != nil {
		log.Printf("Error fetching folders: %v", err)
		return nil, fmt.Errorf("failed to fetch folders: %w", err)
	}

	// set next token if more folders exist
	var nextToken string
	if startIndex+req.Limit < totalCount {
		nextToken = encodeToken(startIndex + req.Limit)
	}

	log.Printf("Returning %d folders, next token: %s", len(paginatedFolders), nextToken)

	return &PaginatedFetchResponse{
		Folders:   paginatedFolders,
		NextToken: nextToken,
	}, nil
}

// fetchFoldersPage simulates an efficient database query
func fetchFoldersPage(orgID uuid.UUID, offset, limit int) ([]*Folder, int, error) {
	allFolders := GetSampleData()

	var orgFolders []*Folder
	for _, folder := range allFolders {
		if folder.OrgId == orgID {
			orgFolders = append(orgFolders, folder)
		}
	}

	totalCount := len(orgFolders)

	if offset >= totalCount {
		return []*Folder{}, totalCount, nil
	}

	end := offset + limit
	if end > totalCount {
		end = totalCount
	}

	return orgFolders[offset:end], totalCount, nil
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
