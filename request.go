package main

// Request type
type Type int

const (
	PathRequest Type = iota
	CompletionRequest
)

type Request struct {
	Type Type
	To   string
	From string
}
