# *p*ipe*l*ine

[![test](https://github.com/lesomnus/pl/actions/workflows/test.yaml/badge.svg)](https://github.com/lesomnus/pl/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/lesomnus/pl)](https://goreportcard.com/report/github.com/lesomnus/pl)
[![codecov](https://codecov.io/gh/lesomnus/pl/branch/main/graph/badge.svg?token=PJ3sRR1Ms0)](https://codecov.io/gh/lesomnus/pl)

Single line expression inspired by pipeline in `text/template`. A pipeline is a sequence of functions separated by `|`. Functions can take arguments, and the result of the previous function is passed to the last argument of the next function. The first word of a pipeline element is the name of the function, and the following words become the function's arguments. Pipeline can be nested by wrapping them with `(...)` in argument position.

## Usage

```go
executor := pl.NewExecutor()
executor.Funcs["sum"] = func(vs ...int) int {
	rst := 0
	for _, v := range vs {
		rst += v
	}

	return rst
}
executor.Props.Set(pl.K("answer"), 42)

rst, err := executor.ExecuteExpr("(sum 1 2 (sum 3 | sum (sum $.answer 5) 6) 7 (sum 8) | sum 9 10)")
if err != nil {
	panic(err)
}

// v == 93
v, ok := rst.(int)
if !ok {
	panic("expected v to be of type int.")
} else if v != 93 {
	panic("expected v to be 93")
}
```


## Syntax

```ebnf
pipeline = '(', function, { '|', function }, ')';
function = name, { { ' ' }*, argument };
name     = identifier;
argument = string | number | reference | pipeline;

identifier = letter { letter | digit | '_' }*;
string     = '"', ? printable characters ?, '"';
number     = integer | floating_point;
reference  = '$', { reference_part }*;

integer        = [ '-' | '+' ], { digit }*;
floating_point = integer, [ '.', { digit }* ];
reference_part = '[', integer, ']' | '.', identifier;

letter = /[a-zA-Z]/;
digit  = /[0-9]/;
```

