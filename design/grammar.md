# The grammatical design of slide
{construct} means 0 or many
(construct} means 1 or many
(construct)x means x of a construct
[construct] means that construct is optional

program -> {statement}
statement -> assignment ';'
statement -> 'if' condition '{' {statement} '}'
statement -> 'forever' '{' {statement} '}'
statement -> 'range' expression '{' {statement} '}'
statement -> 'for' [assignment] ';' [condition] ';' assignment '{' {statement} '}'
statement -> 'call' call ';'
statement -> 'struct' IDENTIFIER '{' {IDENTIFIER TYPE ';'} '}'
statement -> 'fun' IDENTIFIER '(' [IDENTIFIER TYPE] {',' IDENTIFIER TYPE} ')' [TYPE] '{' {statement} '}'
statement -> 'return' [expression] ';'
statement -> 'break' [VALUE] ';'
statement -> 'continue' [VALUE] ';'
condition -> expression that must return a yes or noable value
expression -> [UNARY] VALUE {OPERATOR [UNARY] VALUE}
assignment -> IDENTIFIER [TYPE] '=' expression | 
call -> IDENTIFIER '(' [ expression [{',' expression}]] ')'

VALUE is an IDENTIFIER or a PRIMATIVE or a call
IDENTIFIER is a word, such as a variable name, function name
TYPE is any type that is primitive or created by the user
UNARY is an operator such as + or - that can be used in front of a value, it can also mean a pointer reference or dereference
OPERATOR is any operator such as +, -, or |
PRIMATIVE is usually a literal number or string
