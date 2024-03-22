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

#### Literal evaluators
Several literal evaluators are provided out of the box in `util.go`. To use them, just add them to the language.
Make sure to add your evaluators from narrow to wide match; the statement will be matched to the evaluators in order, 
the first one that matches is used.