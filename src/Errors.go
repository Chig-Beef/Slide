package main

import "strings"

var errorSource []byte

func throwError(line int, msg string) {
	panic(msg + "\n\nLine Below\n" + getSourceLine(line))
}

func getSourceLine(line int) string {
	output := ""

	inLine := false
	curLine := 0

	for i := 0; i < len(errorSource); i++ {
		if inLine {
			if errorSource[i] == '\n' {
				break
			}

			output += string(errorSource[i])
		} else {
			if errorSource[i] == '\n' {
				curLine++
				if curLine == line-1 {
					inLine = true
				}
			}
		}
	}

	return strings.Trim(output, " \t\r\n")
}
