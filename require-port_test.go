package go_http_proxy

import (
	"testing"
)

func Test_required(t *testing.T) {

	type args struct {
		port     int
		protocol string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{name: "1", args: args{port: 0, protocol: "http"}, want: false},
		{name: "2", args: args{port: 65536, protocol: "http"}, want: false},

		{name: "2", args: args{port: 80, protocol: "http"}, want: false},
		{name: "2", args: args{port: 81, protocol: "http"}, want: true},

		{name: "2", args: args{port: 443, protocol: "https"}, want: true},

		{name: "2", args: args{port: 21, protocol: "ftp"}, want: false},
		{name: "2", args: args{port: 23, protocol: "ftp"}, want: true},

		{name: "2", args: args{port: 70, protocol: "gopher"}, want: false},
		{name: "2", args: args{port: 71, protocol: "gopher"}, want: true},

		{name: "2", args: args{port: 70, protocol: "file"}, want: false},
		{name: "2", args: args{port: 71, protocol: "file"}, want: false},

		{name: "2", args: args{port: 71, protocol: "file111"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := required(tt.args.port, tt.args.protocol); got != tt.want {
				t.Errorf("required() = %v, want %v", got, tt.want)
			}
		})
	}
}
