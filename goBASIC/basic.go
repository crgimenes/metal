package basic

import (
	"fmt"
	"io"

	"github.com/crgimenes/lex"
)

// Parse current code
func Parse() (err error) {

	var l lex.Lexer

	l, err = lex.Parse("test.bas")
	if err != nil {
		if err != io.EOF {
			println(err.Error())
			return
		}
	}

	for _, t := range l.Tokens {
		fmt.Printf("%v\t%q\n", t.Type, t.Literal)
	}

	return
}
