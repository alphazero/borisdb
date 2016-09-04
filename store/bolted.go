//    Copyright © 2016 Joubin Houshyar. All rights reserved.
//
//    This file is part of Frankinstore.
//
//    Frankinstore is free software: you can redistribute it and/or modify
//    it under the terms of the GNU Affero General Public License as
//    published by the Free Software Foundation, either version 3 of
//    the License, or (at your option) any later version.
//
//    Frankinstore is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU Affero General Public License for more details.
//
//    You should have received a copy of the GNU Affero General Public
//    License along with Frankinstore.  If not, see <http://www.gnu.org/licenses/>.

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

/// internal ops //////////////////////////////////////////////////////////////

func (p *boltdb) getOpFn(k Key, v *[]byte) func() (interface{}, error) {
	return func() (interface{}, error) {
		e := p.db.View(txViewFn(k, v))
		return nil, e
	}
}

func txViewFn(k Key, v *[]byte) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketId)
		*v = b.Get(k[:])
		if v == nil {
			return NotFoundErr
		}
		return nil
	}
}

func (p *boltdb) putOpFn(k Key, v []byte) func() (interface{}, error) {
	return func() (interface{}, error) {
		e := p.db.Update(txUpdateFn(k, v))
		return nil, e
	}
}

func txUpdateFn(k Key, v []byte) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketId)
		v0 := b.Get(k[:])
		if v0 != nil {
			return NopExistingErr
		}
		err := b.Put(k[:], v)
		return err
	}
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

	key = Key(sha1.Sum(v))
	opkey := key.String()
	_, e := p.putGroup.Do(opkey, p.putOpFn(key, v))
	if e != nil {
		err = e // REVU: not too much time but map boltdb errors to ours
		return
	}

	return
}

// support KVStore.Get
func (p *boltdb) Get(key Key) (value []byte, err error) {

	opkey := key.String()
	_, e := p.getGroup.Do(opkey, p.getOpFn(key, &value))
	if e != nil {
		err = e // REVU: not too much time but map boltdb errors to ours
		return
	}

	return
}
