package main

import (
	// "fmt"
	// "os"
	// "os/user"

	"github.com/ranjdotdev/relang/repl"
)

func main() {
	// user, err := user.Current()
	// if err != nil { panic(err) }
	// fmt.Printf("Hello %s! This is Relang.\n", user.Username)
	// fmt.Printf("Feel free to type in commands\n")
	// repl.Start(os.Stdin, os.Stdout)
	repl.StartOnFile("test.re", "output.yml")
}