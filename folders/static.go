package folders

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gofrs/uuid"
	"github.com/lucasepe/codename"
)

// These are all helper methods and fixed types.
// There's no real need for you to be editting these, but feel free to tweak it to suit your needs.
// If you do make changes here, be ready to discuss why these changes were made.

const dataSetSize = 1000
const DefaultOrgID = "c1556e17-b7c0-45a3-a6ae-9546248fb17a"

type Folder struct {
	// An unique identifier for the folder, must be a valid UUID.
	// For example: '00001d65-d336-485a-8331-7b53f37e8f51'
	Id uuid.UUID `json:"id"`
	// Name associated with folder.
	Name string `json:"name"`
	// The organisation that the folder belongs to.
	OrgId uuid.UUID `json:"org_id"`
	// Whether a folder has been marked as deleted or not.
	Deleted bool `json:"deleted"`
}

func GenerateData() []*Folder {
	rng, _ := codename.DefaultRNG()
	sampleData := []*Folder{}

	for i := 1; i < dataSetSize; i++ {
		orgId := uuid.FromStringOrNil(DefaultOrgID)

		if i%3 == 0 {
			orgId = uuid.Must(uuid.NewV4())
		}

		deleted := rand.Int() % 2

		sampleData = append(sampleData, &Folder{
			Id:      uuid.Must(uuid.NewV4()),
			Name:    codename.Generate(rng, 0),
			OrgId:   orgId,
			Deleted: deleted != 0,
		})
	}

	return sampleData
}

func PrettyPrint(b interface{}) {
	s, _ := json.MarshalIndent(b, "", "\t")
	fmt.Print(string(s))
}

func GetSampleData() []*Folder {
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println(filename)
	basePath := filepath.Dir(filename)
	filePath := filepath.Join(basePath, "sample.json")

	fmt.Println(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	jsonByte, _ := io.ReadAll(file)

	folders := []*Folder{}
	json.Unmarshal(jsonByte, &folders)
	return folders
}
