package filesystem

import "os"

// Exist checks if folder or file exist
func Exist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
	}
	return true, err
}

// Mkdir make directory if directory does not exist
func Mkdir(path string) error {
	if exist, _ := Exist(path); !exist {
		return os.Mkdir(path, 0777)
	}
	return nil
}

// Delete delete file
func Delete(filename string) error {
	return os.Remove(filename)
}
