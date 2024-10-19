package main

type Hoister struct {
	ast *Node
}

// Takes new types (struct, typedef)
// and hoists them first, then it
// hoists function definitions, lastly,
// it returns the rest of the program
func (h *Hoister) hoist() (*Node, *Node, *Node) {

	return nil, nil, nil
}
