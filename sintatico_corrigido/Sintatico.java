import java.io.UnsupportedEncodingException;
import javax.swing.filechooser.FileFilter;
import java.io.FileNotFoundException;
import javax.swing.JFileChooser;
import javax.swing.JOptionPane;
import java.io.BufferedReader;
import java.io.BufferedWriter;
import javax.swing.JTextArea;
import java.io.IOException;
import java.io.FileReader;
import java.io.FileWriter;
import java.io.File;

/**
 * Sintatico - Primeira versao do sintatico
 * 
 * @author Ricardo Ferreira de Oliveira
 * @author Turma de projeto de compiladores 1/2023
 *

gramatica:

<G> ::= 'PROGRAMA' <LISTA> <CMDS> 'FIM'
<LISTA> ::= 'VARIAVEIS' <VARS> ';'
<VARS> ::= <VAR> , <VARS>
<VARS> ::= <VAR> 
<VAR>  ::= <ID>
<CMDS> ::= <CMD> ; <CMDS>
<CMDS> ::= <CMD>
<CMD> ::= <CMD_SE>
<CMD> ::= <CMD_ENQUANTO>
<CMD> ::= <CMD_PARA>
<CMD> ::= <CMD_ATRIBUICAO>
<CMD> ::= <CMD_LER>
<CMD> ::= <CMD_ESCREVER>
<CMD> ::= <CMD_FACA>
<CMD> ::= <CMD_REPITA>
<CMD> ::= <CMD_CASO>
<CMD_SE> ::= 'SE' '(' <CONDICAO> ')' <CMDS> 'FIM_SE' 
<CMD_SE> ::= 'SE' '(' <CONDICAO> ')' <CMDS> 'SENAO' <CMDS> 'FIM_SE' 
<CMD_ENQUANTO> ::= 'ENQUANTO' <CONDICAO> <CMDS> 'FIM_ENQUANTO'
<CMD_PARA> ::= 'PARA' <VAR> '<-' <E> 'ATE' <E> <CMDS> 'FIM_PARA' 
<CMD_ATRIBUICAO> ::= <VAR> '<-' <E>
<CMD_LER> ::= 'LER' '(' <VAR> ')' 
<CMD_ESCREVER> ::= 'ESCREVER' '(' <E> ')'
// <CMD_FACA> ::= 'FACA' <E> 'VEZES' <CMDS> 'FIM_FACA'
// <CMD_REPITA> ::= 'REPITA' <CMDS> 'ATE' <CONDICAO>
// <CMD_CASO> ::= 'CASO' <CASOS> 'FIM_CASO'
// <CASOS> ::= <CASO> ';' <CASOS>
// <CASOS> ::= <CASO>
// <CASO> ::= 'QUANDO' <CONDICAO> 'FACA' <CMDS>
<CONDICAO> ::= <E> '>' <E> 
<CONDICAO> ::= <E> '>=' <E> 
<CONDICAO> ::= <E> '<>' <E> 
<CONDICAO> ::= <E> '<=' <E> 
<CONDICAO> ::= <E> '<' <E> 
<CONDICAO> ::= <E> '==' <E>
<E> ::= <E> + <T>
<E> ::= <E> - <T>
<E> ::= <T>
<T> ::= <T> * <F>
<T> ::= <T> / <F>
<T> ::= <T> % <F>
<T> ::= <F>
<F> ::= -<X>
<F> ::= <X> ** <F>
<F> ::= <X>
<X> ::= '(' <E> ')'
<X> ::= [0-9]+('.'[0-9]+)
<X> ::= <VAR>
<ID> ::= [A-Z]+([A-Z]_[0-9]*)

*/

public class Sintatico {

	  // Lista de tokens	
  static final int T_PROGRAMA        =   1;
  static final int T_FIM             =   2;
  static final int T_VARIAVEIS       =   3;
  static final int T_VIRGULA         =   4;
  static final int T_PONTO_VIRGULA   =   5;
  static final int T_SE              =   6;
  static final int T_SENAO           =   7;
  static final int T_FIM_SE          =   8;
  static final int T_ENQUANTO        =   9;
  static final int T_FIM_ENQUANTO    =  10;
  static final int T_PARA            =  11;
  static final int T_SETA            =  12;
  static final int T_ATE             =  13;
  static final int T_FIM_PARA        =  14;
  static final int T_LER             =  15;
  static final int T_ABRE_PAR        =  16;
  static final int T_FECHA_PAR       =  17;
  static final int T_ESCREVER        =  18;
  static final int T_MAIOR           =  19;
  static final int T_MENOR           =  20;
  static final int T_MAIOR_IGUAL     =  21;
  static final int T_MENOR_IGUAL     =  22;
  static final int T_IGUAL           =  23;
  static final int T_DIFERENTE       =  24;
  static final int T_MAIS            =  25;
  static final int T_MENOS           =  26;
  static final int T_VEZES           =  27;
  static final int T_DIVIDIDO        =  28;
  static final int T_RESTO           =  29;
  static final int T_ELEVADO         =  30;
  static final int T_NUMERO          =  31;
  static final int T_ID              =  32;
  static final int T_FACA            =  33;
  static final int T_VEZES_FACA      =  34;
  static final int T_FIM_FACA        =  35;
  static final int T_REPITA          =  36;
  static final int T_CASO            =  37;
  static final int T_QUANDO          =  38;
  static final int T_FIM_CASO        =  39;
  static final int T_PONTO           =  40;
	  
  static final int T_FIM_FONTE       =  90;
  static final int T_ERRO_LEX        =  98;
  static final int T_NULO            =  99;

  static final int FIM_ARQUIVO       = 226;

  static final int E_SEM_ERROS       =   0;
  static final int E_ERRO_LEXICO     =   1;
  static final int E_ERRO_SINTATICO  =   2;

  // Variaveis que surgem no Lexico
  static File arqFonte;
  static BufferedReader rdFonte;
  static File arqDestino;
  static char   lookAhead;
  static int    token;
  static String lexema;
  static int    ponteiro;
  static String linhaFonte;
  static int    linhaAtual;
  static int    colunaAtual;
  static String mensagemDeErro;
  static StringBuffer tokensIdentificados = new StringBuffer();

  // Variaveis adicionadas para o sintatico
  static StringBuffer 	regrasReconhecidas = new StringBuffer();
  static int 			estadoCompilacao;

  public static void main( String s[] ) throws ErroLexicoException
  {
      try {
          abreArquivo();
          abreDestino();
          linhaAtual     = 0;
          colunaAtual    = 0;
          ponteiro       = 0;
          linhaFonte     = "";
          token          = T_NULO;
          mensagemDeErro = "";
          tokensIdentificados.append( "Tokens reconhecidos: \n\n" );
          regrasReconhecidas.append( "\n\nRegras reconhecidas: \n\n" );
          estadoCompilacao 	= E_SEM_ERROS;

          // posiciono no primeiro token
          movelookAhead();
          buscaProximoToken();

          analiseSintatica();

          exibeSaida();

          gravaSaida( arqDestino );

          fechaFonte();

      } catch( FileNotFoundException fnfe ) {
          JOptionPane.showMessageDialog( null, "Arquivo nao existe!", "FileNotFoundException!", JOptionPane.ERROR_MESSAGE );
      } catch( UnsupportedEncodingException uee ) {
          JOptionPane.showMessageDialog( null, "Erro desconhecido", "UnsupportedEncodingException!", JOptionPane.ERROR_MESSAGE );
      } catch( IOException ioe ) {
          JOptionPane.showMessageDialog( null, "Erro de io: " + ioe.getMessage(), "IOException!", JOptionPane.ERROR_MESSAGE );
      } catch( ErroLexicoException ele ) {
          JOptionPane.showMessageDialog( null, ele.getMessage(), "Erro Lexico Exception!", JOptionPane.ERROR_MESSAGE );
      } catch( ErroSintaticoException ese ) {
          JOptionPane.showMessageDialog( null, ese.getMessage(), "Erro Sint�tico Exception!", JOptionPane.ERROR_MESSAGE );
      } finally {
          System.out.println( "Execucao terminada!" );
      }
  }

  static void analiseSintatica() throws IOException, ErroLexicoException, ErroSintaticoException {

      g();

      if ( estadoCompilacao == E_ERRO_LEXICO ) {
          JOptionPane.showMessageDialog( null, mensagemDeErro, "Erro Lexico!", JOptionPane.ERROR_MESSAGE );
      } else if ( estadoCompilacao == E_ERRO_SINTATICO ) {
          JOptionPane.showMessageDialog( null, mensagemDeErro, "Erro Sintatico!", JOptionPane.ERROR_MESSAGE );
      } else {
          JOptionPane.showMessageDialog( null, "Analise Sintatica terminada sem erros", "Analise Sintatica terminada!", JOptionPane.INFORMATION_MESSAGE );
		  acumulaRegraSintaticaReconhecida( "<G>" );
      }
  }
  
  // <G> ::= 'PROGRAMA' <LISTA> <CMDS> 'FIM'
  private static void g() throws IOException, ErroLexicoException, ErroSintaticoException {
	  if ( token == T_PROGRAMA ) {
		  buscaProximoToken();
		  lista();
		  cmds();
		  if ( token == T_FIM ) {
			  buscaProximoToken();
			  acumulaRegraSintaticaReconhecida( "<G> ::= 'PROGRAMA' <LISTA> <CMDS> 'FIM'" );
		  } else {
			  registraErroSintatico( "Erro Sintatico. Linha: " + linhaAtual + "\nColuna: " + colunaAtual + "\nErro: <" + linhaFonte + ">\n('fim') esperado, mas encontrei: " + lexema );
		  }
	  } else {
		  registraErroSintatico( "Erro Sintatico. Linha: " + linhaAtual + "\nColuna: " + colunaAtual + "\nErro: <" + linhaFonte + ">\n('programa') esperado, mas encontrei: " + lexema );
	  }
  }

  // <LISTA> ::= 'VARIAVEIS' <VARS> ';'
  private static void lista() throws IOException, ErroLexicoException, ErroSintaticoException {
	  if ( token == T_VARIAVEIS ) {
		  buscaProximoToken();
		  vars();
		  if ( token == T_PONTO_VIRGULA ) {
			  buscaProximoToken();
			  acumulaRegraSintaticaReconhecida( "<LISTA> ::= 'VARIAVEIS' <VARS> ';'" );
		  } else {
			  registraErroSintatico( "Erro Sintatico. Linha: " + linhaAtual + "\nColuna: " + colunaAtual + "\nErro: <" + linhaFonte + ">\n';' esperado, mas encontrei: " + lexema );
		  }
	  } else {
		  registraErroSintatico( "Erro Sintatico. Linha: " + linhaAtual + "\nColuna: " + colunaAtual + "\nErro: <" + linhaFonte + ">\n('variaveis') esperado, mas encontrei: " + lexema );
	  }
  }
  
  // <VARS> ::= <VAR> , <VARS> | <VAR> 
  private static void vars() throws IOException, ErroLexicoException, ErroSintaticoException {
	  var();
	  while ( token == T_VIRGULA ) {
		  buscaProximoToken();
		  var();
	  }
	  acumulaRegraSintaticaReconhecida( "<VARS> ::= <VAR> , <VARS> | <VAR>" );
  }
  
  // <VAR> ::= <ID> 
  private static void var() throws IOException, ErroLexicoException, ErroSintaticoException {
      id();
	  acumulaRegraSintaticaReconhecida( "<VAR> ::= <ID>" );
  }
  
  // <ID> ::= [A-Z]+([A-Z]_[0-9])*
  private static void id() throws IOException, ErroLexicoException, ErroSintaticoException {
	if ( token == T_ID ) {
		buscaProximoToken();
		acumulaRegraSintaticaReconhecida( "<ID> ::= [A-Z]+([A-Z]_[0-9])*" );
	} else {
		registraErroSintatico( "Erro Sintatico. Linha: " + linhaAtual + "\nColuna: " + colunaAtual + "\nErro: <" + linhaFonte + ">\nEsperava um identificador. Encontrei: " + lexema );
	}
  }
   
  // <CMDS> ::= <CMD> ; <CMDS> | <CMD>
  private static void cmds() throws IOException, ErroLexicoException, ErroSintaticoException {
	  cmd();
	  while ( token == T_PONTO_VIRGULA ) {
		  buscaProximoToken();
		  cmd();
	  } 
	  acumulaRegraSintaticaReconhecida( "<CMDS> ::= <CMD> ; <CMDS> | <CMD>" );
  }
  
  // <CMD> ::= <CMD_SE>
  // <CMD> ::= <CMD_ENQUANTO>
  // <CMD> ::= <CMD_PARA>
  // <CMD> ::= <CMD_ATRIBUICAO>
  // <CMD> ::= <CMD_LER>
  // <CMD> ::= <CMD_ESCREVER>
  private static void cmd() throws IOException, ErroLexicoException, ErroSintaticoException {
      switch ( token ) {
      case T_SE: cmd_se(); break;
      case T_ENQUANTO: cmd_enquanto(); break;
      case T_PARA: cmd_para(); break;
      case T_ID: cmd_atribuicao(); break;
      case T_LER: cmd_ler(); break;
      case T_ESCREVER: cmd_escrever(); break;
      case T_FACA: cmd_faca(); break;
      case T_REPITA: cmd_repita(); break;
      case T_CASO: cmd_caso(); break;
      default:
          registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\nComando nao identificado va aprender a programar pois encontrei: " + lexema );
      }
	  acumulaRegraSintaticaReconhecida( "<CMD> ::= <CMD_SE>|<CMD_ENQUANTO>|<CMD_PARA>|<CMD_ATRIBUICAO>|<CMD_LER>|<CMD_ESCREVER>" );
  }
  
  // <CMD_SE> ::= 'SE' <CONDICAO> <CMDS> 'FIM_SE' 
  // <CMD_SE> ::= 'SE' <CONDICAO> <CMDS> 'SENAO' <CMDS> 'FIM_SE' 
  private static void cmd_se() throws IOException, ErroLexicoException, ErroSintaticoException {
	  if ( token == T_SE) {
		  buscaProximoToken();
		  if ( token == T_ABRE_PAR ) {
			  buscaProximoToken();
			  condicao();
			  if ( token == T_FECHA_PAR ) {
				  buscaProximoToken();
                  cmds();
				  if ( token == T_SENAO ) {
					  buscaProximoToken();
					  cmds();
				  }
				  if ( token == T_FIM_SE ) {
					  buscaProximoToken();
					  acumulaRegraSintaticaReconhecida( "<CMD_SE> ::= 'SE' <CONDICAO> <CMDS> ( 'FIM_SE'|'SENAO' <CMDS> 'FIM_SE' )" );
				  } else {
					  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'fim_se' esperado mas encontrei: " + lexema );  
				  }
			  } else {
				  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n')' esperado mas encontrei: " + lexema );
			  }
		  } else {
			  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'(' esperado mas encontrei: " + lexema ); 
		  }
	  }
  }
  
  // <CMD_ENQUANTO> ::= 'ENQUANTO' <CONDICAO> <CMDS> 'FIM_ENQUANTO'
  private static void cmd_enquanto() throws IOException, ErroLexicoException, ErroSintaticoException {
	  if ( token == T_ENQUANTO ) {
		  buscaProximoToken();
		  condicao();
		  cmds();
		  if ( token == T_FIM_ENQUANTO ) {
			  buscaProximoToken();
			  acumulaRegraSintaticaReconhecida( "<CMD_ENQUANTO> ::= 'ENQUANTO' <CONDICAO> <CMDS> 'FIM_ENQUANTO'" );
		  } else {
			  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'fim enquanto' esperado mas encontrei: " + lexema );
		  }
	  } else {
			  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'enquanto' esperado mas encontrei: " + lexema ); 
	  }	  
  }

  // <CMD_PARA> ::= 'PARA' <VAR> '<-' <E> 'ATE' <E> <CMDS> 'FIM_PARA' 
  private static void cmd_para() throws IOException, ErroLexicoException, ErroSintaticoException {
	  if ( token == T_PARA ) {
		  buscaProximoToken();
		  var();
		  if ( token == T_SETA ) {
			  buscaProximoToken();
			  e();
			  if ( token == T_ATE ) {
				  buscaProximoToken();
				  e();
				  cmds();
				  if ( token == T_FIM_PARA ) {
					  buscaProximoToken();
					  acumulaRegraSintaticaReconhecida( "<CMD_PARA> ::= 'PARA' <VAR> '<-' <E> 'ATE' <E> <CMDS> 'FIM_PARA'" );
				  } else {
					  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'fim_para' esperado mas encontrei: " + lexema );
				  }
			  } else {
				  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'Ate' esperado mas encontrei: " + lexema );
			  }
		  } else {
			  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'<-' esperado mas encontrei: " + lexema ); 
		  }
	  } else {
		  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'Para' esperado mas encontrei: " + lexema );
	  }
  }  
  
  // <CMD_ATRIBUICAO> ::= <VAR> '<-' <E>
  private static void cmd_atribuicao() throws IOException, ErroLexicoException, ErroSintaticoException {
	  var();
	  if ( token == T_SETA ) {
		  buscaProximoToken();
		  e();
		  acumulaRegraSintaticaReconhecida( "<CMD_ATRIBUICAO> ::= <VAR> '<-' <E>" );
	  } else {
		  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'<-' esperado mas encontrei: " + lexema );		  
	  }
  }
  
  // <CMD_LER> ::= 'LER' '(' <VAR> ')' 
  private static void cmd_ler() throws IOException, ErroLexicoException, ErroSintaticoException {
	  if ( token == T_LER ) {
		  buscaProximoToken();
		  if ( token == T_ABRE_PAR ) {
			  buscaProximoToken();
			  var();
			  if ( token == T_FECHA_PAR ) {
				  buscaProximoToken();
				  acumulaRegraSintaticaReconhecida( "<CMD_LER> ::= 'LER' '(' <VAR> ')'" );
			  } else {
				  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n')' esperado mas encontrei: " + lexema );
			  }
		  } else {
			  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'(' esperado mas encontrei: " + lexema ); 
		  }
	  } else {
		  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'Ler' esperado mas encontrei: " + lexema );
	  }
  }

  // <CMD_ESCREVER> ::= 'ESCREVER' '(' <E> ')'
  private static void cmd_escrever() throws IOException, ErroLexicoException, ErroSintaticoException {
	  if ( token == T_ESCREVER ) {
		  buscaProximoToken();
		  if ( token == T_ABRE_PAR ) {
			  buscaProximoToken();
			  e();
			  if ( token == T_FECHA_PAR ) {
				  buscaProximoToken();
				  acumulaRegraSintaticaReconhecida( "<CMD_ESCREVER> ::= 'ESCREVER' '(' <E> ')'" );
			  } else {
				  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n')' esperado mas encontrei: " + lexema );
			  }
		  } else {
			  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'(' esperado mas encontrei: " + lexema ); 
		  }
	  } else {
		  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'Escrever' esperado mas encontrei: " + lexema );
	  }
  }
  
  // <CMD_FACA> ::= 'FACA' <E> 'VEZES' <CMDS> 'FIM_FACA'
  private static void cmd_faca() throws IOException, ErroLexicoException, ErroSintaticoException {
	  if ( token == T_FACA ) {
		  buscaProximoToken();
		  e();
		  if ( token == T_VEZES_FACA ) {
			  buscaProximoToken();
			  cmds();
			  if ( token == T_FIM_FACA ) {
				  buscaProximoToken();
				  acumulaRegraSintaticaReconhecida( "<CMD_FACA> ::= 'FACA' <E> 'VEZES' <CMDS> 'FIM_FACA'" );
			  } else {
				  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'Fim_faca' esperado mas encontrei: " + lexema );
			  }
		  } else {
			  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'Vezes' esperado mas encontrei: " + lexema ); 
		  }
	  } else {
		  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'Faca' esperado mas encontrei: " + lexema );
	  }
  }

  // <CMD_REPITA> ::= 'REPITA' <CMDS> 'ATE' <CONDICAO>
  private static void cmd_repita() throws IOException, ErroLexicoException, ErroSintaticoException {
	  if ( token == T_REPITA ) {
		  buscaProximoToken();
		  cmds();
		  if ( token == T_ATE ) {
			  buscaProximoToken();
			  condicao();
			  acumulaRegraSintaticaReconhecida( "<CMD_REPITA> ::= 'REPITA' <CMDS> 'ATE' <CONDICAO>" );
		  } else {
			  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'Ate' esperado mas encontrei: " + lexema ); 
		  }
	  } else {
		  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'Repita' esperado mas encontrei: " + lexema );
	  }
  }

  // <CMD_CASO> ::= 'CASO' <CASOS> 'FIM_CASO'
  private static void cmd_caso() throws IOException, ErroLexicoException, ErroSintaticoException {
	  if ( token == T_CASO ) {
		  buscaProximoToken();
		  casos();
		  if ( token == T_FIM_CASO ) {
			  buscaProximoToken();
			  acumulaRegraSintaticaReconhecida( "<CMD_CASO> ::= 'CASO' <CASOS> 'FIM_CASO'" );
		  } else {
			  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'Ate' esperado mas encontrei: " + lexema ); 
		  }
	  } else {
		  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'Repita' esperado mas encontrei: " + lexema );
	  }
  }
  
  // <CASOS> ::= <CASO> '.' <CASOS>
  // <CASOS> ::= <CASO>
  private static void casos() throws IOException, ErroLexicoException, ErroSintaticoException {
	  caso();
	  while ( token == T_PONTO ) {
		  buscaProximoToken();
		  caso();
	  } 
	  acumulaRegraSintaticaReconhecida( "<CMDS> ::= <CMD> ; <CMDS> | <CMD>" );
  }
  
  // <CASO> ::= 'QUANDO' <CONDICAO> 'FACA' <CMDS>
  private static void caso() throws IOException, ErroLexicoException, ErroSintaticoException {

	  if ( token == T_QUANDO ) {
		  buscaProximoToken();
		  condicao();
		  if ( token == T_FACA ) {
			  buscaProximoToken();
			  cmds();
			  acumulaRegraSintaticaReconhecida( "<CASO> ::= 'QUANDO' <CONDICAO> 'FACA' <CMDS>" );
		  } else {
			  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'Faca' esperado mas encontrei: " + lexema ); 
		  }
	  } else {
		  registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n'Quando' esperado mas encontrei: " + lexema );
	  }
	  
  }

  // <CONDICAO> ::= <E> '>' <E> 
  // <CONDICAO> ::= <E> '>=' <E> 
  // <CONDICAO> ::= <E> '<>' <E> 
  // <CONDICAO> ::= <E> '<=' <E> 
  // <CONDICAO> ::= <E> '<' <E> 
  // <CONDICAO> ::= <E> '==' <E>
  private static void condicao() throws ErroLexicoException, IOException, ErroSintaticoException {
	  e();
	  switch ( token ) {
	  case T_MAIOR: buscaProximoToken(); e(); break; 
	  case T_MENOR: buscaProximoToken(); e(); break; 
	  case T_MAIOR_IGUAL: buscaProximoToken(); e(); break; 
	  case T_MENOR_IGUAL: buscaProximoToken(); e(); break; 
	  case T_IGUAL: buscaProximoToken(); e(); break; 
	  case T_DIFERENTE: buscaProximoToken(); e(); break;
	  default: registraErroSintatico( "Erro Sintatico. Linha: " + linhaAtual + "\nColuna: " + colunaAtual + "\nErro: <" + linhaFonte + ">\nEsperava um operador logico. Encontrei: " + lexema );
	  }
	  acumulaRegraSintaticaReconhecida( "<CONDICAO> ::= <E> ('>'|'>='|'<>'|'<='|'<'|'==') <E> " );
  }
  
  // <E> ::= <E> + <T>
  // <E> ::= <E> - <T>
  // <E> ::= <T>
  private static void e() throws IOException, ErroLexicoException, ErroSintaticoException {
	  t();
	  while ( (token == T_MAIS) || (token == T_MENOS) ) {
		  buscaProximoToken();
		  t();
	  }
	  acumulaRegraSintaticaReconhecida( "<E> ::= <E> + <T>|<E> - <T>|<T> " );
  }
  
  // <T> ::= <T> * <F>
  // <T> ::= <T> / <F>
  // <T> ::= <T> % <F>
  // <T> ::= <F>
  private static void t() throws IOException, ErroLexicoException, ErroSintaticoException {
	  f();
	  while ( (token == T_VEZES) || (token == T_DIVIDIDO) || (token == T_RESTO) ) {
		  buscaProximoToken();
		  f();
	  }
	  acumulaRegraSintaticaReconhecida( "<T> ::= <T> * <F>|<T> / <F>|<T> % <F>|<F>" );
  }
  
  // <F> ::= -<F>
  // <F> ::= <X> ** <F>
  // <F> ::= <X>     
  private static void f() throws IOException, ErroLexicoException, ErroSintaticoException {
	  if ( token == T_MENOS ) {
		  buscaProximoToken();
		  f();
	  } else {
		  x();
		  while ( token == T_ELEVADO ) {
			  buscaProximoToken();
	          x();
		  }
	  }
	  acumulaRegraSintaticaReconhecida( "<F> ::= -<F>|<X> ** <F>|<X> " );
	  
  }
  
  // <X> ::= '(' <E> ')'
  // <X> ::= [0-9]+('.'[0-9]+)
  // <X> ::= <VAR>
  private static void x() throws IOException, ErroLexicoException, ErroSintaticoException {
	  switch ( token ) {
	  case T_ID: buscaProximoToken(); acumulaRegraSintaticaReconhecida( "<X> ::= <VAR>" ); break;
	  case T_NUMERO: buscaProximoToken(); acumulaRegraSintaticaReconhecida( "<X> ::= [0-9]+('.'[0-9]+)" ); break;
	  case T_ABRE_PAR: {
	       buscaProximoToken(); 
	       e();
	       if ( token == T_FECHA_PAR ) {
	    	   buscaProximoToken();
	    	   acumulaRegraSintaticaReconhecida( "<X> ::= '(' <E> ')'" );
	       } else {
			   registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\n')' esperado mas encontrei: " + lexema );
	       }
	      } break;
	   default: registraErroSintatico( "Erro Sintatico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\nFator invalido: encontrei: " + lexema );   
	  }
  }
  
  static void fechaFonte() throws IOException
  {
      rdFonte.close();
  }

  static void movelookAhead() throws IOException
  {
	  
    if ( ( ponteiro + 1 ) > linhaFonte.length() ) {

    	linhaAtual++;
        ponteiro = 0;
        
        
        if ( ( linhaFonte = rdFonte.readLine() ) == null ) {
            lookAhead = FIM_ARQUIVO;
        } else {

        	StringBuffer sbLinhaFonte = new StringBuffer( linhaFonte );
        	sbLinhaFonte.append( '\13' ).append( '\10' );
        	linhaFonte = sbLinhaFonte.toString();
        	
            lookAhead = linhaFonte.charAt( ponteiro );
        }
    } else {
        lookAhead = linhaFonte.charAt( ponteiro );
    }

    // Se comentar esse if, eu terei uma linguagem 
    // que diferencia minusculas de maiusculas
    if ( ( lookAhead >= 'a' ) &&
         ( lookAhead <= 'z' ) ) {
        lookAhead = (char) ( lookAhead - 'a' + 'A' );
    }

    ponteiro++;
    colunaAtual = ponteiro + 1;
  }

  static void buscaProximoToken() throws IOException, ErroLexicoException
  {
	//int i, j;
        
    StringBuffer sbLexema = new StringBuffer( "" );

    // Salto espaçoes enters e tabs até o inicio do proximo token
  	while ( ( lookAhead == 9 ) ||
	        ( lookAhead == '\n' ) ||
	        ( lookAhead == 8 ) ||
	        ( lookAhead == 11 ) ||
	        ( lookAhead == 12 ) ||
	        ( lookAhead == '\r' ) ||
	        ( lookAhead == 32 ) )
    {
        movelookAhead();
    }

    /*--------------------------------------------------------------*
     * Caso o primeiro caracter seja alfabetico, procuro capturar a *
     * sequencia de caracteres que se segue a ele e classifica-la   *
     *--------------------------------------------------------------*/
    if ( ( lookAhead >= 'A' ) && ( lookAhead <= 'Z' ) ) {   
        sbLexema.append( lookAhead );
        movelookAhead();

        while ( ( ( lookAhead >= 'A' ) && ( lookAhead <= 'Z' ) ) ||
        		( ( lookAhead >= '0' ) && ( lookAhead <= '9' ) ) || ( lookAhead == '_' ) )
        {
            sbLexema.append( lookAhead );
            movelookAhead();
        }

        lexema = sbLexema.toString();  

        /* Classifico o meu token como palavra reservada ou id */
        if ( lexema.equals( "PROGRAMA" ) )
            token = T_PROGRAMA;
        else if ( lexema.equals( "FIM" ) )
            token = T_FIM;
        else if ( lexema.equals( "VARIAVEIS" ) )
            token = T_VARIAVEIS;
        else if ( lexema.equals( "SE" ) )
            token = T_SE;
        else if ( lexema.equals( "SENAO" ) )
            token = T_SENAO;
        else if ( lexema.equals( "FIM_SE" ) )
            token = T_FIM_SE;
        else if ( lexema.equals( "ENQUANTO" ) )
            token = T_ENQUANTO;
        else if ( lexema.equals( "FIM_ENQUANTO" ) )
            token = T_FIM_ENQUANTO;
        else if ( lexema.equals( "PARA" ) )
            token = T_PARA;
        else if ( lexema.equals( "ATE" ) )
            token = T_ATE;
        else if ( lexema.equals( "FIM_PARA" ) )
            token = T_FIM_PARA;
        else if ( lexema.equals( "LER" ) )
            token = T_LER;
        else if ( lexema.equals( "ESCREVER" ) )
            token = T_ESCREVER;
        else if ( lexema.equals( "FACA" ) )
            token = T_FACA;
        else if ( lexema.equals( "VEZES" ) )
            token = T_VEZES_FACA;
        else if ( lexema.equals( "FIM_FACA" ) )
            token = T_FIM_FACA;
        else if ( lexema.equals( "REPITA" ) )
            token = T_REPITA;
        else if ( lexema.equals( "CASO" ) )
            token = T_CASO;
        else if ( lexema.equals( "QUANDO" ) )
            token = T_QUANDO;
        else if ( lexema.equals( "FIM_CASO" ) )
            token = T_FIM_CASO;
        else {
        	token = T_ID;
        }
    } else if ( ( lookAhead >= '0' ) && ( lookAhead <= '9' ) ) {
        sbLexema.append( lookAhead );
        movelookAhead();
        while ( ( lookAhead >= '0' ) && ( lookAhead <= '9' ) )
        {
            sbLexema.append( lookAhead );
            movelookAhead();
        }
        token = T_NUMERO;    	
    } else if ( lookAhead == '(' ){
        sbLexema.append( lookAhead );
        token = T_ABRE_PAR;    	
        movelookAhead();
    } else if ( lookAhead == ')' ){
        sbLexema.append( lookAhead );
        token = T_FECHA_PAR;    	
        movelookAhead();
    } else if ( lookAhead == ';' ){
        sbLexema.append( lookAhead );
        token = T_PONTO_VIRGULA;    	
        movelookAhead();
    } else if ( lookAhead == ',' ){
        sbLexema.append( lookAhead );
        token = T_VIRGULA;    	
        movelookAhead();
    } else if ( lookAhead == '.' ){
        sbLexema.append( lookAhead );
        token = T_PONTO;    	
        movelookAhead();
    } else if ( lookAhead == '+' ){
        sbLexema.append( lookAhead );
        token = T_MAIS;    	
        movelookAhead();
    } else if ( lookAhead == '-' ){
        sbLexema.append( lookAhead );
        token = T_MENOS;    	
        movelookAhead();
    } else if ( lookAhead == '*' ){
        sbLexema.append( lookAhead );
        movelookAhead();
        if ( lookAhead == '*' ) {
            sbLexema.append( lookAhead );
            movelookAhead();
            token = T_ELEVADO;    	
        } else {
            token = T_VEZES;    	
        }
    } else if ( lookAhead == '/' ){
        sbLexema.append( lookAhead );
        token = T_DIVIDIDO;    	
        movelookAhead();
    } else if ( lookAhead == '%' ){
        sbLexema.append( lookAhead );
        token = T_RESTO;    	
        movelookAhead();
    } else if ( lookAhead == '<' ){
        sbLexema.append( lookAhead );
        movelookAhead();
        if ( lookAhead == '>' ) {
            sbLexema.append( lookAhead );
            movelookAhead();
            token = T_DIFERENTE;
        } else if ( lookAhead == '-' ) {  
            sbLexema.append( lookAhead );
            movelookAhead();
            token = T_SETA;
        } else if ( lookAhead == '=' ) {  
            sbLexema.append( lookAhead );
            movelookAhead();
            token = T_MENOR_IGUAL;
        } else {
            token = T_MENOR;    	
        }
    } else if ( lookAhead == '>' ){
        sbLexema.append( lookAhead );
        movelookAhead();
        if ( lookAhead == '=' ) {
            sbLexema.append( lookAhead );
            movelookAhead();
            token = T_MAIOR_IGUAL;
        } else {
            token = T_MAIOR;    	
        }        
    } else if ( lookAhead == FIM_ARQUIVO ){
         token = T_FIM_FONTE;    	
    } else {
    	token = T_ERRO_LEX;
    	sbLexema.append( lookAhead );
    }
        
    lexema = sbLexema.toString();  
    
    mostraToken();
    
    if ( token == T_ERRO_LEX ) {
    	mensagemDeErro = "Erro Léxico na linha: " + linhaAtual + "\nReconhecido ao atingir a coluna: " + colunaAtual + "\nLinha do Erro: <" + linhaFonte + ">\nToken desconhecido: " + lexema;
    	throw new ErroLexicoException( mensagemDeErro );
    } 
  }
  
  static void mostraToken()
  {

    StringBuffer tokenLexema = new StringBuffer( "" );
    
    switch ( token ) {
    case T_PROGRAMA    : tokenLexema.append( "T_PROGRAMA" ); break;
    case T_FIM    : tokenLexema.append( "T_FIM" ); break;
    case T_VARIAVEIS    : tokenLexema.append( "T_VARIAVEIS" ); break;
    case T_VIRGULA    : tokenLexema.append( "T_VIRGULA" ); break;
    case T_PONTO_VIRGULA    : tokenLexema.append( "T_PONTO_VIRGULA" ); break;
    case T_SE    : tokenLexema.append( "T_SE" ); break;
    case T_SENAO    : tokenLexema.append( "T_SENAO" ); break;
    case T_FIM_SE    : tokenLexema.append( "T_FIM_SE" ); break;
    case T_ENQUANTO    : tokenLexema.append( "T_ENQUANTO" ); break;
    case T_FIM_ENQUANTO    : tokenLexema.append( "T_FIM_ENQUANTO" ); break;
    case T_PARA            : tokenLexema.append( "T_PARA" ); break;
    case T_SETA            : tokenLexema.append( "T_SETA" ); break;
    case T_ATE             : tokenLexema.append( "T_ATE" ); break;
    case T_FIM_PARA        : tokenLexema.append( "T_FIM_PARA" ); break;
    case T_LER             : tokenLexema.append( "T_LER" ); break;
    case T_ABRE_PAR        : tokenLexema.append( "T_ABRE_PAR" ); break;
    case T_FECHA_PAR       : tokenLexema.append( "T_FECHA_PAR" ); break;
    case T_ESCREVER        : tokenLexema.append( "T_ESCREVER" ); break;
    case T_MAIOR           : tokenLexema.append( "T_MAIOR" ); break;
    case T_MENOR           : tokenLexema.append( "T_MENOR" ); break;
    case T_MAIOR_IGUAL     : tokenLexema.append( "T_MAIOR_IGUAL" ); break;
    case T_MENOR_IGUAL     : tokenLexema.append( "T_MENOR_IGUAL" ); break;
    case T_IGUAL           : tokenLexema.append( "T_IGUAL" ); break;
    case T_DIFERENTE       : tokenLexema.append( "T_DIFERENTE" ); break;
    case T_MAIS            : tokenLexema.append( "T_MAIS" ); break;
    case T_MENOS           : tokenLexema.append( "T_MENOS" ); break;
    case T_VEZES           : tokenLexema.append( "T_VEZES" ); break;
    case T_DIVIDIDO        : tokenLexema.append( "T_DIVIDIDO" ); break;
    case T_RESTO           : tokenLexema.append( "T_RESTO" ); break;
    case T_ELEVADO         : tokenLexema.append( "T_ELEVADO" ); break;
    case T_NUMERO          : tokenLexema.append( "T_NUMERO" ); break;
    case T_ID              : tokenLexema.append( "T_ID" ); break;
    case T_FIM_FONTE       : tokenLexema.append( "T_FIM_FONTE" ); break;
    case T_ERRO_LEX        : tokenLexema.append( "T_ERRO_LEX" ); break;
    case T_NULO            : tokenLexema.append( "T_NULO" ); break;
    case T_FACA            : tokenLexema.append( "T_FACA" ); break;
    case T_VEZES_FACA      : tokenLexema.append( "T_VEZES_FACA" ); break;
    case T_FIM_FACA        : tokenLexema.append( "T_FIM_FACA" ); break;
    case T_REPITA          : tokenLexema.append( "T_REPITA" ); break;
    case T_CASO            : tokenLexema.append( "T_CASO" ); break;
    case T_QUANDO          : tokenLexema.append( "T_QUANDO" ); break;
    case T_FIM_CASO        : tokenLexema.append( "T_FIM_CASO" ); break;
    case T_PONTO           : tokenLexema.append( "T_PONTO" ); break;
    
    default                : tokenLexema.append( "N/A" ); break;
    }
    System.out.println( tokenLexema.toString() + " ( " + lexema + " )" );
    acumulaToken( tokenLexema.toString() + " ( " + lexema + " )" );
    tokenLexema.append( lexema );
  }
  
  private static void abreArquivo() {

		JFileChooser fileChooser = new JFileChooser();
		
		fileChooser.setFileSelectionMode( JFileChooser.FILES_ONLY );

		FiltroSab filtro = new FiltroSab();
	    
		fileChooser.addChoosableFileFilter( filtro );
		int result = fileChooser.showOpenDialog( null );
		
		if( result == JFileChooser.CANCEL_OPTION ) {
			return;
		}
		
		arqFonte = fileChooser.getSelectedFile();
		abreFonte( arqFonte ); 	

	}


	private static boolean abreFonte( File fileName ) {

		if( arqFonte == null || fileName.getName().trim().equals( "" ) ) {
			JOptionPane.showMessageDialog( null, "Nome de Arquivo Invalido", "Nome de Arquivo Invalido", JOptionPane.ERROR_MESSAGE );
			return false;
		} else {
			linhaAtual = 1;
	        try {
				FileReader fr = new FileReader( arqFonte );
				rdFonte = new BufferedReader( fr );
			} catch (FileNotFoundException e) {
				e.printStackTrace();
			} 
			return true;
		}
	}

	private static void abreDestino() {

		JFileChooser fileChooser = new JFileChooser();
			
		fileChooser.setFileSelectionMode( JFileChooser.FILES_ONLY );

		FiltroSab filtro = new FiltroSab();
		    
		fileChooser.addChoosableFileFilter( filtro );
		int result = fileChooser.showSaveDialog( null );
			
		if( result == JFileChooser.CANCEL_OPTION ) {
			return;
		}
			
		arqDestino = fileChooser.getSelectedFile();
	}	
	

	private static boolean gravaSaida( File fileName ) {

		if( arqDestino == null || fileName.getName().trim().equals( "" ) ) {
			JOptionPane.showMessageDialog( null, "Nome de Arquivo Invalido", "Nome de Arquivo Invalido", JOptionPane.ERROR_MESSAGE );
			return false;
		} else {
			FileWriter fw;
			try {
				System.out.println( arqDestino.toString() );
				System.out.println( tokensIdentificados.toString() );
				fw = new FileWriter( arqDestino );
				BufferedWriter bfw = new BufferedWriter( fw ); 
				bfw.write( tokensIdentificados.toString() );
				bfw.write( regrasReconhecidas.toString() );
				bfw.close();
				JOptionPane.showMessageDialog( null, "Arquivo Salvo: " + arqDestino, "Salvando Arquivo", JOptionPane.INFORMATION_MESSAGE );
			} catch (IOException e) {
				JOptionPane.showMessageDialog( null, e.getMessage(), "Erro de Entrada/Sa�da", JOptionPane.ERROR_MESSAGE );
			} 
			return true;
		}
	}
	
	public static void exibeTokens() {
		
		JTextArea texto = new JTextArea();
		texto.append( tokensIdentificados.toString() );
		JOptionPane.showMessageDialog(null, texto, "Tokens Identificados (token/lexema)", JOptionPane.INFORMATION_MESSAGE );
	}
	
	
	public static void acumulaRegraSintaticaReconhecida( String regra ) {

		regrasReconhecidas.append( regra );
		regrasReconhecidas.append( "\n" );
		
	}
	
	public static void acumulaToken( String tokenIdentificado ) {

		tokensIdentificados.append( tokenIdentificado );
		tokensIdentificados.append( "\n" );
		
	}
	
    public static void exibeSaida() {

        JTextArea texto = new JTextArea();
        texto.append( tokensIdentificados.toString() );
        JOptionPane.showMessageDialog(null, texto, "Analise Lexica", JOptionPane.INFORMATION_MESSAGE );

        texto.setText( regrasReconhecidas.toString() );
        texto.append( "\n\nStatus da Compilacao:\n\n" );
        texto.append( mensagemDeErro );

        JOptionPane.showMessageDialog(null, texto, "Resumo da Compilacao", JOptionPane.INFORMATION_MESSAGE );
    }
    
    static void registraErroSintatico( String msg ) throws ErroSintaticoException {
        if ( estadoCompilacao == E_SEM_ERROS ) {
            estadoCompilacao = E_ERRO_SINTATICO;
            mensagemDeErro = msg;
        }
        throw new ErroSintaticoException( msg ); 
    }
		
}

/**
 * Classe Interna para criacao de filtro de selecao
 */
class FiltroSab extends FileFilter {

	public boolean accept(File arg0) {
	   	 if(arg0 != null) {
	         if(arg0.isDirectory()) {
	       	  return true;
	         }
	         if( getExtensao(arg0) != null) {
	        	 if ( getExtensao(arg0).equalsIgnoreCase( "grm" ) ) {
		        	 return true;
	        	 }
	         };
	   	 }
	     return false;
	}

	/**
	 * Retorna quais extensoes poderao ser escolhidas
	 */
	public String getDescription() {
		return "*.grm";
	}
	
	/**
	 * Retorna a parte com a extensao de um arquivo
	 */
	public String getExtensao(File arq) {
	if(arq != null) {
		String filename = arq.getName();
	    int i = filename.lastIndexOf('.');
	    if(i>0 && i<filename.length()-1) {
	    	return filename.substring(i+1).toLowerCase();
	    };
	}
		return null;
	}
}
