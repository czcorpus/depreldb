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
	"github.com/czcorpus/scollector/record"
)

type tokenFreqGrouping struct {
	groupByPos    bool
	groupByTT     bool
	groupByDeprel bool
	data          map[record.BinaryKey]record.RawTokenFreq
}

func (rg *tokenFreqGrouping) Iter(yield func(k record.BinaryKey, v record.RawTokenFreq) bool) {
	for k, v := range rg.data {
		if cont := yield(k, v); !cont {
			break
		}
	}
}

func (rg *tokenFreqGrouping) GroupByPos() *tokenFreqGrouping {
	rg.groupByPos = true
	return rg
}

func (rg *tokenFreqGrouping) GroupByTT() *tokenFreqGrouping {
	rg.groupByTT = true
	return rg
}

func (rg *tokenFreqGrouping) GroupByDeprel() *tokenFreqGrouping {
	rg.groupByDeprel = true
	return rg
}

func (rg *tokenFreqGrouping) add(f record.RawTokenFreq) {
	if !rg.groupByTT {
		f.TextType = 0
	}
	if !rg.groupByPos {
		f.PoS = 0
	}
	if !rg.groupByDeprel {
		f.Deprel = 0
	}
	key := f.GroupingKeyBinary()
	curr, ok := rg.data[key]
	if !ok {
		curr = f

	} else {
		curr.Freq += f.Freq
	}
	rg.data[key] = curr
}

func (rg *tokenFreqGrouping) get(key record.BinaryKey) record.RawTokenFreq {
	return rg.data[key]
}

func newTokenFreqGrouping() *tokenFreqGrouping {
	return &tokenFreqGrouping{
		data: make(map[record.BinaryKey]record.RawTokenFreq),
	}
}

// -----------------------------------------

type collFreqGrouping struct {
	groupByPos1    bool
	groupByDeprel1 bool
	groupByPos2    bool
	groupByDeprel2 bool
	groupByTT      bool
	data           map[record.CollBinaryKey]record.RawCollocFreq
}

func (rg *collFreqGrouping) Iter(yield func(k record.CollBinaryKey, v record.RawCollocFreq) bool) {
	for k, v := range rg.data {
		if cont := yield(k, v); !cont {
			break
		}
	}
}

func (rg *collFreqGrouping) GroupByPos1() *collFreqGrouping {
	rg.groupByPos1 = true
	return rg
}

func (rg *collFreqGrouping) GroupByPos2() *collFreqGrouping {
	rg.groupByPos2 = true
	return rg
}

func (rg *collFreqGrouping) GroupByDeprel1() *collFreqGrouping {
	rg.groupByDeprel1 = true
	return rg
}

func (rg *collFreqGrouping) GroupByDeprel2() *collFreqGrouping {
	rg.groupByDeprel2 = true
	return rg
}

func (rg *collFreqGrouping) GroupByTT() *collFreqGrouping {
	rg.groupByTT = true
	return rg
}

func (rg *collFreqGrouping) add(f record.RawCollocFreq) {
	if !rg.groupByTT {
		f.TextType = 0
	}
	if !rg.groupByPos1 {
		f.PoS1 = 0
	}
	if !rg.groupByPos2 {
		f.PoS2 = 0
	}

	if !rg.groupByDeprel1 {
		f.Deprel1 = 0
	}
	if !rg.groupByDeprel2 {
		f.Deprel2 = 0
	}

	key := f.GroupingKeyBinary()
	curr, ok := rg.data[key]
	if !ok {
		curr = f

	} else {
		curr.Freq += f.Freq
	}
	rg.data[key] = curr
}

func newCollFreqGrouping() *collFreqGrouping {
	return &collFreqGrouping{
		data: make(map[record.CollBinaryKey]record.RawCollocFreq),
	}
}
