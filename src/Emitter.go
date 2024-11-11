package main

import "os"

const JOB_GO_EMITTER = "Go Emitter"

type GoEmitter struct {
	types *Node
	funcs *Node
	ast   *Node
}

func (ge *GoEmitter) emit() string {
	output := ""

	return output
}

func (ge *GoEmitter) dump(emitted string) {
	err := os.WriteFile("../output/main.go", []byte(emitted), 0644)
	if err != nil {
		panic(err)
	}
}
