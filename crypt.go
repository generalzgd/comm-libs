/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: crypt.go
 * @time: 2020/1/9 5:20 下午
 * @project: comm-libs
 */

package comm_libs

import (
	`crypto/hmac`
	`crypto/md5`
	`crypto/sha1`
	`fmt`
)

// ////////////////////////////////////////////////////
func Md5(in string) string {
	h := md5.New()
	h.Write([]byte(in))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 加密形式在php中是hash_hmac('sha1',$string,$key)
func HmacSha1(in, key []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(in)
	return mac.Sum(nil)
	//return fmt.Sprintf("%x", mac.Sum(nil))
}
