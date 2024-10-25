package main

import (
	"fmt"
	"strconv"
)

const JOB_LEXER = "Lexer"

type Lexer struct {
	source []byte
	index  int
	line   int
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
		// Automatically sets it to ILLEGAL
		token = Token{}

		switch l.source[l.index] {

		// Unnecessary to compute
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
			token = Token{"%", T_MUL, l.line}

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

		case '+':
			token = Token{"+", T_ADD, l.line}
			if l.peekChar() == '+' {
				token = Token{"++", T_INC, l.line}
				l.index++
			}

		case '-':
			token = Token{"-", T_SUB, l.line}
			if l.peekChar() == '-' {
				token = Token{"--", T_SUB, l.line}
				l.index++
			}

		case '/':
			token = Token{"/", T_DIV, l.line}

			// single-line comment
			if l.peekChar() == '/' {
				for l.peekChar() != '\n' && l.peekChar() != 0 {

					// Make sure line numbers make sense
					if l.peekChar() == '\n' {
						l.line++
					}

					l.index++
				}
				continue
			}

			// multi-line comment
			if l.peekChar() == '*' {
				for l.peekChar() != 0 {

					// Make sure line numbers make sense
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

		case '\'': // Characters
			token = Token{"'", T_CHAR, l.line}
			// TODO: Deal with the case when you have escaped characters, such as newline?

			if l.peekChar() == '\\' {
				l.index += 3
				if l.index >= len(l.source) {
					// TODO: Better error messages
					panic("Expected more source")
				}

				if l.source[l.index] != '\'' {
					fmt.Println(tokens)
					panic("Expected ' (line " + strconv.Itoa(l.line) + ") got " + string(l.source[l.index]))
				}
				token.data += string(l.source[l.index-2]) + string(l.source[l.index-1]) + string(l.source[l.index])
			} else {
				l.index += 2
				if l.index >= len(l.source) {
					// TODO: Better error messages
					panic("Expected more source")
				}

				if l.source[l.index] != '\'' {
					panic("Expected ' (line " + strconv.Itoa(l.line) + ") got " + string(l.source[l.index]))
				}
				token.data += string(l.source[l.index-1]) + string(l.source[l.index])
			}

		case '"': // Strings

			// TODO: Have to deal with comments in string (maybe)

			token = Token{"\"", T_STRING, l.line}

			escaped := false
			for !(l.peekChar() == '"' && !escaped) && l.peekChar() != 0 {
				l.index++
				token.data += string(l.source[l.index])
				if l.source[l.index] == '\\' {
					escaped = true
				} else {
					escaped = false
				}
			}

			if l.peekChar() != '"' {
				panic("Couldn't find end of string")
			}

			l.index++
			token.data += string(l.source[l.index])
		}

		// Numbers
		if l.source[l.index]-'0' <= 9 {
			// We don't know what type of number
			// yet, but we'll figure it out later.
			// At least we know it can only be
			// either integer, float, or illegal
			token = Token{data: string(l.source[l.index]), line: l.line}

			isFloat := false

			for (l.peekChar()-'0' <= 9 || l.peekChar() == '.') && l.peekChar() != 0 {
				l.index++
				token.data += string(l.source[l.index])

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

			// Keywords
			// Other words
		} else if l.source[l.index]-'a' < 26 || l.source[l.index]-'A' < 26 {
			// Don't know specific type yet
			token = Token{data: string(l.source[l.index]), line: l.line}

			for (l.peekChar()-'a' < 26 || l.peekChar()-'A' < 26 || l.peekChar() == '_') && l.peekChar() != 0 {
				l.index++
				token.data += string(l.source[l.index])
			}

			// Keywords
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
			case "switch":
				token.kind = T_SWITCH
			case "case":
				token.kind = T_CASE
			case "default":
				token.kind = T_DEFAULT
			case "true", "false":
				token.kind = T_BOOL

				// Types
			case "byte", "word", "dword", "qword",
				"uint8", "uint16", "uint32", "uint64",
				"uint", "int8", "int16", "int32",
				"int64", "sint", "int", "char",
				"string", "float32", "float64", "double",
				"float", "bool", "any":
				token.kind = T_TYPE

				// Identifiers
			default:
				// HOWEVER, the programmer is allowed
				// to create new types, and the 2 cases
				// that this occurs is in struct and
				// typdefs. It could also be a type if
				// the previous token is new, but the
				// programmer could be lying, so this
				// will be up to semantic analysis. The
				// last case is an identifier after
				// an identifier, which is usually
				// where a type is, so we will also add
				// this case

				// First token in the file (therefore
				// can't be preceeded with the relevant
				// tokens)
				if len(tokens) == 0 {
					token.kind = T_IDENTIFIER
					break
				}

				// Could possibly be a type
				prevToken := tokens[len(tokens)-1]
				if prevToken.kind == T_STRUCT ||
					prevToken.kind == T_TYPEDEF ||
					prevToken.kind == T_NEW ||
					prevToken.kind == T_IDENTIFIER ||
					prevToken.kind == T_MAP {
					token.kind = T_TYPE
				} else {
					token.kind = T_IDENTIFIER
				}
			}
		}

		if token.kind == T_ILLEGAL {
			// TODO: Better error messages
			throwError(JOB_LEXER, "lexing", l.line, "anything else", "ILLEGAL ("+string(l.source[l.index])+")")
		}

		tokens = append(tokens, token)
	}

	return tokens
}
