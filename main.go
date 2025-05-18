package main

import (
	"flag"
	"glox/glox"
	"log"
	"os"
)

func main() {
	fileArg := flag.String("file", "code.lox", "Execute a lox file")
	replArg := flag.Bool("repl", false, "open a repl ")
	flag.Parse()

	if *replArg {
		glox.Repl()
	} else {
		b, err := os.ReadFile(*fileArg)
		if err != nil {
			log.Fatalln("Error occured while reading file", *fileArg, ":", err)
		}
		code := string(b)
		glox.Run(code)
	}
}
