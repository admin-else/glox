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

type Ast struct {
	AstType string
	Nodes   []Node
}

func main() {
	exprs := []Ast{
		{"Expr", []Node{
			{"BinaryExpr", []Arg{
				{"Left", "Expr"},
				{"Operator", "Token"},
				{"Right", "Expr"},
			}}, {"GroupingExpr", []Arg{
				{"Expr", "Expr"},
			}}, {"LiteralExpr", []Arg{
				{"Value", "any"},
			}}, {"UnaryExpr", []Arg{
				{"Operator", "Token"},
				{"Expr", "Expr"},
			}},
		}},
		{"Stmt", []Node{
			{"ExprStmt", []Arg{
				{"Expr", "Expr"},
			}}, {"PrintStmt", []Arg{
				{"Expr", "Expr"},
			}}, {"VarDecl", []Arg{
				{"Name", "Token"},
				{"Initializer", "Expr"},
			}},
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
	err = tmpl.Execute(f, exprs)
	if err != nil {
		panic(err)
	}
}
