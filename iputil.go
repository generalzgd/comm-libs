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
	`io/ioutil`
	`net`
	`net/http`
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


func GetInnerIp() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0{
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if strings.HasPrefix(iface.Name, "docker") || strings.HasPrefix(iface.Name, "w-") ||
			strings.HasPrefix(iface.Name, "isatap") {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return ""
		}
		// fmt.Println(iface.Name, addrs)
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			// if ip[0] == 10 || ip[0] == 172 || (ip[0]==192 && ip[1] == 168) {
			// 	return ip.String()
			// }
			return ip.String()
		}
	}
	return ""
}

func GetExternIp() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(bs)) // IpStr2Ip(string(bs))
}
