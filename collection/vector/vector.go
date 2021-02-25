package vector

import (
	"github.com/muzin/go/rt/system"
	"sync"
)

const (
	MAX_ARRAY_SIZE    = int(^uint(0)>>1) - 8
	INTEGER_MAX_VALUE = int(^uint(0) >> 1)
)

type Vector struct {
	mu sync.Mutex

	elementData []interface{}

	elementCount int

	capacityIncrement int
}

func NewVector() *Vector {
	return NewInstanceOfVector(10, 0)
}

func NewInstanceOfVector(initialCapacity int, capacityIncrement int) *Vector {
	return &Vector{
		elementData:       make([]interface{}, initialCapacity),
		capacityIncrement: capacityIncrement,
	}
}

func (this *Vector) grow() *[]interface{} {
	return this.growByMinCapacity(this.elementCount + 1)
}

func (this *Vector) growByMinCapacity(minCapacity int) *[]interface{} {
	i := make([]interface{}, this.newCapacity(minCapacity))
	// copy 原始值 到 新数组中
	system.ArrayCopy(&this.elementData, 0, &i, 0, len(this.elementData))
	return &i
}

func (this *Vector) newCapacity(minCapacity int) int {
	// overflow-conscious code
	oldCapacity := len(this.elementData)
	newCapacity := oldCapacity
	if this.capacityIncrement > 0 {
		newCapacity += this.capacityIncrement
	} else {
		newCapacity += oldCapacity
	}

	if newCapacity-minCapacity <= 0 {
		if minCapacity < 0 { // overflow
			panic("OutOfMemoryError")
		} else {
			return minCapacity
		}
	}

	var ret = 0
	if newCapacity-MAX_ARRAY_SIZE <= 0 {
		ret = newCapacity
	} else {
		ret = this.hugeCapacity(minCapacity)
	}
	return ret
}

func (this *Vector) hugeCapacity(minCapacity int) int {
	if minCapacity < 0 { // overflow
		panic("OutOfMemoryError()")
	}
	var ret = 0
	if minCapacity > MAX_ARRAY_SIZE {
		ret = INTEGER_MAX_VALUE
	} else {
		ret = MAX_ARRAY_SIZE
	}
	return ret
}

func (this *Vector) Capacity() int {
	this.mu.Lock()
	defer this.mu.Unlock()

	return len(this.elementData)
}

func (this *Vector) Size() int {
	this.mu.Lock()
	defer this.mu.Unlock()

	return this.elementCount
}

func (this *Vector) IsEmpty() bool {
	this.mu.Lock()
	defer this.mu.Unlock()

	return this.elementCount == 0
}

func (this *Vector) IndexOf(o interface{}) int {
	return this.IndexOfWithIndex(o, 0)
}

func (this *Vector) IndexOfWithIndex(o interface{}, index int) int {
	this.mu.Lock()
	defer this.mu.Unlock()

	if o == nil {
		for i := index; i < this.elementCount; i++ {
			if nil == this.elementData[i] {
				return i
			}
		}
	} else {
		for i := index; i < this.elementCount; i++ {
			if o == this.elementData[i] {
				return i
			}
		}
	}
	return -1
}

func (this *Vector) LastIndexOf(o interface{}, index int) int {
	return this.LastIndexOfWithIndex(o, 0)
}

func (this *Vector) LastIndexOfWithIndex(o interface{}, index int) int {
	this.mu.Lock()
	defer this.mu.Unlock()

	if o == nil {
		for i := index; i >= 0; i-- {
			if nil == this.elementData[i] {
				return i
			}
		}
	} else {
		for i := index; i >= 0; i-- {
			if o == this.elementData[i] {
				return i
			}
		}
	}
	return -1
}

func (this *Vector) addToElementData(o *interface{}, elementData *[]interface{}, s int) {
	this.mu.Lock()
	defer this.mu.Unlock()

	//fmt.Printf("addToElementData o p: %v %v\n", o, &o)

	if s == len(*elementData) {
		this.elementData = *(this.grow())
		this.elementData[s] = o
	} else {
		(*elementData)[s] = o
	}

	//fmt.Printf("elementData o p: %v %v %v\n", elementData, &this.elementData, &elementData)
	this.elementCount = s + 1
}

//  Appends the specified element to the end of this Vector.
func (this *Vector) Add(o *interface{}) bool {
	//fmt.Printf("add o p: %v %v\n", o, &o)
	this.addToElementData(o, &this.elementData, this.elementCount)
	return true
}

func (this *Vector) removeElement(o *interface{}) bool {
	this.mu.Lock()
	defer this.mu.Unlock()

	i := this.IndexOf(o)
	if i >= 0 {
		this.removeElementAt(i)
		return true
	} else {
		return false
	}
}

func (this *Vector) removeElementAt(index int) {
	this.mu.Lock()
	defer this.mu.Unlock()

	if index >= this.elementCount {
		panic("ArrayIndexOutOfBoundsException " + string(index) + ">=" + string(this.elementCount))
	} else if index < 0 {
		panic("ArrayIndexOutOfBoundsException " + string(index))
	}

	var j = this.elementCount - index - 1
	if j > 0 {
		system.ArrayCopy(&this.elementData, index+1, &this.elementData, index, j)
	}
	this.elementCount--
	this.elementData[this.elementCount] = nil
}

func (this *Vector) Remove(index int) (val *interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.elementCount == 0 && index == this.elementCount {
		return nil
	}

	if index >= this.elementCount || index < 0 {
		panic("ArrayIndexOutOfBoundsException size: " + string(this.elementCount) + " index: " + string(index))
	}

	oldValue := this.elementData[index]

	//fmt.Printf("oldValue: %v ptr: %v", oldValue, &oldValue)

	var numMoved = this.elementCount - index - 1
	if numMoved > 0 {
		system.ArrayCopy(&this.elementData, index+1, &this.elementData, index, numMoved)
	}

	this.elementCount -= 1
	this.elementData[this.elementCount] = nil // Let gc do its work

	return oldValue.(*interface{})
}

func (this *Vector) removeAllElements() {
	this.mu.Lock()
	defer this.mu.Unlock()

	es := this.elementData
	oldElementCount := this.elementCount
	this.elementCount = 0
	for i := 0; i < oldElementCount; i++ {
		es[i] = nil
	}
}

func (this *Vector) Clear() {
	this.removeAllElements()
}

func (this *Vector) FirstElement() *interface{} {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.elementCount == 0 {
		panic("NoSuchElementException")
	}
	return (this.elementData[0]).(*interface{})
}

func (this *Vector) LastElement() *interface{} {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.elementCount == 0 {
		panic("NoSuchElementException")
	}
	return (this.elementData[this.elementCount-1]).(*interface{})
}

func (this *Vector) Get(index int) *interface{} {
	this.mu.Lock()
	defer this.mu.Unlock()

	if index >= this.elementCount || index < 0 {
		panic("ArrayIndexOutOfBoundsException size: " + string(this.elementCount) + " index: " + string(index))
	}
	return (this.elementData[index]).(*interface{})
}

func (this *Vector) Set(index int, element interface{}) *interface{} {
	this.mu.Lock()
	defer this.mu.Unlock()

	if index >= this.elementCount || index < 0 {
		panic("ArrayIndexOutOfBoundsException size: " + string(this.elementCount) + " index: " + string(index))
	}

	oldValue := this.elementData[index]
	this.elementData[index] = element

	return oldValue.(*interface{})
}
