package main

import (
	"bufio"
	"fmt"
	"os"
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

var arquivoFinal *bufio.Scanner

var lookAhead string
var token int
var lexema string
var ponteiro int
var linhaFonte string
var linhaAtual int
var colunaAtual int
var mensagemDeErro string

//   // Lista de tokens
//   var int T_PROGRAMA        =   1;
//   var int T_FIM             =   2;
//   var int T_VARIAVEIS       =   3;
//   var int T_VIRGULA         =   4;
//   var int T_PONTO_VIRGULA   =   5;
//   var int T_SE              =   6;
//   var int T_SENAO           =   7;
//   var int T_FIM_SE          =   8;
//   var int T_ENQUANTO        =   9;
//   var int T_FIM_ENQUANTO    =  10;
//   var int T_PARA            =  11;
//   var int T_SETA            =  12;
//   var int T_ATE             =  13;
//   var int T_FIM_PARA        =  14;
//   var int T_LER             =  15;
//   var int T_ABRE_PAR        =  16;
//   var int T_FECHA_PAR       =  17;
//   var int T_ESCREVER        =  18;
//   var int T_MAIOR           =  19;
//   var int T_MENOR           =  20;
//   var int T_MAIOR_IGUAL     =  21;
//   var int T_MENOR_IGUAL     =  22;
//   var int T_IGUAL           =  23;
//   var int T_DIFERENTE       =  24;
//   var int T_MAIS            =  25;
//   var int T_MENOS           =  26;
//   var int T_VEZES           =  27;
//   var int T_DIVIDIDO        =  28;
//   var int T_RESTO           =  29;
//   var int T_ELEVADO         =  30;
//   var int T_NUMERO          =  31;
//   var int T_ID              =  32;

//   var int T_FIM_FONTE       =  90;
//   var int T_ERRO_LEX        =  98;
//   var int T_NULO            =  99;

//   var int FIM_ARQUIVO       =  26;

// //   static File arqFonte;
// //   static BufferedReader rdFonte;
// //   static File arqDestino;

//   var char   lookAhead;
//   var int    token;
//   var String lexema;
//   var int    ponteiro;
//   var String linhaFonte;
//   var int    linhaAtual;
//   var int    colunaAtual;
//   var String mensagemDeErro;
//   //var StringBuffer tokensIdentificados = new StringBuffer();

func abreFile() (*bufio.Scanner, error) {
	file, err := os.Open("/Users/I539613/Projects/University/compiler-go/test.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	arquivoFinal = scanner
	//var text string
	// for scanner.Scan() {
	// 	text += scanner.Text() + "\n"
	// 	fmt.Print(text)
	// }

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return scanner, nil

	// fmt.Println(text)

	// return nil
}

func main() {
	file, err := os.Open("/Users/I539613/Projects/University/compiler-go/test.txt")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	arquivoFinal = scanner
	//var text string
	// for scanner.Scan() {
	// 	text += scanner.Text() + "\n"
	// 	fmt.Print(text)
	// }

	if err := scanner.Err(); err != nil {
		return
	}

	//var arquivoFinal *bufio.Scanner
	//fmt.Println("SAD")
	// arquivoFinal, err := abreFile()
	// if err != nil {
	// 	return
	// }

	//fmt.Println(arquivoFinal)

	for arquivoFinal.Scan() {
		//fmt.Println(arquivoFinal.Text())
		buscaProximoToken(arquivoFinal.Text())
	}

	// fmt.Println(text)
	//abreDestino();
	// linhaAtual     = 0;
	// colunaAtual    = 0;
	// ponteiro       = 0;
	// linhaFonte     = "";
	// token          = T_NULO;
	// mensagemDeErro = "";
	// movelookAhead();
	// while ( ( token != T_FIM_FONTE ) && ( token != T_ERRO_LEX ) ) {
	// 		buscaProximoToken();
	// 		mostraToken();
	// }
	// if ( token == T_ERRO_LEX ) {
	// 	JOptionPane.showMessageDialog( null, mensagemDeErro, "Erro L�xico!", JOptionPane.ERROR_MESSAGE );
	// } else {
	// 	JOptionPane.showMessageDialog( null, "An�lise L�xica terminada sem erros l�xicos", "An�lise L�xica terminada!", JOptionPane.INFORMATION_MESSAGE );
	// }
	// exibeTokens();
	// gravaSaida( arqDestino );
	// fechaFonte();
}

// func abreArquivo() error {
// 	// Abrir arquivo
// 	arquivo, err := ioutil.ReadFile("/Users/I539613/Projects/University/compiler-go/test.txt")
// 	if err != nil {
// 		fmt.Println("Erro ao abrir arquivo:", err)
// 		return err
// 	}
// 	fmt.Print(string(arquivo))
// 	arquivoFinal = string(arquivo)
// 	//defer arquivo.Close()

// 	// Ler arquivo
// 	// buffer := make([]byte, 1024)
// 	// for {
// 	// 	n, err := arquivo.Read(buffer)
// 	// 	if err != nil && n == 0 {
// 	// 		break
// 	// 	}
// 	// 	fmt.Print(string(buffer[:n]))
// 	// }
// 	return nil
// }

func buscaProximoToken(line string) {
	for _, v := range line {
		if v == 9 || v == '\n' || v == 8 || v == 11 || v == 12 || v == '\r' || v == 32 {
			continue
		}
		fmt.Println(v)

	}

	//int i, j;

	//StringBuffer sbLexema = new StringBuffer( "" );

	//   // Salto espa�oes enters e tabs at� o inicio do proximo token
	// 	while ( ( lookAhead == 9 ) ||
	// 		  ( lookAhead == '\n' ) ||
	// 		  ( lookAhead == 8 ) ||
	// 		  ( lookAhead == 11 ) ||
	// 		  ( lookAhead == 12 ) ||
	// 		  ( lookAhead == '\r' ) ||
	// 		  ( lookAhead == 32 ) )
	//   {
	// 	  movelookAhead();
	//   }

	//   if !unicode.IsSpace(char) {
	// 	// Se não é um espaço em branco, mover o cursor para trás
	// 	_ = reader.UnreadRune()
	// 	break
	// }

	//   /*--------------------------------------------------------------*
	//    * Caso o primeiro caracter seja alfabetico, procuro capturar a *
	//    * sequencia de caracteres que se segue a ele e classifica-la   *
	//    *--------------------------------------------------------------*/
	//   if ( ( lookAhead >= 'A' ) && ( lookAhead <= 'Z' ) ) {
	// 	  sbLexema.append( lookAhead );
	// 	  movelookAhead();

	// 	  while ( ( ( lookAhead >= 'A' ) && ( lookAhead <= 'Z' ) ) ||
	// 			  ( ( lookAhead >= '0' ) && ( lookAhead <= '9' ) ) || ( lookAhead == '_' ) )
	// 	  {
	// 		  sbLexema.append( lookAhead );
	// 		  movelookAhead();
	// 	  }

	// 	  lexema = sbLexema.toString();

	// 	  /* Classifico o meu token como palavra reservada ou id */
	// 	  if ( lexema.equals( "PROGRAMA" ) )
	// 		  token = T_PROGRAMA;
	// 	  else if ( lexema.equals( "FIM" ) )
	// 		  token = T_FIM;
	// 	  else if ( lexema.equals( "VARIAVEIS" ) )
	// 		  token = T_VARIAVEIS;
	// 	  else if ( lexema.equals( "SE" ) )
	// 		  token = T_SE;
	// 	  else if ( lexema.equals( "SENAO" ) )
	// 		  token = T_SENAO;
	// 	  else if ( lexema.equals( "FIM_SE" ) )
	// 		  token = T_FIM_SE;
	// 	  else if ( lexema.equals( "ENQUANTO" ) )
	// 		  token = T_ENQUANTO;
	// 	  else if ( lexema.equals( "FIM_ENQUANTO" ) )
	// 		  token = T_FIM_ENQUANTO;
	// 	  else if ( lexema.equals( "PARA" ) )
	// 		  token = T_PARA;
	// 	  else if ( lexema.equals( "ATE" ) )
	// 		  token = T_ATE;
	// 	  else if ( lexema.equals( "FIM_PARA" ) )
	// 		  token = T_FIM_PARA;
	// 	  else if ( lexema.equals( "LER" ) )
	// 		  token = T_LER;
	// 	  else if ( lexema.equals( "ESCREVER" ) )
	// 		  token = T_ESCREVER;
	// 	  else {
	// 		  token = T_ID;
	// 	  }
	//   } else if ( ( lookAhead >= '0' ) && ( lookAhead <= '9' ) ) {
	// 	  sbLexema.append( lookAhead );
	// 	  movelookAhead();
	// 	  while ( ( lookAhead >= '0' ) && ( lookAhead <= '9' ) )
	// 	  {
	// 		  sbLexema.append( lookAhead );
	// 		  movelookAhead();
	// 	  }
	// 	  token = T_NUMERO;
	//   } else if ( lookAhead == '(' ){
	// 	  sbLexema.append( lookAhead );
	// 	  token = T_ABRE_PAR;
	// 	  movelookAhead();
	//   } else if ( lookAhead == ')' ){
	// 	  sbLexema.append( lookAhead );
	// 	  token = T_FECHA_PAR;
	// 	  movelookAhead();
	//   } else if ( lookAhead == ';' ){
	// 	  sbLexema.append( lookAhead );
	// 	  token = T_PONTO_VIRGULA;
	// 	  movelookAhead();
	//   } else if ( lookAhead == ',' ){
	// 	  sbLexema.append( lookAhead );
	// 	  token = T_VIRGULA;
	// 	  movelookAhead();
	//   } else if ( lookAhead == '+' ){
	// 	  sbLexema.append( lookAhead );
	// 	  token = T_MAIS;
	// 	  movelookAhead();
	//   } else if ( lookAhead == '-' ){
	// 	  sbLexema.append( lookAhead );
	// 	  token = T_MENOS;
	// 	  movelookAhead();
	//   } else if ( lookAhead == '*' ){
	// 	  sbLexema.append( lookAhead );
	// 	  movelookAhead();
	// 	  if ( lookAhead == '*' ) {
	// 		  sbLexema.append( lookAhead );
	// 		  movelookAhead();
	// 		  token = T_ELEVADO;
	// 	  } else {
	// 		  token = T_VEZES;
	// 	  }
	//   } else if ( lookAhead == '/' ){
	// 	  sbLexema.append( lookAhead );
	// 	  token = T_DIVIDIDO;
	// 	  movelookAhead();
	//   } else if ( lookAhead == '%' ){
	// 	  sbLexema.append( lookAhead );
	// 	  token = T_RESTO;
	// 	  movelookAhead();
	//   } else if ( lookAhead == '<' ){
	// 	  sbLexema.append( lookAhead );
	// 	  movelookAhead();
	// 	  if ( lookAhead == '>' ) {
	// 		  sbLexema.append( lookAhead );
	// 		  movelookAhead();
	// 		  token = T_DIFERENTE;
	// 	  } else if ( lookAhead == '-' ) {
	// 		  sbLexema.append( lookAhead );
	// 		  movelookAhead();
	// 		  token = T_SETA;
	// 	  } else if ( lookAhead == '=' ) {
	// 		  sbLexema.append( lookAhead );
	// 		  movelookAhead();
	// 		  token = T_MENOR_IGUAL;
	// 	  } else {
	// 		  token = T_MENOR;
	// 	  }
	//   } else if ( lookAhead == '>' ){
	// 	  sbLexema.append( lookAhead );
	// 	  movelookAhead();
	// 	  if ( lookAhead == '=' ) {
	// 		  sbLexema.append( lookAhead );
	// 		  movelookAhead();
	// 		  token = T_MAIOR_IGUAL;
	// 	  } else {
	// 		  token = T_MAIOR;
	// 	  }
	//   } else if ( lookAhead == FIM_ARQUIVO ){
	// 	   token = T_FIM_FONTE;
	//   } else {
	// 	  token = T_ERRO_LEX;
	// 	  mensagemDeErro = "Erro L�xico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\nToken desconhecido: " + lookAhead;
	// 	  sbLexema.append( lookAhead );
	//   }

	// lexema = sbLexema.toString();
}

// static void movelookAhead() throws IOException
// {
//   if ( ( ponteiro + 1 ) > linhaFonte.length() ) {
// 	  linhaAtual++;
// 	  ponteiro = 0;
// 	  if ( ( linhaFonte = rdFonte.readLine() ) == null ) {
// 		  lookAhead = FIM_ARQUIVO;
// 	  } else {
// 		  StringBuffer sbLinhaFonte = new StringBuffer( linhaFonte );
// 		  sbLinhaFonte.append( '\13' ).append( '\10' );
// 		  linhaFonte = sbLinhaFonte.toString();
// 		  lookAhead = linhaFonte.charAt( ponteiro );
// 	  }
//   } else {
// 	  lookAhead = linhaFonte.charAt( ponteiro );
//   }
//   if ( ( lookAhead >= 'a' ) &&
// 	   ( lookAhead <= 'z' ) ) {
// 	  lookAhead = (char) ( lookAhead - 'a' + 'A' );
//   }
//   ponteiro++;
//   colunaAtual = ponteiro + 1;
// }

// func movelookAhead() error {
// 	if ponteiro+1 > len(linhaFonte) {
// 		linhaAtual++
// 		ponteiro = 0
// 		linhaFonte, err = rdFonte.ReadString('\n')
// 		if err != nil {
// 			lookAhead = rune(FIM_ARQUIVO)
// 			return err
// 		}
// 		linhaFonte += "\r\n"
// 		lookAhead = rune(linhaFonte[ponteiro])
// 	} else {
// 		lookAhead = rune(linhaFonte[ponteiro])
// 	}
// 	if lookAhead >= 'a' && lookAhead <= 'z' {
// 		lookAhead = lookAhead - 'a' + 'A'
// 	}
// 	ponteiro++
// 	colunaAtual = ponteiro + 1
// 	return nil
// }
