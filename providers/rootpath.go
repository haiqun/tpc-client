package providers

import (
	"os"
	"path/filepath"
	"sync"
)

//@TODO 因为同一个包里面的 init 调用顺序是看编译器排序的，但变量的声明赋值是在 init 之前的,由此来保障配置文件在最前面
// 根路径
var RootPath string

var OnecR sync.Once

// 初始化函数
func GetRootPath() {
	OnecR.Do(func() {
		RootPath = rootPath()
	})
}

/**
 * 获取根目录地址
 * rootPath 会先读取环境变量的 STOCKASSISTANTSRC
 **/
func rootPath() string {
	var rootPath string
	// 检查root_path环境变量是否存在，因为 go run 操作会新建文件导致配置文件项对不上路径
	if rootPath, exist := os.LookupEnv("SYSTEMCOMMUNICATIONSERVICE"); exist {
		return rootPath
	}
	var err error
	rootPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic("根目录地址解析失败:" + err.Error())
	}
	return rootPath
}

