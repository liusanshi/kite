package util

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const isDev = false

//getCurrentPath 获取当前运行的路径
func getCurrentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//GetCurrentPath 获取当前运行的路径
func GetCurrentPath() string {
	if isDev {
		return "E:\\git\\learn\\GO\\kite"
	}
	return getCurrentPath()
}

//FileExists 判断文件路径是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

//IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

//Splite 根据相对路径切分绝对路径
func Splite(fullPath, relativePath string) string {
	relative, _ := os.Stat(relativePath)
	full, err := os.Stat(fullPath)
	rela := []string{}
	path := fullPath
	for err == nil && !os.SameFile(relative, full) {
		rela = append(rela, full.Name())
		path = filepath.Dir(path)
		full, err = os.Stat(path)
	}
	//逆序
	count := len(rela) - 1
	reso := make([]string, count+1)
	for i := count; i >= 0; i-- {
		reso[count-i] = rela[i]
	}
	return filepath.Join(reso...)
}

//Md5 获取文件的Md5
func Md5(filePath string) string {
	if !FileExists(filePath) {
		panic("file not exists")
	}
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	h := md5.New()
	_, err = io.Copy(h, file)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

//IndexOf 查询字符串在字符串数组的位置
func IndexOf(list []string, sub string) int {
	for i, s := range list {
		if s == sub {
			return i
		}
	}
	return -1
}
