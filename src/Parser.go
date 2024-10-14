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

	// TODO: When calling these fucntions,
	// we already assert that the first
	// token is correct, so it might be
	// alright to change the code and
	// make that assumption in the
	// functions

	for t.kind != 0 {

		switch t.kind {
		case T_IDENTIFIER: // Variable declaration
			n = p.variableDeclaration()
			program.children = append(program.children, n)

		case T_IF: // If block
			n = p.ifBlock()
			program.children = append(program.children, n)

		case T_FOREVER: // Forever loop
			n = p.foreverLoop()
			program.children = append(program.children, n)

		case T_RANGE: // Range loop
			n = p.rangeLoop()
			program.children = append(program.children, n)

		case T_FOR: // For loop
			n = p.forLoop()
			program.children = append(program.children, n)

		case T_CALL: // Empty call
			n = p.funcCall()
			program.children = append(program.children, n)

		case T_STRUCT: // Struct definition
			n = p.structDef()
			program.children = append(program.children, n)

		case T_FUN: // Function definition
			n = p.funcDef()
			program.children = append(program.children, n)

		case T_RET: // Return statement
			n = p.retStatement()
			program.children = append(program.children, n)

		case T_BREAK: // Break statement
			n = p.breakStatement()
			program.children = append(program.children, n)

		case T_CONT: // Continue statement
			n = p.contStatement()
			program.children = append(program.children, n)

		case T_ENUM: // Enum definition
			n = p.enumDef()
			program.children = append(program.children, n)

		case T_TYPEDEF: // Type definition
			n = p.typeDef()
			program.children = append(program.children, n)

		default:
			panic(fmt.Sprint("Bad start to statement:", t.kind))
		}

		p.index++

		t = p.curToken()
	}

	return &program
}

func (p *Parser) variableDeclaration() *Node {
	var t Token

	n := Node{kind: N_VAR_DECLARATION}

	n.children = append(n.children, p.assignment())
	p.index++

	t = p.curToken()
	if t.kind != T_SEMICOLON {
		panic("Expected semicolon")
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON, data: t.data})

	return &n
}

func (p *Parser) ifBlock() *Node {
	var t Token

	n := Node{kind: N_VAR_DECLARATION}

	t = p.curToken()
	if t.kind != T_IF {
		panic("Expected if")
	}
	n.children = append(n.children, &Node{kind: N_IF, data: t.data})
	p.index++

	n.children = append(n.children, p.condition())
	p.index++

	n.children = append(n.children, p.block())
	p.index++

	t = p.curToken()
	for t.kind == T_ELIF {
		n.children = append(n.children, &Node{kind: N_ELIF, data: t.data})
		p.index++

		n.children = append(n.children, p.condition())
		p.index++

		n.children = append(n.children, p.block())
		p.index++

		t = p.curToken()
	}

	if t.kind == T_ELSE {
		n.children = append(n.children, &Node{kind: N_ELSE, data: t.data})
		p.index++

		n.children = append(n.children, p.block())
	}

	return &n
}

func (p *Parser) foreverLoop() *Node {
	n := Node{kind: N_FOREVER_LOOP}

	var t Token

	t = p.curToken()
	if t.kind != T_FOREVER {
		panic("Expected forever")
	}
	n.children = append(n.children, &Node{kind: N_FOREVER})
	p.index++

	n.children = append(n.children, p.block())

	return &n
}

func (p *Parser) rangeLoop() *Node {
	n := Node{kind: N_RANGE_LOOP}

	var t Token

	t = p.curToken()
	if t.kind != T_RANGE {
		panic("Expected range")
	}
	n.children = append(n.children, &Node{kind: N_RANGE})
	p.index++

	n.children = append(n.children, p.expression())
	p.index++

	n.children = append(n.children, p.block())

	return &n
}

func (p *Parser) forLoop() *Node {
	n := Node{kind: N_FOR_LOOP}

	var t Token

	t = p.curToken()
	if t.kind != T_RANGE {
		panic("Expected for")
	}
	n.children = append(n.children, &Node{kind: N_FOR})
	p.index++

	// TODO: Try skipping assignment
	n.children = append(n.children, p.assignment())
	p.index++

	t = p.curToken()
	if t.kind != T_SEMICOLON {
		panic("Expected semicolon")
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON})
	p.index++

	// TODO: Try skipping expression
	n.children = append(n.children, p.expression())
	p.index++

	t = p.curToken()
	if t.kind != T_SEMICOLON {
		panic("Expected semicolon")
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON})
	p.index++

	// TODO: Try skipping assignment
	n.children = append(n.children, p.assignment())
	p.index++

	n.children = append(n.children, p.block())

	return &n
}

func (p *Parser) structDef() *Node {
	var t Token

	n := Node{kind: N_STRUCT_DEF}

	p.assignment()
	p.index++

	t = p.curToken()
	if t.kind != T_SEMICOLON {
		panic("Expected semicolon")
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON, data: t.data})

	return &n
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
