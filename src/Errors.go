package main

import (
	"fmt"
	"strconv"
)

func throwError(job string, caller string, line int, expected string, got any) {
	panic("Error in the " + job + "!\n" +
		"When the " + job + " was trying to decipher: " + caller + "\n" +
		"Error found on line " + strconv.Itoa(line) + "\n" +
		"Expected: " + expected + "\n" +
		"Got: " + fmt.Sprint(got),
	)
}
