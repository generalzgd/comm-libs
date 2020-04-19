/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: fmt.go
 * @time: 2020/1/20 9:56 下午
 * @project: comm-libs
 */

package fomater

import (
	`fmt`
	`math`
	`strings`
)

// 格式化内存表达式
func FmtMem(num uint64) string {
	unit := []string{"T", "G", "M", "K", "B"}
	out := make([]string, 0, 5)
	for i := 0; i < len(unit); i++ {
		base := uint64(math.Pow(1024, float64(4-i)))
		val := num / base
		num = num - val*base
		if val > 0 {
			out = append(out, fmt.Sprintf("%d%s", val, unit[i]))
		}
	}
	return strings.Join(out, "")
}