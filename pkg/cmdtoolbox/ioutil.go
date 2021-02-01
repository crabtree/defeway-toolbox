package cmdtoolbox

import (
	"io/ioutil"
	"os"
)

func EnsureDir(dirPath string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		return err
	}

	return err
}

func FileExists(path string) bool {
	if path == "" {
		return false
	}

	f, err := os.Stat(path)
	if err != nil {
		return false
	}

	return f.IsDir() == false
}

func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
