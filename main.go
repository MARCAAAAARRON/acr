package main

import (
	"fmt"
	"os"

	"acr/internal/llm"
	"acr/internal/prompt"
	"acr/internal/retrieval"
	"acr/internal/scanner"
	"acr/internal/scheduler"
)

func main() {
	// 1. scan the repo
	files, err := scanner.Scan(".")
	if err != nil {
		fmt.Println("scan error:", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d files\n", len(files))

	// 2. retrieve chunks relevant to the query
	query := "What does the scanner package do?"
	chunks, err := retrieval.Retrieve(".", files, query)
	if err != nil {
		fmt.Println("retrieval error:", err)
		os.Exit(1)
	}
	fmt.Printf("matched %d chunks\n", len(chunks))

	// 3. schedule chunks within a token budget
	selected := scheduler.Schedule(chunks, 500)
	fmt.Printf("scheduled %d of %d chunks within budget\n", len(selected), len(chunks))

	// 4. build the full prompt
	fullPrompt := prompt.Build(selected, query)

	// 5. ask the model
	reply, err := llm.Ask("bonsai-27b", fullPrompt)
	if err != nil {
		fmt.Println("llm error:", err)
		os.Exit(1)
	}
	fmt.Println("model says:", reply)
}