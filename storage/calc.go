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
	"fmt"
	"math"
	"sort"
)

const (
	rrfConstantD = 60.0
)

// SortByRRF orders items using Reciprocal Rank Fusion
// (https://plg.uwaterloo.ca/%7Egvcormac/cormacksigir09-rrf.pdf)
func SortByRRF(items []Collocation) {
	list1 := make([]Collocation, len(items))
	copy(list1, items)
	sort.Slice(list1, func(i, j int) bool {
		return list1[i].LogDice > list1[j].LogDice
	})

	list2 := make([]Collocation, len(items))
	copy(list2, items)
	sort.Slice(list2, func(i, j int) bool {
		return list2[i].LMI > list2[j].LMI
	})

	list3 := make([]Collocation, len(items))
	copy(list3, items)
	sort.Slice(list3, func(i, j int) bool {
		return list3[i].TScore > list3[j].TScore
	})

	list4 := make([]Collocation, len(items))
	copy(list4, items)
	sort.Slice(list4, func(i, j int) bool {
		return list4[i].LogLikelihood > list4[j].LogLikelihood
	})

	scores := make(map[string]float64)

	for i := range len(items) {
		scores[list1[i].Hash()] += 1.0 / float64((rrfConstantD + i))
		scores[list2[i].Hash()] += 1.0 / float64((rrfConstantD + i))
		scores[list3[i].Hash()] += 1.0 / float64((rrfConstantD + i))
		scores[list4[i].Hash()] += 1.0 / float64((rrfConstantD + i))
	}

	for i := range len(items) {
		items[i].RRFScore = scores[items[i].Hash()]
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].RRFScore > items[j].RRFScore
	})

}

/*
|     |  y  | !y    | total |
|  x  |  a  |  b    | a + b |
|  !x |  c  |  d    | c + d |
|     | a+c | b+d   |  n    |
*/

func LLScore(fxy, fx, fy uint32, n int64) float64 {
	a := float64(fxy)
	b := float64(fx - fxy)
	c := float64(fy - fxy)
	d := float64(n - int64(fx) - int64(fy) + int64(fxy))
	ans := 2 * (a*math.Log(a) + b*math.Log(b) + c*math.Log(c) + d*math.Log(d) -
		(a+b)*math.Log(a+b) - (a+c)*math.Log(a+c) -
		(b+d)*math.Log(b+d) - (c+d)*math.Log(c+d) +
		(a+b+c+d)*math.Log(a+b+c+d))
	if ans > 9005754207 {
		fmt.Printf("GIGA NUMBER: %v, fx: %v, fy: %v, fxy: %v\n", ans, fx, fy, fxy)
	}

	return ans
}
