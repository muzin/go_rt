package system

import "runtime"

// 获取 当前文件的路径
func CurrentFile() string {
	_, file, _, ok := runtime.Caller(1)
	if ok {
		return file
	} else {
		return ""
	}
}
