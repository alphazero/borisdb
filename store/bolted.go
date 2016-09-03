package store

import (
	"fmt"
	"github.com/boltdb/bolt"
)

// default bucket name
// REVU: not very familiar with Bold but a single bucket seems to be OK
const bucketId = "root"

// type encapsulates boltdb instance and other state info as required.
// this type supports store.KVStore.
// this type supports store.Store.
type boltdb struct {
	db *bolt.DB
}

func OpenDb(name string) (Store, error) {
	bdb, e := bolt.Open(name, 0600, nil)
	if e != nil {
		return nil, fmt.Errorf("err - OpenDb - %s", e)
	}

	db := &boltdb{
		db: bdb,
	}
	return db, nil
}

// support Store.Close()
func (p *boltdb) Close() {
	p.db.Close()
}

// support KVStore.Put
func (p *boltdb) Put(v []byte) (Key, error) {
	panic("not implemented")
}

// support KVStore.Get
func (p *boltdb) Get(k Key) ([]byte, error) {
	panic("not implemented")
}
