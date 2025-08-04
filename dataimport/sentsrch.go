package dataimport

import (
	"github.com/czcorpus/cnc-gokit/collections"
	"github.com/czcorpus/scollector/storage"
	"github.com/tomachalek/vertigo/v6"
)

type FreqsCollector interface {
	AddLemma(lemma *vertigo.Token, freq int)
	AddCooc(lemma1, lemma2 *vertigo.Token, freq int, distance int)
	ImportSentence(sent []*vertigo.Token)
	PrintPreview()
	StoreToDb(db *storage.DB, minFreq int) error
}

// ----------------------------

type Searcher struct {
	prevTokens       *collections.CircularList[*vertigo.Token]
	lastTokenIdx     int
	lastSentStartIdx int
	lastSentEndIdx   int
	foundNewSent     bool
	lemmaIdx         int
	parentIdx        int
	deprelIdx        int
	freqs            FreqsCollector
}

func (vf *Searcher) analyzeLastSent() {
	var sentOpen bool
	sent := make([]*vertigo.Token, 0, vf.lastSentEndIdx-vf.lastSentStartIdx+1)
	vf.prevTokens.ForEach(func(i int, item *vertigo.Token) bool {
		if item.Idx == vf.lastSentStartIdx {
			//fmt.Printf("---- s[%d]>>>>\n", vf.lastSentStartIdx)
			sentOpen = true
		}
		if sentOpen {
			sent = append(sent, item)
		}
		if item.Idx == vf.lastSentEndIdx {
			sentOpen = false
			if len(sent) > 0 {
				branches := findLeaves(sent, vf.parentIdx, vf.deprelIdx)
				for _, b := range branches {
					//b.PrintToWord2Vec(vf.lemmaIdx, vf.deprelIdx)
					vf.freqs.ImportSentence(b)
				}
			}
		}
		return true
	})
}

func (vf *Searcher) ProcToken(tk *vertigo.Token, line int, err error) error {
	vf.prevTokens.Append(tk)
	vf.lastTokenIdx = tk.Idx
	if vf.foundNewSent {
		vf.lastSentStartIdx = tk.Idx
		vf.foundNewSent = false
	}
	return nil
}

func (vf *Searcher) ProcStruct(st *vertigo.Structure, line int, err error) error {
	if st.Name == "s" {
		vf.lastSentEndIdx = vf.lastTokenIdx
		vf.analyzeLastSent()
		vf.foundNewSent = true
	}
	return nil
}

func (vf *Searcher) ProcStructClose(st *vertigo.StructureClose, line int, err error) error {
	return nil
}

func NewSearcher(
	size int,
	lemmaIdx, parentAttrIdx, deprelAttrIdx int,
	freqs FreqsCollector,
) *Searcher {
	return &Searcher{
		prevTokens: collections.NewCircularList[*vertigo.Token](size),
		lemmaIdx:   lemmaIdx,
		parentIdx:  parentAttrIdx,
		deprelIdx:  deprelAttrIdx,
		freqs:      freqs,
	}
}
