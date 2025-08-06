package words

import (
	"errors"
	"reflect"
	"testing"
)

func TestPerform(t *testing.T) {
	startWord := &Word{
		Term: "foo",
		LinkedWords: []*Word{
			&Word{
				Term: "foe",
				LinkedWords: []*Word{
					&Word{Term: "fee"},
				},
			},
		},
		Score: 1,
	}

	tests := []struct {
		Description   string
		StartWord     *Word
		EndWord       *Word
		ExpectedError error
		ExpectedChain []string
	}{
		{
			Description:   "when Traverse succesfully finds a path for the words",
			StartWord:     startWord,
			EndWord:       &Word{Term: "fee"},
			ExpectedChain: []string{"foo", "foe", "fee"},
		},
		{
			Description:   "when traverse fails to find a path between the words",
			StartWord:     startWord,
			EndWord:       &Word{Term: "feo"},
			ExpectedError: errors.New("error, no chain found between 'foo' and 'feo'"),
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			traverse := &Traverse{
				StartWord: test.StartWord,
				EndWord:   test.EndWord,
				Results:   make(chan Chain),
			}

			chain, err := traverse.Perform()

			if !reflect.DeepEqual(test.ExpectedError, err) {
				t.Fatalf("expected err to be %v, got %v", test.ExpectedError, err)
			}

			if !reflect.DeepEqual(test.ExpectedChain, chain) {
				t.Fatalf("expected chain to be %v, got %v", test.ExpectedChain, chain)
			}
		})
	}

}
