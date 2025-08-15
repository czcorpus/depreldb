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

import "fmt"

const (
	DeprelAcl         = 0x0001
	DeprelAclRelcl    = 0x0002
	DeprelAdvcl       = 0x0003
	DeprelAdvmod      = 0x0004
	DeprelAdvmodEmph  = 0x0005
	DeprelAmod        = 0x0006
	DeprelAppos       = 0x0007
	DeprelAux         = 0x0008
	DeprelAuxPass     = 0x0009
	DeprelCase        = 0x000a
	DeprelCc          = 0x000b
	DeprelCcomp       = 0x000c
	DeprelClf         = 0x000d
	DeprelCompound    = 0x000e
	DeprelConj        = 0x000f
	DeprelCop         = 0x0010
	DeprelCsubj       = 0x0011
	DeprelCsubjPass   = 0x0012
	DeprelDep         = 0x0013
	DeprelDet         = 0x0014
	DeprelDetNumgov   = 0x0015
	DeprelDetNummod   = 0x0016
	DeprelDiscourse   = 0x0017
	DeprelDislocated  = 0x0018
	DeprelExplPass    = 0x0019
	DeprelExplPv      = 0x001a
	DeprelFixed       = 0x001b
	DeprelFlat        = 0x001c
	DeprelFlatForeign = 0x001d
	DeprelGoeswith    = 0x001e
	DeprelIobj        = 0x001f
	DeprelList        = 0x0020
	DeprelMark        = 0x0021
	DeprelNmod        = 0x0022
	DeprelNsubj       = 0x0023
	DeprelNsubjPass   = 0x0024
	DeprelNummod      = 0x0025
	DeprelNummodGov   = 0x0026
	DeprelObj         = 0x0027
	DeprelObl         = 0x0028
	DeprelOblArg      = 0x0029
	DeprelOrphan      = 0x002a
	DeprelParataxis   = 0x002b
	DeprelPunct       = 0x002c
	DeprelReparandum  = 0x002d
	DeprelRoot        = 0x002e
	DeprelVocative    = 0x002f
	DeprelXcomp       = 0x0030

	PosADJ         = 0x01
	PosADP         = 0x02 // (includes prepositions)
	PosADV         = 0x03
	PosAUX         = 0x04
	PosCCONJ       = 0x05
	PosDET         = 0x06
	PosINTJ        = 0x07
	PosNOUN        = 0x08
	PosNUM         = 0x09
	PosPRON        = 0x0a
	PosPROPN       = 0x0b
	PosPUNCT       = 0x0c
	PosSCONJ       = 0x0d
	PosSYM         = 0x0e
	PosVERB        = 0x0f
	PosX           = 0x10
	PosPART        = 0x11
	PosSCONJ_AUX   = 0x12
	PosPRON_AUX    = 0x13
	PosADP_PRON    = 0x14
	PosVERB_AUX    = 0x15
	PosPROPN_AUX   = 0x16
	PosNOUN_NOUN   = 0x17
	PosX_AUX       = 0x18
	PosNOUN_AUX    = 0x19
	PosPROPN_NOUN  = 0x1a
	PosPART_AUX    = 0x1b
	PosPROPN_PROPN = 0x1c
)

// DeprelMapping allows for mapping between string names/codes
// of core deprel values and their internal byte representation.
// It's native mapping is from strings to bytes but it is also
// handle repeated reverse lookups via caching.
type DeprelMapping struct {
	items    map[string]uint16
	revCache map[uint16]string
	maxValue uint16
}

// Get provides a byte representation based on string name/code.
func (udm *DeprelMapping) Get(key string) (uint16, bool) {
	v, ok := udm.items[key]
	return v, ok
}

// Register registers additional deprel value by automatically
// attaching a new byte value for it.
//
// The method is expected to be used mainly during a corpus/text
// import when scollector is performing syntax tree editing by
// removing/shrinking "useless" nodes.
//
// Calling the method with an already registered key causes panic.
func (udm *DeprelMapping) Register(key string) {
	if _, test := udm.items[key]; test {
		panic(fmt.Errorf("cannot register deprel value - %s is aleady registered", key))
	}
	udm.items[key] = udm.maxValue
	udm.maxValue++
}

func (udm *DeprelMapping) GetRev(val uint16) string {
	v, ok := udm.revCache[val]
	if ok {
		return v
	}
	for k, v := range udm.items {
		if v == val {
			if udm.revCache == nil {
				udm.revCache = make(map[uint16]string)
			}
			udm.revCache[val] = k
			return k
		}
	}
	return ""
}

// AsMap returns the internal mapping representation
// (i.e. string representation => byte code)
func (udm *DeprelMapping) AsMap() map[string]uint16 {
	return udm.items
}

// DeprelMappingFromMap is used for instantiating (possibly extended) deprel
// maps for a specific corpus/dataset based on stored metadata. 
func DeprelMappingFromMap(src map[string]uint16) *DeprelMapping {
	return &DeprelMapping{
		items:    src,
		revCache: map[uint16]string{},
	}
}

var UDDeprelMapping = DeprelMapping{
	maxValue: 0x100,
	items: map[string]uint16{
		"acl":          DeprelAcl,
		"acl:relcl":    DeprelAclRelcl,
		"advcl":        DeprelAdvcl,
		"advmod":       DeprelAdvmod,
		"advmod:emph":  DeprelAdvmodEmph,
		"amod":         DeprelAmod,
		"appos":        DeprelAppos,
		"aux":          DeprelAux,
		"aux:pass":     DeprelAuxPass,
		"case":         DeprelCase,
		"cc":           DeprelCc,
		"ccomp":        DeprelCcomp,
		"clf":          DeprelClf,
		"compound":     DeprelCompound,
		"conj":         DeprelConj,
		"cop":          DeprelCop,
		"csubj":        DeprelCsubj,
		"csubj:pass":   DeprelCsubjPass,
		"dep":          DeprelDep,
		"det":          DeprelDet,
		"det:numgov":   DeprelDetNumgov,
		"det:nummod":   DeprelDetNummod,
		"discourse":    DeprelDiscourse,
		"dislocated":   DeprelDislocated,
		"expl:pass":    DeprelExplPass,
		"expl:pv":      DeprelExplPv,
		"fixed":        DeprelFixed,
		"flat":         DeprelFlat,
		"flat:foreign": DeprelFlatForeign,
		"goeswith":     DeprelGoeswith,
		"iobj":         DeprelIobj,
		"list":         DeprelList,
		"mark":         DeprelMark,
		"nmod":         DeprelNmod,
		"nsubj":        DeprelNsubj,
		"nsubj:pass":   DeprelNsubjPass,
		"nummod":       DeprelNummod,
		"nummod:gov":   DeprelNummodGov,
		"obj":          DeprelObj,
		"obl":          DeprelObl,
		"obl:arg":      DeprelOblArg,
		"orphan":       DeprelOrphan,
		"parataxis":    DeprelParataxis,
		"punct":        DeprelPunct,
		"reparandum":   DeprelReparandum,
		"root":         DeprelRoot,
		"vocative":     DeprelVocative,
		"xcomp":        DeprelXcomp,
		// dynamically generated values should start with 0x0100
	},
}

type posMapping map[string]byte

func (pm posMapping) GetRev(val byte) string {
	for k, v := range pm {
		if v == val {
			return k
		}
	}
	return ""
}

var UDPoSMapping = posMapping{
	"ADJ":         PosADJ,
	"ADP":         PosADP, // (includes prepositions)
	"ADV":         PosADV,
	"AUX":         PosAUX,
	"CCONJ":       PosCCONJ,
	"DET":         PosDET,
	"INTJ":        PosINTJ,
	"NOUN":        PosNOUN,
	"NUM":         PosNUM,
	"PRON":        PosPRON,
	"PROPN":       PosPROPN,
	"PUNCT":       PosPUNCT,
	"SCONJ":       PosSCONJ,
	"SYM":         PosSYM,
	"VERB":        PosVERB,
	"X":           PosX,
	"PART":        PosPART,
	"SCONJ|AUX":   PosSCONJ_AUX,
	"PRON|AUX":    PosPRON_AUX,
	"ADP|PRON":    PosADP_PRON,
	"VERB|AUX":    PosVERB_AUX,
	"PROPN|AUX":   PosPROPN_AUX,
	"NOUN|NOUN":   PosNOUN_NOUN,
	"X|AUX":       PosX_AUX,
	"NOUN|AUX":    PosNOUN_AUX,
	"PROPN|NOUN":  PosPROPN_NOUN,
	"PART|AUX":    PosPART_AUX,
	"PROPN|PROPN": PosPROPN_PROPN,
}
