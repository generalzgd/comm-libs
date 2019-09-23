/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: strutil.go
 * @time: 2019/9/23 14:23
 */
package comm_libs

import (
	`strings`
)

// 首字母大写，不含unicode
func UpCaseString(str string) string {
	if len(str) < 1 {
		return str
	}

	first := string(str[0])
	tail := str[1:]
	return strings.ToUpper(first) + tail
}

func LowCaseString(str string) string {
	if len(str) < 1 {
		return str
	}
	first := string(str[0])
	tail := str[1:]
	return strings.ToLower(first) + tail
}