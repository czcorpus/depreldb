# Scollector

A high-performance Go-based collocation extraction and search tool for linguistic analysis. Scollector processes linguistic data to calculate statistical measures like T-Score and Log-Dice for finding meaningful collocations between lemmas.

## Features

- **Fast collocation search** using BadgerDB for efficient storage and retrieval
- **Statistical measures**: T-Score and Log-Dice calculations
- **Universal Dependencies support**: Full integration with UD POS tags and dependency relations
- **Flexible querying**: Filter by lemma, POS tags, dependency relations, and text types
- **Multiple output formats**: Tabular display or JSON output
- **Memory-efficient**: Optimized key encoding and protobuf serialization

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
2. Build the `search` binary

Alternatively, build manually:
```bash
go build -o search ./cmd/search
```

## Usage

### Basic Search

```bash
./search [options] [db_path] [lemma] [pos] [text_type]
```

### Command Line Options

- `-limit=10` - Maximum number of matching items to show (default: 10)
- `-sort-by=tscore` - Sorting measure: `tscore` or `ldice` (default: tscore)
- `-corpus-size=100000000` - Corpus size for statistical calculations (default: 100,000,000)
- `-collocate-group-by-pos` - Group collocates by their POS tags
- `-collocate-group-by-deprel` - Group collocates by their dependency relations
- `-json-out` - Output results in JSON format instead of tabular format

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
lemma    dep.+PoS    coll       dep.+PoS    T-Score    log-dice    mutual dist.
-------------------------------------------------------------
run      root+VERB   fast       amod+ADJ    12.45      8.23        1.2
run      root+VERB   quickly    advmod+ADV  10.33      7.89        1.5
```

### JSON Output (`-json-out`)
```json
{"lemma1":"run","pos1":"VERB","deprel1":"root","lemma2":"fast","pos2":"ADJ","deprel2":"amod","freq":150,"tscore":12.45,"logdice":8.23,"avgdist":1.2}
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
- **Lemma to ID**: `0x00 + lemma` → `tokenID`
- **Token frequency**: `0x01 + tokenID + pos + textType + deprel` → `TokenDBEntry`
- **Collocation frequency**: `0x02 + token1ID + pos1 + textType + deprel1 + token2ID + pos2 + deprel2` → `CollocDBEntry`
- **Reverse index**: `0x03 + tokenID` → `lemma`

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
├── cmd/search/          # Command-line interface
├── record/              # Data structures and encoding
├── storage/             # BadgerDB storage layer
├── scoll/               # Collocation calculation logic
├── pb/                  # Generated protobuf code
└── dataimport/          # Data import utilities (planned)
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
- **CLI** (`cmd/search/`): Command-line interface for database querying

## Universal Dependencies Support

Scollector fully supports Universal Dependencies v2 standards:

- **POS Tags**: All 17 universal POS tags (ADJ, ADP, ADV, AUX, etc.)
- **Dependency Relations**: Complete set of UD relations (nsubj, obj, amod, etc.)
- **Efficient Encoding**: Byte-level encoding for compact storage

