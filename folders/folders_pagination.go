package folders

import (
	"encoding/base64"
	"fmt"
	"github.com/gofrs/uuid"
	"log"
)

// PaginatedFetchRequest paginated fetch request
type PaginatedFetchRequest struct {
	OrgID  uuid.UUID
	Limit  int    // items per page
	Cursor string // pagination cursor
}

// PaginatedFetchResponse paginated fetch response
type PaginatedFetchResponse struct {
	Folders    []*Folder
	NextCursor string
}

// GetPaginatedFolders retrieves a page of folders for a given organisation
// Input: PaginatedFetchRequest with OrgID, Limit, and Cursor
// Output: PaginatedFetchResponse with matching Folders and NextCursor, or error
func GetPaginatedFolders(req *PaginatedFetchRequest) (*PaginatedFetchResponse, error) {
	log.Printf("Received request for org: %s, limit: %d, cursor: %s", req.OrgID, req.Limit, req.Cursor)

	if req.OrgID == uuid.Nil {
		return nil, fmt.Errorf("invalid org ID: cannot be nil")
	}

	// default limit if not set or negative
	if req.Limit <= 0 {
		req.Limit = 10
		log.Printf("Using default limit: %d", req.Limit)
	}

	// maximum limit to prevent excessive data fetching
	if req.Limit > 100 {
		req.Limit = 100
		log.Printf("Limiting to maximum of 100 items per page")
	}

	// decode cursor
	lastID, err := decodeCursor(req.Cursor)
	if err != nil {
		log.Printf("Error decoding cursor: %v", err)
		lastID = uuid.Nil // fallback to start if cursor is invalid
	}

	// Fetch folders page
	paginatedFolders, nextCursor, err := fetchFoldersPage(req.OrgID, lastID, req.Limit)
	if err != nil {
		log.Printf("Error fetching folders: %v", err)
		return nil, fmt.Errorf("failed to fetch folders: %w", err)
	}

	log.Printf("Returning %d folders, next cursor: %s", len(paginatedFolders), nextCursor)

	return &PaginatedFetchResponse{
		Folders:    paginatedFolders,
		NextCursor: nextCursor,
	}, nil
}

// fetchFoldersPage simulates an efficient database query using cursor-based pagination
// Input: orgID (UUID), lastID (UUID), limit (int)
// Output: slice of Folder pointers, next cursor string, error
func fetchFoldersPage(orgID, lastID uuid.UUID, limit int) ([]*Folder, string, error) {
	allFolders := GetSampleData()

	var result []*Folder
	var lastFolder *Folder
	seenLastID := lastID == uuid.Nil

	for _, folder := range allFolders {
		if folder.OrgId != orgID {
			continue
		}

		if !seenLastID {
			if folder.Id == lastID {
				seenLastID = true
			}
			continue
		}

		result = append(result, folder)
		lastFolder = folder

		if len(result) == limit {
			break
		}
	}

	var nextCursor string
	if lastFolder != nil && lastFolder.Id != lastID {
		nextCursor = encodeCursor(lastFolder.Id)
	}

	return result, nextCursor, nil
}

// encodeCursor encodes the UUID into a base64 string
func encodeCursor(id uuid.UUID) string {
	return base64.StdEncoding.EncodeToString(id.Bytes())
}

// decodeCursor decodes the cursor back into a UUID
func decodeCursor(cursor string) (uuid.UUID, error) {
	if cursor == "" {
		return uuid.Nil, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid cursor format: %w", err)
	}

	if len(decoded) != 16 {
		return uuid.Nil, fmt.Errorf("invalid cursor content: incorrect length")
	}

	// Manual UUID validation
	for i, b := range decoded {
		if (i == 6 && b>>4 != 0x4) || (i == 8 && b>>6 != 0x2) {
			return uuid.Nil, fmt.Errorf("invalid cursor content: not a valid UUID")
		}
	}

	id, err := uuid.FromBytes(decoded)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid cursor content: %w", err)
	}

	return id, nil
}
