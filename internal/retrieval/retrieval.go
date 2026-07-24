package retrieval

import (
	"os"
	"path/filepath"
	"strings"

	"acr/internal/scanner"
)

// Chunk is one file's content, scored by relevance to a query.
type Chunk struct {
	Path    string
	Content string
	Score   int
}

const maxFileSize = 1 << 20 // 1MB safety cap

// isModified checks whether path appears in the modifiedFiles list.
// Paths are normalized to forward slashes before comparing, since git
// always reports paths with "/" regardless of OS, but filepath.Rel on
// Windows produces paths with "\" — without normalizing, these would
// never match even when they refer to the same file.
func isModified(path string, modifiedFiles []string) bool {
	normalizedPath := strings.ReplaceAll(path, "\\", "/")
	for _, m := range modifiedFiles {
		normalizedM := strings.ReplaceAll(m, "\\", "/")
		if normalizedM == normalizedPath {
			return true
		}
	}
	return false
}

// Retrieve scans file contents and scores them against the query's keywords.
// modifiedFiles is the list of git-modified/staged/untracked files (from
// gitctx.ModifiedFiles) -- files in this list get a strong relevance bonus,
// since "you're actively editing this right now" is a stronger signal than
// keyword matching alone. Pass nil if git awareness isn't available.
func Retrieve(root string, files []scanner.FileInfo, query string, modifiedFiles []string) ([]Chunk, error) {
	keywords := strings.Fields(strings.ToLower(query))

	var chunks []Chunk

	for _, f := range files {
		fullPath := root + string(os.PathSeparator) + f.Path

		info, err := os.Stat(fullPath)
		if err != nil {
			return nil, err
		}
		if info.Size() > maxFileSize {
			continue
		}

		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, err
		}
		content := string(data)
		lower := strings.ToLower(content)
		filename := strings.ToLower(filepath.Base(f.Path))

		rawCount := 0
		for _, kw := range keywords {
			rawCount += strings.Count(lower, kw)
		}

		density := 0
		if len(content) > 0 {
			density = int(float64(rawCount) / float64(len(content)) * 1000)
		}

		score := rawCount + density

		for _, kw := range keywords {
			if strings.Contains(filename, kw) {
				score += 50
			}
		}

		if isModified(f.Path, modifiedFiles) {
			score += 100
		}

		if score > 0 {
			chunks = append(chunks, Chunk{
				Path:    f.Path,
				Content: content,
				Score:   score,
			})
		}
	}

	return chunks, nil
}