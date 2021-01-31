package main

import "testing"

func TestParse(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"0+1", "1"},
		{"1+1", "2"},
		{"3+4", "7"},
		{"9+9", "18"},
	}

	for _, tt := range tests {
		i := interpreter{tt.input, 0, token{}}
		if res := i.Parse(); res != tt.want {
			t.Errorf("%v; want %v, got %v", tt.input, tt.want, res)
		}
	}
}
