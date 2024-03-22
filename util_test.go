package pala

import (
	"fmt"
	"math/big"
	"testing"
)

func TestParseInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		error    bool
	}{
		{
			"",
			0,
			true,
		},
		{
			"0",
			0,
			false,
		},
		{
			"1",
			1,
			false,
		},
		{
			"not a number",
			0,
			true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result, err := ParseInt(tt.input)
			if tt.error && err == nil {
				t.Fatalf("expected error but got %d", result)
			} else if result != tt.expected {
				t.Fatalf("epxected %d but got %d", result, tt.expected)
			}
		})
	}
}

func TestParseString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		error    bool
	}{
		{
			"",
			"",
			false,
		},
		{
			"a string",
			"a string",
			false,
		},
		{
			"1",
			"1",
			false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result, err := ParseString(tt.input)
			if tt.error && err == nil {
				t.Fatalf("expected error but got %s", result)
			} else if result != tt.expected {
				t.Fatalf("epxected %s but got %s", result, tt.expected)
			}
		})
	}
}

func TestParseQuotedString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		error    bool
	}{
		{
			"",
			"",
			true,
		},
		{
			"\"quoted string\"",
			"quoted string",
			false,
		},
		{
			"not quoted",
			"",
			true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result, err := ParseQuotedString(tt.input)
			if tt.error && err == nil {
				t.Fatalf("expected error but got %s", result)
			} else if result != tt.expected {
				t.Fatalf("epxected %s but got %s", result, tt.expected)
			}
		})
	}
}

func TestParseRational(t *testing.T) {
	tests := []struct {
		input    string
		expected *big.Rat
		error    bool
	}{
		{
			"",
			nil,
			true,
		},
		{
			" 1/1",
			nil,
			true,
		},
		{
			"1/1",
			big.NewRat(1, 1),
			false,
		},
		{
			"2/2",
			big.NewRat(1, 1),
			false,
		},
		{
			"4/5",
			big.NewRat(4, 5),
			false,
		},
		{
			"0/0",
			big.NewRat(0, 1),
			false,
		},
		{
			"40/0",
			big.NewRat(0, 1),
			false,
		},
		{
			"355/113",
			big.NewRat(355, 113),
			false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result, err := ParseRational(tt.input)
			if tt.error && err == nil {
				t.Fatalf("expected error but got %s", result)
			} else if !tt.error && err != nil {
				t.Fatalf("expected no error but got: %s", err)
			} else if !tt.error && tt.expected.Cmp(result) != 0 {
				t.Fatalf("epxected %s but got %s", result, tt.expected)
			}
		})
	}
}

func TestParseRandomInt(t *testing.T) {
	_, err := ParseRandomInt("invalid")
	if err == nil {
		t.Errorf("expected error")
	}

	_, err = ParseRandomInt("?int")
	if err != nil {
		t.Errorf("did not expect error")
	}
}
