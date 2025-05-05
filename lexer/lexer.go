package lexer

import (
	"fmt"

	"github.com/ranjdotdev/relang/token"
)

type Lexer struct {
	line             string
	ch               byte
	position         int
	nextPosition     int
	lineNumber       int
	mode             int
	errors           []string
	multiLineStr     MultiLineStr
	multiLineComment MultiLineComment
	eof              bool
}

type MultiLineStr struct {
	startLine int
	startCol  int
	value     string
}

type MultiLineComment struct {
	startLine int
	startCol  int
}

const (
	normalMode = iota
	multiLineStringMode
	multiLineCommentMode
)

func New() *Lexer {
	l := &Lexer{
		line:         "",
		ch:           0,
		position:     0,
		nextPosition: 0,
		lineNumber:   0,
		mode:         normalMode,
		errors:       []string{},
		multiLineStr: MultiLineStr{},
		eof:          false,
	}
	l.readChar()
	return l
}

func (l *Lexer) End() {
	if l.mode == multiLineStringMode {
		l.errors = append(l.errors, fmt.Sprintf("%d:%d: syntax error: unclosed multi-line string", l.multiLineStr.startLine, l.multiLineStr.startCol))
		l.mode = normalMode
		l.multiLineStr.value = ""
	} else if l.mode == multiLineCommentMode {
		l.errors = append(l.errors, fmt.Sprintf("%d:%d: syntax error: unclosed multi-line comment", l.multiLineComment.startLine, l.multiLineComment.startCol))
		l.mode = normalMode
	}
	l.eof = true
}

func (l *Lexer) Errors() []string {
	return l.errors
}

func (l *Lexer) ReportErrors() {
	for _, err := range l.errors {
		fmt.Println(err)
	}
}

func (l *Lexer) EOF() bool {
	return l.eof
}

func (l *Lexer) readChar() {
	if l.nextPosition >= len(l.line) {
		l.ch = 0
	} else {
		l.ch = l.line[l.nextPosition]
	}
	l.position = l.nextPosition
	l.nextPosition++
}

func (l *Lexer) SetLine(line string) {
	l.line, l.lineNumber, l.position, l.nextPosition = line, l.lineNumber+1, 0, 0
	l.readChar()
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	startCol := l.position

	if l.mode == multiLineStringMode {
		return l.continueMultiLineString()
	} else if l.mode == multiLineCommentMode {
		return l.continueMultiLineComment()
	}

	if !isKnownIllegal(l.ch) {
		l.skipWhitespace()
	}

	switch l.ch {
	case '=':
		if l.peekNextChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{
				Type:   token.EQ,
				Lexeme: string(ch) + string(l.ch),
				Line:   l.lineNumber,
				Column: l.position,
			}
		} else {
			tok = token.Token{
				Type:   token.ASSIGN,
				Lexeme: string(l.ch),
				Line:   l.lineNumber,
				Column: l.position,
			}
		}
	case '&':
		if l.peekNextChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{
				Type:   token.AND,
				Lexeme: string(ch) + string(l.ch),
				Line:   l.lineNumber,
				Column: l.position,
			}
		} else {
			tok = token.Token{
				Type:   token.ILLEGAL,
				Lexeme: string(l.ch),
				Line:   l.lineNumber,
				Column: l.position,
			}
		}
	case '|':
		if l.peekNextChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{
				Type:   token.OR,
				Lexeme: string(ch) + string(l.ch),
				Line:   l.lineNumber,
				Column: l.position,
			}
		} else {
			tok = token.Token{
				Type:   token.ILLEGAL,
				Lexeme: string(l.ch),
				Line:   l.lineNumber,
				Column: l.position,
			}
		}
	case '+':
		tok = token.Token{
			Type:   token.PLUS,
			Lexeme: string(l.ch),
			Line:   l.lineNumber,
			Column: l.position,
		}
	case '-':
		tok = token.Token{
			Type:   token.MINUS,
			Lexeme: string(l.ch),
			Line:   l.lineNumber,
			Column: l.position,
		}
	case '!':
		if l.peekNextChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{
				Type:   token.NOT_EQ,
				Lexeme: string(ch) + string(l.ch),
				Line:   l.lineNumber,
				Column: l.position,
			}
		} else {
			tok = token.Token{
				Type:   token.BANG,
				Lexeme: string(l.ch),
				Line:   l.lineNumber,
				Column: l.position,
			}
		}
	case '/':
		if l.peekNextChar() == '/' {
			l.readChar()
			l.readChar()

			for l.ch != 0 {
				l.readChar()
			}

			return l.NextToken()
		} else if l.peekNextChar() == '*' {
			l.readChar()
			l.readChar()
			l.mode = multiLineCommentMode
			l.multiLineComment.startLine = l.lineNumber
			l.multiLineComment.startCol = l.position - 2

			for {
				if l.ch == 0 {
					return l.emitEOL()
				}

				if l.ch == '*' && l.peekNextChar() == '/' {
					l.readChar()
					l.readChar()
					l.mode = normalMode
					return l.NextToken()
				}

				l.readChar()
			}
		} else {
			tok = token.Token{
				Type:   token.SLASH,
				Lexeme: string(l.ch),
				Line:   l.lineNumber,
				Column: l.position,
			}
		}
	case '%':
		tok = token.Token{
			Type:   token.MOD,
			Lexeme: string(l.ch),
			Line:   l.lineNumber,
			Column: l.position,
		}
	case '*':
		tok = token.Token{
			Type:   token.ASTERISK,
			Lexeme: string(l.ch),
			Line:   l.lineNumber,
			Column: l.position,
		}
	case '<':
		if l.peekNextChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{
				Type:   token.LEQ,
				Lexeme: string(ch) + string(l.ch),
				Line:   l.lineNumber,
				Column: l.position,
			}
		} else {
			tok = token.Token{
				Type:   token.LT,
				Lexeme: string(l.ch),
				Line:   l.lineNumber,
				Column: l.position,
			}
		}
	case '>':
		if l.peekNextChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{
				Type:   token.GEQ,
				Lexeme: string(ch) + string(l.ch),
				Line:   l.lineNumber,
				Column: l.position,
			}
		} else {
			tok = token.Token{
				Type:   token.GT,
				Lexeme: string(l.ch),
				Line:   l.lineNumber,
				Column: l.position,
			}
		}
	case '"':
		str, err := l.readString()
		if err != nil {
			l.errors = append(l.errors, fmt.Sprintf("%d:%d: syntax error: unclosed string", l.lineNumber, startCol))
			tok = token.Token{
				Type:   token.ILLEGAL,
				Lexeme: str,
				Line:   l.lineNumber,
				Column: l.position,
			}
		} else {
			tok = token.Token{
				Type:   token.STRING,
				Lexeme: str,
				Line:   l.lineNumber,
				Column: l.position,
			}
		}
	case '`':
		l.mode = multiLineStringMode
		l.multiLineStr.startCol = l.position
		l.multiLineStr.startLine = l.lineNumber

		l.readChar()

		for l.ch != 0 {
			if l.ch == '`' {
				l.readChar()
				l.mode = normalMode
				return token.Token{
					Type:   token.STRING,
					Lexeme: string(l.ch),
					Line:   l.multiLineStr.startLine,
					Column: l.multiLineStr.startCol,
				}
			}

			l.multiLineStr.value += string(l.ch)
			l.readChar()
		}

		l.multiLineStr.value += "\n"

		return l.emitEOL()
	case ';':
		tok = token.Token{
			Type:   token.SEMICOLON,
			Lexeme: string(l.ch),
			Line:   l.lineNumber,
			Column: l.position,
		}
	case '(':
		tok = token.Token{
			Type:   token.LPAREN,
			Lexeme: string(l.ch),
			Line:   l.lineNumber,
			Column: l.position,
		}
	case ',':
		tok = token.Token{
			Type:   token.COMMA,
			Lexeme: string(l.ch),
			Line:   l.lineNumber,
			Column: l.position,
		}
	case ')':
		tok = token.Token{
			Type:   token.RPAREN,
			Lexeme: string(l.ch),
			Line:   l.lineNumber,
			Column: l.position,
		}
	case '{':
		tok = token.Token{
			Type:   token.LCBRACE,
			Lexeme: string(l.ch),
			Line:   l.lineNumber,
			Column: l.position,
		}
	case '}':
		tok = token.Token{
			Type:   token.RCBRACE,
			Lexeme: string(l.ch),
			Line:   l.lineNumber,
			Column: l.position,
		}
	case '[':
		tok = token.Token{
			Type:   token.LSBRACE,
			Lexeme: string(l.ch),
			Line:   l.lineNumber,
			Column: l.position,
		}
	case ']':
		tok = token.Token{
			Type:   token.RSBRACE,
			Lexeme: string(l.ch),
			Line:   l.lineNumber,
			Column: l.position,
		}
	case 0:
		return l.emitEOL()
	default:
		if isLetter(l.ch) {
			ident := l.readIdentifier()
			identType := token.LookupIdent(ident)
			tok = token.Token{
				Type:   identType,
				Lexeme: ident,
				Line:   l.lineNumber,
				Column: startCol,
			}
			return tok
		} else if isDigit(l.ch) {
			num := l.readNumber()
			tok = token.Token{
				Type:   token.INT,
				Lexeme: num,
				Line:   l.lineNumber,
				Column: startCol,
			}
			return tok
		} else {
			l.errors = append(l.errors, fmt.Sprintf("%d:%d: syntax error: unexpected literal", l.lineNumber, startCol))
			if isKnownIllegal(l.ch) {
				illegal := l.readKnownIllegal()
				tok = token.Token{
					Type:   token.ILLEGAL,
					Lexeme: illegal,
					Line:   l.lineNumber,
					Column: startCol,
				}
			} else {
				tok = token.Token{
					Type:   token.ILLEGAL,
					Lexeme: string(l.ch),
					Line:   l.lineNumber,
					Column: l.position,
				}
			}
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) continueMultiLineString() token.Token {
	for l.ch != 0 {
		if l.ch == '`' {
			l.readChar()
			l.mode = normalMode
			return token.Token{
				Type:   token.STRING,
				Lexeme: l.multiLineStr.value,
				Line:   l.multiLineStr.startLine,
				Column: l.multiLineStr.startCol,
			}
		}

		l.multiLineStr.value += string(l.ch)
		l.readChar()
	}

	l.multiLineStr.value += "\n"

	return token.Token{
		Type:   token.EOL,
		Lexeme: "",
		Line:   l.multiLineComment.startLine,
		Column: l.multiLineStr.startCol,
	}
}

func (l *Lexer) continueMultiLineComment() token.Token {
	for l.ch != 0 {
		if l.ch == '*' && l.peekNextChar() == '/' {
			l.readChar()
			l.readChar()
			l.mode = normalMode
			l.multiLineComment.startLine = 0
			l.multiLineComment.startCol = 0
			return l.NextToken()
		}
		l.readChar()
	}

	return l.emitEOL()
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.line[position:l.position]
}

func (l *Lexer) readKnownIllegal() string {
	position := l.position
	for l.ch != ' ' && l.ch != '\t' && l.ch != '\n' && l.ch != '\r' && l.ch != ';' {
		l.readChar()
	}
	return l.line[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isKnownIllegal(ch byte) bool {
	return ch == '#'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.line[position:l.position]
}

func (l *Lexer) readString() (string, error) {
	startChar := l.ch
	l.readChar()
	position := l.position

	switch startChar {
	case '"':
		for l.ch != '"' && l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
		str := l.line[position:l.position]
		if l.ch == 0 || l.ch == '\n' {
			return str, fmt.Errorf("inline string started with \" and never closed")
		}
		l.readChar()
		return str, nil
	default:
		return "", fmt.Errorf("readString called with invalid character: %c", startChar)
	}
}

func (l *Lexer) peekNextChar() byte {
	if l.nextPosition >= len(l.line) {
		return 0
	} else {
		return l.line[l.nextPosition]
	}
}

func (l *Lexer) emitEOL() token.Token {
	return token.Token{Type: token.EOL, Lexeme: "", Line: l.lineNumber, Column: l.position}
}
