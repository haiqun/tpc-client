package lib

import "fmt"

// 定义错误报告
func FailOnErr(err error,msg string)  {
	if err != nil {
		panic(fmt.Sprintf("%s:%s",msg,err))
	}
}