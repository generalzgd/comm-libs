/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: send.go
 * @time: 2018/11/23 19:50
 */
package mail

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	url = "http://mail-service.zhanqi.tv/send.php"
)

/*
* todo send mail by curl.
 */
func SendMailByUrl(title, body string, toAddrs ...string) error {
	var toList []string
	for i, v := range toAddrs {
		toList = append(toList, fmt.Sprintf("to[%d][]=%s", i, v))
	}

	params := fmt.Sprintf("%s&title=%s&body=%s", strings.Join(toList, "&"), title, body)
	payload := strings.NewReader(params)
	// payload := strings.NewReader("to[0][]=zhangguodong@bianfeng.com&title=haha&body=test888888888888888888&cc=luowenhui@bianfeng.com")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	bakBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// fmt.Println(res)
	fmt.Println(string(bakBody))
	return nil
}
