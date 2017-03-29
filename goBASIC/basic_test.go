package basic

import (
	"io"
	"testing"
)

func TestParse(t *testing.T) {
	err := Parse()
	if err != nil {
		if err != io.EOF {
			t.Fatal(err)
		}
	}
}
