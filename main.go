package main

import (
	"github.com/ranjdotdev/relang/repl"
)

func main() {
	repl.StartOnFile("test.re", "output.yml")
}
