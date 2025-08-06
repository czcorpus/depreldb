package scoll

import (
	"github.com/czcorpus/scollector/storage"
)

type Calculator struct {
	database *storage.DB
}

func FromDatabase(db *storage.DB) *Calculator {
	return &Calculator{db}
}

func (calc *Calculator) GetCollocations(lemma string, options ...func(opts *CalculationOptions)) ([]storage.Collocation, error) {
	var opts CalculationOptions
	for _, opt := range options {
		opt(&opts)
	}

	return calc.database.CalculateMeasures(
		lemma,
		opts.PoS,
		opts.TextType,
		opts.PrefixSearch,
		opts.Limit,
		opts.SortBy,
		opts.CollocateGroupByPos,
		opts.CollocateGroupByDeprel,
		opts.CollocateGroupByTextType,
	)
}
