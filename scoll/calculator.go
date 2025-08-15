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

import (
	"github.com/czcorpus/scollector/record"
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
		return func(pos1 byte, deprel uint16, pos2 byte, textType byte, dist float64) bool {
			return dist > 0 && deprel == record.DeprelNmod && pos1 == record.PosNOUN
		}
	case NounsModifiedBy:
		return func(pos1 byte, deprel uint16, pos2 byte, textType byte, dist float64) bool {
			return dist < 0 && deprel == record.DeprelNmod && pos2 == record.PosNOUN
		}
	case VerbsObject:
		return func(pos1 byte, deprel uint16, pos2 byte, textType byte, dist float64) bool {
			return dist < 0 && deprel == record.DeprelNsubj && pos2 == record.PosVERB
		}
	case VerbsSubject:
		return func(pos1 byte, deprel uint16, pos2 byte, textType byte, dist float64) bool {
			return dist < 0 && (deprel == record.DeprelObj || deprel == record.DeprelIobj) && pos2 == record.PosVERB
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
		opts.LemmasAsHead,
		opts.Limit,
		opts.SortBy,
		opts.CollocateGroupByPos,
		opts.GroupByDeprel,
		opts.CollocateGroupByTextType,
		customFilter,
	)
}
