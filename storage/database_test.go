/**
 * @version: 1.0.0
 * @author: zgd: general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: database.go
 * @time: 2017/9/15 上午9:48
 */

package storage

import (
	"reflect"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"

	"github.com/generalzgd/comm-libs/conf/svrcfg"
)

func TestDbCtrl_GetDbEngin(t *testing.T) {

	type args struct {
		dbname string
	}
	tests := []struct {
		name    string
		args    args
		want    *xorm.Engine
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "t1",
			args:    args{dbname: "user"},
			want:    &xorm.Engine{},
			wantErr: false,
		},
	}
	p := NewDbCtrl()
	p.AddCfg(map[string]*svrcfg.DatabaseCfg{
		"user": {
			Host:     "192.168.163.184",
			Port:     3306,
			Name:     "user",
			Username: "live",
			Password: "admin",
		},
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := p.GetDbEngin(tt.args.dbname)
			for i := 0; i < 11; i += 1 {
				time.Sleep(time.Second * 7)
				got, err = p.GetDbEngin(tt.args.dbname)
				if (err != nil) != tt.wantErr {
					t.Errorf("DbCtrl.GetDbEngin() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DbCtrl.GetDbEngin() = %v, want %v", got, tt.want)
			}
		})
	}
}
