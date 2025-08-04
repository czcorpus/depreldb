// Copyright 2025 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2025 Institute of the Czech National Corpus,
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

type bidirEncoding map[string]byte

func (be bidirEncoding) GetRev(val byte) string {
	for k, v := range be {
		if v == val {
			return k
		}
	}
	return ""
}

// ------

type hardcodedTextTypes map[string]byte

func (tt hardcodedTextTypes) RawToReadable(val byte) string {
	for k, v := range tt {
		if v == val {
			return k
		}
	}
	return ""
}

func (tt hardcodedTextTypes) ReadableToRaw(val string) byte {
	return tt[val]
}

// ------

type Profile struct {
	Name          string
	LemmaIdx      int
	PosIdx        int
	ParentIdx     int
	DeprelIdx     int
	TextTypesAttr string
	TextTypes     hardcodedTextTypes
}

func (p Profile) IsZero() bool {
	return p.LemmaIdx == 0 && p.PosIdx == 0 && p.ParentIdx == 0 && p.DeprelIdx == 0
}

func FindProfile(name string) Profile {
	switch name {
	case "intercorp_v16ud":
		return Profile{
			Name:          name,
			LemmaIdx:      4,
			PosIdx:        6,
			ParentIdx:     12,
			DeprelIdx:     11,
			TextTypesAttr: "text.txtype",
			TextTypes: map[string]byte{
				"discussions - transcripts": 0x01,
				"drama":                     0x02,
				"fiction":                   0x03,
				"children's lit.":           0x04,
				"journalism - commentaries": 0x05,
				"journalism - news":         0x06,
				"legal texts":               0x07,
				"nonfiction":                0x08,
				"other":                     0x09,
				"poetry":                    0x0a,
				"religious":                 0x0b,
				"subtitles":                 0x0c,
			},
		}
	default:
		return Profile{}
	}
}
