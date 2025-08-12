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

type DeprelMapping struct {
	items    map[string]uint16
	revCache map[uint16]string
	maxValue uint16
}

func (udm *DeprelMapping) Get(key string) (uint16, bool) {
	v, ok := udm.items[key]
	return v, ok
}

func (udm *DeprelMapping) Register(key string) {
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

func (udm *DeprelMapping) AsMap() map[string]uint16 {
	return udm.items
}

func DeprelMappingFromMap(src map[string]uint16) *DeprelMapping {
	return &DeprelMapping{
		items:    src,
		revCache: map[uint16]string{},
	}
}

var UDDeprelMapping = DeprelMapping{
	maxValue: 0x100,
	items: map[string]uint16{
		"acl":          0x0001,
		"acl:relcl":    0x0002,
		"advcl":        0x0003,
		"advmod":       0x0004,
		"advmod:emph":  0x0005,
		"amod":         0x0006,
		"appos":        0x0007,
		"aux":          0x0008,
		"aux:pass":     0x0009,
		"case":         0x000a,
		"cc":           0x000b,
		"ccomp":        0x000c,
		"clf":          0x000d,
		"compound":     0x000e,
		"conj":         0x000f,
		"cop":          0x0010,
		"csubj":        0x0011,
		"csubj:pass":   0x0012,
		"dep":          0x0013,
		"det":          0x0014,
		"det:numgov":   0x0015,
		"det:nummod":   0x0016,
		"discourse":    0x0017,
		"dislocated":   0x0018,
		"expl:pass":    0x0019,
		"expl:pv":      0x001a,
		"fixed":        0x001b,
		"flat":         0x001c,
		"flat:foreign": 0x001d,
		"goeswith":     0x001e,
		"iobj":         0x001f,
		"list":         0x0020,
		"mark":         0x0021,
		"nmod":         0x0022,
		"nsubj":        0x0023,
		"nsubj:pass":   0x0024,
		"nummod":       0x0025,
		"nummod:gov":   0x0026,
		"obj":          0x0027,
		"obl":          0x0028,
		"obl:arg":      0x0029,
		"orphan":       0x002a,
		"parataxis":    0x002b,
		"punct":        0x002c,
		"reparandum":   0x002d,
		"root":         0x002e,
		"vocative":     0x002f,
		"xcomp":        0x0030,
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
	"ADJ":         0x01,
	"ADP":         0x02, // (includes prepositions)
	"ADV":         0x03,
	"AUX":         0x04,
	"CCONJ":       0x05,
	"DET":         0x06,
	"INTJ":        0x07,
	"NOUN":        0x08,
	"NUM":         0x09,
	"PRON":        0x0a,
	"PROPN":       0x0b,
	"PUNCT":       0x0c,
	"SCONJ":       0x0d,
	"SYM":         0x0e,
	"VERB":        0x0f,
	"X":           0x10,
	"PART":        0x11,
	"SCONJ|AUX":   0x12,
	"PRON|AUX":    0x13,
	"ADP|PRON":    0x14,
	"VERB|AUX":    0x15,
	"PROPN|AUX":   0x16,
	"NOUN|NOUN":   0x17,
	"X|AUX":       0x18,
	"NOUN|AUX":    0x19,
	"PROPN|NOUN":  0x1a,
	"PART|AUX":    0x1b,
	"PROPN|PROPN": 0x1c,
}
