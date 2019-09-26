/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: pick.go
 * @time: 2019/9/25 18:21
 */
package json

import (
	"reflect"
	"testing"
)

func TestPickBytes(t *testing.T) {
	type args struct {
		source []byte
		field  []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{
			name: "TestPickBytes_1",
			args: args{
				source: []byte(`{"foo":1,"a":"str","b":0.25}`),
				field:  []byte("b"),
			},
			want: []byte("0.25"),
		},
		{
			name: "TestPickBytes_2",
			args: args{
				source: []byte(`{"foo":1,"a":"str","b":0.25,"d":{"t":234}}`),
				field:  []byte("d"),
			},
			want: []byte(`{"t":234}`),
		},
		{
			name: "TestPickBytes_3",
			args: args{
				source: []byte(`{"foo":1,"a":"str","b":0.25,"d":{"t":234, "a":222}}`),
				field:  []byte("a"),
			},
			want: []byte(`222`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PickBytes(tt.args.source, tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PickBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
