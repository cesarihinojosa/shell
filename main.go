package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// TODO add standardized logging

type Command interface {
	Execute(s []string)
}

var commandRegistry = map[string]Command{
	"echo": Echo{},
	"exit": Exit{},
	"type": Type{},
	"pwd":  Pwd{},
	"cd":   Cd{},
}

type ExternalProgram struct{}

func (e ExternalProgram) Execute(s []string) {
	if len(s) <= 0 {
		return
	}
	if _, ok := os.LookupEnv("PATH"); ok { // TODO repeating code with type. fix this
		paths := os.Getenv("PATH")
		for _, v := range strings.Split(paths, ":") {
			path := filepath.Join(v, s[0])
			info, err := os.Stat(path)
			if err == nil && info.Mode().Perm()&0100 != 0 {
				var cmd *exec.Cmd
				if len(s) == 1 { // TODO do I need this?
					cmd = exec.Command(s[0])
				} else {
					cmd = exec.Command(s[0], s[1:]...)
				}
				cmd.Stdout = os.Stdout // TODO understansd these better
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					fmt.Println(err)
				}
				return
			}
		}
	}
	fmt.Printf("%v: command not found\n", s[0])
}

type Echo struct{}

func (e Echo) Execute(s []string) {
	for _, v := range s[1:] {
		fmt.Print(v, " ") // TODO fix this
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

type Pwd struct{}

func (p Pwd) Execute(s []string) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error()) // TODO fix me later
	}
	fmt.Println(wd)
}

type Cd struct{}

// TODO add "-" support and handle no args correctly
func (c Cd) Execute(s []string) {
	if len(s) < 2 {
		fmt.Println("cd requires one argument")
		return
	}
	if s[1] == "~" {
		s[1] = os.Getenv("HOME")
	}
	err := os.Chdir(s[1])
	if err != nil {
		fmt.Printf("cd: %v: No such file or directory\n", s[1])
	}
}

// func singleQuoteFilter(s []string) []string {

// }

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
				command = ExternalProgram{}
			}
			command.Execute(input) // strategy pattern
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("error!")
		}
	}
}
