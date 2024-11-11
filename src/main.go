package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	args := os.Args

	// Empty
	if len(args) == 1 {
		fmt.Println("No arguments given, try 'slide help' to view commands'")
		return
	}

	cmd := args[1]
	args = args[2:]

	switch cmd {
	case "build":
		// No file
		if len(args) == 0 {
			fmt.Println("No source files given")

			// TODO: What if they gave us a directory? Deal with as such
		} else if len(args) == 1 { // Single file

			// Read the file
			data, err := os.ReadFile(args[0])
			if err != nil {
				fmt.Println("Error while trying to read " + args[0])
				panic(err)
			}

			compile(data)

		} else { // Multi file
			var data []byte

			// Read the files
			for i := 0; i < len(args); i++ {
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
		help(args)
	}
}

func help(args []string) {
	if len(args) == 0 {
		fmt.Println("Possible args are build, or help")
		return
	}

	if len(args) > 1 {
		fmt.Println("Only one command at a time")
		return
	}

	cmd := args[0]

	switch cmd {
	case "help":
		fmt.Println("\"slide help\" takes in any slide command and gives a hint to how it works")
	case "build":
		fmt.Println("\"slide build\" take in a file name and compiles the contents of the file")
		fmt.Println("\"slide build\" can also take in multiple files at once, each as another argument")
	}
}

func compile(source []byte) {
	fmt.Println("Compilation started")
	fmt.Println()

	t := time.Now()

	lexer := Lexer{line: 1}
	parser := Parser{}
	hoister := Hoister{}
	analyser := Analyser{}
	emitter := GoEmitter{}

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
	fmt.Println("types: ", types)
	fmt.Println()
	fmt.Println("funcs: ", funcs)
	fmt.Println()
	fmt.Println("ast: ", ast)
	fmt.Println()

	analyser.types = types
	analyser.funcs = funcs
	analyser.ast = ast
	analyser.analyse()

	emitted := emitter.emit()
	emitter.dump(emitted)

	fmt.Println()
	fmt.Println("Compilation ended in", time.Now().Sub(t))
}
