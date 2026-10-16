package surface

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/internal/state"
	"github.com/socrates-tools/socrates/internal/version"
)

const (
	IdentityRoute = "/identity"
	IdentityTool  = "agora"

	identityCeiling = 4 << 10
	askTimeout      = 2 * time.Second
)

type Identity struct {
	Tool        string `json:"tool"`
	Root        string `json:"root"`
	Version     string `json:"version"`
	Commit      string `json:"commit"`
	Uncommitted bool   `json:"uncommitted"`
}

func (i Identity) SameBuild(other Identity) bool {
	return i.Version == other.Version && i.Commit == other.Commit && i.Uncommitted == other.Uncommitted
}

func thisSurface() (Identity, error) {
	root, err := state.Root()
	if err != nil {
		return Identity{}, err
	}
	built := version.Linked(IdentityTool)
	return Identity{
		Tool:        IdentityTool,
		Root:        root,
		Version:     built.Version,
		Commit:      built.Commit,
		Uncommitted: built.Uncommitted,
	}, nil
}

func Ask(address string) *Caller {
	return Reach(Record{Address: address})
}

func (c *Caller) Identify(ctx context.Context) (Identity, error) {
	asking, cancel := context.WithTimeout(ctx, askTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(asking, http.MethodGet, c.at(IdentityRoute), nil)
	if err != nil {
		return Identity{}, err
	}

	answered, err := c.reading.Do(request)
	if err != nil {
		return Identity{}, fmt.Errorf(reachingTheSurface, err)
	}
	defer func() { _ = answered.Body.Close() }()

	if answered.StatusCode != http.StatusOK {
		return Identity{}, fmt.Errorf(surfaceAnsweredWith, answered.StatusCode)
	}

	var told Identity
	decoder := json.NewDecoder(io.LimitReader(answered.Body, identityCeiling))
	if err := decoder.Decode(&told); err != nil {
		return Identity{}, fmt.Errorf(readingTheAnswer, err)
	}
	return told, nil
}

type Whereabouts struct {
	Standing    bool
	Recorded    bool
	Record      Record
	AnotherRoot string
	AnotherAt   string

	OfAnotherBuild bool
}

func Locate(ctx context.Context, committed string) (Whereabouts, error) {
	own, err := thisSurface()
	if err != nil {
		return Whereabouts{}, err
	}

	var found Whereabouts
	asked := ""
	if port := portToAsk(committed); port != "" {
		asked = port
		if told, err := Ask(port).Identify(ctx); err == nil && told.Tool == IdentityTool {
			if told.Root == own.Root {
				return standingAt(port, told, own), nil
			}
			found.AnotherRoot, found.AnotherAt = told.Root, port
		}
	}

	if recorded, err := ReadRecord(); err == nil && recorded.Address != asked {
		if told, err := Ask(recorded.Address).Identify(ctx); err == nil && told.Tool == IdentityTool && told.Root == own.Root {
			found.Standing, found.Recorded, found.Record = true, true, recorded
			found.OfAnotherBuild = !own.SameBuild(told)
		}
	}
	return found, nil
}

func portToAsk(committed string) string {
	if committed == "" {
		return DefaultAddress()
	}
	if _, port, err := net.SplitHostPort(committed); err != nil || port == "0" {
		return ""
	}
	return committed
}

func standingAt(address string, told, own Identity) Whereabouts {
	of := Whereabouts{Standing: true, OfAnotherBuild: !own.SameBuild(told)}
	recorded, err := ReadRecord()
	if err != nil {
		return of
	}
	recorded.Address = address
	of.Recorded, of.Record = true, recorded
	return of
}

func (w Whereabouts) Reached() (Record, error) {
	if !w.Standing {
		return Record{}, errors.New(noSurfaceIsOpen)
	}
	if !w.Recorded {
		return Record{}, errors.New(surfaceStandsWithNoRecord)
	}
	return w.Record, nil
}

func CheckNotHeldByAnotherRoot(found Whereabouts, file string) error {
	if found.AnotherRoot == "" {
		return nil
	}
	return refusal.New(portIsHeldByAnotherRoot, found.AnotherRoot, found.AnotherAt, BindField, file, AddressOption)
}
