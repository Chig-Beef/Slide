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

	// Other
	a.varStack.push(&Var{kind: V_FUNC, data: "print"})
	a.varStack.push(&Var{kind: V_FUNC, data: "println"})
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
		a.checkMethodReceiver(n.children[1])

		a.varStack.push(&Var{kind: V_FUNC, data: n.children[2].data, ref: n})

		// Parameters
		for i := 4; i < len(n.children)-3; i += 3 {
			typeName, isArray := a.checkValidComplexType(n.children[i+1])
			a.varStack.push(&Var{kind: V_VAR, data: n.children[i].data, datatype: typeName, ref: n.children[i], isArray: isArray})
		}
	} else {
		a.varStack.push(&Var{kind: V_FUNC, data: n.children[1].data, ref: n})

		// Parameters
		for i := 3; i < len(n.children)-3; i += 3 {
			typeName, isArray := a.checkValidComplexType(n.children[i+1])
			a.varStack.push(&Var{kind: V_VAR, data: n.children[i].data, datatype: typeName, ref: n.children[i], isArray: isArray})
		}
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

	a.checkAssignment(n.children[0])
}

func (a *Analyser) checkElementAssignment(n *Node) {
	const FUNC_NAME = "check element assignment"

	a.checkIndexUnary(n.children[0])
	a.checkAssignment(n.children[1])
}

func (a *Analyser) checkIndexUnary(n *Node) {
	const FUNC_NAME = "check index unary"

	a.checkExpression(n.children[1])
}

func (a *Analyser) checkIfBlock(n *Node) {
	const FUNC_NAME = "check if block"

	a.checkExpression(n.children[1])
	a.checkBlock(n.children[2])

	for i := 3; i < len(n.children); i += 3 {
		if n.children[i].kind == N_ELIF {
			a.checkExpression(n.children[i+1])
			a.checkBlock(n.children[i+2])
		}
	}

	i := len(n.children) - 2

	if n.children[i].kind == N_ELSE {
		a.checkBlock(n.children[i+1])
	}
}

func (a *Analyser) checkForeverLoop(n *Node) {
	const FUNC_NAME = "check forever loop"

	a.checkBlock(n.children[len(n.children)-1])
}

func (a *Analyser) checkRangeLoop(n *Node) {
	const FUNC_NAME = "check range loop"

	a.checkExpression(n.children[1])

	a.checkBlock(n.children[len(n.children)-1])
}

func (a *Analyser) checkForLoop(n *Node) {
	const FUNC_NAME = "check for loop"

	i := 1

	if n.children[i].kind != N_SEMICOLON {
		a.checkAssignment(n.children[i])
		i++
	}
	i++

	if n.children[i].kind != N_SEMICOLON {
		a.checkCondition(n.children[i])
		i++
	}
	i++

	if n.children[i].kind != N_BLOCK {
		if n.children[i].kind == N_UNARY_OPERATION {
			a.checkUnaryOperation(n.children[i])
		} else {
			a.checkAssignment(n.children[i])
		}
	}

	a.checkBlock(n.children[len(n.children)-1])
}

func (a *Analyser) checkWhileLoop(n *Node) {
	const FUNC_NAME = "check while loop"

	a.checkCondition(n.children[1])

	a.checkBlock(n.children[len(n.children)-1])
}

// TODO: Check in function
func (a *Analyser) checkRetState(n *Node) {
	const FUNC_NAME = "check ret state"

	if n.children[1].kind != N_SEMICOLON {
		a.checkExpression(n.children[1])
	}
}

// TODO: Check in loop
func (a *Analyser) checkBreakState(n *Node) {
	const FUNC_NAME = "check break state"

	if n.children[1].kind != N_SEMICOLON {
		a.checkValue(n.children[1])
	}
}

// TODO: Check in loop
func (a *Analyser) checkContState(n *Node) {
	const FUNC_NAME = "check cont state"

	if n.children[1].kind != N_SEMICOLON {
		a.checkValue(n.children[1])
	}
}

func (a *Analyser) checkCondition(n *Node) {
	const FUNC_NAME = "check condition"

	a.checkExpression(n)
}

func (a *Analyser) checkExpression(n *Node) {
	const FUNC_NAME = "check expression"

	// TODO: Check operators make sense

	for i := 0; i < len(n.children); i += 2 {
		a.checkValue(n.children[i])
	}
}

func (a *Analyser) checkValue(n *Node) {
	const FUNC_NAME = "check value"

	switch n.kind {
	case N_INT:
	case N_FLOAT:
	case N_STRING:
	case N_CHAR:
	case N_BOOL:
	case N_NIL:
	case N_PROPERTY:
		panic("not implemented")

	case N_STRUCT_NEW:
		a.checkStructNew(n)
	case N_IDENTIFIER:
		a.checkValidIdentifier(FUNC_NAME, n)
	case N_UNARY_OPERATION:
		a.checkUnaryOperation(n)
	case N_MAKE_ARRAY:
		a.checkMakeArray(n)
	case N_L_PAREN:
		a.checkBracketedValue(n)
	case N_CALL:
		a.checkFuncCall(n)
	case N_NEW:
		a.checkStructNew(n)
	default:
		throwError(JOB_ANALYSER, FUNC_NAME, n.line, "value", n.kind)
	}
}

func (a *Analyser) checkAssignment(n *Node) {
	const FUNC_NAME = "check assignment"

	if n.children[1].kind == N_COMPLEX_TYPE {
		a.checkValidComplexType(n.children[1])
	}

	if len(n.children) > 2 {
		a.checkExpression(n.children[len(n.children)-1])
	}

	// Add the variable to the stack
	a.varStack.push(&Var{kind: V_VAR, data: n.children[0].data, ref: n.children[0]})
}

func (a *Analyser) checkLoneCall(n *Node) {
	const FUNC_NAME = "check lone call"

	a.checkFuncCall(n.children[0])
}

func (a *Analyser) checkFuncCall(n *Node) {
	const FUNC_NAME = "check func call"

}

func (a *Analyser) checkStructNew(n *Node) {
	const FUNC_NAME = "check struct new"

	v := a.checkForMatch(n.children[1].data)
	if v == nil {
		throwError(JOB_ANALYSER, FUNC_NAME, n.line, "existing type", n)
	}

	if v.kind != V_TYPE {
		throwError(JOB_ANALYSER, FUNC_NAME, n.line, "valid type", n)
	}

	for i := 3; i < len(n.children)-1; i += 2 {
		a.checkExpression(n.children[i])
	}
}

func (a *Analyser) checkBlock(n *Node) {
	const FUNC_NAME = "check block"

	endSize := a.varStack.length

	for i := 1; i < len(n.children)-1; i++ {
		a.analyseNode(n.children[i])
	}

	for a.varStack.length > endSize {
		a.varStack.pop()
	}
}

func (a *Analyser) checkUnaryOperation(n *Node) {
	const FUNC_NAME = "check unary operation"

}

func (a *Analyser) checkBracketedValue(n *Node) {
	const FUNC_NAME = "check bracketed value"

	a.checkExpression(n.children[1])
}

func (a *Analyser) checkMakeArray(n *Node) {
	const FUNC_NAME = "check make array"

}

func (a *Analyser) checkSwitchState(n *Node) {
	const FUNC_NAME = "check switch state"

	a.checkExpression(n.children[1])

	// TODO: Don't do default check on every
	// single one for performance?
	for i := 3; i < len(n.children)-1; i++ {
		if n.children[i].kind == N_DEFAULT_STATE {
			a.checkDefaultState(n.children[i])
		} else {
			a.checkCaseState(n.children[i])
		}
	}
}

// TODO: More checks on expression
func (a *Analyser) checkCaseState(n *Node) {
	const FUNC_NAME = "check case state"

	a.checkExpression(n.children[1])

	a.checkCaseBlock(n.children[3])
}

func (a *Analyser) checkDefaultState(n *Node) {
	const FUNC_NAME = "check default state"
	a.checkCaseBlock(n.children[2])
}

func (a *Analyser) checkCaseBlock(n *Node) {
	const FUNC_NAME = "check case block"

	endSize := a.varStack.length

	for i := range n.children {
		a.analyseNode(n.children[i])
	}

	for a.varStack.length > endSize {
		a.varStack.pop()
	}
}

func (a *Analyser) checkLoneInc(n *Node) {
	const FUNC_NAME = "check lone inc"

	a.checkValidIdentifier(FUNC_NAME, n.children[1])
}

func (a *Analyser) checkMethodReceiver(n *Node) {
	const FUNC_NAME = "check method receiver"

	typeName, isArray := a.checkValidComplexType(n.children[2])

	// Add the variable to the stack. This
	// will be popped by the function
	// definition
	a.varStack.push(&Var{kind: V_VAR, data: n.children[1].data, datatype: typeName, isArray: isArray})
}

// TODO: Flesh this out into doing actual complex types
// TODO: Map
func (a *Analyser) checkValidComplexType(n *Node) (string, bool) {
	const FUNC_NAME = "check valid complex type"

	if n.kind != N_COMPLEX_TYPE {
		throwError(JOB_ANALYSER, FUNC_NAME, n.line, "complex type", n.kind)
		return "", false
	}

	typeName := ""

	if n.children[0].kind == N_MAP {

		return "map", false
	}

	v := a.checkForMatch(n.children[0].data)
	if v == nil {
		throwError(JOB_ANALYSER, FUNC_NAME, n.line, "valid type", n.children[0])
		return "", false
	}

	if v.kind != V_TYPE {
		throwError(JOB_ANALYSER, FUNC_NAME, n.line, "valid type", n.children[0])
		return "", false
	}

	typeName = n.children[0].data

	// For isArray check
	l := n.children[len(n.children)-1].kind

	return typeName, l == N_EMPTY_BLOCK || l == N_INDEX
}

func (a *Analyser) checkValidIdentifier(FUNC_NAME string, n *Node) {
	v := a.checkForMatch(n.data)

	if v == nil {
		throwError(JOB_ANALYSER, FUNC_NAME, n.line, "existing variable", n)
	}

	if v.kind != V_VAR {
		throwError(JOB_ANALYSER, FUNC_NAME, n.line, "variable", n)
	}
}

func (a *Analyser) checkValidFunction(FUNC_NAME string, n *Node) {
	v := a.checkForMatch(n.data)

	if v == nil {
		throwError(JOB_ANALYSER, FUNC_NAME, n.line, "existing variable", n)
	}

	if v.kind != V_VAR {
		throwError(JOB_ANALYSER, FUNC_NAME, n.line, "variable", n)
	}
}
