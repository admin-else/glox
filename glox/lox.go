package glox

import (
	"encoding/json"
	"fmt"
)

func Run(code string) {
	tokens, err := Scan(code)
	if err != nil {
		panic(err)
	}
	expr, errors := Parse(tokens)
	if len(errors) != 0 {
		panic(errors)
	}
	expr_json, _ := json.Marshal(expr)
	fmt.Printf("%v", string(expr_json))
}
