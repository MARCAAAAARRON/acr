package scheduler

import (
	"sort"

	"acr/internal/retrieval"
)

const charsPerToken = 4

// EstimateTokens gives a rough token count for a string.
func EstimateTokens(s string) int {
	return len(s) / charsPerToken
}

// Schedule picks the highest-scoring chunks that fit within budgetTokens.
func Schedule(chunks []retrieval.Chunk, budgetTokens int) []retrieval.Chunk {
	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].Score > chunks[j].Score
	})

	var selected []retrieval.Chunk
	used := 0

	for _, c := range chunks {
		cost := EstimateTokens(c.Content)
		if used+cost <= budgetTokens {
			selected = append(selected, c)
			used += cost
		}
	}

	return selected
}