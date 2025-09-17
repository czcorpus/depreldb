# DeprelDB

A high-performance Go-based **dependency-based collocation extraction and search library** for linguistic analysis. DeprelDB processes linguistic data to calculate statistical measures like T-Score, Log-Dice, and LMI (Local Mutual Information) for finding meaningful syntactic collocations between lemmas.

## Features

- **Fast collocation search** using BadgerDB with optimized read-only configurations
- **High-performance storage**:
  - **memory-efficient** binary key encoding and optimized grouping algorithms
- **Statistical measures**: T-Score, Log-Dice, and LMI calculations with Reciprocal Rank Fusion (RRF) scoring
- **Universal Dependencies support**: Full integration with UD POS tags and dependency relations
- **Flexible querying**: Filter by lemma, POS tags, dependency relations, and text types
- **Multiple output formats**: Tabular display or JSON output
- **Large dataset optimized**: Handles multi-GB databases with intelligent caching
- **REPL mode**: Interactive query session with CTRL+C support
- **Can be used as a library**

## Installation

### Prerequisites

- Go 1.23.4 or later

### Building

```bash
# Clone the repository
git clone https://github.com/czcorpus/depreldb
cd depreldb

# Build the project
make all
```

This will build:
1. The `scollsrch` binary for querying databases
2. The `mkscolldb` binary for data import

Alternatively, build manually:
```bash
go build -o scollsrch ./cmd/search
```

## Input Data Format

DeprelDB expects linguistic data in **vertical format**, where each token is on a separate line with tab-separated attributes. Sentences are separated by `<s>` structures with possible xml-like attributes.



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
- Custom deprel values

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

- `-limit` - Maximum number of matching items to show (default: 10)
- `-sort-by` - Sorting measure: `tscore`, `ldice`, `lmi`, or `rrf` (default: rrf)
- `-collocate-group-by-pos` - Group collocates by their POS tags
- `-collocate-group-by-deprel` - Group collocates by their dependency relations
- `-collocate-group-by-tt` - Group collocates by their text type
- `-json-out` - Output results in JSON format instead of tabular format
- `-repl` - Run in interactive read-eval-print loop mode (exit with CTRL+C)
- `-log-level` - Set logging level (debug, info, warn, error, default = info)

### Examples

```bash
# Basic search for collocations of "run"
./search /path/to/database.db run

# Search with POS filtering
./search /path/to/database.db run VERB

# Search with custom limits and sorting
./search -limit=20 -sort-by=ldice /path/to/database.db run VERB

# Search using LMI measure
./search -sort-by=lmi /path/to/database.db run VERB

# Search using RRF (default) - combines all measures
./search -sort-by=rrf /path/to/database.db run VERB

# JSON output for programmatic processing
./search -json-out /path/to/database.db run VERB

# Group results by POS and dependency relations
./search -collocate-group-by-pos -collocate-group-by-deprel /path/to/database.db run

# Interactive REPL mode
./search -repl /path/to/database.db
```

## Output Format


### Tabular Output (default)
```
registry  lemma      lemma props.   collocate   collocate props  T-Score  Log-Dice  LMI     RRF Score  mutual dist.
════════  ═════      ════════════   ═════════   ═══════════════  ═══════  ════════  ══════  ═════════  ════════════
-         education  (nmod, -)      of          (-)               45.78    11.29     245.67  0.0821     1.10
-         education  (obj, -)       a           (-)               29.17    9.62      178.43  0.0734     1.10
-         education  (obj, -)       have        (-)               27.51    8.75      156.92  0.0687    -1.00
-         education  (nmod, -)      training    (-)               27.11    9.00      163.45  0.0701     2.00
```

### JSON Output (`-json-out`)
```json
{
  "lemma":{
    "value":"education",
    "pos":"",
    "deprel":"nmod"
  },
  "collocate":{
    "value":"of",
    "pos":"",
    "deprel":""
  },
  "logDice":11.28,
  "tScore":45.78,
  "lmi":245.67,
  "rrfScore":0.0821,
  "mutualDist":1.1,
  "textType":""
}
// etc...

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

### LMI (Local Mutual Information)

Measures pointwise mutual information weighted by co-occurrence frequency:
```
LMI = F(x,y) * log₂(N * F(x,y) / (F(x) * F(y)))
```

### RRF (Reciprocal Rank Fusion)

Combines rankings from T-Score, Log-Dice, and LMI using reciprocal rank fusion for better overall ranking:
```
RRF_score = Σ(1 / (k + rank_i))
```
where k=60 (RRF constant) and rank_i is the rank in each individual measure.

Where:
- `F(x,y)` = frequency of co-occurrence
- `F(x)`, `F(y)` = individual word frequencies
- `N` = corpus size

## Database Schema

DeprelDB uses BadgerDB with highly optimized binary encoding for maximum performance:

- **Binary encoding**: collocation entries encoded in 16 bytes long keys (9 bytes for single lemma frequencies)
- **Frequency and node distance encoded in DB values**
- - 4 bytes for **frequency**, 1 byte for **distance encoding** (0.1 precision; values from -12.7 to +12.7)
- **Efficient result grouping operations** - based on binary keys
- **Read-optimized**: Large block cache (512MB) and index cache (256MB) for fast queries


### Key Types
- **Metadata**: `0x01 + keyID` → JSON metadata (import profile, corpus info)
- **Lemma to ID**: `0x02 + lemma` → `tokenID`
- **Reverse index**: `0x03 + tokenID` → `lemma`
- **Token frequency**: `0x04 + tokenID + pos + textType + deprel` → `freq`
- **Collocation frequency**: `0x05 + [composite key]` → `freq + distance`



## Development

### Project Structure

```
depreldb/
├── cmd/
│   └── mkscolldb/       # An utility for importing corpus vertical files
│   └── search/          # Search command-line interface with REPL mode
├── record/              # Data structures, binary encoding, and key generation
├── storage/             # BadgerDB storage layer
├── scoll/               # High level interface for collocations search
└── dataimport/          # Data import logic
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./storage -v
go test ./record -v
```

