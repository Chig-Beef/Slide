package main

import "os"

const JOB_GO_EMITTER = "Go Emitter"

type GoEmitter struct {
	types *Node
	funcs *Node
	ast   *Node
}

func (ge *GoEmitter) emit() string {
	output := "package main\n\n"

	output += ge.recEmit(ge.types)
	output += ge.recEmit(ge.funcs)

	output += "\nfunc main() {\n"
	output += ge.recEmit(ge.ast)
	output += "\n}\n"

	return output
}

func (ge *GoEmitter) recEmit(n *Node) string {
	output := ""

	switch n.kind {
	case N_PROGRAM:
		for i := 0; i < len(n.children); i++ {
			output += ge.recEmit(n.children[i])
		}
	case N_VAR_DECLARATION:
		for i := 0; i < len(n.children); i++ {
			output += ge.recEmit(n.children[i])
		}
		output += "\n"
	case N_ELEMENT_ASSIGNMENT:
	case N_IF_BLOCK:
		for i := 0; i < len(n.children); i++ {
			output += ge.recEmit(n.children[i]) + " "
		}
		output += "\n"
	case N_FOREVER_LOOP:
		for i := range n.children {
			output += ge.recEmit(n.children[i]) + " "
		}
		output += "\n"
	case N_RANGE_LOOP:
		output = ge.recEmit(n.children[0]) + " range " +
			ge.recEmit(n.children[1]) + " " +
			ge.recEmit(n.children[2])
		output += "\n"
	case N_FOR_LOOP:
		for i := range n.children {
			output += ge.recEmit(n.children[i]) + " "
		}
		output += "\n"
	case N_WHILE_LOOP:
		for i := range n.children {
			output += ge.recEmit(n.children[i]) + " "
		}
		output += "\n"
	case N_STRUCT_DEF:
	case N_FUNC_DEF:
	case N_RET_STATE:
	case N_BREAK_STATE:
	case N_CONT_STATE:
	case N_ENUM_DEF:
	case N_CONDITION:
		for i := range n.children {
			output += ge.recEmit(n.children[i]) + " "
		}
	case N_EXPRESSION:
		for i := range n.children {
			output += ge.recEmit(n.children[i]) + " "
		}
	case N_ASSIGNMENT:
		for i := range n.children {
			output += ge.recEmit(n.children[i]) + " "
		}
	case N_LONE_CALL:
		for i := range n.children {
			output += ge.recEmit(n.children[i]) + " "
		}
	case N_FUNC_CALL:
		for i := range n.children {
			output += ge.recEmit(n.children[i]) + " "
		}
	case N_STRUCT_NEW:
	case N_BLOCK:
		for i := 0; i < len(n.children); i++ {
			output += ge.recEmit(n.children[i])
		}
	case N_NEW_TYPE:
		for i := range n.children {
			output += ge.recEmit(n.children[i]) + " "
		}
		output += "\n"
	case N_UNARY_OPERATION:
		if n.children[0].kind == N_INC || n.children[0].kind == N_DINC {
			for i := 1; i < len(n.children); i++ {
				output += ge.recEmit(n.children[i])
			}
			output += ge.recEmit(n.children[0])
		} else {
			for i := range n.children {
				output += ge.recEmit(n.children[i])
			}
		}
	case N_BRACKETED_VALUE:
	case N_MAKE_ARRAY:
		output += "{"
		for i := 2; i < len(n.children)-1; i++ {
			output += ge.recEmit(n.children[i])
		}
		output += "}"
	case N_COMPLEX_TYPE:
	case N_SWITCH_STATE:
	case N_CASE_STATE:
	case N_DEFAULT_STATE:
	case N_CASE_BLOCK:
		for i := 0; i < len(n.children); i++ {
			output += ge.recEmit(n.children[i])
		}
	case N_LONE_INC:
		for i := 1; i < len(n.children)-1; i++ {
			output += ge.recEmit(n.children[i])
		}
		output += ge.recEmit(n.children[0])
		output += ge.recEmit(n.children[len(n.children)-1])
		output += "\n"
	case N_METHOD_RECEIVER:
		output = ge.recEmit(n.children[0]) +
			ge.recEmit(n.children[1]) + " " +
			ge.recEmit(n.children[2]) +
			ge.recEmit(n.children[3])

	case N_FOR:
		output = "for"
	case N_RANGE:
		output = "for"
	case N_FOREVER:
		output = "for"
	case N_WHILE:
		output = "for"
	case N_IF:
		output = "if"
	case N_ELIF:
		output = "else if"
	case N_ELSE:
		output = "else"
	case N_CALL:
		output = ""
	case N_STRUCT:
		output = "struct"
	case N_FUN:
		output = "func"
	case N_RET:
		output = "return"
	case N_BREAK:
		output = "break"
	case N_CONT:
		output = "continue"
	case N_ENUM:
		output = "enum"
	case N_TYPEDEF:
		output = "type"
	case N_NEW:
		output = "new"
	case N_MAKE:
		output = "make"
	case N_MAP:
		output = "map"
	case N_SWITCH:
		output = "switch"
	case N_CASE:
		output = "case"
	case N_DEFAULT:
		output = "default"
	case N_SEMICOLON:
		output = ";"
	case N_ASSIGN:
		output = ":="
	case N_SEP:
		output = ","
	case N_COLON:
		output = ":"
	case N_L_SQUIRLY:
		output = "{"
	case N_R_SQUIRLY:
		output = "}"
	case N_L_BLOCK:
		output = "["
	case N_R_BLOCK:
		output = "]"
	case N_L_PAREN:
		output = "("
	case N_R_PAREN:
		output = ")"
	case N_ADD:
		output = "+"
	case N_SUB:
		output = "-"
	case N_MUL:
		output = "*"
	case N_DIV:
		output = "/"
	case N_OR:
		output = "|"
	case N_AND:
		output = "&"
	case N_OROR:
		output = "||"
	case N_ANDAND:
		output = "&&"
	case N_EQ:
		output = "=="
	case N_LT:
		output = "<"
	case N_GT:
		output = ">"
	case N_LTEQ:
		output = "<="
	case N_GTEQ:
		output = ">="
	case N_NEQ:
		output = "!="
	case N_MOD:
		output = "%"
	case N_ACCESS:
		output = "."
	case N_XOR:
		output = "^"
	case N_L_SHIFT:
		output = "<<"
	case N_R_SHIFT:
		output = ">>"
	case N_INC:
		output = "++"
	case N_DINC:
		output = "--"
	case N_NOT:
		output = "!"
	case N_REF:
		output = "&"
	case N_DEREF:
		output = "*"
	case N_INDEX:
		for i := 0; i < len(n.children); i++ {
			output += ge.recEmit(n.children[i])
		}
	case N_TYPE:
		output = n.data
	case N_IDENTIFIER:
		output = n.data
	case N_INT:
		output = n.data
	case N_FLOAT:
		output = n.data
	case N_STRING:
		output = n.data
	case N_CHAR:
		output = n.data
	case N_BOOL:
		output = n.data
	case N_NIL:
		output = "nil"

	default:
		panic("Bad node in emitter?")
	}

	return output
}

func (ge *GoEmitter) dump(emitted string) {
	err := os.WriteFile("../output/main.go", []byte(emitted), 0644)
	if err != nil {
		panic(err)
	}
}
