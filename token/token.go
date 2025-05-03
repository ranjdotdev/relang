package token

type TokenType string

type Token struct{
	Type TokenType
	Lexeme string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"
	// Identifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT = "INT"
	STRING = "STRING"
	// Operators
	ASSIGN = "="
	PLUS = "+"
	MINUS = "-"
	BANG = "!"
	ASTERISK = "*"
	SLASH = "/"
	MOD = "%"
	LT = "<"
	GT = ">"
	EQ = "=="
	NOT_EQ = "!="
	LEQ = "<="
	GEQ = ">="
    AND      = "&&"
    OR      = "||"
	// Delimiters
	COMMA = ","
	SEMICOLON = ";"
	LPAREN = "("
	RPAREN = ")"
	LCBRACE = "{"
	RCBRACE = "}"
    LSBRACE   = "["
    RSBRACE   = "]"
	// Keywords
	FUNCTION = "FUNCTION"
	VAR = "VAR"
	TRUE = "TRUE"
	FALSE = "FALSE"
	IF = "IF"
	ELSE = "ELSE"
	RETURN = "RETURN"
	LOOP     = "LOOP"
	NIL     = "NIL"
)

var keywords = map[string]TokenType{
	"fn": FUNCTION,
	"var": VAR,
	"true": TRUE,
	"false": FALSE,
	"if": IF,
	"else": ELSE,
	"return": RETURN,
	"for": LOOP,
	"nil": NIL,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}