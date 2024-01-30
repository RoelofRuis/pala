package pala

import (
	"strings"
	"testing"
)

func Test_FailToParse(t *testing.T) {
	tests := []struct {
		name           string
		program        string
		expectedErrMsg string
	}{
		{
			"operator with wrong argument count",
			"+ 1",
			"operator + expected 2 operands but got 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lang := NewLanguage[*context]()
			lang.BindOperator("+", plus)
			lang.BindLiteralEvaluator(ParseInt)

			parser := NewParser(
				NewLexer(strings.NewReader(tt.program)),
				lang,
			)

			_, err := parser.Parse()
			if err == nil {
				t.Fatalf("expected program to fail to be parsed but it succeeded")
			}

			if err.Error() != tt.expectedErrMsg {
				t.Fatalf("expected error '%s' but got '%s'", tt.expectedErrMsg, err.Error())
			}
		})
	}
}
