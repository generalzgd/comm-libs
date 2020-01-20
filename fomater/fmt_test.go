/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: fmt_test.go.go
 * @time: 2020/1/20 9:57 下午
 * @project: comm-libs
 */

package fomater

import (
	`math`
	`testing`
)

func Test_FmtMem(t *testing.T) {
	type args struct {
		num uint64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name:"Test_fmtMem_1",
			args:args{num:1},
			want:"1B",
		},
		{
			name:"Test_fmtMem_2",
			args:args{num:1024},
			want:"1K",
		},
		{
			name:"Test_fmtMem_3",
			args:args{num:uint64(math.Pow(1024, 1))+1},
			want:"1K1B",
		},
		{
			name:"Test_fmtMem_4",
			args:args{num:uint64(math.Pow(1024, 2))+1024+1},
			want:"1M1K1B",
		},
		{
			name:"Test_fmtMem_5",
			args:args{num:uint64(math.Pow(1024, 3))+uint64(math.Pow(1024,2))+1024+1},
			want:"1G1M1K1B",
		},
		{
			name:"Test_fmtMem_6",
			args:args{num:uint64(math.Pow(1024, 4))},
			want:"1T",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FmtMem(tt.args.num); got != tt.want {
				t.Errorf("fmtMem() = %v, want %v", got, tt.want)
			}
		})
	}
}