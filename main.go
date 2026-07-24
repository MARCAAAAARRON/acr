package main

import (
	"fmt"
	"os"
	"flag"

	"acr/internal/llm"
	"acr/internal/prompt"
	"acr/internal/retrieval"
	"acr/internal/scanner"
	"acr/internal/scheduler"
)

func main() {
	// 1. scan the repo

	query := flag.String("query", "", "the question to ask ACR")
  model := flag.String("model", "qwen2.5-coder-7b-instruct", "which model to use")
  repo := flag.String("repo", ".", "path to the repo to scan")
	budget := flag.Int("budget", 4000, "token budget for scheduling context")
  flag.Parse()

	files, err := scanner.Scan(*repo)
	if err != nil {
		fmt.Println("scan error:", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d files\n", len(files))

	// 2. retrieve chunks relevant to the query
	chunks, err := retrieval.Retrieve(*repo, files, *query)
	if err != nil {
		fmt.Println("retrieval error:", err)
		os.Exit(1)
	}
	fmt.Printf("matched %d chunks\n", len(chunks))

	// 3. schedule chunks within a token budget
	selected := scheduler.Schedule(chunks, *budget)
	fmt.Printf("scheduled %d of %d chunks within budget\n", len(selected), len(chunks))

	// 4. build the full prompt
	fullPrompt := prompt.Build(selected, *query)

	// 5. ask the model
	reply, err := llm.Ask(*model, fullPrompt)
	if err != nil {
		fmt.Println("llm error:", err)
		os.Exit(1)
	}
	fmt.Println("model says:", reply)
}