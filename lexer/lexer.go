package lexer

import (
	"fmt"

	"github.com/ranjdotdev/relang/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	if !isKnownIllegal(l.ch){
		l.skipWhitespace()
	}
	
	
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			ch := l.ch
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		if l.peekChar() == '/' {
			l.readChar() 
			l.readChar()
			
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			
			return l.NextToken()
		} else if l.peekChar() == '*' {
			l.readChar()
			l.readChar()
			
			for {
				if l.ch == 0 {
					tok = newToken(token.ILLEGAL, l.ch)
					return tok
				}
				
				if l.ch == '*' && l.peekChar() == '/' {
					l.readChar()
					l.readChar()
					
					return l.NextToken()
				}
				
				l.readChar()
			}
		} else {
			tok = newToken(token.SLASH, l.ch)
		}
	case '%':
		tok = newToken(token.MOD, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LEQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GEQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case '"':
		str, err := l.readString()
		if err != nil {
			tok.Literal = str
			tok.Type = token.ILLEGAL
		} else {
			tok.Literal = str
			tok.Type = token.STRING
		}
	case '`':
		str, err := l.readString()
		if err != nil {
			tok.Literal = str
			tok.Type = token.ILLEGAL
		} else {
			tok.Literal = str
			tok.Type = token.STRING
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LCBRACE, l.ch)
	case '}':
		tok = newToken(token.RCBRACE, l.ch)
	case '[':
		tok = newToken(token.LSBRACE, l.ch)
	case ']':
		tok = newToken(token.RSBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			if isKnownIllegal(l.ch){
				tok.Literal = l.readKnownIllegal()
				tok.Type = token.ILLEGAL
			} else {
				tok = newToken(token.ILLEGAL, l.ch)
			}
		}
	}
	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readKnownIllegal() string {
    position := l.position
	for l.ch != ' ' && l.ch != '\t' && l.ch != '\n' && l.ch != '\r' && l.ch != ';' {
		l.readChar()
	}
    return l.input[position:l.position]
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
	return l.input[position:l.position]
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
		    str := l.input[position:l.position]
        if l.ch == 0 || l.ch == '\n' {
            return str, fmt.Errorf("inline string started with \" and never closed")
        }
        l.readChar()
        return str, nil
        
    case '`':
        for l.ch != '`' && l.ch != 0 {
            l.readChar()
        }
    		str := l.input[position:l.position]
        if l.ch == 0 {
            return str, fmt.Errorf("multi-line string started with ` and never closed")
        }
        l.readChar()
        return str, nil
    default:
        return "", fmt.Errorf("readString called with invalid character: %c", startChar)
    }
}

func (l *Lexer) peekChar() byte {
  if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}
