/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: unidirectional_test.go.go
 * @time: 2020/1/3 3:47 下午
 * @project: comm-libs
 */

package number

import (
	`math`
	`sync`
	`testing`
)

func TestUnidirectionalNum_AutoIncrease(t *testing.T) {
	type fields struct {
		lock sync.RWMutex
		num  uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		// TODO: Add test cases.
		{
			name: "TestUnidirectionalNum_AutoIncrease",
			fields: fields{
				lock: sync.RWMutex{},
				num:  math.MaxUint64,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &UnidirectionalNum{
				lock: tt.fields.lock,
				num:  tt.fields.num,
			}
			if got := p.AutoIncrease(); got != tt.want {
				t.Errorf("AutoIncrease() = %v, want %v", got, tt.want)
			}
		})
	}
}
