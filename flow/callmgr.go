/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: callmgr.go
 * @time: 2018/8/29 12:28
 */
package flow

import (
	"runtime/debug"
	"sync"
	"time"

	`github.com/astaxie/beego/logs`
)

/*
* todo 异步流水线处理
 */
type WorkFlow struct {
	sync.Mutex
	pippool map[uint32]chan interface{}
}

var (
	ansyCallOnce sync.Once
	ansyCallInst *WorkFlow
)

func GetAnsyFlowInst() *WorkFlow {
	if ansyCallInst == nil {
		ansyCallOnce.Do(func() {
			ansyCallInst = &WorkFlow{
				pippool: make(map[uint32]chan interface{}),
			}
		})
	}
	return ansyCallInst
}

func NewAnsyFlowInst() *WorkFlow {
	inst := &WorkFlow{
		pippool: make(map[uint32]chan interface{}),
	}
	return inst
}

func (p *WorkFlow) newPip(id uint32) chan interface{} {
	p.Lock()
	defer p.Unlock()

	pip, ok := p.pippool[id]
	if !ok {
		pip = make(chan interface{})
		p.pippool[id] = pip
		return pip
	}
	return pip
}

func (p *WorkFlow) getPip(id uint32) (chan interface{}, bool) {
	p.Lock()
	defer p.Unlock()

	pip, ok := p.pippool[id]
	return pip, ok
}

func (p *WorkFlow) delPip(id uint32) {
	p.Lock()
	defer p.Unlock()

	pip, ok := p.pippool[id]
	if ok {
		delete(p.pippool, id)
		close(pip)
	}
}

/*
* todo 异步执行, 在等待消息回来前，会暂停当前gorotine
* 注意：此流水线只适用于一来一回。如果一次请求，有多个响应，请选择其他方案
* @param id 唯一标识id
* @param timeout 设置超时的时长
* @func() 执行前需要处理的方法
* @func(interface) 异步消息回来之后，需要执行的内容同时传入返回的数据 (判空和类型断言一定要先处理),
* @func() 超时需要处理的方法
 */
func (p *WorkFlow) Call(id uint32, timeout time.Duration, prevFun func() error, doneFun func(interface{}), timeoutFun func()) {
	pip := p.newPip(id)

	if err := prevFun(); err != nil {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			logs.Error("WorkFlow.Call panic, askid:", id, "err:", err, "stack:", string(debug.Stack()))
		}
		p.delPip(id)
	}()

	if timeout <= 0 {
		timeout = time.Millisecond * 500
	}

	select {
	case data, ok := <-pip:
		if ok {
			doneFun(data)
		}
	case <-time.After(timeout):
		timeoutFun()
	}
}

/*
* todo 异步执行的消息返回
* if channel is closed, it will raise panic, then by recover to catch the panic
* @param id 唯一标识id
* @param data 返回需要处理的数据
* @return 是否响应数据成功，对应pip不存在则失败
 */
func (p *WorkFlow) OnAnsyBack(id uint32, data interface{}) (succ bool) {
	defer func() {
		if err := recover(); err != nil {
			succ = false
			logs.Error("WorkFlow.OnAnsyBack panic, askid:", id, "err:", err, "stack:", string(debug.Stack()))
		}
	}()

	if pip, ok := p.getPip(id); ok {
		pip <- data
		return true
	}
	return false
}
