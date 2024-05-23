package minecraft

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
)

type ServerConfigObject interface {
	Save(filepath string) error
	Load(filepath string) error
	LockAndGetData() any
	GetData() any
	Lock()
	Unlock()
}

func loadServerConfigObject(
	obj ServerConfigObject,
	objInit func(ServerConfigObject),
	filepath string,
) error {
	file, err := os.Open(filepath)

	if os.IsNotExist(err) {
		// create file if it doesn't exist
		objInit(obj)
		if err := file.Close(); err != nil {
			log.Println(err)
		}

		obj.Unlock()
		err := obj.Save(filepath)
		obj.Lock()

		return err
	} else if err != nil {
		if err := file.Close(); err != nil {
			log.Println(err)
		}

		return err
	}

	if err := json.NewDecoder(file).Decode(obj.GetData()); err != nil {
		if err := file.Close(); err != nil {
			log.Println(err)
		}

		return err
	}

	if err := file.Close(); err != nil {
		log.Println(err)
	}

	return nil
}

func saveServerConfigObject(obj ServerConfigObject, filepath string) error {
	file, err := os.Create(filepath)
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		return err
	}

	err = json.NewEncoder(file).Encode(obj.GetData())
	return err
}

func downloadServerVersion(filepath, version string) error {
	if err := loadVersionsMap("minecraft/data/server-download-links.json"); err != nil {
		return err
	}

	var info VersionInfo
	var ok bool
	if info, ok = versionsMap[version]; !ok {
		return errors.New("unsupported version")
	}

	serverJarPath := path.Join(filepath, fmt.Sprintf("server-%s.jar", version))
	serverJarBytes, err := os.ReadFile(serverJarPath)
	if err != nil {
		if err := exec.Command("wget", info.Link, "-O", serverJarPath).Run(); err != nil {
			return err
		}
		serverJarBytes, err = os.ReadFile(serverJarPath)
		if err != nil {
			return err
		}
	}

	// verify downloaded file using checksum
	md5Bytes := md5.Sum(serverJarBytes)
	if info.Sum == string(md5Bytes[:]) {
		return nil
	}

	return nil
}
