package tcpexample

import (
	"fmt"
	"net"
	"time"

	"goimport.moetang.info/nekoq-api/errorutil"
	"goimport.moetang.info/nekoq-api/network"
	"goimport.moetang.info/nekoq-api/network/tcp"
)

func ServerExample() {

}

func ClientExample(tcpConnStr string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", tcpConnStr, timeout)
	if err != nil {
		return err
	}
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return errorutil.New("connection type is not *net.TCPConn")
	}

	defualtHandler := ClientHandler{}
	channel, _ := tcp.NewTcpChannel(tcpConn, defualtHandler, new(tcp.TcpChannelOption))

	// write
	channel.Write(network.NewSimpleWriteEvent([]byte("GET / HTTP/1.1\r\nHost: localhost\r\nConnection: Close\r\n\r\n")))
	// flush
	channel.Flush(network.NewNoopFlushEvent())

	// inbound

	time.Sleep(2 * time.Second)
	return channel.Close()
}

type ClientHandler struct {
	network.DefaultChannelRawSideHandler
}

func (ClientHandler) OnRead(ch network.Channel, data []byte) {
	fmt.Print(string(data))
}

func (ClientHandler) OnRawWriteOp(ch network.Channel, data []byte) {
}
