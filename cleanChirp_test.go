package main

import (
	"testing"
)

func TestCleanChrip(t *testing.T) {
	type test struct {
		Input string
		Want  string
	}

	tests := []test{
		test{
			Input: "Hello there",
			Want:  "Hello there",
		},
		test{
			Input: "Kerfuffle This fornax and ShaRbeRt",
			Want:  "**** This **** and ****",
		},
		test{
			Input: "Kerfuffle!",
			Want:  "Kerfuffle!",
		},
	}

	for _, test := range tests {
		got := cleanChirpBody(test.Input)
		if got != test.Want {
			t.Errorf("Got: %v\nWant: %v\n", got, test.Want)
		}
	}
}
