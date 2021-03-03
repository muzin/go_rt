package net

import (
	"testing"
)

func TestIsIP(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{name: "default_ipv4", args: args{s: "192.168.1.111"}, want: 4},
		{name: "default_ipv4", args: args{s: "0.0.0.0"}, want: 4},
		{name: "default_ipv6", args: args{s: "fe80::836:16ca:c5d:6c06%en0"}, want: 6},
		{name: "default_ipv6", args: args{s: "::"}, want: 6},
		{name: "default_without", args: args{s: ""}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIP(tt.args.s); got != tt.want {
				t.Errorf("IsIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsLegalPort(t *testing.T) {
	type args struct {
		port int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{name: "default", args: args{port: 0}, want: true},
		{name: "default", args: args{port: 1}, want: true},
		{name: "default", args: args{port: 65535}, want: true},
		{name: "error_65536", args: args{port: 65536}, want: false},
		{name: "error_-1", args: args{port: -1}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLegalPort(tt.args.port); got != tt.want {
				t.Errorf("IsLegalPort() = %v, want %v", got, tt.want)
			}
		})
	}
}
