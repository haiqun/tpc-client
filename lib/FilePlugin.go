package lib

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

// 获取父级目录
func GetParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

// 获取当前目录
func GetCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// 读取证书文件，返回字符串
func Read3(fileName string)  (string){
	f, err := os.Open(fileName)
	if err != nil {
		Info("read file fail")
		return ""
	}
	defer f.Close()
	fd, err := ioutil.ReadAll(f)
	if err != nil {
		Info("read to fd fail")
		return ""
	}
	return string(fd)
}
