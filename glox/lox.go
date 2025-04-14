package glox

import (
	"fmt"
)

func Run(code string) {
	tokens, err := Scan(code)
	if err != nil {
		panic(err)
	}
	expr, errors := Parse(tokens)
	if len(errors) != 0 {
		for _, e := range errors {
			fmt.Println(e)
		}
		return
	}
	err = Interpret(expr)
	if err != nil {
		panic(err)
	}
}
