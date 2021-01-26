package main

import "fmt"

func printSlice(lines []string, name string) {
	for i, line := range lines {
		fmt.Printf("%s[%d]: '%s'\n", name, i, line)
	}
}
