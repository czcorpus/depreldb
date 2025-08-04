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

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/czcorpus/scollector/scoll"
	"github.com/czcorpus/scollector/storage"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func main() {
	limit := flag.Int("limit", 10, "max num. of matching items to show")
	sortBy := flag.String("sort-by", "tscore", "sorting measure (either tscore or ldice)")
	corpusSize := flag.Int("corpus-size", 100000000, "max num. of matching items to show")
	collGroupByPos := flag.Bool("collocate-group-by-pos", false, "if set, then collocates will be split by their PoS")
	collGroupByDeprel := flag.Bool("collocate-group-by-deprel", false, "if set, then collocates will be split by their Deprel value")
	collGroupByTT := flag.Bool("collocate-group-by-tt", false, "if set, then collocates will be split by their text type (registry)")
	jsonOut := flag.Bool("json-out", false, "if set then JSON format will be used to print results")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "search - search for collocations of a provided lemma\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n  %s [options] [db_path] [lemma]\n\t", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	db, err := storage.OpenDB(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: ", err)
		os.Exit(1)
	}
	gbPos := func(opts *scoll.CalculationOptions) {}
	if *collGroupByPos {
		gbPos = scoll.WithCollocateGroupByPos()
	}
	gbDeprel := func(opts *scoll.CalculationOptions) {}
	if *collGroupByDeprel {
		gbDeprel = scoll.WithCollocateGroupByDeprel()
	}
	gbTT := func(opts *scoll.CalculationOptions) {}
	if *collGroupByTT {
		gbTT = scoll.WithCollocateGroupByTextType()
	}

	ans, err := scoll.FromDatabase(db).GetCollocations(
		flag.Arg(1),
		scoll.WithPoS(flag.Arg(2)),
		scoll.WithCorpusSize(*corpusSize),
		scoll.WithTextType(flag.Arg(3)),
		scoll.WithLimit(*limit),
		scoll.WithSortBy(storage.SortingMeasure(*sortBy)),
		gbPos,
		gbDeprel,
		gbTT,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: ", err)
		os.Exit(1)
	}
	if *jsonOut {
		for _, item := range ans {
			out, err := json.Marshal(item)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to json-encode value: %s", err)
				os.Exit(1)
			}
			fmt.Println(string(out))
		}

	} else {
		fmt.Println()

		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()

		tbl := table.New("registry", "lemma", "deprel+PoS", "T-Score", "Log-Dice", "mutual dist.")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
		for _, item := range ans {
			tbl.AddRow(item.AsRow()...)
		}
		tbl.Print()
	}
}
