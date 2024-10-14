package main

import "strconv"

type Token struct {
	// What is looks like in the actual source
	data string

	// What the token is depicting
	kind TokenType
}

func (t Token) String() string {
	return "(" + t.data + " " + strconv.Itoa(int(t.kind)) + ")"
}

type TokenType byte

const (
	T_ILLEGAL TokenType = iota

	// Keywords
	T_FOR
	T_RANGE
	T_FOREVER
	T_IF
	T_ELIF
	T_ELSE
	T_CALL
	T_STRUCT
	T_FUN
	T_RET
	T_BREAK
	T_CONT
	T_ENUM
	T_TYPEDEF

	// Various symbols
	T_SEMICOLON
	T_ASSIGN
	T_SEP

	// Operators
	T_ADD
	T_SUB
	T_MUL
	T_DIV
	T_OR
	T_AND
	T_INC
	T_DINC
	T_EQ
	T_LT
	T_GT
	T_LTEQ
	T_GTEQ
	T_NOT
	T_NEQ
	T_MOD
	T_REF
	T_DEREF
	T_XOR
	T_ACCESS

	// Paired symbols
	T_L_SQUIRLY
	T_R_SQUIRLY
	T_L_BLOCK
	T_R_BLOCK
	T_L_PAREN
	T_R_PAREN

	// Values
	T_TYPE
	T_IDENTIFIER
	T_INT
	T_FLOAT
	T_STRING
	T_CHAR
	T_BOOL
	T_NIL
)
