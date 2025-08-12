// Copyright 2025 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2025 Institute of the Czech National Corpus,
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

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/czcorpus/scollector/dataimport"
	"github.com/czcorpus/scollector/record"
	"github.com/rs/zerolog/log"

	"github.com/czcorpus/cnc-gokit/fs"
	"github.com/czcorpus/cnc-gokit/logging"
	"github.com/czcorpus/scollector/storage"
	"github.com/tomachalek/vertigo/v6"
)

func determineFilesToProc(path string) ([]string, error) {
	isDir, err := fs.IsDir(path)
	if err != nil {
		return []string{}, fmt.Errorf("failed to determine files to process: %w", err)
	}
	ans := make([]string, 0, 50)
	if isDir {
		entries, err := os.ReadDir(path)
		if err != nil {
			return []string{}, fmt.Errorf("failed to list directory contents: %w", err)
		}
		for _, entry := range entries {
			ans = append(ans, filepath.Join(path, entry.Name()))
		}

	} else {
		ans = append(ans, path)
	}
	return ans, nil
}

func runCommand(path, dbPath string, prof storage.Profile, minFreq int, verbose bool) {
	var db *storage.DB
	var err error

	var freqColl dataimport.FreqsCollector
	if dbPath != "" {
		freqColl = dataimport.NewFreqs(
			prof.LemmaIdx,
			prof.PosIdx,
			prof.DeprelIdx,
			prof.TextTypesAttr,
			prof.TextTypes,
		)
		db, err = storage.OpenDBIgnoreMetadata(dbPath, prof.TextTypes)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR: ", err)
			os.Exit(2)
		}

	} else {
		freqColl = dataimport.NewNullFreqs(prof.LemmaIdx, prof.PosIdx, prof.DeprelIdx, verbose)
	}
	proc := dataimport.NewSearcher(
		50, prof.LemmaIdx, prof.PosIdx, prof.ParentIdx, prof.DeprelIdx, freqColl,
	)
	ctx := context.Background()
	files, err := determineFilesToProc(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: ", err)
		os.Exit(2)
	}
	for _, vertFile := range files {
		pConf := vertigo.ParserConf{
			InputFilePath:         vertFile,
			Encoding:              "utf-8",
			StructAttrAccumulator: "comb",
			LogProgressEachNth:    100000,
		}
		fmt.Fprintf(
			os.Stderr,
			"Starting to extract syntax data from file (min freq.: %d) %s\n-------------------\n",
			minFreq, vertFile,
		)
		if parserErr := vertigo.ParseVerticalFile(ctx, &pConf, proc); parserErr != nil {
			fmt.Fprintln(os.Stderr, "ERROR: ", parserErr)
			os.Exit(3)
		}
	}
	freqColl.PrintPreview()

	if err := db.Clear(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to clear existing database: %s\n", err)
	}
	stats, err := freqColl.StoreToDb(db, minFreq)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: ", err)
		os.Exit(2)
	}

	metadata := storage.Metadata{
		CorpusSize:    proc.ImportedCorpusSize(),
		NumCollFreqs:  stats.NumCollFreqs,
		NumLemmaFreqs: stats.NumLemmaFreqs,
		NumLemmas:     stats.NumLemmas,
		ProfileName:   prof.Name,
		DeprelMap:     nil,
	}

	for _, v := range proc.CollectedDeprels() {
		record.UDDeprelMapping.Register(v)
	}
	metadata.DeprelMap = record.UDDeprelMapping.AsMap()
	if err := db.StoreMetadata(metadata); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: ", err)
		os.Exit(4)
	}

	log.Info().
		Int64("corpusSize", metadata.CorpusSize).
		Int("numCollFreqs", metadata.NumCollFreqs).
		Int("numLemmaFreqs", metadata.NumLemmaFreqs).
		Int("numLemmas", metadata.NumLemmas).
		Str("profileName", metadata.ProfileName).
		Msg("collected and stored dataset metadata")
	fmt.Fprintf(
		os.Stderr,
		"import stats - total lemmas: %d, num single lemma freqs: %d, num coll freqs: %d",
		stats.NumLemmas, stats.NumLemmaFreqs, stats.NumCollFreqs,
	)

	db.Close() // this is ok to be called on possible nil

}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "w2vprep - Prepare data for word2vec/wang2vec processing	.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n  %s [options] [vert_path] [db_path]\n\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(
			os.Stderr, "Typical usage:\n  %s -import-profile intercorp_v16ud -min-freq 10 [vert_path] [db_path]\n\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	lemmaIdx := flag.Int("lemma-idx", 2, "vertical file column position where lemma is located")
	posIdx := flag.Int("pos-idx", 5, "vertical file column position where PoS is located (overrides importProfile)")
	parentIdx := flag.Int("parent-idx", 12, "vertical file column position where syntactic parent info is stored (overrides importProfile)")
	deprelIdx := flag.Int("deprel-idx", 11, "vertical file column position where syntactic function is stored (overrides importProfile)")
	iProfile := flag.String("import-profile", "", "select a predefined lemma-idx, pos-idx etc. based on corpus name (e.g. intercorp_v16ud)")
	verbose := flag.Bool("verbose", true, "print more info about program activity")
	minFreq := flag.Int("min-freq", 20, "minimal freq. of collocates to be accepted")
	logLevel := flag.String("log-level", "info", "set log level (debug, info, warn, error)")
	flag.Parse()

	logging.SetupLogging(logging.LoggingConf{
		Level: logging.LogLevel(*logLevel),
	})

	var cprof storage.Profile
	if *iProfile != "" {
		cprof = storage.FindProfile(*iProfile)
		if cprof.IsZero() {
			fmt.Fprintf(os.Stderr, "import profile %s not found", *iProfile)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Using import profile %s\n", *iProfile)

	} else {
		cprof = storage.Profile{
			LemmaIdx:  *lemmaIdx,
			PosIdx:    *posIdx,
			ParentIdx: *parentIdx,
			DeprelIdx: *deprelIdx,
		}
	}
	runCommand(flag.Arg(0), flag.Arg(1), cprof, *minFreq, *verbose)

}
