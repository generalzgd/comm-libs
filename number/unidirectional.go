/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: unidirectional.go
 * @time: 2020/1/3 3:38 下午
 * @project: comm-libs
 */

package number

import (
	`math`
	`sync`
)

// 额外的方法处理，保证前后调用的一致性
type AdditionalHandler func(uint64)

// 单向增加数字
type UnidirectionalNum struct {
	lock sync.RWMutex
	num uint64
}

func (p *UnidirectionalNum) GetNumber() uint64 {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.num
}

// 如果比自己小，则忽略。当且仅当maxuint64时，才允许设置小自己的值
func (p *UnidirectionalNum) SetNumber(v uint64, fs ...AdditionalHandler) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.num == math.MaxUint64 || p.num < v {
		p.num = v

		for _, f := range fs {
			f(p.num)
		}
	}
}

// 超过max.Uint64后，从0开始
func (p *UnidirectionalNum) AutoIncrease(fs ...AdditionalHandler) uint64 {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.num += 1
	for _, f := range fs {
		f(p.num)
	}
	return p.num
}
