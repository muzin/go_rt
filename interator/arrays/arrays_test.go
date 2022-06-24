package arrays

import (
	"fmt"
	"github.com/muzin/go_rt/collection/hash_map"
	"reflect"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	type args struct {
		arr      []interface{}
		iterator func(item interface{}, index int) interface{}
	}

	mapInterator := func(item interface{}, index int) interface{} {
		s := strconv.Itoa(item.(int)) + "_" + strconv.Itoa(index)
		fmt.Println(s)
		return s
	}

	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		// TODO: Add test cases.
		{name: "1",
			args: args{arr: []interface{}{1, 2, 3}, iterator: mapInterator},
			want: []interface{}{"1_0", "2_1", "3_2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.arr, tt.args.iterator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	type args struct {
		arr        []interface{}
		iterator   func(collection interface{}, item interface{}, index int) interface{}
		collection interface{}
	}

	reduceInterator := func(c interface{}, item interface{}, index int) interface{} {
		s := strconv.Itoa(item.(int)) + "_" + strconv.Itoa(index)
		fmt.Println(s)

		hashMap := c.(*hash_map.HashMap)
		hashMap.Put(item.(int), index)

		return c
	}

	wantHashMap := hash_map.NewHashMap()
	wantHashMap.Put(1, 0)
	wantHashMap.Put(2, 1)
	wantHashMap.Put(3, 2)

	tests := []struct {
		name string
		args args
		want interface{}
	}{
		// TODO: Add test cases.
		{name: "1",
			args: args{arr: []interface{}{1, 2, 3}, iterator: reduceInterator, collection: hash_map.NewHashMap()},
			want: wantHashMap,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reduce(tt.args.arr, tt.args.iterator, tt.args.collection); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}
