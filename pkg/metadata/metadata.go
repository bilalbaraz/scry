package metadata

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type DB struct {
	Path string
}

type FileRecord struct {
	Path  string
	Hash  string
	MTime int64
	Size  int64
}

type ChunkRecord struct {
	ID        string
	FilePath  string
	StartLine int
	EndLine   int
	Hash      string
	Content   string
}

type TermRecord struct {
	Term    string
	ChunkID string
	TF      int
}

type ChunkView struct {
	ID        string
	FilePath  string
	StartLine int
	EndLine   int
	Content   string
}

type TermHit struct {
	ChunkID string
	TF      int
}

type Stats struct {
	Files  int
	Chunks int
	Terms  int
}

func Open(path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	if _, err := exec.LookPath("sqlite3"); err != nil {
		return nil, fmt.Errorf("sqlite3 not found: %w", err)
	}
	db := &DB{Path: path}
	if err := db.init(); err != nil {
		return nil, err
	}
	return db, nil
}

func (d *DB) init() error {
	schema := `
CREATE TABLE IF NOT EXISTS files (
  path TEXT PRIMARY KEY,
  hash TEXT NOT NULL,
  mtime INTEGER NOT NULL,
  size INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS chunks (
  id TEXT PRIMARY KEY,
  file_path TEXT NOT NULL,
  start_line INTEGER NOT NULL,
  end_line INTEGER NOT NULL,
  hash TEXT NOT NULL,
  content TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS terms (
  term TEXT NOT NULL,
  chunk_id TEXT NOT NULL,
  tf INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS terms_term_idx ON terms(term);
CREATE INDEX IF NOT EXISTS chunks_file_idx ON chunks(file_path);
`
	return d.runScript(schema)
}

func (d *DB) GetFile(path string) (FileRecord, bool, error) {
	query := fmt.Sprintf("SELECT path, hash, mtime, size FROM files WHERE path = %s;", sqlQuote(path))
	out, err := d.runQuery(query)
	if err != nil {
		return FileRecord{}, false, err
	}
	lines := splitLines(out)
	if len(lines) == 0 {
		return FileRecord{}, false, nil
	}
	fields := strings.Split(lines[0], "\t")
	if len(fields) != 4 {
		return FileRecord{}, false, fmt.Errorf("unexpected columns for files")
	}
	return FileRecord{
		Path:  fields[0],
		Hash:  fields[1],
		MTime: parseInt64(fields[2]),
		Size:  parseInt64(fields[3]),
	}, true, nil
}

func (d *DB) ListFiles() ([]string, error) {
	out, err := d.runQuery("SELECT path FROM files;")
	if err != nil {
		return nil, err
	}
	lines := splitLines(out)
	return lines, nil
}

func (d *DB) GetChunk(id string) (ChunkView, bool, error) {
	query := fmt.Sprintf("SELECT id, file_path, start_line, end_line, content FROM chunks WHERE id = %s;", sqlQuote(id))
	out, err := d.runQuery(query)
	if err != nil {
		return ChunkView{}, false, err
	}
	lines := splitLines(out)
	if len(lines) == 0 {
		return ChunkView{}, false, nil
	}
	fields := strings.Split(lines[0], "\t")
	if len(fields) != 5 {
		return ChunkView{}, false, fmt.Errorf("unexpected columns for chunks")
	}
	return ChunkView{
		ID:        fields[0],
		FilePath:  fields[1],
		StartLine: int(parseInt64(fields[2])),
		EndLine:   int(parseInt64(fields[3])),
		Content:   fields[4],
	}, true, nil
}

func (d *DB) GetChunksByIDs(ids []string) ([]ChunkView, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var quoted []string
	for _, id := range ids {
		quoted = append(quoted, sqlQuote(id))
	}
	query := fmt.Sprintf("SELECT id, file_path, start_line, end_line, content FROM chunks WHERE id IN (%s);", strings.Join(quoted, ","))
	out, err := d.runQuery(query)
	if err != nil {
		return nil, err
	}
	lines := splitLines(out)
	var chunks []ChunkView
	for _, line := range lines {
		fields := strings.Split(line, "\t")
		if len(fields) != 5 {
			return nil, fmt.Errorf("unexpected columns for chunks")
		}
		chunks = append(chunks, ChunkView{
			ID:        fields[0],
			FilePath:  fields[1],
			StartLine: int(parseInt64(fields[2])),
			EndLine:   int(parseInt64(fields[3])),
			Content:   fields[4],
		})
	}
	return chunks, nil
}

func (d *DB) TermHits(term string) ([]TermHit, error) {
	query := fmt.Sprintf("SELECT chunk_id, tf FROM terms WHERE term = %s;", sqlQuote(term))
	out, err := d.runQuery(query)
	if err != nil {
		return nil, err
	}
	lines := splitLines(out)
	var hits []TermHit
	for _, line := range lines {
		fields := strings.Split(line, "\t")
		if len(fields) != 2 {
			return nil, fmt.Errorf("unexpected columns for terms")
		}
		hits = append(hits, TermHit{
			ChunkID: fields[0],
			TF:      int(parseInt64(fields[1])),
		})
	}
	return hits, nil
}

func (d *DB) Stats() (Stats, error) {
	files, err := d.scalarInt("SELECT COUNT(*) FROM files;")
	if err != nil {
		return Stats{}, err
	}
	chunks, err := d.scalarInt("SELECT COUNT(*) FROM chunks;")
	if err != nil {
		return Stats{}, err
	}
	terms, err := d.scalarInt("SELECT COUNT(*) FROM terms;")
	if err != nil {
		return Stats{}, err
	}
	return Stats{Files: files, Chunks: chunks, Terms: terms}, nil
}

func (d *DB) DeleteFile(path string) error {
	script := fmt.Sprintf(`
BEGIN;
DELETE FROM terms WHERE chunk_id IN (SELECT id FROM chunks WHERE file_path = %s);
DELETE FROM chunks WHERE file_path = %s;
DELETE FROM files WHERE path = %s;
COMMIT;
`, sqlQuote(path), sqlQuote(path), sqlQuote(path))
	return d.runScript(script)
}

func (d *DB) ReplaceFileData(fr FileRecord, chunks []ChunkRecord, terms []TermRecord) error {
	var buf bytes.Buffer
	buf.WriteString("BEGIN;\n")
	fmt.Fprintf(&buf, "DELETE FROM terms WHERE chunk_id IN (SELECT id FROM chunks WHERE file_path = %s);\n", sqlQuote(fr.Path))
	fmt.Fprintf(&buf, "DELETE FROM chunks WHERE file_path = %s;\n", sqlQuote(fr.Path))
	fmt.Fprintf(&buf, "INSERT INTO files(path, hash, mtime, size) VALUES(%s, %s, %d, %d)\n", sqlQuote(fr.Path), sqlQuote(fr.Hash), fr.MTime, fr.Size)
	buf.WriteString("ON CONFLICT(path) DO UPDATE SET hash=excluded.hash, mtime=excluded.mtime, size=excluded.size;\n")
	for _, ch := range chunks {
		fmt.Fprintf(&buf, "INSERT INTO chunks(id, file_path, start_line, end_line, hash, content) VALUES(%s, %s, %d, %d, %s, %s);\n",
			sqlQuote(ch.ID), sqlQuote(ch.FilePath), ch.StartLine, ch.EndLine, sqlQuote(ch.Hash), sqlQuote(ch.Content))
	}
	for _, tr := range terms {
		fmt.Fprintf(&buf, "INSERT INTO terms(term, chunk_id, tf) VALUES(%s, %s, %d);\n",
			sqlQuote(tr.Term), sqlQuote(tr.ChunkID), tr.TF)
	}
	buf.WriteString("COMMIT;\n")
	return d.runScript(buf.String())
}

func (d *DB) runQuery(query string) (string, error) {
	cmd := exec.Command("sqlite3", "-batch", "-noheader", "-separator", "\t", d.Path, query)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("sqlite3 query: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func (d *DB) scalarInt(query string) (int, error) {
	out, err := d.runQuery(query)
	if err != nil {
		return 0, err
	}
	lines := splitLines(out)
	if len(lines) == 0 {
		return 0, nil
	}
	return int(parseInt64(lines[0])), nil
}

func (d *DB) runScript(script string) error {
	cmd := exec.Command("sqlite3", d.Path)
	cmd.Stdin = strings.NewReader(script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("sqlite3 script: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func sqlQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

func splitLines(s string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	return lines
}

func parseInt64(s string) int64 {
	var v int64
	_, err := fmt.Sscanf(s, "%d", &v)
	if err != nil && !errors.Is(err, os.ErrInvalid) {
		return 0
	}
	return v
}
