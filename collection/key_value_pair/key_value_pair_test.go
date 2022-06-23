package key_value_pair

import (
	"fmt"
	"testing"
)

func TestKeyValuePair_ContainsKey(t *testing.T) {
	type args struct {
		key string
	}

	keyValuePair := NewKeyValuePair()

	oneArgs := &args{key: "1"}
	twoArgs := &args{key: "2"}
	threeArgs := &args{key: "3"}

	fourArgs := &[]args{
		args{key: "4"},
	}

	keyValuePair.Put("1", oneArgs)
	keyValuePair.Put("2", twoArgs)
	keyValuePair.Put("3", threeArgs)
	keyValuePair.Put("4", fourArgs)

	newOneArgs := keyValuePair.Get("1")
	newTwoArgs := keyValuePair.Get("2")
	newThreeArgs := keyValuePair.Get("3")
	newFourArgs := keyValuePair.Get("4")

	if oneArgs == newOneArgs {
		fmt.Println("one = ")
	}
	if twoArgs == newTwoArgs {
		fmt.Println("two = ")
	}
	if threeArgs == newThreeArgs {
		fmt.Println("three = ")
	}

	if fourArgs == newFourArgs {
		fmt.Println("fourArgs = ")
	}

	tests := []struct {
		name   string
		fields *KeyValuePair
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &KeyValuePair{
				mu:         tt.fields.mu,
				table:      tt.fields.table,
				threshold:  tt.fields.threshold,
				rebuilding: tt.fields.rebuilding,
				delCount:   tt.fields.delCount,
			}
			if got := this.ContainsKey(tt.args.key); got != tt.want {
				t.Errorf("ContainsKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
