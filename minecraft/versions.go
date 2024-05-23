package minecraft

import (
	"encoding/json"
	"log"
	"os"
)

type VersionInfo struct {
	Link string `json:"link"`
	Sum  string `json:"sum"` // md5 sum
}

type VersionMap map[string]VersionInfo

var versionsMap VersionMap

func loadVersionsMap(filepath string) error {
	file, err := os.Open(filepath)
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	if err != nil {
		return err
	}

	return json.NewDecoder(file).Decode(&versionsMap)
}
