/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: iputil.go
 * @time: 2019/9/24 9:56
 */
package comm_libs

import (
	"testing"
)

func TestIp2Long(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		// TODO: Add test cases.
		{
			name: "TestIp2Long",
			args: args{ip: "0.0.0.1"},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ip2Long(tt.args.ip); got != tt.want {
				t.Errorf("Ip2Long() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLong2Ip(t *testing.T) {
	type args struct {
		v uint32
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name:"TestLong2Ip",
			args:args{
				v:1,
			},
			want:"0.0.0.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Long2Ip(tt.args.v); got != tt.want {
				t.Errorf("Long2Ip() = %v, want %v", got, tt.want)
			}
		})
	}
}
