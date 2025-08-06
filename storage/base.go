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
	"encoding/json"
	"fmt"

	"github.com/czcorpus/scollector/record"
	"github.com/dgraph-io/badger/v4"
	"github.com/rs/zerolog/log"
)

// -----

// DB is a wrapper around badger.DB providing concrete
// methods for adding/retrieving collocation information.
type DB struct {
	bdb       *badger.DB
	textTypes record.TextTypeMapper
	Metadata  Metadata
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

func (db *DB) StoreMetadata(data Metadata) error {
	k := record.CreateMetadataKey(record.MetadataKeyImportProfile)
	if err := db.bdb.Update(func(txn *badger.Txn) error {
		rawMetadata, err := json.Marshal(data)
		if err != nil {
			return err
		}
		if err := txn.Set(k, rawMetadata); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to store import profile: %w", err)
	}
	return nil
}

func (db *DB) readMetadata() (Metadata, error) {
	k := record.CreateMetadataKey(record.MetadataKeyImportProfile)
	var result Metadata
	if err := db.bdb.View(func(txn *badger.Txn) error {
		item, err := txn.Get(k)
		if err != nil {
			return err
		}
		item.Value(func(val []byte) error {
			return json.Unmarshal(val, &result)
		})
		return nil
	}); err != nil {
		return result, fmt.Errorf("failed get profile: %w", err)
	}
	return result, nil
}

// --------

// OpenDBIgnoreMetadata opens a BadgerDB database but does not try
// to fetch index metadata from it. It is suitable e.g. for creating
// new databases or rewriting existing ones.
func OpenDBIgnoreMetadata(path string, textTypes record.TextTypeMapper) (*DB, error) {
	db, err := openDB(path, false)
	if err != nil {
		return nil, err
	}
	db.textTypes = textTypes
	log.Info().
		Str("dbPath", path).
		Msg("setting custom text types for opened database")
	return db, nil
}

// OpenDB opens a BadgerDB database with collocations indexes.
// The database must have proper metadata set as otherwise, it won't open.
// For creating a new db, use OpenDBIgnoreMetadata
func OpenDB(path string) (*DB, error) {
	return openDB(path, true)
}

func openDB(path string, loadProfile bool) (*DB, error) {
	opts := badger.DefaultOptions(path).
		// Read-optimized settings for large datasets
		WithValueLogFileSize(1 << 30). // 1GB value log files for better compression
		WithBlockCacheSize(512 << 20). // 512MB block cache
		WithIndexCacheSize(256 << 20). // 256MB index cache
		WithNumMemtables(2).           // Minimal memtables
		WithNumLevelZeroTables(2).     // Minimal level zero tables
		WithLogger(&ZerologWrapper{})

	ans := &DB{}
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open collocations database: %w", err)
	}
	ans.bdb = db

	if loadProfile {
		metadata, err := ans.readMetadata()
		if err != nil {
			return nil, fmt.Errorf("failed to read data import profile: %w", err)
		}
		ans.Metadata = metadata
		prof := FindProfile(metadata.ProfileName)
		if prof.IsZero() {
			log.Warn().
				Str("profile", metadata.ProfileName).
				Msg("unknown import profile, text types mapping won't be available")

		} else {
			log.Info().
				Str("profile", metadata.ProfileName).
				Int64("corpusSize", metadata.CorpusSize).
				Int("numLemmas", metadata.NumLemmas).
				Int("numLemmaFreqs", metadata.NumLemmaFreqs).
				Int("numCollFreqs", metadata.NumCollFreqs).
				Msg("loaded dataset metadata")
		}
		ans.textTypes = prof.TextTypes
	}

	return ans, nil
}
