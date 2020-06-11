package tcp

import "crypto/tls"

type TlsPoolConnInfo struct {
	maxConn int
	conn map[int]*tls.Conn// 连接池设置
}