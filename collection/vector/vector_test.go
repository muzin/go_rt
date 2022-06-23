package vector

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewInstanceOfVector(t *testing.T) {
	type args2 struct {
		key string
	}

	vector := NewVector()

	a := &args2{key: "a"}
	b := &args2{key: "b"}
	c := &args2{key: "c"}

	vector.Add(a)
	vector.Add(b)
	vector.Add(c)

	newa := vector.Get(0)
	newb := vector.Get(1)
	newc := vector.Get(2)

	if a == newa {
		fmt.Println("one = ")
	}
	if b == newb {
		fmt.Println("two = ")
	}
	if c == newc {
		fmt.Println("three = ")
	}

	type args struct {
		initialCapacity   int
		capacityIncrement int
	}
	tests := []struct {
		name string
		args args
		want *Vector
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInstanceOfVector(tt.args.initialCapacity, tt.args.capacityIncrement); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInstanceOfVector() = %v, want %v", got, tt.want)
			}
		})
	}
}
