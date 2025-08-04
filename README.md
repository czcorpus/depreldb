# Scollector

A high-performance Go-based collocation extraction and search tool for linguistic analysis. Scollector processes linguistic data to calculate statistical measures like T-Score and Log-Dice for finding meaningful syntactic collocations between lemmas.

## Features

- **Fast collocation search** using BadgerDB for efficient storage and retrieval
- **Statistical measures**: T-Score and Log-Dice calculations
- **Universal Dependencies support**: Full integration with UD POS tags and dependency relations
- **Flexible querying**: Filter by lemma, POS tags, dependency relations, and text types
- **Multiple output formats**: Tabular display or JSON output
- **Memory-efficient**: Optimized key encoding and protobuf serialization
- **Can be used as a library**

## Installation

### Prerequisites

- Go 1.23.4 or later
- Protocol Buffer compiler (`protoc`)
- `protoc-gen-go` plugin

#### Installing dependencies on Ubuntu/Debian:
```bash
# Install protobuf compiler
sudo apt-get install protobuf-compiler

# Install protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

### Building

```bash
# Clone the repository
git clone https://github.com/czcorpus/scollector
cd scollector

# Build the project
make all
```

This will:
1. Generate Go code from protobuf definitions
2. Build the `search` binary for querying databases
3. Build the `mkscolldb` binary for data import

Alternatively, build manually:
```bash
go build -o search ./cmd/search
go build -o mkscolldb ./cmd/mkscolldb
```

## Input Data Format

Scollector expects linguistic data in **vertical format**, where each token is on a separate line with tab-separated attributes. Sentences are separated by `<s>` structures.



### Import Profiles

Import profiles define the column structure of your vertical files. Predefined profiles include:

- **intercorp_v16ud**: InterCorp v16 with Universal Dependencies
- Add custom profiles in `storage/profiles.go`

Each profile specifies:
- Lemma column position
- POS tag column position
- Dependency relation column position
- Syntactic parent column position
- Text type mappings

## Usage

### Data Import

Before searching, you need to import linguistic data into the database using the `mkscolldb` tool:

```bash
./mkscolldb [options] [vert_path] [db_path]
```

#### Import Options

- `-import-profile=NAME` - Use predefined corpus profile (e.g., "intercorp_v16ud")
- `-lemma-idx=2` - Column position of lemma in vertical file (default: 2)
- `-pos-idx=5` - Column position of POS tag (default: 5)
- `-parent-idx=12` - Column position of syntactic parent info (default: 12)
- `-deprel-idx=11` - Column position of dependency relation (default: 11)
- `-min-freq=20` - Minimal frequency of collocates to accept (default: 20)
- `-syntax-mode` - Enable syntactic variant extraction (default: true)
- `-verbose` - Print detailed activity information (default: true)
- `-log-level=info` - Set logging level (debug, info, warn, error)

#### Import Examples

```bash
# Import using predefined profile
./mkscolldb -import-profile intercorp_v16ud -min-freq 10 /path/to/corpus.vert /path/to/database.db

# Import with custom column positions
./mkscolldb -lemma-idx 1 -pos-idx 3 -min-freq 5 /path/to/corpus.vert /path/to/database.db

# Import from directory of vertical files
./mkscolldb -import-profile intercorp_v16ud /path/to/corpus/dir/ /path/to/database.db
```

### Basic Search

```bash
./scollsrch [options] [db_path] [lemma] [pos] [text_type]
```

### Command Line Options

- `-limit=10` - Maximum number of matching items to show (default: 10)
- `-sort-by=tscore` - Sorting measure: `tscore` or `ldice` (default: tscore)
- `-corpus-size=100000000` - Corpus size for statistical calculations (default: 100,000,000)
- `-collocate-group-by-pos` - Group collocates by their POS tags
- `-collocate-group-by-deprel` - Group collocates by their dependency relations
- `-collocate-group-by-tt` - Group collocates by their text type
- `-json-out` - Output results in JSON format instead of tabular format
- `-log-level=info` - Set logging level (debug, info, warn, error)

### Examples

```bash
# Basic search for collocations of "run"
./search /path/to/database.db run

# Search with POS filtering
./search /path/to/database.db run VERB

# Search with custom limits and sorting
./search -limit=20 -sort-by=ldice /path/to/database.db run VERB

# JSON output for programmatic processing
./search -json-out /path/to/database.db run VERB

# Group results by POS and dependency relations
./search -collocate-group-by-pos -collocate-group-by-deprel /path/to/database.db run
```

## Output Format

### Tabular Output (default)
```
registry    lemma    lemma props.    collocate    collocate props    T-Score    Log-Dice    mutual dist.
══════════════════════════════════════════════════════════════════════════════════════════════════════
-           run      (root, VERB)    fast         (amod, ADJ)        12.45      8.23        1.2
-           run      (root, VERB)    quickly      (advmod, ADV)      10.33      7.89        1.5
```

### JSON Output (`-json-out`)
```json
{"lemma":{"value":"run","pos":"VERB","deprel":"root"},"collocate":{"value":"fast","pos":"ADJ","deprel":"amod"},"logDice":8.23,"tScore":12.45,"mutualDist":1.2,"textType":""}
```

## Statistical Measures

### T-Score
Measures the confidence of word association:
```
T-Score = (F(x,y) - F(x)*F(y)/N) / √F(x,y)
```

### Log-Dice
Measures the strength of association between words:
```
Log-Dice = 14.0 + log₂(2*F(x,y)/(F(x)+F(y)))
```

Where:
- `F(x,y)` = frequency of co-occurrence
- `F(x)`, `F(y)` = individual word frequencies
- `N` = corpus size

## Database Schema

Scollector uses BadgerDB with an optimized key encoding scheme:

### Key Types
- **Metadata**: `0x01 + keyID` → metadata values (e.g., import profile)
- **Lemma to ID**: `0x02 + lemma` → `tokenID`
- **Reverse index**: `0x03 + tokenID` → `lemma`
- **Token frequency**: `0x04 + tokenID + pos + textType + deprel` → `TokenDBEntry`
- **Collocation frequency**: `0x05 + token1ID + pos1 + textType + deprel1 + token2ID + pos2 + deprel2` → `CollocDBEntry`

### Data Structures

#### TokenFreq
Single token frequency with linguistic metadata:
```go
type TokenFreq struct {
    Lemma    string
    PoS      UDPoS      // Universal Dependencies POS
    Deprel   UDDeprel   // Universal Dependencies relation
    Freq     int
    TextType TextType
}
```

#### CollocFreq
Collocation frequency between two tokens:
```go
type CollocFreq struct {
    Lemma1, Lemma2   string
    PoS1, PoS2       UDPoS
    Deprel1, Deprel2 UDDeprel
    Freq             int
    AVGDist          float32  // Average distance between tokens
    TextType         TextType
}
```

## Development

### Project Structure

```
scollector/
├── cmd/
│   ├── search/          # Search command-line interface
│   └── mkscolldb/       # Data import command-line interface
├── record/              # Data structures and encoding
├── storage/             # BadgerDB storage layer
├── scoll/               # Collocation calculation logic
├── pb/                  # Generated protobuf code
└── dataimport/          # Data import utilities
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./storage -v
go test ./record -v
```

### Key Components

- **Storage Layer** (`storage/`): BadgerDB wrapper with optimized read/write operations
- **Record Types** (`record/`): Core data structures and key encoding functions
- **Calculator** (`scoll/`): Statistical measure calculations and query processing
- **Search CLI** (`cmd/search/`): Command-line interface for database querying
- **Import CLI** (`cmd/mkscolldb/`): Command-line interface for data import
- **Data Import** (`dataimport/`): Utilities for processing vertical format linguistic data

## Universal Dependencies Support

Scollector fully supports Universal Dependencies v2 standards:

- **POS Tags**: All 17 universal POS tags (ADJ, ADP, ADV, AUX, etc.)
- **Dependency Relations**: Complete set of UD relations (nsubj, obj, amod, etc.)
- **Efficient Encoding**: Byte-level encoding for compact storage

