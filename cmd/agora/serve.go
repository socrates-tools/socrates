package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/socrates-devtool/socrates/internal/config"
	"github.com/socrates-devtool/socrates/internal/document"
	"github.com/socrates-devtool/socrates/internal/process"
	"github.com/socrates-devtool/socrates/internal/refusal"
	"github.com/socrates-devtool/socrates/internal/surface"
	"github.com/spf13/cobra"
)

const (
	servingVariable = "AGORA_SERVING"
	standardInput   = "-"
	addressFlag     = "address"
	sessionFlag     = "session"

	reachInterval = 20 * time.Millisecond
	reachDeadline = 20 * time.Second
)

func newServeCommand() *cobra.Command {
	var from, address, session string

	command := &cobra.Command{
		Use:   "serve",
		Short: serveShort,
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			read, err := config.Load()
			if err != nil {
				return err
			}
			configured := read.Surface
			if !command.Flags().Changed(addressFlag) && configured.Bind != "" {
				address = configured.Bind
			}
			if aPortWasCommittedTo(command, configured) {
				if err := surface.CheckLoopback(address); err != nil {
					return err
				}
			} else {
				address = ""
			}
			if configured.ReachedAt != "" {
				if err := surface.CheckOpenable(configured.ReachedAt); err != nil {
					return err
				}
			}

			if inTheForeground() {
				return hold(address)
			}

			if from == "" {
				return refusal.New(nameTheDocument)
			}
			raw, err := readDocument(command, from)
			if err != nil {
				return err
			}
			if _, err := document.Read(raw); err != nil {
				return err
			}
			file, err := config.Path()
			if err != nil {
				return err
			}
			return start(command, raw, address, stated{reachedAt: configured.ReachedAt, file: file}, session)
		},
	}

	command.Flags().StringVar(&from, "document", "", documentUsage)
	command.Flags().StringVar(&address, addressFlag, surface.DefaultAddress, addressUsage)
	command.Flags().StringVar(&session, sessionFlag, "", sessionUsage)
	return command
}

func inTheForeground() bool {
	return os.Getenv(servingVariable) != ""
}

func aPortWasCommittedTo(command *cobra.Command, configured config.Surface) bool {
	return command.Flags().Changed(addressFlag) || configured.Bind != "" || configured.ReachedAt != ""
}

func whereToBind(committed string) (string, error) {
	if committed != "" {
		if err := surface.CheckFree(committed); err != nil {
			return "", err
		}
		return committed, nil
	}
	return surface.FirstFreeIn(surface.WindowAddresses())
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

func hold(address string) error {
	server, err := surface.Serve(address, surface.IdlePeriod)
	if err != nil {
		return err
	}

	if err := surface.WriteRecord(surface.Record{
		Address: server.Address(),
		Token:   server.Control(),
		Process: os.Getpid(),
	}); err != nil {
		return err
	}

	server.Wait()
	return surface.ForgetRecord()
}

func start(command *cobra.Command, raw []byte, address string, said stated, session string) error {
	if open, err := surface.ReadRecord(); err == nil {
		reaching, cancel := context.WithTimeout(command.Context(), reachTimeout())
		defer cancel()

		caller := surface.Reach(open)
		if caller.Reachable(reaching) == nil {
			named, err := onTheSurfaceStanding(reaching, caller, raw, session)
			if err != nil {
				return err
			}
			handOver(command, open, said, named)
			return nil
		}
	}

	if session != "" {
		return refusal.New(noSurfaceHoldsThatSession, session)
	}
	if err := surface.ForgetRecord(); err != nil {
		return err
	}
	binding, err := whereToBind(address)
	if err != nil {
		return err
	}

	binary, err := os.Executable()
	if err != nil {
		return fmt.Errorf(readingThisBinary, err)
	}

	spawned := exec.Command(binary, "serve", "--"+addressFlag, binding)
	spawned.Env = append(os.Environ(), servingVariable+"=yes")
	spawned.Stdin = nil
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

	opening, cancel := context.WithTimeout(command.Context(), reachTimeout())
	defer cancel()

	opened, err := surface.Reach(open).Open(opening, raw)
	if err != nil {
		return err
	}
	handOver(command, open, said, opened.Session)
	return nil
}

type stated struct {
	reachedAt string
	file      string
}

func handOver(command *cobra.Command, open surface.Record, said stated, session string) {
	fmt.Fprintln(command.OutOrStdout(), whereThePageIs(open, said.reachedAt, session))
	if said.reachedAt != "" {
		return
	}
	fmt.Fprintf(command.ErrOrStderr(), pageOpensOnThisMachineOnly+"\n", open.Address, surface.ReachedAtField, said.file)
}

func onTheSurfaceStanding(ctx context.Context, caller *surface.Caller, raw []byte, session string) (string, error) {
	if session == "" {
		opened, err := caller.Open(ctx, raw)
		if err != nil {
			return "", err
		}
		return opened.Session, nil
	}
	if err := caller.Revise(ctx, session, raw); err != nil {
		return "", err
	}
	return session, nil
}

func whereThePageIs(open surface.Record, reachedAt, session string) string {
	reached := reachedAt
	if reached == "" {
		reached = "http://" + open.Address
	}
	return strings.TrimSuffix(reached, "/") + surface.SessionPrefix + session
}

func waitUntilReachable(ctx context.Context) (surface.Record, error) {
	waiting, cancel := context.WithTimeout(ctx, reachDeadline)
	defer cancel()

	for {
		if open, err := surface.ReadRecord(); err == nil {
			reaching, done := context.WithTimeout(waiting, reachTimeout())
			err := surface.Reach(open).Reachable(reaching)
			done()
			if err == nil {
				return open, nil
			}
		}

		select {
		case <-waiting.Done():
			return surface.Record{}, fmt.Errorf(surfaceNeverAnswered, reachDeadline)
		case <-time.After(reachInterval):
		}
	}
}

func reachTimeout() time.Duration {
	return 2 * time.Second
}
