package fidelity

import "github.com/socrates-devtool/socrates/internal/findings"

type Request struct {
	Name   string
	Prompt string
	Schema []byte
}

func Requests() []Request {
	return []Request{
		{
			Name:   "findings-review",
			Prompt: "Read parse.go in your working directory and report what you find as the schema requires. Report nothing you are not sure of.",
			Schema: []byte(findings.Schema),
		},
		{
			Name:   "plain-answer",
			Prompt: "Answer with the word ready and nothing else.",
		},
	}
}
