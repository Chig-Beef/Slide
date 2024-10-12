package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	// Empty
	if len(args) == 0 {
		fmt.Println("No arguments given, try 'slide help' to view commands'")
		return
	}

	switch args[1] {
	case "build":
		// No file
		if len(args) == 1 {
			fmt.Println("No source files given")
			return

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
	lexer := Lexer{}
	//parser := Parser{}
	//analyser := Analyser{}
	//emitter := GoEmitter{}

	lexer.source = source
	lexed := lexer.lex()
	fmt.Println(lexed)
	//parsed := parser.parse(lexed)
	//analysed := analyser.analyse(parsed)
	//emitted := emitter.emit(analysed)
	//emitter.dump(emitted)
}
