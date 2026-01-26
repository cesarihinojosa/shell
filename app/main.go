package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	// TODO: Uncomment the code below to pass the first stage

	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("$ ")
		if scanner.Scan() {
			input := strings.Fields(scanner.Text())
			if len(input) == 0 {
				continue
			} else if input[0] == "exit" {
				os.Exit(0)
			} else if input[0] == "echo" {
				for _, v := range input[1:] {
					fmt.Print(v, " ")
				}
				fmt.Println()
			} else {
				fmt.Printf("%v: command not found\n", input[0])
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("error!")
		}
	}
}
