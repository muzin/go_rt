package key_value_pair

import (
	"sync"
)

type KeyValuePair struct {
	mu sync.RWMutex

	table map[string]interface{}

	// 重新组建 Map 的阀值
	// 当 阀值 达到 删除次数开始重建 map，清理内存
	threshold int

	// 是否正在重建
	rebuilding bool

	// 删除 map中元素的次数
	delCount int
}

// 初始化
func (this *KeyValuePair) init() {

	this.table = make(map[string]interface{})

	// 重新组建 Map 的阀值
	this.threshold = 100

}

func NewKeyValuePair() *KeyValuePair {
	hashMap := &KeyValuePair{}
	hashMap.init()
	return hashMap
}

func (this *KeyValuePair) Size() int {
	this.mu.RLock()
	defer this.mu.RUnlock()

	i := len(this.table)
	return i
}

func (this *KeyValuePair) IsEmpty() bool {
	size := this.Size()
	if size == 0 {
		return true
	} else {
		return false
	}
}

func (this *KeyValuePair) Get(key string) interface{} {
	this.mu.RLock()
	defer this.mu.RUnlock()

	i := this.table[key]
	return i
}

func (this *KeyValuePair) Put(key string, value interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.table[key] = value
}

func (this *KeyValuePair) Remove(key string) {
	this.mu.Lock()
	defer this.mu.Unlock()

	_, contains := this.table[key]
	if contains {
		delete(this.table, key)
		this.delCount++
		// 每次移除后，检验是否重新组建Map
		isRebuild := this.IsRebuild()
		if isRebuild {
			this.rebuildMap()
		}
	}
}

func (this *KeyValuePair) RemoveValue(value interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()

	for k, v := range this.table {
		if v == value {
			delete(this.table, k)
			this.delCount++
			// 每次移除后，检验是否重新组建Map
			isRebuild := this.IsRebuild()
			if isRebuild {
				this.rebuildMap()
			}
		}
	}

}

func (this *KeyValuePair) ContainsKey(key string) bool {
	this.mu.RLock()
	defer this.mu.RUnlock()

	_, ok := this.table[key]
	if ok {
		return true
	} else {
		return false
	}
}

func (this *KeyValuePair) ContainsValue(value interface{}) bool {
	this.mu.RLock()
	defer this.mu.RUnlock()

	for _, v := range this.table {
		if v == value {
			return true
		}
	}
	return false
}

func (this *KeyValuePair) Keys() []string {
	this.mu.RLock()
	defer this.mu.RUnlock()

	keys := make([]string, 0)
	for k, _ := range this.table {
		keys = append(keys, k)
	}
	return keys
}

func (this *KeyValuePair) Values() []interface{} {
	this.mu.RLock()
	defer this.mu.RUnlock()

	values := make([]interface{}, 0)
	for _, v := range this.table {
		values = append(values, v)
	}
	return values
}

func (this *KeyValuePair) SetRebuildThreshold(value int) {
	this.threshold = value
}

func (this *KeyValuePair) GetThreshold() int {
	return this.threshold
}

// 是否 重新 组建
func (this *KeyValuePair) IsRebuild() bool {
	return this.isRebuild()
}

// 是否 重新 组建
func (this *KeyValuePair) isRebuild() bool {
	ret := false
	if this.delCount >= this.threshold {
		ret = true
	} else {
		ret = false
	}
	return ret
}

// 是否 正在重新组建s
func (this *KeyValuePair) IsRebuilding() bool {
	return this.rebuilding
}

// 重建 map
func (this *KeyValuePair) rebuildMap() {

	this.rebuilding = true

	// 创建新的map
	newMap := make(map[string]interface{})

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

func (this *KeyValuePair) Clear() {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.rebuilding = true

	// 创建新的map
	newMap := make(map[string]interface{})

	this.table = newMap

	// 重置 delCount
	this.delCount = 0

	this.rebuilding = false
}

func (this *KeyValuePair) Destroy() {
	go func() {
		this.table = nil
	}()
}
