// Copyright 2025 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2025 Department of Linguistics,
//                Faculty of Arts, Charles University
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"fmt"

	"github.com/czcorpus/scollector/record"
	"github.com/dgraph-io/badger/v4"
)

// -----

// DB is a wrapper around badger.DB providing concrete
// methods for adding/retrieving collocation information.
type DB struct {
	bdb       *badger.DB
	textTypes record.TextTypeMapper
}

// Close closes the internal Badger database.
// It is necessary to perform the close especially
// in cases of data writing.
// It is possible to call the method on nil instance
// or on an uninitialized DB object, in which case
// it is a NOP.
func (db *DB) Close() error {
	if db != nil && db.bdb != nil {
		return db.bdb.Close()
	}
	return nil
}

func (db *DB) Clear() error {
	return db.bdb.DropAll()
}

func (db *DB) Size() (int64, int64) {
	return db.bdb.Size()
}

func (db *DB) StoreImportProfile(profileName string) error {
	k := record.CreateMetadataKey(record.MetadataKeyImportProfile)
	if err := db.bdb.Update(func(txn *badger.Txn) error {
		if err := txn.Set(k, []byte(profileName)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to store import profile: %w", err)
	}
	return nil
}

func (db *DB) ReadImportProfile() (string, error) {
	k := record.CreateMetadataKey(record.MetadataKeyImportProfile)
	var result string
	if err := db.bdb.View(func(txn *badger.Txn) error {
		item, err := txn.Get(k)
		if err != nil {
			return err
		}
		item.Value(func(val []byte) error {
			result = string(val)
			return nil
		})
		return nil
	}); err != nil {
		return result, fmt.Errorf("failed to store profile: %w", err)
	}
	return result, nil
}
