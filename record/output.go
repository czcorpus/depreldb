package record

import (
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

// -------------------

type RawCollocFreq struct {
	Token1ID uint32
	PoS1     byte
	Deprel1  byte
	Token2ID uint32
	PoS2     byte
	Deprel2  byte
	Freq     uint32
	AVGDist  int32
	TextType byte
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
