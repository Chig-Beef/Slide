package main

import "fmt"

type Parser struct {
	source []Token
	index  int
}

// The main program
func (p *Parser) parse() *Node {
	program := Node{kind: N_PROGRAM}

	switch p.source[p.index].kind {
	case T_IDENTIFIER: // Variable declaration
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

	default:
		panic(fmt.Sprint("Bad start to statement:", p.source[p.index]))
	}

	return &program
}
