package main

import (
	"fmt"
	"time"

	"github.com/ranjdotdev/relang/repl"
)

func main() {
	start := time.Now()

	repl.Tokenize("test.re", "output.yml")

	elapsed := time.Since(start)
	fmt.Printf("\nThe tokenization only took: %s\n", elapsed)
}
