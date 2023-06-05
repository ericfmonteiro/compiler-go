package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

const (
	T_PROGRAM      int = 1
	T_END              = 2
	T_VARIABLES        = 3
	T_COLON            = 4
	T_SEMI_COLON       = 5
	T_IF               = 6
	T_ELSE             = 7
	T_END_IF           = 8
	T_WHILE            = 9
	T_END_WHILE        = 10
	T_FOR              = 11
	T_ARROW            = 12
	T_TO               = 13
	T_END_FOR          = 14
	T_READ             = 15
	T_OPEN_PAR         = 16
	T_CLOSE_PAR        = 17
	T_WRITE            = 18
	T_BIGGER           = 19
	T_LESS             = 20
	T_BIGGER_EQUAL     = 21
	T_LESS_EQUAL       = 22
	T_EQUAL            = 23
	T_DIFFERENT        = 24
	T_PLUS             = 25
	T_MINUS            = 26
	T_TIMES            = 27
	T_DIVIDED          = 28
	T_REST             = 29
	T_RAISED           = 30
	T_NUMBER           = 31
	T_ID               = 32
	T_FUNCTION         = 33
	T_END_FUNCTION     = 34
	T_RETURN           = 35
	T_CALL             = 36
	T_FUNC_CALL        = 37
	T_END_SOURCE       = 90
	T_ERROR_LEX        = 98
	T_NULL             = 99
	END_FILE           = 226
)

var finalFile *bufio.Reader
var lookAhead rune
var token int
var lexema string
var pointer int
var fontLine string
var currentLine int
var currentColumn int

func main() {
	// Get file path
	path, err := os.Getwd()
	if err != nil {
		return
	}

	// Get file to be compiled
	file, err := os.Open(path + "/code.txt")
	if err != nil {
		return
	}
	defer file.Close()

	finalFile = bufio.NewReader(file)

	currentLine = 0
	currentColumn = 0
	pointer = 0
	fontLine = ""
	token = T_NULL

	// First token position
	moveLookAhead()
	searchNextToken()

	// Lexical and syntax analisys
	syntaxAnalisys()

	// Log if there is no lexical error
	fmt.Println("No lexical and syntax errors!")
}

func syntaxAnalisys() {
	g()
}

// <G> ::= 'PROGRAM' <LIST> <CMDS> 'END'
func g() {
	if token == T_PROGRAM {
		searchNextToken()
		list()
		functions() // VERIFY
		cmds()
		if token == T_END {
			searchNextToken()
		} else {
			handleSyntaxError("END")
		}
	} else {
		handleSyntaxError("PROGRAM")
	}
}

// <list> ::= 'VARIABLES' <VARS> ';'
func list() {
	if token == T_VARIABLES {
		searchNextToken()
		vars()
		if token == T_SEMI_COLON {
			searchNextToken()
		} else {
			handleSyntaxError(";")
		}
	} else {
		handleSyntaxError("VARIABLES")
	}
}

// <VARS> ::= <VAR> , <VARS> | <VAR>
func vars() {
	varFunc()
	for token == T_COLON {
		searchNextToken()
		varFunc()
	}
}

// <VAR> ::= <ID>
func varFunc() {
	id()
}

// <ID> ::= [A-Z]+([A-Z]_[0-9])*
func id() {
	if token == T_ID {
		searchNextToken()
	} else {
		handleSyntaxError("IDENTIFIER")
	}
}

// <FUNCS> ::= <FUNC> ';' <FUNCS>
func functions() {
	function()
	for token == T_SEMI_COLON {
		searchNextToken()
		function()
	}
}

// <FUNC> ::= 'FUNCTION' <ID> '(' <VARS> ')' <CMDS> 'RETURN' <E> 'END_FUNTION'
func function() {
	if token == T_FUNCTION {
		searchNextToken()
		id()
		if token == T_OPEN_PAR {
			searchNextToken()
			vars()
			if token == T_CLOSE_PAR {
				searchNextToken()
				cmds()
				if token == T_RETURN {
					searchNextToken()
					e()
					if token == T_END_FUNCTION {
						searchNextToken()
					} else {
						handleSyntaxError("END_FUNCTION")
					}
				} else {
					handleSyntaxError("RETURN")
				}
			} else {
				handleSyntaxError(")")
			}
		} else {
			handleSyntaxError("(")
		}
	}
}

// <CMDS> ::= <CMD> ; <CMDS> | <CMD>
func cmds() {
	cmd()
	for token == T_SEMI_COLON {
		searchNextToken()
		cmd()
	}
}

// <CMD> ::= <CMD_IF>
// <CMD> ::= <CMD_WHILE>
// <CMD> ::= <CMD_FOR>
// <CMD> ::= <CMD_ASSIGNMENT>
// <CMD> ::= <CMD_READ>
// <CMD> ::= <CMD_WRITE>
// <CMD> ::= <CMD_CALL>
func cmd() {
	switch token {
	case T_IF:
		cmd_if()
		break
	case T_WHILE:
		cmd_while()
		break
	case T_FOR:
		cmd_for()
		break
	case T_ID:
		cmd_assignment()
		break
	case T_READ:
		cmd_read()
		break
	case T_WRITE:
		cmd_write()
		break
	case T_CALL:
		cmd_call()
		break
	default:
		handleSyntaxError("COMMAND NOT IDENTIFIED")
	}
}

// <CMD_IF> ::= 'IF' '(' <CONDITION> ')' <CMDS> 'END_IF'
// <CMD_IF> ::= 'IF' '(' <CONDITION> ')' <CMDS> 'ELSE' <CMDS> 'END_IF'
func cmd_if() {
	if token == T_IF {
		searchNextToken()
		if token == T_OPEN_PAR {
			searchNextToken()
			condition()
			if token == T_CLOSE_PAR {
				searchNextToken()
				cmds()
				if token == T_ELSE {
					searchNextToken()
					cmds()
				}
				if token == T_END_IF {
					searchNextToken()
				} else {
					handleSyntaxError("END_IF")
				}
			} else {
				handleSyntaxError(")")
			}
		} else {
			handleSyntaxError("(")
		}
	}
}

// <CMD_WHILE> ::= 'WHILE' <CONDITION> <CMDS> 'END_WHILE'
func cmd_while() {
	if token == T_WHILE {
		searchNextToken()
		condition()
		cmds()
		if token == T_END_WHILE {
			searchNextToken()
		} else {
			handleSyntaxError("END_WHILE")
		}
	} else {
		handleSyntaxError("WHILE")
	}
}

// <CMD_FOR> ::= 'FOR' <VAR> '<-' <E> 'TO' <E> <CMDS> 'END_FOR'
func cmd_for() {
	if token == T_FOR {
		searchNextToken()
		varFunc()
		if token == T_ARROW {
			searchNextToken()
			e()
			if token == T_TO {
				searchNextToken()
				e()
				cmds()
				if token == T_END_FOR {
					searchNextToken()
				} else {
					handleSyntaxError("END_FOR")
				}
			} else {
				handleSyntaxError("TO")
			}
		} else {
			handleSyntaxError("(<-)")
		}
	} else {
		handleSyntaxError("FOR")
	}
}

// <CMD_ASSIGNMENT> ::= <VAR> '<-' <E>
func cmd_assignment() {
	varFunc()
	if token == T_ARROW {
		searchNextToken()
		e()
	} else {
		handleSyntaxError("<-")
	}
}

// <CMD_READ> ::= 'READ' '(' <VAR> ')'
func cmd_read() {
	if token == T_READ {
		searchNextToken()
		if token == T_OPEN_PAR {
			searchNextToken()
			varFunc()
			if token == T_CLOSE_PAR {
				searchNextToken()
			} else {
				handleSyntaxError(")")
			}
		} else {
			handleSyntaxError("(")
		}
	} else {
		handleSyntaxError("READ")
	}
}

// <CMD_WRITE> ::= 'WRITE' '(' <E> ')'
func cmd_write() {
	if token == T_WRITE {
		searchNextToken()
		if token == T_OPEN_PAR {
			searchNextToken()
			e()
			if token == T_CLOSE_PAR {
				searchNextToken()
			} else {
				handleSyntaxError(")")
			}
		} else {
			handleSyntaxError("(")
		}
	} else {
		handleSyntaxError("WRITE")
	}
}

// <CMD_CALL> ::= 'CALL' <VAR> <FUNC_CALL>
func cmd_call() {
	if token == T_CALL {
		searchNextToken()
		varFunc()
		function_call()
	} else {
		handleSyntaxError("CALL")
	}
}

// <FUNC_CALL> ::= <ID> '(' <VARS> ')'
func function_call() {
	if token == T_ID {
		searchNextToken()
		if token == T_OPEN_PAR {
			searchNextToken()
			vars()
			if token == T_CLOSE_PAR {
				searchNextToken()
			} else {
				handleSyntaxError(")")
			}
		} else {
			handleSyntaxError("(")
		}
	}
}

// <CONDITION> ::= <E> '>' <E>
// <CONDITION> ::= <E> '>=' <E>
// <CONDITION> ::= <E> '!-' <E>
// <CONDITION> ::= <E> '<=' <E>
// <CONDITION> ::= <E> '<' <E>
// <CONDITION> ::= <E> '==' <E>
func condition() {
	e()
	switch token {
	case T_BIGGER:
		searchNextToken()
		e()
		break
	case T_LESS:
		searchNextToken()
		e()
		break
	case T_BIGGER_EQUAL:
		searchNextToken()
		e()
		break
	case T_LESS_EQUAL:
		searchNextToken()
		e()
		break
	case T_EQUAL:
		searchNextToken()
		e()
		break
	case T_DIFFERENT:
		searchNextToken()
		e()
		break
	default:
		handleSyntaxError("OPERATOR")
	}
}

// <E> ::= <E> + <T>
// <E> ::= <E> - <T>
// <E> ::= <T>
func e() {
	t()
	for (token == T_PLUS) || (token == T_MINUS) {
		searchNextToken()
		t()
	}
}

// <T> ::= <T> * <F>
// <T> ::= <T> / <F>
// <T> ::= <T> % <F>
// <T> ::= <F>
func t() {
	f()
	for (token == T_TIMES) || (token == T_DIVIDED) || (token == T_REST) {
		searchNextToken()
		f()
	}
}

// <F> ::= -<F>
// <F> ::= <X> ** <F>
// <F> ::= <X>
func f() {
	if token == T_MINUS {
		searchNextToken()
		f()
	} else {
		x()
		for token == T_RAISED {
			searchNextToken()
			x()
		}
	}
}

// <X> ::= '(' <E> ')'
// <X> ::= [0-9]+('.'[0-9]+)
// <X> ::= <VAR>
func x() {
	switch token {
	case T_ID:
		searchNextToken()
	case T_NUMBER:
		searchNextToken()
	case T_OPEN_PAR:
		{
			searchNextToken()
			e()
			if token == T_CLOSE_PAR {
				searchNextToken()
			} else {
				handleSyntaxError(")")
			}
		}
	default:
		handleSyntaxError("INVALID FACTOR")
	}
}

func searchNextToken() {
	auxLexema := ""
	for lookAhead == 9 || lookAhead == '\n' || lookAhead == 8 || lookAhead == 11 || lookAhead == 12 || lookAhead == '\r' || lookAhead == 32 {
		moveLookAhead()
	}

	if (lookAhead >= 'A') && (lookAhead <= 'Z') {
		auxLexema += string(lookAhead)
		moveLookAhead()

		for ((lookAhead >= 'A') && (lookAhead <= 'Z')) || ((lookAhead >= '0') && (lookAhead <= '9')) || (lookAhead == '_') {
			auxLexema += string(lookAhead)
			moveLookAhead()
		}

		lexema = auxLexema

		switch lexema {
		case "PROGRAM":
			token = T_PROGRAM
		case "END":
			token = T_END
		case "VARIABLES":
			token = T_VARIABLES
		case "IF":
			token = T_IF
		case "ELSE":
			token = T_ELSE
		case "END_IF":
			token = T_END_IF
		case "WHILE":
			token = T_WHILE
		case "END_WHILE":
			token = T_END_WHILE
		case "FOR":
			token = T_FOR
		case "TO":
			token = T_TO
		case "END_FOR":
			token = T_END_FOR
		case "READ":
			token = T_READ
		case "WRITE":
			token = T_WRITE
		case "FUNCTION":
			token = T_FUNCTION
		case "END_FUNCTION":
			token = T_END_FUNCTION
		case "RETURN":
			token = T_RETURN
		case "CALL":
			token = T_CALL
		default:
			token = T_ID
		}

	} else if (lookAhead >= '0') && (lookAhead <= '9') {
		auxLexema += string(lookAhead)
		moveLookAhead()
		for (lookAhead >= '0') && (lookAhead <= '9') {
			auxLexema += string(lookAhead)
			moveLookAhead()
		}
		token = T_NUMBER
	} else if lookAhead == '(' {
		auxLexema += string(lookAhead)
		token = T_OPEN_PAR
		moveLookAhead()
	} else if lookAhead == ')' {
		auxLexema += string(lookAhead)
		token = T_CLOSE_PAR
		moveLookAhead()
	} else if lookAhead == ';' {
		auxLexema += string(lookAhead)
		token = T_SEMI_COLON
		moveLookAhead()
	} else if lookAhead == ',' {
		auxLexema += string(lookAhead)
		token = T_COLON
		moveLookAhead()
	} else if lookAhead == '+' {
		auxLexema += string(lookAhead)
		token = T_PLUS
		moveLookAhead()
	} else if lookAhead == '-' {
		auxLexema += string(lookAhead)
		token = T_MINUS
		moveLookAhead()
	} else if lookAhead == '*' {
		auxLexema += string(lookAhead)
		moveLookAhead()
		if lookAhead == '*' {
			auxLexema += string(lookAhead)
			moveLookAhead()
			token = T_RAISED
		} else {
			token = T_TIMES
		}
	} else if lookAhead == '/' {
		auxLexema += string(lookAhead)
		token = T_DIVIDED
		moveLookAhead()
	} else if lookAhead == '%' {
		auxLexema += string(lookAhead)
		token = T_REST
		moveLookAhead()
	} else if lookAhead == '<' {
		auxLexema += string(lookAhead)
		moveLookAhead()
		if lookAhead == '-' {
			auxLexema += string(lookAhead)
			moveLookAhead()
			token = T_ARROW
		} else if lookAhead == '=' {
			auxLexema += string(lookAhead)
			moveLookAhead()
			token = T_LESS_EQUAL
		} else {
			token = T_LESS
		}
	} else if lookAhead == '>' {
		auxLexema += string(lookAhead)
		moveLookAhead()
		if lookAhead == '=' {
			auxLexema += string(lookAhead)
			moveLookAhead()
			token = T_BIGGER_EQUAL
		} else {
			token = T_BIGGER
		}
	} else if lookAhead == '!' {
		auxLexema += string(lookAhead)
		moveLookAhead()
		if lookAhead == '=' {
			auxLexema += string(lookAhead)
			moveLookAhead()
			token = T_DIFFERENT
		}
	} else if lookAhead == END_FILE {
		token = T_END_SOURCE
	} else {
		token = T_ERROR_LEX
		auxLexema += string(lookAhead)
		handleLexicalError()
	}
	lexema = auxLexema
}

func moveLookAhead() {
	if pointer+1 > len(fontLine) {
		currentLine++
		var err error
		pointer = 0
		fontLine, err = finalFile.ReadString('\n')
		if err == io.EOF && fontLine == "" {
			lookAhead = END_FILE
		} else {
			fontLine += "\r\n"
			lookAhead = rune(fontLine[pointer])
		}
	} else {
		lookAhead = rune(fontLine[pointer])
	}
	if unicode.IsLetter(lookAhead) {
		lookAhead = unicode.ToUpper(lookAhead)
	}
	pointer++
	currentColumn = pointer + 1
}

func handleSyntaxError(expected string) {
	fmt.Println("Syntax Error. \nLine: " + fmt.Sprint(currentLine) + "\nColumn: " + fmt.Sprint(currentColumn) + "\nError: \n" + fontLine + expected + " expected, but found: " + lexema)
	os.Exit(1)
}

func handleLexicalError() {
	fmt.Println("Lexical Error. \nLine: " + fmt.Sprint(currentLine) + "\nColumn: " + fmt.Sprint(currentColumn) + "\nError: \n" + fontLine + "Token unknown: " + string(lookAhead))
	os.Exit(1)
}
