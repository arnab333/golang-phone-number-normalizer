package main

import "testing"

// type normalizeTestCase struct {
// 	input string
// 	want  string
// }

func TestNormalize(t *testing.T) {
	// testCases := []normalizeTestCase{

	// }
	testCases := []struct {
		input string
		want  string
	}{
		{input: "1234567890", want: "1234567890"},
		{input: "123 456 7891", want: "1234567891"},
		{input: "(123) 456 7892", want: "1234567892"},
		{input: "(123) 456-7893", want: "1234567893"},
		{input: "123-456-7894", want: "1234567894"},
		{input: "(123)456-7892", want: "1234567892"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := normalize(tc.input)
			if actual != tc.want {
				t.Errorf("got %s; want %s", actual, tc.want)
			}
		})
	}
}
