## PArsed LAnguage

A simple lexer/parser that generates a runnable program.

It allows for dynamically binding operation and literal parsers.

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

See the comments in `pala_test.go` for an example on how to use.
