/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: checkresponser.go
 * @time: 2017/8/14 17:57
 */
package consul

import (
	"net"
	"net/http"
	"net/url"

	`github.com/astaxie/beego/logs`
)

// http://127.0.0.1:12308/check
func CreateCheckResponser_Http(addr string) error {

	uri, err := url.Parse(addr)
	if err != nil {
		return err
	}

	http.HandleFunc(uri.Path, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("check success"))
	})

	go http.ListenAndServe(uri.Host, nil)
	return nil
}

// 支持tcp响应
func CreateCheckResponser_Tcp(addr string) error {
	uri, err := url.Parse(addr)
	if err != nil {
		return err
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", uri.Host)
	if err != nil {
		return err
	}

	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := listen.AcceptTCP()
			if err != nil {
				logs.Info("conn err: ", err)
				break
			}

			conn.Close()
		}
	}()

	return nil
}
