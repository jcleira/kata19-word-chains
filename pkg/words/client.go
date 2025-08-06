package words

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Client is words package core struct, define & containsall the needed
// settings to perform Word operations.
//
// This "Client" idea came to my mind from the standard library http.Client
// struct. http.Client is a struct that provides a single entrypoint for the
// main operations on the http package. I'm used to that pattern for my packages
// and words.Client will provide all word's features needed for the Kata.
type Client struct {
	Words []*Word
}

// NewClient loads a feed of words, an initialize them.
//
// words: An io.Reader with a words dictionary, word per line.
//
// If there is no words on the provided words io.Reader it will return an error
//
// Returns a reference to the created Client or an error if any.
func NewClient(words io.Reader) (*Client, error) {
	client := &Client{
		Words: []*Word{},
	}

	scanner := bufio.NewScanner(words)
	for scanner.Scan() {
		// TODO - Word would need its own constructor
		word := scanner.Text()

		if len(word) != 3 {
			continue
		}

		client.Words = append(client.Words, &Word{Term: strings.ToLower(word)})
	}

	if len(client.Words) == 0 {
		return nil, errors.New("error, no words found on the given input")
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	for i := 0; i < len(client.Words); i++ {
		go client.Words[i].Link(client.Words)
	}

	return client, nil
}

// GetChain will perform a chain lookup from start to end words.
//
// It will return an []string with the word chain if it does find the path
// between those two words, otherwise it will return one of the following
// errors:
//
// - The start word is not found on the dictionary.
// - The end word is not found on the dictionary.
// - There is no chain between the two words.
func (c *Client) GetChain(start, end string) ([]string, error) {
	var startWord, endWord *Word

	for i := 0; i < len(c.Words); i++ {
		if c.Words[i].Term == start {
			startWord = c.Words[i]
			break
		}
	}

	if startWord == nil {
		return nil, fmt.Errorf("error, start word '%s' not found", start)
	}

	for i := 0; i < len(c.Words); i++ {
		if c.Words[i].Term == end {
			endWord = c.Words[i]
			break
		}
	}

	if endWord == nil {
		return nil, fmt.Errorf("error, end word '%s' not found", end)
	}

	startWord.CalcScore(endWord.Term)

	traverse := &Traverse{
		StartWord: startWord,
		EndWord:   endWord,
		Results:   make(chan Chain),
	}

	return traverse.Perform()
}
