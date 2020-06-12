package tcp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log"
	"path/filepath"
	"sync"
	"time"
	"tcp_client/providers"
	"tcp_client/lib"
)


const (
	cMaxRetryCount = 20
)

type Remote struct {
	Index int
	retryCount int8 // 重连次数
	//comms []*communication.Communication
	Comms map[int]*Communication
	mus []sync.Mutex
}

/**
 * 通信
 */
type Communication struct {
	id int // 唯一 id ，客户端可用
	//cas // 添加 id 的原子操作
	Conn *tls.Conn
	logger *log.Logger
	callbackChannel sync.Map//map[uint64]chan string, 保证并发写安全 // 回调通道
	heartbeat time.Time //心跳时间
}

/**
 * 初始化
 */
func (r *Remote) ClientRun(ctx context.Context,index int) {
	for r.retryCount < cMaxRetryCount {
		r.retryCount ++
		err := r.createConn1(index)
		if err == nil {
			r.retryCount = 0
			break
		}
		providers.Logger.Errorf("创建第[%d]次连接失败:%s", r.retryCount, err.Error())
		time.Sleep(1 * time.Second)
	}

	if r.retryCount >= cMaxRetryCount {
		panic("远程服务器无法建立连接")
	}
	providers.Logger.Info("成功创建服务器连接")
}


/**
 * 创建连接
 */

func (r *Remote) createConn1(index int) error {
	fileCrt := filepath.Join(providers.RootPath, "config") + "/ca.crt";
	rootPEM := lib.Read3(fileCrt)
	if len(rootPEM) == 0 {
		providers.Logger.Errorf("证书读取有误 %s", rootPEM)
		return errors.New("证书读取有误")
	}
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		panic("failed to parse root certificate")
	}
	conn, err := tls.Dial("tcp", "127.0.0.1:"+ providers.Config.GetString("tcp.port"), &tls.Config{
		RootCAs: roots,
	})
	if err != nil {
		panic("failed to connect: " + err.Error())
	}
	r.Comms[index].Conn = conn
	return nil
	//defer conn.Close()
}
