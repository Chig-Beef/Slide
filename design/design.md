# Slide
## Static typing
The primitive types will be:
1. int (and variants)
2. bool
3. string
4. Floating points
5. Char

Complex data types will be:
1. Array
2. Slice
3. Map
4. Structs

## Arrays and slices
```
x []byte = [1, 2, 3+4];
y [3]string = ["a", "b", "c"];
z [2]^float = [`1.1, `2.2];

// Prefixing with the index makes it easier to differentiate between other
// cases, however, this might be changed in future
[0]x = 7;

```

## Maps
```
// Creates and initialises a map
y [string]int;
["a"]y = 3;
```

## Type inference
Variables can be declared with the type, or just infered.
This does not require the walrus operator.
Start off without inference for the time being
Logically search through operations to understand a type?

Sometimes it is also necessary to create a type alias.
```
// Alias is now an int
typedef Alias int;
```

## Variable definitions
```
x = 0; // Will set x to the generic int type

x byte = 0; // Will set x to the byte type (uint8)

x = -1; // Will set x to the generic sint type

x = "hello"; // Will set x to string type
x = 2; // Error, incompatible types

```

## Numeric types
Using compile time checks, the generics will use the smallest possible integer
or conversion to make sure everythign works out

`num` covers all integers and floats

### Integer
`int` covers all integers

#### Unsigned
`uint` covers all unsigned integers
`byte` or `uint8` is an unsigned 8-bit integer
`word` or `uint16` is an unsigned 16-bit integer
`dword` or `uint32` is an unsigned 32-bit integer
`qword` or `uint64` is an unsigned 64-bit integer

#### Signed
`sint` covers all signed integers
`int8` is a signed 8-bit integer
`int16` is a signed 16-bit integer
`int32` is a signed 32-bit integer
`int64` is a signed 64-bit integer

### Float
`float` covers all floats
`float32` is a 32-bit float
`float64` or `double` is a 64-bit float

## Function definitions
```
fun sum(a int, b int) int {
    return a + b;
}
```

## Errors as values
Errors have a string and a code to make them easier to track.
The code will contain flags, such as whether the giver of the error thinks it
is recoverable, and so forth.
Actually I think we should do error wrapping?
This is probably a better way to maintain a single return value system.
There's probably extra syntax or so on needed, but we can leave that for now.

## Enums
```
enum Days {
    monday,
    tuesday,
    wednesday,
    thursday,
    friday,
    saturday,
    sunday,
}
```

## Type conversions
Type conversions are single functions
May be easier to have as methods, then the creator of types and structs can
write their own.
For primitives, each type is prefixed with "to", so it would be tostring.

## Compilation
Lexing is entirely completed before moving on.
Parsing is entirely completed before moving on.
Semantic analysis is entirely completed before moving on.
Optimizaion is entirely completed before moving on.
Emitting to Go? C? NASM?
Maybe start with Go, however, since the logic of emission should entirely be in
the emitter, this should be easy to replace with a different language.
Correctly done, lexing could be done on another thread, multi-threaded
compilation could be a huge time save.

## String manipulation
Strings can be concatonated to each other using the + operator.
Strings can be multiplied by a non-negative integer? (however this could return an error?)

## Ifs
```
if cond {

} elif cond {

} else {

}
```

## Fors
```
for x = 2; cond; x = 7 {

}

// Will figure out as an iterator
range 7 {

}

forever {
    continue;
    break; // Shorthand for break 0;
    return;
}

forever {
    forever {
        break 1;
    }
}
```

## Structs
```
struct Point {
    x int
    y int
    id string
}
```

## Pointers
` is a reference (value prime)
^ is a dereference (pull up the value)
Using these symbols that aren't used for anythign else will increase
speed as it avoids ambiguity
```
x int^ = `7;
```

## Pattern matching

## Interfaces

## Nil

## Multiple files
Since functions, structs, types, and methods are hoisted, there is no need to
worry about importing files like python. Order doesn't matter, as as there is a
single main function.
On another note, no regular statements are usable outside of the main function.
Struct definitions, function definitions, enum definitions, typedefs, can be outside main only.

## Packages

## In-built functions

## Other-language interop

## Calling functions
```
call println(6);
x int = add(2, 3);
```

## Arithmatic
```
x = 1 + 2

y = 3 * 7

z = 2 / 2

a = 5 - 1

z = 5 | 3

z = 1 & 3
```

Do everything

## Standard Library
The standard library is an interesting problem considering the ability to slide
in different emitters.
Because of this, there's almost a need for multiple implementations of it.
This isn't fully true, as there are parts of the standard that can be fully
implemented in slide.
For example, some basic mathematics can be implemented in slide, also array manipulation.
However, for example, accessing the network will need access to the emitters own standard.

Slide standard:
1. [ ] math
2. [ ] slices
3. [ ] strconv
4. [ ] strings
