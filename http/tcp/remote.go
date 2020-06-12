package tcp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log"
	"path/filepath"
	"strconv"
	"sync"
	"time"
	"tcp_client/providers"
	"tcp_client/lib"
)


const (
	cMaxRetryCount = 20
	cMaxConnetCount = 20
)

var cCurrentNumber = 1

type Remote struct {
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
func ClientRun(ctx context.Context) {
	// 设置初始tcp-client 创建的链接数
	providers.Logger.Info("开始创建tcp的client链接")
	r := Remote{

	}
	cannelTag := 0
	go func() {
		for  {
			select {
			case <-ctx.Done():
				cannelTag = 1
			default:
				
			}
		}
	}()
	
	for cCurrentNumber <= cMaxConnetCount {
		// 直接退出这个类
		if cannelTag == 1 {
			return
		}
		// 判断是否超过重试次数
		retryCount :=0
		providers.Logger.Info("创建服务器连接"+strconv.Itoa(retryCount))
		for retryCount < cMaxRetryCount {
			retryCount ++
			err := r.createConn(cCurrentNumber)
			if err == nil {
				retryCount = 0
				break
			}
			providers.Logger.Errorf("创建第[%d]次连接失败:%s", retryCount, err.Error())
			time.Sleep(1 * time.Second)
		}
		// 第一次链接并且失败了，就panic
		if retryCount >= cMaxRetryCount && cCurrentNumber == 1{
			panic("远程服务器无法建立连接")
		}
		cCurrentNumber++
	}

}


/**
 * 创建连接
 */

func (r *Remote)createConn(index int) error {
	fileCrt := filepath.Join(providers.RootPath, "config") + "/ca.crt";
	providers.Logger.Info("证书路径:"+fileCrt)
	rootPEM := lib.Read3(fileCrt)
	if len(rootPEM) == 0 {
		providers.Logger.Errorf("证书读取有误 %s", rootPEM)
		return errors.New("证书读取有误")
	}
	providers.Logger.Info("证书路径1")
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		panic("failed to parse root certificate")
	}
	providers.Logger.Info("证书路径2")
	conn, err := tls.Dial("tcp", "127.0.0.1:"+ providers.Config.GetString("tcp.port"), &tls.Config{
		RootCAs: roots,
	})
	providers.Logger.Info("证书路径3")
	if err != nil {
		panic("failed to connect: " + err.Error())
	}
	providers.Logger.Info("证书路径4")
	r.Comms[index].Conn = conn
	providers.Logger.Infof("发起链接成功 %s", r.Comms)
	return nil
	//defer conn.Close()
}
