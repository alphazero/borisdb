package store

import (
	"crypto/sha1"
	"fmt"
)

// value blob keys are sha1 digests
const KeySize = sha1.Size

// Keys are immutable byte arrays
type Key [KeySize]byte

// type defines the interface for a content addressable k/v store.
type KVStore interface {
	// Adds value blob 'val' to store. Returns computed key.
	Put(val []byte) (Key, error)
	// Gets the specified value for 'key', if any.
	Get(key Key) ([]byte, error)
	// Closes the store
	Close()
}

type Store interface {
	KVStore
}

func OpenDb(name string) (Store, error) {
	return nil, fmt.Errorf("OpenDB not implemented!")
}
