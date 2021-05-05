package main

import (
	"bytes"
	"errors"
	"github.com/z7zmey/php-parser/pkg/ast"
	"github.com/z7zmey/php-parser/pkg/conf"
	phperrors "github.com/z7zmey/php-parser/pkg/errors"
	"github.com/z7zmey/php-parser/pkg/parser"
	"github.com/z7zmey/php-parser/pkg/version"
)

func Parse(code []byte) (ast.Vertex, error) {
	code = bytes.TrimSpace(code)
	if bytes.HasPrefix(code, []byte("<?")) || bytes.HasPrefix(code, []byte("<?php")) {
		return ParseCode(code)
	}
	return ParseStmtList(code)
}

func ParseCode(code []byte) (*ast.Root, error) {
	phpVersion, err := version.New("7.4")
	if err != nil {
		return nil, err
	}

	code = append(code, []byte("\n;")...)

	var parserErrors []*phperrors.Error
	rootNode, err := parser.Parse(code, conf.Config{
		Version: phpVersion,
		ErrorHandlerFunc: func(e *phperrors.Error) {
			parserErrors = append(parserErrors, e)
		},
	})
	if err != nil {
		return nil, err
	}
	if len(parserErrors) != 0 {
		return nil, errors.New(parserErrors[0].String())
	}

	res := rootNode.(*ast.Root)
	stmtLast := res.Stmts[len(res.Stmts)-1]

	// 处理最后一个空 stmt
	if _, ok := stmtLast.(*ast.StmtNop); ok {
		res.Stmts = res.Stmts[:len(res.Stmts)-1]
	}

	return res, nil
}

func ParseStmtList(code []byte) (*ast.Root, error) {
	code = append([]byte("<?php "), code...)
	return ParseCode(code)
}
