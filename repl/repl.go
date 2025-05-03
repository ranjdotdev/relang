package repl

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ranjdotdev/relang/lexer"
	"github.com/ranjdotdev/relang/token"
)

func Tokenize(inPath string, outPath string) {
	file, err := os.Open(filepath.Join("io", inPath))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	outFile, err := os.Create(filepath.Join("io", outPath))
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	scanner := bufio.NewScanner(file)
	l := lexer.New()

	for scanner.Scan() {
		line := scanner.Text()
		l.SetLine(line)

		for {
			tok := l.NextToken()

			if tok.Type == token.EOL {
				break
			}

			tokenString := fmt.Sprintf("\"%s\":\n   type: %s\n   line: %d\n   column: %d\n\n",
				tok.Lexeme, tok.Type, tok.Line, tok.Column)
			_, err := writer.WriteString(tokenString)
			if err != nil {
				panic(err)
			}
		}
	}

	l.End()
	l.ReportErrors()

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(writer, "Error reading file: %v\n", err)
	}
}
