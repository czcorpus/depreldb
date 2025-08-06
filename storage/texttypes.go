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

// PreconfTextTypeMapping represents a type providing mapping
// between text types encoded as byte values and their actual human
// readable string value.
type PreconfTextTypeMapping struct {
	data map[string]byte
}

func (mapping *PreconfTextTypeMapping) RawToReadable(val byte) string {
	for k, v := range mapping.data {
		if v == val {
			return k
		}
	}
	return ""
}

func (mapping *PreconfTextTypeMapping) ReadableToRaw(val string) byte {
	return mapping.data[val]
}

func NewPreconfTextTypeMapping(data map[string]byte) *PreconfTextTypeMapping {
	normData := data
	if normData == nil {
		normData = map[string]byte{}
	}
	return &PreconfTextTypeMapping{
		data: normData,
	}
}
