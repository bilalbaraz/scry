package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func FileHash(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func ChunkHash(fileHash string, startLine, endLine int, text string) string {
	h := sha256.New()
	fmt.Fprintf(h, "%s:%d:%d:\n", fileHash, startLine, endLine)
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}
