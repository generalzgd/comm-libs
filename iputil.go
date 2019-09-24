/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: iputil.go
 * @time: 2019/9/24 9:56
 */
package comm_libs

import (
	`fmt`
	`net`
	`strings`
)

/*
* 小端转换成数字
* ipv6会溢出
 */
func Ip2Long(ip string) uint32 {
	ipv := net.ParseIP(strings.TrimSpace(ip))
	out := uint32(0)
	if ipv != nil {
		l := len(ipv)
		for i, v := range ipv {
			out |= uint32(v) << uint((l-1-i)*8)
		}
	}
	return out
}

/*
* 小端转换成字符串
* ipv6会溢出
 */
func Long2Ip(v uint32) string {
	v1 := v & 0xFF000000 >> 24
	v2 := v & 0x00FF0000 >> 16
	v3 := v & 0x0000FF00 >> 8
	v4 := v & 0x000000FF

	res := fmt.Sprintf("%d.%d.%d.%d", v1, v2, v3, v4)
	return res
}