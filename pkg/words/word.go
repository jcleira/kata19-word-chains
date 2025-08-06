package words

// Word represents an unique real world word.
type Word struct {
	Term        string
	LinkedWords []*Word
	Score       int
}

// Link would seek linked words for the given word on a word's dictionary.
//
// A linkable word should match on char length and just have one letter different
// from the given word.
//
// TODO: Link should cleanup any previous Link call.
func (w *Word) Link(words []*Word) {
	for i := 0; i < len(words); i++ {
		if w.isLinkable(*words[i]) {
			w.LinkedWords = append(w.LinkedWords, words[i])
		}
	}
}

// CalcScore method will calculate a score for the given word based on the
// same chars at the same position compared with a target word.
//
// target: The give word as string to compare to
//
// Ex:
// Given foo & feo
// f o o
// 1 0 1 = 2 Score
// f e o
//
// Given foo & bar
// f o o
// 0 0 0 = 0 Score
// b a r
//
// Returns nothing.
func (w *Word) CalcScore(target string) {
	w.Score = 0
	for i := 0; i < len(w.Term); i++ {
		if w.Term[i] == target[i] {
			w.Score += 1
		}
	}
}

// isLinkable will check if a given word is linkable or not with the current
// word.
//
// A linkable word should match on char length and just have one letter different
// from the given word.
//
// Examples of linkable words for foo:
// - boo, foe, feo, too
//
// Examples of non-linkable words for foo:
// - foo, bar, faa, foobar, barfoo
//
// Returns true if the given word is linkable, false otherwise.
func (w *Word) isLinkable(word Word) bool {
	if w.Term == word.Term {
		return false
	}

	if len(w.Term) != len(word.Term) {
		return false
	}

	differences := 0
	for i := 0; i < len(w.Term); i++ {
		if w.Term[i] == word.Term[i] {
			continue
		}

		if differences == 1 {
			return false
		}

		differences++
	}

	return true

}
