package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("$ ")
		if scanner.Scan() {
			input := strings.Fields(scanner.Text())
			if len(input) == 0 {
				continue
			}
			command, ok := commandRegistry[input[0]]
			if !ok {
				command = ExternalProgram{}
			}
			command.Execute(input) // strategy pattern
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("error!")
		}
	}
}
