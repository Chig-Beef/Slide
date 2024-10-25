# The Slide Programming Language
Slide is a relatively simple programming language the has the ability to slide in a variety
of emitters, allowing it to take advantage of a variety of performance and libraries from
other technologies.

## Todo
Lexer
- [x] Leftshift operator
- [x] Rightshift operator
- [ ] Map init
- [x] Switches
- [ ] +=, -=, etc
- [x] While

Parse
- [x] Leftshift operator
- [x] Rightshift operator
- [ ] Map init
- [x] Switches
- [ ] Increment and decrement can be entire assignment
- [ ] +=, -=, etc
- [x] While

Semantic Analysis
- [ ] Identifier exists
- [ ] Type exists
- [ ] Correct types used in expression
- [ ] Correct type returned in function
- [ ] Correct type assigned to variable
- [ ] Breaks only in loop
- [ ] Continue only in loop
- [ ] Returns only in functions
- [ ] Properties exist on struct
- [ ] Only dereference pointers
- [ ] Only index array
- [ ] Only access on structs
- [ ] Can't compare arraylists
- [ ] Check array index is int
- [ ] Check argument types match
- [ ] Check for out-of-bound on array
- [ ] Type inference
- [ ] Only nils for pointers
- [ ] Exhaustive switches
- [ ] Switches on same time

Tests
- [ ] Create proper tests
- [x] Linked lists
- [x] Pointer to structs
- [ ] Map of structs
- [ ] Array of structs
- [ ] Struct with array, arraylist, map, and pointer properties
- [ ] Map init
- [ ] Complete compiler test
- [ ] Switch test

Emitter
- [ ] Golang emitter
