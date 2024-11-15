package main

type Lexer struct {
	source []byte
	index  int
	line   int
}

type TokenType int

const (
	T_ILLEGAL TokenType = iota
	T_FOR
	T_RANGE
	T_FOREVER
	T_WHILE
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
	T_SWITCH
	T_CASE
	T_DEFAULT
	T_CONST
	T_SEMICOLON
	T_ASSIGN
	T_SEP
	T_COLON
	T_ADD
	T_SUB
	T_MUL
	T_DIV
	T_OR
	T_AND
	T_OROR
	T_ANDAND
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
	T_L_SHIFT
	T_R_SHIFT
	T_L_SQUIRLY
	T_R_SQUIRLY
	T_L_BLOCK
	T_R_BLOCK
	T_L_PAREN
	T_R_PAREN
	T_TYPE
	T_IDENTIFIER
	T_INT
	T_FLOAT
	T_STRING
	T_CHAR
	T_BOOL
	T_NIL
)

type Token struct {
	data string
	kind TokenType
	line int
}

type NodeType int

const (
	N_ILLEGAL NodeType = iota
	N_PROGRAM
	N_VAR_DECLARATION
	N_ELEMENT_ASSIGNMENT
	N_IF_BLOCK
	N_FOREVER_LOOP
	N_RANGE_LOOP
	N_FOR_LOOP
	N_WHILE_LOOP
	N_STRUCT_DEF
	N_FUNC_DEF
	N_RET_STATE
	N_BREAK_STATE
	N_CONT_STATE
	N_ENUM_DEF
	N_CONDITION
	N_EXPRESSION
	N_ASSIGNMENT
	N_LONE_CALL
	N_FUNC_CALL
	N_STRUCT_NEW
	N_BLOCK
	N_NEW_TYPE
	N_UNARY_OPERATION
	N_BRACKETED_VALUE
	N_MAKE_ARRAY
	N_COMPLEX_TYPE
	N_SWITCH_STATE
	N_CASE_STATE
	N_DEFAULT_STATE
	N_CASE_BLOCK
	N_LONE_INC
	N_METHOD_RECEIVER
	N_EMPTY_BLOCK
	N_PROPERTY
	N_CONSTANT
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
	N_SEMICOLON
	N_ASSIGN
	N_SEP
	N_COLON
	N_L_SQUIRLY
	N_R_SQUIRLY
	N_L_BLOCK
	N_R_BLOCK
	N_L_PAREN
	N_R_PAREN
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
	N_INC
	N_DINC
	N_NOT
	N_REF
	N_DEREF
	N_INDEX
	N_TYPE
	N_IDENTIFIER
	N_INT
	N_FLOAT
	N_STRING
	N_CHAR
	N_BOOL
	N_NIL
)

type Node struct {
	kind     NodeType
	children []*Node
	data     string
	line     int
}

func (l *Lexer) peekChar() byte {
	if l.index >= len(l.source)-1 {
		return 0
	}
	return l.source[l.index+1]
}

func (l *Lexer) lex() []Token {
	var tokens []Token
	var token Token
	for ; l.index < len(l.source); l.index++ {
		token = Token{}
		switch l.source[l.index] {
		case ' ':
			continue
		case '\t':
			continue
		case '\n':
			l.line++
			continue
		case '\r':
			continue
		case ';':
			token = Token{";", T_SEMICOLON, l.line}
		case '.':
			token = Token{".", T_ACCESS, l.line}
		case ':':
			token = Token{":", T_COLON, l.line}
		case '~':
			token = Token{"~", T_XOR, l.line}
		case '`':
			token = Token{"`", T_REF, l.line}
		case '^':
			token = Token{"^", T_DEREF, l.line}
		case '{':
			token = Token{"{", T_L_SQUIRLY, l.line}
		case '}':
			token = Token{"}", T_R_SQUIRLY, l.line}
		case '[':
			token = Token{"[", T_L_BLOCK, l.line}
		case ']':
			token = Token{"]", T_R_BLOCK, l.line}
		case '(':
			token = Token{"(", T_L_PAREN, l.line}
		case ')':
			token = Token{")", T_R_PAREN, l.line}
		case ',':
			token = Token{",", T_SEP, l.line}
		case '%':
			token = Token{"%", T_MOD, l.line}
		case '*':
			token = Token{"*", T_MUL, l.line}
		case '|':
			token = Token{"|", T_OR, l.line}

			if l.peekChar() == '|' {
				token = Token{"||", T_OROR, l.line}
				l.index++
			}
		case '&':
			token = Token{"&", T_AND, l.line}

			if l.peekChar() == '&' {
				token = Token{"&&", T_ANDAND, l.line}
				l.index++
			}
		case '=':
			token = Token{"=", T_ASSIGN, l.line}

			if l.peekChar() == '=' {
				token = Token{"==", T_EQ, l.line}
				l.index++
			}
		case '!':
			token = Token{"!", T_NOT, l.line}

			if l.peekChar() == '=' {
				token = Token{"!=", T_NEQ, l.line}
				l.index++
			}
		case '+':
			token = Token{"+", T_ADD, l.line}

			if l.peekChar() == '+' {
				token = Token{"++", T_INC, l.line}
				l.index++
			}
		case '-':
			token = Token{"-", T_SUB, l.line}

			if l.peekChar() == '-' {
				token = Token{"--", T_DINC, l.line}
				l.index++
			}
		case '<':
			token = Token{"<", T_LT, l.line}

			if l.peekChar() == '=' {
				token = Token{"<=", T_LTEQ, l.line}
				l.index++
			} else if l.peekChar() == '<' {
				token = Token{"<<", T_L_SHIFT, l.line}
				l.index++
			}
		case '>':
			token = Token{">", T_GT, l.line}

			if l.peekChar() == '=' {
				token = Token{">=", T_GTEQ, l.line}
				l.index++
			} else if l.peekChar() == '>' {
				token = Token{">>", T_R_SHIFT, l.line}
				l.index++
			}
		case '/':
			token = Token{"/", T_DIV, l.line}

			if l.peekChar() == '/' {
				for l.peekChar() != '\n' && l.peekChar() != 0 {
					if l.peekChar() == '\n' {
						l.line++
					}
					l.index++
				}
				continue
			}

			if l.peekChar() == '*' {
				for l.peekChar() != 0 {
					if l.peekChar() == '\n' {
						l.line++
					}

					if l.peekChar() == '/' && l.source[l.index] == '*' {
						l.index++
						break
					}
					l.index++
				}
				continue
			}
		case '\'':
			token = Token{"'", T_CHAR, l.line}

			if l.peekChar() == '\\' {
				l.index = l.index + 3

				if l.index >= len(l.source) {
					l.throwError("character", l.line, "more source", "end of source")
				}

				if l.source[l.index] != '\'' {
					l.throwError("character", l.line, "' (single quote)", string(l.source[l.index]))
				}
				token.data = token.data + string(l.source[l.index-2]) + string(l.source[l.index-1]) + string(l.source[l.index])
			} else {
				l.index = l.index + 2

				if l.index >= len(l.source) {
					l.throwError("character", l.line, "more source", "end of source")
				}

				if l.source[l.index] != '\'' {
					l.throwError("character", l.line, "' (single quote)", string(l.source[l.index]))
				}
				token.data = token.data + string(l.source[l.index-1]) + string(l.source[l.index])
			}
		case '"':
			token = Token{"\"", T_STRING, l.line}
			escaped := false
			for !(l.peekChar() == '"' && !escaped) && l.peekChar() != 0 {
				l.index++
				token.data = token.data + string(l.source[l.index])

				if l.source[l.index] == '\\' {
					escaped = true
				} else {
					escaped = false
				}
			}

			if l.peekChar() != '"' {
				l.throwError("string", l.line, "\" (end of string)", "end of source")
			}
			l.index++
			token.data = token.data + string(l.source[l.index])
		}
		if l.source[l.index]-'0' <= 9 {
			token = Token{string(l.source[l.index]), T_ILLEGAL, l.line}
			isFloat := false
			for (l.peekChar()-'0' <= 9 || l.peekChar() == '.') && l.peekChar() != 0 {
				l.index++
				token.data = token.data + string(l.source[l.index])

				if l.source[l.index] == '.' {
					isFloat = true
				}
			}

			if isFloat {
				if token.data[len(token.data)-1] == '.' {
					panic("Floats cannot end with '.'")
				}
				token.kind = T_FLOAT
			} else {
				token.kind = T_INT
			}
		} else if l.source[l.index]-'a' < 26 || l.source[l.index]-'A' < 26 {
			token = Token{string(l.source[l.index]), T_ILLEGAL, l.line}
			for (l.peekChar()-'a' < 26 || l.peekChar()-'A' < 26 || l.peekChar() == '_') && l.peekChar() != 0 {
				l.index++
				token.data = token.data + string(l.source[l.index])
			}
			switch token.data {
			case "for":
				token.kind = T_FOR
			case "range":
				token.kind = T_RANGE
			case "forever":
				token.kind = T_FOREVER
			case "while":
				token.kind = T_WHILE
			case "if":
				token.kind = T_IF
			case "elif":
				token.kind = T_ELIF
			case "else":
				token.kind = T_ELSE
			case "call":
				token.kind = T_CALL
			case "struct":
				token.kind = T_STRUCT
			case "fun":
				token.kind = T_FUN
			case "return":
				token.kind = T_RET
			case "break":
				token.kind = T_BREAK
			case "continue":
				token.kind = T_CONT
			case "enum":
				token.kind = T_ENUM
			case "nil":
				token.kind = T_NIL
			case "typedef":
				token.kind = T_TYPEDEF
			case "new":
				token.kind = T_NEW
			case "make":
				token.kind = T_MAKE
			case "map":
				token.kind = T_MAP
			case "switch":
				token.kind = T_SWITCH
			case "case":
				token.kind = T_CASE
			case "default":
				token.kind = T_DEFAULT
			case "const":
				token.kind = T_CONST
			case "true":
				token.kind = T_BOOL
			case "false":
				token.kind = T_BOOL
			case "byte":
				token.kind = T_TYPE
			case "word":
				token.kind = T_TYPE
			case "dword":
				token.kind = T_TYPE
			case "qword":
				token.kind = T_TYPE
			case "uint8":
				token.kind = T_TYPE
			case "uint16":
				token.kind = T_TYPE
			case "uint32":
				token.kind = T_TYPE
			case "uint64":
				token.kind = T_TYPE
			case "uint":
				token.kind = T_TYPE
			case "int8":
				token.kind = T_TYPE
			case "int16":
				token.kind = T_TYPE
			case "int32":
				token.kind = T_TYPE
			case "int64":
				token.kind = T_TYPE
			case "sint":
				token.kind = T_TYPE
			case "int":
				token.kind = T_TYPE
			case "char":
				token.kind = T_TYPE
			case "string":
				token.kind = T_TYPE
			case "float32":
				token.kind = T_TYPE
			case "float64":
				token.kind = T_TYPE
			case "double":
				token.kind = T_TYPE
			case "float":
				token.kind = T_TYPE
			case "bool":
				token.kind = T_TYPE
			case "any":
				token.kind = T_TYPE
			default:

				if len(tokens) == 0 {
					token.kind = T_IDENTIFIER
					break
				}
				prevToken := tokens[len(tokens)-1]

				if prevToken.kind == T_STRUCT || prevToken.kind == T_TYPEDEF || prevToken.kind == T_NEW || prevToken.kind == T_IDENTIFIER || prevToken.kind == T_MAP {
					token.kind = T_TYPE
				} else {
					token.kind = T_IDENTIFIER
				}
			}
		}

		if token.kind == T_ILLEGAL {
			l.throwError("lexing", l.line, "anything else", "ILLEGAL ("+string(l.source[l.index])+")")
		}
		tokens = append(tokens, token)
	}
	return tokens
}

func (l *Lexer) throwError(caller string, line int, expected string, got string) {
	panic("Error in the " + JOB_LEXER + "!\n" + "When the " + JOB_LEXER + " was trying to decipher: " + caller + "\n" + "Error found on line " + string(line) + "\n" + "Expected: " + expected + "\n" + "Got: " + got)
}

func (t TokenType) String() string {
	switch t {
	case T_ILLEGAL:
		return "ILLEGAL"
	case T_FOR:
		return "FOR"
	case T_RANGE:
		return "RANGE"
	case T_FOREVER:
		return "FOREVER"
	case T_WHILE:
		return "WHILE"
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
	case T_SWITCH:
		return "SWITCH"
	case T_CASE:
		return "CASE"
	case T_DEFAULT:
		return "DEFAULT"
	case T_CONST:
		return "CONST"
	case T_SEMICOLON:
		return "SEMICOLON"
	case T_ASSIGN:
		return "ASSIGN"
	case T_SEP:
		return "SEP"
	case T_COLON:
		return "COLON"
	case T_L_SHIFT:
		return "T_L_SHIFT"
	case T_R_SHIFT:
		return "T_R_SHIFT"
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
	case T_OROR:
		return "OROR"
	case T_ANDAND:
		return "ANDAND"
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

func (t Token) String() string {
	return "(" + t.data + " " + t.kind.String() + ")"
}

func (n *Node) String() string {
	output := n.kind.String() + ": " + n.data + "\n"
	var c *Node
	for i := 0; i < len(n.children); i++ {
		c = n.children[i]
		output = output + c.StringRec(1)
	}
	return output
}

func (n *Node) StringRec(indent int) string {
	output := getIndent(indent) + n.kind.String() + ": " + n.data + "\n"
	var c *Node
	for i := 0; i < len(n.children); i++ {
		c = n.children[i]
		output = output + c.StringRec(indent+1)
	}
	return output
}

func getIndent(indent int) string {
	output := ""
	for range indent {
		output = output + "\t"
	}
	return output
}

func (n NodeType) String() string {
	switch n {
	case N_ILLEGAL:
		return "ILLEGAL"
	case N_PROGRAM:
		return "PROGRAM"
	case N_VAR_DECLARATION:
		return "VAR_DECLARATION"
	case N_ELEMENT_ASSIGNMENT:
		return "ELEMENT_ASSIGNMENT"
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
	case N_SEMICOLON:
		return "SEMICOLON"
	case N_ASSIGN:
		return "ASSIGN"
	case N_SEP:
		return "SEP"
	case N_COLON:
		return "COLON"
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

const JOB_LEXER = "Lexer"

type float = float64

func main() {

}
