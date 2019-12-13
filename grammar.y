%{
package nanogo
%}

%union{
  val interface{}
  declaration *Declaration
  declarations []*Declaration
  typ Type
  types []Type
  expression Expression
  expressions []Expression
  statement Statement
  statements []Statement
  block Block
  name_and_type struct{
    Name string
    Type Type
  }
  name_and_types []struct{
    Name string
    Type Type
  }
}

%token<> INT64
%token<> FLOAT64
%token<> BOOL
%token<val> BOOL_VALUE
%token<val> INT_VALUE
%token<val> FLOAT_VALUE
%token<> FUNC
%token<> LEFT_BRACE
%token<> RIGHT_BRACE
%token<> LEFT_PARENTHESIS
%token<> RIGHT_PARENTHESIS
%token<> LEFT_BRACKET
%token<> RIGHT_BRACKET
%token<> MINUS
%token<> PLUS
%token<> ASTERISK
%token<> SLASH
%token<> EQUAL
%token<> EQUAL_EQUAL
%token<> EXCLAMATION_EQUAL
%token<> RETURN
%token<> COMMA
%token<> IF
%token<> LESS
%token<> LESS_EQUAL
%token<> GREATER
%token<> GREATER_EQUAL
%token<> VAR
%token<> FOR
%token<> PRINT
%token<val> IDENTIFIER

%left EQUAL_EQUAL EXCLAMATION_EQUAL LESS LESS_EQUAL GREATER GREATER_EQUAL
%left PLUS MINUS
%left ASTERISK SLASH

%type<> program
%type<declaration> declaration
%type<typ> type
%type<types> types
%type<expression> expression
%type<expressions> expressions
%type<statement> statement
%type<statements> statements
%type<block> block
%type<name_and_type> name_and_type
%type<name_and_types> name_and_types

%start program

%%

program:
  { yylex.(*lexer).result = &Program{} }
| program declaration
  { yylex.(*lexer).result.Declarations = append(yylex.(*lexer).result.Declarations, $2) }
| program FUNC IDENTIFIER LEFT_PARENTHESIS name_and_types RIGHT_PARENTHESIS LEFT_PARENTHESIS types RIGHT_PARENTHESIS block
  {
    typ := &FunctionType{Return: $8}
    args := []string{}
    for _, nameAndType := range $5 {
      args = append(args, nameAndType.Name)
      typ.Args = append(typ.Args, nameAndType.Type)
    }
    yylex.(*lexer).result.Declarations = append(
      yylex.(*lexer).result.Declarations, &Declaration{$3.(string), typ})
    yylex.(*lexer).result.Assignments = append(
      yylex.(*lexer).result.Assignments, &Assignment{$3.(string), &Function{typ, args, $10}})
  }

declaration: VAR IDENTIFIER type
  { $$ = &Declaration{$2.(string), $3} }

type: INT64
  { $$ = &IntType{} }
| FLOAT64
  { $$ = &FloatType{} }
| BOOL
  { $$ = &BoolType{} }
| FUNC LEFT_PARENTHESIS types RIGHT_PARENTHESIS LEFT_PARENTHESIS types RIGHT_PARENTHESIS
  { $$ = &FunctionType{$3, $6} }

types:
  { $$ = []Type{} }
| type
  { $$ = []Type{$1} }
| types COMMA type
  { $$ = append($1, $3) }

name_and_type: IDENTIFIER type
  { $$ = struct{ Name string; Type Type }{$1.(string), $2} }

name_and_types:
  { $$ = []struct{ Name string; Type Type}{} }
| name_and_type
  { $$ = []struct{ Name string; Type Type}{$1} }
| name_and_types COMMA name_and_type
  { $$ = append($1, $3) }

statement: VAR IDENTIFIER type
  { $$ = &Declaration{$2.(string), $3} }
| IDENTIFIER EQUAL expression
  { $$ = &Assignment{$1.(string), $3} }
| RETURN expression
  { $$ = &Return{$2} }
| IF expression block
  { $$ = &If{$2, $3} }
| FOR expression block
  { $$ = &For{$2, $3} }
| PRINT LEFT_PARENTHESIS expression RIGHT_PARENTHESIS
  { $$ = &Print{$3} }
| block
  { $$ = $1 }
| IDENTIFIER LEFT_PARENTHESIS expressions RIGHT_PARENTHESIS
  { $$ = &Application{&Variable{$1.(string)}, $3} }

statements:
  { $$ = []Statement{} }
| statements statement
  { $$ = append($1, $2) }

block: LEFT_BRACE statements RIGHT_BRACE
  { $$ = Block($2) }

expression: expression PLUS expression
  { $$ = &Add{$1, $3} }
| expression MINUS expression
  { $$ = &Sub{$1, $3} }
| expression ASTERISK expression
  { $$ = &Mul{$1, $3} }
| expression SLASH expression
  { $$ = &Div{$1, $3} }
| INT_VALUE
  { $$ = &Int{$1.(int64)} }
| FLOAT_VALUE
  { $$ = &Float{$1.(float64)} }
| BOOL_VALUE
  { $$ = &Bool{$1.(bool)} }
| IDENTIFIER
  { $$ = &Variable{$1.(string)} }
| LEFT_PARENTHESIS expression RIGHT_PARENTHESIS
  { $$ = $2 }
| IDENTIFIER LEFT_PARENTHESIS expressions RIGHT_PARENTHESIS
  { $$ = &Application{&Variable{$1.(string)}, $3} }
| expression EQUAL_EQUAL expression
  { $$ = &Equal{$1, $3} }
| expression EXCLAMATION_EQUAL expression
  { $$ = &Not{&Equal{$1, $3}} }
| expression LESS expression
  { $$ = &LessThan{$1, $3} } 
| expression LESS_EQUAL expression
  { $$ = &Not{&LessThan{$3, $1}} }
| expression GREATER expression
  { $$ = &LessThan{$3, $1} }
| expression GREATER_EQUAL expression
  { $$ = &Not{&LessThan{$1, $3}} }
| FUNC LEFT_PARENTHESIS name_and_types RIGHT_PARENTHESIS LEFT_PARENTHESIS types RIGHT_PARENTHESIS block
  {
    typ := &FunctionType{Return: $6}
    args := []string{}
    for _, nameAndType := range $3 {
      args = append(args, nameAndType.Name)
      typ.Args = append(typ.Args, nameAndType.Type)
    }
    
    $$ = &Function{typ, args, $8}
  }

expressions:
  { $$ = []Expression{} }
| expression
  { $$ = []Expression{$1} }
| expressions COMMA expression
  { $$ = append($1, $3) }
