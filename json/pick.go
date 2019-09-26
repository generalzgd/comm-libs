/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: pick.go
 * @time: 2019/9/25 18:21
 */
package json

var (
	startSymble = map[byte]struct{}{
		[]byte("{")[0]: {},
		[]byte("[")[0]: {},
	}

	endSymble = map[byte]struct{}{
		[]byte("}")[0]: {},
		[]byte("]")[0]: {},
	}

	colon        = []byte(":")[0]  // 58
	comma        = []byte(",")[0]  // 44
	blank        = []byte(" ")[0]  // 32
	quot         = []byte("\"")[0] // 34
	slash        = []byte("\\")[0] //
	beginBrace   = []byte("{")[0]  //
	endBrace     = []byte("}")[0]  //
	beginBracket = []byte("[")[0]  //
	endBracket   = []byte("]")[0]
	pair         = map[byte]byte{
		beginBrace:   endBrace,
		beginBracket: endBracket,
	}
)

func isStartSymble(c byte) bool {
	_, ok := startSymble[c]
	return ok
}

func isEndSymble(c byte) bool {
	_, ok := endSymble[c]
	return ok
}

func isPair(f, l byte) bool {
	if v, ok := pair[f]; ok && v == l {
		return true
	}
	return false
}

func strEqual(s, d []byte) (bool, int) {
	if len(s) != len(d) {
		return false, 0
	}
	for i := 0; i < len(s); i++ {
		if s[i] != d[i] {
			return false, i
		}
	}
	return true, len(s)
}

// 获取json串中，第一个field字段对应的value
// 不支持同名field，多值获取
func PickBytes(source []byte, field []byte) []byte {
	fieldLen := len(field)
	end := len(source)
	for i := 0; i < end-fieldLen; i++ {

		eq, n := strEqual(source[i:i+fieldLen], field)
		i += n
		if !eq {
			continue
		}

		findColon := false
		for k := 0; k < end; k++ {
			char := source[i+k]
			if char == blank || char == quot {
				continue
			} else if char == colon {
				i += k
				findColon = true
				break
			}
		}
		if findColon {
			val := findValue(source[i:])
			// str := string(val)
			// fmt.Println("output: ", str)
			return val
		}
	}
	return nil
}

func nextQuot(source []byte) int {
	begin := 0
	for i, char := range source {
		if i == 0 && char == quot {
			begin = 1
			continue
		}
		if char == quot {
			if i-1 > begin && source[i-1] == slash { // 引号作为内容
				continue
			}
			return i
		}
	}
	return 0
}

func findValue(source []byte) []byte {
	// fmt.Println(string(source))
	startIdx := 0
	for i := 0; i < len(source); i++ {
		char := source[i]
		// str := string([]byte{char})
		// fmt.Printf(str)
		switch char {
		case blank:
			continue
		case colon:
			startIdx = i + 1 // 标记开始
			continue
		case quot: // 引号
			id := nextQuot(source[i:])
			i += id
		case comma, endBrace:
			// 没有引号的时候，获取到逗号
			if startIdx < i {
				return source[startIdx:i]
			}
		default:
			// 没有引号的时候，获取到对象符号
			if isStartSymble(char) {
				return findObjectValue(source[i:])
			}
		}
	}
	return nil
}

func findObjectValue(source []byte) []byte {
	if len(source) <0 || !isStartSymble(source[0]) {
		return nil
	}
	stack := make([]byte,0,10)
	for i:=0;i<len(source);i++ {
		char := source[i]
		if isStartSymble(char) {
			stack = append(stack, char) // 压栈
			continue
		}
		if char == quot {
			id := nextQuot(source[i:]) // 调过引号里的括号
			i += id
			continue
		}
		if isEndSymble(char) {
			if len(stack) < 1{
				return nil // 异常
			}
			last := stack[len(stack)-1]
			if isPair(last, char) {
				stack = stack[:len(stack)-1] // 退栈
			}
			if len(stack) == 0 {
				return source[ : i+1]
			}
		}
	}
	return nil
}