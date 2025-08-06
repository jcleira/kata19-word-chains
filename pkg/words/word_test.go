package words

import (
	"reflect"
	"testing"
)

func TestLink(t *testing.T) {
	tests := []struct {
		Description         string
		Word                Word
		Words               []*Word
		ExpectedLinkedWords []*Word
	}{
		{
			Description: "when words succesfully links to other words",
			Word:        Word{Term: "foo"},
			Words: []*Word{
				&Word{Term: "foo"},
				&Word{Term: "bar"},
				&Word{Term: "foe"},
				&Word{Term: "boo"},
				&Word{Term: "feo"},
				&Word{Term: "foobar"},
				&Word{Term: "barfoo"},
			},
			ExpectedLinkedWords: []*Word{
				&Word{Term: "foe"},
				&Word{Term: "boo"},
				&Word{Term: "feo"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			test.Word.Link(test.Words)

			if !reflect.DeepEqual(test.ExpectedLinkedWords, test.Word.LinkedWords) {
				t.Fatalf("expected LinkedWords to be %v, got %v",
					test.ExpectedLinkedWords,
					test.Word.LinkedWords)
			}
		})
	}
}

func TestCalcScore(t *testing.T) {
	tests := []struct {
		Description   string
		Word          Word
		EndWord       string
		ExpectedScore int
	}{
		{
			Description:   "when CalcScore succesfully calculates an score for two words",
			Word:          Word{Term: "foo"},
			EndWord:       "feo",
			ExpectedScore: 2,
		},
		{
			Description:   "when CalcScore succesfully calculates an score for two words",
			Word:          Word{Term: "foo"},
			EndWord:       "bar",
			ExpectedScore: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			test.Word.CalcScore(test.EndWord)

			if !reflect.DeepEqual(test.Word.Score, test.ExpectedScore) {
				t.Fatalf("expected ExpectedScore to be %v, got %v",
					test.ExpectedScore,
					test.Word.Score)
			}
		})
	}
}
