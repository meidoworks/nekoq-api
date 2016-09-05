package network

import (
	"io"
	"log"
)

type EVENT_TYPE int

const (
	EVENT_READ EVENT_TYPE = iota
	EVENT_WRITE
	EVENT_FLUSH
	EVENT_PANIC
	EVENT_USER
)

type ChannelRawSideHandler interface {
	Active(Channel)
	Inactive(Channel)
	OnRead(Channel, []byte)
	OnError(channel Channel, eventType EVENT_TYPE, err error)
	OnUserEvent(channel Channel, userEvent *UserEvent)

	OnWrite(Channel, WriteEvent) error
	OnFlush(Channel, FlushEvent) error

	// callback when following operations triggered
	OnRawReadOp(Channel, []byte, error)
	OnRawWriteOp(Channel, []byte)
	OnRawFlushOp(Channel)
}

type ChannelSideOutboundTrigger interface {
	Write(WriteEvent) (submitted bool, err error)
	Flush(FlushEvent) (submitted bool, err error)
	Close() error
	IsActive() bool
	IsInactive() bool
}

type Channel interface {
	ChannelSideOutboundTrigger

	SetAttachment(attachment interface{})
	GetAttachment() interface{}
}

type OutboundEvent interface {
	GetType() EVENT_TYPE
	Process(w io.Writer) (int, error)
}

type WriteEvent struct {
	Data     []byte
	Callback func(writeEvent WriteEvent, err error)
}

func (this WriteEvent) GetType() EVENT_TYPE {
	return EVENT_WRITE
}

func (this WriteEvent) Process(w io.Writer) (int, error) {
	return w.Write(this.Data)
}

type FlushEvent struct {
	Callback func(flushEvent FlushEvent, err error)
}

func (this FlushEvent) GetType() EVENT_TYPE {
	return EVENT_FLUSH
}

func (this FlushEvent) Process(w io.Writer) (int, error) {
	return 0, nil
}

type UserEvent struct {
}

func WrapWriteEvent(data []byte, cb func(writeEvent WriteEvent, err error)) WriteEvent {
	return WriteEvent{
		Data:     data,
		Callback: cb,
	}
}

func NewSimpleWriteEvent(data []byte) WriteEvent {
	return WriteEvent{
		Data:     data,
		Callback: func(writeEvent WriteEvent, err error) {},
	}
}

func WrapFlushEvent(cb func(flushEvent FlushEvent, err error)) FlushEvent {
	return FlushEvent{
		Callback: cb,
	}
}

var noopFlushEvent = FlushEvent{
	Callback: func(flushEvent FlushEvent, err error) {},
}

func NewNoopFlushEvent() FlushEvent {
	return noopFlushEvent
}

//================

type DefaultChannelRawSideHandler struct {
}

var _ ChannelRawSideHandler = DefaultChannelRawSideHandler{}

func (DefaultChannelRawSideHandler) Active(channel Channel) {
}

func (DefaultChannelRawSideHandler) Inactive(channel Channel) {
}

func (DefaultChannelRawSideHandler) OnRead(ch Channel, data []byte) {
}

func (DefaultChannelRawSideHandler) OnError(channel Channel, eventType EVENT_TYPE, err error) {
	switch eventType {
	case EVENT_READ:
		if channel.IsActive() {
			log.Println("default handler:", err)
			channel.Close()
		}
	}
}

func (DefaultChannelRawSideHandler) OnWrite(ch Channel, we WriteEvent) error {
	return nil
}

func (DefaultChannelRawSideHandler) OnFlush(Channel, FlushEvent) error {
	return nil
}

func (DefaultChannelRawSideHandler) OnClose(Channel) error {
	return nil
}

func (DefaultChannelRawSideHandler) OnUserEvent(channel Channel, userEvent *UserEvent) {
}

func (DefaultChannelRawSideHandler) OnRawReadOp(Channel, []byte, error) {
}

func (DefaultChannelRawSideHandler) OnRawWriteOp(Channel, []byte) {
}

func (DefaultChannelRawSideHandler) OnRawFlushOp(Channel) {
}
