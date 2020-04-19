/**
 * @version: 1.0.0
 * @author: zgd: general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: redis.go
 * @time: 2017/9/15 下午1:56
 */

package storage

import (
	"reflect"
	"testing"

	"github.com/generalzgd/comm-libs/conf/svrcfg"
)

func TestRedisCtrl_GetRdsLink(t *testing.T) {

	type args struct {
		rdsName string
	}
	tests := []struct {
		name    string
		args    args
		want    *RedisLink
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "t1",
			args:    args{rdsName: "cache"},
			want:    &RedisLink{},
			wantErr: false,
		},
	}
	p := NewRedisCtrl()
	p.AddCfg(map[string]*svrcfg.RedisCfg{
		"cache": {
			Name:    "cache",
			Address: "192.168.163.200",
			Port:    6379,
		},
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := p.GetRdsLink(tt.args.rdsName)
			got.Close()

			got, err = p.GetRdsLink(tt.args.rdsName)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisCtrl.GetRdsLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedisCtrl.GetRdsLink() = %v, want %v", got, tt.want)
			}
			got.Close()
		})
	}
}
