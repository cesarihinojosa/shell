package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Command interface {
	Execute(s []string)
}

var commandRegistry = map[string]Command{
	"echo": Echo{},
	"exit": Exit{},
	"type": Type{},
}

type InvalidCommand struct{}

func (i InvalidCommand) Execute(s []string) {
	if len(s) > 0 {
		fmt.Printf("%v: command not found\n", s[0])
	} else {
		fmt.Println()
	}
}

type Echo struct{}

func (e Echo) Execute(s []string) {
	for _, v := range s[1:] {
		fmt.Print(v, " ") // fix this
	}
	fmt.Println()
}

type Exit struct{}

func (e Exit) Execute(s []string) {
	os.Exit(0)
}

type Type struct{}

func (t Type) Execute(s []string) {
	if len(s) < 2 {
		fmt.Println("type requires at least one argument") //standardize this later
		return
	}
	_, ok := commandRegistry[s[1]]
	if !ok {
		if _, ok := os.LookupEnv("PATH"); ok {
			paths := os.Getenv("PATH") // put this into variable??
			for _, v := range strings.Split(paths, ":") {
				path := filepath.Join(v, s[1])
				info, err := os.Stat(path)
				if err == nil && info.Mode().Perm()&0100 != 0 {
					fmt.Printf("%v is %v\n", s[1], path)
					return
				}
			}
		}
		fmt.Printf("%v: not found\n", s[1])
	} else {
		fmt.Printf("%v is a shell builtin\n", s[1])
	}
}

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
			}

			command, ok := commandRegistry[input[0]]
			if !ok {
				command = InvalidCommand{}
			}
			command.Execute(input)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("error!")
		}
	}
}
