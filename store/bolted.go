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
	"github.com/alphazero/borisdb/singleflight"
	"github.com/boltdb/bolt"
)

// count of segments in key-space
const segmentCnt = 8

// metabucket
var dbinfo = []byte("dbinfo")

// type encapsulates boltdb instance and other state info as required.
// this type supports store.KVStore.
// this type supports store.Store.
type boltdb struct {
	db        *bolt.DB
	metaGroup *singleflight.Group
	putGroup  []*singleflight.Group
	getGroup  []*singleflight.Group
}

func OpenDb(name string) (Store, error) {
	bdb, e := bolt.Open(name, 0600, nil)
	if e != nil {
		return nil, fmt.Errorf("err - OpenDb - %s", e)
	}

	// create the store and return
	db := &boltdb{
		db:        bdb,
		metaGroup: &singleflight.Group{},
		putGroup:  make([]*singleflight.Group, segmentCnt),
		getGroup:  make([]*singleflight.Group, segmentCnt),
	}

	e = db.init()
	return db, e
}

func (p *boltdb) init() error {
	for i := 0; i < segmentCnt; i++ {
		p.putGroup[i] = &singleflight.Group{}
		p.getGroup[i] = &singleflight.Group{}
		// create the single toplevel bucket
		bid := []byte(fmt.Sprintf("bucket-%d", i))
		if e := p.db.Update(createBucketFn(bid)); e != nil {
			return e
		}
	}
	if e := p.db.Update(createBucketFn(dbinfo)); e != nil {
		return e
	}
	return nil
}

func createBucketFn(bid []byte) func(*bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists(bid)
		if e != nil {
			return fmt.Errorf("failed to create bucket: %s", e)
		}
		return nil
	}
}

/// interface: Store //////////////////////////////////////////////////////////

// support Store.Close()
func (p *boltdb) Close() error {
	return p.db.Close()
}

// support Store.Info
func (p *boltdb) Info() (value []byte, err error) {
	dbinfo, e := p.metaGroup.Do("update", p.dbinfoUpdateOpFn(0))
	if e != nil {
		err = e // REVU: not too much time but map boltdb errors to ours
		return
	}

	return dbinfo.([]byte), e
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
	gid := segmentFor(key)
	opkey := key.String()
	_, e := p.putGroup[gid].Do(opkey, p.putOpFn(key, v))
	if e != nil {
		err = e // REVU: not too much time but map boltdb errors to ours
		return
	}
	_, e = p.metaGroup.Do("update", p.dbinfoUpdateOpFn(len(v)))
	if e != nil {
		err = e // REVU: not too much time but map boltdb errors to ours
		return
	}

	return
}

// support KVStore.Get
func (p *boltdb) Get(key Key) (value []byte, err error) {
	gid := segmentFor(key)
	opkey := key.String()
	v, e := p.getGroup[gid].Do(opkey, p.getOpFn(key))
	return v.([]byte), e
}

// support KVStore Del
func (p *boltdb) Del(key Key) (value []byte, err error) {
	gid := segmentFor(key)
	opkey := key.String()
	v, e := p.getGroup[gid].Do(opkey, p.delOpFn(key))
	return v.([]byte), e
}

/// internal ops //////////////////////////////////////////////////////////////

func segmentFor(k Key) int {
	return int(k[0] & 0x7)
}

func bucketIdFor(segment int) []byte {
	return []byte(fmt.Sprintf("bucket-%d", segment))
}

/* Get */

func (p *boltdb) getOpFn(k Key) func() (interface{}, error) {
	return func() (interface{}, error) {
		var v []byte
		e := p.db.View(txViewFn(k, &v))
		return v, e
	}
}

func txViewFn(k Key, v *[]byte) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		bid := bucketIdFor(segmentFor(k))
		b := tx.Bucket(bid)
		*v = b.Get(k[:])
		if *v == nil {
			return NotFoundErr
		}
		return nil
	}
}

/* Put */

func (p *boltdb) putOpFn(k Key, v []byte) func() (interface{}, error) {
	return func() (interface{}, error) {
		e := p.db.Update(txUpdateFn(k, v))
		return nil, e
	}
}

func txUpdateFn(k Key, v []byte) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		bid := bucketIdFor(segmentFor(k))
		b := tx.Bucket(bid)
		v0 := b.Get(k[:])
		if v0 != nil {
			return fmt.Errorf("%s - %s", ExistingErr, k.String())
		}
		return b.Put(k[:], v)
	}
}

/* Del */

func (p *boltdb) delOpFn(k Key) func() (interface{}, error) {
	return func() (interface{}, error) {
		var v []byte
		e := p.db.Update(txRemoveFn(k, &v))
		return v, e
	}
}

func txRemoveFn(k Key, v *[]byte) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		/* TODO delete
		bid := bucketIdFor(segmentFor(k))
		b := tx.Bucket(bid)
		*v = b.Get(k[:])
		if *v == nil {
			return NotFoundErr
		}
		return nil
		*/
		return fmt.Errorf("txRemove opfn not implemented")
	}
}

/* update dbinfo */

func (p *boltdb) dbinfoUpdateOpFn(size int) func() (interface{}, error) {
	var infostr string
	return func() (interface{}, error) {
		e := p.db.Update(txUpdateInfoFn(size, &infostr))
		return []byte(infostr), e
	}
}

func txUpdateInfoFn(size int, infostr *string) func(tx *bolt.Tx) error {
	var objcntKey = []byte("object-cnt")
	var sizeKey = []byte("size")
	return func(tx *bolt.Tx) error {
		b := tx.Bucket(dbinfo)

		var v0 []byte
		var e error

		// update size
		v0 = b.Get(sizeKey)
		var totsize = toInt64(v0)
		if size > 0 {
			totsize += int64(size)
			e = b.Put(sizeKey, toByte8(totsize))
			if e != nil {
				return e
			}
		}

		// update object count
		v0 = b.Get(objcntKey)
		var cnt = toInt32(v0)
		if size > 0 {
			cnt++
			e = b.Put(objcntKey, toByte4(cnt))
			if e != nil {
				return e
			}
		}

		*infostr = fmt.Sprintf("dbinfo: object-cnt:%d - totsize:%d\n", cnt, totsize)
		return nil
	}
}

/// temp //////////////////////////////////////////////////////////////////////

func toInt64(b []byte) int64 {
	if b == nil || len(b) < 8 {
		return 0
	}
	return int64(b[0]) |
		int64(b[1])<<8 |
		int64(b[2])<<16 |
		int64(b[3])<<24 |
		int64(b[4])<<32 |
		int64(b[5])<<40 |
		int64(b[6])<<48 |
		int64(b[7])<<56

}
func toInt32(b []byte) int32 {
	if b == nil || len(b) < 4 {
		return 0
	}
	return int32(b[0]) |
		int32(b[1])<<8 |
		int32(b[2])<<16 |
		int32(b[3])<<24
}
func toByte4(n int32) []byte {
	var b = make([]byte, 4)
	b[0] = byte(n & 0xff)
	b[1] = byte((n >> 8) & 0xff)
	b[2] = byte((n >> 16) & 0xff)
	b[3] = byte((n >> 24) & 0xff)
	return b
}

func toByte8(n int64) []byte {
	var b = make([]byte, 8)
	b[0] = byte(n & 0xff)
	b[1] = byte((n >> 8) & 0xff)
	b[2] = byte((n >> 16) & 0xff)
	b[3] = byte((n >> 24) & 0xff)
	b[4] = byte((n >> 32) & 0xff)
	b[5] = byte((n >> 40) & 0xff)
	b[6] = byte((n >> 48) & 0xff)
	b[7] = byte((n >> 56) & 0xff)
	return b
}
