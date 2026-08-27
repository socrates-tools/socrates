package e2e_test

import (
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func (s *surfaced) configure(content string) {
	s.t.Helper()

	directory := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "agora")
	if err := os.MkdirAll(directory, 0o700); err != nil {
		s.t.Fatalf("making the directory agora reads its configuration from: %v", err)
	}
	if err := os.WriteFile(filepath.Join(directory, "config.yaml"), []byte(content), 0o600); err != nil {
		s.t.Fatalf("writing what an operator stated: %v", err)
	}
}

func aFreeLoopbackAddress(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("finding a loopback address nothing holds: %v", err)
	}
	address := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatalf("releasing the loopback address that was found: %v", err)
	}
	return address
}

func TestWithNothingStatedTheSurfaceBindsWhatTheOperatingSystemChose(t *testing.T) {
	held := standing(t)

	printed := held.serve(oneDecisionDocument)

	open := held.session()
	if printed != open.Address {
		t.Errorf("agora serve printed %q and bound %q", printed, open.Address)
	}

	host, port, err := net.SplitHostPort(open.Address)
	if err != nil {
		t.Fatalf("the surface bound %q: %v", open.Address, err)
	}
	if stated := net.ParseIP(host); stated == nil || !stated.IsLoopback() {
		t.Errorf("the surface bound the host %q, and it binds a loopback address only", host)
	}
	if port == "0" {
		t.Errorf("the surface reports the port it was asked for rather than the one it was given: %q", open.Address)
	}
}

const reachedAt = "https://agora.example:8443/"

func TestServePrintsTheAddressStatedOverTheAddressBound(t *testing.T) {
	held := standing(t)
	held.configure("reached_at: " + reachedAt + "\n")

	printed := held.serve(oneDecisionDocument)

	open := held.session()
	if printed != reachedAt {
		t.Errorf("agora serve printed %q, and an operator stated the surface is reached at %q", printed, reachedAt)
	}
	if printed == open.Address {
		t.Errorf("agora serve printed the address it bound, %q", open.Address)
	}
}

func TestASurfaceThatWasRevisedPrintsTheAddressStatedToo(t *testing.T) {
	// scenario: serve prints the address stated over the address bound
	held := standing(t)
	held.configure("reached_at: " + reachedAt + "\n")

	started := held.serve(oneDecisionDocument)
	revised := held.serve(oneDecisionDocument)

	if started != revised {
		t.Errorf("starting a surface printed %q and revising a live one printed %q", started, revised)
	}
	if revised != reachedAt {
		t.Errorf("revising a live surface printed %q, and an operator stated the surface is reached at %q", revised, reachedAt)
	}
}

func TestWithNothingStatedBothExitsPrintTheAddressBound(t *testing.T) {
	// scenario: serve returns before the person has read anything
	held := standing(t)

	started := held.serve(oneDecisionDocument)
	bound := held.session().Address
	revised := held.serve(oneDecisionDocument)

	if started != bound {
		t.Errorf("starting a surface printed %q and bound %q", started, bound)
	}
	if revised != bound {
		t.Errorf("revising a live surface printed %q and it is bound to %q", revised, bound)
	}
}

func TestAnAddressAPersonCouldNotOpenIsRefused(t *testing.T) {
	held := standing(t)
	held.configure("reached_at: agora.example:8443\n")

	printed, failures, code := held.reading(oneDecisionDocument, "serve", "--document", "-")
	if code == 0 {
		t.Fatalf("an address a person could not open was handed back at %q", printed)
	}
	if !strings.Contains(failures, "agora.example:8443") {
		t.Errorf("refusing an address a person could not open does not name it: %s", failures)
	}
	if _, err := os.Stat(filepath.Join(held.home, "session.json")); err == nil {
		t.Error("a refused address left a session behind, and it starts no server")
	}
}

func TestAnAddressNothingAnswersOnIsStillServed(t *testing.T) {
	held := standing(t)

	unserved := "http://" + aFreeLoopbackAddress(t) + "/"
	held.configure("reached_at: " + unserved + "\n")

	printed := held.serve(oneDecisionDocument)
	if printed != unserved {
		t.Errorf("agora serve printed %q, and an operator stated the surface is reached at %q", printed, unserved)
	}

	answered, err := http.Get("http://" + held.session().Address + "/")
	if err != nil {
		t.Fatalf("reaching the address the surface bound: %v", err)
	}
	defer answered.Body.Close()
	if answered.StatusCode != http.StatusOK {
		t.Errorf("the page the surface bound answered with %d, and the address it states is not one it serves", answered.StatusCode)
	}
}

func TestTheAddressTheSurfaceIsReachedAtIsNotAHostItAdmits(t *testing.T) {
	held := standing(t)
	held.configure("reached_at: " + reachedAt + "\n")

	if printed := held.serve(oneDecisionDocument); printed != reachedAt {
		t.Fatalf("agora serve printed %q, and an operator stated the surface is reached at %q", printed, reachedAt)
	}

	request, err := http.NewRequest(http.MethodGet, "http://"+held.session().Address+"/", nil)
	if err != nil {
		t.Fatalf("addressing the surface: %v", err)
	}
	request.Host = "agora.example:8443"

	answered, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("reaching the surface: %v", err)
	}
	defer answered.Body.Close()

	if answered.StatusCode != http.StatusMisdirectedRequest {
		t.Errorf("a request naming the address the surface is reached at was answered with %d", answered.StatusCode)
	}
	said, err := io.ReadAll(answered.Body)
	if err != nil {
		t.Fatalf("reading what the surface said: %v", err)
	}
	if strings.Contains(string(said), "one thing to settle") {
		t.Errorf("a request naming the address the surface is reached at reached a handler:\n%s", said)
	}
}

func TestAConfiguredBindThatIsNotLoopbackIsRefusedToo(t *testing.T) {
	held := standing(t)
	held.configure("bind: 0.0.0.0:8080\n")

	printed, failures, code := held.reading(oneDecisionDocument, "serve", "--document", "-")
	if code == 0 {
		t.Fatalf("a configured bind on every interface was served at %q", printed)
	}
	if !strings.Contains(failures, "0.0.0.0:8080") {
		t.Errorf("refusing a configured bind that is not loopback does not name it: %s", failures)
	}
	if _, err := os.Stat(filepath.Join(held.home, "session.json")); err == nil {
		t.Error("a refused configured bind left a session behind, and it starts no server")
	}
}

func TestAStatedFlagBeatsAConfiguredBind(t *testing.T) {
	held := standing(t)

	configured := aFreeLoopbackAddress(t)
	passed := aFreeLoopbackAddress(t)
	held.configure("bind: " + configured + "\n")

	printed, failures, code := held.reading(oneDecisionDocument, "serve", "--document", "-", "--address", passed)
	if code != 0 {
		t.Fatalf("agora serve exited %d: %s", code, failures)
	}

	if open := held.session(); open.Address != passed {
		t.Errorf("a caller passed %q, an operator stated %q, and the surface bound %q", passed, configured, open.Address)
	}
	if strings.TrimSpace(printed) != passed {
		t.Errorf("agora serve printed %q, and the caller passed %q", strings.TrimSpace(printed), passed)
	}
}

func TestAConfiguredBindIsWhatTheSurfaceBinds(t *testing.T) {
	held := standing(t)

	wanted := aFreeLoopbackAddress(t)
	held.configure("bind: " + wanted + "\n")

	printed := held.serve(oneDecisionDocument)

	if open := held.session(); open.Address != wanted {
		t.Errorf("an operator stated the bind %q and the surface bound %q", wanted, open.Address)
	}
	if printed != wanted {
		t.Errorf("agora serve printed %q, and an operator stated the bind %q", printed, wanted)
	}
}

func inFrontOf(t *testing.T, bound string) string {
	t.Helper()

	upstream, err := url.Parse("http://" + bound)
	if err != nil {
		t.Fatalf("addressing the surface the proxy stands in front of: %v", err)
	}

	forwarding := httputil.NewSingleHostReverseProxy(upstream)
	forwarding.FlushInterval = -1
	director := forwarding.Director
	forwarding.Director = func(request *http.Request) {
		director(request)
		request.Host = bound
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("binding the proxy: %v", err)
	}
	standing := &http.Server{Handler: forwarding, ReadHeaderTimeout: 10 * time.Second}
	go func() { _ = standing.Serve(listener) }()
	t.Cleanup(func() { _ = standing.Close() })

	return listener.Addr().String()
}

func TestAPersonAnswersThroughSomethingStandingInFrontOfTheSurface(t *testing.T) {
	const inTheirOwnWords = "the lock goes first, and the cache reads under it"

	held := standing(t)

	bound := aFreeLoopbackAddress(t)
	through := inFrontOf(t, bound)
	held.configure("bind: " + bound + "\nreached_at: http://" + through + "/\n")

	printed := held.serve(oneDecisionDocument)
	if printed != "http://"+through+"/" {
		t.Fatalf("agora serve printed %q, and the surface is reached at %q", printed, "http://"+through+"/")
	}
	if held.session().Address != bound {
		t.Fatalf("the surface bound %q, and the proxy stands in front of %q", held.session().Address, bound)
	}

	page := opening(t, through)
	page.evaluate(chose("c-1", 0), nil)
	page.write(within("c-1", ".words-field"), inTheirOwnWords)
	page.click(within("c-1", ".commit"))
	page.click(".finish")
	page.waitFor(`document.querySelector('.send').disabled === true`)

	answered, failures, code := held.run("poll")
	if code != 0 {
		t.Fatalf("agora poll exited %d: %s", code, failures)
	}
	if !strings.Contains(answered, "patch c-1") {
		t.Errorf("what agora poll printed does not carry the command the option carries:\n%s", answered)
	}
	if !strings.Contains(answered, inTheirOwnWords) {
		t.Errorf("what agora poll printed does not carry what the person wrote:\n%s", answered)
	}

	nothingLeftTheMachine(t, through, page.requests.all())
}
