package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jcleira/kata19-word-chains/pkg/words"
)

func main() {
	startWord := flag.String("start", "dog", "Start word for the chain")
	endWord := flag.String("end", "cat", "Endword for the chain")
	flag.Parse()

	if *startWord == "" {
		fmt.Fprint(os.Stderr, "no start word provided, use -start=_word_\n")
		os.Exit(1)
	}

	if *endWord == "" {
		fmt.Fprint(os.Stderr, "no start word provided, use -start=_word_\n")
		os.Exit(1)
	}

	dictionaryFile, err := os.Open("./dictionary.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening the dictionary file: %v\n", err)
		os.Exit(1)
	}
	defer dictionaryFile.Close()

	fmt.Println("Loading dictionary file...")

	client, err := words.NewClient(dictionaryFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading the dictionary file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Searching chains...")

	chain, err := client.GetChain(*startWord, *endWord)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	if chain == nil {
		return
	}

	fmt.Println("Best chain found:")
	for _, word := range chain {
		fmt.Println(word)
	}
}
