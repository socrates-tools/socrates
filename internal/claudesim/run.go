package claudesim

import (
	"io"

	"github.com/socrates-tools/socrates-review/internal/simulation"
)

const defaultRejection = "error: unknown option '--settings'"

type wire struct{}

func Run(arguments []string, input io.Reader, out, failures io.Writer, environment []string) int {
	return simulation.Run(wire{}, arguments, input, out, failures, environment)
}

func (wire) Read(arguments []string, invocation *simulation.Invocation) error {
	parsed, err := parseArguments(arguments)
	if err != nil {
		return err
	}

	invocation.Model = parsed.model
	invocation.Effort = parsed.effort
	invocation.Schema = parsed.schema
	invocation.SettingSources = parsed.settingSources
	invocation.Settings = parsed.settings
	invocation.Tools = parsed.tools
	return nil
}

func (wire) Write(out io.Writer, response simulation.Response, invocation simulation.Invocation) error {
	return writeStream(out, response, invocation)
}

func (wire) Rejection() string {
	return defaultRejection
}
