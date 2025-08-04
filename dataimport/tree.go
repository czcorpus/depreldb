package dataimport

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/czcorpus/cnc-gokit/collections"
	"github.com/rs/zerolog/log"
	"github.com/tomachalek/vertigo/v6"
)

type branch []*vertigo.Token

func (br branch) String() string {
	var buff strings.Builder
	for i, v := range br {
		if i > 0 {
			buff.WriteString(" -> ")
		}
		buff.WriteString(v.Word)
	}
	return buff.String()
}

func (br branch) PrintToWord2Vec(lemmaxIdx, deprelIdx int) {
	for _, tk := range br {
		fmt.Printf("%s_%s ", tk.PosAttrByIndex(lemmaxIdx), tk.PosAttrByIndex(deprelIdx))
	}
	fmt.Println()
}

func asBranch(sent []*vertigo.Token, parentAttrIdx int) branch {
	ans := make(branch, 0, len(sent)+2)
	for _, tk := range sent {
		for v := range strings.SplitSeq(tk.PosAttrByIndex(parentAttrIdx), "|") {
			tmp := *tk
			if parentAttrIdx >= len(tmp.Attrs) {
				log.Error().
					Int("numColumns", len(tmp.Attrs)+1).
					Int("colIdx", parentAttrIdx).
					Msg("failed to convert sentence to branch - column not found")
				return branch{}
			}
			tmp.Attrs[parentAttrIdx-1] = v
			ans = append(ans, &tmp)
		}
	}
	return ans
}

func findLeaves(sent []*vertigo.Token, parentAttrIdx, deprelIdx int) []branch {
	syntSent := asBranch(sent, parentAttrIdx)
	allToks := collections.NewSet[int]()
	parents := collections.NewSet[int]()
	parentMap := make(map[int]int)
	for i, t := range syntSent {
		allToks.Add(i)
		par := t.PosAttrByIndex(parentAttrIdx)
		realPars := make([]int, 0, 1)
		if par != "" {
			pars := strings.Split(par, "|")
			for _, realPar := range pars {
				strings.TrimPrefix(par, "+")
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
	branches := make([]branch, 0, 10)
	for v := range allToks.Sub(parents).Iterate {
		path := make(branch, 0, 20)

		currNode := v
		visited := collections.NewSet[int]()
		var ok bool
		for currNode >= 0 {
			if visited.Contains(currNode) {
				log.Error().Msg("detected cycle, skipping")
				break
			}
			visited.Add(currNode)
			syntTok := syntSent[currNode]
			currNode, ok = parentMap[currNode]
			if !ok {
				log.Error().
					Int("elementIdx", v).
					Int("idx", syntSent[0].Idx).
					Msg("broken syntax tree path for element, skipping the path")
				continue
			}
			if syntTok.PosAttrByIndex(deprelIdx) != "punct" &&
				syntTok.PosAttrByIndex(deprelIdx) != "cc" {
				path = append(path, syntTok)
			}
		}
		branches = append(branches, path)
	}
	return branches
}
