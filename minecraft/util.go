package minecraft

import (
	"io"
	"log"
	"os"
	"text/template"

	"github.com/raian621/minecraft-server-controller/api"
)

// func downloadServerVersion(filepath, version string) error {
// 	file, err := os.Open("minecraft/data/server-download-links.json")
// 	defer func() {
// 		if err := file.Close(); err != nil {
// 			log.Println(err)
// 		}
// 	}()

// 	var versionsMap VersionMap
// 	if err != nil {
// 		return err
// 	}

// 	if err := versionsMap.Load(file); err != nil {
// 		return err
// 	}

// 	var (
// 		info VersionInfo
// 		ok   bool
// 	)
// 	if info, ok = versionsMap[version]; !ok {
// 		return errors.New("unsupported version")
// 	}

// 	serverJarPath := path.Join(filepath, fmt.Sprintf("server-%s.jar", version))
// 	serverJarBytes, err := os.ReadFile(serverJarPath)
// 	if err != nil {
// 		if err := exec.Command("wget", info.Link, "-O", serverJarPath).Run(); err != nil {
// 			return err
// 		}
// 		serverJarBytes, err = os.ReadFile(serverJarPath)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	// verify downloaded file using checksum
// 	md5Bytes := md5.Sum(serverJarBytes)
// 	if info.Sum == string(md5Bytes[:]) {
// 		return nil
// 	}

// 	return nil
// }

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		log.Println("unexpected error closing file:", err)
	}
}

func saveServerPropertiesTemplate(props *api.ServerProperties, tmplFilepath, propertiesFilepath string) error {
	tmpl, err := template.New("server.properties.tmpl").Parse(tmplFilepath)
	if err != nil {
		return err
	}

	file, err := os.Create(propertiesFilepath)
	if err != nil {
		return err
	}
	defer closeFile(file)
	if err := tmpl.Execute(file, *props); err != nil {
		return err
	}

	return nil
}

func loadJSON(
	loadFn func(file io.Reader) error,
	saveFn func(file io.Writer) error,
	createFn func(),
	filepath string,
) error {
	file, err := os.Open(filepath)
	if os.IsNotExist(err) {
		file, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer closeFile(file)
		createFn()
		err = saveFn(file)
		if err != nil {
			return err
		}
		return nil
	}
	defer closeFile(file)
	return loadFn(file)
}

func ref[T any](v T) *T { return &v }
