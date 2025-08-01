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

type keyValMapping map[string]byte

func (udm keyValMapping) GetRev(val byte) string {
	for k, v := range udm {
		if v == val {
			return k
		}
	}
	return ""
}

var UDDeprelMapping = keyValMapping{
	"acl":          0x01,
	"acl:relcl":    0x02,
	"advcl":        0x03,
	"advmod":       0x04,
	"advmod:emph":  0x05,
	"amod":         0x06,
	"appos":        0x07,
	"aux":          0x08,
	"aux:pass":     0x09,
	"case":         0x0a,
	"cc":           0x0b,
	"ccomp":        0x0c,
	"clf":          0x0d,
	"compound":     0x0e,
	"conj":         0x0f,
	"cop":          0x10,
	"csubj":        0x11,
	"csubj:pass":   0x12,
	"dep":          0x13,
	"det":          0x14,
	"det:numgov":   0x15,
	"det:nummod":   0x16,
	"discourse":    0x17,
	"dislocated":   0x18,
	"expl:pass":    0x19,
	"expl:pv":      0x1a,
	"fixed":        0x1b,
	"flat":         0x1c,
	"flat:foreign": 0x1d,
	"goeswith":     0x1e,
	"iobj":         0x1f,
	"list":         0x20,
	"mark":         0x21,
	"nmod":         0x22,
	"nsubj":        0x23,
	"nsubj:pass":   0x24,
	"nummod":       0x25,
	"nummod:gov":   0x26,
	"obj":          0x27,
	"obl":          0x28,
	"obl:arg":      0x29,
	"orphan":       0x2a,
	"parataxis":    0x2b,
	"punct":        0x2c,
	"reparandum":   0x2d,
	"root":         0x2e,
	"vocative":     0x2f,
	"xcomp":        0x30,
}

var UDPoSMapping = keyValMapping{
	"ADJ":   0x01,
	"ADP":   0x02,
	"ADV":   0x03,
	"AUX":   0x04,
	"CCONJ": 0x05,
	"DET":   0x06,
	"INTJ":  0x07,
	"NOUN":  0x08,
	"NUM":   0x09,
	"PRON":  0x0a,
	"PROPN": 0x0b,
	"PUNCT": 0x0c,
	"SCONJ": 0x0d,
	"SYM":   0x0e,
	"VERB":  0x0f,
	"X":     0x10,
}
