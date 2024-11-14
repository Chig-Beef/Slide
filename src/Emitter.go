package main

import (
	"fmt"
	"os"
)

const JOB_GO_EMITTER = "Go Emitter"

type GoEmitter struct {
	types  *Node
	consts *Node
	funcs  *Node
	ast    *Node

	varType *Node // Used on making arrays
	inConst bool  // Used to not use the walrus operator when assigning to const
}

func (ge *GoEmitter) emit() string {
	// Translate things, to the correct
	// thing for Golang
	ge.prePass(ge.funcs)
	ge.prePass(ge.ast)

	output := "package main\n\n"

	output += ge.recEmit(ge.types)
	output += ge.recEmit(ge.funcs)
	output += ge.recEmit(ge.consts)

	output += "\ntype float = float64\n"

	output += "\nfunc main() {\n"
	output += ge.recEmit(ge.ast)
	output += "\n}\n"

	return output
}

func (ge *GoEmitter) prePass(n *Node) {
	switch n.kind {
	case N_PROGRAM:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_VAR_DECLARATION:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_IF_BLOCK:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_FOREVER_LOOP:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_RANGE_LOOP:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_FOR_LOOP:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_WHILE_LOOP:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_FUNC_DEF:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_RET_STATE:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_BREAK_STATE:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_CONT_STATE:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_CONDITION:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_EXPRESSION:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_ASSIGNMENT:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_LONE_CALL:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_FUNC_CALL:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_BLOCK:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_NEW_TYPE:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_UNARY_OPERATION:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_MAKE_ARRAY:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_EMPTY_BLOCK:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_COMPLEX_TYPE:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_SWITCH_STATE:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_CASE_STATE:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_DEFAULT_STATE:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_CASE_BLOCK:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_LONE_INC:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_METHOD_RECEIVER:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_ENUM_DEF:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_STRUCT_NEW:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_BRACKETED_VALUE:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_ELEMENT_ASSIGNMENT:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_STRUCT_DEF:
		for i := 0; i < len(n.children); i++ {
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

		// Now curr should be the identifier at
		// the end
		if curr.data != "len" {
			return
		}

		// Get rid of ".len"
		parent.children = parent.children[:1]
		// Is a single identifier rather than a
		// property, so just clean up the tree
		// a bit
		*parent = *parent.children[0]

		value := *n

		// The function call that will be used
		p := Node{kind: N_FUNC_CALL, line: line}
		p.children = []*Node{
			{line: line, kind: N_CALL, data: "call"},
			{line: line, kind: N_IDENTIFIER, data: "len"},
			{line: line, kind: N_L_PAREN, data: "("},
			{line: line, kind: N_EXPRESSION, children: []*Node{&value}},
			{line: line, kind: N_R_PAREN, data: ")"},
		}

		// Place the new node back in the AST
		*n = p
	case N_CONSTANT:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_INDEX:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
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
	case N_INT:
	case N_FLOAT:
	case N_STRING:
	case N_CHAR:
	case N_BOOL:
	case N_NIL:
	}
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
	case N_IF_BLOCK:
		output = "\n"
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
	case N_FUNC_DEF:
		for i := range n.children {
			output += ge.recEmit(n.children[i]) + " "
		}
		output += "\n\n"
	case N_RET_STATE:
		for i := range n.children {
			output += ge.recEmit(n.children[i]) + " "
		}
		output += "\n"
	case N_BREAK_STATE:
		for i := range n.children {
			output += ge.recEmit(n.children[i]) + " "
		}
		output += "\n"
	case N_CONT_STATE:
		for i := range n.children {
			output += ge.recEmit(n.children[i]) + " "
		}
		output += "\n"
	case N_CONDITION:
		for i := range n.children {
			output += ge.recEmit(n.children[i]) + " "
		}
	case N_EXPRESSION:
		for i := range n.children {
			output += ge.recEmit(n.children[i]) + " "
		}
	case N_ASSIGNMENT:
		isFirstShow := false

		if len(n.children) == 2 {
			// No definition, just declaration
			output = "var " + ge.recEmit(n.children[0]) + " " + ge.recEmit(n.children[1])

		} else {
			for i := 0; i < len(n.children); i++ {
				if n.children[i].kind == N_COMPLEX_TYPE {
					isFirstShow = true
					ge.varType = n.children[i]
				} else if n.children[i].kind == N_ASSIGN && isFirstShow {
					if ge.inConst {
						output += "= "
					} else {
						output += ":= "
					}
				} else {
					output += ge.recEmit(n.children[i]) + " "
				}
			}
		}
	case N_LONE_CALL:
		for i := range n.children {
			output += ge.recEmit(n.children[i]) + " "
		}
	case N_FUNC_CALL:
		for i := range n.children {
			output += ge.recEmit(n.children[i]) + " "
		}
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
		if n.children[0].kind == N_INC || n.children[0].kind == N_DINC || n.children[0].kind == N_INDEX {
			for i := 1; i < len(n.children); i++ {
				output += ge.recEmit(n.children[i])
			}
			output += ge.recEmit(n.children[0])
		} else {
			for i := range n.children {
				output += ge.recEmit(n.children[i])
			}
		}
	case N_MAKE_ARRAY:
		output = ge.recEmit(ge.varType)

		output += "{"
		for i := 2; i < len(n.children)-1; i++ {
			output += ge.recEmit(n.children[i])
		}
		output += "}"
	case N_EMPTY_BLOCK:
		for i := 0; i < len(n.children); i++ {
			output += ge.recEmit(n.children[i])
		}
	case N_COMPLEX_TYPE:
		if n.children[0].kind == N_MAP {
			output = "map" +
				ge.recEmit(n.children[2]) +
				ge.recEmit(n.children[3]) +
				ge.recEmit(n.children[4]) +
				ge.recEmit(n.children[1])
		} else {
			for i := len(n.children) - 1; i >= 0; i-- {
				output += ge.recEmit(n.children[i])
			}
		}
	case N_SWITCH_STATE:
		for i := 0; i < len(n.children); i++ {
			output += ge.recEmit(n.children[i]) + " "
		}
	case N_CASE_STATE:
		for i := 0; i < len(n.children); i++ {
			output += ge.recEmit(n.children[i]) + " "
		}
	case N_DEFAULT_STATE:
		for i := 0; i < len(n.children); i++ {
			output += ge.recEmit(n.children[i])
		}
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
	case N_ENUM_DEF: // TODO: What happens when we get an empty enum?
		typeName := ge.recEmit(n.children[1])
		output = "type " + typeName + " int\n"

		output += "\nconst (\n"
		output += ge.recEmit(n.children[3]) + " " + typeName + " = iota\n"
		for i := 5; i < len(n.children)-1; i += 2 {
			output += ge.recEmit(n.children[i]) + "\n"
		}
		output += ")\n"
	case N_STRUCT_NEW:
		output = ge.recEmit(n.children[1]) + "{"

		for i := 3; i < len(n.children)-1; i++ {
			output += ge.recEmit(n.children[i])
		}

		output += "}"
	case N_BRACKETED_VALUE:
		for i := 0; i < len(n.children); i++ {
			output += ge.recEmit(n.children[i])
		}
	case N_ELEMENT_ASSIGNMENT:
		output = ge.recEmit(n.children[1].children[0]) + " " + ge.recEmit(n.children[0])
		for i := 1; i < len(n.children[1].children); i++ {
			output += " " + ge.recEmit(n.children[1].children[i])
		}
		output += "\n"
	case N_STRUCT_DEF:
		// First line
		output = "type " + ge.recEmit(n.children[1]) + " struct {"
		// Props
		for i := 3; i < len(n.children)-1; i += 2 {
			output += "\n" + ge.recEmit(n.children[i]) + " " + ge.recEmit(n.children[i+1])
		}
		// Closing
		output += "\n}\n\n"
	case N_PROPERTY:
		// Check for a len property
		if n.children[2].data == "len" {
			output += "len(" + ge.recEmit(n.children[0]) + ")"
		} else {
			for i := 0; i < len(n.children); i++ {
				output += ge.recEmit(n.children[i])
			}
		}
	case N_CONSTANT:
		ge.inConst = true
		for i := 0; i < len(n.children); i++ {
			output += ge.recEmit(n.children[i]) + " "
		}
		ge.inConst = false
	case N_INDEX:
		for i := 0; i < len(n.children); i++ {
			output += ge.recEmit(n.children[i])
		}
	case N_CONST:
		output = "const"
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
		output = "="
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
