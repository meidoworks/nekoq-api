package network

import "import.moetang.info/go/nekoq-api/errorutil"

var _ERROR_CHANNEL_CLOSED = errorutil.New("channel closed <- network <- nekoq-api")

func ErrChannelClosed() error {
	return _ERROR_CHANNEL_CLOSED
}

var _ERROR_CHANNEL_QUEUE_NOT_READY = errorutil.New("write queue not ready <- network <- nekoq-api")

func ErrChannelQueueNotReady() error {
	return _ERROR_CHANNEL_QUEUE_NOT_READY
}

var _ERROR_UNKNOWN = errorutil.New("unknown <- network <- nekoq-api")

func ErrUnknown() error {
	return _ERROR_UNKNOWN
}
