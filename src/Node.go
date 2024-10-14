package main

// Tokens get converted into nodes,
// with a slight twist in that they
// have children, creating a tree
// structure
type Node struct {
	kind     NodeType
	children []*Node
	data     string
}

type NodeType int

const (
	N_ILLEGAL TokenType = iota

	// Higher level structures
	N_PROGRAM         // Top level structure, holds everything
	N_VAR_DECLARATION // Declaration (or def) of a variable
	N_IF_BLOCK        // If (elif else) block
	N_FOREVER_LOOP    // Loop with no end
	N_RANGE_LOOP      // Loop with fixed end
	N_FOR_LOOP        // Boomer loop
	N_CALL_STATE      // Statement that is just a call to a function
	N_STRUCT_DEF      // Definition of a struct
	N_FUNC_DEF        // Definition of a function
	N_RET_STATE       // Return
	N_BREAK_STATE     // Break
	N_CONT_STATE      // Continue
	N_ENUM_DEF        // Enum definition
	N_CONDITION       // Expression that returns a bool
	N_EXPRESSION      // Mathematical operation, etc
	N_ASSIGNMENT      // An identifier being assigned a value
	N_FUNC_CALL       // Calling a functions
	N_BLOCK           // A bunch of statements

	// Keywords
	N_FOR
	N_RANGE
	N_FOREVER
	N_IF
	N_ELIF
	N_ELSE
	N_CALL
	N_STRUCT
	N_FUN
	N_RET
	N_BREAK
	N_CONT
	N_ENUM
	N_TYPEDEF

	// Various symbols
	N_SEMICOLON
	N_ASSIGN
	N_SEP
	N_REF
	N_DEREF
	N_XOR

	// Paired symbols
	N_L_SQUIRLY
	N_R_SQUIRLY
	N_L_BLOCK
	N_R_BLOCK
	N_L_PAREN
	N_R_PAREN

	// Operators
	N_ADD
	N_SUB
	N_MUL
	N_DIV
	N_OR
	N_AND
	N_INC
	N_DINC
	N_EQ
	N_LT
	N_GT
	N_LTEQ
	N_GTEQ
	N_NOT
	N_NEQ
	N_MOD
	N_ACCESS

	// Values
	N_TYPE
	N_IDENTIFIER
	N_INT
	N_FLOAT
	N_STRING
	N_CHAR
	N_BOOL
	N_NIL
)
