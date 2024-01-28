package pala

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

// An example program showing the basic capabilities.
// Some things to take note of:
// 1. A line that starts with a variable ($a, $b, $c) is an assignment to that variable
// 2. An operator might also be provided with a context that can be mutated but remains invisible in the program syntax.
//
// The program does the following:
// $a is assigned the result of the min operator applied to the list with elements 2, 3, 4
// $b is assigned the result of the + operator applied to $a and 4
// $c is assigned the result of the * operator applied to $b and 7
// The echo operator is applied to $c
var program = "$a min [2 3 4]\n" +
	"$b + $a 4\n" +
	"$c * $b 7\n" +
	"echo $c"

func TestParser(t *testing.T) {
	// Create a new Language (optionally with a context)
	lang := NewLanguage[*context]()

	// Bind operators, these calls will panic if the functions are of the wrong type.
	lang.BindOperator("+", plus)
	lang.BindOperator("*", mul)
	lang.BindOperator("min", smallest)
	lang.BindOperator("echo", echo)
	// Bind literal parsers, this call will panic if the function is of the wrong type.
	lang.BindLiteralEvaluator(ParseInt)

	// Create a lexer from any source.
	lexer := NewLexer(strings.NewReader(program))

	// Combine lexer and language in a parser.
	p := NewParser(lexer, lang)

	// Run the parser, obtaining either the constructed program or a parsing error.
	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("expected program to be parsed:\n%s", err)
	}

	// Instantiate the context
	ctx := &context{}
	// Run the program with the context.
	prog.Run(ctx)

	// Further process the modified context
	for _, logLine := range ctx.Log {
		fmt.Println(logLine)
	}
}

// Below are the example functions for a very simple mini language.

type context struct {
	Log []string
}

func (c context) String() string {
	return strings.Join(c.Log, "\n")
}

func smallest(c *context, a []int) int {
	var log []string
	for _, num := range a {
		log = append(log, fmt.Sprintf("%d", num))
	}
	c.Log = append(c.Log, fmt.Sprintf("finding min of [%s]", strings.Join(log, ",")))
	n := math.MaxInt
	for _, num := range a {
		if num < n {
			n = num
		}
	}
	return n
}

func mul(c *context, a, b int) int {
	c.Log = append(c.Log, fmt.Sprintf("multiplied %d and %d", a, b))
	return a * b
}

func plus(c *context, a, b int) int {
	c.Log = append(c.Log, fmt.Sprintf("added %d and %d", a, b))
	return a + b
}

func neg(c *context, a int) int {
	c.Log = append(c.Log, fmt.Sprintf("negated %d", a))
	return -a
}

func debug(c *context) {
	c.Log = append(c.Log, "debug")
}

func echo(s any) {
	fmt.Printf("%+v\n", s)
}

func shortest(lists [][]int) []int {
	m := math.MaxInt
	var res []int
	for _, list := range lists {
		if len(list) < m {
			m = len(list)
			res = list
		}
	}
	return res
}
