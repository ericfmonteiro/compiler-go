package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
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

// Semantics
var lastLexema string
var symbolTable = make(map[string]int)
var node_1 NodeSemanticStack
var node_2 NodeSemanticStack
var semanticStack SemanticStack

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
	fmt.Println("No lexical, syntax or semantic errors!")
}

func syntaxAnalisys() {
	g()
}

// <G> ::= 'PROGRAM' <LIST> <CMDS> 'END'
func g() {
	if token == T_PROGRAM {
		searchNextToken()
		list()
		functions()
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

// <LIST> ::= 'VARIABLES' <VARS> ';'
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
	semanticAnalisys(2)
}

// <VAR> ::= <ID>
func variable() {
	id()
	semanticAnalisys(4)
}

// <VARS> ::= <VAR> , <VARS> | <VAR>
func varsParam() {
	variable()
	for token == T_COLON {
		searchNextToken()
		variable()
	}
}

// <ID> ::= [A-Z]+([A-Z]_[0-9])*
func id() {
	if token == T_ID {
		searchNextToken()
	} else {
		handleSyntaxError("IDENTIFIER")
	}
}

// ADD VALIDATION IN SEMANTICS
// <FUNCS> ::= <FUNC> ';' <FUNCS>
func functions() {
	function()
	for token == T_SEMI_COLON {
		searchNextToken()
		function()
	}
}

// ADD VALIDATION IN SEMANTICS
// <FUNC> ::= 'FUNCTION' <ID> '(' <VARS> ')' <CMDS> 'RETURN' <E> 'END_FUNTION'
func function() {
	if token == T_FUNCTION {
		searchNextToken()
		id()
		if token == T_OPEN_PAR {
			searchNextToken()
			varsParam()
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
			semanticAnalisys(17)
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
		semanticAnalisys(15)
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
		variable()
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
	variable()
	if token == T_ARROW {
		searchNextToken()
		e()
		semanticAnalisys(3)
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
			variable()
			if token == T_CLOSE_PAR {
				searchNextToken()
				semanticAnalisys(14)
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
				semanticAnalisys(25)
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
		variable()
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
			varsParam()
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
		semanticAnalisys(19)
		break
	case T_LESS:
		searchNextToken()
		e()
		semanticAnalisys(20)
		break
	case T_BIGGER_EQUAL:
		searchNextToken()
		e()
		semanticAnalisys(21)
		break
	case T_LESS_EQUAL:
		searchNextToken()
		e()
		semanticAnalisys(22)
		break
	case T_EQUAL:
		searchNextToken()
		e()
		semanticAnalisys(23)
		break
	case T_DIFFERENT:
		searchNextToken()
		e()
		semanticAnalisys(24)
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
		switch token {
		case T_PLUS:
			searchNextToken()
			t()
			semanticAnalisys(5)
		case T_MINUS:
			searchNextToken()
			t()
			semanticAnalisys(6)
		}
	}
}

// <T> ::= <T> * <F>
// <T> ::= <T> / <F>
// <T> ::= <T> % <F>
// <T> ::= <F>
func t() {
	f()
	for (token == T_TIMES) || (token == T_DIVIDED) || (token == T_REST) {
		switch token {
		case T_TIMES:
			searchNextToken()
			t()
			semanticAnalisys(7)
		case T_DIVIDED:
			searchNextToken()
			t()
			semanticAnalisys(8)
		case T_REST:
			searchNextToken()
			t()
			semanticAnalisys(9)
		}
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
			semanticAnalisys(10)
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
		semanticAnalisys(11)
	case T_NUMBER:
		searchNextToken()
		semanticAnalisys(12)
	case T_OPEN_PAR:
		{
			searchNextToken()
			e()
			if token == T_CLOSE_PAR {
				searchNextToken()
			} else {
				handleSyntaxError(")")
			}
			semanticAnalisys(13)
		}
	default:
		handleSyntaxError("INVALID FACTOR")
	}
}

func searchNextToken() {
	auxLexema := ""

	if lexema != "" {
		lastLexema = lexema
	}

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

func semanticAnalisys(ruleNumber int) {
	switch ruleNumber {
	case 2:
		insertInSymbolTable(lastLexema)
	case 3:
		node_2 = semanticStack.pop()
		node_1 = semanticStack.pop()
	case 4:
		if isInSymbolTable(lastLexema) {
			semanticStack.push(lastLexema, 4)
		}
	case 5:
		node_2 = semanticStack.pop()
		node_1 = semanticStack.pop()
		semanticStack.push(node_1.getLowerCaseCode()+"+"+node_2.getLowerCaseCode(), 5)
		break
	case 6:
		node_2 = semanticStack.pop()
		node_1 = semanticStack.pop()
		semanticStack.push(node_1.getLowerCaseCode()+"-"+node_2.getLowerCaseCode(), 6)
		break
	case 7:
		node_2 = semanticStack.pop()
		node_1 = semanticStack.pop()
		semanticStack.push(node_1.getLowerCaseCode()+"*"+node_2.getLowerCaseCode(), 7)
		break
	case 8:
		node_2 = semanticStack.pop()
		node_1 = semanticStack.pop()
		semanticStack.push(node_1.getLowerCaseCode()+"/"+node_2.getLowerCaseCode(), 8)
		break
	case 9:
		node_2 = semanticStack.pop()
		node_1 = semanticStack.pop()
		semanticStack.push(node_1.getLowerCaseCode()+"%"+node_2.getLowerCaseCode(), 9)
		break
	case 10:
		node_2 = semanticStack.pop()
		node_1 = semanticStack.pop()
		semanticStack.push(node_1.getLowerCaseCode()+"**"+node_2.getLowerCaseCode(), 10)
		break
	case 11:
		if isInSymbolTable(lastLexema) {
			semanticStack.push(lastLexema, 11)
		}
		break
	case 12:
		semanticStack.push(lastLexema, 12)
		break
	case 13:
		node_1 = semanticStack.pop()
		semanticStack.push("("+node_1.getLowerCaseCode()+")", 13)
		break
	case 14:
		node_1 = semanticStack.pop()
		break
	case 15:
		node_1 = semanticStack.pop()
		break
	case 17:
		node_1 = semanticStack.pop()
		break
	case 19:
		node_2 = semanticStack.pop()
		node_1 = semanticStack.pop()
		semanticStack.push(node_1.getLowerCaseCode()+" > "+node_2.getLowerCaseCode(), 19)
		break
	case 20:
		node_2 = semanticStack.pop()
		node_1 = semanticStack.pop()
		semanticStack.push(node_1.getLowerCaseCode()+" < "+node_2.getLowerCaseCode(), 20)
		break
	case 21:
		node_2 = semanticStack.pop()
		node_1 = semanticStack.pop()
		semanticStack.push(node_1.getLowerCaseCode()+" >= "+node_2.getLowerCaseCode(), 21)
		break
	case 22:
		node_2 = semanticStack.pop()
		node_1 = semanticStack.pop()
		semanticStack.push(node_1.getLowerCaseCode()+" <= "+node_2.getLowerCaseCode(), 22)
		break
	case 23:
		node_2 = semanticStack.pop()
		node_1 = semanticStack.pop()
		semanticStack.push(node_1.getLowerCaseCode()+" == "+node_2.getLowerCaseCode(), 23)
		break
	case 24:
		node_2 = semanticStack.pop()
		node_1 = semanticStack.pop()
		semanticStack.push(node_1.getLowerCaseCode()+" != "+node_2.getLowerCaseCode(), 24)
		break
	case 25:
		node_1 = semanticStack.pop()
		break
	}
}

func isInSymbolTable(lastLexema string) bool {
	_, ok := symbolTable[lastLexema]
	if !ok {
		handleSemanticNotDeclaredError(lastLexema)
		return false
	} else {
		return true
	}
}

func insertInSymbolTable(lastLexema string) {
	_, ok := symbolTable[lastLexema]
	if ok {
		handleSemanticAlreadyDeclaredError(lastLexema)
	} else {
		symbolTable[lastLexema] = 0
	}
}

func handleSemanticNotDeclaredError(lastLexema string) {
	fmt.Println("Semantic Error. \nLine: " + fmt.Sprint(currentLine) + " \nError: \nVariable " + lastLexema + " is not declared!")
	os.Exit(1)
}

func handleSemanticAlreadyDeclaredError(lastLexema string) {
	fmt.Println("Semantic Error. \nLine: " + fmt.Sprint(currentLine) + " \nError: \nVariable " + lastLexema + " is already declared!")
	os.Exit(1)
}

type SemanticStack struct {
	stack []NodeSemanticStack
}

func (p *SemanticStack) pop() NodeSemanticStack {
	var nps NodeSemanticStack
	if len(p.stack) > 0 {
		nps = p.stack[len(p.stack)-1]
		p.stack = p.stack[:len(p.stack)-1]
	}

	return nps
}

func (p *SemanticStack) push(c string, r int) NodeSemanticStack {
	nps := NodeSemanticStack{code: c, implementedRule: r}
	p.stack = append(p.stack, nps)

	return nps
}

type NodeSemanticStack struct {
	code            string
	implementedRule int
}

func (nps *NodeSemanticStack) getLowerCaseCode() string {
	return strings.ToLower(nps.code)
}
