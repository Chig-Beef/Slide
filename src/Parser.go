package main

import (
	"fmt"
	"strconv"
)

const JOB_PARSER = "Parser"

// TODO: Turn the function last called
// in the compiler into a stack so that
// recursion debuggins isn't do bad

type Parser struct {
	source []Token
	tok    Token
	index  int
}

func (p *Parser) nextToken() {
	t := Token{}

	if p.index < len(p.source) {
		t = p.source[p.index]
	}

	p.index++

	p.tok = t
}

func (p *Parser) peekToken() Token {
	if p.index < len(p.source) {
		return p.source[p.index]
	}

	// Illegal token
	return Token{}
}

// The main program
func (p *Parser) parse() *Node {
	const FUNC_NAME = "program"

	program := Node{kind: N_PROGRAM}

	var n *Node

	p.nextToken()

	// TODO: When calling these functions,
	// we already assert that the first
	// token is correct, so it might be
	// alright to change the code and
	// make that assumption in the
	// functions

	for p.tok.kind != T_ILLEGAL {

		switch p.tok.kind {
		case T_IDENTIFIER, T_L_BLOCK, T_DEREF: // Variable declaration
			n = p.variableDeclaration()
			program.children = append(program.children, n)

		case T_INC, T_DINC: // Lone increment
			n = p.loneIncrement()
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

		case T_WHILE: // While loop
			n = p.whileLoop()
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

		case T_SWITCH: // Switch statement
			n = p.switchStatement()
			program.children = append(program.children, n)

		case T_CONST: //  Constant definition
			n = p.constantStatement()
			program.children = append(program.children, n)

		default:
			throwError(p.tok.line, fmt.Sprint("Bad start to statement: ", p.tok.kind, " on line ", p.tok.line))
		}

		p.nextToken()
	}

	return &program
}

func (p *Parser) loneIncrement() *Node {
	const FUNC_NAME = "lone increment"

	n := Node{kind: N_LONE_INC, line: p.tok.line}

	if p.tok.kind == T_INC {
		n.children = append(n.children, &Node{kind: N_INC, data: p.tok.data, line: p.tok.line})
	} else if p.tok.kind == T_DINC {
		n.children = append(n.children, &Node{kind: N_DINC, data: p.tok.data, line: p.tok.line})
	} else {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "inc or dinc", p.tok))
	}
	p.nextToken()

	if p.tok.kind != T_IDENTIFIER {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok))
	}

	if p.peekToken().kind == T_ACCESS {
		// Property
		n.children = append(n.children, p.property())
		p.nextToken()
	} else {
		// Standalone
		n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: p.tok.data, line: p.tok.line})
		p.nextToken()
	}

	if p.tok.kind != T_SEMICOLON {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "semicolon", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON, data: p.tok.data, line: p.tok.line})

	return &n
}

func (p *Parser) constantStatement() *Node {
	const FUNC_NAME = "constant statement"

	n := Node{kind: N_CONSTANT, line: p.tok.line}

	if p.tok.kind != T_CONST {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "const", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_CONST, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	n.children = append(n.children, p.variableDeclaration())

	return &n
}

func (p *Parser) variableDeclaration() *Node {
	const FUNC_NAME = "variable declaration"

	n := Node{kind: N_VAR_DECLARATION, line: p.tok.line}

	// Element assignment
	if p.tok.kind == T_L_BLOCK {
		n.children = append(n.children, p.indexUnary())
		p.nextToken()
	} else if p.tok.kind == T_DEREF {
		n.children = append(n.children, &Node{kind: N_DEREF, data: p.tok.data, line: p.tok.line})
		p.nextToken()
	}

	n.children = append(n.children, p.assignment())

	p.nextToken()
	if p.tok.kind != T_SEMICOLON {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "semicolon", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON, data: p.tok.data, line: p.tok.line})

	return &n
}

func (p *Parser) switchStatement() *Node {
	const FUNC_NAME = "switch statement"

	n := Node{kind: N_SWITCH_STATE, line: p.tok.line}

	if p.tok.kind != T_SWITCH {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "switch", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_SWITCH, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	n.children = append(n.children, p.expression())
	p.nextToken()

	if p.tok.kind != T_L_SQUIRLY {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "left squirly", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_L_SQUIRLY, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	// Cases here
	for p.tok.kind == T_CASE {
		n.children = append(n.children, p.caseStatement())
		p.nextToken()
	}

	if p.tok.kind == T_DEFAULT {
		n.children = append(n.children, p.defaultStatement())
		p.nextToken()
	}

	// End switch
	if p.tok.kind != T_R_SQUIRLY {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "right squirly", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_R_SQUIRLY, data: p.tok.data, line: p.tok.line})

	return &n
}

func (p *Parser) caseStatement() *Node {
	const FUNC_NAME = "case statement"

	n := Node{kind: N_CASE_STATE, line: p.tok.line}

	if p.tok.kind != T_CASE {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "case", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_CASE, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	n.children = append(n.children, p.expression())
	p.nextToken()

	for p.tok.kind == T_SEP {
		n.children = append(n.children, &Node{kind: N_SEP, data: p.tok.data, line: p.tok.line})
		p.nextToken()

		n.children = append(n.children, p.expression())
		p.nextToken()
	}

	if p.tok.kind != T_COLON {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "colon", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_COLON, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	n.children = append(n.children, p.caseBlock())

	return &n
}

func (p *Parser) defaultStatement() *Node {
	const FUNC_NAME = "default statement"

	n := Node{kind: N_DEFAULT_STATE, line: p.tok.line}

	if p.tok.kind != T_DEFAULT {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "default", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_DEFAULT, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	if p.tok.kind != T_COLON {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "colon", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_COLON, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	n.children = append(n.children, p.caseBlock())

	return &n
}

// NOTE: Is also used for default
func (p *Parser) caseBlock() *Node {
	const FUNC_NAME = "case block"

	block := Node{kind: N_CASE_BLOCK, line: p.tok.line}

	var n *Node

	for p.tok.kind != T_ILLEGAL && p.tok.kind != T_CASE && p.tok.kind != T_DEFAULT && p.tok.kind != T_R_SQUIRLY {
		switch p.tok.kind {
		case T_IDENTIFIER, T_L_BLOCK, T_DEREF: // Variable declaration
			n = p.variableDeclaration()
			block.children = append(block.children, n)

		case T_INC, T_DINC: // Lone increment
			n = p.loneIncrement()
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

		case T_WHILE: // While loop
			n = p.whileLoop()
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

		case T_SWITCH: // Switch statement
			n = p.switchStatement()
			block.children = append(block.children, n)

		default:
			throwError(p.tok.line, fmt.Sprint("Bad start to statement: ", p.tok.kind, " on line ", p.tok.line))
		}

		p.nextToken()
	}

	// We end on a case, default, or
	// whatever, we want to move back one
	// so that the outer function can
	// discover it themselves
	p.index--

	return &block
}

func (p *Parser) ifBlock() *Node {
	const FUNC_NAME = "if block"

	n := Node{kind: N_IF_BLOCK, line: p.tok.line}

	if p.tok.kind != T_IF {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "if", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_IF, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	n.children = append(n.children, p.condition())
	p.nextToken()

	n.children = append(n.children, p.block())
	p.nextToken()

	for p.tok.kind == T_ELIF {
		n.children = append(n.children, &Node{kind: N_ELIF, data: p.tok.data, line: p.tok.line})
		p.nextToken()

		n.children = append(n.children, p.condition())
		p.nextToken()

		n.children = append(n.children, p.block())
		p.nextToken()
	}

	if p.tok.kind == T_ELSE {
		n.children = append(n.children, &Node{kind: N_ELSE, data: p.tok.data, line: p.tok.line})
		p.nextToken()

		n.children = append(n.children, p.block())
	} else {
		// Rather do this than mess around with
		// a whole bunch of peeks
		p.index--
	}

	return &n
}

func (p *Parser) foreverLoop() *Node {
	const FUNC_NAME = "forever loop"

	n := Node{kind: N_FOREVER_LOOP, line: p.tok.line}

	if p.tok.kind != T_FOREVER {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "forever", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_FOREVER, line: p.tok.line})
	p.nextToken()

	n.children = append(n.children, p.block())

	return &n
}

func (p *Parser) rangeLoop() *Node {
	const FUNC_NAME = "range loop"

	n := Node{kind: N_RANGE_LOOP, line: p.tok.line}

	if p.tok.kind != T_RANGE {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "range", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_RANGE, line: p.tok.line})
	p.nextToken()

	n.children = append(n.children, p.expression())
	p.nextToken()

	n.children = append(n.children, p.block())

	return &n
}

func (p *Parser) forLoop() *Node {
	const FUNC_NAME = "for loop"

	n := Node{kind: N_FOR_LOOP, line: p.tok.line}

	// for
	if p.tok.kind != T_FOR {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "for", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_FOR, line: p.tok.line})
	p.nextToken()

	// for i int = 0
	if p.tok.kind != T_SEMICOLON {
		n.children = append(n.children, p.assignment())
		p.nextToken()
	}

	// for i int = 0;
	if p.tok.kind != T_SEMICOLON {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "semicolon", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON, line: p.tok.line})
	p.nextToken()

	// for i int = 0; i < 10
	if p.tok.kind != T_SEMICOLON {
		n.children = append(n.children, p.expression())
		p.nextToken()
	}

	// for i int = 0; i < 10;
	if p.tok.kind != T_SEMICOLON {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "semicolon", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON, line: p.tok.line})
	p.nextToken()

	// for i int = 0; i < 10; i = i + 1
	if p.tok.kind != T_L_SQUIRLY {
		if p.tok.kind == T_INC || p.tok.kind == T_DINC {
			inc := Node{kind: N_UNARY_OPERATION, line: p.tok.line}

			if p.tok.kind == T_INC {
				inc.children = append(inc.children, &Node{kind: N_INC, data: p.tok.data, line: p.tok.line})
			} else if p.tok.kind == T_DINC {
				inc.children = append(inc.children, &Node{kind: N_DINC, data: p.tok.data, line: p.tok.line})
			} else {
				throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "inc or dinc", p.tok))
			}
			p.nextToken()

			if p.peekToken().kind == T_ACCESS {
				inc.children = append(inc.children, p.property())
			} else {
				inc.children = append(inc.children, &Node{kind: N_IDENTIFIER, data: p.tok.data, line: p.tok.line})
			}
			inc.children = append(inc.children, &Node{kind: N_IDENTIFIER, data: p.tok.data, line: p.tok.line})
			n.children = append(n.children, &inc)
		} else {
			n.children = append(n.children, p.assignment())
		}
		p.nextToken()
	}

	// for i int = 0; i < 10; i = i + 1 {}
	n.children = append(n.children, p.block())

	return &n
}

func (p *Parser) whileLoop() *Node {
	const FUNC_NAME = "while loop"

	n := Node{kind: N_WHILE_LOOP, line: p.tok.line}

	// while
	if p.tok.kind != T_WHILE {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "while", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_WHILE, line: p.tok.line})
	p.nextToken()

	// while i < 10
	n.children = append(n.children, p.condition())
	p.nextToken()

	// while i < 10 {}
	n.children = append(n.children, p.block())

	return &n
}

func (p *Parser) structDef() *Node {
	const FUNC_NAME = "struct definition"

	n := Node{kind: N_STRUCT_DEF, line: p.tok.line}

	if p.tok.kind != T_STRUCT {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "struct", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_STRUCT, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	// NOTE: The lexer has figured out for
	// the parser that the name for the
	// struct is a type, and therefore we
	// check for type (over identifier)
	if p.tok.kind != T_TYPE {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "type", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_TYPE, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	if p.tok.kind != T_L_SQUIRLY {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "left squirly", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_L_SQUIRLY, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	for p.tok.kind != T_R_SQUIRLY && p.tok.kind != T_ILLEGAL {
		if p.tok.kind != T_IDENTIFIER {
			throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok))
		}
		n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: p.tok.data, line: p.tok.line})
		p.nextToken()

		n.children = append(n.children, p.complexType())
		p.nextToken()
	}

	if p.tok.kind != T_R_SQUIRLY {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "right squirly", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_R_SQUIRLY, data: p.tok.data, line: p.tok.line})

	return &n
}

func (p *Parser) funcDef() *Node {
	const FUNC_NAME = "function definition"

	n := Node{kind: N_FUNC_DEF, line: p.tok.line}

	if p.tok.kind != T_FUN {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "fun", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_FUN, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	// Method on struct
	if p.tok.kind == T_L_PAREN {
		n.children = append(n.children, p.methodReceiver())
		p.nextToken()
	}

	if p.tok.kind != T_IDENTIFIER {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	if p.tok.kind != T_L_PAREN {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "left paren", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_L_PAREN, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	if p.tok.kind == T_IDENTIFIER {
		n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: p.tok.data, line: p.tok.line})
		p.nextToken()

		n.children = append(n.children, p.complexType())
		p.nextToken()

		for p.tok.kind == T_SEP {
			n.children = append(n.children, &Node{kind: N_SEP, data: p.tok.data, line: p.tok.line})
			p.nextToken()

			if p.tok.kind != T_IDENTIFIER {
				throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok))
			}
			n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: p.tok.data, line: p.tok.line})
			p.nextToken()

			n.children = append(n.children, p.complexType())
			p.nextToken()
		}
	}

	if p.tok.kind != T_R_PAREN {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "right paren", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_R_PAREN, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	// Return type?
	if p.tok.kind != T_L_SQUIRLY {
		n.children = append(n.children, p.complexType())
		p.nextToken()
	}

	n.children = append(n.children, p.block())

	return &n
}

func (p *Parser) methodReceiver() *Node {
	const FUNC_NAME = "method receiver"

	n := Node{kind: N_METHOD_RECEIVER, line: p.tok.line}

	if p.tok.kind != T_L_PAREN {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "left paren", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_L_PAREN, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	if p.tok.kind != T_IDENTIFIER {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	n.children = append(n.children, p.complexType())
	p.nextToken()

	if p.tok.kind != T_R_PAREN {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "right paren", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_R_PAREN, data: p.tok.data, line: p.tok.line})

	return &n
}

func (p *Parser) retStatement() *Node {
	const FUNC_NAME = "return statement"

	n := Node{kind: N_RET_STATE, line: p.tok.line}

	if p.tok.kind != T_RET {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "return", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_RET, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	// don't do the extra value
	if p.tok.kind == T_SEMICOLON {
		n.children = append(n.children, &Node{kind: N_SEMICOLON, data: p.tok.data, line: p.tok.line})
		return &n
	}

	n.children = append(n.children, p.expression())
	p.nextToken()

	if p.tok.kind != T_SEMICOLON {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "semicolon", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON, data: p.tok.data, line: p.tok.line})

	return &n
}

func (p *Parser) breakStatement() *Node {
	const FUNC_NAME = "break statement"

	n := Node{kind: N_BREAK_STATE, line: p.tok.line}

	if p.tok.kind != T_BREAK {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "break", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_BREAK, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	// don't do the extra value
	if p.tok.kind == T_SEMICOLON {
		n.children = append(n.children, &Node{kind: N_SEMICOLON, data: p.tok.data, line: p.tok.line})
		return &n
	}

	n.children = append(n.children, p.expression())
	p.nextToken()

	if p.tok.kind != T_SEMICOLON {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "semicolon", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON, data: p.tok.data, line: p.tok.line})

	return &n
}

func (p *Parser) contStatement() *Node {
	const FUNC_NAME = "continue statement"

	n := Node{kind: N_CONT_STATE, line: p.tok.line}

	if p.tok.kind != T_CONT {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "continue", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_CONT, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	// don't do the extra value
	// (same as "continue 0")
	if p.tok.kind == T_SEMICOLON {
		n.children = append(n.children, &Node{kind: N_SEMICOLON, data: p.tok.data, line: p.tok.line})
		return &n
	}

	n.children = append(n.children, p.expression())
	p.nextToken()

	if p.tok.kind != T_SEMICOLON {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "semicolon", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON, data: p.tok.data, line: p.tok.line})

	return &n
}

func (p *Parser) enumDef() *Node {
	const FUNC_NAME = "enum definition"

	n := Node{kind: N_ENUM_DEF, line: p.tok.line}

	if p.tok.kind != T_ENUM {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "enum", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_ENUM, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	if p.tok.kind != T_IDENTIFIER {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	if p.tok.kind != T_L_SQUIRLY {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "left squirly", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_L_SQUIRLY, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	for p.tok.kind != T_R_SQUIRLY && p.tok.kind != T_ILLEGAL {
		if p.tok.kind != T_IDENTIFIER {
			throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok))
		}
		n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: p.tok.data, line: p.tok.line})
		p.nextToken()

		if p.tok.kind != T_SEP {
			throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "separator", p.tok))
		}
		n.children = append(n.children, &Node{kind: N_SEP, data: p.tok.data, line: p.tok.line})
		p.nextToken()
	}

	if p.tok.kind != T_R_SQUIRLY {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "right squirly", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_R_SQUIRLY, data: p.tok.data, line: p.tok.line})

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

	n := Node{kind: N_EXPRESSION, line: p.tok.line}

	// 7
	n.children = append(n.children, p.value())
	p.nextToken()

	v = p.operator()
	for v != nil {
		n.children = append(n.children, v)
		p.nextToken()
		n.children = append(n.children, p.value())
		p.nextToken()
		v = p.operator()
	}

	p.index--

	return &n
}

func (p *Parser) assignment() *Node {
	const FUNC_NAME = "assignment"

	n := Node{kind: N_ASSIGNMENT, line: p.tok.line}

	// x
	if p.tok.kind != T_IDENTIFIER {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok))
	}

	if p.peekToken().kind == T_ACCESS {
		n.children = append(n.children, p.property())
		p.nextToken()
	} else {
		n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: p.tok.data, line: p.tok.line})
		p.nextToken()
	}

	isMap := p.tok.kind == T_MAP

	// We have a type?
	if p.tok.kind != T_ASSIGN {
		n.children = append(n.children, p.complexType())

		if isMap {
			return &n
		}

		p.nextToken()
	}

	// Now we MUST have an assign
	if p.tok.kind != T_ASSIGN {
		if p.tok.kind != T_SEMICOLON {
			throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "assign or semicolon", p.tok))
		}

		p.index--

		return &n
	}

	n.children = append(n.children, &Node{kind: N_ASSIGN, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	// x int = 3 + 7
	n.children = append(n.children, p.expression())

	return &n
}

func (p *Parser) loneCall() *Node {
	const FUNC_NAME = "lone call"

	n := Node{kind: N_LONE_CALL, line: p.tok.line}

	n.children = append(n.children, p.funcCall())
	p.nextToken()

	if p.tok.kind != T_SEMICOLON {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "semicolon", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON, data: p.tok.data, line: p.tok.line})

	return &n
}

func (p *Parser) funcCall() *Node {
	const FUNC_NAME = "function call"

	n := Node{kind: N_FUNC_CALL, line: p.tok.line}

	if p.tok.kind != T_CALL {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "call", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_CALL, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	if p.tok.kind != T_IDENTIFIER {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok))
	}

	if p.peekToken().kind == T_ACCESS {
		n.children = append(n.children, p.property())
		p.nextToken()
	} else {
		n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: p.tok.data, line: p.tok.line})
		p.nextToken()
	}

	if p.tok.kind != T_L_PAREN {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "left paren", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_L_PAREN, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	if p.tok.kind != T_R_PAREN {
		n.children = append(n.children, p.expression())
		p.nextToken()

		for p.tok.kind == T_SEP {
			n.children = append(n.children, &Node{kind: N_SEP, data: p.tok.data, line: p.tok.line})
			p.nextToken()

			n.children = append(n.children, p.expression())
			p.nextToken()
		}
	}

	if p.tok.kind != T_R_PAREN {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "right paren", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_R_PAREN, data: p.tok.data, line: p.tok.line})

	return &n
}

func (p *Parser) structNew() *Node {
	const FUNC_NAME = "new struct"

	n := Node{kind: N_STRUCT_NEW, line: p.tok.line}

	if p.tok.kind != T_NEW {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "new", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_NEW, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	if p.tok.kind != T_TYPE {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "type", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_TYPE, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	if p.tok.kind != T_L_PAREN {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "left paren", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_L_PAREN, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	if p.tok.kind != T_R_PAREN {
		n.children = append(n.children, p.expression())
		p.nextToken()

		for p.tok.kind == T_SEP {
			n.children = append(n.children, &Node{kind: N_SEP, data: p.tok.data, line: p.tok.line})
			p.nextToken()

			n.children = append(n.children, p.expression())
			p.nextToken()
		}
	}

	if p.tok.kind != T_R_PAREN {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "right paren", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_R_PAREN, data: p.tok.data, line: p.tok.line})

	return &n
}

func (p *Parser) block() *Node {
	const FUNC_NAME = "block"

	block := Node{kind: N_BLOCK, line: p.tok.line}

	var n *Node

	if p.tok.kind != T_L_SQUIRLY {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "left squirly", p.tok))
	}
	block.children = append(block.children, &Node{kind: N_L_SQUIRLY, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	for p.tok.kind != T_ILLEGAL && p.tok.kind != T_R_SQUIRLY {
		switch p.tok.kind {
		case T_IDENTIFIER, T_L_BLOCK, T_DEREF: // Variable declaration
			n = p.variableDeclaration()
			block.children = append(block.children, n)

		case T_INC, T_DINC: // Lone increment
			n = p.loneIncrement()
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

		case T_WHILE: // While loop
			n = p.whileLoop()
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

		case T_SWITCH: // Switch statement
			n = p.switchStatement()
			block.children = append(block.children, n)

		default:
			throwError(p.tok.line, fmt.Sprint("Bad start to statement: ", p.tok.kind, " on line ", p.tok.line))
		}

		p.nextToken()
	}

	if p.tok.kind != T_R_SQUIRLY {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "right squirly", p.tok))
	}
	block.children = append(block.children, &Node{kind: N_R_SQUIRLY, data: p.tok.data, line: p.tok.line})

	return &block
}

func (p *Parser) typeDef() *Node {
	const FUNC_NAME = "type definition"

	n := Node{kind: N_NEW_TYPE, line: p.tok.line}

	if p.tok.kind != T_TYPEDEF {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "typedef", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_TYPEDEF, line: p.tok.line})
	p.nextToken()

	if p.tok.kind != T_TYPE {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "type", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_TYPE, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	if p.tok.kind != T_TYPE {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "type", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_TYPE, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	if p.tok.kind != T_SEMICOLON {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "semicolon", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_SEMICOLON, line: p.tok.line})

	return &n
}

// TODO: Dealing with unary -?
// TODO: What happens with p.x (accessors)
func (p *Parser) value() *Node {
	const FUNC_NAME = "value"

	var n Node

	switch p.tok.kind {
	case T_IDENTIFIER:
		if p.peekToken().kind == T_ACCESS {
			return p.property()
		} else {
			return &Node{kind: N_IDENTIFIER, data: p.tok.data, line: p.tok.line}
		}
	case T_INT:
		return &Node{kind: N_INT, data: p.tok.data, line: p.tok.line}
	case T_FLOAT:
		return &Node{kind: N_FLOAT, data: p.tok.data, line: p.tok.line}
	case T_STRING:
		return &Node{kind: N_STRING, data: p.tok.data, line: p.tok.line}
	case T_CHAR:
		return &Node{kind: N_CHAR, data: p.tok.data, line: p.tok.line}
	case T_BOOL:
		return &Node{kind: N_BOOL, data: p.tok.data, line: p.tok.line}
	case T_NIL:
		return &Node{kind: N_NIL, data: p.tok.data, line: p.tok.line}

		// Unary cases
	case T_NOT:
		n = Node{kind: N_UNARY_OPERATION, data: p.tok.data, line: p.tok.line}
		n.children = append(n.children, &Node{kind: N_NOT, data: p.tok.data, line: p.tok.line})
		p.nextToken()

		n.children = append(n.children, p.value())
		return &n

	case T_INC:
		n = Node{kind: N_UNARY_OPERATION, data: p.tok.data, line: p.tok.line}
		n.children = append(n.children, &Node{kind: N_INC, data: p.tok.data, line: p.tok.line})
		p.nextToken()

		n.children = append(n.children, p.value())
		return &n

	case T_DINC:
		n = Node{kind: N_UNARY_OPERATION, data: p.tok.data, line: p.tok.line}
		n.children = append(n.children, &Node{kind: N_DINC, data: p.tok.data, line: p.tok.line})
		p.nextToken()

		n.children = append(n.children, p.value())
		return &n

	case T_REF:
		n = Node{kind: N_UNARY_OPERATION, data: p.tok.data, line: p.tok.line}
		n.children = append(n.children, &Node{kind: N_REF, data: p.tok.data, line: p.tok.line})
		p.nextToken()

		n.children = append(n.children, p.value())
		return &n

	case T_DEREF:
		n = Node{kind: N_UNARY_OPERATION, data: p.tok.data, line: p.tok.line}
		n.children = append(n.children, &Node{kind: N_DEREF, data: p.tok.data, line: p.tok.line})
		p.nextToken()

		n.children = append(n.children, p.value())
		return &n

	case T_L_BLOCK:
		n = Node{kind: N_UNARY_OPERATION, data: p.tok.data, line: p.tok.line}
		n.children = append(n.children, p.indexUnary())
		p.nextToken()

		n.children = append(n.children, p.value())
		return &n

		// We're using an array or slice
	case T_MAKE:
		return p.makeArray()

		// Someones trying to do some bedmas
	case T_L_PAREN:
		n = Node{kind: N_BRACKETED_VALUE, line: p.tok.line}

		n.children = append(n.children, &Node{kind: N_L_PAREN, data: p.tok.data, line: p.tok.line})
		p.nextToken()

		n.children = append(n.children, p.expression())
		p.nextToken()

		if p.tok.kind != T_R_PAREN {
			throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "right paren", p.tok))
		}
		n.children = append(n.children, &Node{kind: N_R_PAREN, data: p.tok.data, line: p.tok.line})
		return &n

		// Calling a function to use as a value
	case T_CALL:
		return p.funcCall()

		// Creating a new struct
	case T_NEW:
		return p.structNew()

	default:
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "unary or value", p.tok))

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

	switch p.tok.kind {
	case T_ADD:
		return &Node{kind: N_ADD, data: p.tok.data, line: p.tok.line}
	case T_XOR:
		return &Node{kind: N_XOR, data: p.tok.data, line: p.tok.line}
	case T_ACCESS:
		return &Node{kind: N_ACCESS, data: p.tok.data, line: p.tok.line}
	case T_NEQ:
		return &Node{kind: N_NEQ, data: p.tok.data, line: p.tok.line}
	case T_MOD:
		return &Node{kind: N_MOD, data: p.tok.data, line: p.tok.line}
	case T_EQ:
		return &Node{kind: N_EQ, data: p.tok.data, line: p.tok.line}
	case T_LT:
		return &Node{kind: N_LT, data: p.tok.data, line: p.tok.line}
	case T_GT:
		return &Node{kind: N_GT, data: p.tok.data, line: p.tok.line}
	case T_LTEQ:
		return &Node{kind: N_LTEQ, data: p.tok.data, line: p.tok.line}
	case T_GTEQ:
		return &Node{kind: N_GTEQ, data: p.tok.data, line: p.tok.line}
	case T_SUB:
		return &Node{kind: N_SUB, data: p.tok.data, line: p.tok.line}
	case T_MUL:
		return &Node{kind: N_MUL, data: p.tok.data, line: p.tok.line}
	case T_DIV:
		return &Node{kind: N_DIV, data: p.tok.data, line: p.tok.line}
	case T_OR:
		return &Node{kind: N_OR, data: p.tok.data, line: p.tok.line}
	case T_AND:
		return &Node{kind: N_AND, data: p.tok.data, line: p.tok.line}
	case T_OROR:
		return &Node{kind: N_OROR, data: p.tok.data, line: p.tok.line}
	case T_ANDAND:
		return &Node{kind: N_ANDAND, data: p.tok.data, line: p.tok.line}
	case T_L_SHIFT:
		return &Node{kind: N_L_SHIFT, data: p.tok.data, line: p.tok.line}
	case T_R_SHIFT:
		return &Node{kind: N_R_SHIFT, data: p.tok.data, line: p.tok.line}

	default:
		return nil
	}
}

func (p *Parser) indexUnary() *Node {
	const FUNC_NAME = "index unary"

	n := Node{kind: N_INDEX, line: p.tok.line}

	if p.tok.kind != T_L_BLOCK {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "left block", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_L_BLOCK, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	n.children = append(n.children, p.expression())
	p.nextToken()

	if p.tok.kind != T_R_BLOCK {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "right block", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_R_BLOCK, data: p.tok.data, line: p.tok.line})

	return &n
}

func (p *Parser) makeArray() *Node {
	const FUNC_NAME = "make array"

	n := Node{kind: N_MAKE_ARRAY, line: p.tok.line}

	if p.tok.kind != T_MAKE {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "make", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_MAKE, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	if p.tok.kind != T_L_BLOCK {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "left block", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_L_BLOCK, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	for p.tok.kind != T_R_BLOCK {
		n.children = append(n.children, p.expression())
		p.nextToken()

		if p.tok.kind != T_SEP {
			break // Last element
		}
		n.children = append(n.children, &Node{kind: N_SEP, data: p.tok.data, line: p.tok.line})
		p.nextToken()
	}

	if p.tok.kind != T_R_BLOCK {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "right block", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_R_BLOCK, data: p.tok.data, line: p.tok.line})

	return &n
}

func (p *Parser) complexType() *Node {
	const FUNC_NAME = "complex type"

	n := Node{kind: N_COMPLEX_TYPE, line: p.tok.line}

	// First check if we're making a map
	if p.tok.kind == T_MAP {
		n.children = append(n.children, &Node{kind: N_MAP, data: p.tok.data, line: p.tok.line})
		p.nextToken()

		if p.tok.kind != T_TYPE {
			throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "type", p.tok))
		}
		n.children = append(n.children, &Node{kind: N_TYPE, data: p.tok.data, line: p.tok.line})
		p.nextToken()

		if p.tok.kind != T_L_BLOCK {
			throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "left block", p.tok))
		}
		n.children = append(n.children, &Node{kind: N_L_BLOCK, data: p.tok.data, line: p.tok.line})
		p.nextToken()

		// NOTE: The logic for determining
		// if an identifier inside the brackets
		// is a type or not isn't so clean, so
		// we left it as an identifier in the
		// lexer, but here, we can finally
		// change it to a type
		if p.tok.kind != T_IDENTIFIER && p.tok.kind != T_TYPE {
			throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "identifier (but really a type)", p.tok))
		}
		n.children = append(n.children, &Node{kind: N_TYPE, data: p.tok.data, line: p.tok.line})
		p.nextToken()

		if p.tok.kind != T_R_BLOCK {
			throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "right block", p.tok))
		}
		n.children = append(n.children, &Node{kind: N_R_BLOCK, data: p.tok.data, line: p.tok.line})

		// Can't assign to map straight away
		// (for now because bruh)
		return &n

		// Otherwise we might have a more normal type
	} else if p.tok.kind == T_TYPE || p.tok.kind == T_IDENTIFIER {
		n.children = append(n.children, &Node{kind: N_TYPE, data: p.tok.data, line: p.tok.line})

		pt := p.peekToken()

		// Is it a pointer type?
		if pt.kind == T_DEREF {
			p.nextToken()

			n.children = append(n.children, &Node{kind: N_DEREF, data: p.tok.data, line: p.tok.line})
			pt = p.peekToken()
		}

		if pt.kind == T_L_BLOCK { // Some larger type
			p.nextToken()

			// ArrayList
			if p.peekToken().kind == T_R_BLOCK {
				n.children = append(n.children, p.emptyBlock())
			} else { // Array
				n.children = append(n.children, p.indexUnary())
			}
		}
	}

	return &n
}

func (p *Parser) emptyBlock() *Node {
	const FUNC_NAME = "empty block"

	n := Node{kind: N_EMPTY_BLOCK, line: p.tok.line}

	if p.tok.kind != T_L_BLOCK {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "left block", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_L_BLOCK, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	if p.tok.kind != T_R_BLOCK {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "right block", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_R_BLOCK, data: p.tok.data, line: p.tok.line})

	return &n
}

func (p *Parser) property() *Node {
	const FUNC_NAME = "property"

	n := Node{kind: N_PROPERTY, line: p.tok.line}

	if p.tok.kind != T_IDENTIFIER {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	if p.tok.kind != T_ACCESS {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "access", p.tok))
	}
	n.children = append(n.children, &Node{kind: N_ACCESS, data: p.tok.data, line: p.tok.line})
	p.nextToken()

	if p.tok.kind != T_IDENTIFIER {
		throwError(p.tok.line, p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok))
	}

	if p.peekToken().kind == T_ACCESS {
		n.children = append(n.children, p.property())
	} else {
		n.children = append(n.children, &Node{kind: N_IDENTIFIER, data: p.tok.data, line: p.tok.line})
	}

	return &n
}

func (p *Parser) throwError(caller string, line int, expected string, got Token) string {
	return "Error in the " + JOB_PARSER + "!\n" +
		"When the " + JOB_PARSER + " was trying to decipher: " + caller + "\n" +
		"Error found on line " + strconv.Itoa(line) + "\n" +
		"Expected: " + expected + "\n" +
		"Got: " + fmt.Sprint(got)

}
