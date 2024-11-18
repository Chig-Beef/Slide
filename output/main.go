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

type Parser struct {
	source []Token
	tok    Token
	index  int
}

type Hoister struct {
	ast *Node
}

type VarType int

const (
	V_VAR VarType = iota
	V_TYPE
	V_FUNC
	V_MAP
)

type Var struct {
	kind     VarType
	data     string
	datatype string
	ref      *Node
	props    []*Var
	key      *Var
	value    *Var
	isArray  bool
}

type VarFrame struct {
	value *Var
	prev  *VarFrame
}

type VarStack struct {
	tail   *VarFrame
	length int
}

type GoEmitter struct {
	types   *Node
	consts  *Node
	funcs   *Node
	ast     *Node
	varType *Node
	inConst bool
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

func (p *Parser) nextToken() {
	var t Token

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
	return Token{}
}

func (p *Parser) parse() *Node {
	program := Node{N_PROGRAM}
	var n *Node
	p.nextToken()
	for p.tok.kind != T_ILLEGAL {
		switch p.tok.kind {
		case T_IDENTIFIER:
			n = p.variableDeclaration()
			program = children.append(program, n)
		case T_DINC:
			n = p.loneIncrement()
			program = children.append(program, n)
		case T_INC:
			n = p.loneIncrement()
			program = children.append(program, n)
		case T_IF:
			n = p.ifBlock()
			program = children.append(program, n)
		case T_FOREVER:
			n = p.foreverLoop()
			program = children.append(program, n)
		case T_RANGE:
			n = p.rangeLoop()
			program = children.append(program, n)
		case T_FOR:
			n = p.forLoop()
			program = children.append(program, n)
		case T_WHILE:
			n = p.whileLoop()
			program = children.append(program, n)
		case T_CALL:
			n = p.loneCall()
			program = children.append(program, n)
		case T_STRUCT:
			n = p.structDef()
			program = children.append(program, n)
		case T_FUN:
			n = p.funcDef()
			program = children.append(program, n)
		case T_RET:
			n = p.retStatement()
			program = children.append(program, n)
		case T_BREAK:
			n = p.breakStatement()
			program = children.append(program, n)
		case T_CONT:
			n = p.contStatement()
			program = children.append(program, n)
		case T_ENUM:
			n = p.enumDef()
			program = children.append(program, n)
		case T_TYPEDEF:
			n = p.typeDef()
			program = children.append(program, n)
		case T_L_BLOCK:
			n = p.elementAssignment()
			program = children.append(program, n)
		case T_SWITCH:
			n = p.switchStatement()
			program = children.append(program, n)
		case T_CONST:
			n = p.constantStatement()
			program = children.append(program, n)
		default:
			panic("Bad start to statement: " + p.tok.kind.String() + " on line " + string(p.tok.line))
		}
		p.nextToken()
	}
	return &program
}

func (p *Parser) loneIncrement() *Node {
	FUNC_NAME := "lone increment"
	n := Node{N_LONE_INC, p.tok.line}

	if p.tok.kind == T_INC {
		n = children.append(n, &Node{N_INC, p.tok.data, p.tok.line})
	} else if p.tok.kind == T_DINC {
		n = children.append(n, &Node{N_DINC, p.tok.data, p.tok.line})
	} else {
		p.throwError(FUNC_NAME, p.tok.line, "inc or dinc", p.tok)
	}
	p.nextToken()
	if p.tok.kind != T_IDENTIFIER {
		p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok)
	}

	if p.peekToken().kind == T_ACCESS {
		n = children.append(n, p.property())
		p.nextToken()
	} else {
		n = children.append(n, &Node{N_IDENTIFIER, p.tok.data, p.tok.line})
		p.nextToken()
	}

	if p.tok.kind != T_SEMICOLON {
		p.throwError(FUNC_NAME, p.tok.line, "semicolon", p.tok)
	}
	n = children.append(n, &Node{N_SEMICOLON, p.tok.data, p.tok.line})
	return &n
}

func (p *Parser) constantStatement() *Node {
	FUNC_NAME := "constant statement"
	n := Node{N_CONSTANT, p.tok.line}

	if p.tok.kind != T_CONST {
		p.throwError(FUNC_NAME, p.tok.line, "const", p.tok)
	}
	n = children.append(n, &Node{N_CONST, p.tok.data, p.tok.line})
	p.nextToken()
	n = children.append(n, p.variableDeclaration())
	return &n
}

func (p *Parser) variableDeclaration() *Node {
	FUNC_NAME := "variable declaration"
	n := Node{N_VAR_DECLARATION, p.tok.line}
	n = children.append(n, p.assignment())
	p.nextToken()
	if p.tok.kind != T_SEMICOLON {
		p.throwError(FUNC_NAME, p.tok.line, "semicolon", p.tok)
	}
	n = children.append(n, &Node{N_SEMICOLON, p.tok.data, p.tok.line})
	return &n
}

func (p *Parser) switchStatement() *Node {
	FUNC_NAME := "switch statement"
	n := Node{N_SWITCH_STATE, p.tok.line}

	if p.tok.kind != T_SWITCH {
		p.throwError(FUNC_NAME, p.tok.line, "switch", p.tok)
	}
	n = children.append(n, &Node{N_SWITCH, p.tok.data, p.tok.line})
	p.nextToken()
	n = children.append(n, p.expression())
	p.nextToken()
	if p.tok.kind != T_L_SQUIRLY {
		p.throwError(FUNC_NAME, p.tok.line, "left squirly", p.tok)
	}
	n = children.append(n, &Node{N_L_SQUIRLY, p.tok.data, p.tok.line})
	p.nextToken()
	for p.tok.kind == T_CASE {
		n = children.append(n, p.caseStatement())
		p.nextToken()
	}

	if p.tok.kind == T_DEFAULT {
		n = children.append(n, p.defaultStatement())
		p.nextToken()
	}

	if p.tok.kind != T_R_SQUIRLY {
		p.throwError(FUNC_NAME, p.tok.line, "right squirly", p.tok)
	}
	n = children.append(n, &Node{N_R_SQUIRLY, p.tok.data, p.tok.line})
	return &n
}

func (p *Parser) caseStatement() *Node {
	FUNC_NAME := "case statement"
	n := Node{N_CASE_STATE, p.tok.line}

	if p.tok.kind != T_CASE {
		p.throwError(FUNC_NAME, p.tok.line, "case", p.tok)
	}
	n = children.append(n, &Node{N_CASE, p.tok.data, p.tok.line})
	p.nextToken()
	n = children.append(n, p.expression())
	p.nextToken()
	for p.tok.kind == T_SEP {
		n = children.append(n, &Node{N_SEP, p.tok.data, p.tok.line})
		p.nextToken()
		n = children.append(n, p.expression())
		p.nextToken()
	}

	if p.tok.kind != T_COLON {
		p.throwError(FUNC_NAME, p.tok.line, "colon", p.tok)
	}
	n = children.append(n, &Node{N_COLON, p.tok.data, p.tok.line})
	p.nextToken()
	n = children.append(n, p.caseBlock())
	return &n
}

func (p *Parser) defaultStatement() *Node {
	FUNC_NAME := "default statement"
	n := Node{N_DEFAULT_STATE, p.tok.line}

	if p.tok.kind != T_DEFAULT {
		p.throwError(FUNC_NAME, p.tok.line, "default", p.tok)
	}
	n = children.append(n, &Node{N_DEFAULT, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_COLON {
		p.throwError(FUNC_NAME, p.tok.line, "colon", p.tok)
	}
	n = children.append(n, &Node{N_COLON, p.tok.data, p.tok.line})
	p.nextToken()
	n = children.append(n, p.caseBlock())
	return &n
}

func (p *Parser) caseBlock() *Node {
	block := Node{N_CASE_BLOCK, p.tok.line}
	var n *Node
	for p.tok.kind != T_ILLEGAL && p.tok.kind != T_CASE && p.tok.kind != T_DEFAULT && p.tok.kind != T_R_SQUIRLY {
		switch p.tok.kind {
		case T_IDENTIFIER:
			n = p.variableDeclaration()
			block = children.append(block, n)
		case T_INC:
			n = p.loneIncrement()
			block = children.append(block, n)
		case T_DINC:
			n = p.loneIncrement()
			block = children.append(block, n)
		case T_IF:
			n = p.ifBlock()
			block = children.append(block, n)
		case T_FOREVER:
			n = p.foreverLoop()
			block = children.append(block, n)
		case T_RANGE:
			n = p.rangeLoop()
			block = children.append(block, n)
		case T_FOR:
			n = p.forLoop()
			block = children.append(block, n)
		case T_WHILE:
			n = p.whileLoop()
			block = children.append(block, n)
		case T_CALL:
			n = p.loneCall()
			block = children.append(block, n)
		case T_STRUCT:
			n = p.structDef()
			block = children.append(block, n)
		case T_FUN:
			n = p.funcDef()
			block = children.append(block, n)
		case T_RET:
			n = p.retStatement()
			block = children.append(block, n)
		case T_BREAK:
			n = p.breakStatement()
			block = children.append(block, n)
		case T_CONT:
			n = p.contStatement()
			block = children.append(block, n)
		case T_ENUM:
			n = p.enumDef()
			block = children.append(block, n)
		case T_TYPEDEF:
			n = p.typeDef()
			block = children.append(block, n)
		case T_L_BLOCK:
			n = p.elementAssignment()
			block = children.append(block, n)
		case T_SWITCH:
			n = p.switchStatement()
			block = children.append(block, n)
		default:
			panic("Bad start to statement: " + p.tok.kind.String() + " on line " + string(p.tok.line))
		}
		p.nextToken()
	}
	p.index--
	return &block
}

func (p *Parser) elementAssignment() *Node {
	FUNC_NAME := "element assignment"
	n := Node{N_ELEMENT_ASSIGNMENT, p.tok.line}
	n = children.append(n, p.indexUnary())
	p.nextToken()
	n = children.append(n, p.assignment())
	p.nextToken()
	if p.tok.kind != T_SEMICOLON {
		p.throwError(FUNC_NAME, p.tok.line, "semicolon", p.tok)
	}
	n = children.append(n, &Node{N_SEMICOLON, p.tok.data, p.tok.line})
	return &n
}

func (p *Parser) ifBlock() *Node {
	FUNC_NAME := "if block"
	n := Node{N_IF_BLOCK, p.tok.line}

	if p.tok.kind != T_IF {
		p.throwError(FUNC_NAME, p.tok.line, "if", p.tok)
	}
	n = children.append(n, &Node{N_IF, p.tok.data, p.tok.line})
	p.nextToken()
	n = children.append(n, p.condition())
	p.nextToken()
	n = children.append(n, p.block())
	p.nextToken()
	for p.tok.kind == T_ELIF {
		n = children.append(n, &Node{N_ELIF, p.tok.data, p.tok.line})
		p.nextToken()
		n = children.append(n, p.condition())
		p.nextToken()
		n = children.append(n, p.block())
		p.nextToken()
	}

	if p.tok.kind == T_ELSE {
		n = children.append(n, &Node{N_ELSE, p.tok.data, p.tok.line})
		p.nextToken()
		n = children.append(n, p.block())
	} else {
		p.index--
	}
	return &n
}

func (p *Parser) foreverLoop() *Node {
	FUNC_NAME := "forever loop"
	n := Node{N_FOREVER_LOOP, p.tok.line}

	if p.tok.kind != T_FOREVER {
		p.throwError(FUNC_NAME, p.tok.line, "forever", p.tok)
	}
	n = children.append(n, &Node{N_FOREVER, p.tok.line})
	p.nextToken()
	n = children.append(n, p.block())
	return &n
}

func (p *Parser) rangeLoop() *Node {
	FUNC_NAME := "range loop"
	n := Node{N_RANGE_LOOP, p.tok.line}

	if p.tok.kind != T_RANGE {
		p.throwError(FUNC_NAME, p.tok.line, "range", p.tok)
	}
	n = children.append(n, &Node{N_RANGE, p.tok.line})
	p.nextToken()
	n = children.append(n, p.expression())
	p.nextToken()
	n = children.append(n, p.block())
	return &n
}

func (p *Parser) forLoop() *Node {
	FUNC_NAME := "for loop"
	n := Node{N_FOR_LOOP, p.tok.line}

	if p.tok.kind != T_FOR {
		p.throwError(FUNC_NAME, p.tok.line, "for", p.tok)
	}
	n = children.append(n, &Node{N_FOR, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_SEMICOLON {
		n = children.append(n, p.assignment())
		p.nextToken()
	}

	if p.tok.kind != T_SEMICOLON {
		p.throwError(FUNC_NAME, p.tok.line, "semicolon", p.tok)
	}
	n = children.append(n, &Node{N_SEMICOLON, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_SEMICOLON {
		n = children.append(n, p.expression())
		p.nextToken()
	}

	if p.tok.kind != T_SEMICOLON {
		p.throwError(FUNC_NAME, p.tok.line, "semicolon", p.tok)
	}
	n = children.append(n, &Node{N_SEMICOLON, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_L_SQUIRLY {
		if p.tok.kind == T_INC || p.tok.kind == T_DINC {
			inc := Node{N_UNARY_OPERATION, p.tok.line}

			if p.tok.kind == T_INC {
				inc = children.append(inc, &Node{N_INC, p.tok.data, p.tok.line})
			} else if p.tok.kind == T_DINC {
				inc = children.append(inc, &Node{N_DINC, p.tok.data, p.tok.line})
			} else {
				p.throwError(FUNC_NAME, p.tok.line, "inc or dinc", p.tok)
			}
			p.nextToken()
			if p.peekToken().kind == T_ACCESS {
				inc = children.append(inc, p.property())
			} else {
				inc = children.append(inc, &Node{N_IDENTIFIER, p.tok.data, p.tok.line})
			}
			inc = children.append(inc, &Node{N_IDENTIFIER, p.tok.data, p.tok.line})
			n = children.append(n, &inc)
		} else {
			n = children.append(n, p.assignment())
		}
		p.nextToken()
	}
	n = children.append(n, p.block())
	return &n
}

func (p *Parser) whileLoop() *Node {
	FUNC_NAME := "while loop"
	n := Node{N_WHILE_LOOP, p.tok.line}

	if p.tok.kind != T_WHILE {
		p.throwError(FUNC_NAME, p.tok.line, "while", p.tok)
	}
	n = children.append(n, &Node{N_WHILE, p.tok.line})
	p.nextToken()
	n = children.append(n, p.condition())
	p.nextToken()
	n = children.append(n, p.block())
	return &n
}

func (p *Parser) structDef() *Node {
	FUNC_NAME := "struct definition"
	n := Node{N_STRUCT_DEF, p.tok.line}

	if p.tok.kind != T_STRUCT {
		p.throwError(FUNC_NAME, p.tok.line, "struct", p.tok)
	}
	n = children.append(n, &Node{N_STRUCT, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_TYPE {
		p.throwError(FUNC_NAME, p.tok.line, "type", p.tok)
	}
	n = children.append(n, &Node{N_TYPE, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_L_SQUIRLY {
		p.throwError(FUNC_NAME, p.tok.line, "left squirly", p.tok)
	}
	n = children.append(n, &Node{N_L_SQUIRLY, p.tok.data, p.tok.line})
	p.nextToken()
	for p.tok.kind != T_R_SQUIRLY && p.tok.kind != T_ILLEGAL {
		if p.tok.kind != T_IDENTIFIER {
			p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok)
		}
		n = children.append(n, &Node{N_IDENTIFIER, p.tok.data, p.tok.line})
		p.nextToken()
		n = children.append(n, p.complexType())
		p.nextToken()
	}

	if p.tok.kind != T_R_SQUIRLY {
		p.throwError(FUNC_NAME, p.tok.line, "right squirly", p.tok)
	}
	n = children.append(n, &Node{N_R_SQUIRLY, p.tok.data, p.tok.line})
	return &n
}

func (p *Parser) funcDef() *Node {
	FUNC_NAME := "function definition"
	n := Node{N_FUNC_DEF, p.tok.line}

	if p.tok.kind != T_FUN {
		p.throwError(FUNC_NAME, p.tok.line, "fun", p.tok)
	}
	n = children.append(n, &Node{N_FUN, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind == T_L_PAREN {
		n = children.append(n, p.methodReceiver())
		p.nextToken()
	}

	if p.tok.kind != T_IDENTIFIER {
		p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok)
	}
	n = children.append(n, &Node{N_IDENTIFIER, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_L_PAREN {
		p.throwError(FUNC_NAME, p.tok.line, "left paren", p.tok)
	}
	n = children.append(n, &Node{N_L_PAREN, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind == T_IDENTIFIER {
		n = children.append(n, &Node{N_IDENTIFIER, p.tok.data, p.tok.line})
		p.nextToken()
		n = children.append(n, p.complexType())
		p.nextToken()
		for p.tok.kind == T_SEP {
			n = children.append(n, &Node{N_SEP, p.tok.data, p.tok.line})
			p.nextToken()
			if p.tok.kind != T_IDENTIFIER {
				p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok)
			}
			n = children.append(n, &Node{N_IDENTIFIER, p.tok.data, p.tok.line})
			p.nextToken()
			n = children.append(n, p.complexType())
			p.nextToken()
		}
	}

	if p.tok.kind != T_R_PAREN {
		p.throwError(FUNC_NAME, p.tok.line, "right paren", p.tok)
	}
	n = children.append(n, &Node{N_R_PAREN, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_L_SQUIRLY {
		n = children.append(n, p.complexType())
		p.nextToken()
	}
	n = children.append(n, p.block())
	return &n
}

func (p *Parser) methodReceiver() *Node {
	FUNC_NAME := "method receiver"
	n := Node{N_METHOD_RECEIVER, p.tok.line}

	if p.tok.kind != T_L_PAREN {
		p.throwError(FUNC_NAME, p.tok.line, "left paren", p.tok)
	}
	n = children.append(n, &Node{N_L_PAREN, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_IDENTIFIER {
		p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok)
	}
	n = children.append(n, &Node{N_IDENTIFIER, p.tok.data, p.tok.line})
	p.nextToken()
	n = children.append(n, p.complexType())
	p.nextToken()
	if p.tok.kind != T_R_PAREN {
		p.throwError(FUNC_NAME, p.tok.line, "right paren", p.tok)
	}
	n = children.append(n, &Node{N_R_PAREN, p.tok.data, p.tok.line})
	return &n
}

func (p *Parser) retStatement() *Node {
	FUNC_NAME := "return statement"
	n := Node{N_RET_STATE, p.tok.line}

	if p.tok.kind != T_RET {
		p.throwError(FUNC_NAME, p.tok.line, "return", p.tok)
	}
	n = children.append(n, &Node{N_RET, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind == T_SEMICOLON {
		n = children.append(n, &Node{N_SEMICOLON, p.tok.data, p.tok.line})
		return &n
	}
	n = children.append(n, p.expression())
	p.nextToken()
	if p.tok.kind != T_SEMICOLON {
		p.throwError(FUNC_NAME, p.tok.line, "semicolon", p.tok)
	}
	n = children.append(n, &Node{N_SEMICOLON, p.tok.data, p.tok.line})
	return &n
}

func (p *Parser) breakStatement() *Node {
	FUNC_NAME := "break statement"
	n := Node{N_BREAK_STATE, p.tok.line}

	if p.tok.kind != T_BREAK {
		p.throwError(FUNC_NAME, p.tok.line, "break", p.tok)
	}
	n = children.append(n, &Node{N_BREAK, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind == T_SEMICOLON {
		n = children.append(n, &Node{N_SEMICOLON, p.tok.data, p.tok.line})
		return &n
	}
	n = children.append(n, p.expression())
	p.nextToken()
	if p.tok.kind != T_SEMICOLON {
		p.throwError(FUNC_NAME, p.tok.line, "semicolon", p.tok)
	}
	n = children.append(n, &Node{N_SEMICOLON, p.tok.data, p.tok.line})
	return &n
}

func (p *Parser) contStatement() *Node {
	FUNC_NAME := "continue statement"
	n := Node{N_CONT_STATE, p.tok.line}

	if p.tok.kind != T_CONT {
		p.throwError(FUNC_NAME, p.tok.line, "continue", p.tok)
	}
	n = children.append(n, &Node{N_CONT, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind == T_SEMICOLON {
		n = children.append(n, &Node{N_SEMICOLON, p.tok.data, p.tok.line})
		return &n
	}
	n = children.append(n, p.expression())
	p.nextToken()
	if p.tok.kind != T_SEMICOLON {
		p.throwError(FUNC_NAME, p.tok.line, "semicolon", p.tok)
	}
	n = children.append(n, &Node{N_SEMICOLON, p.tok.data, p.tok.line})
	return &n
}

func (p *Parser) enumDef() *Node {
	FUNC_NAME := "enum definition"
	n := Node{N_ENUM_DEF, p.tok.line}

	if p.tok.kind != T_ENUM {
		p.throwError(FUNC_NAME, p.tok.line, "enum", p.tok)
	}
	n = children.append(n, &Node{N_ENUM, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_IDENTIFIER {
		p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok)
	}
	n = children.append(n, &Node{N_IDENTIFIER, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_L_SQUIRLY {
		p.throwError(FUNC_NAME, p.tok.line, "left squirly", p.tok)
	}
	n = children.append(n, &Node{N_L_SQUIRLY, p.tok.data, p.tok.line})
	p.nextToken()
	for p.tok.kind != T_R_SQUIRLY && p.tok.kind != T_ILLEGAL {
		if p.tok.kind != T_IDENTIFIER {
			p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok)
		}
		n = children.append(n, &Node{N_IDENTIFIER, p.tok.data, p.tok.line})
		p.nextToken()
		if p.tok.kind != T_SEP {
			p.throwError(FUNC_NAME, p.tok.line, "separator", p.tok)
		}
		n = children.append(n, &Node{N_SEP, p.tok.data, p.tok.line})
		p.nextToken()
	}

	if p.tok.kind != T_R_SQUIRLY {
		p.throwError(FUNC_NAME, p.tok.line, "right squirly", p.tok)
	}
	n = children.append(n, &Node{N_R_SQUIRLY, p.tok.data, p.tok.line})
	return &n
}

func (p *Parser) condition() *Node {
	return p.expression()
}

func (p *Parser) expression() *Node {
	var v *Node
	n := Node{N_EXPRESSION, p.tok.line}
	n = children.append(n, p.value())
	p.nextToken()
	v = p.operator()
	for v != nil {
		n = children.append(n, v)
		p.nextToken()
		n = children.append(n, p.value())
		p.nextToken()
		v = p.operator()
	}
	p.index--
	return &n
}

func (p *Parser) assignment() *Node {
	FUNC_NAME := "assignment"
	n := Node{N_ASSIGNMENT, p.tok.line}

	if p.tok.kind != T_IDENTIFIER {
		p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok)
	}

	if p.peekToken().kind == T_ACCESS {
		n = children.append(n, p.property())
		p.nextToken()
	} else {
		n = children.append(n, &Node{N_IDENTIFIER, p.tok.data, p.tok.line})
		p.nextToken()
	}
	isMap := p.tok.kind == T_MAP

	if p.tok.kind != T_ASSIGN {
		n = children.append(n, p.complexType())

		if isMap {
			return &n
		}
		p.nextToken()
	}

	if p.tok.kind != T_ASSIGN {
		if p.tok.kind != T_SEMICOLON {
			p.throwError(FUNC_NAME, p.tok.line, "assign or semicolon", p.tok)
		}
		p.index--
		return &n
	}
	n = children.append(n, &Node{N_ASSIGN, p.tok.data, p.tok.line})
	p.nextToken()
	n = children.append(n, p.expression())
	return &n
}

func (p *Parser) loneCall() *Node {
	FUNC_NAME := "lone call"
	n := Node{N_LONE_CALL, p.tok.line}
	n = children.append(n, p.funcCall())
	p.nextToken()
	if p.tok.kind != T_SEMICOLON {
		p.throwError(FUNC_NAME, p.tok.line, "semicolon", p.tok)
	}
	n = children.append(n, &Node{N_SEMICOLON, p.tok.data, p.tok.line})
	return &n
}

func (p *Parser) funcCall() *Node {
	FUNC_NAME := "function call"
	n := Node{N_FUNC_NAME, p.tok.line}

	if p.tok.kind != T_CALL {
		p.throwError(FUNC_NAME, p.tok.line, "call", p.tok)
	}
	n = children.append(n, &Node{N_CALL, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_IDENTIFIER {
		p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok)
	}

	if p.peekToken().kind == T_ACCESS {
		n = children.append(n, p.property())
		p.nextToken()
	} else {
		n = children.append(n, &Node{N_IDENTIFIER, p.tok.data, p.tok.line})
		p.nextToken()
	}

	if p.tok.kind != T_L_PAREN {
		p.throwError(FUNC_NAME, p.tok.line, "left paren", p.tok)
	}
	n = children.append(n, &Node{N_L_PAREN, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_R_PAREN {
		n = children.append(n, p.expression())
		p.nextToken()
		for p.tok.kind == T_SEP {
			n = children.append(n, &Node{N_SEP, p.tok.data, p.tok.line})
			p.nextToken()
			n = children.append(n, p.expression())
			p.nextToken()
		}
	}

	if p.tok.kind != T_R_PAREN {
		p.throwError(FUNC_NAME, p.tok.line, "right paren", p.tok)
	}
	n = children.append(n, &Node{N_R_PAREN, p.tok.data, p.tok.line})
	return &n
}

func (p *Parser) structNew() *Node {
	FUNC_NAME := "new struct"
	n := Node{N_STRUCT_NEW, p.tok.line}

	if p.tok.kind != T_NEW {
		p.throwError(FUNC_NAME, p.tok.line, "new", p.tok)
	}
	n = children.append(n, &Node{N_NEW, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_TYPE {
		p.throwError(FUNC_NAME, p.tok.line, "type", p.tok)
	}
	n = children.append(n, &Node{N_TYPE, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_L_PAREN {
		p.throwError(FUNC_NAME, p.tok.line, "left paren", p.tok)
	}
	n = children.append(n, &Node{N_L_PAREN, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_R_PAREN {
		n = children.append(n, p.expression())
		p.nextToken()
		for p.tok.kind == T_SEP {
			n = children.append(n, &Node{N_SEP, p.tok.data, p.tok.line})
			p.nextToken()
			n = children.append(n, p.expression())
			p.nextToken()
		}
	}

	if p.tok.kind != T_R_PAREN {
		p.throwError(FUNC_NAME, p.tok.line, "right paren", p.tok)
	}
	n = children.append(n, &Node{N_R_PAREN, p.tok.data, p.tok.line})
	return &n
}

func (p *Parser) block() *Node {
	FUNC_NAME := "block"
	block := Node{N_BLOCK, p.tok.line}
	var n *Node

	if p.tok.kind != T_L_SQUIRLY {
		p.throwError(FUNC_NAME, p.tok.line, "left squirly", p.tok)
	}
	block = children.append(block, &Node{N_L_SQUIRLY, p.tok.data, p.tok.line})
	p.nextToken()
	for p.tok.kind != T_ILLEGAL && p.tok.kind != T_R_SQUIRLY {
		switch p.tok.kind {
		case T_IDENTIFIER:
			n = p.variableDeclaration()
			block = children.append(block, n)
		case T_INC:
			n = p.loneIncrement()
			block = children.append(block, n)
		case T_DINC:
			n = p.loneIncrement()
			block = children.append(block, n)
		case T_IF:
			n = p.ifBlock()
			block = children.append(block, n)
		case T_FOREVER:
			n = p.foreverLoop()
			block = children.append(block, n)
		case T_RANGE:
			n = p.rangeLoop()
			block = children.append(block, n)
		case T_FOR:
			n = p.forLoop()
			block = children.append(block, n)
		case T_WHILE:
			n = p.whileLoop()
			block = children.append(block, n)
		case T_CALL:
			n = p.loneCall()
			block = children.append(block, n)
		case T_STRUCT:
			n = p.structDef()
			block = children.append(block, n)
		case T_FUN:
			n = p.funcDef()
			block = children.append(block, n)
		case T_RET:
			n = p.retStatement()
			block = children.append(block, n)
		case T_BREAK:
			n = p.breakStatement()
			block = children.append(block, n)
		case T_CONT:
			n = p.contStatement()
			block = children.append(block, n)
		case T_ENUM:
			n = p.enumDef()
			block = children.append(block, n)
		case T_TYPEDEF:
			n = p.typeDef()
			block = children.append(block, n)
		case T_L_BLOCK:
			n = p.elementAssignment()
			block = children.append(block, n)
		case T_SWITCH:
			n = p.switchStatement()
			block = children.append(block, n)
		default:
			panic("Bad start to statement: " + p.tok.kind.String() + " on line " + string(p.tok.line))
		}
		p.nextToken()
	}

	if p.tok.kind != T_R_SQUIRLY {
		p.throwError(FUNC_NAME, p.tok.line, "right squirly", p.tok)
	}
	block = children.append(block, &Node{N_R_SQUIRLY, p.tok.data, p.tok.line})
	return &block
}

func (p *Parser) typeDef() *Node {
	FUNC_NAME := "type definition"
	n := Node{N_NEW_TYPE, p.tok.line}

	if p.tok.kind != T_TYPEDEF {
		p.throwError(FUNC_NAME, p.tok.line, "typedef", p.tok)
	}
	n = children.append(n, &Node{N_TYPEDEF, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_TYPE {
		p.throwError(FUNC_NAME, p.tok.line, "type", p.tok)
	}
	n = children.append(n, &Node{N_TYPE, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_TYPE {
		p.throwError(FUNC_NAME, p.tok.line, "type", p.tok)
	}
	n = children.append(n, &Node{N_TYPE, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_SEMICOLON {
		p.throwError(FUNC_NAME, p.tok.line, "semicolon", p.tok)
	}
	n = children.append(n, &Node{N_SEMICOLON, p.tok.line})
	return &n
}

func (p *Parser) value() *Node {
	FUNC_NAME := "value"
	var n Node
	switch p.tok.kind {
	case T_IDENTIFIER:

		if p.peekToken().kind == T_ACCESS {
			return p.property()
		} else {
			return &Node{N_IDENTIFIER, p.tok.data, p.tok.line}
		}
	case T_INT:
		return &Node{N_INT, p.tok.data, p.tok.line}
	case T_FLOAT:
		return &Node{N_FLOAT, p.tok.data, p.tok.line}
	case T_STRING:
		return &Node{N_STRING, p.tok.data, p.tok.line}
	case T_CHAR:
		return &Node{N_CHAR, p.tok.data, p.tok.line}
	case T_BOOL:
		return &Node{N_BOOL, p.tok.data, p.tok.line}
	case T_NIL:
		return &Node{N_NIL, p.tok.data, p.tok.line}
	case T_NOT:
		n = Node{N_UNARY_OPERATION, p.tok.data, p.tok.line}
		n = children.append(n, &Node{N_NOT, p.tok.data, p.tok.line})
		p.nextToken()
		n = children.append(n, p.value())
		return &n
	case T_INC:
		n = Node{N_UNARY_OPERATION, p.tok.data, p.tok.line}
		n = children.append(n, &Node{N_INC, p.tok.data, p.tok.line})
		p.nextToken()
		n = children.append(n, p.value())
		return &n
	case T_DINC:
		n = Node{N_UNARY_OPERATION, p.tok.data, p.tok.line}
		n = children.append(n, &Node{N_DINC, p.tok.data, p.tok.line})
		p.nextToken()
		n = children.append(n, p.value())
		return &n
	case T_REF:
		n = Node{N_UNARY_OPERATION, p.tok.data, p.tok.line}
		n = children.append(n, &Node{N_REF, p.tok.data, p.tok.line})
		p.nextToken()
		n = children.append(n, p.value())
		return &n
	case T_DEREF:
		n = Node{N_UNARY_OPERATION, p.tok.data, p.tok.line}
		n = children.append(n, &Node{N_DEREF, p.tok.data, p.tok.line})
		p.nextToken()
		n = children.append(n, p.value())
		return &n
	case T_L_BLOCK:
		n = Node{N_UNARY_OPERATION, p.tok.data, p.tok.line}
		n = children.append(n, p.indexUnary())
		p.nextToken()
		n = children.append(n, p.value())
		return &n
	case T_MAKE:
		return p.makeArray()
	case T_L_PAREN:
		n = Node{N_BRACKETED_VALUE, p.tok.line}
		n = children.append(n, &Node{N_L_PAREN, p.tok.data, p.tok.line})
		p.nextToken()
		n = children.append(n, p.expression())
		p.nextToken()
		if p.tok.kind != T_R_PAREN {
			p.throwError(FUNC_NAME, p.tok.line, "right paren", p.tok)
		}
		n = children.append(n, &Node{N_R_PAREN, p.tok.data, p.tok.line})
		return &n
	case T_CALL:
		return p.funcCall()
	case T_NEW:
		return p.structNew()
	default:
		p.throwError(FUNC_NAME, p.tok.line, "unary or value", p.tok)
		return nil
	}
}

func (p *Parser) operator() *Node {
	switch p.tok.kind {
	case T_ADD:
		return &Node{N_ADD, p.tok.data, p.tok.line}
	case T_XOR:
		return &Node{N_XOR, p.tok.data, p.tok.line}
	case T_ACCESS:
		return &Node{N_ACCESS, p.tok.data, p.tok.line}
	case T_NEQ:
		return &Node{N_NEQ, p.tok.data, p.tok.line}
	case T_MOD:
		return &Node{N_MOD, p.tok.data, p.tok.line}
	case T_EQ:
		return &Node{N_EQ, p.tok.data, p.tok.line}
	case T_LT:
		return &Node{N_LT, p.tok.data, p.tok.line}
	case T_GT:
		return &Node{N_GT, p.tok.data, p.tok.line}
	case T_LTEQ:
		return &Node{N_LTEQ, p.tok.data, p.tok.line}
	case T_GTEQ:
		return &Node{N_GTEQ, p.tok.data, p.tok.line}
	case T_SUB:
		return &Node{N_SUB, p.tok.data, p.tok.line}
	case T_MUL:
		return &Node{N_MUL, p.tok.data, p.tok.line}
	case T_DIV:
		return &Node{N_DIV, p.tok.data, p.tok.line}
	case T_OR:
		return &Node{N_OR, p.tok.data, p.tok.line}
	case T_AND:
		return &Node{N_AND, p.tok.data, p.tok.line}
	case T_OROR:
		return &Node{N_OROR, p.tok.data, p.tok.line}
	case T_ANDAND:
		return &Node{N_ANDAND, p.tok.data, p.tok.line}
	case T_L_SHIFT:
		return &Node{N_L_SHIFT, p.tok.data, p.tok.line}
	case T_R_SHIFT:
		return &Node{N_R_SHIFT, p.tok.data, p.tok.line}
	default:
		return nil
	}
}

func (p *Parser) indexUnary() *Node {
	FUNC_NAME := "index unary"
	n := Node{N_INDEX, p.tok.line}

	if p.tok.kind != T_L_BLOCK {
		p.throwError(FUNC_NAME, p.tok.line, "left block", p.tok)
	}
	n = children.append(n, &Node{N_L_BLOCK, p.tok.data, p.tok.line})
	p.nextToken()
	n = children.append(n, p.expression())
	p.nextToken()
	if p.tok.kind != T_R_BLOCK {
		p.throwError(FUNC_NAME, p.tok.line, "right block", p.tok)
	}
	n = children.append(n, &Node{N_R_BLOCK, p.tok.data, p.tok.line})
	return &n
}

func (p *Parser) makeArray() *Node {
	FUNC_NAME := "make array"
	n := Node{N_MAKE_ARRAY, p.tok.line}

	if p.tok.kind != T_MAKE {
		p.throwError(FUNC_NAME, p.tok.line, "make", p.tok)
	}
	n = children.append(n, &Node{N_MAKE, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_L_BLOCK {
		p.throwError(FUNC_NAME, p.tok.line, "left block", p.tok)
	}
	n = children.append(n, &Node{N_L_BLOCK, p.tok.data, p.tok.line})
	p.nextToken()
	for p.tok.kind != T_R_BLOCK {
		n = children.append(n, p.expression())
		p.nextToken()
		if p.tok.kind != T_SEP {
			break
		}
		n = children.append(n, &Node{N_SEP, p.tok.data, p.tok.line})
		p.nextToken()
	}

	if p.tok.kind != T_R_BLOCK {
		p.throwError(FUNC_NAME, p.tok.line, "right block", p.tok)
	}
	n = children.append(n, &Node{N_R_BLOCK, p.tok.data, p.tok.line})
	return &n
}

func (p *Parser) complexType() *Node {
	FUNC_NAME := "complex type"
	n := Node{N_COMPLEX_TYPE, p.tok.line}

	if p.tok.kind == T_MAP {
		n = children.append(n, &Node{N_MAP, p.tok.data, p.tok.line})
		p.nextToken()
		if p.tok.kind != T_TYPE {
			p.throwError(FUNC_NAME, p.tok.line, "type", p.tok)
		}
		n = children.append(n, &Node{N_TYPE, p.tok.data, p.tok.line})
		p.nextToken()
		if p.tok.kind != T_L_BLOCK {
			p.throwError(FUNC_NAME, p.tok.line, "left block", p.tok)
		}
		n = children.append(n, &Node{N_L_BLOCK, p.tok.data, p.tok.line})
		p.nextToken()
		if p.tok.kind != T_IDENTIFIER && p.tok.kind != T_TYPE {
			p.throwError(FUNC_NAME, p.tok.line, "identifier (but really a type)", p.tok)
		}
		n = children.append(n, &Node{N_TYPE, p.tok.data, p.tok.line})
		p.nextToken()
		if p.tok.kind != T_R_BLOCK {
			p.throwError(FUNC_NAME, p.tok.line, "right block", p.tok)
		}
		n = children.append(n, &Node{N_R_BLOCK, p.tok.data, p.tok.line})
		return &n
	} else if p.tok.kind == T_TYPE || p.tok.kind == T_IDENTIFIER {
		n = children.append(n, &Node{N_TYPE, p.tok.data, p.tok.line})
		pt := p.peekToken()

		if pt.kind == T_DEREF {
			p.nextToken()
			n = children.append(n, &Node{N_DEREF, p.tok.data, p.tok.line})
			pt = p.peekToken()
		}

		if pt.kind == T_L_BLOCK {
			p.nextToken()
			if p.peekToken().kind == T_R_BLOCK {
				n = children.append(n, p.emptyBlock())
			} else {
				n = children.append(n, p.indexUnary())
			}
		}
	}
	return &n
}

func (p *Parser) emptyBlock() *Node {
	FUNC_NAME := "empty block"
	n := Node{N_EMPTY_BLOCK, p.tok.line}

	if p.tok.kind != T_L_BLOCK {
		p.throwError(FUNC_NAME, p.tok.line, "left block", p.tok)
	}
	n = children.append(n, &Node{N_L_BLOCK, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_R_BLOCK {
		p.throwError(FUNC_NAME, p.tok.line, "right block", p.tok)
	}
	n = children.append(n, &Node{N_R_BLOCK, p.tok.data, p.tok.line})
	return &n
}

func (p *Parser) property() *Node {
	FUNC_NAME := "property"
	n := Node{N_PROPERTY, p.tok.line}

	if p.tok.kind != T_IDENTIFIER {
		p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok)
	}
	n = children.append(n, &Node{N_IDENTIFIER, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_ACCESS {
		p.throwError(FUNC_NAME, p.tok.line, "access", p.tok)
	}
	n = children.append(n, &Node{N_ACCESS, p.tok.data, p.tok.line})
	p.nextToken()
	if p.tok.kind != T_IDENTIFIER {
		p.throwError(FUNC_NAME, p.tok.line, "identifier", p.tok)
	}

	if p.peekToken().kind == T_ACCESS {
		n = children.append(n, p.property())
	} else {
		n = children.append(n, &Node{N_IDENTIFIER, p.tok.data, p.tok.line})
	}
	return &n
}

func (p *Parser) throwError(caller string, line int, expected string, got Token) {
	panic("Error in the " + JOB_PARSER + "!\n" + "When the " + JOB_PARSER + " was trying to decipher: " + caller + "\n" + "Error found on line " + string(line) + "\n" + "Expected: " + expected + "\n" + "Got: " + got.String())
}

func (h *Hoister) hoist() []*Node {
	types := &Node{N_PROGRAM}
	consts := &Node{N_PROGRAM}
	funcs := &Node{N_PROGRAM}
	for i := 0; i < len(h.ast.children); i++ {
		c := h.ast.children[i]

		if c.kind == N_NEW_TYPE {
			types = children.append(types, c)
		} else if c.kind == N_STRUCT_DEF {
			types = children.append(types, c)
		} else if c.kind == N_ENUM_DEF {
			types = children.append(types, c)
		} else if c.kind == N_FUNC_DEF {
			funcs = children.append(funcs, c)
		} else if c.kind == N_CONSTANT {
			consts = children.append(consts, c)
		} else {
			continue
		}
		h.ast.children = slices.Delete(h.ast.children, i, i+1)
		i--
	}
	return *Node{types, consts, funcs, h.ast}
}

func (v *Var) String() string {
	return v.kind.String() + " " + v.data + " " + v.datatype + " " + string(v.isArray)
}

func (v VarType) String() string {
	switch v {
	case V_VAR:
		return "VAR"
	case V_TYPE:
		return "TYPE"
	case V_FUNC:
		return "FUNC"
	case V_MAP:
		return "MAP"
	default:
		return "UNKNOWN"
	}
}

func (vs *VarStack) push(n *Var) {
	vf := &VarFrame{n, nil}
	vs.length++

	if vs.length == 1 {
		vs.tail = vf
		return
	}
	vf.prev = vs.tail
	vs.tail = vf
}

func (vs *VarStack) pop() *Var {
	if vs.length == 0 {
		return nil
	}
	vs.length--

	if vs.length == 0 {
		tail := vs.tail
		vs.tail = nil
		return tail.value
	}
	tail := vs.tail
	vs.tail = vs.tail.prev
	return tail.value
}

func (vs *VarStack) peek() *Var {
	if vs.tail == nil {
		return nil
	}
	return vs.tail.value
}

func (ge *GoEmitter) emit() string {
	ge.prePass(ge.funcs)
	ge.prePass(ge.ast)
	output := "package main\n\n"
	output = output + ge.recEmit(ge.types)
	output = output + ge.recEmit(ge.funcs)
	output = output + ge.recEmit(ge.consts)
	output = output + "\ntype float = float64\n"
	output = output + "\nfunc main() {\n"
	output = output + ge.recEmit(ge.ast)
	output = output + "\n}\n"
	return output
}

func (ge *GoEmitter) prePass(n *Node) {
	switch n.kind {
	case N_PROGRAM:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_VAR_DECLARATION:
		ge.prePass(n.children[0])
	case N_IF_BLOCK:
		for i := 1; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_FOREVER_LOOP:
		ge.prePass(n.children[1])
	case N_RANGE_LOOP:
		ge.prePass(n.children[1])
		ge.prePass(n.children[2])
	case N_FOR_LOOP:
		for i := 1; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_WHILE_LOOP:
		ge.prePass(n.children[1])
		ge.prePass(n.children[2])
	case N_FUNC_DEF:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_RET_STATE:
		ge.prePass(n.children[1])
	case N_BREAK_STATE:
		ge.prePass(n.children[1])
	case N_CONT_STATE:
		ge.prePass(n.children[1])
	case N_CONDITION:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_EXPRESSION:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_ASSIGNMENT:
		ge.prePass(n.children[0])
		ge.prePass(n.children[1])
		ge.prePass(n.children[len(n.children)-1])
	case N_LONE_CALL:
		curr := n.children.children[0][1]
		for len(curr.children) == 3 {
			curr = curr.children[2]
		}

		if curr.data == "append" {
			fc := *n.children[0]
			p := *fc.children[1]
			line := n.line
			newChildren := []*Node{}
			for i := 0; i < len(fc.children)-1; i++ {
				newChildren = append(newChildren, fc.children[i])
			}
			fc = children.append(fc, &Node{line, N_SEP, ","})
			fc = children.append(fc, fc.children[3])
			fc = children.append(fc, &Node{line, N_R_PAREN, ")"})
			[3]fc.children = &Node{line, N_EXPRESSION, int{fc.children.children[1][0]}}
			[1]fc.children = fc.children.children[1][2]
			finalParent := Node{line, N_VAR_DECLARATION}
			assign := Node{line, N_ASSIGNMENT}
			finalParent.children = Node{&assign, &Node{line, N_SEMICOLON, ";"}}
			assign.children = Node{p.children[0], &Node{line, N_ASSIGN, "="}, &Node{line, N_EXPRESSION, Node{&fc}}}
			*n = finalParent
		}
		ge.prePass(n.children[0])
	case N_FUNC_CALL:
		for i := 1; i < len(n.children)-1; i++ {
			ge.prePass(n.children[i])
		}
	case N_BLOCK:
		for i := 1; i < len(n.children)-1; i++ {
			ge.prePass(n.children[i])
		}
	case N_NEW_TYPE:
		ge.prePass(n.children[2])
	case N_UNARY_OPERATION:
		ge.prePass(n.children[0])
		ge.prePass(n.children[1])
	case N_MAKE_ARRAY:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_COMPLEX_TYPE:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_SWITCH_STATE:
		for i := 1; i < len(n.children)-1; i++ {
			ge.prePass(n.children[i])
		}
	case N_CASE_STATE:
		ge.prePass(n.children[1])
		ge.prePass(n.children[3])
	case N_DEFAULT_STATE:
		ge.prePass(n.children[2])
	case N_CASE_BLOCK:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_LONE_INC:
		ge.prePass(n.children[1])
	case N_METHOD_RECEIVER:
		ge.prePass(n.children[1])
		ge.prePass(n.children[2])
	case N_ENUM_DEF:
		for i := 1; i < len(n.children)-1; i++ {
			ge.prePass(n.children[i])
		}
	case N_STRUCT_NEW:
		for i := 1; i < len(n.children)-1; i++ {
			ge.prePass(n.children[i])
		}
	case N_BRACKETED_VALUE:
		ge.prePass(n.children[1])
	case N_ELEMENT_ASSIGNMENT:
		ge.prePass(n.children[0])
		ge.prePass(n.children[1])
	case N_STRUCT_DEF:
		for i := 1; i < len(n.children)-1; i++ {
			ge.prePass(n.children[i])
		}
	case N_PROPERTY:
		line := n.line
		var parent *Node
		curr := n
		for len(curr.children) == 3 {
			parent = curr
			curr = curr.children[2]
		}

		if curr.data == "len" {
			newChildren := []*Node{}
			newChildren = append(newChildren, parent.children[1])
			*parent = *parent.children[0]
			value := *n
			p := Node{N_FUNC_CALL, line}
			p.children = Node{&Node{line, N_CALL, "call"}, &Node{line, N_IDENTIFIER, "len"}, &Node{line, N_L_PAREN, "("}, &Node{line, N_EXPRESSION, Node{&value}}, &Node{line, N_R_PAREN, ")"}}
			*n = p
			ge.prePass(n)
			return
		}
		ge.prePass(n.children[0])
		ge.prePass(n.children[2])
	case N_CONSTANT:
		ge.prePass(n.children[1])
	case N_INDEX:
		ge.prePass(n.children[1])
	case N_EMPTY_BLOCK:
	case N_CONST:
	case N_FOR:
	case N_RANGE:
	case N_FOREVER:
	case N_WHILE:
	case N_IF:
	case N_ELIF:
	case N_ELSE:
	case N_CALL:
	case N_STRUCT:
	case N_FUN:
	case N_RET:
	case N_BREAK:
	case N_CONT:
	case N_ENUM:
	case N_TYPEDEF:
	case N_NEW:
	case N_MAKE:
	case N_MAP:
	case N_SWITCH:
	case N_CASE:
	case N_DEFAULT:
	case N_SEMICOLON:
	case N_ASSIGN:
	case N_SEP:
	case N_COLON:
	case N_L_SQUIRLY:
	case N_R_SQUIRLY:
	case N_L_BLOCK:
	case N_R_BLOCK:
	case N_L_PAREN:
	case N_R_PAREN:
	case N_ADD:
	case N_SUB:
	case N_MUL:
	case N_DIV:
	case N_OR:
	case N_AND:
	case N_OROR:
	case N_ANDAND:
	case N_EQ:
	case N_LT:
	case N_GT:
	case N_LTEQ:
	case N_GTEQ:
	case N_NEQ:
	case N_MOD:
	case N_ACCESS:
	case N_XOR:
	case N_L_SHIFT:
	case N_R_SHIFT:
	case N_INC:
	case N_DINC:
	case N_NOT:
	case N_REF:
	case N_DEREF:
	case N_TYPE:
	case N_IDENTIFIER:
		switch n.data {
		case "tobyte":
			n.data = "byte"
		case "toword":
			n.data = "word"
		case "todword":
			n.data = "dword"
		case "toqword":
			n.data = "qword"
		case "touint8":
			n.data = "uint8"
		case "touint16":
			n.data = "uint16"
		case "touint32":
			n.data = "uint32"
		case "touint64":
			n.data = "uint64"
		case "touint":
			n.data = "uint"
		case "toint8":
			n.data = "int8"
		case "toint16":
			n.data = "int16"
		case "toint32":
			n.data = "int32"
		case "toint64":
			n.data = "int64"
		case "tosint":
			n.data = "sint"
		case "toint":
			n.data = "int"
		case "tochar":
			n.data = "char"
		case "tostring":
			n.data = "string"
		case "tofloat32":
			n.data = "float32"
		case "tofloat64":
			n.data = "float64"
		case "todouble":
			n.data = "double"
		case "tofloat":
			n.data = "float"
		case "tobool":
			n.data = "bool"
		}
	case N_INT:
	case N_FLOAT:
	case N_STRING:
	case N_CHAR:
	case N_BOOL:
	case N_NIL:
	}
}

func (ge *GoEmitter) recEmit(n *Node) string {
	var output string
	switch n.kind {
	case N_PROGRAM:
		for i := 0; i < len(n.children); i++ {
			output = output + ge.recEmit(n.children[i])
		}
	case N_VAR_DECLARATION:
		output = output + ge.recEmit(n.children[0]) + ge.recEmit(n.children[1]) + "\n"
	case N_IF_BLOCK:
		output = output + "\n"
		for i := 0; i < len(n.children); i++ {
			output = output + ge.recEmit(n.children[i]) + " "
		}
		output = output + "\n"
	case N_FOREVER_LOOP:
		output = output + ge.recEmit(n.children[0]) + " " + ge.recEmit(n.children[1]) + "\n"
	case N_RANGE_LOOP:
		output = output + ge.recEmit(n.children[0]) + " range " + ge.recEmit(n.children[1]) + " " + ge.recEmit(n.children[2]) + "\n"
	case N_FOR_LOOP:
		for i := 0; i < len(n.children); i++ {
			output = output + ge.recEmit(n.children[i]) + " "
		}
		output = output + "\n"
	case N_WHILE_LOOP:
		output = output + ge.recEmit(n.children[0]) + " " + ge.recEmit(n.children[1]) + " " + ge.recEmit(n.children[2]) + "\n"
	case N_FUNC_DEF:
		for i := 0; i < len(n.children); i++ {
			output = output + ge.recEmit(n.children[i]) + " "
		}
		output = output + "\n\n"
	case N_RET_STATE:
		for i := 0; i < len(n.children); i++ {
			output = output + ge.recEmit(n.children[i]) + " "
		}
		output = output + "\n"
	case N_BREAK_STATE:
		for i := 0; i < len(n.children); i++ {
			output = output + ge.recEmit(n.children[i]) + " "
		}
		output = output + "\n"
	case N_CONT_STATE:
		for i := 0; i < len(n.children); i++ {
			output = output + ge.recEmit(n.children[i]) + " "
		}
		output = output + "\n"
	case N_CONDITION:
		for i := 0; i < len(n.children); i++ {
			output = output + ge.recEmit(n.children[i]) + " "
		}
	case N_EXPRESSION:
		for i := 0; i < len(n.children); i++ {
			output = output + ge.recEmit(n.children[i]) + " "
		}
	case N_ASSIGNMENT:
		isFirstShow := false

		if len(n.children) == 2 {
			output = output + "var " + ge.recEmit(n.children[0]) + " " + ge.recEmit(n.children[1])
		} else {
			for i := 0; i < len(n.children); i++ {
				if n.children.kind[i] == N_COMPLEX_TYPE {
					isFirstShow = true
					ge.varType = n.children[i]
				} else if n.children.kind[i] == N_ASSIGN && isFirstShow {
					if ge.inConst {
						output = output + "= "
					} else {
						output = output + ":= "
					}
				} else {
					output = output + ge.recEmit(n.children[i]) + " "
				}
			}
		}
	case N_LONE_CALL:
		output = output + ge.recEmit(n.children[0]) + " " + ge.recEmit(n.children[1])
	case N_FUNC_CALL:
		for i := 0; i < len(n.children); i++ {
			output = output + ge.recEmit(n.children[i]) + " "
		}
	case N_BLOCK:
		for i := 0; i < len(n.children); i++ {
			output = output + ge.recEmit(n.children[i])
		}
	case N_NEW_TYPE:
		output = output + ge.recEmit(n.children[0]) + " " + ge.recEmit(n.children[1]) + " " + ge.recEmit(n.children[2]) + " " + ge.recEmit(n.children[3]) + "\n"
	case N_UNARY_OPERATION:

		if n.children.kind[0] == N_INC || n.children.kind[0] == N_DINC || n.children.kind[0] == N_INDEX {
			output = output + ge.recEmit(n.children[1]) + ge.recEmit(n.children[0])
		} else {
			output = output + ge.recEmit(n.children[0]) + ge.recEmit(n.children[1])
		}
	case N_MAKE_ARRAY:
		output = output + ge.recEmit(ge.varType)
		output = output + "{"
		for i := 2; i < len(n.children)-1; i++ {
			output = output + ge.recEmit(n.children[i])
		}
		output = output + "}"
	case N_EMPTY_BLOCK:
		output = output + "[]"
	case N_COMPLEX_TYPE:

		if n.children.kind[0] == N_MAP {
			output = "map" + ge.recEmit(n.children[2]) + ge.recEmit(n.children[3]) + ge.recEmit(n.children[4]) + ge.recEmit(n.children[1])
		} else {
			for i := len(n.children) - 1; i >= 0; i-- {
				output = output + ge.recEmit(n.children[i])
			}
		}
	case N_SWITCH_STATE:
		for i := 0; i < len(n.children); i++ {
			output = output + ge.recEmit(n.children[i]) + " "
		}
	case N_CASE_STATE:
		output = output + ge.recEmit(n.children[0]) + " " + ge.recEmit(n.children[1]) + ":\n" + ge.recEmit(n.children[3])
	case N_DEFAULT_STATE:
		output = output + ge.recEmit(n.children[0]) + " " + ":\n" + ge.recEmit(n.children[2])
	case N_CASE_BLOCK:
		for i := 0; i < len(n.children); i++ {
			output = output + ge.recEmit(n.children[i])
		}
	case N_LONE_INC:
		output = output + ge.recEmit(n.children[1]) + ge.recEmit(n.children[0]) + ge.recEmit(n.children[2]) + "\n"
	case N_METHOD_RECEIVER:
		output = output + ge.recEmit(n.children[0]) + ge.recEmit(n.children[1]) + " " + ge.recEmit(n.children[2]) + ge.recEmit(n.children[3])
	case N_ENUM_DEF:
		typeName := ge.recEmit(n.children[1])
		output = output + "type " + typeName + " int\n"
		output = output + "\nconst (\n"
		output = output + ge.recEmit(n.children[3]) + " " + typeName + " = iota\n"
		for i := 5; i < len(n.children)-1; i = i + 2 {
			output = output + ge.recEmit(n.children[i]) + "\n"
		}
		output = output + ")\n"
	case N_STRUCT_NEW:
		output = output + ge.recEmit(n.children[1]) + "{"
		for i := 3; i < len(n.children)-1; i++ {
			output = output + ge.recEmit(n.children[i])
		}
		output = output + "}"
	case N_BRACKETED_VALUE:
		output = output + ge.recEmit(n.children[0]) + ge.recEmit(n.children[1]) + ge.recEmit(n.children[2])
	case N_ELEMENT_ASSIGNMENT:
		output = output + ge.recEmit(n.children[0]) + " " + ge.recEmit(n.children[1]) + ge.recEmit(n.children[2]) + "\n"
	case N_STRUCT_DEF:
		output = output + "type " + ge.recEmit(n.children[1]) + " struct {"
		for i := 3; i < len(n.children)-1; i = i + 2 {
			output = output + "\n" + ge.recEmit(n.children[i]) + " " + ge.recEmit(n.children[i+1])
		}
		output = output + "\n}\n\n"
	case N_PROPERTY:
		output = output + ge.recEmit(n.children[0]) + ge.recEmit(n.children[1]) + ge.recEmit(n.children[2])
	case N_CONSTANT:
		ge.inConst = true
		output = output + ge.recEmit(n.children[0]) + " " + ge.recEmit(n.children[1])
		ge.inConst = false
	case N_INDEX:
		output = output + ge.recEmit(n.children[0]) + ge.recEmit(n.children[1]) + ge.recEmit(n.children[2])
	case N_CONST:
		output = output + "const"
	case N_FOR:
		output = output + "for"
	case N_RANGE:
		output = output + "for"
	case N_FOREVER:
		output = output + "for"
	case N_WHILE:
		output = output + "for"
	case N_IF:
		output = output + "if"
	case N_ELIF:
		output = output + "else if"
	case N_ELSE:
		output = output + "else"
	case N_CALL:
		output = output + ""
	case N_STRUCT:
		output = output + "struct"
	case N_FUN:
		output = output + "func"
	case N_RET:
		output = output + "return"
	case N_BREAK:
		output = output + "break"
	case N_CONT:
		output = output + "continue"
	case N_ENUM:
		output = output + "enum"
	case N_TYPEDEF:
		output = output + "type"
	case N_NEW:
		output = output + "new"
	case N_MAKE:
		output = output + "make"
	case N_MAP:
		output = output + "map"
	case N_SWITCH:
		output = output + "switch"
	case N_CASE:
		output = output + "case"
	case N_DEFAULT:
		output = output + "default"
	case N_SEMICOLON:
		output = output + ";"
	case N_ASSIGN:
		output = output + "="
	case N_SEP:
		output = output + ","
	case N_COLON:
		output = output + ":"
	case N_L_SQUIRLY:
		output = output + "{"
	case N_R_SQUIRLY:
		output = output + "}"
	case N_L_BLOCK:
		output = output + "["
	case N_R_BLOCK:
		output = output + "]"
	case N_L_PAREN:
		output = output + "("
	case N_R_PAREN:
		output = output + ")"
	case N_ADD:
		output = output + "+"
	case N_SUB:
		output = output + "-"
	case N_MUL:
		output = output + "*"
	case N_DIV:
		output = output + "/"
	case N_OR:
		output = output + "|"
	case N_AND:
		output = output + "&"
	case N_OROR:
		output = output + "||"
	case N_ANDAND:
		output = output + "&&"
	case N_EQ:
		output = output + "=="
	case N_LT:
		output = output + "<"
	case N_GT:
		output = output + ">"
	case N_LTEQ:
		output = output + "<="
	case N_GTEQ:
		output = output + ">="
	case N_NEQ:
		output = output + "!="
	case N_MOD:
		output = output + "%"
	case N_ACCESS:
		output = output + "."
	case N_XOR:
		output = output + "^"
	case N_L_SHIFT:
		output = output + "<<"
	case N_R_SHIFT:
		output = output + ">>"
	case N_INC:
		output = output + "++"
	case N_DINC:
		output = output + "--"
	case N_NOT:
		output = output + "!"
	case N_REF:
		output = output + "&"
	case N_DEREF:
		output = output + "*"
	case N_TYPE:
		output = output + n.data
	case N_IDENTIFIER:
		output = output + n.data
	case N_INT:
		output = output + n.data
	case N_FLOAT:
		output = output + n.data
	case N_STRING:
		output = output + n.data
	case N_CHAR:
		output = output + n.data
	case N_BOOL:
		output = output + n.data
	case N_NIL:
		output = output + "nil"
	default:
		panic("Bad node in emitter?")
	}
	return output
}

func (ge *GoEmitter) dump(emitted string) {
	err := os.WriteFile("../output/main.go", emitted, 0644)

	if err != nil {
		panic(err)
	}
}

const JOB_LEXER = "Lexer"
const JOB_PARSER = "Parser"
const JOB_HOISTER = "Hoister"
const JOB_GO_EMITTER = "Go Emitter"

type float = float64

func main() {

}
