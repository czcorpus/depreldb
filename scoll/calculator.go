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

func createPredefinedSearchFilter(srch PredefinedSearch) storage.SearchFilter {
	switch srch {
	case ModifiersOf:
		return func(pos1, deprel1, pos2, deprel2, textType byte) bool {
			// TODO implement the logic
			return false
		}
	case NounsModifiedBy:
		return func(pos1, deprel1, pos2, deprel2, textType byte) bool {
			// TODO implement the logic
			return false
		}
	case VerbsObject:
		return func(pos1, deprel1, pos2, deprel2, textType byte) bool {
			// TODO implement the logic
			return false
		}
	case VerbsSubject:
		return func(pos1, deprel1, pos2, deprel2, textType byte) bool {
			// TODO implement the logic
			return false
		}
	default:
		return nil
	}
}

func (calc *Calculator) GetCollocations(lemma string, options ...func(opts *CalculationOptions)) ([]storage.Collocation, error) {
	var opts CalculationOptions
	for _, opt := range options {
		opt(&opts)
	}
	customFilter := createPredefinedSearchFilter(opts.PredefinedSearch)
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
		customFilter,
	)
}
