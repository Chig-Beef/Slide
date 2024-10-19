package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	args := os.Args[1:]

	// Empty
	if len(args) == 0 {
		fmt.Println("No arguments given, try 'slide help' to view commands'")
		return
	}

	switch args[0] {
	case "build":
		// No file
		if len(args) == 1 {
			fmt.Println("No source files given")

		} else if len(args) == 2 { // Single file

			// Read the file
			data, err := os.ReadFile(args[1])
			if err != nil {
				fmt.Println("Error while trying to read " + args[1])
				panic(err)
			}

			compile(data)

		} else { // Multi file
			var data []byte

			// Read the files
			for i := 1; i < len(args); i++ {
				newData, err := os.ReadFile(args[i])
				if err != nil {
					fmt.Println("Error while trying to read " + args[i])
					panic(err)
				}
				data = append(data, '\n')
				data = append(data, newData...)
			}

			compile(data)
		}

	case "help":
		help(args[1:])
	}
}

func help(args []string) {
	if len(args) == 0 {
		fmt.Println("Possible args are build, or help")
		return
	}
}

func compile(source []byte) {
	fmt.Println("Compilation started")
	fmt.Println()

	t := time.Now()

	lexer := Lexer{line: 1}
	parser := Parser{}
	hoister := Hoister{}
	//analyser := Analyser{}
	//emitter := GoEmitter{}

	lexer.source = source
	lexed := lexer.lex()
	fmt.Println(lexed)
	fmt.Println()

	parser.source = lexed
	parsed := parser.parse()
	fmt.Println(parsed)
	fmt.Println()

	hoister.ast = parsed
	types, funcs, ast := hoister.hoist()
	fmt.Println(types)
	fmt.Println()
	fmt.Println(funcs)
	fmt.Println()
	fmt.Println(ast)
	fmt.Println()

	//analysed := analyser.analyse()
	//emitted := emitter.emit()
	//emitter.dump(emitted)

	fmt.Println()
	fmt.Println("Compilation ended in", time.Now().Sub(t))
}
