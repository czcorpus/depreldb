package scoll

import "github.com/czcorpus/scollector/storage"

type CalculationOptions struct {
	PrefixSearch           bool
	PoS                    string
	TextType               string
	CorpusSize             int
	Limit                  int
	SortBy                 storage.SortingMeasure
	CollocateGroupByPos    bool
	CollocateGroupByDeprel bool
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

func WithCorpusSize(size int) func(opts *CalculationOptions) {
	return func(opts *CalculationOptions) {
		opts.CorpusSize = size
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
