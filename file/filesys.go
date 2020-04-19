/**
 * @version: 1.0.0
 * @author: Administrator:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: tools.go
 * @time: 2017/4/10 14:11
 */

package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	LogsDir  = "logs"
	CacheDir = "cache"
	// CrashDir = "crash"
	// RecoverDir = "recover"
)

// 检测文件夹
func DetectDir(dirArgs ...string) {
	path := filepath.Dir(os.Args[0])
	for _, dirName := range dirArgs {
		dirPath := path + "/" + dirName
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			os.MkdirAll(dirPath, os.ModePerm)
			fmt.Println("Dir created: ", dirPath)
		}
	}
}

func GetCacheDir() string {
	return filepath.Dir(os.Args[0]) + "/" + CacheDir
}

func GetCacheFile(fileName string) string {
	return filepath.Dir(os.Args[0]) + "/" + CacheDir + "/" + fileName
}

func GetSvrAppDir() string {
	t := filepath.Base(filepath.Dir(os.Args[0]))
	return t
}

/*
* todo 判断文件是否存在，并获取文件信息
 */
func FileExists(path string) (bool, os.FileInfo) {
	info, err := os.Stat(path)
	if err == nil {
		return true, info
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, info
}

func SaveFile(content []byte, path string, keepOld bool) error {
	// 建目录
	dir := filepath.Dir(path)
	if ok, _ := FileExists(dir); !ok {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	if keepOld {
		if ok, _ := FileExists(path); ok {
			if err := os.Rename(path, path+"."+time.Now().Format("200601021504")); err != nil {
				return err
			}
		}
	}

	return ioutil.WriteFile(path, content, os.ModePerm)
}

func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func RemoveFile(path string) error {
	if ok, _ := FileExists(path); ok {
		return os.Remove(path)
	}
	return nil
}
