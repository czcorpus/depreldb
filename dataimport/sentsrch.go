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

package dataimport

import (
	"github.com/czcorpus/cnc-gokit/collections"
	"github.com/czcorpus/scollector/storage"
	"github.com/tomachalek/vertigo/v6"
)

type FreqsCollector interface {
	AddLemma(lemma *vertigo.Token, freq int)
	AddCooc(lemma1, lemma2 *vertigo.Token, freq int, distance int)
	ImportSentence(sent []*vertigo.Token)
	PrintPreview()
	StoreToDb(db *storage.DB, minFreq int) (storage.ImportStats, error)
}

// ----------------------------

type Searcher struct {
	prevTokens       *collections.CircularList[*vertigo.Token]
	lastTokenIdx     int
	lastSentStartIdx int
	lastSentEndIdx   int
	foundNewSent     bool
	lemmaIdx         int
	parentIdx        int
	deprelIdx        int
	freqs            FreqsCollector
	corpusSize       int64
}

func (vf *Searcher) analyzeLastSent() {
	var sentOpen bool
	sent := make([]*vertigo.Token, 0, vf.lastSentEndIdx-vf.lastSentStartIdx+1)
	vf.prevTokens.ForEach(func(i int, item *vertigo.Token) bool {
		if item.Idx == vf.lastSentStartIdx {
			sentOpen = true
		}
		if sentOpen {
			sent = append(sent, item)
		}
		if item.Idx == vf.lastSentEndIdx {
			sentOpen = false
			if len(sent) > 0 {
				vf.corpusSize += int64(len(sent))
				branches := findPathsToRoot(sent, vf.parentIdx, vf.deprelIdx)
				for _, b := range branches {
					vf.freqs.ImportSentence(b)
				}
			}
		}
		return true
	})
}

func (vf *Searcher) ProcToken(tk *vertigo.Token, line int, err error) error {
	vf.prevTokens.Append(tk)
	vf.lastTokenIdx = tk.Idx
	if vf.foundNewSent {
		vf.lastSentStartIdx = tk.Idx
		vf.foundNewSent = false
	}
	return nil
}

func (vf *Searcher) ProcStruct(st *vertigo.Structure, line int, err error) error {
	if st.Name == "s" {
		vf.lastSentEndIdx = vf.lastTokenIdx
		vf.analyzeLastSent()
		vf.foundNewSent = true
	}
	return nil
}

func (vf *Searcher) ProcStructClose(st *vertigo.StructureClose, line int, err error) error {
	return nil
}

func (vf *Searcher) ImportedCorpusSize() int64 {
	return vf.corpusSize
}

func NewSearcher(
	size int,
	lemmaIdx, parentAttrIdx, deprelAttrIdx int,
	freqs FreqsCollector,
) *Searcher {
	return &Searcher{
		prevTokens: collections.NewCircularList[*vertigo.Token](size),
		lemmaIdx:   lemmaIdx,
		parentIdx:  parentAttrIdx,
		deprelIdx:  deprelAttrIdx,
		freqs:      freqs,
	}
}
