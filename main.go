package main

import (
	"flag"
	"fmt"
	"os"

	"acr/internal/gitctx"
	"acr/internal/llm"
	"acr/internal/prompt"
	"acr/internal/retrieval"
	"acr/internal/scanner"
	"acr/internal/scheduler"
)

func main() {
	query := flag.String("query", "", "the question to ask ACR")
	model := flag.String("model", "qwen2.5-coder-7b-instruct", "which model to use")
	repo := flag.String("repo", ".", "path to the repo to scan")
	budget := flag.Int("budget", 4000, "token budget for scheduling context")
	explain := flag.Bool("explain", false, "show scheduler decisions")
	flag.Parse()

	files, err := scanner.Scan(*repo)
	if err != nil {
		fmt.Println("scan error:", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d files\n", len(files))

	modifiedFiles := gitctx.ModifiedFiles(*repo)
	if len(modifiedFiles) > 0 {
		fmt.Printf("Detected %d modified/staged/untracked files via git\n", len(modifiedFiles))
	}

	chunks, err := retrieval.Retrieve(*repo, files, *query, modifiedFiles)
	if err != nil {
		fmt.Println("retrieval error:", err)
		os.Exit(1)
	}
	fmt.Printf("matched %d chunks\n", len(chunks))

	selected, decisions := scheduler.Schedule(chunks, *budget)
	fmt.Printf("scheduled %d of %d chunks within budget\n", len(selected), len(chunks))

	if *explain {
		fmt.Println("\n--- scheduler decisions ---")
		used := 0
		for _, d := range decisions {
			status := "SKIPPED"
			if d.Included {
				status = "included"
				used += d.Tokens
			}
			fmt.Printf("%-50s score=%-4d tokens=%-6d %s (running total: %d/%d)\n",
				d.Chunk.Path, d.Chunk.Score, d.Tokens, status, used, *budget)
		}
		fmt.Println("---------------------------")
	}

	fullPrompt := prompt.Build(selected, *query)

	reply, err := llm.Ask(*model, fullPrompt)
	if err != nil {
		fmt.Println("llm error:", err)
		os.Exit(1)
	}
	fmt.Println("model says:", reply)
}