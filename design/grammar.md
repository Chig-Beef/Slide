# The grammatical design of slide
{construct} means 0 or many
(construct} means 1 or many
(construct)x means x of a construct
[construct] means that construct is optional

program -> {statement}
variableDeclaration -> assignment ';'
elementAssignment -> indexUnary assignment ';'
ifBlock -> 'if' condition block ['elif' condition block] ['else' block]
newType -> 'typedef' IDENTIFIER TYPE ';'
foreverLoop -> 'forever' block 
rangeLoop -> 'range' expression block
forLoop -> 'for' [assignment] ';' [condition] ';' [assignment | (('++' | '--' IDENTIFIER))] block
funcDef -> 'fun' [methodReceiver] IDENTIFIER '(' [IDENTIFIER complexType] {',' IDENTIFIER complexType} ')' complexType block
methodReceiver -> '(' IDENTIFIER complexType ')'
block -> '{' {statement} '}'
expression -> value {OPERATOR value}
value -> VALUE | (UNARY value) | makeArray | bracketedValue | funcCall | structNew
loneIncrement -> ('++' | '--') PROPERTY ';'
switchStatement

funcCall -> 'call'  IDENTIFIER '(' [expression [{',' expression}]] ')' ';'
structDef -> 'struct' IDENTIFIER '{' {IDENTIFIER TYPE ';'} '}'
retStatement -> 'return' [expression] ';'
breakStatement -> 'break' [VALUE] ';'
contStatement -> 'continue' [VALUE] ';'
enumDef -> 'enum' IDENTIFIER '{' {IDENTIFIER ','} '}'
condition -> expression that must return a yes or noable value
assignment -> IDENTIFIER [TYPE] '=' expression 

PTYPE is TYPE but that can include a dereference
VALUE is an IDENTIFIER or a PRIMATIVE
IDENTIFIER is a word, such as a variable name, function name
TYPE is any type that is primitive or created by the user
UNARY is an operator such as + or - that can be used in front of a value, it can also mean a pointer reference or dereference
OPERATOR is any operator such as +, -, or |
PRIMATIVE is usually a literal number or string
