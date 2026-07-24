package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"acr/internal/gitctx"
	"acr/internal/llm"
	"acr/internal/prompt"
	"acr/internal/retrieval"
	"acr/internal/scanner"
	"acr/internal/scheduler"
)

func main() {
	model := flag.String("model", "qwen2.5-coder-7b-instruct", "which model to use")
	repo := flag.String("repo", ".", "path to the repo to scan")
	budget := flag.Int("budget", 4000, "token budget for scheduling context")
	explain := flag.Bool("explain", false, "show scheduler decisions")
	flag.Parse()

	// scan ONCE at startup, reused for every question this session
	files, err := scanner.Scan(*repo)
	if err != nil {
		fmt.Println("scan error:", err)
		os.Exit(1)
	}
	fmt.Printf("Scanned %d files in %s. Type your question, or 'exit' to quit.\n", len(files), *repo)

	input := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\n> ")
		if !input.Scan() {
			break // EOF (Ctrl+D / Ctrl+Z) ends the session
		}
		query := strings.TrimSpace(input.Text())

		if query == "" {
			continue
		}
		if query == "exit" || query == "quit" {
			break
		}

		// git status is re-checked each question, since files you're
		// touching can change mid-session
		modifiedFiles := gitctx.ModifiedFiles(*repo)

		chunks, err := retrieval.Retrieve(*repo, files, query, modifiedFiles)
		if err != nil {
			fmt.Println("retrieval error:", err)
			continue
		}

		selected, decisions := scheduler.Schedule(chunks, *budget)
		fmt.Printf("scheduled %d of %d chunks within budget\n", len(selected), len(chunks))

		if *explain {
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
		}

		fullPrompt := prompt.Build(selected, query)

		reply, err := llm.Ask(*model, fullPrompt)
		if err != nil {
			fmt.Println("llm error:", err)
			continue
		}
		fmt.Println("\n" + reply)
	}

	fmt.Println("Session ended.")
}