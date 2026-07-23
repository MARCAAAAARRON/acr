package main

import (
	"fmt"
	// "acr/internal/scanner"
	// "acr/internal/retrieval"
	// "acr/internal/scheduler"
	"acr/internal/llm"
	"os"
)

func main(){
	reply, err := llm.Ask("bonsai-27b", "Say hello and confirm you're working.")
if err != nil {
	fmt.Println("llm error:", err)
	os.Exit(1)
}
fmt.Println("model says:", reply)
}
