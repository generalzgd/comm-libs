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

import (
	`errors`
	`fmt`
)



// ////////////////////////////////////////////////////////////////////////////////////////
type IErrCode interface {
	error
	fmt.Stringer
	GetError() error
	SetError(e error) IErrCode
	GetErrMsg() string
	SetErrMsg(s string) IErrCode
	GetCode() int
	SetCode(c int) IErrCode
}

type ErrCode struct {
	err  error
	msg  string
	code int
}

// 实现error接口
func (p *ErrCode) Error() string {
	return p.String()
}

func (p *ErrCode) String() string {
	return fmt.Sprintf("<code:%v msg:%v error:%v>", p.code, p.msg, p.err)
}

func (p *ErrCode) GetError() error {
	return p.err
}

func (p *ErrCode) SetError(e error) IErrCode {
	out := *p
	out.err = e
	return &out
}

func (p *ErrCode) GetErrMsg() string {
	return p.msg
}

// 复制一个新对象
func (p *ErrCode) SetErrMsg(s string) IErrCode {
	out := *p
	out.msg = s
	return &out
}

// 复制一个新对象
func (p *ErrCode) SetCode(c int) IErrCode {
	out := *p
	out.code = c
	return &out
}

func (p *ErrCode) GetCode() int {
	return p.code
}

func MakeErr(errStr, msg string, code int) IErrCode {
	return &ErrCode{err: errors.New(errStr), msg: msg, code: code}
}

func MakeFromErr(err error) IErrCode {
	return &ErrCode{err: err, code: 1, msg: err.Error()}
}
