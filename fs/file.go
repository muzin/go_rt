package fs

import "os"

// 检查是否存在
func Exist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func Mkdir(path string, perm os.FileMode) (bool, error) {
	err := os.Mkdir(path, perm)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func Mkdirs(path string, perm os.FileMode) (bool, error) {
	err := os.MkdirAll(path, perm)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
