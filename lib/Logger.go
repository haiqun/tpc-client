package lib

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func init(){
	// 设置日志输出的格式
	log.SetFormatter(&log.JSONFormatter{})
	// 设置日输出的地点
	logPath := GetParentDirectory(GetCurrentDirectory())+"/log/"+ time.Now().Format("2006-01-02")+".log"
	fileObj,err := os.OpenFile(logPath,os.O_CREATE|os.O_WRONLY|os.O_APPEND,0644)
	if err != nil {
		fmt.Println(" open file failed err: ",err)
		return
	}
	log.SetOutput(fileObj)
	// 设置日志的级别 - 判断日志级别
	log.SetLevel(log.DebugLevel)
}

func Info(message string)  {
	log.Info(message)
}

func Out(message string)  {
	// 控制台输出
	log.SetOutput(os.Stdout)
	log.Debug(message)
}

func Trace(message string)  {
	log.Trace(message)
}
