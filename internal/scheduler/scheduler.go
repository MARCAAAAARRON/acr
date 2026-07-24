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

// Decision records what happened to one chunk during scheduling.
type Decision struct {
	Chunk    retrieval.Chunk
	Tokens   int
	Included bool
}

// Schedule picks the highest-scoring chunks that fit within budgetTokens.
// It returns both the selected chunks and a full record of every decision made,
// so callers can inspect (and optionally print) why each chunk was kept or cut.
func Schedule(chunks []retrieval.Chunk, budgetTokens int) ([]retrieval.Chunk, []Decision) {
	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].Score > chunks[j].Score
	})

	var selected []retrieval.Chunk
	var decisions []Decision
	used := 0

	for _, c := range chunks {
		cost := EstimateTokens(c.Content)
		include := used+cost <= budgetTokens

		decisions = append(decisions, Decision{
			Chunk:    c,
			Tokens:   cost,
			Included: include,
		})

		if include {
			selected = append(selected, c)
			used += cost
		}
	}

	return selected, decisions
}	