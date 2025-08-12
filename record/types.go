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
	"strings"
)

type UDDeprel struct {
	Readable string
	Raw      uint16
}

func (d UDDeprel) IsValid() bool {
	return d.Raw >= 0x01
}

func (d UDDeprel) AsUint16() uint16 {
	return d.Raw
}

func (d UDDeprel) String() string {
	return d.Readable
}

func (d UDDeprel) AsUint32() uint32 {
	return uint32(d.Raw)
}

func ImportUDDeprel(v string) UDDeprel {
	repr, ok := UDDeprelMapping.Get(strings.ToLower(v))
	if !ok {
		return UDDeprel{Raw: 0x00, Readable: v}
	}
	return UDDeprel{Raw: uint16(repr), Readable: v}
}

func UDDeprelFromUint16(v uint16) UDDeprel {
	readable := UDDeprelMapping.GetRev(v)
	if readable != "" {
		return UDDeprel{Raw: v, Readable: readable}
	}
	return UDDeprel{}
}

// ----

type UDPoS struct {
	Readable string
	Raw      byte
}

func (pos UDPoS) Byte() byte {
	return pos.Raw
}

func (pos UDPoS) String() string {
	return pos.Readable
}

func (pos UDPoS) IsValid() bool {
	return pos.Raw >= 0x01 && pos.Raw <= 0x10
}

func UDPosFromByte(v byte) UDPoS {
	readable := UDPoSMapping.GetRev(v)
	if readable != "" {
		return UDPoS{Raw: v, Readable: readable}
	}
	return UDPoS{}
}

func ImportUDPoS(v string) UDPoS {
	repr, ok := UDPoSMapping[strings.ToUpper(v)]
	if !ok {
		return UDPoS{Raw: 0x00, Readable: v}
	}
	return UDPoS{Raw: repr, Readable: v}
}

// ------

type TextType struct {
	Readable string
	Raw      byte
}

func (tt TextType) Byte() byte {
	return tt.Raw
}

func (tt TextType) String() string {
	return tt.Readable
}

func (tt TextType) AsUint32() uint32 {
	return uint32(tt.Raw)
}

func (tt TextType) IsValid() bool {
	return tt.Raw >= 0x01 && tt.Readable != ""
}

// -------

type TextTypeMapper interface {
	RawToReadable(val byte) string
	ReadableToRaw(val string) byte
}
