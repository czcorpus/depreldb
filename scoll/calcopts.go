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

package scoll

import "github.com/czcorpus/scollector/storage"

const (

	// ModifiersOf represents CQL chunk [p_lemma="team" & deprel="nmod" & upos="NOUN"]
	ModifiersOf PredefinedSearch = "modifiers-of"

	// NounsModifiedBy represents CQL chunk [lemma="team" & deprel="nmod" & p_upos="NOUN"]
	NounsModifiedBy PredefinedSearch = "nouns-modified-by"

	// VerbsSubject represents CQL chunk [lemma="team" & deprel="nsubj" & p_upos="VERB"]
	VerbsSubject PredefinedSearch = "verbs-subject"

	// VerbsObject represents CQL chunk [lemma="team" & deprel="obj|iobj" & p_upos="VERB"]
	VerbsObject PredefinedSearch = "verbs-object"
)

type PredefinedSearch string

type CalculationOptions struct {
	PrefixSearch             bool
	PoS                      string
	TextType                 string
	Limit                    int
	SortBy                   storage.SortingMeasure
	CollocateGroupByPos      bool
	CollocateGroupByDeprel   bool
	CollocateGroupByTextType bool
	LemmaGroupByDeprel       bool
	PredefinedSearch         PredefinedSearch
}

func WithPoS(pos string) func(opts *CalculationOptions) {
	return func(opts *CalculationOptions) {
		opts.PoS = pos
	}
}

func WithTextType(tt string) func(opts *CalculationOptions) {
	return func(opts *CalculationOptions) {
		opts.TextType = tt
	}
}

func WithLimit(lim int) func(opts *CalculationOptions) {
	return func(opts *CalculationOptions) {
		opts.Limit = lim
	}
}

func WithSortBy(measure storage.SortingMeasure) func(opts *CalculationOptions) {
	return func(opts *CalculationOptions) {
		opts.SortBy = measure
	}
}

func WithPrefixSearch() func(opts *CalculationOptions) {
	return func(opts *CalculationOptions) {
		opts.PrefixSearch = true
	}
}

func WithCollocateGroupByPos() func(opts *CalculationOptions) {
	return func(opts *CalculationOptions) {
		opts.CollocateGroupByPos = true
	}
}

func WithCollocateGroupByDeprel() func(opts *CalculationOptions) {
	return func(opts *CalculationOptions) {
		opts.CollocateGroupByDeprel = true
	}
}

func WithCollocateGroupByTextType() func(opts *CalculationOptions) {
	return func(opts *CalculationOptions) {
		opts.CollocateGroupByTextType = true
	}
}

func WithLemmaGroupByDeprel() func(opts *CalculationOptions) {
	return func(opts *CalculationOptions) {
		opts.LemmaGroupByDeprel = true
	}
}

func WithPredefinedSearch(srch PredefinedSearch) func(opts *CalculationOptions) {
	return func(opts *CalculationOptions) {
		opts.PredefinedSearch = srch
	}
}

// WithNOP is a convenience function which sets no option and
// can be used as an alternative to boolean With... functions
// with no argument.
func WithNOP() func(opts *CalculationOptions) {
	return func(opts *CalculationOptions) {
	}
}
