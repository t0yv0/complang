package parser

type token struct {
	t      any
	offset int
	length int
}

type symbol string
