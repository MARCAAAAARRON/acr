package main

import (
	"fmt"
	"acr/internal/scanner"
	"acr/internal/retrieval"
	"acr/internal/scheduler"
	"os"
)

func main(){
	files, err := scanner.Scan(".")
	if err != nil{
		fmt.Println("scan error:", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d files\n", len(files))
	for _, f := range files{
		fmt.Println(f.Path)
	}

	chunks, err := retrieval.Retrieve(".", files, "scanner")
	if err != nil{
		fmt.Println("retrieval error:", err)
		os.Exit(1)
	}

	fmt.Printf("matched %d chunks\n", len(chunks))
	for _, c := range chunks{
		fmt.Printf("%s (score %d)\n", c.Path, c.Score)
	}

	selected := scheduler.Schedule(chunks, 500)
	fmt.Printf("scheduled %d of %d chunks within budget\n", len(selected), len(chunks))
	for _, c := range selected {
		fmt.Printf("%s (score %d, ~%d tokens)\n", c.Path, c.Score, scheduler.EstimateTokens(c.Content))
	}
}
