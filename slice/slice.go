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

import (
	`math`
)

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
func IsEveryUniqueUint16(in []uint16) bool {
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

func IsEveryUniqueUint(in []uint) bool {
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

func IsEveryUniqueInt(in []int) bool {
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

func MaxUint64(in []uint64) uint64 {
	max := uint64(0)
	for _, it := range in {
		if it > max {
			max = it
		}
	}
	return max
}

func MinUint64(in []uint64) uint64 {
	min := uint64(math.MaxUint64)
	for _, it := range in {
		if it < min {
			min = it
		}
	}
	return min
}

func MaxUint32(in []uint32) uint32 {
	max := uint32(0)
	for _, it := range in {
		if it > max {
			max = it
		}
	}
	return max
}

func MinUint32(in []uint32) uint32 {
	min := uint32(math.MaxInt32)
	for _, it := range in {
		if it < min {
			min = it
		}
	}
	return min
}

func MaxInt(in []int) int {
	max := int(0)
	for _, it := range in {
		if it > max {
			max = it
		}
	}
	return max
}

func MinInt(in []int) int {
	min := int(math.MaxInt64)
	for _, it := range in {
		if it < min {
			min = it
		}
	}
	return min
}

