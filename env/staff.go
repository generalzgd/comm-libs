/**
 * @version: 1.0.0
 * @author: zgd: general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: staff.go
 * @time: 2017/4/28 下午3:13
 */

package env

import (
	mathRand "math/rand"
	"net"
	"os"
	"strings"
)

var (
	// 随机数种子
	// r = mathRand.New(mathRand.NewSource(time.Now().UnixNano()))

	baseStr = "abcdefghijklmnopqrstuvwxyzZBCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	baseLen = len(baseStr)

	env = ""
)

func init() {
	// mathRand.Seed(time.Now().UnixNano())
}

func RandomInt(min int, max int) int {
	if min == max {
		return min
	} else if min > max {
		min, max = max, min
	}
	return mathRand.Intn(max-min) + min
}

func RandomString(length int) []byte {

	str := make([]byte, 0, length)
	for i := 0; i < length; i++ {
		idx := mathRand.Intn(baseLen - 1)
		str = append(str, byte(baseStr[idx]))
	}
	return str
}

const (
	ENV_DEV    = "dev"
	ENV_BETA   = "beta"
	ENV_ONLINE = "online"
)

// 根据主机名获取对应的环境变量字符串，dev/beta/online
func GetEnvName() string {
	if len(env) > 0 {
		return env
	}

	e := os.Getenv("GOENV")
	if len(e) > 0 {
		env = e
		return env
	}

	hostName, _ := os.Hostname()
	if hostName == `BF-ZhanqiTest` {
		env = ENV_DEV
		return env
	}
	if hostName == `BF-ZhanqiBeta` {
		env = ENV_BETA
		return env
	}
	if strings.Index(hostName, `local`) > -1 {
		env = ENV_DEV
		return env
	}
	if strings.Index(hostName, `BF-`) > -1 {
		env = ENV_DEV
		return env
	}
	env = ENV_ONLINE
	return env
}

// 启动容器时，获取一些必要的信息来测试
func GetContainerInfo() (info string, host string, list []string) {
	var err error
	info = "container info:"
	host, _ = os.Hostname()
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, address := range addrList {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				list = append(list, ipNet.IP.String())
			}
		}
	}

	return
}

// 首字母大写，不含unicode
/*func UpCaseString(str string) string {
	if len(str) < 1 {
		return str
	}

	first := string(str[0])
	tail := str[1:]
	return strings.ToUpper(first) + tail
}*/

/*func LowCaseString(str string) string {
	if len(str) < 1 {
		return str
	}
	first := string(str[0])
	tail := str[1:]
	return strings.ToLower(first) + tail
}*/
