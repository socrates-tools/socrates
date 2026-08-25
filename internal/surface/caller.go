package surface

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const reachTimeout = 5 * time.Second

type Caller struct {
	open    Session
	reading *http.Client
}

func Reach(open Session) *Caller {
	return &Caller{open: open, reading: &http.Client{}}
}

func (c *Caller) Reachable(ctx context.Context) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.at("/"), nil)
	if err != nil {
		return err
	}

	answered, err := c.reading.Do(request)
	if err != nil {
		return err
	}
	defer answered.Body.Close()
	_, _ = io.Copy(io.Discard, answered.Body)

	if answered.StatusCode != http.StatusOK {
		return fmt.Errorf(surfaceAnsweredWith, answered.StatusCode)
	}
	return nil
}

func (c *Caller) Await(ctx context.Context, within time.Duration) (Answered, bool, error) {
	waiting := &http.Client{Timeout: within + reachTimeout}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.at("/answer?within="+within.String()), nil)
	if err != nil {
		return Answered{}, false, err
	}
	request.Header.Set(TokenHeader, c.open.Token)

	answered, err := waiting.Do(request)
	if err != nil {
		return Answered{}, false, fmt.Errorf(reachingTheSurface, err)
	}
	defer answered.Body.Close()

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

func (c *Caller) Revise(ctx context.Context, raw []byte) error {
	return c.post(ctx, "/document", raw)
}

func (c *Caller) End(ctx context.Context) error {
	return c.post(ctx, "/close", nil)
}

func (c *Caller) Stop(ctx context.Context) error {
	return c.post(ctx, "/stop", nil)
}

func (c *Caller) post(ctx context.Context, route string, body []byte) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.at(route), bytes.NewReader(body))
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

	if answered.StatusCode >= http.StatusBadRequest {
		said, _ := io.ReadAll(io.LimitReader(answered.Body, bodyCeiling))
		return fmt.Errorf(surfaceRefused, answered.StatusCode, bytes.TrimSpace(said))
	}
	return nil
}

func (c *Caller) at(route string) string {
	return "http://" + c.open.Address + route
}
