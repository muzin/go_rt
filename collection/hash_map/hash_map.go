package hash_map

import (
	"sync"
)

type HashMap struct {
	mu sync.RWMutex

	table map[interface{}]interface{}

	// 重新组建 Map 的阀值
	// 当 阀值 达到 删除次数开始重建 map，清理内存
	threshold int

	// 是否正在重建
	rebuilding bool

	// 删除 map中元素的次数
	delCount int
}

// 初始化
func (this *HashMap) init() {

	this.table = make(map[interface{}]interface{})

	// 重新组建 Map 的阀值
	this.threshold = 100

}

func NewHashMap() *HashMap {
	hashMap := &HashMap{}
	hashMap.init()
	return hashMap
}

func (this *HashMap) Size() int {
	i := len(this.table)
	return i
}

func (this *HashMap) IsEmpty() bool {
	size := this.Size()
	if size == 0 {
		return true
	} else {
		return false
	}
}

func (this *HashMap) Get(key interface{}) interface{} {
	this.mu.RLock()
	defer this.mu.RUnlock()

	return this.table[key]
}

func (this *HashMap) Put(key interface{}, value interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.table[key] = value
}

func (this *HashMap) Remove(key interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()

	delete(this.table, key)
	this.delCount++
	// 每次移除后，检验是否重新组建Map
	isRebuild := this.IsRebuild()
	if isRebuild {
		this.rebuildMap()
	}
}

func (this *HashMap) ContainsKey(key interface{}) bool {
	this.mu.RLock()
	defer this.mu.RUnlock()

	_, ok := this.table[key]
	if ok {
		return true
	} else {
		return false
	}
}

func (this *HashMap) ContainsValue(value interface{}) bool {
	this.mu.RLock()
	defer this.mu.RUnlock()

	for _, v := range this.table {
		if v == value {
			return true
		}
	}
	return false
}

func (this *HashMap) Keys() []interface{} {
	this.mu.RLock()
	defer this.mu.RUnlock()

	keys := make([]interface{}, 0)
	for k, _ := range this.table {
		keys = append(keys, k)
	}
	return keys
}

func (this *HashMap) Values() []interface{} {
	this.mu.RLock()
	defer this.mu.RUnlock()

	values := make([]interface{}, 0)
	for _, v := range this.table {
		values = append(values, v)
	}
	return values
}

func (this *HashMap) SetRebuildThreshold(value int) {
	this.threshold = value
}

func (this *HashMap) GetThreshold() int {
	return this.threshold
}

// 是否 重新 组建
func (this *HashMap) IsRebuild() bool {
	ret := false
	if this.delCount >= this.threshold {
		ret = true
	} else {
		ret = false
	}
	return ret
}

// 是否 正在重新组建s
func (this *HashMap) IsRebuilding() bool {
	return this.rebuilding
}

// 重建 map
func (this *HashMap) rebuildMap() {

	this.rebuilding = true

	// 创建新的map
	newMap := make(map[interface{}]interface{})

	// 将 旧 map 的值 拷贝到 新 map 中
	for k, v := range this.table {
		newMap[k] = v
	}

	// 删除 map 中 的每一项
	for k, _ := range this.table {
		delete(this.table, k)
	}

	this.table = newMap

	// 重置 delCount
	this.delCount = 0

	this.rebuilding = false

}

func (this *HashMap) Destroy() {
	go func() {
		this.table = nil
	}()
}
