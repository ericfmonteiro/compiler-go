package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

const (
	T_PROGRAMA      int = 1
	T_FIM               = 2
	T_VARIAVEIS         = 3
	T_VIRGULA           = 4
	T_PONTO_VIRGULA     = 5
	T_SE                = 6
	T_SENAO             = 7
	T_FIM_SE            = 8
	T_ENQUANTO          = 9
	T_FIM_ENQUANTO      = 10
	T_PARA              = 11
	T_SETA              = 12
	T_ATE               = 13
	T_FIM_PARA          = 14
	T_LER               = 15
	T_ABRE_PAR          = 16
	T_FECHA_PAR         = 17
	T_ESCREVER          = 18
	T_MAIOR             = 19
	T_MENOR             = 20
	T_MAIOR_IGUAL       = 21
	T_MENOR_IGUAL       = 22
	T_IGUAL             = 23
	T_DIFERENTE         = 24
	T_MAIS              = 25
	T_MENOS             = 26
	T_VEZES             = 27
	T_DIVIDIDO          = 28
	T_RESTO             = 29
	T_ELEVADO           = 30
	T_NUMERO            = 31
	T_ID                = 32
	T_FIM_FONTE         = 90
	T_ERRO_LEX          = 98
	T_NULO              = 99
	FIM_ARQUIVO         = 26
)

var arquivoFinal *bufio.Reader
var lookAhead rune
var token int
var lexema string
var ponteiro int
var linhaFonte string
var linhaAtual int
var colunaAtual int
var mensagemDeErro string

var tokensIdentificados string

func main() {
	file, err := os.Open("/Users/I539613/Projects/University/compiler-go/test.txt")
	if err != nil {
		return
	}
	defer file.Close()

	arquivoFinal = bufio.NewReader(file)

	linhaAtual = 0
	colunaAtual = 0
	ponteiro = 0
	linhaFonte = ""
	token = T_NULO
	mensagemDeErro = ""

	moveLookAhead()

	for token != T_FIM_FONTE && token != T_ERRO_LEX {
		searchNextToken()
		mostraToken()
	}

	if token == T_ERRO_LEX {
		fmt.Println(mensagemDeErro)
	} else {
		fmt.Println("Sem erro léxico!")
	}

}

func searchNextToken() error {
	auxLexema := ""
	// Salto espa�oes enters e tabs at� o inicio do proximo token
	for lookAhead == 9 || lookAhead == '\n' || lookAhead == 8 || lookAhead == 11 || lookAhead == 12 || lookAhead == '\r' || lookAhead == 32 {
		moveLookAhead()
	}

	/*--------------------------------------------------------------*
	 * Caso o primeiro caracter seja alfabetico, procuro capturar a *
	 * sequencia de caracteres que se segue a ele e classifica-la   *
	 *--------------------------------------------------------------*/
	if (lookAhead >= 'A') && (lookAhead <= 'Z') {
		auxLexema += string(lookAhead)
		moveLookAhead()

		for ((lookAhead >= 'A') && (lookAhead <= 'Z')) || ((lookAhead >= '0') && (lookAhead <= '9')) || (lookAhead == '_') {
			auxLexema += string(lookAhead)
			moveLookAhead()
		}

		lexema = auxLexema

		// Classify the token as a reserved word or an identifier
		//var token int // Assuming token is an integer type
		switch lexema {
		case "PROGRAMA":
			token = T_PROGRAMA
		case "FIM":
			token = T_FIM
		case "VARIAVEIS":
			token = T_VARIAVEIS
		case "SE":
			token = T_SE
		case "SENAO":
			token = T_SENAO
		case "FIM_SE":
			token = T_FIM_SE
		case "ENQUANTO":
			token = T_ENQUANTO
		case "FIM_ENQUANTO":
			token = T_FIM_ENQUANTO
		case "PARA":
			token = T_PARA
		case "ATE":
			token = T_ATE
		case "FIM_PARA":
			token = T_FIM_PARA
		case "LER":
			token = T_LER
		case "ESCREVER":
			token = T_ESCREVER
		default:
			token = T_ID // TODO ver pq todo token que deveria ser erro léxico estar retornando um T_ID
		}

	} else if (lookAhead >= '0') && (lookAhead <= '9') {
		auxLexema += string(lookAhead)
		moveLookAhead()
		for (lookAhead >= '0') && (lookAhead <= '9') {
			auxLexema += string(lookAhead)
			moveLookAhead()
		}
		token = T_NUMERO
	} else if lookAhead == '(' {
		auxLexema += string(lookAhead)
		token = T_ABRE_PAR
		moveLookAhead()
	} else if lookAhead == ')' {
		auxLexema += string(lookAhead)
		token = T_FECHA_PAR
		moveLookAhead()
	} else if lookAhead == ';' {
		auxLexema += string(lookAhead)
		token = T_PONTO_VIRGULA
		moveLookAhead()
	} else if lookAhead == ',' {
		auxLexema += string(lookAhead)
		token = T_VIRGULA
		moveLookAhead()
	} else if lookAhead == '+' {
		auxLexema += string(lookAhead)
		token = T_MAIS
		moveLookAhead()
	} else if lookAhead == '-' {
		auxLexema += string(lookAhead)
		token = T_MENOS
		moveLookAhead()
	} else if lookAhead == '*' {
		auxLexema += string(lookAhead)
		moveLookAhead()
		if lookAhead == '*' {
			auxLexema += string(lookAhead)
			moveLookAhead()
			token = T_ELEVADO
		} else {
			token = T_VEZES
		}
	} else if lookAhead == '/' {
		auxLexema += string(lookAhead)
		token = T_DIVIDIDO
		moveLookAhead()
	} else if lookAhead == '%' {
		auxLexema += string(lookAhead)
		token = T_RESTO
		moveLookAhead()
	} else if lookAhead == '<' {
		auxLexema += string(lookAhead)
		moveLookAhead()
		if lookAhead == '>' {
			auxLexema += string(lookAhead)
			moveLookAhead()
			token = T_DIFERENTE
		} else if lookAhead == '-' {
			auxLexema += string(lookAhead)
			moveLookAhead()
			token = T_SETA
		} else if lookAhead == '=' {
			auxLexema += string(lookAhead)
			moveLookAhead()
			token = T_MENOR_IGUAL
		} else {
			token = T_MENOR
		}
	} else if lookAhead == '>' {
		auxLexema += string(lookAhead)
		moveLookAhead()
		if lookAhead == '=' {
			auxLexema += string(lookAhead)
			moveLookAhead()
			token = T_MAIOR_IGUAL
		} else {
			token = T_MAIOR
		}
	} else if lookAhead == FIM_ARQUIVO {
		token = T_FIM_FONTE
	} else {
		token = T_ERRO_LEX // EXAMPLE OF LEXICAL ERROR: &
		mensagemDeErro = "Erro Léxico na linha: " + fmt.Sprint(linhaAtual) + "\nReconhecido ao atingir a coluna: " + fmt.Sprint(colunaAtual) + "\nLinha do Erro: <" + string(linhaFonte) + ">\nToken desconhecido: " + string(lookAhead)
		auxLexema += string(lookAhead)
	}
	lexema = auxLexema
	fmt.Println(token)
	return nil
}

func moveLookAhead() error {
	if ponteiro+1 > len(linhaFonte) {
		linhaAtual++
		var err error
		ponteiro = 0
		linhaFonte, err = arquivoFinal.ReadString('\n')
		if err == io.EOF && linhaFonte == "" {
			lookAhead = FIM_ARQUIVO
		} else {
			linhaFonte += "\r\n"
			lookAhead = rune(linhaFonte[ponteiro])
		}
	} else {
		lookAhead = rune(linhaFonte[ponteiro])
	}
	if unicode.IsLetter(lookAhead) {
		lookAhead = unicode.ToUpper(lookAhead)
	}
	ponteiro++
	colunaAtual = ponteiro + 1

	return nil
}

func mostraToken() {
	var tokenLexema string
	switch token {
	case T_PROGRAMA:
		tokenLexema = "T_PROGRAMA"
	case T_FIM:
		tokenLexema = "T_FIM"
	case T_VARIAVEIS:
		tokenLexema = "T_VARIAVEIS"
	case T_VIRGULA:
		tokenLexema = "T_VIRGULA"
	case T_PONTO_VIRGULA:
		tokenLexema = "T_PONTO_VIRGULA"
	case T_SE:
		tokenLexema = "T_SE"
	case T_SENAO:
		tokenLexema = "T_SENAO"
	case T_FIM_SE:
		tokenLexema = "T_FIM_SE"
	case T_ENQUANTO:
		tokenLexema = "T_ENQUANTO"
	case T_FIM_ENQUANTO:
		tokenLexema = "T_FIM_ENQUANTO"
	case T_PARA:
		tokenLexema = "T_PARA"
	case T_SETA:
		tokenLexema = "T_SETA"
	case T_ATE:
		tokenLexema = "T_ATE"
	case T_FIM_PARA:
		tokenLexema = "T_FIM_PARA"
	case T_LER:
		tokenLexema = "T_LER"
	case T_ABRE_PAR:
		tokenLexema = "T_ABRE_PAR"
	case T_FECHA_PAR:
		tokenLexema = "T_FECHA_PAR"
	case T_ESCREVER:
		tokenLexema = "T_ESCREVER"
	case T_MAIOR:
		tokenLexema = "T_MAIOR"
	case T_MENOR:
		tokenLexema = "T_MENOR"
	case T_MAIOR_IGUAL:
		tokenLexema = "T_MAIOR_IGUAL"
	case T_MENOR_IGUAL:
		tokenLexema = "T_MENOR_IGUAL"
	case T_IGUAL:
		tokenLexema = "T_IGUAL"
	case T_DIFERENTE:
		tokenLexema = "T_DIFERENTE"
	case T_MAIS:
		tokenLexema = "T_MAIS"
	case T_MENOS:
		tokenLexema = "T_MENOS"
	case T_VEZES:
		tokenLexema = "T_VEZES"
	case T_DIVIDIDO:
		tokenLexema = "T_DIVIDIDO"
	case T_RESTO:
		tokenLexema = "T_RESTO"
	case T_ELEVADO:
		tokenLexema = "T_ELEVADO"
	case T_NUMERO:
		tokenLexema = "T_NUMERO"
	case T_ID:
		tokenLexema = "T_ID"
	case T_FIM_FONTE:
		tokenLexema = "T_FIM_FONTE"
	case T_ERRO_LEX:
		tokenLexema = "T_ERRO_LEX"
	case T_NULO:
		tokenLexema = "T_NULO"
	default:
		tokenLexema = "N/A"
	}
	fmt.Println(tokenLexema + " ( " + lexema + " )")
	acumulaToken(tokenLexema + " ( " + lexema + " )")
	tokenLexema += lexema
}

func acumulaToken(tokenIdentificado string) {
	tokensIdentificados += tokenIdentificado
	tokensIdentificados += "\n"
}

// func validateStringLexema() {
// 	if lexema == "PROGRAMA" {
// 		token = T_PROGRAMA
// 	} else if lexema == "FIM" {
// 		token = T_FIM
// 	} else if lexema == "VARIAVEIS" {
// 		token = T_VARIAVEIS
// 	} else if lexema == "SE" {
// 		token = T_SE
// 	} else if lexema == "SENAO" {
// 		token = T_SENAO
// 	} else if lexema == "FIM_SE" {
// 		token = T_FIM_SE
// 	} else if lexema == "ENQUANTO" {
// 		token = T_ENQUANTO
// 	} else if lexema == "FIM_ENQUANTO" {
// 		token = T_FIM_ENQUANTO
// 	} else if lexema == "PARA" {
// 		token = T_PARA
// 	} else if lexema == "ATE" {
// 		token = T_ATE
// 	} else if lexema == "FIM_PARA" {
// 		token = T_FIM_PARA
// 	} else if lexema == "LER" {
// 		token = T_LER
// 	} else if lexema == "ESCREVER" {
// 		token = T_ESCREVER
// 	} else {
// 		token = T_ID
// 	}
// }
