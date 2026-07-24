package prompt

import (
	"strings"

	"acr/internal/retrieval"
)

// Build assembles a full prompt string from selected chunks and the user's query.
func Build(chunks []retrieval.Chunk, query string) string {
	var b strings.Builder

	b.WriteString("You are a helpful coding assistant. Use the following file contents to answer the question.\n\n")

	for _, c := range chunks {
		b.WriteString("--- file: ")
		b.WriteString(c.Path)
		b.WriteString(" ---\n")
		b.WriteString(c.Content)
		b.WriteString("\n\n")
	}

	b.WriteString("Question: ")
	b.WriteString(query)

	return b.String()
}