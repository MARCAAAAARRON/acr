package retrieval

import (
	"os"
	"path/filepath"
	"testing"
	"strings"

	"acr/internal/scanner"
)

func TestFilenameMatchScoresHigherThanBodyOnly(t *testing.T) {
	// set up a temporary directory with two fake files
	dir := t.TempDir()

	mainContent := "package main\nfunc main() {}\n" // barely mentions "main"
	os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainContent), 0644)

	otherContent := "package other\n// main main main main main main\n" // repeats "main" a lot
	os.WriteFile(filepath.Join(dir, "other.go"), []byte(otherContent), 0644)

	files := []scanner.FileInfo{
		{Path: "main.go", Ext: "go"},
		{Path: "other.go", Ext: "go"},
	}

	chunks, err := Retrieve(dir, files, "what does main.go do")
	if err != nil {
		t.Fatalf("Retrieve failed: %v", err)
	}

	var mainScore, otherScore int
	for _, c := range chunks {
		if c.Path == "main.go" {
			mainScore = c.Score
		}
		if c.Path == "other.go" {
			otherScore = c.Score
		}
	}

	if mainScore <= otherScore {
		t.Errorf("expected main.go (filename match) to outscore other.go (body-only matches), got main.go=%d other.go=%d", mainScore, otherScore)
	}
}

func TestDensityScoringFavorsShortRelevantFiles(t *testing.T) {
	dir := t.TempDir()

	// short file: the query word appears twice, in very little surrounding text
	shortContent := "widget widget\n"
	os.WriteFile(filepath.Join(dir, "short.go"), []byte(shortContent), 0644)

	// long file: same two mentions of "widget", but diluted with a lot of
	// unrelated filler text — a pure raw-count scorer would call this a tie
	longContent := "widget\n" + strings.Repeat("filler text unrelated to anything\n", 50) + "widget\n"
	os.WriteFile(filepath.Join(dir, "long.go"), []byte(longContent), 0644)

	files := []scanner.FileInfo{
		{Path: "short.go", Ext: "go"},
		{Path: "long.go", Ext: "go"},
	}

	// query deliberately doesn't match either filename, so only body scoring
	// (raw count + density) is being tested here
	chunks, err := Retrieve(dir, files, "widget")
	if err != nil {
		t.Fatalf("Retrieve failed: %v", err)
	}

	var shortScore, longScore int
	for _, c := range chunks {
		if c.Path == "short.go" {
			shortScore = c.Score
		}
		if c.Path == "long.go" {
			longScore = c.Score
		}
	}

	if shortScore <= longScore {
		t.Errorf("expected short.go (dense with matches) to outscore long.go (same raw count, diluted), got short.go=%d long.go=%d", shortScore, longScore)
	}
}