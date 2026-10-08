package surface_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/internal/state"
)

func TestARequestWithNoTokenIsToldWhoTheSurfaceIsAndNothingItCouldActOn(t *testing.T) {
	// scenario: The port is asked who stands there
	held := serving(t, twoDecisions(), 0)
	root, err := state.Root()
	if err != nil {
		t.Fatalf("resolving the root the surface writes under: %v", err)
	}

	code, said := held.fetch(t, held.at(surface.IdentityRoute))
	if code != http.StatusOK {
		t.Fatalf("a request with no token was answered with %d: %s", code, said)
	}

	var told map[string]any
	if err := json.Unmarshal([]byte(said), &told); err != nil {
		t.Fatalf("reading what the surface said about itself: %v", err)
	}
	if told["tool"] != "agora" {
		t.Errorf("the surface says it is %v, and it is agora", told["tool"])
	}
	if told["root"] != root {
		t.Errorf("the surface says it writes under %v, and it writes under %s", told["root"], root)
	}
	for _, named := range []string{"version", "commit", "uncommitted"} {
		if _, stated := told[named]; !stated {
			t.Errorf("the surface does not say %s of the build it was linked from: %s", named, said)
		}
	}
	for given := range told {
		switch given {
		case "tool", "root", "version", "commit", "uncommitted":
		default:
			t.Errorf("the surface tells %q to a request with no token, and that is more than who it is", given)
		}
	}
	for _, secret := range []string{held.server.Control(), held.token, held.session} {
		if strings.Contains(said, secret) {
			t.Errorf("the surface tells %q to a request with no token", secret)
		}
	}
}

func TestTheAnswerToWhoStandsHereIsRefusedToARequestNamingAnotherHost(t *testing.T) {
	// scenario: The identity answer is held to the host guard
	held := serving(t, twoDecisions(), 0)

	request, err := http.NewRequest(http.MethodGet, held.at(surface.IdentityRoute), nil)
	if err != nil {
		t.Fatalf("addressing the surface: %v", err)
	}
	request.Host = "somewhere.test"

	answered, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("reaching the surface: %v", err)
	}
	defer func() { _ = answered.Body.Close() }()

	if answered.StatusCode != http.StatusMisdirectedRequest {
		t.Errorf("who the surface is was told to a request naming another host: %d", answered.StatusCode)
	}
	said, _ := io.ReadAll(answered.Body)
	if strings.Contains(string(said), "agora") && strings.Contains(string(said), "root") {
		t.Errorf("a refused request was told who the surface is: %s", said)
	}
}

func TestAskingAnAddressIsAnsweredWithWhatTheSurfaceSaid(t *testing.T) {
	held := serving(t, twoDecisions(), 0)
	root, err := state.Root()
	if err != nil {
		t.Fatalf("resolving the root the surface writes under: %v", err)
	}

	identity, err := surface.Ask(held.server.Address()).Identify(context.Background())
	if err != nil {
		t.Fatalf("asking the surface who it is: %v", err)
	}
	if identity.Tool != surface.IdentityTool || identity.Root != root {
		t.Errorf("the surface is %q under %q, and it is agora under %q", identity.Tool, identity.Root, root)
	}
}

func TestAskingSomethingThatIsNotAgoraIsRefusedRatherThanTrusted(t *testing.T) {
	for named, handler := range map[string]http.HandlerFunc{
		"a server answering something else": func(writer http.ResponseWriter, _ *http.Request) { _, _ = io.WriteString(writer, "hello") },
		"a server answering nothing found":  func(writer http.ResponseWriter, _ *http.Request) { http.NotFound(writer, nil) },
		"a server answering json that is not an identity": func(writer http.ResponseWriter, _ *http.Request) {
			_, _ = io.WriteString(writer, `{"unrelated": true}`)
		},
	} {
		t.Run(named, func(t *testing.T) {
			foreign := httptest.NewServer(handler)
			t.Cleanup(foreign.Close)

			identity, err := surface.Ask(strings.TrimPrefix(foreign.URL, "http://")).Identify(context.Background())
			if err == nil && identity.Tool == surface.IdentityTool {
				t.Errorf("%s was taken for agora: %+v", named, identity)
			}
		})
	}
}

func TestAskingAnAddressNothingHoldsIsAnError(t *testing.T) {
	held := serving(t, twoDecisions(), 0)
	address := held.server.Address()
	held.server.Stop()
	held.server.Wait()

	if _, err := surface.Ask(address).Identify(context.Background()); err == nil {
		t.Errorf("asking %s, where nothing answers, was answered", address)
	}
}

func TestLocatingAskedForTheCommittedPortIgnoresARecordNamingAnotherAddress(t *testing.T) {
	// scenario: A record naming an address where nothing answers does not hide the surface on the port
	held := serving(t, twoDecisions(), 0)
	writeRecord(t, surface.Record{Address: "127.0.0.1:1", Token: held.server.Control(), Process: 1})

	found := locating(t, held.server.Address())
	if !found.Standing || !found.Recorded {
		t.Fatalf("the surface standing on %s was not found beside a record naming another address: %+v", held.server.Address(), found)
	}
	if want := held.server.Address(); found.Record.Address != want {
		t.Errorf("the surface was found at %q, and it stands at %q", found.Record.Address, want)
	}
	if found.Record.Token != held.server.Control() {
		t.Error("the token the record keeps was not carried to the surface found")
	}
}

func TestLocatingWithNoPortToAskReachesTheAddressTheRecordNames(t *testing.T) {
	// scenario: A bind of port 0 is found through the address the record names
	held := serving(t, twoDecisions(), 0)
	writeRecord(t, surface.Record{Address: held.server.Address(), Token: held.server.Control(), Process: 1})

	found := locating(t, aPortTheOperatingSystemChooses)
	if !found.Standing || found.Record.Address != held.server.Address() {
		t.Errorf("the surface standing at %s was not found through the record: %+v", held.server.Address(), found)
	}
}

func TestLocatingTrustsNoRecordWhereNothingAnswersAsThisRoot(t *testing.T) {
	// scenario: A record naming an address where nothing answers is not a surface
	t.Setenv(state.HomeVariable, t.TempDir())
	writeRecord(t, surface.Record{Address: "127.0.0.1:1", Token: "stale", Process: 1})

	found := locating(t, aPortTheOperatingSystemChooses)
	if found.Standing {
		t.Errorf("a record alone was taken for a surface: %+v", found)
	}
}

func TestLocatingNamesTheRootOfAnotherSurfaceHoldingThePort(t *testing.T) {
	// scenario: A port held by an agora surface of another root is named
	held := serving(t, twoDecisions(), 0)
	other, err := state.Root()
	if err != nil {
		t.Fatalf("resolving the root the surface writes under: %v", err)
	}
	t.Setenv(state.HomeVariable, t.TempDir())

	found := locating(t, held.server.Address())
	if found.Standing {
		t.Errorf("a surface of root %s was taken for this root's", other)
	}
	if found.AnotherRoot != other || found.AnotherAt != held.server.Address() {
		t.Errorf("the surface of another root was reported as %q at %q, want %q at %q", found.AnotherRoot, found.AnotherAt, other, held.server.Address())
	}
	if err := surface.CheckNotHeldByAnotherRoot(found, "agora.yaml"); err == nil {
		t.Error("a port held by a surface of another root was not refused")
	} else {
		for _, named := range []string{other, held.server.Address(), surface.BindField, "agora.yaml", surface.AddressOption} {
			if !strings.Contains(err.Error(), named) {
				t.Errorf("refusing a port held by another root does not name %q: %v", named, err)
			}
		}
	}
}

func TestLocatingFindsNothingWhereNothingHoldsThePort(t *testing.T) {
	held := serving(t, twoDecisions(), 0)
	address := held.server.Address()
	held.server.Stop()
	held.server.Wait()

	found := locating(t, address)
	if diff := cmp.Diff(surface.Whereabouts{}, found); diff != "" {
		t.Errorf("a port nothing holds was reported as holding something (-want +got):\n%s", diff)
	}
}

func TestAStandingSurfaceWithNoRecordIsFoundButNotReachable(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	found := locating(t, held.server.Address())
	if !found.Standing || found.Recorded {
		t.Fatalf("a surface with no record was found as %+v", found)
	}
	if _, err := found.Reached(); err == nil {
		t.Error("a surface whose token nobody holds was handed out as reachable")
	}
}

func writeRecord(t *testing.T, open surface.Record) {
	t.Helper()

	if err := surface.WriteRecord(open); err != nil {
		t.Fatalf("writing what the surface is reached at: %v", err)
	}
}

func locating(t *testing.T, committed string) surface.Whereabouts {
	t.Helper()

	asking, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	found, err := surface.Locate(asking, committed)
	if err != nil {
		t.Fatalf("locating the surface: %v", err)
	}
	return found
}
