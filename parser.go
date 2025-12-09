// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"

	"golang.org/x/tools/txtar"
)

// ASTNode represents a node in the AST tree for display
type ASTNode struct {
	Label       string
	Children    []*ASTNode
	IndentLevel int
	Collapsed   bool
}

// ParseTxtar parses txtar content and returns AST nodes
func ParseTxtar(content string) ([]*ASTNode, error) {
	ar := txtar.Parse([]byte(content))

	var nodes []*ASTNode

	for _, file := range ar.Files {
		if !strings.HasSuffix(file.Name, ".go") {
			continue
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, file.Name, file.Data, parser.ParseComments)
		if err != nil {
			nodes = append(nodes, &ASTNode{
				Label:       fmt.Sprintf("%s (error: %v)", file.Name, err),
				IndentLevel: 1,
			})
			continue
		}

		fileNode := &ASTNode{
			Label:       fmt.Sprintf("File: %s", file.Name),
			IndentLevel: 1,
		}
		fileNode.Children = astToNodes(f, 2)
		nodes = append(nodes, fileNode)
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("no .go files found in txtar content")
	}

	return nodes, nil
}

// astToNodes converts an AST node to our display nodes
func astToNodes(node ast.Node, level int) []*ASTNode {
	if node == nil {
		return nil
	}

	var nodes []*ASTNode

	switch n := node.(type) {
	case *ast.File:
		// Package name
		nodes = append(nodes, &ASTNode{
			Label:       fmt.Sprintf("Package: %s", n.Name.Name),
			IndentLevel: level,
		})

		// Imports
		if len(n.Imports) > 0 {
			importsNode := &ASTNode{
				Label:       "Imports",
				IndentLevel: level,
			}
			for _, imp := range n.Imports {
				path := imp.Path.Value
				importsNode.Children = append(importsNode.Children, &ASTNode{
					Label:       fmt.Sprintf("Import: %s", path),
					IndentLevel: level + 1,
				})
			}
			nodes = append(nodes, importsNode)
		}

		// Declarations
		for _, decl := range n.Decls {
			nodes = append(nodes, declToNode(decl, level)...)
		}

	default:
		nodes = append(nodes, &ASTNode{
			Label:       fmt.Sprintf("%T", node),
			IndentLevel: level,
		})
	}

	return nodes
}

// declToNode converts a declaration to display nodes
func declToNode(decl ast.Decl, level int) []*ASTNode {
	var nodes []*ASTNode

	switch d := decl.(type) {
	case *ast.GenDecl:
		nodes = append(nodes, genDeclToNode(d, level)...)
	case *ast.FuncDecl:
		nodes = append(nodes, funcDeclToNode(d, level))
	default:
		nodes = append(nodes, &ASTNode{
			Label:       fmt.Sprintf("Decl: %T", decl),
			IndentLevel: level,
		})
	}

	return nodes
}

// genDeclToNode converts a general declaration to display nodes
func genDeclToNode(d *ast.GenDecl, level int) []*ASTNode {
	var nodes []*ASTNode

	switch d.Tok {
	case token.TYPE:
		for _, spec := range d.Specs {
			if ts, ok := spec.(*ast.TypeSpec); ok {
				typeNode := &ASTNode{
					Label:       fmt.Sprintf("Type: %s", ts.Name.Name),
					IndentLevel: level,
				}
				typeNode.Children = typeSpecToNodes(ts, level+1)
				nodes = append(nodes, typeNode)
			}
		}
	case token.CONST:
		constNode := &ASTNode{
			Label:       "Const",
			IndentLevel: level,
		}
		for _, spec := range d.Specs {
			if vs, ok := spec.(*ast.ValueSpec); ok {
				for _, name := range vs.Names {
					constNode.Children = append(constNode.Children, &ASTNode{
						Label:       fmt.Sprintf("Const: %s", name.Name),
						IndentLevel: level + 1,
					})
				}
			}
		}
		if len(constNode.Children) > 0 {
			nodes = append(nodes, constNode)
		}
	case token.VAR:
		varNode := &ASTNode{
			Label:       "Var",
			IndentLevel: level,
		}
		for _, spec := range d.Specs {
			if vs, ok := spec.(*ast.ValueSpec); ok {
				for _, name := range vs.Names {
					child := &ASTNode{
						Label:       fmt.Sprintf("Var: %s", name.Name),
						IndentLevel: level + 1,
					}
					if vs.Type != nil {
						child.Children = append(child.Children, &ASTNode{
							Label:       fmt.Sprintf("Type: %s", exprToString(vs.Type)),
							IndentLevel: level + 2,
						})
					}
					varNode.Children = append(varNode.Children, child)
				}
			}
		}
		if len(varNode.Children) > 0 {
			nodes = append(nodes, varNode)
		}
	}

	return nodes
}

// typeSpecToNodes converts a type specification to display nodes
func typeSpecToNodes(ts *ast.TypeSpec, level int) []*ASTNode {
	var nodes []*ASTNode

	switch t := ts.Type.(type) {
	case *ast.StructType:
		structNode := &ASTNode{
			Label:       "StructType",
			IndentLevel: level,
		}
		if t.Fields != nil {
			for _, field := range t.Fields.List {
				fieldNode := fieldToNode(field, level+1)
				structNode.Children = append(structNode.Children, fieldNode)
			}
		}
		nodes = append(nodes, structNode)

	case *ast.InterfaceType:
		ifaceNode := &ASTNode{
			Label:       "InterfaceType",
			IndentLevel: level,
		}
		if t.Methods != nil {
			for _, method := range t.Methods.List {
				methodNode := fieldToNode(method, level+1)
				ifaceNode.Children = append(ifaceNode.Children, methodNode)
			}
		}
		nodes = append(nodes, ifaceNode)

	default:
		nodes = append(nodes, &ASTNode{
			Label:       fmt.Sprintf("TypeExpr: %s", exprToString(ts.Type)),
			IndentLevel: level,
		})
	}

	return nodes
}

// fieldToNode converts a field to a display node
func fieldToNode(field *ast.Field, level int) *ASTNode {
	var name string
	if len(field.Names) > 0 {
		names := make([]string, len(field.Names))
		for i, n := range field.Names {
			names[i] = n.Name
		}
		name = strings.Join(names, ", ")
	} else {
		name = "(embedded)"
	}

	node := &ASTNode{
		Label:       fmt.Sprintf("Field: %s", name),
		IndentLevel: level,
	}

	node.Children = append(node.Children, &ASTNode{
		Label:       fmt.Sprintf("Type: %s", exprToString(field.Type)),
		IndentLevel: level + 1,
	})

	if field.Tag != nil {
		node.Children = append(node.Children, &ASTNode{
			Label:       fmt.Sprintf("Tag: %s", field.Tag.Value),
			IndentLevel: level + 1,
		})
	}

	return node
}

// funcDeclToNode converts a function declaration to a display node
func funcDeclToNode(f *ast.FuncDecl, level int) *ASTNode {
	var label string
	if f.Recv != nil && len(f.Recv.List) > 0 {
		recv := f.Recv.List[0]
		recvType := exprToString(recv.Type)
		label = fmt.Sprintf("Method: (%s) %s", recvType, f.Name.Name)
	} else {
		label = fmt.Sprintf("Func: %s", f.Name.Name)
	}

	node := &ASTNode{
		Label:       label,
		IndentLevel: level,
	}

	// Parameters
	if f.Type.Params != nil && len(f.Type.Params.List) > 0 {
		paramsNode := &ASTNode{
			Label:       "Params",
			IndentLevel: level + 1,
		}
		for _, param := range f.Type.Params.List {
			paramsNode.Children = append(paramsNode.Children, fieldToNode(param, level+2))
		}
		node.Children = append(node.Children, paramsNode)
	}

	// Results
	if f.Type.Results != nil && len(f.Type.Results.List) > 0 {
		resultsNode := &ASTNode{
			Label:       "Results",
			IndentLevel: level + 1,
		}
		for _, result := range f.Type.Results.List {
			resultsNode.Children = append(resultsNode.Children, fieldToNode(result, level+2))
		}
		node.Children = append(node.Children, resultsNode)
	}

	// Body
	if f.Body != nil {
		bodyNode := stmtToNode(f.Body, level+1)
		node.Children = append(node.Children, bodyNode)
	}

	return node
}

// stmtToNode converts a statement to a display node
func stmtToNode(stmt ast.Stmt, level int) *ASTNode {
	if stmt == nil {
		return nil
	}

	node := &ASTNode{
		Label:       reflect.TypeOf(stmt).String(),
		IndentLevel: level,
	}

	switch s := stmt.(type) {
	case *ast.BlockStmt:
		node.Label = "BlockStmt"
		for _, child := range s.List {
			if childNode := stmtToNode(child, level+1); childNode != nil {
				node.Children = append(node.Children, childNode)
			}
		}

	case *ast.ExprStmt:
		node.Label = "ExprStmt"
		node.Children = append(node.Children, exprToNode(s.X, level+1))

	case *ast.AssignStmt:
		node.Label = fmt.Sprintf("AssignStmt (%s)", s.Tok.String())
		for _, lhs := range s.Lhs {
			node.Children = append(node.Children, exprToNode(lhs, level+1))
		}
		for _, rhs := range s.Rhs {
			node.Children = append(node.Children, exprToNode(rhs, level+1))
		}

	case *ast.ReturnStmt:
		node.Label = "ReturnStmt"
		for _, result := range s.Results {
			node.Children = append(node.Children, exprToNode(result, level+1))
		}

	case *ast.IfStmt:
		node.Label = "IfStmt"
		if s.Init != nil {
			node.Children = append(node.Children, stmtToNode(s.Init, level+1))
		}
		node.Children = append(node.Children, exprToNode(s.Cond, level+1))
		node.Children = append(node.Children, stmtToNode(s.Body, level+1))
		if s.Else != nil {
			node.Children = append(node.Children, stmtToNode(s.Else, level+1))
		}

	case *ast.ForStmt:
		node.Label = "ForStmt"
		if s.Init != nil {
			node.Children = append(node.Children, stmtToNode(s.Init, level+1))
		}
		if s.Cond != nil {
			node.Children = append(node.Children, exprToNode(s.Cond, level+1))
		}
		if s.Post != nil {
			node.Children = append(node.Children, stmtToNode(s.Post, level+1))
		}
		node.Children = append(node.Children, stmtToNode(s.Body, level+1))

	case *ast.RangeStmt:
		node.Label = "RangeStmt"
		if s.Key != nil {
			node.Children = append(node.Children, exprToNode(s.Key, level+1))
		}
		if s.Value != nil {
			node.Children = append(node.Children, exprToNode(s.Value, level+1))
		}
		node.Children = append(node.Children, exprToNode(s.X, level+1))
		node.Children = append(node.Children, stmtToNode(s.Body, level+1))

	case *ast.DeclStmt:
		node.Label = "DeclStmt"
		declNodes := declToNode(s.Decl, level+1)
		node.Children = append(node.Children, declNodes...)

	case *ast.DeferStmt:
		node.Label = "DeferStmt"
		node.Children = append(node.Children, exprToNode(s.Call, level+1))

	case *ast.GoStmt:
		node.Label = "GoStmt"
		node.Children = append(node.Children, exprToNode(s.Call, level+1))

	case *ast.SwitchStmt:
		node.Label = "SwitchStmt"
		if s.Init != nil {
			node.Children = append(node.Children, stmtToNode(s.Init, level+1))
		}
		if s.Tag != nil {
			node.Children = append(node.Children, exprToNode(s.Tag, level+1))
		}
		node.Children = append(node.Children, stmtToNode(s.Body, level+1))

	case *ast.CaseClause:
		if len(s.List) == 0 {
			node.Label = "CaseClause (default)"
		} else {
			node.Label = "CaseClause"
			for _, expr := range s.List {
				node.Children = append(node.Children, exprToNode(expr, level+1))
			}
		}
		for _, stmt := range s.Body {
			node.Children = append(node.Children, stmtToNode(stmt, level+1))
		}

	case *ast.IncDecStmt:
		node.Label = fmt.Sprintf("IncDecStmt (%s)", s.Tok.String())
		node.Children = append(node.Children, exprToNode(s.X, level+1))

	case *ast.BranchStmt:
		if s.Label != nil {
			node.Label = fmt.Sprintf("BranchStmt (%s %s)", s.Tok.String(), s.Label.Name)
		} else {
			node.Label = fmt.Sprintf("BranchStmt (%s)", s.Tok.String())
		}
	}

	return node
}

// exprToNode converts an expression to a display node
func exprToNode(expr ast.Expr, level int) *ASTNode {
	if expr == nil {
		return nil
	}

	node := &ASTNode{
		Label:       exprToString(expr),
		IndentLevel: level,
	}

	switch e := expr.(type) {
	case *ast.CallExpr:
		node.Label = "CallExpr"
		node.Children = append(node.Children, &ASTNode{
			Label:       fmt.Sprintf("Fun: %s", exprToString(e.Fun)),
			IndentLevel: level + 1,
		})
		if len(e.Args) > 0 {
			argsNode := &ASTNode{
				Label:       "Args",
				IndentLevel: level + 1,
			}
			for _, arg := range e.Args {
				argsNode.Children = append(argsNode.Children, exprToNode(arg, level+2))
			}
			node.Children = append(node.Children, argsNode)
		}

	case *ast.BinaryExpr:
		node.Label = fmt.Sprintf("BinaryExpr (%s)", e.Op.String())
		node.Children = append(node.Children, exprToNode(e.X, level+1))
		node.Children = append(node.Children, exprToNode(e.Y, level+1))

	case *ast.UnaryExpr:
		node.Label = fmt.Sprintf("UnaryExpr (%s)", e.Op.String())
		node.Children = append(node.Children, exprToNode(e.X, level+1))

	case *ast.SelectorExpr:
		node.Label = fmt.Sprintf("SelectorExpr: %s.%s", exprToString(e.X), e.Sel.Name)

	case *ast.IndexExpr:
		node.Label = "IndexExpr"
		node.Children = append(node.Children, exprToNode(e.X, level+1))
		node.Children = append(node.Children, exprToNode(e.Index, level+1))

	case *ast.CompositeLit:
		node.Label = fmt.Sprintf("CompositeLit: %s", exprToString(e.Type))
		for _, elt := range e.Elts {
			node.Children = append(node.Children, exprToNode(elt, level+1))
		}

	case *ast.FuncLit:
		node.Label = "FuncLit"
		if e.Body != nil {
			node.Children = append(node.Children, stmtToNode(e.Body, level+1))
		}

	case *ast.KeyValueExpr:
		node.Label = "KeyValueExpr"
		node.Children = append(node.Children, &ASTNode{
			Label:       fmt.Sprintf("Key: %s", exprToString(e.Key)),
			IndentLevel: level + 1,
		})
		node.Children = append(node.Children, exprToNode(e.Value, level+1))

	case *ast.TypeAssertExpr:
		node.Label = "TypeAssertExpr"
		node.Children = append(node.Children, exprToNode(e.X, level+1))
		if e.Type != nil {
			node.Children = append(node.Children, &ASTNode{
				Label:       fmt.Sprintf("Type: %s", exprToString(e.Type)),
				IndentLevel: level + 1,
			})
		}

	case *ast.StarExpr:
		node.Label = "StarExpr"
		node.Children = append(node.Children, exprToNode(e.X, level+1))

	case *ast.SliceExpr:
		node.Label = "SliceExpr"
		node.Children = append(node.Children, exprToNode(e.X, level+1))
		if e.Low != nil {
			node.Children = append(node.Children, exprToNode(e.Low, level+1))
		}
		if e.High != nil {
			node.Children = append(node.Children, exprToNode(e.High, level+1))
		}
	}

	return node
}

// exprToString converts an expression to a string representation
func exprToString(expr ast.Expr) string {
	if expr == nil {
		return ""
	}

	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.BasicLit:
		return e.Value
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", exprToString(e.X), e.Sel.Name)
	case *ast.StarExpr:
		return "*" + exprToString(e.X)
	case *ast.ArrayType:
		if e.Len != nil {
			return fmt.Sprintf("[%s]%s", exprToString(e.Len), exprToString(e.Elt))
		}
		return "[]" + exprToString(e.Elt)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", exprToString(e.Key), exprToString(e.Value))
	case *ast.ChanType:
		switch e.Dir {
		case ast.SEND:
			return "chan<- " + exprToString(e.Value)
		case ast.RECV:
			return "<-chan " + exprToString(e.Value)
		default:
			return "chan " + exprToString(e.Value)
		}
	case *ast.FuncType:
		return "func(...)"
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.StructType:
		return "struct{...}"
	case *ast.Ellipsis:
		return "..." + exprToString(e.Elt)
	case *ast.CallExpr:
		return fmt.Sprintf("%s(...)", exprToString(e.Fun))
	case *ast.IndexExpr:
		return fmt.Sprintf("%s[%s]", exprToString(e.X), exprToString(e.Index))
	case *ast.IndexListExpr:
		return fmt.Sprintf("%s[...]", exprToString(e.X))
	case *ast.BinaryExpr:
		return fmt.Sprintf("%s %s %s", exprToString(e.X), e.Op.String(), exprToString(e.Y))
	case *ast.UnaryExpr:
		return fmt.Sprintf("%s%s", e.Op.String(), exprToString(e.X))
	case *ast.ParenExpr:
		return fmt.Sprintf("(%s)", exprToString(e.X))
	case *ast.CompositeLit:
		if e.Type != nil {
			return fmt.Sprintf("%s{...}", exprToString(e.Type))
		}
		return "{...}"
	case *ast.FuncLit:
		return "func(){...}"
	default:
		return fmt.Sprintf("%T", expr)
	}
}

// FlattenNodes flattens the tree structure for list display
func FlattenNodes(nodes []*ASTNode) []*ASTNode {
	var result []*ASTNode
	for _, node := range nodes {
		result = append(result, flattenNode(node)...)
	}
	return result
}

func flattenNode(node *ASTNode) []*ASTNode {
	result := []*ASTNode{node}
	if !node.Collapsed {
		for _, child := range node.Children {
			result = append(result, flattenNode(child)...)
		}
	}
	return result
}
