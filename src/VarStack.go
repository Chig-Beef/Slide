package main

type VarFrame struct {
	value *Node
	prev  *VarFrame
}

type VarStack struct {
	tail   *VarFrame
	length int
}

func (vs *VarStack) push(n *Node) {
	vf := &VarFrame{n, nil}

	vs.length++

	if vs.length == 1 {
		vs.tail = vf
		return
	}

	vf.prev = vs.tail
	vs.tail = vf
}

func (vs *VarStack) pop() *Node {
	if vs.length == 0 {
		return nil
	}

	vs.length--

	if vs.length == 0 {
		tail := vs.tail
		vs.tail = nil
		return tail.value
	}

	tail := vs.tail
	vs.tail = vs.tail.prev
	return tail.value
}

func (vs *VarStack) peek() *Node {
	// Length == 0
	if vs.tail == nil {
		return nil
	}

	return vs.tail.value
}
