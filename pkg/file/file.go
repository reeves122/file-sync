package file

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
)

func copySourceFiles(files []string, sourceDir, destDir string) error {
	for _, f := range files {
		sourcePath := filepath.Join(sourceDir, f)
		destPath := filepath.Join(destDir, f)
		log.Debugf("Copying %s to %s", sourcePath, destPath)
		if err := copyFile(sourcePath, destPath); err != nil {
			log.Error("error copying files from source")
			return err
		}
	}
	return nil
}

func copyFile(source, dest string) error {
	input, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	if baseDir, _ := filepath.Split(dest); baseDir != "" {
		if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(dest, input, 0644)
	return err
}
