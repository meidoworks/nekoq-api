package db

type SimpleDB interface {
	Close() error

	Get(key []byte) ([]byte, bool, error)
	Put(key, data []byte) error
	RangeGetFrom(key []byte, limit int) (keys [][]byte, values [][]byte, err error)
}
