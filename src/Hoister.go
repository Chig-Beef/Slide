package main

import "slices"

type Hoister struct {
	ast *Node
}

// Takes new types (struct, typedef)
// and hoists them first, then it
// hoists function definitions, lastly,
// it returns the rest of the program
func (h *Hoister) hoist() (*Node, *Node, *Node) {
	types := &Node{kind: N_PROGRAM}
	funcs := &Node{kind: N_PROGRAM}

	h.walkAndHoist(h.ast, nil, types, funcs)

	return types, funcs, h.ast
}

// Will clean node of type defs and
// funcs, and add these to the relevant
// nodes. Returns true if it was
// hoisted
func (h *Hoister) walkAndHoist(node *Node, parent *Node, types *Node, funcs *Node) bool {
	// BASE CASES

	// Can't recurse
	if len(node.children) == 0 {
		return false
	}

	// Found a typedef
	if node.kind == N_NEW_TYPE {
		// Add the child to types
		types.children = append(types.children, node)

		// Remove the child from the ast
		i := slices.Index(parent.children, node)
		if i == -1 {
			panic("Hoister couldn't find child in parent (doesn't make sense, does it?")
		}
		parent.children = slices.Delete(parent.children, i, i+1)

		return true
	}

	// Found a struct def
	if node.kind == N_STRUCT_DEF {
		// Add the child to types
		types.children = append(types.children, node)

		// Remove the child from the ast
		i := slices.Index(parent.children, node)
		if i == -1 {
			panic("Hoister couldn't find child in parent (doesn't make sense, does it?")
		}
		parent.children = slices.Delete(parent.children, i, i+1)

		return true
	}

	// Founc an enum def
	if node.kind == N_ENUM_DEF {
		// Add the child to types
		types.children = append(types.children, node)

		// Remove the child from the ast
		i := slices.Index(parent.children, node)
		if i == -1 {
			panic("Hoister couldn't find child in parent (doesn't make sense, does it?")
		}
		parent.children = slices.Delete(parent.children, i, i+1)

		return true
	}

	// Found a function def
	if node.kind == N_FUNC_DEF {
		// Add the child to types
		funcs.children = append(funcs.children, node)

		// Remove the child from the ast
		i := slices.Index(parent.children, node)
		if i == -1 {
			panic("Hoister couldn't find child in parent (doesn't make sense, does it?")
		}
		parent.children = slices.Delete(parent.children, i, i+1)

		return true
	}

	// RECURSION CASE
	for i := 0; i < len(node.children); i++ {
		if h.walkAndHoist(node.children[i], node, types, funcs) {
			i--
		}
	}

	return false
}
