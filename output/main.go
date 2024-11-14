package main

type Lexer struct {
	source []byte
	index  int
	line   int
}

type TokenType byte
type Token struct {
	data string
	kind TokenType
	line int
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
	source := l.source
	sourceLen := len(source)
	for ; l.index < sourceLen; l.index = l.index + 1 {
		token = Token{}
		switch l.source[l.index] {
		case ' ':
			continue
		case '\t':
			continue
		case '\n':
			l.line = l.line + 1
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
		case '|':
			token = Token{"|", T_OR, l.line}
		case '&':
			token = Token{"&", T_AND, l.line}
		case ',':
			token = Token{",", T_SEP, l.line}
		case '%':
			token = Token{"%", T_MUL, l.line}
		case '*':
			token = Token{"*", T_MUL, l.line}
		case '=':
			token = Token{"=", T_ASSIGN, l.line}

			if l.peekChar() == '=' {
				token = Token{"==", T_EQ, l.line}
				l.index = l.index + 1
			}
		case '!':
			token = Token{"!", T_NOT, l.line}

			if l.peekChar() == '=' {
				token = Token{"!=", T_NEQ, l.line}
				l.index = l.index + 1
			}
		case '<':
			token = Token{"<", T_LT, l.line}

			if l.peekChar() == '=' {
				token = Token{"<=", T_LTEQ, l.line}
				l.index = l.index + 1
			} else if l.peekChar() == '<' {
				token = Token{"<<", T_L_SHIFT, l.line}
				l.index = l.index + 1
			}
		case '>':
			token = Token{">", T_GT, l.line}

			if l.peekChar() == '=' {
				token = Token{">=", T_GTEQ, l.line}
				l.index = l.index + 1
			} else if l.peekChar() == '>' {
				token = Token{">>", T_R_SHIFT, l.line}
				l.index = l.index + 1
			}
		case '+':
			token = Token{"+", T_ADD, l.line}

			if l.peekChar() == '+' {
				token = Token{"++", T_INC, l.line}
				l.index = l.index + 1
			}
		case '-':
			token = Token{"-", T_SUB, l.line}

			if l.peekChar() == '-' {
				token = Token{"--", T_SUB, l.line}
				l.index = l.index + 1
			}
		case '/':
			token = Token{"/", T_DIV, l.line}

			if l.peekChar() == '/' {
				for l.peekChar() != '\n' && l.peekChar() != 0 {
					if l.peekChar() == '\n' {
						l.line = l.line + 1
					}
					l.index = l.index + 1
				}
				continue
			}

			if l.peekChar() == '*' {
				for l.peekChar() != 0 {
					if l.peekChar() == '\n' {
						l.line = l.line + 1
					}

					if l.peekChar() == '/' && l.source[l.index] == '*' {
						l.index = l.index + 1
						break
					}
					l.index = l.index + 1
				}
				continue
			}
		case '\'':
			token = Token{"'", T_CHAR, l.line}
			l.index = l.index + 1

			if l.index == sourceLen {
				panic("Expected more source")
			}

			if l.source[l.index] != '\'' {
				panic("Expected '")
			}
			token.data = token.data + tostring(l.source[l.index-1]) + tostring(l.source[l.index])
		case '"':
			token = Token{"\"", T_STRING, l.line}
			for l.peekChar() != '"' && l.peekChar() != 0 {
				l.index = l.index + 1
				token.data = token.data + tostring(l.source[l.index])
			}

			if l.peekChar() != '"' {
				panic("Couldn't find end of string")
			}
			l.index = l.index + 1
			token.data = token.data + tostring(l.source[l.index])
		}
		if l.source[l.index]-'0' <= 9 {
			token = Token{tostring(l.source[l.index]), T_ILLEGAL, l.line}
			isFloat := false
			for (l.peekChar()-'0' <= 9 || l.peekChar() == '.') && l.peekChar() != 0 {
				l.index = l.index + 1
				token.data = token.data + tostring(l.source[l.index])

				if l.source[l.index] == '.' {
					isFloat = true
				}
			}

			if isFloat {
				tokenData := token.data

				if token.data[len(tokenData)-1] == '.' {
					panic("Floats cannot end with '.'")
				}
				token.kind = T_FLOAT
			} else {
				token.kind = T_INT
			}
		} else if l.source[l.index]-'a' < 26 || l.source[l.index]-'A' < 26 {
			token = Token{tostring(l.source[l.index]), T_ILLEGAL, l.line}
			for (l.peekChar()-'a' < 26 || l.peekChar()-'A' < 26 || l.peekChar() == '_') && l.peekChar() != 0 {
				l.index = l.index + 1
				token.data = token.data + tostring(l.source[l.index])
			}
			switch token.data {
			case "for":
				token.kind = T_FOR
			case "range":
				token.kind = T_RANGE
			case "forever":
				token.kind = T_FOREVER
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
			throwError(JOB_LEXER, "lexing", l.line, "anything else", "ILLEGAL ("+tostring(l.source[l.index])+")")
		}
		tokens.append(token)
	}
	return tokens
}

func throwError(job string, caller string, line int, expected string, got any) {
	panic("Error in the " + job + "!\n" + "When the " + job + " was trying to decipher: " + caller + "\n" + "Error found on line " + tostring(line) + "\n" + "Expected: " + expected + "\n" + "Got: " + tostring(got))
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

const JOB_LEXER = "Lexer"
const T_ILLEGAL = 0
const T_FOR = 1
const T_RANGE = 2
const T_FOREVER = 3
const T_WHILE = 4
const T_IF = 5
const T_ELIF = 6
const T_ELSE = 7
const T_CALL = 8
const T_STRUCT = 9
const T_FUN = 10
const T_RET = 11
const T_BREAK = 12
const T_CONT = 13
const T_ENUM = 14
const T_TYPEDEF = 15
const T_NEW = 16
const T_MAKE = 17
const T_MAP = 18
const T_SWITCH = 19
const T_CASE = 20
const T_DEFAULT = 21
const T_SEMICOLON = 22
const T_ASSIGN = 23
const T_SEP = 24
const T_COLON = 25
const T_ADD = 26
const T_SUB = 27
const T_MUL = 28
const T_DIV = 29
const T_OR = 30
const T_AND = 31
const T_OROR = 32
const T_ANDAND = 33
const T_EQ = 34
const T_LT = 35
const T_GT = 36
const T_LTEQ = 37
const T_GTEQ = 38
const T_NEQ = 39
const T_MOD = 40
const T_XOR = 41
const T_ACCESS = 42
const T_NOT = 43
const T_INC = 44
const T_DINC = 45
const T_REF = 46
const T_DEREF = 47
const T_L_SHIFT = 48
const T_R_SHIFT = 49
const T_L_SQUIRLY = 50
const T_R_SQUIRLY = 51
const T_L_BLOCK = 52
const T_R_BLOCK = 53
const T_L_PAREN = 54
const T_R_PAREN = 55
const T_TYPE = 56
const T_IDENTIFIER = 57
const T_INT = 58
const T_FLOAT = 59
const T_STRING = 60
const T_CHAR = 61
const T_BOOL = 62
const T_NIL = 63

type float = float64

func main() {

}
