package main

import "slices"

const JOB_HOISTER = "Hoister"

type Hoister struct {
	ast *Node
}

// Takes new types (struct, typedef)
// and hoists them first, then it
// hoists function definitions, lastly,
// it returns the rest of the program.
// Even though an ast is a tree, since
// all definitions should be on the
// global scope, we only have to check
// the program's children
func (h *Hoister) hoist() (*Node, *Node, *Node, *Node) {
	types := &Node{kind: N_PROGRAM}
	consts := &Node{kind: N_PROGRAM}
	funcs := &Node{kind: N_PROGRAM}

	for i := 0; i < len(h.ast.children); i++ {
		c := h.ast.children[i]

		if c.kind == N_NEW_TYPE {
			// Add the child to types
			types.children = append(types.children, c)

		} else if c.kind == N_STRUCT_DEF {
			// Add the child to types
			types.children = append(types.children, c)

		} else if c.kind == N_ENUM_DEF {
			// Add the child to types
			types.children = append(types.children, c)

		} else if c.kind == N_FUNC_DEF {
			// Add the child to funcs
			funcs.children = append(funcs.children, c)

		} else if c.kind == N_CONSTANT {
			// Add the child to consts
			consts.children = append(consts.children, c)

		} else {
			continue
		}

		// Remove the child from the ast
		h.ast.children = slices.Delete(h.ast.children, i, i+1)

		// Cancel out the ++ later
		i--
	}

	return types, consts, funcs, h.ast
}
