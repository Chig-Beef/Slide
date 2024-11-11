package main

type Var struct {
	kind     VarType // Overarching type (identifier, func, etc)
	data     string  // Name of it
	datatype string  // The actual type of the variable
	ref      *Node   // Where is this from?
	props    []*Var  // For structs, their properties
	key      *Var    // For maps, the type for the key
	value    *Var    // For maps, the type for the value
	isArray  bool    // Will have length property
}

type VarType int

const (
	V_VAR VarType = iota
	V_TYPE
	V_FUNC
	V_MAP
)

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
