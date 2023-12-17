# complang

What if Go was interactive? Suppose you work in a Go-dominated domain and you have some data and
objects to explore interactively. If only you could drop them in a scripting REPL... Now if this
really describes you and you are looking for a mature solution, you may be interested in
[go-pry](https://github.com/d4l3k/go-pry). Complang is a toy-level exploration in this space.

Complang is a toy scripting language designed for:

- interactive shell-like use
- context-aware code completion
- being extensible in Go by binding Go values to it

## Example

See [complang-bare](./cmd/complang-bare/main.go) for an example Go-extended complang REPL.

```
go build ./cmd/complang-bare/
./complang-bare

» $digits
one:
    text: "1"
three:
    text: "3"
two:
    text: "2"

» $digits t<TAB>
two three

» $digits three
"3"

» $x = three
» $digits $x
"3"

» [$x | $digits $x]
<Closure:$x>

» [$x | $digits $x] one
"1"

» [$x $y | $y $x] three $digits
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
    closure

simpleValue
    null
    number
    bool
    string
    symbol
```

Custom values can be implemented in Go to override key interactions with the interpreter:

```go
type Value interface {
	Message(context.Context, Value) Value
}
```

The are several special messages that encode interactions with the interpreter.

- ShowMessage asks the object to produce a StringValue to display itself in the REPL

- RunMessage asks the object to run any deferred side-effects and return the final value

- CompleteRequest queries which messages the object supports responding to

Note that `Message` evaluation should not have side-effects except when responding to the
RunMessage. This helps the REPL perform side-effect free dynamic completion while avoiding
side-effects until you press enter.

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
    lambdaBlockExpr

lambdaBlockExpr
    [ expr* ]
    [ symbol* | expr* ]

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
    '['
    ']'
    '|'

symbol
    [_a-zA-Z][-_a-zA-Z0-9:/]*

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
