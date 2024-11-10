package main

const JOB_ANALYSER = "Analyser"

type Analyser struct {
	types *Node
	funcs *Node
	ast   *Node

	// The stack that all variables are on,
	// all types and functions are added
	// first, and then each variable is
	// pushed and popped as needed
	varStack *VarStack
}

func (a *Analyser) analyse() {
	a.varStack = &VarStack{}

	a.initVarStack()

	for i := 0; i < len(a.types.children); i++ {
		a.analyseType(a.types.children[i])
	}

	for i := 0; i < len(a.funcs.children); i++ {
		a.analyseFunc(a.funcs.children[i])
	}

	a.analyseNode(a.ast)
}

func (a *Analyser) initVarStack() {
	// Basic types are added so that we can
	// use them
	a.varStack.push(&Var{kind: V_TYPE, data: "byte"})
	a.varStack.push(&Var{kind: V_TYPE, data: "word"})
	a.varStack.push(&Var{kind: V_TYPE, data: "dword"})
	a.varStack.push(&Var{kind: V_TYPE, data: "qword"})
	a.varStack.push(&Var{kind: V_TYPE, data: "uint8"})
	a.varStack.push(&Var{kind: V_TYPE, data: "uint16"})
	a.varStack.push(&Var{kind: V_TYPE, data: "uint32"})
	a.varStack.push(&Var{kind: V_TYPE, data: "uint64"})
	a.varStack.push(&Var{kind: V_TYPE, data: "uint"})
	a.varStack.push(&Var{kind: V_TYPE, data: "int8"})
	a.varStack.push(&Var{kind: V_TYPE, data: "int16"})
	a.varStack.push(&Var{kind: V_TYPE, data: "int32"})
	a.varStack.push(&Var{kind: V_TYPE, data: "int64"})
	a.varStack.push(&Var{kind: V_TYPE, data: "sint"})
	a.varStack.push(&Var{kind: V_TYPE, data: "int"})
	a.varStack.push(&Var{kind: V_TYPE, data: "char"})
	a.varStack.push(&Var{kind: V_TYPE, data: "string"})
	a.varStack.push(&Var{kind: V_TYPE, data: "float32"})
	a.varStack.push(&Var{kind: V_TYPE, data: "float64"})
	a.varStack.push(&Var{kind: V_TYPE, data: "double"})
	a.varStack.push(&Var{kind: V_TYPE, data: "float"})
	a.varStack.push(&Var{kind: V_TYPE, data: "bool"})
	a.varStack.push(&Var{kind: V_TYPE, data: "any"})

	// Conversion functions
	a.varStack.push(&Var{kind: V_FUNC, data: "tobyte"})
	a.varStack.push(&Var{kind: V_FUNC, data: "toword"})
	a.varStack.push(&Var{kind: V_FUNC, data: "todword"})
	a.varStack.push(&Var{kind: V_FUNC, data: "toqword"})
	a.varStack.push(&Var{kind: V_FUNC, data: "touint8"})
	a.varStack.push(&Var{kind: V_FUNC, data: "touint16"})
	a.varStack.push(&Var{kind: V_FUNC, data: "touint32"})
	a.varStack.push(&Var{kind: V_FUNC, data: "touint64"})
	a.varStack.push(&Var{kind: V_FUNC, data: "touint"})
	a.varStack.push(&Var{kind: V_FUNC, data: "toint8"})
	a.varStack.push(&Var{kind: V_FUNC, data: "toint16"})
	a.varStack.push(&Var{kind: V_FUNC, data: "toint32"})
	a.varStack.push(&Var{kind: V_FUNC, data: "toint64"})
	a.varStack.push(&Var{kind: V_FUNC, data: "tosint"})
	a.varStack.push(&Var{kind: V_FUNC, data: "toint"})
	a.varStack.push(&Var{kind: V_FUNC, data: "tochar"})
	a.varStack.push(&Var{kind: V_FUNC, data: "tostring"})
	a.varStack.push(&Var{kind: V_FUNC, data: "tofloat32"})
	a.varStack.push(&Var{kind: V_FUNC, data: "tofloat64"})
	a.varStack.push(&Var{kind: V_FUNC, data: "todouble"})
	a.varStack.push(&Var{kind: V_FUNC, data: "tofloat"})
	a.varStack.push(&Var{kind: V_FUNC, data: "tobool"})
}

// Searches the stack until a match for
// an identifier is found
func (a *Analyser) checkForMatch(query string) *Var {
	curr := a.varStack.tail

	for i := 0; i < a.varStack.length; i++ {
		// Match?
		if curr.value.data == query {
			return curr.value
		}

		// Shouldn't really happen, but still
		// important I guess
		if curr.prev == nil {
			return nil
		}

		curr = curr.prev
	}

	// Not found
	return nil
}

func (a *Analyser) analyseType(n *Node) {
	const FUNC_NAME = "analyse type"

	var typeName string
	var aliasName string

	switch n.kind {
	case N_NEW_TYPE:
		typeName = n.children[1].data
		aliasName = n.children[2].data

		// Make sure the identifier doesn't
		// already exist
		v := a.checkForMatch(typeName)
		if v != nil {
			throwError(JOB_ANALYSER, FUNC_NAME, n.line, "the identifier to not exist", v.data+" exists")
		}

		// Make sure the type to alias exists
		v = a.checkForMatch(aliasName)
		if v == nil {
			throwError(JOB_ANALYSER, FUNC_NAME, n.line, "the specified type to exist", aliasName+" doesn't exist")
		}

		if v.kind != V_TYPE {
			throwError(JOB_ANALYSER, FUNC_NAME, n.line, "the specified type isn't a type", aliasName+" is actually a "+v.kind.String())
		}

		a.varStack.push(&Var{kind: V_TYPE, datatype: aliasName, data: typeName, ref: n.children[1]})

	case N_STRUCT_DEF:
		typeName = n.children[1].data
		a.varStack.push(&Var{kind: V_TYPE, data: typeName, ref: n.children[1]})

		// We need to modify the props of the
		// struct first
		s := a.varStack.tail.value

		// Attributes
		for i := 3; i < len(n.children)-1; i += 2 {
			s.props = append(s.props, &Var{kind: V_VAR, datatype: n.children[i+1].data, data: n.children[i].data, ref: n.children[i]})
		}

	case N_ENUM_DEF:
		typeName = n.children[1].data

		// The enum type
		a.varStack.push(&Var{kind: V_TYPE, data: typeName, ref: n.children[1]})

		// The identifiers based on this type
		for i := 3; i < len(n.children)-1; i += 2 {
			a.varStack.push(&Var{kind: V_VAR, datatype: typeName, data: n.children[i].data, ref: n.children[i]})
		}
	default:
		throwError(JOB_ANALYSER, FUNC_NAME, n.line, "Invalid statement went through type check", n.kind)
	}
}

func (a *Analyser) analyseFunc(n *Node) {
	const FUNC_NAME = "analyse func"

	if n.kind != N_FUNC_DEF {
		throwError(JOB_ANALYSER, FUNC_NAME, n.line, "Invalid statement went through func check", n.kind)
		return
	}

	// At the end of this function, the
	// stack should be one larger than it's
	// original size, 1 being the function
	endSize := a.varStack.length + 1

	if n.children[1].kind == N_METHOD_RECEIVER {

		return
	}

	a.varStack.push(&Var{kind: V_FUNC, data: n.children[1].data, ref: n})

	// Parameters
	for i := 3; i < len(n.children)-3; i += 3 {
		typeName := a.checkValidComplexType(n.children[i+1])
		a.varStack.push(&Var{kind: V_VAR, data: n.children[i].data, datatype: typeName, ref: n.children[i]})
	}

	// The function body
	a.analyseNode(n.children[len(n.children)-1])

	// Unrolling everything
	for a.varStack.length > endSize {
		a.varStack.pop()
	}
}

func (a *Analyser) analyseNode(n *Node) {
	const FUNC_NAME = "analyse node"

	switch n.kind {
	case N_ILLEGAL:
		throwError(JOB_ANALYSER, FUNC_NAME, n.line, "illegal found!", n.kind)
	case N_PROGRAM:
		for i := range n.children {
			a.analyseNode(n.children[i])
		}
	case N_VAR_DECLARATION:
		a.checkVarDeclaration(n)
	case N_ELEMENT_ASSIGNMENT:
		a.checkElementAssignment(n)
	case N_IF_BLOCK:
		a.checkIfBlock(n)
	case N_FOREVER_LOOP:
		a.checkForeverLoop(n)
	case N_RANGE_LOOP:
		a.checkRangeLoop(n)
	case N_FOR_LOOP:
		a.checkForLoop(n)
	case N_WHILE_LOOP:
		a.checkWhileLoop(n)
	case N_RET_STATE:
		a.checkRetState(n)
	case N_BREAK_STATE:
		a.checkBreakState(n)
	case N_CONT_STATE:
		a.checkContState(n)
	case N_LONE_CALL:
		a.checkLoneCall(n)
	case N_SWITCH_STATE:
		a.checkSwitchState(n)
	case N_LONE_INC:
		a.checkLoneInc(n)
	case N_BLOCK:
		a.checkBlock(n)

	default:
		throwError(JOB_ANALYSER, FUNC_NAME, n.line, "any other start to node", n.kind)
	}
}

func (a *Analyser) checkVarDeclaration(n *Node) {
	const FUNC_NAME = "check var declaration"
}

func (a *Analyser) checkElementAssignment(n *Node) {
	const FUNC_NAME = "check element assignment"

}

func (a *Analyser) checkIfBlock(n *Node) {
	const FUNC_NAME = "check if block"

}

func (a *Analyser) checkForeverLoop(n *Node) {
	const FUNC_NAME = "check forever loop"

}

func (a *Analyser) checkRangeLoop(n *Node) {
	const FUNC_NAME = "check range loop"

}

func (a *Analyser) checkForLoop(n *Node) {
	const FUNC_NAME = "check for loop"

}

func (a *Analyser) checkWhileLoop(n *Node) {
	const FUNC_NAME = "check while loop"

}

func (a *Analyser) checkRetState(n *Node) {
	const FUNC_NAME = "check ret state"

}

func (a *Analyser) checkBreakState(n *Node) {
	const FUNC_NAME = "check break state"

}

func (a *Analyser) checkContState(n *Node) {
	const FUNC_NAME = "check cont state"

}

func (a *Analyser) checkCondition(n *Node) {
	const FUNC_NAME = "check condition"

}

func (a *Analyser) checkExpression(n *Node) {
	const FUNC_NAME = "check expression"

}

func (a *Analyser) checkAssignment(n *Node) {
	const FUNC_NAME = "check assignment"

}

func (a *Analyser) checkLoneCall(n *Node) {
	const FUNC_NAME = "check lone call"

}

func (a *Analyser) checkFuncCall(n *Node) {
	const FUNC_NAME = "check func call"

}

func (a *Analyser) checkStructNew(n *Node) {
	const FUNC_NAME = "check struct new"

}

func (a *Analyser) checkBlock(n *Node) {
	const FUNC_NAME = "check block"

}

func (a *Analyser) checkUnaryOperation(n *Node) {
	const FUNC_NAME = "check unary operation"

}

func (a *Analyser) checkBracketedValue(n *Node) {
	const FUNC_NAME = "check bracketed value"

}

func (a *Analyser) checkMakeArray(n *Node) {
	const FUNC_NAME = "check make array"

}

func (a *Analyser) checkSwitchState(n *Node) {
	const FUNC_NAME = "check switch state"

}

func (a *Analyser) checkCaseState(n *Node) {
	const FUNC_NAME = "check case state"

}

func (a *Analyser) checkDefaultState(n *Node) {
	const FUNC_NAME = "check default state"

}

func (a *Analyser) checkCaseBlock(n *Node) {
	const FUNC_NAME = "check case block"

}

func (a *Analyser) checkLoneInc(n *Node) {
	const FUNC_NAME = "check lone inc"

	a.checkValidIdentifier(FUNC_NAME, n.children[1], n.children[1].data)
}

func (a *Analyser) checkMethodReceiver(n *Node) {
	const FUNC_NAME = "check method receiver"

}

// TODO: Flesh this out into doing actual complex types
func (a *Analyser) checkValidComplexType(n *Node) string {
	const FUNC_NAME = "check valid complex type"

	if n.kind != N_COMPLEX_TYPE {
		throwError(JOB_ANALYSER, FUNC_NAME, n.line, "complex type", n.kind)
		return ""
	}

	typeName := ""

	v := a.checkForMatch(n.children[0].data)
	if v == nil {
		throwError(JOB_ANALYSER, FUNC_NAME, n.line, "valid type", n.children[0].data)
		return ""
	}

	typeName = n.children[0].data

	return typeName
}

func (a *Analyser) checkValidIdentifier(FUNC_NAME string, n *Node, identifier string) {
	v := a.checkForMatch(identifier)

	if v == nil {
		throwError(JOB_ANALYSER, FUNC_NAME, n.line, "variable", n.kind)
	}

	if v.kind != V_VAR {
		throwError(JOB_ANALYSER, FUNC_NAME, n.line, "variable", n.kind)
	}
}
