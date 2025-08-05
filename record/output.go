package record

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type RawTokenFreq struct {
	TokenID  uint32
	PoS      byte
	Deprel   byte
	Freq     uint32
	TextType byte
}

// BinaryKey represents a binary grouping key for high-performance map operations
type BinaryKey [8]byte

// GroupingKeyBinary creates a binary key (8 bytes) instead of string key
// Layout: [TokenID:4][PoS:1][Deprel:1][TextType:1][padding:1]
func (rtf RawTokenFreq) GroupingKeyBinary() BinaryKey {
	var key BinaryKey
	binary.LittleEndian.PutUint32(key[0:4], rtf.TokenID)
	key[4] = rtf.PoS
	key[5] = rtf.Deprel
	key[6] = rtf.TextType
	// key[7] is padding/unused
	return key
}

// GroupingKey creates a string key - kept for backward compatibility but slower
func (rtf RawTokenFreq) GroupingKey() string {
	var keyBuff strings.Builder
	keyBuff.WriteString(fmt.Sprintf("%d", rtf.TokenID))
	keyBuff.WriteString("|")
	keyBuff.WriteString(fmt.Sprintf("%x", rtf.PoS))
	keyBuff.WriteString("|")
	keyBuff.WriteString(fmt.Sprintf("%x", rtf.Deprel))
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
	keyBuff.WriteByte(hexChar(rtf.Deprel >> 4))
	keyBuff.WriteByte(hexChar(rtf.Deprel & 0xF))
	keyBuff.WriteByte('|')
	keyBuff.WriteByte(hexChar(rtf.TextType >> 4))
	keyBuff.WriteByte(hexChar(rtf.TextType & 0xF))
	return keyBuff.String()
}

// -------------------

type RawCollocFreq struct {
	Token1ID uint32
	PoS1     byte
	Deprel1  byte
	Token2ID uint32
	PoS2     byte
	Deprel2  byte
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
	key[5] = rcf.Deprel1
	binary.LittleEndian.PutUint32(key[6:10], rcf.Token2ID)
	key[10] = rcf.PoS2
	key[11] = rcf.Deprel2
	key[12] = rcf.TextType
	// key[13:16] are padding/unused
	return key
}

// GroupingKeyLemma1Binary creates a binary key for first lemma grouping (8 bytes)
func (rcf RawCollocFreq) GroupingKeyLemma1Binary() BinaryKey {
	var key BinaryKey
	binary.LittleEndian.PutUint32(key[0:4], rcf.Token1ID)
	key[4] = rcf.PoS1
	key[5] = rcf.Deprel1
	key[6] = rcf.TextType
	// key[7] is padding/unused
	return key
}

// GroupingKeyLemma2Binary creates a binary key for second lemma grouping (8 bytes)
func (rcf RawCollocFreq) GroupingKeyLemma2Binary() BinaryKey {
	var key BinaryKey
	binary.LittleEndian.PutUint32(key[0:4], rcf.Token2ID)
	key[4] = rcf.PoS2
	key[5] = rcf.Deprel2
	key[6] = rcf.TextType
	// key[7] is padding/unused
	return key
}

func (rcf RawCollocFreq) GroupingKey() string {

	var keyBuff strings.Builder
	keyBuff.WriteString(fmt.Sprintf("%d", rcf.Token1ID))
	keyBuff.WriteString("|")
	keyBuff.WriteString(fmt.Sprintf("%x", rcf.PoS1))
	keyBuff.WriteString("|")
	keyBuff.WriteString(fmt.Sprintf("%x", rcf.Deprel1))
	keyBuff.WriteString(fmt.Sprintf("|%d", rcf.Token2ID))
	keyBuff.WriteString("|")
	keyBuff.WriteString(fmt.Sprintf("%x", rcf.PoS2))
	keyBuff.WriteString("|")
	keyBuff.WriteString(fmt.Sprintf("%x", rcf.Deprel2))
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
	keyBuff.WriteString(fmt.Sprintf("%x", rcf.Deprel1))
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
	keyBuff.WriteString(fmt.Sprintf("%x", rcf.Deprel2))
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
