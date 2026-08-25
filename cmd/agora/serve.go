package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/rui-ktei/socrates/internal/document"
	"github.com/rui-ktei/socrates/internal/process"
	"github.com/rui-ktei/socrates/internal/refusal"
	"github.com/rui-ktei/socrates/internal/surface"
	"github.com/spf13/cobra"
)

const (
	servingVariable = "AGORA_SERVING"
	standardInput   = "-"

	reachInterval = 20 * time.Millisecond
	reachDeadline = 20 * time.Second
)

func newServeCommand() *cobra.Command {
	var from, address string

	command := &cobra.Command{
		Use:   "serve",
		Short: serveShort,
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if from == "" {
				return refusal.New(nameTheDocument)
			}

			raw, err := readDocument(command, from)
			if err != nil {
				return err
			}
			read, err := document.Read(raw)
			if err != nil {
				return err
			}
			if err := surface.CheckLoopback(address); err != nil {
				return err
			}

			if inTheForeground() {
				return hold(read, raw, address)
			}
			return start(command, raw, address)
		},
	}

	command.Flags().StringVar(&from, "document", "", documentUsage)
	command.Flags().StringVar(&address, "address", surface.DefaultAddress, addressUsage)
	return command
}

func inTheForeground() bool {
	return os.Getenv(servingVariable) != ""
}

func readDocument(command *cobra.Command, from string) ([]byte, error) {
	if from == standardInput {
		raw, err := io.ReadAll(command.InOrStdin())
		if err != nil {
			return nil, fmt.Errorf(readingTheDocument, err)
		}
		return raw, nil
	}

	raw, err := os.ReadFile(from)
	if err != nil {
		return nil, fmt.Errorf(readingTheDocument, err)
	}
	return raw, nil
}

func hold(read document.Document, raw []byte, address string) error {
	token, err := surface.MintToken()
	if err != nil {
		return err
	}

	server, err := surface.Serve(read, address, token, surface.IdlePeriod)
	if err != nil {
		return err
	}

	if err := surface.WriteDocument(raw); err != nil {
		return err
	}
	if err := surface.WriteSession(surface.Session{
		Address: server.Address(),
		Token:   token,
		Process: os.Getpid(),
	}); err != nil {
		return err
	}

	server.Wait()
	return surface.ForgetSession()
}

func start(command *cobra.Command, raw []byte, address string) error {
	if open, err := surface.ReadSession(); err == nil {
		reaching, cancel := context.WithTimeout(command.Context(), reachTimeout())
		defer cancel()

		caller := surface.Reach(open)
		if caller.Reachable(reaching) == nil {
			if err := caller.Revise(reaching, raw); err != nil {
				return err
			}
			fmt.Fprintln(command.OutOrStdout(), open.Address)
			return nil
		}
	}

	if err := surface.ForgetSession(); err != nil {
		return err
	}

	binary, err := os.Executable()
	if err != nil {
		return fmt.Errorf(readingThisBinary, err)
	}

	spawned := exec.Command(binary, "serve", "--document", standardInput, "--address", address)
	spawned.Env = append(os.Environ(), servingVariable+"=yes")
	spawned.Stdin = bytes.NewReader(raw)
	spawned.Stdout = nil
	spawned.Stderr = nil
	process.Detach(spawned)

	if err := spawned.Start(); err != nil {
		return fmt.Errorf(startingTheSurface, err)
	}
	go func() { _ = spawned.Wait() }()

	open, err := waitUntilReachable(command.Context())
	if err != nil {
		return err
	}
	fmt.Fprintln(command.OutOrStdout(), open.Address)
	return nil
}

func waitUntilReachable(ctx context.Context) (surface.Session, error) {
	waiting, cancel := context.WithTimeout(ctx, reachDeadline)
	defer cancel()

	for {
		if open, err := surface.ReadSession(); err == nil {
			reaching, done := context.WithTimeout(waiting, reachTimeout())
			err := surface.Reach(open).Reachable(reaching)
			done()
			if err == nil {
				return open, nil
			}
		}

		select {
		case <-waiting.Done():
			return surface.Session{}, fmt.Errorf(surfaceNeverAnswered, reachDeadline)
		case <-time.After(reachInterval):
		}
	}
}

func reachTimeout() time.Duration {
	return 2 * time.Second
}
