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

		if v.kind != V_TYPE && v.kind != V_ENUM {
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
		a.varStack.push(&Var{kind: V_ENUM, data: typeName, ref: n.children[1]})

		// The identifiers based on this type
		for i := 3; i < len(n.children)-1; i += 2 {
			a.varStack.push(&Var{kind: V_VAR, datatype: typeName, data: n.children[i].data, ref: n.children[i]})
		}
	}
}

func (a *Analyser) analyseFunc(n *Node) {

}

func (a *Analyser) analyseNode(n *Node) {

}
