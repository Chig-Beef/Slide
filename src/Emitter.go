package main

import (
	"os"
	"strings"
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
		// Only need to do the assignment
		ge.prePass(n.children[0])
	case N_IF_BLOCK:
		// Skip over that first if keyword
		for i := 1; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_FOREVER_LOOP:
		// Only do the block
		ge.prePass(n.children[1])
	case N_RANGE_LOOP:
		// Expression
		ge.prePass(n.children[1])
		// Block
		ge.prePass(n.children[2])
	case N_FOR_LOOP:
		for i := 1; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_WHILE_LOOP:
		// Condition
		ge.prePass(n.children[1])
		// Block
		ge.prePass(n.children[2])
	case N_FUNC_DEF:
		for i := 0; i < len(n.children); i++ {
			ge.prePass(n.children[i])
		}
	case N_RET_STATE:
		// Expression
		ge.prePass(n.children[1])
	case N_BREAK_STATE:
		// Value
		ge.prePass(n.children[1])
	case N_CONT_STATE:
		// Value
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
		// Everything but the assign
		ge.prePass(n.children[0])
		ge.prePass(n.children[1])
		ge.prePass(n.children[len(n.children)-1])
	case N_LONE_CALL:
		// Using append?
		curr := n.children[0].children[1]
		for len(curr.children) == 3 {
			curr = curr.children[2]
		}

		// TODO: What if it's var.property.append?
		if curr.data == "append" {
			fc := *n.children[0]
			p := *fc.children[1]

			// Start constructing new tree
			line := n.line

			// Move the item into the correct position
			fc.children = fc.children[:len(fc.children)-1]
			fc.children = append(fc.children, &Node{line: line, kind: N_SEP, data: ","})
			fc.children = append(fc.children, fc.children[3])
			fc.children = append(fc.children, &Node{line: line, kind: N_R_PAREN, data: ")"})

			// Move the list into the correct positino
			fc.children[3] = &Node{line: line, kind: N_EXPRESSION, children: []*Node{
				fc.children[1].children[0],
			}}

			// Move the append into the correct position
			fc.children[1] = fc.children[1].children[2]

			finalParent := Node{line: line, kind: N_VAR_DECLARATION}
			assign := Node{line: line, kind: N_ASSIGNMENT}
			finalParent.children = []*Node{
				&assign,
				{line: line, kind: N_SEMICOLON, data: ";"},
			}

			assign.children = []*Node{
				p.children[0],
				{line: line, kind: N_ASSIGN, data: "="},
				{line: line, kind: N_EXPRESSION, children: []*Node{
					&fc,
				}},
			}

			// Apply the changes
			*n = finalParent
		}

		// In both the regular lone call case
		// and the append case this works
		ge.prePass(n.children[0])
	case N_FUNC_CALL:
		// Trim first and last
		for i := 1; i < len(n.children)-1; i++ {
			ge.prePass(n.children[i])
		}
	case N_BLOCK:
		// Trim first and last
		for i := 1; i < len(n.children)-1; i++ {
			ge.prePass(n.children[i])
		}
	case N_NEW_TYPE:
		// Only check the type that we are
		// aliasing
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
		// Only check thingy and case block
		ge.prePass(n.children[1])
		ge.prePass(n.children[3])
	case N_DEFAULT_STATE:
		// Only case block
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

		// Now curr should be the identifier at
		// the end
		if curr.data == "len" {
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

			// Retry as the correct type
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
	// Faster string manipulation, instead
	// of creating a copy on each iteration
	// of a loop, it just mutates the
	// slice.
	output := strings.Builder{}

	switch n.kind {
	case N_PROGRAM:
		output.Grow(20)
		for i := 0; i < len(n.children); i++ {
			output.WriteString(ge.recEmit(n.children[i]))
		}
	case N_VAR_DECLARATION:
		output.WriteString(ge.recEmit(n.children[0]) + ge.recEmit(n.children[1]) + "\n")
	case N_IF_BLOCK:
		output.WriteString("\n")
		for i := 0; i < len(n.children); i++ {
			output.WriteString(ge.recEmit(n.children[i]) + " ")
		}
		output.WriteString("\n")
	case N_FOREVER_LOOP:
		output.WriteString(ge.recEmit(n.children[0]) + " " +
			ge.recEmit(n.children[1]) + "\n")
	case N_RANGE_LOOP:
		output.WriteString(ge.recEmit(n.children[0]) + " range " +
			ge.recEmit(n.children[1]) + " " +
			ge.recEmit(n.children[2]) + "\n")
	case N_FOR_LOOP:
		for i := range n.children {
			output.WriteString(ge.recEmit(n.children[i]) + " ")
		}
		output.WriteString("\n")
	case N_WHILE_LOOP:
		output.WriteString(ge.recEmit(n.children[0]) + " " +
			ge.recEmit(n.children[1]) + " " +
			ge.recEmit(n.children[2]) + "\n")
	case N_FUNC_DEF:
		for i := range n.children {
			output.WriteString(ge.recEmit(n.children[i]) + " ")
		}
		output.WriteString("\n\n")
	case N_RET_STATE:
		for i := range n.children {
			output.WriteString(ge.recEmit(n.children[i]) + " ")
		}
		output.WriteString("\n")
	case N_BREAK_STATE:
		for i := range n.children {
			output.WriteString(ge.recEmit(n.children[i]) + " ")
		}
		output.WriteString("\n")
	case N_CONT_STATE:
		for i := range n.children {
			output.WriteString(ge.recEmit(n.children[i]) + " ")
		}
		output.WriteString("\n")
	case N_CONDITION:
		for i := range n.children {
			output.WriteString(ge.recEmit(n.children[i]) + " ")
		}
	case N_EXPRESSION:
		for i := range n.children {
			output.WriteString(ge.recEmit(n.children[i]) + " ")
		}
	case N_ASSIGNMENT:
		isFirstShow := false

		if len(n.children) == 2 {
			// No definition, just declaration
			output.WriteString("var " + ge.recEmit(n.children[0]) + " " + ge.recEmit(n.children[1]))

		} else {
			for i := 0; i < len(n.children); i++ {
				if n.children[i].kind == N_COMPLEX_TYPE {
					isFirstShow = true
					ge.varType = n.children[i]
				} else if n.children[i].kind == N_ASSIGN && isFirstShow {
					if ge.inConst {
						output.WriteString("= ")
					} else {
						output.WriteString(":= ")
					}
				} else {
					output.WriteString(ge.recEmit(n.children[i]) + " ")
				}
			}
		}
	case N_LONE_CALL:
		output.WriteString(ge.recEmit(n.children[0]) + " " +
			ge.recEmit(n.children[1]))
	case N_FUNC_CALL:
		for i := range n.children {
			output.WriteString(ge.recEmit(n.children[i]) + " ")
		}
	case N_BLOCK:
		for i := 0; i < len(n.children); i++ {
			output.WriteString(ge.recEmit(n.children[i]))
		}
	case N_NEW_TYPE:
		output.WriteString(ge.recEmit(n.children[0]) + " " +
			ge.recEmit(n.children[1]) + " " +
			ge.recEmit(n.children[2]) + " " +
			ge.recEmit(n.children[3]) + "\n")
	case N_UNARY_OPERATION:
		if n.children[0].kind == N_INC || n.children[0].kind == N_DINC || n.children[0].kind == N_INDEX {
			output.WriteString(ge.recEmit(n.children[1]) +
				ge.recEmit(n.children[0]))
		} else {
			output.WriteString(ge.recEmit(n.children[0]) +
				ge.recEmit(n.children[1]))
		}
	case N_MAKE_ARRAY:
		output.WriteString(ge.recEmit(ge.varType))

		output.WriteString("{")
		for i := 2; i < len(n.children)-1; i++ {
			output.WriteString(ge.recEmit(n.children[i]))
		}
		output.WriteString("}")
	case N_EMPTY_BLOCK:
		output.WriteString("[]")
	case N_COMPLEX_TYPE:
		if n.children[0].kind == N_MAP {
			output.WriteString("map" +
				ge.recEmit(n.children[2]) +
				ge.recEmit(n.children[3]) +
				ge.recEmit(n.children[4]) +
				ge.recEmit(n.children[1]))
		} else {
			for i := len(n.children) - 1; i >= 0; i-- {
				output.WriteString(ge.recEmit(n.children[i]))
			}
		}
	case N_SWITCH_STATE:
		for i := 0; i < len(n.children); i++ {
			output.WriteString(ge.recEmit(n.children[i]) + " ")
		}
	case N_CASE_STATE:
		output.WriteString(ge.recEmit(n.children[0]) + " " +
			ge.recEmit(n.children[1]) +
			":\n" +
			ge.recEmit(n.children[3]))
	case N_DEFAULT_STATE:
		output.WriteString(ge.recEmit(n.children[0]) + " " +
			":\n" +
			ge.recEmit(n.children[2]))
	case N_CASE_BLOCK:
		for i := 0; i < len(n.children); i++ {
			output.WriteString(ge.recEmit(n.children[i]))
		}
	case N_LONE_INC:
		output.WriteString(ge.recEmit(n.children[1]) +
			ge.recEmit(n.children[0]) +
			ge.recEmit(n.children[2]) + "\n")
	case N_METHOD_RECEIVER:
		output.WriteString(ge.recEmit(n.children[0]) +
			ge.recEmit(n.children[1]) + " " +
			ge.recEmit(n.children[2]) +
			ge.recEmit(n.children[3]))
	case N_ENUM_DEF: // TODO: What happens when we get an empty enum?
		typeName := ge.recEmit(n.children[1])
		output.WriteString("type " + typeName + " int\n")

		output.WriteString("\nconst (\n")
		output.WriteString(ge.recEmit(n.children[3]) + " " + typeName + " = iota\n")
		for i := 5; i < len(n.children)-1; i += 2 {
			output.WriteString(ge.recEmit(n.children[i]) + "\n")
		}
		output.WriteString(")\n")
	case N_STRUCT_NEW:
		output.WriteString(ge.recEmit(n.children[1]) + "{")

		for i := 3; i < len(n.children)-1; i++ {
			output.WriteString(ge.recEmit(n.children[i]))
		}

		output.WriteString("}")
	case N_BRACKETED_VALUE:
		output.WriteString(ge.recEmit(n.children[0]) +
			ge.recEmit(n.children[1]) +
			ge.recEmit(n.children[2]))
	case N_ELEMENT_ASSIGNMENT:
		output.WriteString(ge.recEmit(n.children[0]) + " " +
			ge.recEmit(n.children[1]) +
			ge.recEmit(n.children[2]) + "\n")
	case N_STRUCT_DEF:
		// First line
		output.WriteString("type " + ge.recEmit(n.children[1]) + " struct {")
		// Props
		for i := 3; i < len(n.children)-1; i += 2 {
			output.WriteString("\n" + ge.recEmit(n.children[i]) + " " + ge.recEmit(n.children[i+1]))
		}
		// Closing
		output.WriteString("\n}\n\n")
	case N_PROPERTY:
		output.WriteString(ge.recEmit(n.children[0]) +
			ge.recEmit(n.children[1]) +
			ge.recEmit(n.children[2]))
	case N_CONSTANT:
		ge.inConst = true
		output.WriteString(ge.recEmit(n.children[0]) + " " +
			ge.recEmit(n.children[1]))
		ge.inConst = false
	case N_INDEX:
		output.WriteString(ge.recEmit(n.children[0]) +
			ge.recEmit(n.children[1]) +
			ge.recEmit(n.children[2]))
	case N_CONST:
		output.WriteString("const")
	case N_FOR:
		output.WriteString("for")
	case N_RANGE:
		output.WriteString("for")
	case N_FOREVER:
		output.WriteString("for")
	case N_WHILE:
		output.WriteString("for")
	case N_IF:
		output.WriteString("if")
	case N_ELIF:
		output.WriteString("else if")
	case N_ELSE:
		output.WriteString("else")
	case N_CALL:
		output.WriteString("")
	case N_STRUCT:
		output.WriteString("struct")
	case N_FUN:
		output.WriteString("func")
	case N_RET:
		output.WriteString("return")
	case N_BREAK:
		output.WriteString("break")
	case N_CONT:
		output.WriteString("continue")
	case N_ENUM:
		output.WriteString("enum")
	case N_TYPEDEF:
		output.WriteString("type")
	case N_NEW:
		output.WriteString("new")
	case N_MAKE:
		output.WriteString("make")
	case N_MAP:
		output.WriteString("map")
	case N_SWITCH:
		output.WriteString("switch")
	case N_CASE:
		output.WriteString("case")
	case N_DEFAULT:
		output.WriteString("default")
	case N_SEMICOLON:
		output.WriteString(";")
	case N_ASSIGN:
		output.WriteString("=")
	case N_SEP:
		output.WriteString(",")
	case N_COLON:
		output.WriteString(":")
	case N_L_SQUIRLY:
		output.WriteString("{")
	case N_R_SQUIRLY:
		output.WriteString("}")
	case N_L_BLOCK:
		output.WriteString("[")
	case N_R_BLOCK:
		output.WriteString("]")
	case N_L_PAREN:
		output.WriteString("(")
	case N_R_PAREN:
		output.WriteString(")")
	case N_ADD:
		output.WriteString("+")
	case N_SUB:
		output.WriteString("-")
	case N_MUL:
		output.WriteString("*")
	case N_DIV:
		output.WriteString("/")
	case N_OR:
		output.WriteString("|")
	case N_AND:
		output.WriteString("&")
	case N_OROR:
		output.WriteString("||")
	case N_ANDAND:
		output.WriteString("&&")
	case N_EQ:
		output.WriteString("==")
	case N_LT:
		output.WriteString("<")
	case N_GT:
		output.WriteString(">")
	case N_LTEQ:
		output.WriteString("<=")
	case N_GTEQ:
		output.WriteString(">=")
	case N_NEQ:
		output.WriteString("!=")
	case N_MOD:
		output.WriteString("%")
	case N_ACCESS:
		output.WriteString(".")
	case N_XOR:
		output.WriteString("^")
	case N_L_SHIFT:
		output.WriteString("<<")
	case N_R_SHIFT:
		output.WriteString(">>")
	case N_INC:
		output.WriteString("++")
	case N_DINC:
		output.WriteString("--")
	case N_NOT:
		output.WriteString("!")
	case N_REF:
		output.WriteString("&")
	case N_DEREF:
		output.WriteString("*")
	case N_TYPE:
		output.WriteString(n.data)
	case N_IDENTIFIER:
		output.WriteString(n.data)
	case N_INT:
		output.WriteString(n.data)
	case N_FLOAT:
		output.WriteString(n.data)
	case N_STRING:
		output.WriteString(n.data)
	case N_CHAR:
		output.WriteString(n.data)
	case N_BOOL:
		output.WriteString(n.data)
	case N_NIL:
		output.WriteString("nil")

	default:
		panic("Bad node in emitter?")
	}

	return output.String()
}

func (ge *GoEmitter) dump(emitted string) {
	err := os.WriteFile("../output/main.go", []byte(emitted), 0644)
	if err != nil {
		panic(err)
	}
}
