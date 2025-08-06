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

func WithPredefinedSearch(srch PredefinedSearch) func(opts *CalculationOptions) {
	return func(opts *CalculationOptions) {
		opts.PredefinedSearch = srch
	}
}
