package claudesim

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	exitSuccess = 0
	exitUsage   = 1
)

func Run(arguments []string, input io.Reader, out, failures io.Writer, environment []string) int {
	parsed, err := parseArguments(arguments)
	if err != nil {
		fmt.Fprintf(failures, "error: %s\n", err)
		return exitUsage
	}

	prompt, err := io.ReadAll(input)
	if err != nil {
		fmt.Fprintf(failures, "error: reading the prompt: %s\n", err)
		return exitUsage
	}

	directory, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(failures, "error: locating the working directory: %s\n", err)
		return exitUsage
	}

	invocation := Invocation{
		Process:        os.Getpid(),
		Arguments:      arguments,
		Prompt:         string(prompt),
		Directory:      directory,
		Environment:    recordedEnvironment(environment),
		Model:          parsed.model,
		Effort:         parsed.effort,
		Schema:         parsed.schema,
		SettingSources: parsed.settingSources,
		Tools:          parsed.tools,
		MCPServers:     parsed.mcpServers,
		StrictMCP:      parsed.strictMCP,
	}

	scenario, err := LoadScenario(os.Getenv(ScenarioVariable))
	if err != nil {
		fmt.Fprintf(failures, "error: %s\n", err)
		return exitUsage
	}
	response := scenario.Match(invocation.Prompt)

	invocation.WritesRefused = attemptWrites(directory, response)

	child, err := startChild(response)
	if err != nil {
		fmt.Fprintf(failures, "error: %s\n", err)
		return exitUsage
	}
	invocation.Child = child

	if err := Append(os.Getenv(LogVariable), invocation); err != nil {
		fmt.Fprintf(failures, "error: %s\n", err)
		return exitUsage
	}

	if response.Behaviour == RefuseStrictness && parsed.strictMCP {
		fmt.Fprintln(failures, strictnessRefused)
		return exitUsage
	}

	if response.Behaviour == SayNothingAtAll {
		time.Sleep(response.hold())
		return exitSuccess
	}

	if err := writeStream(out, response, invocation); err != nil {
		fmt.Fprintf(failures, "error: %s\n", err)
		return exitUsage
	}

	if response.Behaviour == ExitNonZero {
		fmt.Fprintln(failures, response.Stderr)
		return response.Exit
	}

	return exitSuccess
}

func attemptWrites(directory string, response Response) []string {
	var refused []string
	for _, write := range response.Writes {
		if err := attemptWrite(directory, write); err != nil {
			refused = append(refused, fmt.Sprintf("%s: %s", write.Path, err))
		}
	}
	return refused
}

func attemptWrite(directory string, write Write) error {
	path := filepath.Join(directory, write.Path)

	if write.DefeatsTheSeal {
		if err := clearTheWayTo(directory, path); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating the directory for %s: %w", write.Path, err)
	}
	if err := os.WriteFile(path, []byte(write.Content), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", write.Path, err)
	}
	return nil
}

func clearTheWayTo(directory, path string) error {
	for at := filepath.Dir(path); strings.HasPrefix(at, directory); at = filepath.Dir(at) {
		if err := makeWritable(at); err != nil {
			return err
		}
	}
	return makeWritable(path)
}

func makeWritable(path string) error {
	known, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	if err := os.Chmod(path, known.Mode().Perm()|0o200); err != nil {
		return fmt.Errorf("clearing the way to %s: %w", path, err)
	}
	return nil
}

func startChild(response Response) (int, error) {
	if len(response.ChildCommand) == 0 {
		return 0, nil
	}

	command := exec.Command(response.ChildCommand[0], response.ChildCommand[1:]...)
	if err := command.Start(); err != nil {
		return 0, fmt.Errorf("starting the child %v: %w", response.ChildCommand, err)
	}
	return command.Process.Pid, nil
}
