package controllers

import (
	"github.com/gin-gonic/gin"
	"tcp-tls-project/http/tcp"
	"tcp-tls-project/providers"
	"time"
)

type Test struct {
	Controller
}

func (test *Test)TestA(c *gin.Context)  {
	providers.Logger.Info("test12312312")
	t := tcp.GetTcpCommunication(20)
	go t.ServerRun()
	tcpConn,ids := t.GetTcpConn()
	time.Sleep(time.Second * 10 )
	if tcpConn == nil {
		providers.Logger.Errorf("GetTcpConn err")
		return
	}
	providers.Logger.Info("GetTcpConn success")
	//发送
	tcpConn.TcpAisle[ids].Send <- "test"
	//接受
	mgs := <-tcpConn.TcpAisle[ids].Accept
	providers.Logger.Infof("接受到tcp-client的值:%s",mgs)
}