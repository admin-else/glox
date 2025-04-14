package main

import (
	"html/template"
	"os"
)

type Arg struct {
	Name string
	Type string
}

type Node struct {
	Name string
	Args []Arg
}

func main() {
	nodes := []Node{
		{"BinaryExpr", []Arg{
			{"Left", "Expr"},
			{"Operator", "TokenType"},
			{"Right", "Expr"},
		}}, {"GroupingExpr", []Arg{
			{"Expr", "Expr"},
		}}, {"LiteralExpr", []Arg{
			{"Value", "any"},
		}}, {"UnaryExpr", []Arg{
			{"Operator", "TokenType"},
			{"Expr", "Expr"},
		}}}
	tmplSrc, err := os.ReadFile("../generateast/ast.go.tmpl")
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New("ast").Parse(string(tmplSrc))
	if err != nil {
		panic(err)
	}
	f, err := os.Create("ast.go")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, nodes)
	if err != nil {
		panic(err)
	}
}
