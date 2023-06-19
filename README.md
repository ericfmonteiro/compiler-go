# compiler-go
Compiler based on an academic activity grammar prompt in Golang

# How to run
To run this code, you must have golang installed in your machine. After that, execute the command go run main.go to execute the main file.

# Semantic
The semantic part validates if the variables are declared or not.


# Gramatic

<G> ::= 'PROGRAM' <LIST> <FUNCS> <CMDS> 'END'
<LIST> ::= 'VARIABLES' <VARS> ';'
<VARS> ::= <VAR> ',' <VARS>
<VARS> ::= <VAR>
<VAR> ::= <ID>
<FUNCS> ::= <FUNC> ';' <FUNCS>
<FUNCS> ::= <FUNC>
<FUNC> ::= 'FUNCTION' <ID> '(' <VARS> ')' <CMDS> 'RETURN' <E> 'END_FUNTION'
<CMDS> ::= <CMD> ';' <CMDS>
<CMDS> ::= <CMD>
<CMD> ::= <CMD_IF>
<CMD> ::= <CMD_WHILE>
<CMD> ::= <CMD_FOR>
<CMD> ::= <CMD_ASSIGNMENT>
<CMD> ::= <CMD_READ>
<CMD> ::= <CMD_WRITE>
<CMD_IF> ::= 'IF' '(' <CONDITION> ')' <CMDS> 'END_IF'
<CMD_IF> ::= 'IF' '(' <CONDITION> ')' <CMDS> 'ELSE' <CMDS> 'END_IF'
<CMD_WHILE> ::= 'WHILE' <CONDITION> <CMDS> 'END_WHILE'
<CMD_FOR> ::= 'FOR' <VAR> '<-' <E> 'TO' <E> <CMDS> 'END_FOR'
<CMD_ASSIGNMENT> ::= <VAR> '<-' <E>
<CMD_READ> ::= 'READ' '(' <VAR> ')'
<CMD_WRITE> ::= 'WRITE' '(' <E> ')'
<CMD_CALL> ::= 'CALL' <VAR> <FUNC_CALL>
<FUNC_CALL> ::= <ID> '(' <VARS> ')'
<CONDITION> ::= <E> '>' <E>
<CONDITION> ::= <E> '>=' <E>
<CONDITION> ::= <E> '!=' <E>
<CONDITION> ::= <E> '<=' <E>
<CONDITION> ::= <E> '<' <E>
<CONDITION> ::= <E> '==' <E>
<E> ::= <E> '+' <T>
<E> ::= <E> '-' <T>
<E> ::= <T>
<T> ::= <T> '*' <F>
<T> ::= <T> '/' <F>
<T> ::= <T> '%' <F>
<T> ::= <F>
<F> ::= '-' <X>
<F> ::= <X> '**' <F>
<F> ::= <X>
<X> ::= '(' <E> ')'
<X> ::= [0-9]+('.'[0-9]+)
<X> ::= <VAR>
<ID> ::= [A-Z]+([A-Z]_[0-9])