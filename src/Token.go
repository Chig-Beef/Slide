package main

type Token struct {
	// What is looks like in the actual source
	data string

	// What the token is depicting
	kind TokenType

	// What line this token comes from
	line int
}

func (t Token) String() string {
	return "(" + t.data + " " + t.kind.String() + ")"
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
	T_NEW
	T_MAKE
	T_MAP

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
	T_EQ
	T_LT
	T_GT
	T_LTEQ
	T_GTEQ
	T_NEQ
	T_MOD
	T_XOR
	T_ACCESS
	T_NOT
	T_INC
	T_DINC
	T_REF
	T_DEREF

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

func (t TokenType) String() string {
	switch t {
	case T_ILLEGAL:
		return "ILLEGAL"

		// Keywords
	case T_FOR:
		return "FOR"
	case T_RANGE:
		return "RANGE"
	case T_FOREVER:
		return "FOREVER"
	case T_IF:
		return "IF"
	case T_ELIF:
		return "ELIF"
	case T_ELSE:
		return "ELSE"
	case T_CALL:
		return "CALL"
	case T_STRUCT:
		return "STRUCT"
	case T_FUN:
		return "FUN"
	case T_RET:
		return "RET"
	case T_BREAK:
		return "BREAK"
	case T_CONT:
		return "CONT"
	case T_ENUM:
		return "ENUM"
	case T_TYPEDEF:
		return "TYPEDEF"
	case T_NEW:
		return "NEW"
	case T_MAKE:
		return "MAKE"
	case T_MAP:
		return "MAP"

	case T_SEMICOLON:
		return "SEMICOLON"
	case T_ASSIGN:
		return "ASSIGN"
	case T_SEP:
		return "SEP"
	case T_ADD:
		return "ADD"
	case T_SUB:
		return "SUB"
	case T_MUL:
		return "MUL"
	case T_DIV:
		return "DIV"
	case T_OR:
		return "OR"
	case T_AND:
		return "AND"
	case T_EQ:
		return "EQ"
	case T_LT:
		return "LT"
	case T_GT:
		return "GT"
	case T_LTEQ:
		return "LTEQ"
	case T_GTEQ:
		return "GTEQ"
	case T_NEQ:
		return "NEQ"
	case T_MOD:
		return "MOD"
	case T_XOR:
		return "XOR"
	case T_ACCESS:
		return "ACCESS"
	case T_NOT:
		return "NOT"
	case T_INC:
		return "INC"
	case T_DINC:
		return "DINC"
	case T_REF:
		return "REF"
	case T_DEREF:
		return "DEREF"
	case T_L_SQUIRLY:
		return "L_SQUIRLY"
	case T_R_SQUIRLY:
		return "R_SQUIRLY"
	case T_L_BLOCK:
		return "L_BLOCK"
	case T_R_BLOCK:
		return "R_BLOCK"
	case T_L_PAREN:
		return "L_PAREN"
	case T_R_PAREN:
		return "R_PAREN"
	case T_TYPE:
		return "TYPE"
	case T_IDENTIFIER:
		return "IDENTIFIER"
	case T_INT:
		return "INT"
	case T_FLOAT:
		return "FLOAT"
	case T_STRING:
		return "STRING"
	case T_CHAR:
		return "CHAR"
	case T_BOOL:
		return "BOOL"
	case T_NIL:
		return "NIL"
	default:
		return "UNKNOWN"
	}
}
