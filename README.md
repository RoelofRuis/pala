## PArsed LAnguage

A simple lexer/parser that generates an interpretable AST.

It allows for dynamically binding operation and literal parsers.

The grammar of a program has the form:
```
Program    <- Expression ( '\n' Expression )* '\n'?
Expression <- Variable ( Variable | List | Literal | Operator ) | Operator
Operation  <- Literal ( List | Variable | Literal )*
List       <- '[' ( Variable | Literal )+ ']'
Variable   <- '$.+'
Comment    <- '#.+
Literal    <- '.+'
```

See the comments in `parser_test.go` for an example on how to use.
