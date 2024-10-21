package main

import (
	"fmt"
)

const JOB_PARSER = "Parser"

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

func (p *Parser) peekToken() Token {
	if p.index+1 < len(p.source) {
		return p.source[p.index+1]
	}

	// Illegal token
	return Token{}
}

// The main program
func (p *Parser) parse() *Node {
	const FUNC_NAME = "program"

	program := Node{kind: N_PROGRAM}

	var n *Node

	t := p.curToken()

	// TODO: When calling these fucntions,
	// we already assert that the first
	// token is correct, so it might be
	// alright to change the code and
	// make that assumption in the
	// functions

	for t.kind != T_ILLEGAL {

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
			n = p.loneCall()
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

		case T_L_BLOCK: // Assigning to an element in an array (or map?)
			n = p.elementAssignment()
			program.children = append(program.children, n)

		default:
			panic(fmt.Sprint("Bad start to statement: ", t.kind, "on line", t.line))
		}

		p.index++

		t = p.curToken()
	}

	return &program
}

func (p *Parser) variableDeclaration() *Node {
	const FUNC_NAME = "variable declaration"

	var t Token

	n := Node{kind: N_VAR_DECLARATION}

	n.children = append(n.children, p.assignment())
	p.index++

	t = p.curToken()
	if t.kind != T_SEMICOLON {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "semicolon", t)
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON, data: t.data})

	return &n
}

func (p *Parser) elementAssignment() *Node {
	const FUNC_NAME = "element assignment"

	var t Token

	n := Node{kind: N_ELEMENT_ASSIGNMENT}

	n.children = append(n.children, p.indexUnary())
	p.index++

	n.children = append(n.children, p.assignment())
	p.index++

	t = p.curToken()
	if t.kind != T_SEMICOLON {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "semicolon", t)
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON, data: t.data})

	return &n
}

func (p *Parser) ifBlock() *Node {
	const FUNC_NAME = "if block"

	var t Token

	n := Node{kind: N_IF_BLOCK}

	t = p.curToken()
	if t.kind != T_IF {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "if", t)
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
	} else {
		p.index--
	}

	return &n
}

func (p *Parser) foreverLoop() *Node {
	const FUNC_NAME = "forever loop"

	n := Node{kind: N_FOREVER_LOOP}

	var t Token

	t = p.curToken()
	if t.kind != T_FOREVER {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "forever", t)
	}
	n.children = append(n.children, &Node{kind: N_FOREVER})
	p.index++

	n.children = append(n.children, p.block())

	return &n
}

func (p *Parser) rangeLoop() *Node {
	const FUNC_NAME = "range loop"

	n := Node{kind: N_RANGE_LOOP}

	var t Token

	t = p.curToken()
	if t.kind != T_RANGE {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "range", t)
	}
	n.children = append(n.children, &Node{kind: N_RANGE})
	p.index++

	n.children = append(n.children, p.expression())
	p.index++

	n.children = append(n.children, p.block())

	return &n
}

func (p *Parser) forLoop() *Node {
	const FUNC_NAME = "for loop"

	n := Node{kind: N_FOR_LOOP}

	var t Token

	t = p.curToken()
	if t.kind != T_FOR {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "for", t)
	}
	n.children = append(n.children, &Node{kind: N_FOR})
	p.index++

	// TODO: Try skipping assignment
	n.children = append(n.children, p.assignment())
	p.index++

	t = p.curToken()
	if t.kind != T_SEMICOLON {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "semicolon", t)
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON})
	p.index++

	// TODO: Try skipping expression
	n.children = append(n.children, p.expression())
	p.index++

	t = p.curToken()
	if t.kind != T_SEMICOLON {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "semicolon", t)
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
	const FUNC_NAME = "struct definition"

	var t Token

	n := Node{kind: N_STRUCT_DEF}

	t = p.curToken()
	if t.kind != T_STRUCT {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "struct", t)
	}
	n.children = append(n.children, &Node{kind: N_STRUCT, data: t.data})
	p.index++

	// NOTE: The lexer has figured out for
	// the parser that the name for the
	// struct is a type, and therefore we
	// check for type (over identifier)
	t = p.curToken()
	if t.kind != T_TYPE {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "type", t)
	}
	n.children = append(n.children, &Node{kind: N_TYPE, data: t.data})
	p.index++

	t = p.curToken()
	if t.kind != T_L_SQUIRLY {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "left squirly", t)
	}
	n.children = append(n.children, &Node{kind: N_L_SQUIRLY, data: t.data})
	p.index++

	t = p.curToken()

	for t.kind != T_R_SQUIRLY && t.kind != T_ILLEGAL {

		if t.kind != T_IDENTIFIER {
			throwError(JOB_PARSER, FUNC_NAME, t.line, "identifier", t)
		}
		n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: t.data})
		p.index++

		n.children = append(n.children, p.complexType())
		p.index++
		t = p.curToken()
	}

	if t.kind != T_R_SQUIRLY {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "right squirly", t)
	}
	n.children = append(n.children, &Node{kind: N_R_SQUIRLY, data: t.data})

	return &n
}

func (p *Parser) funcDef() *Node {
	const FUNC_NAME = "function definition"

	var t Token

	n := Node{kind: N_FUNC_DEF}

	t = p.curToken()
	if t.kind != T_FUN {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "fun", t)
	}
	n.children = append(n.children, &Node{kind: N_FUN, data: t.data})
	p.index++
	t = p.curToken()

	// Method on struct
	if t.kind == T_L_PAREN {
		n.children = append(n.children, &Node{kind: N_L_PAREN, data: t.data})
		p.index++

		t = p.curToken()
		if t.kind != T_IDENTIFIER {
			throwError(JOB_PARSER, FUNC_NAME, t.line, "identifier", t)
		}
		n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: t.data})
		p.index++

		n.children = append(n.children, p.complexType())
		p.index++

		t = p.curToken()
		if t.kind != T_R_PAREN {
			throwError(JOB_PARSER, FUNC_NAME, t.line, "right paren", t)
		}
		n.children = append(n.children, &Node{kind: N_R_PAREN, data: t.data})
		p.index++
		t = p.curToken()
	}

	if t.kind != T_IDENTIFIER {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "identifier", t)
	}
	n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: t.data})
	p.index++

	t = p.curToken()
	if t.kind != T_L_PAREN {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "left paren", t)
	}
	n.children = append(n.children, &Node{kind: N_L_PAREN, data: t.data})
	p.index++

	t = p.curToken()
	if t.kind == T_IDENTIFIER {
		n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: t.data})
		p.index++

		n.children = append(n.children, p.complexType())
		p.index++

		t = p.curToken()
		for t.kind == T_SEP {
			n.children = append(n.children, &Node{kind: N_SEP, data: t.data})
			p.index++

			t = p.curToken()
			if t.kind != T_IDENTIFIER {
				throwError(JOB_PARSER, FUNC_NAME, t.line, "identifier", t)
			}
			n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: t.data})
			p.index++

			t = p.curToken()
			if t.kind != T_TYPE {
				throwError(JOB_PARSER, FUNC_NAME, t.line, "type", t)
			}
			n.children = append(n.children, &Node{kind: N_TYPE, data: t.data})
			p.index++

			t = p.curToken()
		}
	}

	if t.kind != T_R_PAREN {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "right paren", t)
	}
	n.children = append(n.children, &Node{kind: N_R_PAREN, data: t.data})
	p.index++

	// Return rype?
	t = p.curToken()
	if t.kind != T_L_SQUIRLY {
		n.children = append(n.children, p.complexType())
		p.index++
	}

	n.children = append(n.children, p.block())

	return &n
}

func (p *Parser) retStatement() *Node {
	const FUNC_NAME = "return statement"

	var t Token

	n := Node{kind: N_RET_STATE}

	t = p.curToken()
	if t.kind != T_RET {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "return", t)
	}
	n.children = append(n.children, &Node{kind: N_RET, data: t.data})
	p.index++

	// don't do the extra value
	t = p.curToken()
	if t.kind == T_SEMICOLON {
		n.children = append(n.children, &Node{kind: N_SEMICOLON, data: t.data})
		return &n
	}

	n.children = append(n.children, p.expression())
	p.index++

	t = p.curToken()
	if t.kind != T_SEMICOLON {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "semicolon", t)
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON, data: t.data})

	return &n
}

func (p *Parser) breakStatement() *Node {
	const FUNC_NAME = "break statement"

	var t Token

	n := Node{kind: N_BREAK_STATE}

	t = p.curToken()
	if t.kind != T_BREAK {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "break", t)
	}
	n.children = append(n.children, &Node{kind: N_BREAK, data: t.data})
	p.index++

	// don't do the extra value
	t = p.curToken()
	if t.kind == T_SEMICOLON {
		n.children = append(n.children, &Node{kind: N_SEMICOLON, data: t.data})
		return &n
	}

	n.children = append(n.children, p.expression())
	p.index++

	t = p.curToken()
	if t.kind != T_SEMICOLON {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "semicolon", t)
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON, data: t.data})

	return &n
}

func (p *Parser) contStatement() *Node {
	const FUNC_NAME = "continue statement"

	var t Token

	n := Node{kind: N_CONT_STATE}

	t = p.curToken()
	if t.kind != T_CONT {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "continue", t)
	}
	n.children = append(n.children, &Node{kind: N_CONT, data: t.data})
	p.index++

	// don't do the extra value
	// (same as "continue 0")
	t = p.curToken()
	if t.kind == T_SEMICOLON {
		n.children = append(n.children, &Node{kind: N_SEMICOLON, data: t.data})
		return &n
	}

	n.children = append(n.children, p.expression())
	p.index++

	t = p.curToken()
	if t.kind != T_SEMICOLON {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "semicolon", t)
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON, data: t.data})

	return &n
}

func (p *Parser) enumDef() *Node {
	const FUNC_NAME = "enum definition"

	var t Token

	n := Node{kind: N_ENUM_DEF}

	t = p.curToken()
	if t.kind != T_ENUM {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "enum", t)
	}
	n.children = append(n.children, &Node{kind: N_ENUM, data: t.data})
	p.index++

	t = p.curToken()
	if t.kind != T_IDENTIFIER {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "identifier", t)
	}
	n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: t.data})
	p.index++

	t = p.curToken()
	if t.kind != T_L_SQUIRLY {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "left squirly", t)
	}
	n.children = append(n.children, &Node{kind: N_L_SQUIRLY, data: t.data})
	p.index++

	t = p.curToken()

	for t.kind != T_R_SQUIRLY && t.kind != T_ILLEGAL {

		if t.kind != T_IDENTIFIER {
			throwError(JOB_PARSER, FUNC_NAME, t.line, "identifier", t)
		}
		n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: t.data})
		p.index++

		t = p.curToken()
		if t.kind != T_SEP {
			throwError(JOB_PARSER, FUNC_NAME, t.line, "separator", t)
		}
		n.children = append(n.children, &Node{kind: N_SEP, data: t.data})
		p.index++

		t = p.curToken()
	}

	if t.kind != T_R_SQUIRLY {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "right squirly", t)
	}
	n.children = append(n.children, &Node{kind: N_R_SQUIRLY, data: t.data})

	return &n
}

// For now they are identical, but
// further logic will be applied later
func (p *Parser) condition() *Node {
	const FUNC_NAME = "condition"

	return p.expression()
}

func (p *Parser) expression() *Node {
	const FUNC_NAME = "expression"

	//var t Token
	var v *Node

	n := Node{kind: N_EXPRESSION}

	// 7
	n.children = append(n.children, p.value())
	p.index++

	v = p.operator()
	for v != nil {
		n.children = append(n.children, v)
		p.index++
		n.children = append(n.children, p.value())
		p.index++
		v = p.operator()
		v = p.operator()
	}

	p.index--

	return &n
}

func (p *Parser) assignment() *Node {
	const FUNC_NAME = "assignment"

	var t Token

	n := Node{kind: N_ASSIGNMENT}

	// x
	t = p.curToken()
	if t.kind != T_IDENTIFIER {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "identifier", t)
	}
	n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: t.data})
	p.index++

	t = p.curToken()

	// Access
	if t.kind == T_ACCESS {
		n.children = append(n.children, &Node{kind: N_ACCESS, data: t.data})
		p.index++
		t = p.curToken()

		if t.kind != T_IDENTIFIER {
			throwError(JOB_PARSER, FUNC_NAME, t.line, "identifier", t)
		}
		n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: t.data})
		p.index++
		t = p.curToken()
	}

	// We have a type?
	if t.kind != T_ASSIGN {
		n.children = append(n.children, p.complexType())

		if t.kind == T_MAP {
			return &n
		}

		p.index++
	}

	// Now we MUST have an assign
	t = p.curToken()
	if t.kind != T_ASSIGN {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "assign", t)
	}

	n.children = append(n.children, &Node{kind: N_ASSIGN, data: t.data})
	p.index++

	// x int = 3 + 7
	n.children = append(n.children, p.expression())

	return &n
}

func (p *Parser) loneCall() *Node {
	const FUNC_NAME = "lone call"

	var t Token

	n := Node{kind: N_LONE_CALL}

	n.children = append(n.children, p.funcCall())

	p.index++

	t = p.curToken()
	if t.kind != T_SEMICOLON {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "semicolon", t)
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON, data: t.data})

	return &n
}

func (p *Parser) funcCall() *Node {
	const FUNC_NAME = "function call"

	var t Token

	n := Node{kind: N_FUNC_CALL}

	t = p.curToken()
	if t.kind != T_CALL {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "call", t)
	}
	n.children = append(n.children, &Node{kind: N_CALL, data: t.data})
	p.index++

	t = p.curToken()
	if t.kind != T_IDENTIFIER {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "identifier", t)
	}
	n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: t.data})
	p.index++

	t = p.curToken()
	if t.kind != T_L_PAREN {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "left paren", t)
	}
	n.children = append(n.children, &Node{kind: N_L_PAREN, data: t.data})
	p.index++

	t = p.curToken()

	if t.kind != T_R_PAREN {
		n.children = append(n.children, p.value())
		p.index++

		t = p.curToken()
		for t.kind == T_SEP {
			n.children = append(n.children, &Node{kind: N_SEP, data: t.data})
			p.index++

			n.children = append(n.children, p.value())
			p.index++

			t = p.curToken()
		}
	}

	t = p.curToken()
	if t.kind != T_R_PAREN {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "right paren", t)
	}
	n.children = append(n.children, &Node{kind: N_R_PAREN, data: t.data})

	return &n
}

func (p *Parser) structNew() *Node {
	const FUNC_NAME = "new struct"

	var t Token

	n := Node{kind: N_STRUCT_NEW}

	t = p.curToken()
	if t.kind != T_NEW {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "new", t)
	}
	n.children = append(n.children, &Node{kind: N_NEW, data: t.data})
	p.index++

	t = p.curToken()
	if t.kind != T_TYPE {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "type", t)
	}
	n.children = append(n.children, &Node{kind: N_TYPE, data: t.data})
	p.index++

	t = p.curToken()
	if t.kind != T_L_PAREN {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "left paren", t)
	}
	n.children = append(n.children, &Node{kind: N_L_PAREN, data: t.data})
	p.index++

	t = p.curToken()

	if t.kind != T_R_PAREN {
		n.children = append(n.children, p.value())
		p.index++

		t = p.curToken()
		for t.kind == T_SEP {
			n.children = append(n.children, &Node{kind: N_SEP, data: t.data})
			p.index++

			n.children = append(n.children, p.value())
			p.index++

			t = p.curToken()
		}
	}

	t = p.curToken()
	if t.kind != T_R_PAREN {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "right paren", t)
	}
	n.children = append(n.children, &Node{kind: N_R_PAREN, data: t.data})

	return &n
}

func (p *Parser) block() *Node {
	const FUNC_NAME = "block"

	block := Node{kind: N_BLOCK}

	var n *Node

	t := p.curToken()

	if t.kind != T_L_SQUIRLY {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "left squirly", t)
	}
	block.children = append(block.children, &Node{kind: N_L_SQUIRLY, data: t.data})
	p.index++

	t = p.curToken()

	for t.kind != T_ILLEGAL && t.kind != T_R_SQUIRLY {
		switch t.kind {
		case T_IDENTIFIER: // Variable declaration
			n = p.variableDeclaration()
			block.children = append(block.children, n)

		case T_IF: // If block
			n = p.ifBlock()
			block.children = append(block.children, n)

		case T_FOREVER: // Forever loop
			n = p.foreverLoop()
			block.children = append(block.children, n)

		case T_RANGE: // Range loop
			n = p.rangeLoop()
			block.children = append(block.children, n)

		case T_FOR: // For loop
			n = p.forLoop()
			block.children = append(block.children, n)

		case T_CALL: // Empty call
			n = p.loneCall()
			block.children = append(block.children, n)

		case T_STRUCT: // Struct definition
			n = p.structDef()
			block.children = append(block.children, n)

		case T_FUN: // Function definition
			n = p.funcDef()
			block.children = append(block.children, n)

		case T_RET: // Return statement
			n = p.retStatement()
			block.children = append(block.children, n)

		case T_BREAK: // Break statement
			n = p.breakStatement()
			block.children = append(block.children, n)

		case T_CONT: // Continue statement
			n = p.contStatement()
			block.children = append(block.children, n)

		case T_ENUM: // Enum definition
			n = p.enumDef()
			block.children = append(block.children, n)

		case T_TYPEDEF: // Type definition
			n = p.typeDef()
			block.children = append(block.children, n)

		case T_L_BLOCK: // Assigning to an element in an array (or map?)
			n = p.elementAssignment()
			block.children = append(block.children, n)

		default:
			panic(fmt.Sprint("Bad start to statement: ", t.kind, " on line ", t.line))
		}

		p.index++

		t = p.curToken()
	}

	if t.kind != T_R_SQUIRLY {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "right squirly", t)
	}
	block.children = append(block.children, &Node{kind: N_R_SQUIRLY, data: t.data})

	return &block
}

func (p *Parser) typeDef() *Node {
	const FUNC_NAME = "type definition"

	n := Node{kind: N_NEW_TYPE}

	var t Token

	t = p.curToken()
	if t.kind != T_TYPEDEF {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "typedef", t)
	}
	n.children = append(n.children, &Node{kind: N_TYPEDEF})
	p.index++

	t = p.curToken()
	if t.kind != T_TYPE {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "type", t)
	}
	n.children = append(n.children, &Node{kind: N_TYPE})
	p.index++

	t = p.curToken()
	if t.kind != T_TYPE {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "type", t)
	}
	n.children = append(n.children, &Node{kind: N_TYPE})
	p.index++

	t = p.curToken()
	if t.kind != T_SEMICOLON {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "semicolon", t)
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON})

	return &n
}

// TODO: Dealing with unary -?
// TODO: What happens with p.x (accessors)
func (p *Parser) value() *Node {
	const FUNC_NAME = "value"

	var t Token
	var n Node

	t = p.curToken()

	switch t.kind {
	case T_IDENTIFIER:
		return &Node{kind: N_IDENTIFIER, data: t.data}
	case T_INT:
		return &Node{kind: N_INT, data: t.data}
	case T_FLOAT:
		return &Node{kind: N_FLOAT, data: t.data}
	case T_STRING:
		return &Node{kind: N_STRING, data: t.data}
	case T_CHAR:
		return &Node{kind: N_CHAR, data: t.data}
	case T_BOOL:
		return &Node{kind: N_BOOL, data: t.data}
	case T_NIL:
		return &Node{kind: N_NIL, data: t.data}

		// Unary cases
	case T_NOT:
		n = Node{kind: N_UNARY_OPERATION, data: t.data}
		n.children = append(n.children, &Node{kind: N_NOT, data: t.data})
		p.index++

		n.children = append(n.children, p.value())

		return &n

	case T_INC:
		n = Node{kind: N_UNARY_OPERATION, data: t.data}
		n.children = append(n.children, &Node{kind: N_INC, data: t.data})
		p.index++

		n.children = append(n.children, p.value())

		return &n
	case T_DINC:
		n = Node{kind: N_UNARY_OPERATION, data: t.data}
		n.children = append(n.children, &Node{kind: N_DINC, data: t.data})
		p.index++

		n.children = append(n.children, p.value())

		return &n
	case T_REF:
		n = Node{kind: N_UNARY_OPERATION, data: t.data}
		n.children = append(n.children, &Node{kind: N_REF, data: t.data})
		p.index++

		n.children = append(n.children, p.value())

		return &n
	case T_DEREF:
		n = Node{kind: N_UNARY_OPERATION, data: t.data}
		n.children = append(n.children, &Node{kind: N_DEREF, data: t.data})
		p.index++

		n.children = append(n.children, p.value())

		return &n
	case T_L_BLOCK:
		n = Node{kind: N_UNARY_OPERATION, data: t.data}
		n.children = append(n.children, p.indexUnary())
		p.index++

		n.children = append(n.children, p.value())

		return &n

		// We're using an array or slice
	case T_MAKE:
		return p.makeArray()

		// Someones trying to do some bedmas
	case T_L_PAREN:
		n = Node{kind: N_BRACKETED_VALUE}

		n.children = append(n.children, &Node{kind: N_L_PAREN, data: t.data})

		p.index++
		n.children = append(n.children, p.expression())

		p.index++
		t = p.curToken()
		if t.kind != T_R_PAREN {
			throwError(JOB_PARSER, FUNC_NAME, t.line, "right paren", t)
		}
		n.children = append(n.children, &Node{kind: N_R_PAREN, data: t.data})

		return &n

		// Calling a function to use as a value
	case T_CALL:
		panic("Not implemented")

		// Creating a new struct
	case T_NEW:
		return p.structNew()

	default:
		throwError(JOB_PARSER, FUNC_NAME, t.line, "unary or value", t)

		// This never executes because
		// throwError panics
		return nil
	}
}

// Operator doesn't panic but returns
// nil, making it an unsafe operation,
// however, this is necessary to signal
// to expression that the expression is
// continued
func (p *Parser) operator() *Node {
	const FUNC_NAME = "operator"

	var t Token

	t = p.curToken()

	switch t.kind {
	case T_ADD:
		return &Node{kind: N_ADD, data: t.data}
	case T_XOR:
		return &Node{kind: N_XOR, data: t.data}
	case T_ACCESS:
		return &Node{kind: N_ACCESS, data: t.data}
	case T_NEQ:
		return &Node{kind: N_NEQ, data: t.data}
	case T_MOD:
		return &Node{kind: N_MOD, data: t.data}
	case T_EQ:
		return &Node{kind: N_EQ, data: t.data}
	case T_LT:
		return &Node{kind: N_LT, data: t.data}
	case T_GT:
		return &Node{kind: N_GT, data: t.data}
	case T_LTEQ:
		return &Node{kind: N_LTEQ, data: t.data}
	case T_GTEQ:
		return &Node{kind: N_GTEQ, data: t.data}
	case T_SUB:
		return &Node{kind: N_SUB, data: t.data}
	case T_MUL:
		return &Node{kind: N_MUL, data: t.data}
	case T_DIV:
		return &Node{kind: N_DIV, data: t.data}
	case T_OR:
		return &Node{kind: N_OR, data: t.data}
	case T_AND:
		return &Node{kind: N_AND, data: t.data}
	case T_L_SHIFT:
		return &Node{kind: N_L_SHIFT, data: t.data}
	case T_R_SHIFT:
		return &Node{kind: N_R_SHIFT, data: t.data}

	default:
		return nil
	}
}

func (p *Parser) indexUnary() *Node {
	const FUNC_NAME = "index unary"

	n := Node{kind: N_INDEX}

	var t Token

	t = p.curToken()
	if t.kind != T_L_BLOCK {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "left block", t)
	}
	n.children = append(n.children, &Node{kind: N_L_BLOCK, data: t.data})

	p.index++

	n.children = append(n.children, p.expression())

	p.index++
	t = p.curToken()
	if t.kind != T_R_BLOCK {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "right block", t)
	}
	n.children = append(n.children, &Node{kind: N_R_BLOCK, data: t.data})

	return &n
}

func (p *Parser) makeArray() *Node {
	const FUNC_NAME = "make array"

	n := Node{kind: N_MAKE_ARRAY}

	var t Token

	t = p.curToken()
	if t.kind != T_MAKE {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "make", t)
	}
	n.children = append(n.children, &Node{kind: N_MAKE, data: t.data})

	p.index++
	t = p.curToken()
	if t.kind != T_L_BLOCK {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "left block", t)
	}
	n.children = append(n.children, &Node{kind: N_L_BLOCK, data: t.data})

	p.index++
	t = p.curToken()
	for t.kind != T_R_BLOCK {

		n.children = append(n.children, p.expression())

		p.index++
		t = p.curToken()
		if t.kind != T_SEP {
			break // Last element
		}
		n.children = append(n.children, &Node{kind: N_SEP, data: t.data})

		p.index++
		t = p.curToken()
	}

	if t.kind != T_R_BLOCK {
		throwError(JOB_PARSER, FUNC_NAME, t.line, "right block", t)
	}
	n.children = append(n.children, &Node{kind: N_R_BLOCK, data: t.data})

	return &n
}

func (p *Parser) complexType() *Node {
	const FUNC_NAME = "complex type"

	n := Node{kind: N_COMPLEX_TYPE}

	var t Token

	t = p.curToken()

	// First check if we're making a map
	if t.kind == T_MAP {
		n.children = append(n.children, &Node{kind: N_MAP, data: t.data})
		p.index++

		t = p.curToken()
		if t.kind != T_TYPE {
			throwError(JOB_PARSER, FUNC_NAME, t.line, "type", t)
		}
		n.children = append(n.children, &Node{kind: N_TYPE, data: t.data})
		p.index++

		t = p.curToken()
		if t.kind != T_L_BLOCK {
			throwError(JOB_PARSER, FUNC_NAME, t.line, "left block", t)
		}
		n.children = append(n.children, &Node{kind: N_L_BLOCK, data: t.data})
		p.index++

		// NOTE: The logic for determining
		// if an identifier inside the brackets
		// is a type or not isn't so clean, so
		// we left it as an identifier in the
		// lexer, but here, we can finally
		// change it to a type
		t = p.curToken()
		if t.kind != T_IDENTIFIER && t.kind != T_TYPE {
			throwError(JOB_PARSER, FUNC_NAME, t.line, "identifier (but really a type)", t)
		}
		n.children = append(n.children, &Node{kind: N_TYPE, data: t.data})
		p.index++

		t = p.curToken()
		if t.kind != T_R_BLOCK {
			throwError(JOB_PARSER, FUNC_NAME, t.line, "right block", t)
		}
		n.children = append(n.children, &Node{kind: N_R_BLOCK, data: t.data})

		// Can't assign to map straight away
		// (for now because bruh)
		return &n

		// Otherwise we might have a more normal type
	} else if t.kind == T_TYPE {
		n.children = append(n.children, &Node{kind: N_TYPE, data: t.data})

		pt := p.peekToken()

		// Is it a pointer type?
		if pt.kind == T_DEREF {
			p.index++
			t = p.curToken()

			n.children = append(n.children, &Node{kind: N_DEREF, data: t.data})
			pt = p.peekToken()
		}

		if pt.kind == T_L_BLOCK { // Some larger type
			p.index++
			t = p.curToken()

			// ArrayList
			if p.peekToken().kind == T_R_BLOCK {
				n.children = append(n.children, &Node{kind: N_L_BLOCK, data: t.data})
				p.index++
				t = p.curToken()
				n.children = append(n.children, &Node{kind: N_R_BLOCK, data: t.data})
			} else { // Array
				n.children = append(n.children, p.indexUnary())
			}
		}
	}

	return &n
}
