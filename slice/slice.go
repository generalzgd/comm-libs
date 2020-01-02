/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: slice.go
 * @time: 2019/12/27 1:53 下午
 * @project: comm-libs
 */

package slice

// 是否所有元素都相同
func IsAllSameStr(in []string) bool {
	l := len(in)
	for i:= 1;i<l;i++{
		if in[0] != in[i] {
			return false
		}
	}
	return true
}

func IsAllSameInt(in []int) bool {
	l := len(in)
	for i:= 1;i<l;i++{
		if in[0] != in[i] {
			return false
		}
	}
	return true
}

func IsAllSameUint64(in []uint64) bool {
	l := len(in)
	for i:= 1;i<l;i++{
		if in[0] != in[i] {
			return false
		}
	}
	return true
}

// 检测每个元素都是唯一的
func IsEveryUniqueUint16(in []uint32) bool {
	l := len(in)
	for i:=0; i<l; i++ {
		for j:=0; j<l; j++ {
			if j == i {
				continue
			}
			if in[j] == in[i] {
				return false
			}
		}
	}
	return true
}

