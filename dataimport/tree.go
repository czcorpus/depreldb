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
	"strconv"
	"strings"

	"github.com/czcorpus/cnc-gokit/collections"
	"github.com/rs/zerolog/log"
	"github.com/tomachalek/vertigo/v6"
)

type expandedSent []*vertigo.Token

func (br expandedSent) String() string {
	var buff strings.Builder
	for i, v := range br {
		if i > 0 {
			buff.WriteString(" -> ")
		}
		buff.WriteString(v.Word)
	}
	return buff.String()
}

func (br expandedSent) PrintToWord2Vec(lemmaxIdx, deprelIdx int) {
	for _, tk := range br {
		fmt.Printf("%s_%s ", tk.PosAttrByIndex(lemmaxIdx), tk.PosAttrByIndex(deprelIdx))
	}
	fmt.Println()
}

func asExpandedSent(sent []*vertigo.Token, parentAttrIdx int) expandedSent {
	ans := make(expandedSent, 0, len(sent)+2)
	for _, tk := range sent {
		for v := range strings.SplitSeq(tk.PosAttrByIndex(parentAttrIdx), "|") {
			tmp := *tk
			tmp.Attrs = make([]string, len(tmp.Attrs))
			copy(tmp.Attrs, tk.Attrs)
			if parentAttrIdx >= len(tmp.Attrs) {
				log.Error().
					Int("numColumns", len(tmp.Attrs)+1).
					Int("colIdx", parentAttrIdx).
					Msg("failed to convert sentence to branch - column not found")
				return expandedSent{}
			}
			tmp.Attrs[parentAttrIdx-1] = v
			ans = append(ans, &tmp)
		}
	}
	return ans
}

func isBlocklistedRel(rel string) bool {
	return rel == "punct" || rel == "cc" || strings.HasPrefix(rel, "det") || strings.HasPrefix(rel, "aux") ||
		rel == "cop" || rel == "mark" || strings.HasPrefix(rel, "expl") || rel == "discourse" ||
		rel == "goeswith" || rel == "reparandum" || rel == "orphan" || rel == "list" || rel == "vocative" ||
		rel == "dep"
}

func logCyclePath(path expandedSent, cycleToken *vertigo.Token, parentIdx int) {
	tmp := make([]string, len(path)+1)
	for i, v := range path {
		tmp[i] = fmt.Sprintf("%s (parent: %s)", v.Word, cycleToken.PosAttrByIndex(parentIdx))
	}
	tmp[len(tmp)-1] = fmt.Sprintf("%s (parent: %s)", cycleToken.Word, cycleToken.PosAttrByIndex(parentIdx))
	log.Error().
		Strs("cyclePath", tmp).
		Int("cycleElmIdx", cycleToken.Idx).
		Strs("parentMap", collections.SliceMap(path, func(v *vertigo.Token, i int) string { return v.Word })).
		Msg("detected cycle, skipping")
}

// --------

type visitedNode struct {
	isMultival bool
	idx        int
}

func (vn visitedNode) ID() string {
	return fmt.Sprintf("%d:%t", vn.idx, vn.isMultival)
}

func (vn visitedNode) valid() bool {
	return vn.idx > -1
}

func findPathsToRoot(
	sent []*vertigo.Token,
	lemmaIdx, posIdx, parentAttrIdx, deprelIdx int,
	deprelCollector *collections.Set[string],
) []expandedSent {
	syntSent := asExpandedSent(sent, parentAttrIdx)
	allToks := collections.NewSet[int]()
	parents := collections.NewSet[int]()
	parentMap := make(map[int]int)
	for i, t := range syntSent {
		allToks.Add(i)
		par := t.PosAttrByIndex(parentAttrIdx)
		realPars := make([]int, 0, 1)
		if par != "" {
			// we must deal with multivalues (val1|val2) which should
			// be split into two nodes
			for realPar := range strings.SplitSeq(par, "|") {
				par = strings.TrimPrefix(par, "+")
				iPar, err := strconv.Atoi(realPar)
				if err != nil {
					log.Error().Err(err).Str("value", par).Msg("failed to parse attribute 'parent', skipping")
					continue
				}
				realPars = append(realPars, iPar)
			}

		} else {
			realPars = append(realPars, -1)
		}
		for _, rp := range realPars {
			parents.Add(i + rp)

			if rp != 0 {
				parentMap[i] = i + rp

			} else {
				parentMap[i] = -1
			}
		}
	}
	branches := make([]expandedSent, 0, 10)
	for v := range allToks.Sub(parents).Iterate {
		path := make(expandedSent, 0, 20)

		currNode := visitedNode{
			idx:        v,
			isMultival: strings.Contains(syntSent[v].PosAttrByIndex(parentAttrIdx), "|"),
		}
		visited := collections.NewHSet[visitedNode]()
		for currNode.valid() {
			if visited.Contains(currNode) && !currNode.isMultival {
				logCyclePath(path, syntSent[currNode.idx], parentAttrIdx)
				break
			}
			visited.Add(currNode)
			syntTok := syntSent[currNode.idx]

			parentNode := visitedNode{}
			parentNodeIdx, ok := parentMap[currNode.idx]
			if !ok {
				log.Error().
					Int("elementIdx", v).
					Int("idx", syntSent[0].Idx).
					Msg("broken syntax tree path - unknown parent, taking partial path")
				break
			}
			parentNode.idx = parentNodeIdx
			if parentNode.valid() {
				currNode.isMultival = strings.Contains(syntSent[parentNode.idx].PosAttrByIndex(parentAttrIdx), "|")
			}

			if isBlocklistedRel(syntTok.PosAttrByIndex(deprelIdx)) {
				// NOP

			} else if parentNode.valid() && syntTok.PosAttrByIndex(posIdx) == "ADP" {
				if syntSent[parentNode.idx].PosAttrByIndex(deprelIdx) == "obl" {
					syntSent[parentNode.idx].Attrs[deprelIdx-1] = "obl:" + syntTok.PosAttrByIndex(lemmaIdx)
					deprelCollector.Add(syntSent[parentNode.idx].Attrs[deprelIdx-1])
					log.Debug().
						Str("word", syntSent[parentNode.idx].Word).
						Str("deprel", syntSent[parentNode.idx].Attrs[deprelIdx-1]).
						Msgf("merged ADP+case into parent obl:%s", syntTok.PosAttrByIndex(lemmaIdx))
				}

			} else {
				path = append(path, syntTok)
			}

			currNode = parentNode

		}
		branches = append(branches, path)
	}
	return branches
}
