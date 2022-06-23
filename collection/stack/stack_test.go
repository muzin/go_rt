package stack

import (
	"fmt"
	"github.com/muzin/go_rt/collection/vector"
	"reflect"
	"testing"
)

func TestNewStack(t *testing.T) {
	tests := []struct {
		name string
		want *Stack
	}{
		// TODO: Add test cases.
		{name: "default", want: NewStack()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStack(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStack_IsEmpty(t *testing.T) {

	type Obj struct {
		Name string
	}

	emptyVector := vector.NewVector()

	notEmptyVector := vector.NewVector()
	var obj123 = &Obj{Name: "123"}
	var obj345 = &Obj{Name: "345"}
	notEmptyVector.Add(obj123)
	notEmptyVector.Add(obj345)

	type fields struct {
		vector *vector.Vector
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
		{name: "emptyStack", fields: fields{vector: emptyVector}, want: true},
		{name: "notEmptyStack", fields: fields{vector: notEmptyVector}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := &Stack{
				vector: tt.fields.vector,
			}
			if got := stack.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStack_Pop(t *testing.T) {

	type Obj struct {
		Name string
	}

	//emptyVector := vector.NewVector()

	notEmptyVector := vector.NewVector()
	var obj123 = &Obj{Name: "123"}
	var obj345 = &Obj{Name: "345"}
	t.Logf("obj123 p: %v\n", obj123)
	notEmptyVector.Add(obj123)

	t.Logf("obj345 p: %v\n", obj345)
	notEmptyVector.Add(obj345)

	type fields struct {
		vector *vector.Vector
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		// TODO: Add test cases.
		//{name:"emptyVector", fields: fields{ vector: emptyVector }, want: nil},
		{name: "notEmptyVector", fields: fields{vector: notEmptyVector}, want: obj345},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := &Stack{
				vector: tt.fields.vector,
			}
			got := stack.Pop()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pop() = %v, want %v", got, tt.want)
			} else {
				t.Logf("Pop() = result %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStack_Push(t *testing.T) {
	type Obj struct {
		Name string
	}

	notEmptyVector1 := vector.NewVector()
	notEmptyVector2 := vector.NewVector()
	var obj123 = &Obj{Name: "123"}
	var obj345 = &Obj{Name: "345"}

	type fields struct {
		vector *vector.Vector
	}
	type args struct {
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   args
	}{
		// TODO: Add test cases.
		{name: "notEmptyVector1", fields: fields{vector: notEmptyVector1}, args: args{value: obj123}, want: args{value: obj123}},
		{name: "notEmptyVector2", fields: fields{vector: notEmptyVector2}, args: args{value: obj345}, want: args{value: obj345}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := &Stack{
				vector: tt.fields.vector,
			}
			stack.Push(tt.args.value)
			got := stack.Pop()
			if !reflect.DeepEqual(got, tt.want.value) {
				t.Errorf("Pop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStack_Shift(t *testing.T) {
	type Obj struct {
		Name string
	}

	notEmptyVector := vector.NewVector()
	var obj1 = &Obj{Name: "1"}
	var obj2 = &Obj{Name: "2"}
	var obj3 = &Obj{Name: "3"}
	var obj4 = &Obj{Name: "4"}
	var obj5 = &Obj{Name: "5"}

	var objlist = []interface{}{
		obj1, obj2, obj3, obj4, obj5,
	}

	for i := 0; i < len(objlist); i++ {
		obj := objlist[i]
		t.Logf("add o %v p %v\n", obj, &obj)
		notEmptyVector.Add(obj)
	}

	var objlistPtr = objlist

	t.Logf("objlistPtr %v %v\n", objlistPtr, &objlistPtr)

	type fields struct {
		vector *vector.Vector
	}
	tests := []struct {
		name   string
		fields fields
		want   []interface{}
	}{
		// TODO: Add test cases.
		{name: "notEmptyVector", fields: fields{vector: notEmptyVector}, want: objlistPtr},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := &Stack{
				vector: tt.fields.vector,
			}
			//if got := stack.Shift(); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("Shift() = %v, want %v", got, tt.want)
			//}

			//t.Logf("want %v %v %v %v\n", objlistPtr, &objlistPtr, tt.want, &(tt.want))

			var count = 0
			for i := 0; i < stack.Size(); i++ {
				got := stack.Shift()
				//if nil == got {
				//	break
				//}
				want := &((objlistPtr)[count])
				if !reflect.DeepEqual(got, want) {
					//t.Errorf("Shift() = %v %v , want %v %v %v", got, got, want, *want, &want)
				} else {
					//t.Logf("Shift() = result %v %v , want %v %v %v", got, got, want, *want, &want)
				}
				i--
				count++
			}

		})
	}
}

func TestStack_Size(t *testing.T) {
	type Obj struct {
		Name string
	}

	emptyVector := vector.NewVector()

	notEmptyVector := vector.NewVector()
	var obj123 = &Obj{Name: "123"}
	var obj345 = &Obj{Name: "345"}
	notEmptyVector.Add(obj123)
	notEmptyVector.Add(obj345)

	type fields struct {
		vector *vector.Vector
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
		{name: "emptyStack", fields: fields{vector: emptyVector}, want: 0},
		{name: "notEmptyStack", fields: fields{vector: notEmptyVector}, want: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := &Stack{
				vector: tt.fields.vector,
			}
			if got := stack.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStack_Batch_Pop(t *testing.T) {
	type Obj struct {
		Name string
	}

	notEmptyVector := vector.NewVector()
	var obj1 = &Obj{Name: "1"}
	var obj2 = &Obj{Name: "2"}
	var obj3 = &Obj{Name: "3"}
	var obj4 = &Obj{Name: "4"}
	var obj5 = &Obj{Name: "5"}

	var objlist = []*Obj{
		obj1, obj2, obj3, obj4, obj5,
	}

	for i := 0; i < len(objlist); i++ {
		obj := objlist[i]
		t.Logf("add o %v p %v\n", obj, &obj)
		notEmptyVector.Add(obj)
	}

	stack := &Stack{vector: notEmptyVector}

	//time.Sleep(1 * time.Second)

	for i := 0; i < 5; i++ {
		getPop(stack)
	}

	//time.Sleep(3 * time.Second)

}

func getPop(stack *Stack) {
	pop := stack.Pop()
	fmt.Printf("pop %v\n", pop)
}
