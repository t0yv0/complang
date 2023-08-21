# complang

A toy scripting language designed for:

- interactive shell-like use
- context-aware code completion
- being extensible in Go

## Example

See [complang-bare](./cmd/complang-bare/main.go) for an example Go-extended complang REPL.

```
go build ./cmd/complang-bare/
./complang-bare

» $digits
{one: "1", three: "3", two: "2"}

» $digits t<TAB>
three two

» $digits three
"3"

» $x = three
» $digits $x
"3"
```

## Semantics

Spaces denote object message send (inspired by Smalltalk), for example the following is bit like `foo.subfield()` in JS:

    » $foo subfield

The interesting bit around which complang is designed is context-aware completion. Completion is activated when pressing
a TAB at the end of an incomplete line:

    » $foo subfield bar<<TAB>>

The interpreter will parse the line as a `query`, and evaluate the `expr` part (`$foo subfield` above) to an object
value `v`, then dispatch the `symbol` part (`bar` above) to the value `v` to find completions specific to the object.

Since completions are often repeated, evaluation of expressions must be free of side-effects. Instead, expressions may
evaluate to descriptions of side-effects to be performed by the interpreter when submitted (think IO Monad in Haskell):

    » $foo subfield barbell<<ENTER>>

If `$foo subfield barbell` evaluates to the moral equivalent of `print("barbell")` effect, this will give:

    » $foo subfield barbell<<ENTER>>
    barbell

To make it ergonomic to work in shell-like contexts symbols are self-evaluating:

    » sym
    sym

To distinguish bound symbols the syntax requires sigils:

    » $foo
    fooValue

The binding form is as follows, and it simply modifies the global environment `Map[Symbol,Value]`:

    » $foo = fooValue

### Values

The space of values looks like this:

```
value
    simpleValue
    map[symbol,value]
    []value
    customValue

simpleValue
    null
    number
    bool
    string
    symbol
```

Custom values can be implemented in Go to overide key interactions with the interpreter:

```go
type CustomValue interface {
	Message(arg Value) Value
	CompleteSymbol(query Symbol) []Symbol
    Run() Value
	Show() string
}
```

Such values can have custom completion via `CompleteSymbol`.

Note that `Message` should not have side-effects. Instead, side-effects should be evaluated as part of `Run`.

## Syntax

### Expressions

```
expr
    simpleExpr
    expr simpleExpr

simpleExpr
    literal
    ref
    '(' expr ')'

literal
    null
    symbol
    bool
    string
    number

query
    expr symbol
    ref

stmt
    expr
    ref '=' expr
```

### Tokens

Borrowing lexical structure from the JSON grammar:

```
token
    symbol
    ref
    string
    number
    bool
    "null"
    '('
    ')'
    '='

symbol
    [_a-zA-Z][-_a-zA-Z0-9]*

ref
    [$] [-_a-zA-Z0-9]*

bool
    "true"
    "false"

number
    0
    [-]?[1-9][0-9]*

string
    ["] char* ["]

char
    [^"\\]
    "\" escape

escape
    '"'
    '\'
    '/'
    'b'
    'f'
    'n'
    'r'
    't'
```
