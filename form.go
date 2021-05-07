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

func NewForm(args... string) *Form {
	f := &Form{
		prefix: "form",
	}

	if len(args) > 0 {
		f.code = args[0]
	}

	if len(args) > 1 {
		f.prefix = args[1]
	}

	return f
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

	root, err := Parse([]byte(f.code))
	if err != nil {
		return nil, err
	}

	f.root = root.(*ast.Root)

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

	return buf.String(), nil
}

// 根据 fields 更新 ExprAssign
func (f *Form) UpdateExprAssign(expr *ast.ExprAssign, fields []Field) {
	for _, field := range fields {
		// 变量名匹配
		isSameVar := field.Name == f.GetExprValue(expr.Var)

		if expr, ok := expr.Expr.(*ast.ExprArray); ok && isSameVar{
			for _, item := range expr.Items {
				key, val := f.GetItemExpr(item)
				if f.GetExprValue(key) == "value" {
					f.SetExprValue(val, field.Value)

					return
				}
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

	if exprArray, ok := expr.Expr.(*ast.ExprArray); ok {
		field.SetName(f.GetExprValue(expr.Var))

		for _, item := range exprArray.Items {
			field.Set(f.GetItemValues(item))
		}
	}

	return field
}

func (f *Form) ParseCode(code string) ([]Field, error) {
	f.code = code
	f.root = nil

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
