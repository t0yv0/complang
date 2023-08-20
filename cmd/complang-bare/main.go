package main

import (
	"log"

	"github.com/t0yv0/complang/repl"
	"github.com/t0yv0/complang/value"
)

func main() {
	err := repl.ReadEvalPrintLoop(repl.ReadEvalPrintLoopOptions{
		HistoryFile: "/tmp/complang-bare-readline.history",
		InitialEnvironment: map[value.Symbol]value.Value{
			value.NewSymbol("$digits"): &value.MapValue{
				Value: map[value.Symbol]value.Value{
					value.NewSymbol("one"):   &value.StringValue{Value: "1"},
					value.NewSymbol("two"):   &value.StringValue{Value: "2"},
					value.NewSymbol("three"): &value.StringValue{Value: "3"},
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
