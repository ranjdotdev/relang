package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/ranjdotdev/relang/lexer"
	"github.com/ranjdotdev/relang/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer){
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}

func StartOnFile(inPath string, outPath string) {
	data, err := os.ReadFile(fmt.Sprintf("./io/%s", inPath))
	if err != nil {panic(err)}

	outFile, err := os.Create(fmt.Sprintf("./io/%s", outPath))
	if err != nil {panic(err)}
	
	defer outFile.Close()
	
	writer := bufio.NewWriter(outFile)
	defer writer.Flush()
	
	l := lexer.New(string(data))
	
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		tokenString := fmt.Sprintf("\"%s\":\n  type: %s\n", tok.Literal, tok.Type)
		_, err := writer.WriteString(tokenString)
		if err != nil {
			panic(err)
		}
	}
}