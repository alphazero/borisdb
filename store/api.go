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

// package defines the general API for a content addressable k/v store
// and a default implementation using boltdb
package store

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

// api constants
const (
	DefaultDb = "boris.db"
)

// Errors & Warnings
var (
	ExistingErr      = fmt.Errorf("existing entry")
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

func (k Key) String() string {
	return hex.EncodeToString(k[:])
}

// type defines the interface for a content addressable k/v store.
type KVStore interface {
	// Adds value blob 'val' to store. Returns computed key.
	Put(val []byte) (Key, error)
	// Gets the specified value for 'key', if any.
	Get(key Key) ([]byte, error)
	// Dels the specified value for 'key', if any.
	Del(key Key) ([]byte, error)
}

// type defines the general store and data semantics of the storage engine.
type Store interface {
	KVStore
	// Closes the store
	Close()
	Info() ([]byte, error)
}
