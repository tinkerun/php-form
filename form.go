package main

import (
	"bytes"
	"fmt"
	"github.com/z7zmey/php-parser/pkg/ast"
	"strings"
)

type Form struct {
	prefix string
}

func NewForm() *Form {
	return &Form{
		prefix: "tinkerun",
	}
}

func (f *Form) WithPrefix(s string) *Form {
	f.prefix = s
	return f
}

func (f *Form) Prefix() string {
	if strings.HasPrefix(f.prefix, "$") {
		return f.prefix
	}

	return fmt.Sprintf("$%s", f.prefix)
}

func (f Form) ParseInputs()  {
	
}

func (f *Form) ParseCode(code string) ([]Input, error) {
	root, err := Parse([]byte(code))
	if err != nil {
		return nil, err
	}

	stmtList := root.(*ast.Root).Stmts

	var res []Input

	for _, stmt := range stmtList {
		// 如果是表达式
		if expr, ok := stmt.(*ast.StmtExpression); ok {

			input := NewInput()

			assign := expr.Expr.(*ast.ExprAssign)

			matched := bytes.HasPrefix(f.ExprValue(assign.Var), []byte(f.Prefix()))

			if expr, ok := assign.Expr.(*ast.ExprArray); ok && matched {
				for _, item := range expr.Items {
					input.Set(f.ItemValues(item))
				}
			}

			if !input.IsEmpty() {
				res = append(res, *input)
			}
		}
	}

	return res, nil
}

func (f *Form) TrimQuotes()  {

}

func (f *Form) ItemValues(item ast.Vertex) ([]byte, []byte) {
	expr := item.(*ast.ExprArrayItem)
	return f.ExprValue(expr.Key), f.ExprValue(expr.Val)
}

func (f *Form) ExprValue(expr ast.Vertex) []byte {
	switch v := expr.(type) {
	case *ast.ExprVariable:
		return v.Name.(*ast.Identifier).Value
	case *ast.ScalarString:
		return v.Value[1:len(v.Value)-1]
	case *ast.ScalarLnumber:
		return v.Value
	case *ast.ScalarDnumber:
		return v.Value
	case *ast.ExprConstFetch:
		return v.Const.(*ast.Name).Parts[0].(*ast.NamePart).Value
	}

	return nil
}
