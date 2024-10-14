package main

import "fmt"

type Parser struct {
	source []Token
	index  int
}

func (p *Parser) curToken() Token {
	if p.index < len(p.source) {
		return p.source[p.index]
	}

	// Illegal token
	return Token{}
}

// The main program
func (p *Parser) parse() *Node {
	program := Node{kind: N_PROGRAM}

	var n *Node

	t := p.curToken()

	switch t.kind {
	case T_IDENTIFIER: // Variable declaration
		n = p.variableDeclaration()
		program.children = append(program.children, n)

	case T_IF: // If block
	case T_FOREVER: // Forever loop
	case T_RANGE: // Range loop
	case T_FOR: // For loop
	case T_CALL: // Empty call
	case T_STRUCT: // Struct definition
	case T_FUN: // Function definition
	case T_RET: // Return statement
	case T_BREAK: // Break statement
	case T_CONT: // Continue statement
	case T_ENUM: // Enum definition
	case T_TYPEDEF: // Type definition
		n = p.typeDef()
		program.children = append(program.children, n)

	default:
		panic(fmt.Sprint("Bad start to statement:", t.kind))
	}

	return &program
}

func (p *Parser) variableDeclaration() *Node {
	var t Token

	n := Node{kind: N_VAR_DECLARATION}

	p.assignment()
	p.index++

	t = p.curToken()
	if t.kind != T_SEMICOLON {
		panic("Expected semicolon")
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON, data: t.data})

	return &n
}

func (p *Parser) ifBlock() *Node {
	return nil
}

func (p *Parser) foreverLoop() *Node {
	return nil
}

func (p *Parser) rangeLoop() *Node {
	return nil
}

func (p *Parser) forLoop() *Node {
	return nil
}

func (p *Parser) callStatement() *Node {
	return nil
}

func (p *Parser) structDef() *Node {
	return nil
}

func (p *Parser) funcDef() *Node {
	return nil
}

func (p *Parser) retStatement() *Node {
	return nil
}

func (p *Parser) breakStatement() *Node {
	return nil
}

func (p *Parser) contStatement() *Node {
	return nil
}

func (p *Parser) enumDef() *Node {
	return nil
}

// For now they are identical, but
// further logic will be applied later
func (p *Parser) condition() *Node {
	return p.expression()
}

func (p *Parser) expression() *Node {
	return nil
}

func (p *Parser) assignment() *Node {
	var t Token

	n := Node{kind: N_ASSIGNMENT}

	// x
	t = p.curToken()
	if t.kind != T_IDENTIFIER {
		panic("Expected identifier")
	}
	n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: t.data})
	p.index++

	// x int
	t = p.curToken()
	if t.kind == T_TYPE {
		n.children = append(n.children, &Node{kind: N_TYPE, data: t.data})
		p.index++
	}

	// x int =
	t = p.curToken()
	if t.kind != T_ASSIGN {
		panic("Expected identifier")
	}
	n.children = append(n.children, &Node{kind: N_ASSIGN, data: t.data})
	p.index++

	// x int = 3 + 7
	n.children = append(n.children, p.expression())

	return nil
}

func (p *Parser) funcCall() *Node {
	return nil
}

func (p *Parser) block() *Node {
	return nil
}

func (p *Parser) typeDef() *Node {
	n := Node{kind: N_NEW_TYPE}

	var t Token

	t = p.curToken()
	if t.kind != T_TYPEDEF {
		panic("Expected typedef")
	}
	n.children = append(n.children, &Node{kind: N_TYPEDEF})
	p.index++

	t = p.curToken()
	if t.kind != T_IDENTIFIER {
		panic("Expected identifier")
	}
	n.children = append(n.children, &Node{kind: N_IDENTIFIER})
	p.index++

	t = p.curToken()
	if t.kind != T_TYPE {
		panic("Expected type")
	}
	n.children = append(n.children, &Node{kind: N_TYPE})
	p.index++

	t = p.curToken()
	if t.kind != T_SEMICOLON {
		panic("Expected semicolon")
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON})

	return nil
}
