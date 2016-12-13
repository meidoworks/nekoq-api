package db

type SequenceKey struct {
	key  string
	hash int
}

func (this SequenceKey) HashCode() int32 {
	return int32(this.hash)
}

func MakeSequenceKey(key string) SequenceKey {
	seq := SequenceKey{
		key: key,
	}
	seq.hash = bkdrHash([]byte(key))
	return seq
}

func bkdrHash(char []byte) int {
	seed := 131
	hash := 0

	for i := 0; i < len(char); i++ {
		hash = hash*seed + int(char[i])
	}

	return (hash & 0x7FFFFFFF)
}

type AtomicDB interface {
	Close() error

	Incr(key SequenceKey, step int64) (start, end int64, err error)

	AtomicGet(key SequenceKey) ([]byte, bool, error)
	AtomicSet(key SequenceKey, val []byte) error
	CompareAndSet(key SequenceKey, oldVal, newVal []byte) (swap bool, err error)
}
