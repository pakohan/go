package main

import "testing"

func TestPlural(t *testing.T) {
	testCases := map[string]string{
		"car":    "cars",
		"boy":    "boys",
		"theory": "theories",
		"":       "",
		"a":      "as",
		"ay":     "ays",
		"y":      "ys",
	}

	for input, expectedOutput := range testCases {
		output := plural(input)
		if output != expectedOutput {
			t.Errorf("expected plural of '%s' to be '%s', got '%s'", input, expectedOutput, output)
		}
	}
}
