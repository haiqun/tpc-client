package lib

import (
	"io/ioutil"
	"os"
	"tcp_client/providers"
)


// 读取证书文件，返回字符串
func Read3(fileName string)  (string){
	f, err := os.Open(fileName)
	if err != nil {
		providers.Logger.Info("read file fail")
		return ""
	}
	defer f.Close()
	fd, err := ioutil.ReadAll(f)
	if err != nil {
		providers.Logger.Info("read to fd fail")
		return ""
	}
	return string(fd)

}
