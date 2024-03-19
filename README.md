## PArsed LAnguage

A simple lexer/parser that generates a runnable program.

The grammar of a program has the form:
```
Program    <- Expression ( '\n' Expression )* '\n'?
Expression <- Assignment | Operation | Comment
Assignment <- Variable ( Variable | List | Literal | Operation )
Operation  <- Literal ( List | Variable | Literal )*
List       <- '[' ( Variable | Literal )+ ']'
Variable   <- '$.+'
Comment    <- '#.+'
Literal    <- '.+'
```

Optionally, the operands of an operation may be wrapped in parentheses `()` to allow them to be on multiple lines.

While the grammar of the language is fixed, the operators and literals are bound dynamically. Types are however still
enforced.

See the comments in `pala_test.go` for an example on how to use.
