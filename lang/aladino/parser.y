%{
// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

var base int

func setAST(l AladinoLexer, root Expr) {
    l.(*AladinoLex).ast = root
}
%}

// fields inside this union end up as the fields in a structure known
// as ${PREFIX}SymType, of which a reference is passed to the lexer.
%union{
    str string
    int int
    ast Expr
    astList []Expr
    bool bool
}

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct
%type <ast> expr
%type <astList> expr_list

// same for terminals
%token <str> TIMESTAMP RELATIVETIMESTAMP IDENTIFIER STRINGLITERAL TK_CMPOP TK_LAMBDA
%token <int> NUMBER
%token <bool> TRUE
%token <bool> FALSE

%left TK_OR
%left TK_AND
%left TK_EQ TK_NEQ TK_CMPOP
%left TK_NOT

%%

prog :
      expr { setAST(Aladinolex, $1) }
;

expr :
      TK_NOT expr        { $$ = BuildNotOp($2) }
    | expr TK_AND expr   { $$ = BuildAndOp($1, $3) }
    | expr TK_OR expr    { $$ = BuildOrOp($1, $3) }
    | expr TK_EQ expr    { $$ = BuildEqOp($1, $3) }
    | expr TK_NEQ expr   { $$ = BuildNeqOp($1, $3) }
    | expr TK_CMPOP expr { $$ = BuildCmpOp($1, $2, $3) }
    | '(' expr ')'       { $$ = $2 }
    | TIMESTAMP          { $$ = BuildTimeConst($1) }
    | RELATIVETIMESTAMP  { $$ = BuildRelativeTimeConst($1) }
    | NUMBER             { $$ = BuildIntConst($1) }
    | STRINGLITERAL      { $$ = BuildStringConst($1) }
    | '[' expr_list ']'  { $$ = BuildArray($2) }
    | '$' IDENTIFIER     { $$ = BuildVariable($2) }
    | TRUE               { $$ = BuildBoolConst(true) }
    | FALSE              { $$ = BuildBoolConst(false) }
    | '$' IDENTIFIER '(' expr_list ')' 
        { $$ = BuildFunctionCall(BuildVariable($2), $4) }
    | '(' expr_list TK_LAMBDA expr  ')'      { $$ = BuildLambda($2, $4) }
;

expr_list :
      expr ',' expr_list  { $$ = append([]Expr{$1}, $3...) }
    | expr                { $$ = []Expr{$1} }
    |                     { $$ = []Expr{} }
;

%%      /*  start  of  programs  */
