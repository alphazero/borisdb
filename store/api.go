package store

import (
	"crypto/sha1"
	"fmt"
)

// Errors & Warnings
var (
	NopExistingErr   = fmt.Errorf("existing entry")
	NotFoundErr      = fmt.Errorf("entry not found")
	DataCorruptedErr = fmt.Errorf("data corrupted")
	DiskFullErr      = fmt.Errorf("disk full error")
	NilValueErr      = fmt.Errorf("nil value error")
	ZeroValueErr     = fmt.Errorf("zero value error")
	InvalidKeyErr    = fmt.Errorf("key is not compliant to spec.")
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
}

// type defines the general store and data semantics of the storage engine.
type Store interface {
	KVStore
	// Closes the store
	Close()
}
