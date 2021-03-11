package system

import (
	"runtime"
	"strings"
)

// 获取 当前文件的路径
func CurrentFileName() string {
	_, file, _, ok := runtime.Caller(1)
	if ok {
		return file
	} else {
		return ""
	}
}

// 获取 当前文件的路径
func CurrentDirName() string {
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		index := strings.LastIndex(filename, "/")
		if index > 0 {
			return filename[0:index]
		} else {
			return filename
		}
	} else {
		return ""
	}
}
