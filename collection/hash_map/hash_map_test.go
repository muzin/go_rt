package hash_map

import (
	"fmt"
	"runtime"
	"testing"
)

func TestHashMap_ContainsKey(t *testing.T) {
	t.Run("测试 HashMap", func(t *testing.T) {

		hashMap := NewHashMap()

		for i := 0; i < 5000; i++ {
			hashMap.Put(i, i)
		}
		runtime.GC()
		_printMemStats("add 5000")

		size := hashMap.Size()
		t.Logf("%v\n", size)

		for i := 0; i < 500; i++ {
			hashMap.Remove(i)
		}

		runtime.GC()
		_printMemStats("remove 500")

		size = hashMap.Size()
		t.Logf("%v\n", size)

		for i := 500; i < 4500; i++ {
			hashMap.Remove(i)
		}

		runtime.GC()
		_printMemStats("remove 4000")

		size = hashMap.Size()
		t.Logf("%v\n", size)

	})
}

func _printMemStats(mag string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%v：memory = %vKB, GC Times = %v\n", mag, m.Alloc/1024, m.NumGC)
}
