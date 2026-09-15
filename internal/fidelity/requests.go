package fidelity

import (
	"github.com/socrates-tools/socrates/internal/agent"
	"github.com/socrates-tools/socrates/internal/findings"
)

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

func Asked(vendor agent.Vendor, request Request) string {
	if vendor.Tier != agent.Parsed || len(request.Schema) == 0 {
		return request.Prompt
	}
	return request.Prompt + theSchemaItMustSatisfy + string(request.Schema) + "\n\n" + findings.DeliveryInstruction()
}
