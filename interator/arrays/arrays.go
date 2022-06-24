package arrays

// Map
//
// 遍历元素返回新的数组集合
//
// @param arr []interface{}		数组
//
// @param iterator func(item interface{}, index int) interface{}		迭代器
//
// @param index int	索引
func Map(arr []interface{}, iterator func(item interface{}, index int) interface{}) []interface{} {
	if arr != nil {
		arrLen := len(arr)
		newArr := make([]interface{}, arrLen)
		for i := 0; i < arrLen; i++ {
			newArr[i] = iterator(arr[i], i)
		}
		return newArr
	} else {
		return make([]interface{}, 0)
	}
}

// Reduce
//
// 遍历元素返回新的集合
//
// @param arr []interface{}		数组
//
// @param iterator func(collection interface{}, item interface{}, index int) interface{}		迭代器
//
// @param collection interface{}	初始集合
func Reduce(arr []interface{},
	iterator func(collection interface{}, item interface{}, index int) interface{},
	collection interface{}) interface{} {
	if arr != nil {
		arrLen := len(arr)
		for i := 0; i < arrLen; i++ {
			collection = iterator(collection, arr[i], i)
		}
	}
	return collection
}
