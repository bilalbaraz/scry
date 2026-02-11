# AGENTS.md

You are working in this repository as a coding agent.

Principles
- Make small, reviewable commits.
- Prefer Go standard library; keep deps minimal.
- CLI must be stable and ergonomic.

Project goal
- Build `scry` CLI for a local-first codebase memory engine:
  - `scry index` creates/updates local indexes
  - `scry search` does hybrid search
  - `scry ask` does RAG with citations
  - `scry impact` shows change impact

Build & test
- `go test ./...`
- `golangci-lint run` (if configured)

Output expectations
- Always propose a plan before coding.
- For each change: explain files touched + why.
- Add minimal unit tests for non-trivial logic.
