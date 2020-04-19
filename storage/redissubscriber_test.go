/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: subscriber_test.go.go
 * @time: 2020/2/17 3:30 下午
 * @project: svr-frame
 */

package storage

import (
	`reflect`
	`testing`

	`github.com/generalzgd/comm-libs/conf/svrcfg`
)

func TestNewSubscriber(t *testing.T) {
	type args struct {
		cfg svrcfg.RedisCfg
	}
	tests := []struct {
		name string
		args args
		want *RedisSubscriber
	}{
		// TODO: Add test cases.
		{
			name: "TestNewSubscriber",
			args: args{cfg: svrcfg.RedisCfg{}},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSubscriber(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Logf("NewSubscriber() = %v, want %v", got, tt.want)
			}
		})
	}
}
