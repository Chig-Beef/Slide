# The grammatical design of slide
{construct} means 0 or many
(construct} means 1 or many
(construct)x means x of a construct
[construct] means that construct is optional

program -> {statement}
variableDeclaration -> assignment ';'
newType -> 'typedef' IDENTIFIER TYPE
ifBlock -> 'if' condition block ['elif' condition block] ['else' block]
foreverLoop -> 'forever' block 
rangeLoop -> 'range' expression block
forLoop -> 'for' [assignment] ';' [condition] ';' assignment block
funcCall -> 'call' IDENTIFIER '(' [ expression [{',' expression}]] ')' ';'
structDef -> 'struct' IDENTIFIER '{' {IDENTIFIER TYPE ';'} '}'
funcDef -> 'fun' IDENTIFIER '(' [IDENTIFIER PTYPE] {',' IDENTIFIER PTYPE} ')' [PTYPE] block
retStatement -> 'return' [expression] ';'
breakStatement -> 'break' [VALUE] ';'
contStatement -> 'continue' [VALUE] ';'
enumDef -> 'enum' IDENTIFIER '{' {IDENTIFIER ','} '}'
condition -> expression that must return a yes or noable value
expression -> [UNARY] VALUE {OPERATOR [UNARY] VALUE}
assignment -> IDENTIFIER [TYPE] '=' expression | 
block -> '{' {statement} '}'

PTYPE is TYPE but that can include a dereference
VALUE is an IDENTIFIER or a PRIMATIVE or a call
IDENTIFIER is a word, such as a variable name, function name
TYPE is any type that is primitive or created by the user
UNARY is an operator such as + or - that can be used in front of a value, it can also mean a pointer reference or dereference
OPERATOR is any operator such as +, -, or |
PRIMATIVE is usually a literal number or string
