package network

import "io"

type wrappedReader struct {
	r       io.Reader
	ch      Channel
	handler ChannelRawSideHandler
}

func (this wrappedReader) Read(p []byte) (n int, err error) {
	n, err = this.r.Read(p)
	this.handler.OnRawReadOp(this.ch, p[:n], err)
	return n, err
}

func WrapReader(r io.Reader, ch Channel, handler ChannelRawSideHandler) io.Reader {
	return wrappedReader{
		r:       r,
		ch:      ch,
		handler: handler,
	}
}

type wrappedFlusher interface {
	Flush() error
}

type wrappedWriterIf interface {
	io.Writer
	wrappedFlusher
}

type wrappedWriter struct {
	w       io.Writer
	f       wrappedFlusher
	ch      Channel
	handler ChannelRawSideHandler
}

func (this wrappedWriter) Write(p []byte) (n int, err error) {
	this.handler.OnRawWriteOp(this.ch, p)
	return this.w.Write(p)
}

func (this wrappedWriter) Flush() error {
	this.handler.OnRawFlushOp(this.ch)
	return this.f.Flush()
}

func WrapWriterAndFlush(wf wrappedWriterIf, ch Channel, handler ChannelRawSideHandler) wrappedWriterIf {
	return wrappedWriter{
		w:       wf,
		f:       wf,
		ch:      ch,
		handler: handler,
	}
}

func WrapWriter(w io.Writer, ch Channel, handler ChannelRawSideHandler) io.Writer {
	return wrappedWriter{
		w:       w,
		f:       nil,
		ch:      ch,
		handler: handler,
	}
}
