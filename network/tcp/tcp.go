package tcp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"runtime"
	"sync/atomic"

	"import.moetang.info/go/nekoq-api/errorutil"
	"import.moetang.info/go/nekoq-api/network"
)

type wrappedTcpWriter interface {
	io.Writer
	Flush() error
}

type TcpChannelOption struct {
	ReadBufferSize int
	WriteQueueSize int
}

func (this *TcpChannelOption) GetReadBufferSize() int {
	if this.ReadBufferSize <= 0 {
		return 4096
	} else {
		return this.ReadBufferSize
	}
}

func (this *TcpChannelOption) GetWriteQueueSize() int {
	if this.WriteQueueSize <= 0 {
		return 64
	} else {
		return this.WriteQueueSize
	}
}

func NewTcpChannel(tcpConn *net.TCPConn, handler network.ChannelRawSideHandler, option *TcpChannelOption) network.Channel {
	tcpCh := &tcpChannel{
		tcpConn:               tcpConn,
		buf:                   make([]byte, option.GetReadBufferSize()),
		handler:               handler,
		option:                option,
		queue:                 make(chan network.OutboundEvent, option.GetWriteQueueSize()),
		isClose:               0,
		writeCnt:              0,
		closeComplete:         false,
		writeLoopCloseChannel: make(chan bool, 1),
	}
	tcpCh.bufWriter = network.WrapWriterAndFlush(bufio.NewWriter(tcpConn), tcpCh, handler)

	go tcpCh.inboundTask()
	go tcpCh.outboundTask()

	handler.Active(tcpCh)

	return tcpCh
}

type tcpChannel struct {
	buf                   []byte
	tcpConn               *net.TCPConn
	bufWriter             wrappedTcpWriter
	handler               network.ChannelRawSideHandler
	attachment            interface{}
	option                *TcpChannelOption
	queue                 chan network.OutboundEvent
	isClose               int32
	writeCnt              int32
	closeComplete         bool
	writeLoopCloseChannel chan bool
}

func (this *tcpChannel) decWriteCnt() {
	old := atomic.LoadInt32(&this.writeCnt)
	for !atomic.CompareAndSwapInt32(&this.writeCnt, old, old-1) {
		old = atomic.LoadInt32(&this.writeCnt)
	}
}

func (this *tcpChannel) Close() error {
	// set inactive / set nil to queue
	if !atomic.CompareAndSwapInt32(&this.isClose, 0, 1) {
		return network.ErrChannelClosed()
	}

	// wait until no write operation submitted
	wcnt := atomic.LoadInt32(&this.writeCnt)
	if wcnt >= 0 {
		for !atomic.CompareAndSwapInt32(&this.writeCnt, 0, -1) {
			runtime.Gosched()
			wcnt := atomic.LoadInt32(&this.writeCnt)
			if wcnt < 0 {
				break
			}
		}
	}

	this.queue = nil
	this.closeComplete = true
	err := this.tcpConn.Close()
	if err != nil {
		return err
	}
	close(this.writeLoopCloseChannel)
	return err
}

func (this *tcpChannel) Write(ev network.WriteEvent) (submitted bool, err error) {
	return this.fireOutboundEvent(ev)
}

func (this *tcpChannel) Flush(ev network.FlushEvent) (submitted bool, err error) {
	return this.fireOutboundEvent(ev)
}

func (this *tcpChannel) fireOutboundEvent(ev network.OutboundEvent) (bool, error) {
	// add entry count. outbound exit check.
	if this.IsInactive() {
		return false, network.ErrChannelClosed()
	}

	// write cnt for close
	old := atomic.LoadInt32(&this.writeCnt)
	if old < 0 {
		return false, network.ErrChannelClosed()
	}
	for !atomic.CompareAndSwapInt32(&this.writeCnt, old, old+1) {
		old = atomic.LoadInt32(&this.writeCnt)
		if old < 0 {
			return false, network.ErrChannelClosed()
		}
	}

	select {
	case this.queue <- ev:
		this.decWriteCnt()
		return true, nil
	default:
		this.decWriteCnt()
		return false, network.ErrChannelQueueNotReady()
	}
	this.decWriteCnt()
	return false, network.ErrUnknown()
}

func (this *tcpChannel) IsActive() bool {
	return atomic.LoadInt32(&this.isClose) == 0
}

func (this *tcpChannel) IsInactive() bool {
	return atomic.LoadInt32(&this.isClose) == 1
}

func (this *tcpChannel) SetAttachment(attachment interface{}) {
	this.attachment = attachment
}

func (this *tcpChannel) GetAttachment() interface{} {
	return this.attachment
}

func (this *tcpChannel) isReadLoop() bool {
	return true
}

func (this *tcpChannel) isWriteLoop() bool {
	return this.IsActive()
}

// exit condition: read err != nil && length == 9
func (this *tcpChannel) inboundTask() {
	defer func() {
		if e := recover(); e != nil {
			this.handler.OnError(this, network.EVENT_PANIC, errorutil.New(fmt.Sprint(e)))
		}
	}()
	reader := network.WrapReader(this.tcpConn, this, this.handler)
	buf := this.buf
	handler := this.handler
	for this.isReadLoop() {
		n, err := reader.Read(buf)
		if err != nil {
			if n > 0 {
				handler.OnRead(this, buf[:n])
			}
			handler.OnError(this, network.EVENT_READ, err)
			if n == 0 {
				break
			}
			continue
		} else {
			handler.OnRead(this, buf[:n])
			continue
		}
	}
}

// exit condition: close
func (this *tcpChannel) outboundTask() {
	defer func() {
		if e := recover(); e != nil {
			this.handler.OnError(this, network.EVENT_PANIC, errorutil.New(fmt.Sprint(e)))
		}
	}()
	queue := this.queue
	w := this.bufWriter
	handler := this.handler
	wlclose := this.writeLoopCloseChannel
MAIN_WRITE_LOOP:
	for this.isWriteLoop() {
		var ev network.OutboundEvent
		select {
		case ev = <-queue:
		case <-wlclose:
			break MAIN_WRITE_LOOP
		}
		switch ev.GetType() {
		case network.EVENT_WRITE:
			err := handler.OnWrite(this, ev.(network.WriteEvent))
			if err != nil {
				handler.OnError(this, network.EVENT_WRITE, err)
			}
			_, err = ev.Process(w)
			if err != nil {
				handler.OnError(this, network.EVENT_WRITE, err)
			}
		case network.EVENT_FLUSH:
			err := w.Flush()
			if err != nil {
				handler.OnError(this, network.EVENT_FLUSH, err)
			}
		}
	LOOP:
		for this.isWriteLoop() {
			select {
			case ev = <-queue:
			default:
				break LOOP
			}
			switch ev.GetType() {
			case network.EVENT_WRITE:
				err := handler.OnWrite(this, ev.(network.WriteEvent))
				if err != nil {
					handler.OnError(this, network.EVENT_WRITE, err)
				}
				_, err = ev.Process(w)
				if err != nil {
					handler.OnError(this, network.EVENT_WRITE, err)
				}
			case network.EVENT_FLUSH:
				err := w.Flush()
				if err != nil {
					handler.OnError(this, network.EVENT_FLUSH, err)
				}
			}
		}
		if this.isWriteLoop() {
			err := w.Flush()
			if err != nil {
				handler.OnError(this, network.EVENT_FLUSH, err)
			}
		}
	}
	// clean submitted task
	for {
		if this.closeComplete {
		CLEAN_LOOP:
			for {
				select {
				case e := <-queue:
					switch e.GetType() {
					case network.EVENT_FLUSH:
						fe := e.(network.FlushEvent)
						f := fe.Callback
						if f != nil {
							f(fe, network.ErrChannelClosed())
						}
					case network.EVENT_WRITE:
						we := e.(network.WriteEvent)
						f := we.Callback
						if f != nil {
							f(we, network.ErrChannelClosed())
						}
					}
				default:
					break CLEAN_LOOP
				}
			}
			break
		}
		runtime.Gosched()
	}
}

//TODO support: timer, stream, multiplexing
