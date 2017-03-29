package basic

import (
	"fmt"
	"io"
	"strings"

	"github.com/crgimenes/lex"
)

// TokenParsers array
var TokenParsers = []lex.TokenFunction{
	lex.Ident,
	lex.NewLine,
	lex.NotImplemented,
}

// Parse current code
func Parse() (err error) {

	code := `10 print "test"
20 goto 10`

	var lex lex.Lexer
	lex.MaxParseID = len(TokenParsers)

	err = lex.Run(strings.NewReader(code))
	if err != nil {
		if err != io.EOF {
			println(err.Error())
			return
		}
	}

	for _, t := range lex.Tokens {
		fmt.Printf("%v\t%q\n", t.Type, t.Literal)
	}

	return
}
