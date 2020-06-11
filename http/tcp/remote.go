package tcp

import (
	"crypto/rand"
	"crypto/tls"
	"net"
	"path/filepath"
	"tcp-tls-project/lib"
	"tcp-tls-project/providers"
	"time"
)

/**
 * 启动tcp - 监听服务
 */
func (t *TcpCommunication)ServerRun()  {
	// 将连接保存起来
	path := filepath.Join(providers.RootPath, "config")
	crt, err := tls.LoadX509KeyPair(path+"/ca.crt", path+"/ca.key")
	if err != nil {
		panic("证书认证有误！err ：%s" + err.Error())
	}
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = []tls.Certificate{crt}
	tlsConfig.Time = time.Now
	tlsConfig.Rand = rand.Reader
	l, err := tls.Listen("tcp", "127.0.0.1:"+ providers.Config.GetString("tcp.port"), tlsConfig)

	if err != nil {
		panic("监听失败,原因:" + err.Error())
	}
	providers.Logger.Info("服务启动成功!开始监听:" + providers.Config.GetString("tcp.port"))
	// 心跳监测 todo

	for {
		conn, err := l.Accept()// 这个是连接过来的tcp连接
		if err != nil {
			providers.Logger.Errorf("接收失败，失败原因:%s", err.Error())
			continue
		}
		// ip 监测 todo
		providers.Logger.Infof("有新链接进入:%+v", conn.RemoteAddr().String())
		go t.createCommunication(conn) // 将连接保存起来
	}
	providers.Logger.Info("监听关闭")
	// 监听关闭
	l.Close()
}

// 同一个tcp内部的，通讯携程的概念
type TcpAisle struct {
	Send chan string
	Accept chan string
}

// 一个tcp连接的实例
type TcpConnect struct {
	Conn net.Conn // 当前通讯连接
	Occupancy int // 当前tcp被多个业务端使用
	TcpAisle map[string]TcpAisle
}

// 一个tpc连接池
type TcpCommunication struct {
	ConnPool map[int]TcpConnect
}

func (tc *TcpConnect)TcpSendMsg(msg string) bool {
	//发送给客户端
	_, err := tc.Conn.Write([]byte(msg))
	if err != nil {
		providers.Logger.Errorf("TcpSendMsg failed: %s",err.Error())
		return false
	}
	return true
}


//var TC TcpCommunication // tcp的连接池
var currentConn = 1  // 初始化 - 当前的连接数

func GetTcpCommunication(i int) *TcpCommunication {
	tc := TcpCommunication{
		ConnPool: make(map[int]TcpConnect,i),
	}
	return &tc
}

/**
 * 保存tcp的连接
 */
func (t *TcpCommunication)createCommunication(conn net.Conn) {
	providers.Logger.Infof("createCommunication:%+v", conn.RemoteAddr().String())
	tcpC := TcpConnect{
		Conn: conn,
		TcpAisle: make(map[string]TcpAisle,100), // todo 默认每个tcp能给100个服务端调用
	}
	t.ConnPool[currentConn] = tcpC
	currentConn ++
}

func (t *TcpCommunication)GetTcpConn() (tcp *TcpConnect,ids string){
	if currentConn == 1 {
		return nil,""
	}
	n := lib.GetRandNum(currentConn)
	providers.Logger.Infof("GetTcpConn 获取链接n:%d",n)
	tcpC ,ok := t.ConnPool[n]
	if !ok {
		return nil,""
	}

	ids = string(lib.Krand(10,1000))
	tcpC.TcpAisle[ids] = TcpAisle{
		Send: make(chan string,100),
		Accept: make(chan string,100),
	}
	// 监听被调用的tcp通道
	go sendChannel(&tcpC,ids)
	go acceptanceChannel(&tcpC,ids)
	return &tcpC,ids
}


// 监听发送给tcp-client的channel，有信息马上发送
func sendChannel(tcp *TcpConnect,ids string) {
	conn := tcp.Conn
	//buffer := make([]byte, 1024)
	for {
		select {
			case sendData := <-tcp.TcpAisle[ids].Send:
				//发送给客户端
				_, err := conn.Write([]byte(sendData))
				if err != nil {
					providers.Logger.Info("Tcp sendData : %s ",err.Error())
				}
		}
	}
}

// 监听tcp-client发送过来的信息，马上发送响应给 tcp.TcpAisle[ids].Accept
func acceptanceChannel(tcp *TcpConnect,ids string)  {
	buffer := make([]byte, 1024)
	conn := tcp.Conn
	for {
		len, err := conn.Read(buffer)
		if err != nil {
			providers.Logger.Info("acceptanceChannel : %s ",err.Error())
			break
		}
		tcp.TcpAisle[ids].Accept <- string(buffer[:len])
	}
}


