/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: staff_test.go
 * @time: 2018/7/20 18:33
 */
package env

import "testing"

func TestRandomInt(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(RandomInt(-5, 5))
	}
}
