package cmdtoolbox

import "os"

func EnsureDir(dirPath string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		return err
	}

	return err
}
