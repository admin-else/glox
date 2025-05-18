package glox

import (
	"bufio"
	"fmt"
	"os"
)

func Repl() {
	i := &Interpreter{
		&Enviorment{
			values:    map[string]any{},
			enclosing: nil,
		},
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">>> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading line", err)
			os.Exit(1)
		}

		stmts, errs := ParseCode(input)
		if len(errs) != 0 {
			for _, err := range errs {
				fmt.Println(err)
			}
			os.Exit(1)
		}
		for _, stmt := range stmts {
			_, err := i.execute(stmt)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}
func Run(code string) {
	stmts, errors := ParseCode(code)
	if len(errors) != 0 {
		for _, e := range errors {
			fmt.Println(e)
		}
		return
	}
	err := Interpret(stmts)
	if err != nil {
		panic(err)
	}
}
