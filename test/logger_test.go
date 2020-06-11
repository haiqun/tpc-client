package test

import (
	"tcp-tls-project/providers"
	"testing"
)

func TestGetLogger(t *testing.T) {
	// 调用不报错
	providers.GetLogger()
	//if providers.Logger == nil {
	//	t.Errorf("got : %v,want :%v ", providers.Logger,"* seelog.LoggerInterface")
	//}
	// 可以进行写入 todo
}