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

// tokenIDSequence is a generator of unique sequential integer identifiers (1, 2, ...)
// for storing lemmas (along with PoS if available)
type tokenIDSequence struct {
	value uint32
	cache map[string]uint32 // key is a hashed mix of lemma and PoS
}

// next generates next ID in the stored sequence.
// Please note that calling the method with the same lemma produces
// new ID each time. To test if a lemma has already been registered,
// use recall().
func (tseq *tokenIDSequence) next(lemmaHash string) uint32 {
	tseq.value++
	if tseq.value == 0 {
		panic("tokenIDSequence overflow")
	}
	tseq.cache[lemmaHash] = tseq.value
	return tseq.value
}

func (tseq *tokenIDSequence) nextIfNotFound(lemmaHash string) (uint32, bool) {
	nextID := tseq.recall(lemmaHash)
	found := true
	if nextID == 0 {
		nextID = tseq.next(lemmaHash)
		found = false
	}
	return nextID, found
}

// recall returns ID of an already registered lemma. If not found,
// zero is returned (numbers are generated from 1 so the distinction is clear)
func (tseq *tokenIDSequence) recall(lemmaHash string) uint32 {
	// zero means = not found (we serve ids from 1)
	return tseq.cache[lemmaHash]
}

// NewTokenIDSequence creates a properly initialized
// ID sequence generator
func NewTokenIDSequence() *tokenIDSequence {
	return &tokenIDSequence{
		value: 0,
		cache: make(map[string]uint32),
	}
}

// --------------

func (db *DB) StoreSingleTokenFreqTx(txn *badger.Txn, tokenID uint32, freq record.TokenFreq) error {
	key := record.TokenFreqKey(tokenID, freq.PoS.Byte(), freq.TextType.Byte(), freq.Deprel.AsUint16())
	encoded := record.EncodeTokenValue(uint32(freq.Freq))
	return txn.Set(key, encoded)
}

func (db *DB) StorePairTokenFreqTx(txn *badger.Txn, token1ID, token2ID uint32, collFreq record.CollocFreq) error {
	key := record.CollFreqKey(
		token1ID, collFreq.PoS1.Byte(), collFreq.TextType.Byte(), collFreq.Deprel1.AsUint16(),
		token2ID, collFreq.PoS2.Byte(), collFreq.Deprel2.AsUint16())
	encoded := record.EncodeCollocValue(uint32(collFreq.Freq), collFreq.AVGDist)
	return txn.Set(key, encoded)
}

func (db *DB) CreateTransaction() *badger.Txn {
	return db.bdb.NewTransaction(true)
}

func (db *DB) StoreLemmaTx(txn *badger.Txn, lemma record.TokenFreq, tokenID uint32) error {
	key := record.EncodeLemmaKey(lemma)
	value := record.TokenIDToBytes(tokenID)
	if err := txn.Set(key, value); err != nil {
		return err
	}
	// Store tokenID -> lemma mapping (reverse index)
	idKey := record.TokenIDToRevIndexKey(tokenID)
	return txn.Set(idKey, []byte(lemma.Lemma))
}

type ImportStats struct {
	NumCollFreqs  int
	NumLemmaFreqs int
	NumLemmas     int
}

func (db *DB) StoreData(
	tidSeq *tokenIDSequence,
	singleFreqs map[record.GroupingKey]record.TokenFreq,
	pairFreqs map[record.GroupingKey]record.CollocFreq,
	minPairFreq int,
) (ImportStats, error) {
	var res ImportStats
	// use singleFreqs as source of lemmas and create indexes
	for _, lemmaEntry := range singleFreqs {

		err := db.bdb.Update(func(txn *badger.Txn) error {
			nextId, alreadyStored := tidSeq.nextIfNotFound(lemmaEntry.LemmaKey())
			if alreadyStored {
				return nil
			}
			if err := db.StoreLemmaTx(txn, lemmaEntry, nextId); err != nil {
				return err
			}
			res.NumLemmas++
			return nil
		})
		if err != nil {
			return res, fmt.Errorf("failed to store lemma: %w", err)
		}
	}

	// Process single token frequencies
	for _, lemmaEntry := range singleFreqs {
		err := db.bdb.Update(func(txn *badger.Txn) error {
			if err := db.StoreSingleTokenFreqTx(txn, tidSeq.recall(lemmaEntry.LemmaKey()), lemmaEntry); err != nil {
				return err
			}
			res.NumLemmaFreqs++
			return nil
		})
		if err != nil {
			return res, fmt.Errorf("failed to store single freq: %w", err)
		}
	}

	// Process pair frequencies
	for _, pairFreq := range pairFreqs {
		if pairFreq.Freq < minPairFreq {
			continue
		}
		err := db.bdb.Update(func(txn *badger.Txn) error {
			if err := db.StorePairTokenFreqTx(
				txn,
				tidSeq.recall(pairFreq.Lemma1Key()),
				tidSeq.recall(pairFreq.Lemma2Key()),
				pairFreq,
			); err != nil {
				return err
			}
			res.NumCollFreqs++
			return nil
		})
		if err != nil {
			return res, fmt.Errorf("failed to store pair freq: %w", err)
		}
	}

	return res, nil
}
