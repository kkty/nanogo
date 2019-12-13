package nanogo

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

type lexer struct {
	program string
	result  *Program
}

func atoi(s string) int64 {
	i, err := strconv.Atoi(s)

	if err != nil {
		log.Fatal(err)
	}

	return int64(i)
}

func atof(s string) float32 {
	f, err := strconv.ParseFloat(s, 32)

	if err != nil {
		log.Fatal(err)
	}

	return float32(f)
}

func (l *lexer) Lex(lval *yySymType) int {
	advance := func(i int) {
		l.program = l.program[i:]
	}

	hasPrefix := func(s string) bool {
		return strings.HasPrefix(l.program, s)
	}

	skipWhitespaces := func() bool {
		modified := false
		for hasPrefix(" ") || hasPrefix("\n") || hasPrefix("\t") {
			modified = true
			advance(1)
		}
		return modified
	}

	for skipWhitespaces() {
	}

	if len(l.program) == 0 {
		// 0 stands for EOF.
		return 0
	}

	patterns := []struct {
		pattern string
		token   int
		f       func(s string)
	}{
		{"func", FUNC, nil},
		{"{", LEFT_BRACE, nil},
		{"}", RIGHT_BRACE, nil},
		{"\\(", LEFT_PARENTHESIS, nil},
		{"\\)", RIGHT_PARENTHESIS, nil},
		{",", COMMA, nil},
		{"int64", INT64, nil},
		{"float64", FLOAT64, nil},
		{"bool", BOOL, nil},
		{"true", BOOL_VALUE, func(s string) { lval.val = true }},
		{"false", BOOL_VALUE, func(s string) { lval.val = false }},
		{"[0-9]+", INT_VALUE, func(s string) { lval.val = atoi(s) }},
		{"[0-9]+(\\.[0-9]*)?([eE][\\+\\-]?[0-9]+)?", FLOAT_VALUE, func(s string) { lval.val = atof(s) }},
		{"-", MINUS, nil},
		{"\\+", PLUS, nil},
		{"\\*", ASTERISK, nil},
		{"/", SLASH, nil},
		{"=", EQUAL, nil},
		{"==", EQUAL_EQUAL, nil},
		{"!=", EXCLAMATION_EQUAL, nil},
		{"<", LESS, nil},
		{"<=", LESS_EQUAL, nil},
		{">", GREATER, nil},
		{">=", GREATER_EQUAL, nil},
		{"if", IF, nil},
		{"for", FOR, nil},
		{"return", RETURN, nil},
		{"var", VAR, nil},
		{"print", PRINT, nil},
		{"[a-z][0-9a-zA-Z_]*", IDENTIFIER, func(s string) { lval.val = s }},
	}

	longestMatch := struct {
		pattern string
		found   string
		token   int
		f       func(s string)
	}{}

	for _, pattern := range patterns {
		found := regexp.MustCompile("^" + pattern.pattern).FindString(l.program)

		if len(found) > len(longestMatch.found) {
			longestMatch.pattern = pattern.pattern
			longestMatch.token = pattern.token
			longestMatch.found = found
			longestMatch.f = pattern.f
		}
	}

	if longestMatch.pattern == "" {
		log.Fatal("no matching token")
	}

	if f := longestMatch.f; f != nil {
		f(longestMatch.found)
	}

	advance(len(longestMatch.found))

	return longestMatch.token
}

func (l *lexer) Error(e string) {
	log.Fatal(e)
}

func Parse(program string) *Program {
	l := lexer{program: program}
	yyParse(&l)
	return l.result
}
