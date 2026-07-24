package scheduler

import (
	"testing"

	"acr/internal/retrieval"
)

func TestMinScoreFloorExcludesLowRelevanceChunks(t *testing.T) {
	chunks := []retrieval.Chunk{
		{Path: "low.go", Content: "short content", Score: 5},   // below minScore (10)
		{Path: "high.go", Content: "short content", Score: 50}, // above minScore
	}

	selected, _ := Schedule(chunks, 1000) // generous budget, so only score should matter here

	var lowIncluded, highIncluded bool
	for _, c := range selected {
		if c.Path == "low.go" {
			lowIncluded = true
		}
		if c.Path == "high.go" {
			highIncluded = true
		}
	}

	if lowIncluded {
		t.Errorf("expected low.go (score 5, below minScore) to be excluded, but it was included")
	}
	if !highIncluded {
		t.Errorf("expected high.go (score 50, above minScore) to be included, but it was excluded")
	}
}