package main

import (
	"context"
	"log"
	"tcp_client/http/tcp"
	"tcp_client/providers"
	"time"
)

var CtxCannel context.CancelFunc

func init()  {
	// 初始化变量
	providers.GetRootPath()
	providers.GetLogger()
	providers.GetConfig()
}

func main()  {
	// 链接
	tcpCom  := tcp.Remote{
		Index:20,
		Comms: make(map[int]*tcp.Communication,1),
	}
	var ctx context.Context
	ctx,CtxCannel = context.WithCancel(context.Background())
	tcpCom.ClientRun(ctx,20)
	// 发送发送信息
	comms := tcpCom.Comms[tcpCom.Index]
	conn := comms.Conn
	defer conn.Close()

	for  {
		// 写入信息
		msg := "test"
		conn.Write([]byte(msg))
		n, err := conn.Write([]byte("hello\n"))
		if err != nil {
			log.Println(n, err)
			return
		}

		// 读取信息
		buf := make([]byte, 100)
		n, err = conn.Read(buf)
		if err != nil {
			providers.Logger.Error("读取信息有误:",err)
			time.Sleep(time.Second*2)
			continue
		}
		providers.Logger.Info("读取信息内容:",string(buf[:n]))
	}
	// 退出

}