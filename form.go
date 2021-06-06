package main

import (
	"bytes"
	"fmt"
	"github.com/z7zmey/php-parser/pkg/ast"
	"github.com/z7zmey/php-parser/pkg/visitor/printer"
	"strings"
)

type Form struct {
	prefix string
	code   string
	root   *ast.Root
}

func NewForm() *Form {
	return &Form{
		prefix: "form",
	}
}

func (f *Form) SetPrefix(prefix string) {
	f.prefix = prefix
}

func (f *Form) SetCode(code string) {
	f.code = code
	f.root = nil
}

func (f *Form) IsSnippetCode() bool {
	code := strings.TrimSpace(f.code)
	return !(strings.HasPrefix(code, "<?") || strings.HasPrefix(code, "<?php"))
}

func (f *Form) Prefix() string {
	if strings.HasPrefix(f.prefix, "$") {
		return f.prefix
	}

	return fmt.Sprintf("$%s", f.prefix)
}

func (f *Form) Root() (*ast.Root, error) {
	if f.root != nil {
		return f.root, nil
	}

	code := f.code
	// 如果是代码片段
	if f.IsSnippetCode() {
		code = fmt.Sprintf("<?php %s", f.code)
	}

	root, err := ParseCode([]byte(code))
	if err != nil {
		return nil, err
	}

	f.root = root
	return f.root, nil
}

// 根据输入的 fields 生成代码
func (f *Form) Stringify(fields []Field) (string, error) {
	root, err := f.Root()
	if err != nil {
		return "", err
	}

	stmtList := root.Stmts

	for _, stmt := range stmtList {
		// 如果是表达式
		if expr, ok := stmt.(*ast.StmtExpression); ok {
			if exprAssign, ok := expr.Expr.(*ast.ExprAssign); ok {
				if f.IsExprAssignMatched(exprAssign) {
					f.UpdateExprAssign(exprAssign, fields)
				}
			}
		}
	}

	buf := &bytes.Buffer{}

	pr := printer.NewPrinter(buf)
	root.Accept(pr)

	res := buf.String()
	if f.IsSnippetCode() {
		res = strings.Replace(res, "<?php ", "", 1)
	}

	return res, nil
}

// 根据 fields 更新 ExprAssign
func (f *Form) UpdateExprAssign(expr *ast.ExprAssign, fields []Field) {
	for _, field := range fields {
		// 变量名匹配
		if field.Name == f.GetExprValue(expr.Var) {
			switch e := expr.Expr.(type) {
			case *ast.ExprArray:
				for _, item := range e.Items {
					key, val := f.GetItemExpr(item)
					if f.GetExprValue(key) == "value" {
						f.SetExprValue(val, field.Value)

						return
					}
				}
			default:
				f.SetExprValue(e, field.Value)
			}
		}

	}
}

// 表达式的变量是否匹配前缀
func (f *Form) IsExprAssignMatched(expr *ast.ExprAssign) bool {
	return strings.HasPrefix(f.GetExprValue(expr.Var), f.Prefix())
}

// 根据 ExprAssign 获取 Input
func (f *Form) ParseExprAssign(expr *ast.ExprAssign) *Field {
	field := NewField()

	field.SetName(f.GetExprValue(expr.Var))

	switch e := expr.Expr.(type) {
	case *ast.ExprArray:
		for _, item := range e.Items {
			field.Set(f.GetItemValues(item))
		}
	case *ast.ScalarString, *ast.ScalarLnumber, *ast.ScalarDnumber, *ast.ExprConstFetch:
		field.SetValue(f.GetExprValue(e))
	}

	return field
}

func (f *Form) ParseCode(code string) ([]Field, error) {
	f.SetCode(code)
	return f.Parse()
}

// 解析代码获取 fields
func (f *Form) Parse() ([]Field, error) {
	root, err := f.Root()
	if err != nil {
		return nil, err
	}

	stmtList := root.Stmts

	var fields []Field

	for _, stmt := range stmtList {
		// 如果是表达式
		if expr, ok := stmt.(*ast.StmtExpression); ok {
			if exprAssign, ok := expr.Expr.(*ast.ExprAssign); ok {
				if f.IsExprAssignMatched(exprAssign) {
					field := f.ParseExprAssign(exprAssign)
					if !field.IsEmpty() {
						fields = append(fields, *field)
					}
				}
			}
		}
	}

	return fields, nil
}

// 获取 ExprArrayItem 的 kv 表达式
func (f *Form) GetItemExpr(item ast.Vertex) (ast.Vertex, ast.Vertex) {
	expr := item.(*ast.ExprArrayItem)
	return expr.Key, expr.Val
}

// 获取 ExprArray 中的数组数据
func (f *Form) GetExprArrayValue(expr *ast.ExprArray) interface{} {
	// 没有 key 的情况，返回数组
	if keyCheck, _ := f.GetItemExpr(expr.Items[0]); keyCheck == nil {
		var res []interface{}
		for _, item := range expr.Items {
			_, val := f.GetItemValues(item)
			if val != "" {
				res = append(res, val)
			}
		}

		return res
	}

	// 返回对象
	res := make(map[string]interface{})
	for _, item := range expr.Items {
		key, val := f.GetItemValues(item)
		if key != "" {
			res[key] = val
		}
	}

	return res
}

// 获取 ExprArrayItem 中的 kv 值
func (f *Form) GetItemValues(item ast.Vertex) (string, interface{}) {
	keyExpr, valExpr := f.GetItemExpr(item)

	key := f.GetExprValue(keyExpr)

	if arrayExpr, ok := valExpr.(*ast.ExprArray); ok {
		return key, f.GetExprArrayValue(arrayExpr)
	}

	return key, f.GetExprValue(valExpr)
}

// 获取表达式的值
func (f *Form) GetExprValue(expr ast.Vertex) string {
	var res []byte

	switch v := expr.(type) {
	case *ast.ExprVariable:
		res = v.Name.(*ast.Identifier).Value
	case *ast.ScalarString:
		res = v.Value[1 : len(v.Value)-1]
	case *ast.ScalarLnumber:
		res = v.Value
	case *ast.ScalarDnumber:
		res = v.Value
	case *ast.ExprConstFetch:
		res = v.Const.(*ast.Name).Parts[0].(*ast.NamePart).Value
	case *ast.ExprClosure, *ast.ExprArrowFunction:
		buf := &bytes.Buffer{}
		pr := printer.NewPrinter(buf)
		v.Accept(pr)
		res = buf.Bytes()
		res = bytes.TrimLeft(res, "<?php  ")
	}

	return string(res)
}

// 设置表达式的值
func (f Form) SetExprValue(expr ast.Vertex, value string) {
	data := []byte(value)

	switch v := expr.(type) {
	case *ast.ScalarString:
		v.Value = append([]byte{'\''}, data...)
		v.Value = append(v.Value, '\'')
		v.StringTkn.Value = v.Value
	case *ast.ScalarLnumber:
		v.Value = data
		v.NumberTkn.Value = data
	case *ast.ScalarDnumber:
		v.Value = data
		v.NumberTkn.Value = data
	case *ast.ExprConstFetch:
		v.Const.(*ast.Name).Parts[0].(*ast.NamePart).Value = data
		v.Const.(*ast.Name).Parts[0].(*ast.NamePart).StringTkn.Value = data
	}
}
