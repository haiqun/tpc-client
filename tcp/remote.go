package tcp

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"tcp-tls-project/lib"
	"time"
)

func HandleClientConnect(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Receive Connect Request From ", conn.RemoteAddr().String())
	buffer := make([]byte, 1024)
	for {
		len, err := conn.Read(buffer)
		if err != nil {
			log.Println(err.Error())
			break
		}
		fmt.Printf("Receive Data: %s\n", string(buffer[:len]))
		fmt.Println(conn.RemoteAddr().String() )
		//发送给客户端
		_, err = conn.Write([]byte("服务器收到数据:" + string(buffer[:len]) + "  " + conn.RemoteAddr().String()))
		if err != nil {
			break
		}
	}
	log.Println("Client " + conn.RemoteAddr().String() + " Connection Closed.....")
}

func ServerRun()  {
	path := lib.GetParentDirectory(lib.GetCurrentDirectory())+"/config/";
	crt, err := tls.LoadX509KeyPair(path+"ca.crt", path+"ca.key")
	if err != nil {
		log.Fatalln(err.Error())
	}
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = []tls.Certificate{crt}
	tlsConfig.Time = time.Now
	tlsConfig.Rand = rand.Reader
	l, err := tls.Listen("tcp", "fhq.com:8999", tlsConfig)
	if err != nil {
		log.Fatalln(err.Error())
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		} else {
			go HandleClientConnect(conn)
		}
	}
}

func ClientStart() * TlsPoolConnInfo {
	var m TlsPoolConnInfo
	m = TlsPoolConnInfo{
		maxConn:20,
		conn:make(map[int]* tls.Conn,20),
	}
	m.createConnClient()
	return &m
}

func (m * TlsPoolConnInfo) GetConnClient() *tls.Conn  {
	n := lib.RandInt(1,m.maxConn)
	log.Println(n)
	if len(m.conn) == 0 {
		return nil
	}
	conn, ok := m.conn[n];
	if ok {
		return conn
	}else{
		return  nil
	}
}


func (m * TlsPoolConnInfo) createConnClient()  {
	fileCrt := lib.GetParentDirectory(lib.GetCurrentDirectory())+"/config/ca.crt";
	rootPEM := lib.Read3(fileCrt)
	if len(rootPEM) == 0 {
		log.Fatalf("证书读取有误 %s",rootPEM)
	}
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		panic("failed to parse root certificate")
	}
	for j:=0;j<=m.maxConn;j++ {
		conn, err := tls.Dial("tcp", "127.0.0.1:17080", &tls.Config{
			RootCAs: roots,
		})
		if err != nil {
			panic("failed to connect: " + err.Error())
		}
		m.conn[j] = conn
	}
}
