package main

import "fmt"

type Lexer struct {
	source []byte
	index  int
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
			continue
		case '\r':
			continue

		case ';':
			token = Token{";", T_SEMICOLON}

		case '{':
			token = Token{"{", T_L_SQUIRLY}

		case '}':
			token = Token{"}", T_R_SQUIRLY}

		case '[':
			token = Token{"[", T_L_BLOCK}

		case ']':
			token = Token{"]", T_R_BLOCK}

		case '=':
			token = Token{"=", T_ASSIGN}
			if l.peekChar() == '=' {
				token = Token{"==", T_EQ}
				l.index++
			}

		case '!':
			token = Token{"!", T_NOT}
			if l.peekChar() == '=' {
				token = Token{"!=", T_NEQ}
				l.index++
			}

		case '%':
			token = Token{"%", T_MUL}

		case '*':
			token = Token{"*", T_MUL}

		case '/': // TODO: Implement comments
			token = Token{"/", T_DIV}

		case '|':
			token = Token{"|", T_OR}

		case '&':
			token = Token{"&", T_AND}

		case ',':
			token = Token{",", T_SEP}

		case '<':
			token = Token{"<", T_LT}
			if l.peekChar() == '=' {
				token = Token{"<=", T_LTEQ}
				l.index++
			}

		case '>':
			token = Token{">", T_GT}
			if l.peekChar() == '=' {
				token = Token{">=", T_GTEQ}
				l.index++
			}

		case '+':
			token = Token{"+", T_ADD}
			if l.peekChar() == '+' {
				token = Token{"++", T_INC}
				l.index++
			}

		case '-':
			token = Token{"-", T_SUB}
			if l.peekChar() == '-' {
				token = Token{"--", T_SUB}
				l.index++
			}

		case '\'': // Characters
			token = Token{"'", T_CHAR}

			l.index += 2
			if l.index == len(l.source) {
				// TODO: Better error messages
				panic("Expected more source")
			}

			if l.source[l.index] != '\'' {
				panic("Expected '")
			}

			token.data += string(l.source[l.index-1]) + string(l.source[l.index])

		case '"': // Strings

			// TODO: Have to deal with comments in string (maybe)

			token = Token{"\"", T_STRING}

			for l.peekChar() != '"' && l.peekChar() != 0 {
				l.index++
				token.data += string(l.source[l.index])
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
			token = Token{data: string(l.source[l.index])}

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
		} else if l.source[l.index]-'a' <= 26 || l.source[l.index]-'A' <= 26 {
			// Don't know specific type yet
			token = Token{data: string(l.source[l.index])}

			for (l.source[l.index]-'a' <= 26 || l.source[l.index]-'A' <= 26 || l.source[l.index] == '_') && l.peekChar() != 0 {
				l.index++
				token.data += string(l.source[l.index])
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
			case "true", "false":
				token.kind = T_BOOL
			case "byte", "word", "dword", "qword",
				"uint8", "uint16", "uint32", "uint64",
				"uint", "int8", "int16", "int32",
				"int64", "sint", "int", "char",
				"string", "float32", "float64", "double",
				"float", "bool":
				token.kind = T_TYPE
			default:
				token.kind = T_IDENTIFIER
			}
		}

		if token.kind == T_ILLEGAL {
			// TODO: Better error messages
			fmt.Println(l.source[l.index], string(l.source[l.index]))
			panic("ILLEGAL token found")
		}

		tokens = append(tokens, token)
	}

	return tokens
}
