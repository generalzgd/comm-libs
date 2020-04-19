/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: send_test.go
 * @time: 2018/11/23 20:01
 */
package mail

import "testing"

func TestSendMailByUrl(t *testing.T) {
	SendMailByUrl("牛逼", "小灰灰太牛逼了", "zhangguodong@bianfeng.com", "luowenhui@bianfeng.com")
}
