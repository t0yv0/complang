package main

import (
	"context"
	"fmt"
	"log"

	cl "github.com/t0yv0/complang"
	"github.com/t0yv0/complang/repl"
)

func main() {
	ctx := context.Background()
	err := repl.ReadEvalPrintLoop(ctx, repl.ReadEvalPrintLoopOptions{
		HistoryFile: "/tmp/complang-bare-readline.history",
		InitialEnvironment: map[string]cl.Value{
			"$digits": cl.MapValue(map[string]cl.Value{
				"one":   cl.StringValue{Text: "1"},
				"two":   cl.StringValue{Text: "2"},
				"three": cl.StringValue{Text: "3"},
			}),
			"$something": cl.BindValue(map[string]interface{}{
				"string": "stringValue",
				"bool":   true,
				"slice":  []string{"a", "b", "c"},
				"map":    map[string]string{"one": "1"},
				"structure": exampleStruct{
					X: true,
					Y: "ok",
				},
			}),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

type exampleStruct struct {
	X bool
	Y string
}

func (s exampleStruct) String() string {
	return fmt.Sprintf("exampleStruct{X:%v,Y:%q}", s.X, s.Y)
}

func (s exampleStruct) WithArg(arg string) string {
	return s.Y + arg
}

func (s exampleStruct) With2Args(arg1, arg2 string) []string {
	return []string{s.Y, arg1, arg2}
}

func (s exampleStruct) Print() {
	fmt.Println(s.String())
}
