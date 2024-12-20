package main

// Tokens get converted into nodes,
// with a slight twist in that they
// have children, creating a tree
// structure
type Node struct {
	kind     NodeType
	children []*Node
	data     string
	line     int
}

type NodeType int

func (n *Node) String() string {
	output := n.kind.String() + ": " + n.data + "\n"

	for i := range n.children {
		output += n.children[i].StringRec(1)
	}

	return output
}

func (n *Node) StringRec(indent int) string {
	output := getIndent(indent) + n.kind.String() + ": " + n.data + "\n"

	for i := range n.children {
		output += n.children[i].StringRec(indent + 1)
	}

	return output
}

func getIndent(indent int) string {
	output := ""

	for range indent {
		output += "\t"
	}

	return output
}

const (
	N_ILLEGAL NodeType = iota

	// Higher level structures
	N_PROGRAM         // Top level structure, holds everything
	N_VAR_DECLARATION // Declaration (or def) of a variable
	N_IF_BLOCK        // If (elif else) block
	N_FOREVER_LOOP    // Loop with no end
	N_RANGE_LOOP      // Loop with fixed end
	N_FOR_LOOP        // Boomer loop
	N_WHILE_LOOP      // Loop with only a condition (can be recreated with for loops)
	N_STRUCT_DEF      // Definition of a struct
	N_FUNC_DEF        // Definition of a function
	N_RET_STATE       // Return
	N_BREAK_STATE     // Break
	N_CONT_STATE      // Continue
	N_ENUM_DEF        // Enum definition
	N_CONDITION       // Expression that returns a bool
	N_EXPRESSION      // Mathematical operation, etc
	N_ASSIGNMENT      // An identifier being assigned a value
	N_LONE_CALL       // A lone function call as a statement
	N_FUNC_CALL       // Calling a functions
	N_STRUCT_NEW      // Creating a new struct
	N_BLOCK           // A bunch of statements
	N_NEW_TYPE        // Typdef (the statement)
	N_UNARY_OPERATION // When you use a ref, deref, not, etc
	N_BRACKETED_VALUE // An expression in parens (for bedmas purposes)
	N_MAKE_ARRAY      // The syntax of creating an array
	N_COMPLEX_TYPE    // A node that can encompass any type
	N_SWITCH_STATE    // A whole switch statement
	N_CASE_STATE      // The case line in a switch
	N_DEFAULT_STATE   // The default line in a switch
	N_CASE_BLOCK      // The block that is contained within a switch statement
	N_LONE_INC        // Can also be a lone decrement
	N_METHOD_RECEIVER // Brackets before a function name that shows that it's a method
	N_EMPTY_BLOCK     // Empty block brackets for denoting slices
	N_PROPERTY        // A chain of accesses
	N_CONSTANT        // A constant variable declaration

	// Keywords
	N_FOR
	N_RANGE
	N_FOREVER
	N_WHILE
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
	N_NEW
	N_MAKE
	N_MAP
	N_SWITCH
	N_CASE
	N_DEFAULT
	N_CONST

	// Various symbols
	N_SEMICOLON
	N_ASSIGN
	N_SEP
	N_COLON

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
	N_OROR
	N_ANDAND
	N_EQ
	N_LT
	N_GT
	N_LTEQ
	N_GTEQ
	N_NEQ
	N_MOD
	N_ACCESS
	N_XOR
	N_L_SHIFT
	N_R_SHIFT

	// Unary
	N_INC
	N_DINC
	N_NOT
	N_REF
	N_DEREF
	N_INDEX

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

func (n NodeType) String() string {
	switch n {
	case N_ILLEGAL:
		return "ILLEGAL"
	case N_PROGRAM:
		return "PROGRAM"
	case N_VAR_DECLARATION:
		return "VAR_DECLARATION"
	case N_IF_BLOCK:
		return "IF_BLOCK"
	case N_FOREVER_LOOP:
		return "FOREVER_LOOP"
	case N_RANGE_LOOP:
		return "RANGE_LOOP"
	case N_FOR_LOOP:
		return "FOR_LOOP"
	case N_WHILE_LOOP:
		return "WHILE_LOOP"
	case N_STRUCT_DEF:
		return "STRUCT_DEF"
	case N_FUNC_DEF:
		return "FUNC_DEF"
	case N_RET_STATE:
		return "RET_STATE"
	case N_BREAK_STATE:
		return "BREAK_STATE"
	case N_CONT_STATE:
		return "CONT_STATE"
	case N_ENUM_DEF:
		return "ENUM_DEF"
	case N_CONDITION:
		return "CONDITION"
	case N_EXPRESSION:
		return "EXPRESSION"
	case N_ASSIGNMENT:
		return "ASSIGNMENT"
	case N_FUNC_CALL:
		return "FUNC_CALL"
	case N_LONE_CALL:
		return "N_LONE_CALL"
	case N_STRUCT_NEW:
		return "STRUCT_NEW"
	case N_BLOCK:
		return "BLOCK"
	case N_UNARY_OPERATION:
		return "UNARY_OPERATION"
	case N_BRACKETED_VALUE:
		return "BRACKETED_VALUE"
	case N_NEW_TYPE:
		return "NEW_TYPE"
	case N_MAKE_ARRAY:
		return "MAKE_ARRAY"
	case N_COMPLEX_TYPE:
		return "COMPLEX_TYPE"
	case N_SWITCH_STATE:
		return "SWITCH_STATE"
	case N_CASE_STATE:
		return "CASE_STATE"
	case N_DEFAULT_STATE:
		return "DEFAULT_STATE"
	case N_CASE_BLOCK:
		return "CASE_BLOCK"
	case N_LONE_INC:
		return "LONE_INC"
	case N_METHOD_RECEIVER:
		return "METHOD_RECEIVER"
	case N_EMPTY_BLOCK:
		return "EMPTY_BLOCK"
	case N_PROPERTY:
		return "PROPERTY"
	case N_CONSTANT:
		return "CONSTANT"

	// Keywords
	case N_FOR:
		return "FOR"
	case N_RANGE:
		return "RANGE"
	case N_FOREVER:
		return "FOREVER"
	case N_WHILE:
		return "WHILE"
	case N_IF:
		return "IF"
	case N_ELIF:
		return "ELIF"
	case N_ELSE:
		return "ELSE"
	case N_CALL:
		return "CALL"
	case N_STRUCT:
		return "STRUCT"
	case N_FUN:
		return "FUN"
	case N_RET:
		return "RET"
	case N_BREAK:
		return "BREAK"
	case N_CONT:
		return "CONT"
	case N_ENUM:
		return "ENUM"
	case N_TYPEDEF:
		return "TYPEDEF"
	case N_NEW:
		return "NEW"
	case N_MAKE:
		return "MAKE"
	case N_MAP:
		return "MAP"
	case N_SWITCH:
		return "SWITCH"
	case N_CASE:
		return "CASE"
	case N_DEFAULT:
		return "DEFAULT"
	case N_CONST:
		return "CONST"

	// Paired symbols
	case N_L_SQUIRLY:
		return "L_SQUIRLY"
	case N_R_SQUIRLY:
		return "R_SQUIRLY"
	case N_L_BLOCK:
		return "L_BLOCK"
	case N_R_BLOCK:
		return "R_BLOCK"
	case N_L_PAREN:
		return "L_PAREN"
	case N_R_PAREN:
		return "R_PAREN"

	// Various symbols
	case N_SEMICOLON:
		return "SEMICOLON"
	case N_ASSIGN:
		return "ASSIGN"
	case N_SEP:
		return "SEP"
	case N_COLON:
		return "COLON"

	// Operator
	case N_XOR:
		return "XOR"
	case N_ADD:
		return "ADD"
	case N_SUB:
		return "SUB"
	case N_MUL:
		return "MUL"
	case N_DIV:
		return "DIV"
	case N_OR:
		return "OR"
	case N_AND:
		return "AND"
	case N_OROR:
		return "OROR"
	case N_ANDAND:
		return "ANDAND"
	case N_EQ:
		return "EQ"
	case N_LT:
		return "LT"
	case N_GT:
		return "GT"
	case N_LTEQ:
		return "LTEQ"
	case N_GTEQ:
		return "GTEQ"
	case N_L_SHIFT:
		return "N_L_SHIFT"
	case N_R_SHIFT:
		return "N_R_SHIFT"
	case N_NEQ:
		return "NEQ"
	case N_MOD:
		return "MOD"
	case N_ACCESS:
		return "ACCESS"

	// Unary
	case N_REF:
		return "REF"
	case N_DEREF:
		return "DEREF"
	case N_INDEX:
		return "INDEX"
	case N_NOT:
		return "NOT"
	case N_INC:
		return "INC"
	case N_DINC:
		return "DINC"

	case N_TYPE:
		return "TYPE"
	case N_IDENTIFIER:
		return "IDENTIFIER"
	case N_INT:
		return "INT"
	case N_FLOAT:
		return "FLOAT"
	case N_STRING:
		return "STRING"
	case N_CHAR:
		return "CHAR"
	case N_BOOL:
		return "BOOL"
	case N_NIL:
		return "NIL"
	default:
		return "UNKNOWN"
	}
}
