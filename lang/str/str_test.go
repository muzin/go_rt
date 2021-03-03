package str

import "testing"

func TestIsBlank(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{name: "not_only_blank", args: args{str: " 1 "}, want: false},
		{name: "only_blank", args: args{str: "  "}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBlank(tt.args.str); got != tt.want {
				t.Errorf("IsBlank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotBlank(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{name: "not_only_blank", args: args{str: " 1 "}, want: true},
		{name: "only_blank", args: args{str: "  "}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotBlank(tt.args.str); got != tt.want {
				t.Errorf("IsNotBlank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrim(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{name: "src_value", args: args{str: "1"}, want: "1"},
		{name: "both_blank", args: args{str: " 1 "}, want: "1"},
		{name: "left_blank", args: args{str: " 1"}, want: "1"},
		{name: "right_blank", args: args{str: "1 "}, want: "1"},
		{name: "both_blank_and_center_blank", args: args{str: " 1 1 "}, want: "1 1"},
		{name: "all_blank", args: args{str: "  "}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Trim(tt.args.str); got != tt.want {
				t.Errorf("Trim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{name: "not_only_empty", args: args{str: " 1 "}, want: false},
		{name: "only_empty", args: args{str: ""}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.args.str); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotEmpty(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{name: "not_only_empty", args: args{str: " 1 "}, want: true},
		{name: "only_empty", args: args{str: ""}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotEmpty(tt.args.str); got != tt.want {
				t.Errorf("IsNotEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
