package surface

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/socrates-devtool/socrates/internal/refusal"
)

const reachTimeout = 5 * time.Second

type Caller struct {
	open    Record
	reading *http.Client
}

func Reach(open Record) *Caller {
	return &Caller{open: open, reading: &http.Client{}}
}

func (c *Caller) Reachable(ctx context.Context) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.at("/"), nil)
	if err != nil {
		return err
	}
	request.Header.Set(TokenHeader, c.open.Token)

	answered, err := c.reading.Do(request)
	if err != nil {
		return err
	}
	defer answered.Body.Close()
	_, _ = io.Copy(io.Discard, answered.Body)

	if answered.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf(surfaceAnsweredWith, answered.StatusCode)
	}
	if answered.Header.Get(SurfaceHeader) != c.open.Address {
		return fmt.Errorf(addressIsNotTheSurface, c.open.Address)
	}
	return nil
}

func (c *Caller) Open(ctx context.Context, raw []byte) (Opened, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.at(SessionPrefix), bytes.NewReader(raw))
	if err != nil {
		return Opened{}, err
	}
	request.Header.Set(TokenHeader, c.open.Token)
	request.Header.Set("content-type", "application/json")

	answered, err := c.reading.Do(request)
	if err != nil {
		return Opened{}, fmt.Errorf(reachingTheSurface, err)
	}
	defer answered.Body.Close()

	if answered.StatusCode >= http.StatusBadRequest {
		said, _ := io.ReadAll(io.LimitReader(answered.Body, BodyCeiling))
		return Opened{}, fmt.Errorf(surfaceRefused, answered.StatusCode, bytes.TrimSpace(said))
	}

	var opened Opened
	if err := json.NewDecoder(answered.Body).Decode(&opened); err != nil {
		return Opened{}, fmt.Errorf(openingTheSession, err)
	}
	return opened, nil
}

func (c *Caller) Await(ctx context.Context, session string, within time.Duration) (Answered, bool, error) {
	waiting := &http.Client{Timeout: within + reachTimeout}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.under(session, "/answer?within="+within.String()), nil)
	if err != nil {
		return Answered{}, false, err
	}
	request.Header.Set(TokenHeader, c.open.Token)

	answered, err := waiting.Do(request)
	if err != nil {
		return Answered{}, false, fmt.Errorf(reachingTheSurface, err)
	}
	defer answered.Body.Close()

	if answered.StatusCode == http.StatusNotFound {
		return Answered{}, false, refusal.New(surfaceHoldsNoSuchSession, session)
	}
	if answered.StatusCode == http.StatusNoContent {
		return Answered{}, false, nil
	}
	if answered.StatusCode != http.StatusOK {
		return Answered{}, false, fmt.Errorf(surfaceAnsweredWith, answered.StatusCode)
	}

	var read Answered
	if err := json.NewDecoder(answered.Body).Decode(&read); err != nil {
		return Answered{}, false, fmt.Errorf(readingTheAnswer, err)
	}
	return read, true, nil
}

func (c *Caller) Revise(ctx context.Context, session string, raw []byte) error {
	return c.post(ctx, session, c.under(session, "/document"), raw)
}

func (c *Caller) End(ctx context.Context, session string) error {
	return c.post(ctx, session, c.under(session, "/close"), nil)
}

func (c *Caller) Stop(ctx context.Context) error {
	return c.post(ctx, "", c.at("/stop"), nil)
}

func (c *Caller) post(ctx context.Context, session, address string, body []byte) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, address, bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set(TokenHeader, c.open.Token)
	request.Header.Set("content-type", "application/json")

	answered, err := c.reading.Do(request)
	if err != nil {
		return fmt.Errorf(reachingTheSurface, err)
	}
	defer answered.Body.Close()

	if answered.StatusCode == http.StatusNotFound && session != "" {
		return refusal.New(surfaceHoldsNoSuchSession, session)
	}
	if answered.StatusCode >= http.StatusBadRequest {
		said, _ := io.ReadAll(io.LimitReader(answered.Body, BodyCeiling))
		return fmt.Errorf(surfaceRefused, answered.StatusCode, bytes.TrimSpace(said))
	}
	return nil
}

func (c *Caller) at(route string) string {
	return "http://" + c.open.Address + route
}

func (c *Caller) under(session, route string) string {
	return c.at(SessionPrefix + session + route)
}
