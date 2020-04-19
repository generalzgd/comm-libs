/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: code.go
 * @time: 2017/6/28 11:18
 */

package ecode

import "testing"

func TestErrCode_GetCode(t *testing.T) {
	type fields struct {
		errStr string
		msg    string
		code   int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
		{
			name: "t1",
			fields: fields{
				errStr: "test err",
				msg:    "好测一下",
				code:   1,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := MakeErr(tt.fields.errStr, tt.fields.msg, tt.fields.code)
			if got := p.SetCode(2).GetCode(); got != tt.want {
				t.Errorf("ErrCode.GetCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
