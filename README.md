# scry

![GitHub Release](https://img.shields.io/github/v/release/bilalbaraz/scry)
![Codecov](https://img.shields.io/codecov/c/gh/bilalbaraz/scry?logo=codecov&label=codecov)

**Local-first codebase memory CLI.**

`scry` builds a local index of your repository and lets you search and ask questions with evidence-only snippets. It is designed to keep your code on your machine and remain simple, predictable, and hackable.

---

## What it does (demo-style)

You point `scry` at a repo, it indexes files, and then you can search or ask questions:

```
# Index the repo (incremental by default)
scry index

# Lexical search
scry search "indexing pipeline"

# Ask a question (extractive evidence-based)
scry ask "ignore nasil calisiyor" --k 2
```

`scry ask` does **not** generate summaries or reasoning. It returns the best matching evidence snippets with citations.

---

## Why scry exists (design philosophy)

- **Local-first by default**: no code leaves your machine.
- **Deterministic and inspectable**: lexical scoring + explicit rules.
- **Small dependencies**: standard library whenever possible.
- **CLI first**: stable, scriptable commands.

---

## Current capabilities (accurate MVP)

**Implemented**
- Local-first indexing with `.scry/` workspace
- Incremental indexing using file + chunk hashing
- Go + Markdown chunking
- Lexical search (TF-based inverted index)
- Extractive `scry ask` with evidence snippets
- Index status reporting

**Ask ranking improvements currently in place**
- Relevance filtering (query term presence)
- Path boosting and penalization
- Minimum term match preference (2+ terms if possible)
- Evidence trimming (top 1–2 snippets)
- Stricter “I don’t know” with hint

---

## Honest limitations (please read)

`scry` is early-stage and intentionally simple:

- **No vector embeddings**
- **No ML reranking**
- **No semantic reasoning**
- **No LLM summarization**
- `ask` is **lexical-only** and **extractive** (snippets only)

If you need semantic search or summarization, those are future milestones.

---

## Install

There is no packaged release yet. Build from source:

```
go build -o scry ./cmd/scry
```

---

## Usage

### Index

```
# Incremental index
./scry index

# Full rebuild
./scry index --clean

# JSON progress
./scry index --json
```

### Search

```
./scry search "scan rules" --limit 5
./scry search "ignore pattern" --json
```

### Ask (extractive evidence)

```
./scry ask "ignore nasil calisiyor" --k 2
./scry ask "scan rules" --k 4 --json
```

### Status

```
./scry status
./scry status --json
```

---

## Output format (ask)

Human output:

```
Answer:
Found 2 relevant evidence chunk(s).

Evidence:
[1] pkg/scan/scan.go:17-20
    type Scanner struct {
        Root    string
        Matcher *ignore.Matcher
    }
```

JSON output (line-delimited objects):

```
{"type":"answer","text":"Found 2 relevant evidence chunk(s)."}
{"type":"evidence","id":1,"snippet":"...","path":"pkg/scan/scan.go","start_line":17,"end_line":20}
{"type":"summary","k":2}
```

---

## Configuration

A placeholder config file `.memengine.yml` is supported but not yet parsed into a full schema. This will evolve.

---

## Repository structure

- `cmd/scry/` CLI entrypoints
- `pkg/` indexing, parsing, metadata, search
- `internal/query/ask/` ask ranking pipeline
- `docs/architecture.md` architecture notes

---

## Development

```
# tests
make test
# or
go test ./...
```

---

## Roadmap (high-level)

- Vector embeddings + hybrid search
- Reranking and semantic matching
- More language parsers
- Impact analysis (`scry impact`)

---

## License

MIT. See `LICENSE`.
