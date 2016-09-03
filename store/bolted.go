package store

import (
	"crypto/sha1"
	"fmt"
	"frankinstore/singleflight"
	"github.com/boltdb/bolt"
)

// default bucket name
// REVU: not very familiar with Bold but a single bucket seems to be OK
var bucketId = []byte("root")

// type encapsulates boltdb instance and other state info as required.
// this type supports store.KVStore.
// this type supports store.Store.
type boltdb struct {
	db       *bolt.DB
	putGroup *singleflight.Group
	getGroup *singleflight.Group
}

func OpenDb(name string) (Store, error) {
	bdb, e := bolt.Open(name, 0600, nil)
	if e != nil {
		return nil, fmt.Errorf("err - OpenDb - %s", e)
	}

	// create the single toplevel bucket
	bdb.Update(func(tx *bolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists(bucketId)
		if e != nil {
			return fmt.Errorf("err - failed to create bucket: %s", e)
		}
		return nil
	})

	// create the store and return
	db := &boltdb{
		db:       bdb,
		putGroup: &singleflight.Group{},
		getGroup: &singleflight.Group{},
	}

	return db, nil
}

/// interface: Store //////////////////////////////////////////////////////////

// support Store.Close()
func (p *boltdb) Close() {
	p.db.Close()
}

/// interface: KVStore ////////////////////////////////////////////////////////

// support KVStore.Put
// computes sha1 hash of value and stores the blob.
// nil or zerovalue values are not accepted.
func (p *boltdb) Put(v []byte) (key Key, err error) {
	/* assert constraints */
	if v == nil {
		err = NilValueErr
		return
	}
	if len(v) == 0 {
		err = ZeroValueErr
		return
	}

	// compute key
	key = Key(sha1.Sum(v))

	/* TODO: use singleflight here */

	// singleflight insures concurrent putts for the same key
	// result in a single call to the db.
	_, e := p.putGroup.Do(string(key.String()), func() (interface{}, error) {
		e := p.db.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket(bucketId)
			existing := bucket.Get(key[:])
			if existing != nil {
				return NopExistingErr
			}
			err := bucket.Put(key[:], v)
			return err
		})
		return nil, e
	})

	if e != nil {
		err = e // REVU: not too much time but map boltdb errors to ours
		return
	}

	return
}

// support KVStore.Get
func (p *boltdb) Get(k Key) (value []byte, err error) {
	/* assert constraints */
	if len(k) != KeySize {
		err = InvalidKeyErr
		return
	}

	// singleflight insures concurrent gets for the same key
	// result in a single call to the db.
	_, e := p.getGroup.Do(string(k.String()), func() (interface{}, error) {
		e := p.db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket(bucketId)
			value = bucket.Get(k[:])
			if value == nil {
				return NotFoundErr
			}
			return nil
		})
		return value, e
	})

	if e != nil {
		err = e // REVU: not too much time but map boltdb errors to ours
		return
	}

	return
}
