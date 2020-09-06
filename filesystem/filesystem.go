package filesystem

import (
	"errors"
	"image-storage/app/errs"
	"io"
	"mime/multipart"
	"os"
)

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

// Delete delete file
func DeleteDir(dir string) error {
	return os.RemoveAll(dir)
}

// WriteOutputFile Writes the Output image file
func WriteImage(dir, fileName string, records multipart.File) (err error) {

	exist, _ := Exist(dir)
	if !exist {
		err = errors.New(errs.ErrAlbumNotExist)
		return err
	}

	existimg, _ := Exist(fileName)
	if existimg {
		err = errors.New(errs.ErrImageExist)
		return err
	}

	file, err := os.Create(fileName)
	defer file.Close()

	if _, err = io.Copy(file, records); err != nil {
		return err
	}
	return nil
}
