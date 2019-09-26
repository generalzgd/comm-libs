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
package pickjson

import (
	`encoding/json`
	`reflect`
	"testing"
	`time`
)

func TestStandLibJson(t *testing.T)  {
	parseJson := func(src []byte, field string) interface{} {
		out := map[string]interface{}{}
		json.Unmarshal(src, &out)
		return out[field]
	}
	begin := time.Now()
	for i:=0; i<100000; i++{
		parseJson([]byte(`{"foo":1,"a":"str","b":0.25}`), "b")
	}
	t.Logf("use time:%v", time.Since(begin)) // use time:810.0463ms
}

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
			want: []byte("0.25"), // use time:174.0099ms
		},
		{
			name: "TestPickBytes_2",
			args: args{
				source: []byte(`{"foo":1,"a":"str","b":0.25,"d":{"t":234}}`),
				field:  []byte("d"),
			},
			want: []byte(`{"t":234}`), // use time:334.0191ms
		},
		{
			name: "TestPickBytes_3",
			args: args{
				source: []byte(`{"foo":1,"a":"str","b":0.25,"d":{"t":234, "a":222}}`),
				field:  []byte("a"),
			},
			want: []byte(`"str"`), // use time:168.0096ms
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			begin := time.Now()
			for i:=0;i<100000;i++ {
				if got := PickBytes(tt.args.source, tt.args.field); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("PickBytes() = %v, want %v", got, tt.want)
				}
			}
			t.Logf("use time:%v", time.Since(begin))
		})
	}
}
