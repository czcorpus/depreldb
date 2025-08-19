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

package record

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type RawTokenFreq struct {
	TokenID  uint32
	PoS      byte
	Freq     uint32
	TextType byte
}

// BinaryKey represents a binary grouping key for high-performance map operations
type BinaryKey [6]byte

// GroupingKeyBinary creates a binary key (8 bytes) instead of string key
// Layout: [TokenID:4][PoS:1][Deprel:2][TextType:1][padding:1]
func (rtf RawTokenFreq) GroupingKeyBinary() BinaryKey {
	var key BinaryKey
	binary.LittleEndian.PutUint32(key[0:4], rtf.TokenID)
	key[4] = rtf.PoS
	key[5] = rtf.TextType
	return key
}

// GroupingKey creates a string key - kept for backward compatibility but slower
func (rtf RawTokenFreq) GroupingKey() string {
	var keyBuff strings.Builder
	keyBuff.WriteString(fmt.Sprintf("%d", rtf.TokenID))
	keyBuff.WriteString("|")
	keyBuff.WriteString(fmt.Sprintf("%x", rtf.PoS))
	keyBuff.WriteString("|")
	keyBuff.WriteString(fmt.Sprintf("%x", rtf.TextType))
	return keyBuff.String()
}

// GroupingKeyOptimized creates an optimized string key using faster number formatting
func (rtf RawTokenFreq) GroupingKeyOptimized() string {
	// Pre-allocate with estimated capacity to reduce allocations
	var keyBuff strings.Builder
	keyBuff.Grow(16) // Estimate: 4-8 digits + separators + hex values

	// Avoid fmt.Sprintf for simple integer formatting
	keyBuff.WriteString(uitoa(uint64(rtf.TokenID)))
	keyBuff.WriteByte('|')
	keyBuff.WriteByte(hexChar(rtf.PoS >> 4))
	keyBuff.WriteByte(hexChar(rtf.PoS & 0xF))
	keyBuff.WriteByte('|')
	keyBuff.WriteByte(hexChar(rtf.TextType >> 4))
	keyBuff.WriteByte(hexChar(rtf.TextType & 0xF))
	return keyBuff.String()
}

// -------------------

type RawCollocFreq struct {
	Token1ID uint32
	PoS1     byte
	Deprel   uint16
	Token2ID uint32
	PoS2     byte
	Freq     uint32
	AVGDist  float64
	TextType byte
}

// CollBinaryKey represents a binary grouping key for collocation data (16 bytes)
type CollBinaryKey [16]byte

// GroupingKeyBinary creates a binary key for full collocation grouping
// Layout: [Token1ID:4][PoS1:1][Deprel1:1][Token2ID:4][PoS2:1][Deprel2:1][TextType:1][padding:3]
func (rcf RawCollocFreq) GroupingKeyBinary() CollBinaryKey {
	var key CollBinaryKey
	binary.LittleEndian.PutUint32(key[0:4], rcf.Token1ID)
	key[4] = rcf.PoS1
	binary.LittleEndian.PutUint16(key[5:7], rcf.Deprel)
	binary.LittleEndian.PutUint32(key[7:11], rcf.Token2ID)
	key[11] = rcf.PoS2
	key[12] = rcf.TextType
	// key[13:16] is padding/unused
	return key
}

// GroupingKeyLemma1Binary creates a binary key for first lemma grouping (8 bytes)
func (rcf RawCollocFreq) GroupingKeyLemma1Binary() BinaryKey {
	var key BinaryKey
	binary.LittleEndian.PutUint32(key[0:4], rcf.Token1ID)
	key[4] = rcf.PoS1
	key[5] = rcf.TextType
	return key
}

// GroupingKeyLemma2Binary creates a binary key for second lemma grouping (8 bytes)
func (rcf RawCollocFreq) GroupingKeyLemma2Binary() BinaryKey {
	var key BinaryKey
	binary.LittleEndian.PutUint32(key[0:4], rcf.Token2ID)
	key[4] = rcf.PoS2
	key[5] = rcf.TextType
	return key
}

func (rcf RawCollocFreq) GroupingKey() string {

	var keyBuff strings.Builder
	if rcf.AVGDist > 0 {
		keyBuff.WriteString("H")

	} else {
		keyBuff.WriteString("D")
	}
	keyBuff.WriteString(fmt.Sprintf("%d", rcf.Token1ID))
	keyBuff.WriteString("|")
	keyBuff.WriteString(fmt.Sprintf("%x", rcf.PoS1))
	keyBuff.WriteString("|")
	keyBuff.WriteString(fmt.Sprintf("%x", rcf.Deprel))
	keyBuff.WriteString(fmt.Sprintf("|%d", rcf.Token2ID))
	keyBuff.WriteString("|")
	keyBuff.WriteString(fmt.Sprintf("%x", rcf.PoS2))
	keyBuff.WriteString("|")
	keyBuff.WriteString(fmt.Sprintf("%x", rcf.TextType))
	return keyBuff.String()
}

func (rcf RawCollocFreq) GroupingKeyLemma1() string {

	var keyBuff strings.Builder
	keyBuff.WriteString(fmt.Sprintf("%d", rcf.Token1ID))
	keyBuff.WriteString("|")
	keyBuff.WriteString(fmt.Sprintf("%x", rcf.PoS1))
	keyBuff.WriteString("|")
	keyBuff.WriteString(fmt.Sprintf("%x", rcf.Deprel))
	keyBuff.WriteString("|")
	keyBuff.WriteString(fmt.Sprintf("%x", rcf.TextType))
	return keyBuff.String()
}

func (rcf RawCollocFreq) GroupingKeyLemma2() string {

	var keyBuff strings.Builder
	keyBuff.WriteString(fmt.Sprintf("%d", rcf.Token2ID))
	keyBuff.WriteString("|")
	keyBuff.WriteString(fmt.Sprintf("%x", rcf.PoS2))
	keyBuff.WriteString("|")
	keyBuff.WriteString(fmt.Sprintf("%x", rcf.Deprel))
	keyBuff.WriteString("|")
	keyBuff.WriteString(fmt.Sprintf("%x", rcf.TextType))
	return keyBuff.String()
}

// Helper functions for optimized string formatting
func hexChar(b byte) byte {
	if b < 10 {
		return '0' + b
	}
	return 'a' + b - 10
}

// Fast unsigned integer to string conversion
func uitoa(u uint64) string {
	if u == 0 {
		return "0"
	}

	// Estimate buffer size
	var buf [20]byte // uint64 max is 20 digits
	i := len(buf)

	for u > 0 {
		i--
		buf[i] = byte(u%10) + '0'
		u /= 10
	}

	return string(buf[i:])
}
