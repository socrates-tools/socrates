package claudesim

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

type flag struct {
	names      []string
	takesValue bool
}

var knownFlags = []flag{
	{names: []string{"-p", "--print"}},
	{names: []string{"--verbose"}},
	{names: []string{"--dangerously-skip-permissions"}},
	{names: []string{"--fork-session"}},
	{names: []string{"--output-format"}, takesValue: true},
	{names: []string{"--input-format"}, takesValue: true},
	{names: []string{"--json-schema"}, takesValue: true},
	{names: []string{"--model"}, takesValue: true},
	{names: []string{"--effort"}, takesValue: true},
	{names: []string{"--setting-sources"}, takesValue: true},
	{names: []string{"--settings"}, takesValue: true},
	{names: []string{"--tools"}, takesValue: true},
	{names: []string{"--mcp-config"}, takesValue: true},
	{names: []string{"--strict-mcp-config"}},
	{names: []string{"--allowedTools", "--allowed-tools"}, takesValue: true},
	{names: []string{"--disallowedTools", "--disallowed-tools"}, takesValue: true},
	{names: []string{"--permission-mode"}, takesValue: true},
	{names: []string{"--session-id"}, takesValue: true},
	{names: []string{"-r", "--resume"}, takesValue: true},
	{names: []string{"--add-dir"}, takesValue: true},
	{names: []string{"--system-prompt"}, takesValue: true},
	{names: []string{"--append-system-prompt"}, takesValue: true},
	{names: []string{"--max-budget-usd"}, takesValue: true},
}

var outputFormats = []string{"text", "json", "stream-json"}

var efforts = []string{"low", "medium", "high", "xhigh", "max"}

var settingSources = []string{"user", "project", "local"}

type request struct {
	printMode      bool
	verbose        bool
	outputFormat   string
	schema         string
	model          string
	effort         string
	settingSources string
	settings       string
	tools          string
	mcpServers     string
	strictMCP      bool
}

func parseArguments(arguments []string) (request, error) {
	parsed := request{}
	values := map[string]string{}

	for position := 0; position < len(arguments); position++ {
		argument := arguments[position]

		name, inlineValue, inline := strings.Cut(argument, "=")
		if !strings.HasPrefix(argument, "-") {
			return request{}, fmt.Errorf("claudesim reads the prompt from standard input, and was given the argument %q", argument)
		}

		known, found := lookupFlag(name)
		if !found {
			return request{}, fmt.Errorf("unknown option '%s'", name)
		}

		if !known.takesValue {
			if inline {
				return request{}, fmt.Errorf("option '%s' does not take a value", name)
			}
			values[known.names[0]] = ""
			continue
		}

		value := inlineValue
		if !inline {
			if position+1 >= len(arguments) {
				return request{}, fmt.Errorf("option '%s' argument missing", name)
			}
			position++
			value = arguments[position]
		}
		values[known.names[0]] = value
	}

	_, parsed.printMode = values["-p"]
	_, parsed.verbose = values["--verbose"]
	parsed.outputFormat = values["--output-format"]
	parsed.schema = values["--json-schema"]
	parsed.model = values["--model"]
	parsed.effort = values["--effort"]
	parsed.settingSources = values["--setting-sources"]
	parsed.settings = values["--settings"]
	parsed.tools = values["--tools"]
	parsed.mcpServers = values["--mcp-config"]
	_, parsed.strictMCP = values["--strict-mcp-config"]

	if err := validate(parsed); err != nil {
		return request{}, err
	}
	return parsed, nil
}

func validate(parsed request) error {
	if !parsed.printMode {
		return fmt.Errorf("claudesim answers print mode only, and -p was not given")
	}
	if parsed.outputFormat == "" {
		return fmt.Errorf("claudesim answers a streaming stream, and --output-format was not given")
	}
	if !slices.Contains(outputFormats, parsed.outputFormat) {
		return fmt.Errorf("option '--output-format <format>' argument '%s' is invalid. Allowed choices are %s.", parsed.outputFormat, strings.Join(outputFormats, ", "))
	}
	if parsed.outputFormat != "stream-json" {
		return fmt.Errorf("claudesim speaks stream-json only, and was asked for %s", parsed.outputFormat)
	}
	if parsed.effort != "" && !slices.Contains(efforts, parsed.effort) {
		return fmt.Errorf("option '--effort <level>' argument '%s' is invalid. Allowed choices are %s.", parsed.effort, strings.Join(efforts, ", "))
	}
	for _, source := range strings.Split(parsed.settingSources, ",") {
		trimmed := strings.TrimSpace(source)
		if trimmed == "" {
			continue
		}
		if !slices.Contains(settingSources, trimmed) {
			return fmt.Errorf("option '--setting-sources <sources>' argument '%s' is invalid. Allowed choices are %s.", parsed.settingSources, strings.Join(settingSources, ", "))
		}
	}
	if parsed.schema != "" && !json.Valid([]byte(parsed.schema)) {
		return fmt.Errorf("option '--json-schema <schema>' argument is not JSON")
	}
	return nil
}

func lookupFlag(name string) (flag, bool) {
	for _, known := range knownFlags {
		if slices.Contains(known.names, name) {
			return known, true
		}
	}
	return flag{}, false
}
