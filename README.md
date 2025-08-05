# Scollector

A high-performance Go-based collocation extraction and search tool for linguistic analysis. Scollector processes linguistic data to calculate statistical measures like T-Score and Log-Dice for finding meaningful syntactic collocations between lemmas.

## Features

- **Ultra-fast collocation search** using BadgerDB with optimized read-only configurations
- **High-performance storage**:
  - **memory-efficient** binary key encoding and optimized grouping algorithms
- **Statistical measures**: T-Score and Log-Dice calculations
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
git clone https://github.com/czcorpus/scollector
cd scollector

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

Scollector expects linguistic data in **vertical format**, where each token is on a separate line with tab-separated attributes. Sentences are separated by `<s>` structures with possible xml-like attributes.



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
- `-repl` - Run in interactive read-eval-print loop mode (exit with CTRL+C)
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

# Interactive REPL mode
./search -repl /path/to/database.db
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

Scollector uses BadgerDB with highly optimized binary encoding for maximum performance:

### Storage Format
- **Binary encoding**: Custom 5-byte format
- **Compact keys**: Binary keys instead of strings for efficient grouping operations
- **Read-optimized**: Large block cache (512MB) and index cache (256MB) for fast queries
- **Distance encoding**: 1-byte encoding for syntactic distances with 0.1 precision (-12.7 to +12.7)

### Key Types
- **Metadata**: `0x01 + keyID` → JSON metadata (import profile, corpus info)
- **Lemma to ID**: `0x02 + lemma` → `tokenID` (4 bytes)
- **Reverse index**: `0x03 + tokenID` → `lemma` (string)
- **Token frequency**: `0x04 + tokenID + pos + textType + deprel` → `freq` (4 bytes)
- **Collocation frequency**: `0x05 + [composite key]` → `freq + distance` (5 bytes)

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

## Performance Optimizations

Scollector implements several high-performance optimizations:

### Binary Encoding
- **5-byte collocation values**: `freq (4 bytes) + distance (1 byte)` vs 10-15 bytes with Protocol Buffers
- **4-byte token values**: Direct uint32 encoding vs Protocol Buffer overhead
- **Distance compression**: Float64 distances compressed to single byte with 0.1 precision

### Binary Keys
- **8-byte token keys**: `[TokenID:4][PoS:1][Deprel:1][TextType:1][padding:1]`
- **16-byte collocation keys**: Two token keys combined for full collocation grouping
- **Zero heap allocations**: Stack-allocated array keys vs string concatenation
- **50%+ performance improvement** in grouping operations

### Database Configuration
- **Read-only mode**: Eliminates write locks and journal overhead
- **Large caches**: 512MB block cache + 256MB index cache for hot data
- **Optimized BadgerDB settings**: 1GB value logs, minimal memtables
- **Intelligent caching**: Multi-level caching for frequently accessed tokens during result collection

## Development

### Project Structure

```
scollector/
├── cmd/
│   └── search/          # Search command-line interface with REPL mode
├── record/              # Data structures, binary encoding, and key generation
├── storage/             # BadgerDB storage layer with performance optimizations
├── scoll/               # Collocation calculation logic
└── dataimport/          # Data import utilities (if available)
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

- **Storage Layer** (`storage/`): BadgerDB wrapper with read-optimized configurations and binary encoding
- **Record Types** (`record/`): Binary data structures, custom encoding functions, and high-performance key generation
- **Calculator** (`scoll/`): Statistical measure calculations with optimized grouping algorithms
- **Search CLI** (`cmd/search/`): Command-line interface with interactive REPL mode and CTRL+C support
- **Data Import** (`dataimport/`): Utilities for processing vertical format linguistic data (if available)

## Universal Dependencies Support

Scollector fully supports Universal Dependencies v2 standards:

- **POS Tags**: All 17 universal POS tags (ADJ, ADP, ADV, AUX, etc.)
- **Dependency Relations**: Complete set of UD relations (nsubj, obj, amod, etc.)
- **Efficient Encoding**: Byte-level encoding for compact storage

