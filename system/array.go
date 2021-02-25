package system

// 数组拷贝
func ArrayCopy(src *[]interface{}, srcPos int, dest *[]interface{}, destPos int, length int) {
	srcSlice := (*src)[srcPos:]
	srcSliceLen := len(srcSlice)

	destLen := len(*dest)

	for i := 0; i < srcSliceLen; i++ {
		newDestPos := destPos + i
		if newDestPos < destLen && i < length {
			(*dest)[newDestPos] = srcSlice[i]
		}
	}

}

// int 数组拷贝
func IntArrayCopy(src *[]int, srcPos int, dest *[]int, destPos int, length int) {
	srcSlice := (*src)[srcPos:]
	srcSliceLen := len(srcSlice)

	destLen := len(*dest)

	for i := 0; i < srcSliceLen; i++ {
		newDestPos := destPos + i
		if newDestPos < destLen && i < length {
			(*dest)[newDestPos] = srcSlice[i]
		}
	}

}

// byte 数组拷贝
func ByteArrayCopy(src *[]byte, srcPos int, dest *[]byte, destPos int, length int) {
	srcSlice := (*src)[srcPos:]
	srcSliceLen := len(srcSlice)

	destLen := len(*dest)

	for i := 0; i < srcSliceLen; i++ {
		newDestPos := destPos + i
		if newDestPos < destLen && i < length {
			(*dest)[newDestPos] = srcSlice[i]
		}
	}

}
