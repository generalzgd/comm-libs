/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: ev.go
 * @time: 2020/4/19 9:27 上午
 * @project: comm-libs
 */

package event

// 异步消息通道队列
type EventChannel chan *EventObj

// 异步消息通道中的数据结构
type EventObj struct {
	Type  int
	Data  interface{}
	Data2 interface{}
}

func NewEventObj(eventId int, data interface{}) *EventObj {
	e := new(EventObj)
	e.Type = eventId
	e.Data = data
	return e
}
