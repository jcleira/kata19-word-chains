package words

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// TODO: probably chains should not be []string but []Word
// Chain would keep the word chain between to words
//
// Example: for "foo" and "fee" => []string{"foo", "foe", fee}
type Chain []string

// Traverse represents the chain's search between two words.
//
// We would be using a channel (Results channel) in order to provide an unique
// point to return and process successfull chains. That chain's process will
// look for the shortest (fewer nodes chain) that will be store on the
// ShortestChain attribute.
type Traverse struct {
	StartWord     *Word
	EndWord       *Word
	Results       chan Chain
	ShortestChain Chain

	sync.Mutex
}

// Perform starts a Traverse search on the words with the given configuration.
//
// It will trigger the recursive traverse.step function that will traverse all
// the linked words from the start word looking for suitable chains.
//
// It will return a suitable Chain if the Traverse's Start and End words are
// connected, or an error otherwise.
func (t *Traverse) Perform() ([]string, error) {
	go t.collectResults()

	t.step(t.StartWord, Chain{})
	// TODO this sleep MUST be removed, actually its needed to let the concurrent
	// code to gather results before exiting.
	time.Sleep(1 * time.Second)

	t.Lock()
	defer t.Unlock()
	defer close(t.Results)

	if t.ShortestChain == nil {
		return nil, fmt.Errorf("error, no chain found between '%s' and '%s'",
			t.StartWord.Term, t.EndWord.Term)
	}

	return t.ShortestChain, nil
}

// step performs a Traversal step on a word that is long linked with the start
// word.
//
// step currently keeps track of the word in the chain, finalize the chain if
// the // step's word is the end's word or perform new steps for each step
// word's LinkedWords.
//
// There is a check that before creating a new step for a linked word to
// confirm that it is not already on the chain.
//
// step doesn't return anything, it does uses the Traverse.Result channel to
// provide results if needed.
func (t *Traverse) step(stepWord *Word, chain Chain) {
	chain = append(chain, stepWord.Term)

	if stepWord.Term == t.EndWord.Term {
		t.Results <- chain
		return
	}

	bestScore := math.MaxFloat64
	var bestLinkedWord *Word

	for i := 0; i < len(stepWord.LinkedWords); i++ {
		alreadyOnChain := false
		for _, word := range chain {
			if word == stepWord.LinkedWords[i].Term {
				alreadyOnChain = true
				break
			}
		}
		if alreadyOnChain {
			continue
		}

		linkedWordScore := getWordsScore(t.EndWord.Term, stepWord.LinkedWords[i].Term)
		if bestScore <= linkedWordScore {
			continue
		}

		bestScore = linkedWordScore
		bestLinkedWord = stepWord.LinkedWords[i]
	}

	if bestLinkedWord == nil {
		return
	}
	t.step(bestLinkedWord, chain)
}

func getWordsScore(target, current string) float64 {
	totalScore := .0
	for i := 0; i < len(current); i++ {
		totalScore = math.Abs(float64(current[i]) - float64(target[i]))
	}

	return totalScore
}

// collectResults is executed as goroutine and is an infinite loop that would
// collect and process Traverse.Result channel messages from the Traverse's
// steps.
//
// The infinite loop will end when the Traverse finalizes all chains.
func (t *Traverse) collectResults() {
	for {
		select {
		case chain, ok := <-t.Results:
			if !ok {
				break
			}

			t.Lock()
			if t.ShortestChain == nil {
				t.ShortestChain = chain
			} else if len(chain) < len(t.ShortestChain) {
				t.ShortestChain = chain
			}
			t.Unlock()
		}
	}
}
