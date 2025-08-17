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
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/czcorpus/cnc-gokit/logging"
	"github.com/czcorpus/scollector/scoll"
	"github.com/czcorpus/scollector/storage"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type srchCommand struct {
	lemma    string
	pos      string
	textType string
}

func evalREPLCommand(cmd string) srchCommand {
	items := strings.Split(strings.TrimSpace(cmd), " ")
	ans := srchCommand{lemma: items[0]}
	if len(items) > 1 && items[1] != "-" {
		ans.pos = items[1]
	}
	if len(items) > 2 && items[2] != "-" {
		ans.textType = items[2]
	}
	return ans
}

func main() {
	limit := flag.Int("limit", 10, "max num. of matching items to show")
	sortBy := flag.String("sort-by", "rrf", "sorting measure (either tscore or ldice)")
	collGroupByPos := flag.Bool("collocate-group-by-pos", false, "if set, then collocates will be split by their PoS")
	groupByDeprel := flag.Bool("group-by-deprel", false, "if set, then collocates will be split by their Deprel variants")
	collGroupByTT := flag.Bool("collocate-group-by-tt", false, "if set, then collocates will be split by their text type (registry)")
	predefinedSearch := flag.String("predefined-search", "", "use predefined search (modifiers-of, nouns-modified-by, verbs-subject, verbs-object)")
	jsonOut := flag.Bool("json-out", false, "if set then JSON format will be used to print results")
	logLevel := flag.String("log-level", "info", "set log level (debug, info, warn, error)")
	repl := flag.Bool("repl", false, "if set, then the search will run in an infinite read-eval-print loop (until Ctrl+C is pressed)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "search - search for collocations of a provided lemma\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n  %s [options] [db_path] [lemma]\n\t", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	logging.SetupLogging(logging.LoggingConf{
		Level: logging.LogLevel(*logLevel),
	})

	db, err := storage.OpenDB(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: ", err)
		os.Exit(1)
	}
	gbPos := scoll.WithNOP()
	if *collGroupByPos {
		gbPos = scoll.WithCollocateGroupByPos()
	}
	gbDeprel := scoll.WithNOP()
	if *groupByDeprel {
		gbDeprel = scoll.WithGroupByDeprel()
	}
	gbTT := scoll.WithNOP()
	if *collGroupByTT {
		gbTT = scoll.WithCollocateGroupByTextType()
	}

	gbPredSrch := scoll.WithNOP()
	if *predefinedSearch != "" {
		tmp := scoll.PredefinedSearch(*predefinedSearch)
		if !tmp.Validate() {
			fmt.Fprintf(os.Stderr, "uknown predefined search: %s", *predefinedSearch)
			os.Exit(1)
		}
		gbPredSrch = scoll.WithPredefinedSearch(tmp)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cmdReader := bufio.NewReader(os.Stdin)

	currCommand := srchCommand{
		lemma:    flag.Arg(1),
		pos:      flag.Arg(2),
		textType: flag.Arg(3),
	}

	for {

		if *repl && currCommand.lemma == "" {
			fmt.Println("\nenter a query (lemma [optional PoS] [optional TT]):")
			cmdChan := make(chan string, 1)
			go func() {
				cmd, _ := cmdReader.ReadString('\n')
				cmdChan <- cmd
			}()

			select {
			case <-ctx.Done():
				fmt.Println("\nExiting...")
				return
			case cmd := <-cmdChan:
				currCommand = evalREPLCommand(cmd)
			}
		}

		if currCommand.lemma == "" {
			fmt.Println("no query entered")
			continue
		}

		ans, err := scoll.FromDatabase(db).GetCollocations(
			currCommand.lemma,
			scoll.WithPoS(currCommand.pos),
			scoll.WithTextType(currCommand.textType),
			scoll.WithLimit(*limit),
			scoll.WithSortBy(storage.SortingMeasure(*sortBy)),
			gbPos,
			gbDeprel,
			gbTT,
			gbPredSrch,
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

			if len(ans) > 0 {
				headerFmt := color.New(color.FgGreen).SprintfFunc()
				columnFmt := color.New(color.FgHiMagenta).SprintfFunc()

				tbl := table.New(
					"registry",
					"lemma",
					"dependency",
					"collocate",
					"T-Score",
					"Log-Dice",
					"LMI",
					"LL",
					"RRF",
					"mutual dist.",
				)
				tbl.
					WithHeaderFormatter(headerFmt).
					WithFirstColumnFormatter(columnFmt).
					WithHeaderSeparatorRow('\u2550')
				for _, item := range ans {
					tbl.AddRow(item.AsRow()...)
				}
				tbl.Print()

			} else {
				fmt.Println("-- NO RESULT --")
			}
		}

		if *repl {
			currCommand = srchCommand{}

		} else {
			return
		}
	}
}
