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
	"fmt"
	"os"

	"github.com/czcorpus/scollector/record"
	"github.com/czcorpus/scollector/storage"
	"github.com/rs/zerolog/log"
	"github.com/tomachalek/vertigo/v6"
)

type freqs struct {
	LemmaIdx     int
	PosIdx       int
	DeprelIdx    int
	TextTypeAttr string
	Single       map[record.GroupingKey]record.TokenFreq
	Double       map[record.GroupingKey]record.CollocFreq
	TTMapping    map[string]byte
}

func (f *freqs) newCollocFreq(token1, token2 *vertigo.Token, freq int, distance int) record.CollocFreq {
	return record.CollocFreq{
		Lemma1:  token1.PosAttrByIndex(f.LemmaIdx),
		PoS1:    record.ImportUDPoS(token1.PosAttrByIndex(f.PosIdx)),
		Deprel1: record.ImportUDDeprel(token1.PosAttrByIndex(f.DeprelIdx)),
		Lemma2:  token2.PosAttrByIndex(f.LemmaIdx),
		PoS2:    record.ImportUDPoS(token2.PosAttrByIndex(f.PosIdx)),
		Deprel2: record.ImportUDDeprel(token2.PosAttrByIndex(f.DeprelIdx)),
		TextType: record.TextType{
			Readable: token1.StructAttrs[f.TextTypeAttr],
			Raw:      f.TTMapping[token1.StructAttrs[f.TextTypeAttr]],
		},
		Freq:    freq,
		AVGDist: float64(distance),
	}
}

func (f *freqs) AddLemma(token *vertigo.Token, freq int) {
	newEntry := record.TokenFreq{
		Lemma:  token.PosAttrByIndex(f.LemmaIdx),
		PoS:    record.ImportUDPoS(token.PosAttrByIndex(f.PosIdx)),
		Deprel: record.ImportUDDeprel(token.PosAttrByIndex(f.DeprelIdx)),
		Freq:   freq,
		TextType: record.TextType{
			Readable: token.StructAttrs[f.TextTypeAttr],
			Raw:      f.TTMapping[token.StructAttrs[f.TextTypeAttr]],
		},
	}
	curr, ok := f.Single[newEntry.Key()]
	if !ok {
		curr = newEntry
	}
	curr.UpdateFreq(freq)
	f.Single[curr.Key()] = curr
}

func (f *freqs) validateTT(token *vertigo.Token) {
	_, ok := f.TTMapping[token.StructAttrs[f.TextTypeAttr]]
	if !ok {
		log.Warn().
			Str("confAttribute", f.TextTypeAttr).
			Str("sourceValue", token.StructAttrs[f.TextTypeAttr]).
			Msg("cannot map text type value")
	}
}

func (f *freqs) AddCooc(token1, token2 *vertigo.Token, freq int, distance int) {
	newEntry := f.newCollocFreq(token1, token2, 0, 0)
	entryKey := newEntry.Key()
	curr, ok := f.Double[entryKey]
	if !ok {
		curr = newEntry
	}
	curr.UpdateFreqAndDist(freq, distance)
	f.Double[entryKey] = curr
}

func (f *freqs) ImportSentence(sent []*vertigo.Token) {
	if len(sent) > 0 {
		f.validateTT(sent[0]) // just shows a warning in case of missing tt values
	}
	for i, tok := range sent {
		f.AddLemma(tok, 1)
		for j := max(0, i-2); j < min(i+2, len(sent)); j++ {
			if i == j {
				continue
			}
			f.AddCooc(tok, sent[j], 1, i-j)
		}
	}
}

func (f *freqs) PrintPreview() {
	i := 0
	for k, v := range f.Single {
		fmt.Fprintf(os.Stderr, "%s => %v\n", k, v)
		if i > 10 {
			break
		}
		i++
	}
	i = 0
	for k, v := range f.Double {
		fmt.Fprintf(os.Stderr, "%s => %v\n", k, v)
		if i > 10 {
			break
		}
		i++
	}
}

func (f *freqs) StoreToDb(db *storage.DB, minFreq int) (storage.ImportStats, error) {
	seq := storage.NewTokenIDSequence()
	return db.StoreData(seq, f.Single, f.Double, minFreq)
}

func NewFreqs(lemmaIdx, posIdx, deprelIdx int, ttAttr string, ttMapping map[string]byte) *freqs {
	return &freqs{
		LemmaIdx:     lemmaIdx,
		DeprelIdx:    deprelIdx,
		PosIdx:       posIdx,
		Single:       make(map[record.GroupingKey]record.TokenFreq),
		Double:       make(map[record.GroupingKey]record.CollocFreq),
		TextTypeAttr: ttAttr,
		TTMapping:    ttMapping,
	}
}

// --------------------------------

type nullFreqs struct {
	verbose      bool
	lemmaIdx     int
	posIdx       int
	deprelIdx    int
	textTypeAttr string
}

func (f *nullFreqs) AddLemma(lemma *vertigo.Token, freq int) {
}

func (f *nullFreqs) AddCooc(lemma1, lemma2 *vertigo.Token, freq int, distance int) {
	if f.verbose {
		fmt.Fprintf(
			os.Stderr,
			"addCooc(%s (%s, %s), %v (%s, %s), %d, %d)\n",
			lemma1.PosAttrByIndex(f.lemmaIdx), lemma1.PosAttrByIndex(f.posIdx), lemma1.PosAttrByIndex(f.deprelIdx),
			lemma2.PosAttrByIndex(f.lemmaIdx), lemma2.PosAttrByIndex(f.posIdx), lemma2.PosAttrByIndex(f.deprelIdx),
			freq, distance,
		)
	}
}

func (f *nullFreqs) ImportSentence(sent []*vertigo.Token) {
	for i, tok := range sent {
		f.AddLemma(tok, 1)
		for j := max(0, i-2); j < min(i+2, len(sent)); j++ {
			if i == j {
				continue
			}
			f.AddCooc(tok, sent[j], 1, i-j)
		}
	}
}

func (f *nullFreqs) PrintPreview() {
	fmt.Println("NullFreqs ...")
}

func (f *nullFreqs) StoreToDb(db *storage.DB, minFreq int) (storage.ImportStats, error) {
	return storage.ImportStats{}, nil
}

func NewNullFreqs(
	lemmaIdx int,
	posIdx int,
	deprelIdx int,
	verbose bool,
) *nullFreqs {
	return &nullFreqs{
		verbose:   verbose,
		lemmaIdx:  lemmaIdx,
		posIdx:    posIdx,
		deprelIdx: deprelIdx,
	}
}
