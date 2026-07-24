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

// Retrieve scans file contents and scores them against the query's keywords.
func Retrieve(root string, files []scanner.FileInfo, query string) ([]Chunk, error) {
	keywords := strings.Fields(strings.ToLower(query))

	var chunks []Chunk

	for _, f := range files {
		fullPath := root + string(os.PathSeparator) + f.Path

		info, err := os.Stat(fullPath)
		if err != nil {
			return nil, err
		}
		if info.Size() > maxFileSize {
			continue // skip huge files, don't fail the whole run
		}

		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, err
		}
		content := string(data)
		lower := strings.ToLower(content)
		filename := strings.ToLower(filepath.Base(f.Path))

		// raw keyword count — rewards files that mention query terms a lot
		rawCount := 0
		for _, kw := range keywords {
			rawCount += strings.Count(lower, kw)
		}

		// density bonus — rewards files where matches make up a large
		// proportion of the content, so short-but-relevant files (like a
		// small .lua script) aren't buried by longer files that just have
		// more raw text for matches to accumulate in.
		density := 0
		if len(content) > 0 {
			density = int(float64(rawCount) / float64(len(content)) * 1000)
		}

		score := rawCount + density

		// filename match bonus — a query mentioning the file's own name
		// is a very strong relevance signal, stronger than body text.
		for _, kw := range keywords {
			if strings.Contains(filename, kw) {
				score += 50
			}
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