package test

import (
	"tcp-tls-project/providers"
	"testing"
)

func TestGetRootPath(t *testing.T) {
	providers.GetRootPath() // todo 返回值是执行环境的对应的值
	//want := "/var/folders/d3"
	////  测试失败输出错误提示
	//index := strings.Index(got,want)
	//if index == -1{
	//	t.Errorf("got : %s,want :%s ,index: %v", got, want,index)
	//}
	//
}