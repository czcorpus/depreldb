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
	"testing"

	"github.com/czcorpus/scollector/record"
	"github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/assert"
)

func TestStoreData(t *testing.T) {
	opts := badger.DefaultOptions("").WithInMemory(true)
	bdb, err := badger.Open(opts)
	assert.NoError(t, err, "Failed to open in-memory database")
	defer bdb.Close()

	textTypeMapping := &PreconfTextTypeMapping{
		data: map[string]byte{"test": 0x01},
	}
	db := &DB{bdb: bdb, textTypes: textTypeMapping}
	tidSeq := NewTokenIDSequence()

	// Create test data for singleFreqs
	singleFreqs := map[record.GroupingKey]record.TokenFreq{
		"test1": {
			Lemma:    "word",
			PoS:      record.UDPosFromByte(0x01),    // ADJ
			Deprel:   record.UDDeprelFromByte(0x01), // amod
			Freq:     100,
			TextType: record.TextType{Raw: 0x01, Readable: "test"},
		},
		"test2": {
			Lemma:    "example",
			PoS:      record.UDPosFromByte(0x02),    // NOUN
			Deprel:   record.UDDeprelFromByte(0x02), // nsubj
			Freq:     50,
			TextType: record.TextType{Raw: 0x01, Readable: "test"},
		},
		"test3": {
			Lemma:    "run",
			PoS:      record.UDPosFromByte(0x03),    // VERB
			Deprel:   record.UDDeprelFromByte(0x03), // root
			Freq:     75,
			TextType: record.TextType{Raw: 0x01, Readable: "test"},
		},
	}

	// Create test data for pairFreqs
	pairFreqs := map[record.GroupingKey]record.CollocFreq{
		"pair1": {
			Lemma1:   "word",
			PoS1:     record.UDPosFromByte(0x01),
			Deprel1:  record.UDDeprelFromByte(0x01),
			Lemma2:   "example",
			PoS2:     record.UDPosFromByte(0x02),
			Deprel2:  record.UDDeprelFromByte(0x02),
			Freq:     25,
			AVGDist:  1.5,
			TextType: record.TextType{Raw: 0x01, Readable: "test"},
		},
		"pair2": {
			Lemma1:   "example",
			PoS1:     record.UDPosFromByte(0x02),
			Deprel1:  record.UDDeprelFromByte(0x02),
			Lemma2:   "run",
			PoS2:     record.UDPosFromByte(0x03),
			Deprel2:  record.UDDeprelFromByte(0x03),
			Freq:     10,
			AVGDist:  2.0,
			TextType: record.TextType{Raw: 0x01, Readable: "test"},
		},
		"pair3": {
			Lemma1:   "word",
			PoS1:     record.UDPosFromByte(0x01),
			Deprel1:  record.UDDeprelFromByte(0x01),
			Lemma2:   "run",
			PoS2:     record.UDPosFromByte(0x03),
			Deprel2:  record.UDDeprelFromByte(0x03),
			Freq:     5,
			AVGDist:  1.7,
			TextType: record.TextType{Raw: 0x01, Readable: "test"},
		},
	}

	// Call StoreData
	minPairFreq := 5
	_, err = db.StoreData(tidSeq, singleFreqs, pairFreqs, minPairFreq)
	assert.NoError(t, err, "StoreData should not fail")

	// Verify singleFreqs were stored correctly
	t.Run("VerifySingleFreqs", func(t *testing.T) {
		for _, freq := range singleFreqs {
			// Get lemma ID
			tokenID, err := db.GetLemmaID(freq)
			assert.NoError(t, err, "Should get lemma ID for %s", freq.Lemma)
			assert.NotZero(t, tokenID, "Token ID should be non-zero for lemma %s", freq.Lemma)

			// Get stored frequency
			storedFreqInfos, err := db.GetSingleTokenFreq(tokenID, freq.PoS.Byte(), freq.TextType.Byte(), 0)
			assert.NoError(t, err, "Should get stored frequency for %s", freq.Lemma)
			assert.NotEmpty(t, storedFreqInfos, "Should get at least one frequency entry for %s", freq.Lemma)

			// Sum up all frequencies for this token
			totalFreq := record.SumTokenFreqs(storedFreqInfos)
			assert.Equal(t, freq.Freq, totalFreq, "Total frequency should match for %s", freq.Lemma)

			// Check that at least one entry has the correct deprel
			found := false
			for _, storedFreq := range storedFreqInfos {
				if storedFreq.Deprel == freq.Deprel {
					found = true
					break
				}
			}
			assert.True(t, found, "Should find entry with correct deprel for %s", freq.Lemma)

			// Verify reverse lookup
			lemma, err := db.GetLemmaByID(tokenID)
			assert.NoError(t, err, "Should get lemma by ID for %s", freq.Lemma)
			assert.Equal(t, freq.Lemma, lemma, "Reverse lookup should return correct lemma")
		}
	})

	// Verify pairFreqs were stored correctly
	t.Run("VerifyPairFreqs", func(t *testing.T) {
		for _, pairFreq := range pairFreqs {
			// Skip pairs below minimum frequency threshold
			if pairFreq.Freq < minPairFreq {
				continue
			}

			// Get token IDs for both lemmas using the token ID sequence cache
			tokenID1 := tidSeq.recall(pairFreq.Lemma1Key())
			assert.NotZero(t, tokenID1, "Should get token ID for lemma1 %s", pairFreq.Lemma1)
			tokenID2 := tidSeq.recall(pairFreq.Lemma2Key())
			assert.NotZero(t, tokenID2, "Should get token ID for lemma2 %s", pairFreq.Lemma2)

			// Try to retrieve the pair frequency directly from BadgerDB
			err := db.bdb.View(func(txn *badger.Txn) error {
				key := record.CollFreqKey(
					tokenID1, pairFreq.PoS1.Byte(), pairFreq.TextType.Byte(), pairFreq.Deprel1.Byte(), tokenID2, pairFreq.PoS2.Byte(), pairFreq.Deprel2.Byte())
				item, err := txn.Get(key)
				if err != nil {
					return err
				}

				return item.Value(func(val []byte) error {
					collValue := record.DecodeCollocValue(val)

					assert.Equal(t, uint32(pairFreq.Freq), collValue.Freq,
						"Pair frequency should match for %s-%s", pairFreq.Lemma1, pairFreq.Lemma2)

					// Verify average distance (already decoded)
					decodedDist := collValue.Dist
					expectedDist := pairFreq.AVGDist
					assert.InDelta(t, expectedDist, decodedDist, 0.1,
						"Pair distance should match for %s-%s within 0.1 precision", pairFreq.Lemma1, pairFreq.Lemma2)

					// Note: Dependency relations are encoded in the key, not stored in the protobuf entry

					return nil
				})
			})
			assert.NoError(t, err, "Should retrieve pair frequency for %s-%s",
				pairFreq.Lemma1, pairFreq.Lemma2)
		}
	})

	// Test edge case: pair frequency below minimum threshold should not be stored
	t.Run("VerifyMinFreqFiltering", func(t *testing.T) {
		// Add a pair with frequency below threshold
		lowFreqPair := record.CollocFreq{
			Lemma1:   "low",
			PoS1:     record.UDPosFromByte(0x01),
			Deprel1:  record.UDDeprelFromByte(0x01),
			Lemma2:   "freq",
			PoS2:     record.UDPosFromByte(0x02),
			Deprel2:  record.UDDeprelFromByte(0x02),
			Freq:     2, // Below minPairFreq of 5
			AVGDist:  10.0,
			TextType: record.TextType{Raw: 0x01, Readable: "test"},
		}

		lowFreqSingles := map[record.GroupingKey]record.TokenFreq{
			"low": {
				Lemma:    "low",
				PoS:      record.UDPosFromByte(0x01),
				Deprel:   record.UDDeprelFromByte(0x01),
				Freq:     10,
				TextType: record.TextType{Raw: 0x01, Readable: "test"},
			},
			"freq": {
				Lemma:    "freq",
				PoS:      record.UDPosFromByte(0x02),
				Deprel:   record.UDDeprelFromByte(0x02),
				Freq:     10,
				TextType: record.TextType{Raw: 0x01, Readable: "test"},
			},
		}

		lowFreqPairs := map[record.GroupingKey]record.CollocFreq{
			"lowfreq": lowFreqPair,
		}

		// Store the low frequency data
		_, err = db.StoreData(tidSeq, lowFreqSingles, lowFreqPairs, minPairFreq)
		assert.NoError(t, err, "StoreData should not fail for low frequency test")

		// Verify the pair was NOT stored (should not be found)
		// Get token IDs using the sequence cache
		tokenID1 := tidSeq.recall(lowFreqPair.Lemma1Key())
		tokenID2 := tidSeq.recall(lowFreqPair.Lemma2Key())

		err = db.bdb.View(func(txn *badger.Txn) error {
			key := record.CollFreqKey(
				tokenID1, lowFreqPair.PoS1.Byte(), lowFreqPair.TextType.Byte(), lowFreqPair.Deprel1.Byte(), tokenID2, lowFreqPair.PoS2.Byte(), lowFreqPair.Deprel2.Byte())
			_, err := txn.Get(key)
			return err
		})

		// We expect ErrKeyNotFound since the pair frequency was below threshold
		assert.Equal(t, badger.ErrKeyNotFound, err,
			"Pair with frequency %d should not be stored (below minPairFreq %d)",
			lowFreqPair.Freq, minPairFreq)
	})
}
