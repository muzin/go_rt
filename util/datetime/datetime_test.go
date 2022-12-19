package datetime

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestFormat(t *testing.T) {

	times := "2020-09-18 15:04:05"
	s, _ := time.Parse("2006-01-02 15:04:05", times)

	fmt.Println(s)

	type args struct {
		date   time.Time
		format string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		//{name:"value_nil", args:args{ date: nil, format: "yyyy-MM-dd hh:mm:ss:m" }, want: "" },
		{name: "value_standand", args: args{date: s, format: "yyyy-MM-dd hh:mm:ss"}, want: times},
		{name: "value_standand_exist_ms", args: args{date: s, format: "yyyy-MM-dd hh:mm:ss.ms"}, want: times + ".000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format(tt.args.date, tt.args.format); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			} else {
				t.Logf("Format() = date: %v, result: %v", tt.args.date, got)
			}
		})
	}
}

func TestParse(t *testing.T) {

	times := "2020-09-18 15:04:05"
	s, _ := time.Parse("2006-01-02 15:04:05", times)

	fmt.Println(s)

	type args struct {
		datestr string
		format  string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "default", args: args{datestr: times, format: "yyyy-MM-dd hh:mm:ss"}, want: s, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.datestr, tt.args.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			} else {
				t.Logf("Parse() args = %v, result =  %v, want = %v", tt.args, got, tt.want)
			}
		})
	}
}
