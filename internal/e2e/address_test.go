package e2e_test

import (
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/ktei/socrates/internal/surface"
	"github.com/ktei/socrates/internal/testing/source"
)

func (s *surfaced) configure(stated string) {
	s.t.Helper()

	if err := os.MkdirAll(s.home, 0o700); err != nil {
		s.t.Fatalf("making the home the one configuration file is read from: %v", err)
	}
	if err := os.WriteFile(filepath.Join(s.home, "config.yaml"), []byte(asASurfaceBlock(stated)), 0o600); err != nil {
		s.t.Fatalf("writing what an operator stated: %v", err)
	}
}

func asASurfaceBlock(stated string) string {
	block := "surface:\n"
	for _, line := range strings.Split(strings.TrimRight(stated, "\n"), "\n") {
		block += "  " + line + "\n"
	}
	return block
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

func thePageAt(reached, session string) string {
	return strings.TrimSuffix(reached, "/") + surface.SessionPrefix + session
}

func TestTheSurfaceReportsThePortItWasGivenRatherThanTheOneItAskedFor(t *testing.T) {
	held := standing(t)

	printed := held.serve(oneDecisionDocument)

	open := held.record()
	if want := thePageAt("http://"+open.Address, held.opened); printed != want {
		t.Errorf("agora serve printed %q and bound %q, so the page is at %q", printed, open.Address, want)
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

func untilThePortIsFree(t *testing.T, address string) {
	t.Helper()

	deadline := time.Now().Add(10 * time.Second)
	for {
		listener, err := net.Listen("tcp", address)
		if err == nil {
			_ = listener.Close()
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("the port %s is still held after the surface was stopped: %v", address, err)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func requireTheFixedPortIsFree(t *testing.T) {
	t.Helper()

	listener, err := net.Listen("tcp", surface.DefaultAddress)
	if err != nil {
		t.Skipf("something else holds %s, so what the surface binds when nothing states one cannot be read here: %v", surface.DefaultAddress, err)
	}
	if err := listener.Close(); err != nil {
		t.Fatalf("releasing %s: %v", surface.DefaultAddress, err)
	}
}

func holdingTheWindowDownTo(t *testing.T, ports int) []string {
	t.Helper()

	addresses := surface.WindowAddresses()
	for _, address := range addresses[:ports] {
		listener, err := net.Listen("tcp", address)
		if err != nil {
			t.Skipf("something else already holds %s, so a window walked past a held port cannot be read here: %v", address, err)
		}
		t.Cleanup(func() { _ = listener.Close() })
	}
	return addresses
}

func TestAHeldDefaultIsTheNextPortDownRatherThanARefusal(t *testing.T) {
	addresses := holdingTheWindowDownTo(t, 1)

	held := standing(t)
	held.bindsWhatIsConfigured()

	printed := held.serve(oneDecisionDocument)

	if bound := held.record().Address; bound != addresses[1] {
		t.Fatalf("the port at the top of the window is held and the surface bound %q, and the next port down is %q", bound, addresses[1])
	}
	if want := thePageAt("http://"+addresses[1], held.opened); printed != want {
		t.Errorf("agora serve printed %q, and the page it served is at %q", printed, want)
	}

	answered, err := http.Get(printed)
	if err != nil {
		t.Fatalf("reaching the page the surface walked down to: %v", err)
	}
	defer func() { _ = answered.Body.Close() }()
	if answered.StatusCode != http.StatusOK {
		t.Errorf("the page at %s answered with %d", printed, answered.StatusCode)
	}
}

func TestTheWindowIsWalkedFromTheTopDownRatherThanAtRandom(t *testing.T) {
	addresses := holdingTheWindowDownTo(t, 3)

	held := standing(t)
	held.bindsWhatIsConfigured()

	held.serve(oneDecisionDocument)
	bound := held.record().Address

	if bound != addresses[3] {
		t.Fatalf("the three ports at the top of the window are held and the surface bound %q, and the highest free one is %q", bound, addresses[3])
	}

	if _, failures, code := held.run("stop"); code != 0 {
		t.Fatalf("agora stop exited %d: %s", code, failures)
	}
	untilThePortIsFree(t, bound)

	held.serve(oneDecisionDocument)
	if again := held.record().Address; again != bound {
		t.Errorf("a surface started again after the first stopped bound %q, and the first bound %q", again, bound)
	}
}

func TestAWholeWindowHeldIsARefusalThatNamesTheWindow(t *testing.T) {
	addresses := holdingTheWindowDownTo(t, len(surface.WindowAddresses()))

	held := standing(t)
	held.bindsWhatIsConfigured()

	printed, failures, code := held.reading(oneDecisionDocument, "serve", "--document", "-")
	if code == 0 {
		t.Fatalf("every port of the window is held and a page was served at %q", printed)
	}
	for _, named := range []string{addresses[0], addresses[len(addresses)-1]} {
		if !strings.Contains(failures, named) {
			t.Errorf("refusing a whole window held does not name %s:\n%s", named, failures)
		}
	}
	if _, err := os.Stat(filepath.Join(held.surfaceHome(), "surface.json")); err == nil {
		t.Error("a refused window left a surface behind, and it starts no server")
	}
}

func aHeldLoopbackAddress(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("holding a loopback address: %v", err)
	}
	t.Cleanup(func() { _ = listener.Close() })
	return listener.Addr().String()
}

func windowAddressesNothingHolds(t *testing.T) []string {
	t.Helper()

	var free []string
	for _, address := range surface.WindowAddresses() {
		listener, err := net.Listen("tcp", address)
		if err != nil {
			continue
		}
		_ = listener.Close()
		free = append(free, address)
	}
	return free
}

func (s *surfaced) boundNothing(named string, free []string) {
	s.t.Helper()

	if _, err := os.Stat(filepath.Join(s.surfaceHome(), "surface.json")); err == nil {
		s.t.Errorf("refusing %s left a surface behind, and it starts no server", named)
	}
	for _, address := range free {
		if address == named {
			continue
		}
		listener, err := net.Listen("tcp", address)
		if err != nil {
			s.t.Errorf("refusing %s bound %s instead of taking no port at all", named, address)
			continue
		}
		_ = listener.Close()
	}
}

func TestAStatedAddressThatIsHeldIsRefusedRatherThanWalked(t *testing.T) {
	taken := aHeldLoopbackAddress(t)

	held := standing(t)

	free := windowAddressesNothingHolds(t)
	printed, failures, code := held.reading(oneDecisionDocument, "serve", "--document", "-", "--address", taken)
	if code == 0 {
		t.Fatalf("a caller stated the held address %s and a page was served at %q", taken, printed)
	}
	if !strings.Contains(failures, taken) {
		t.Errorf("refusing a stated address that is held does not name it: %s", failures)
	}
	held.boundNothing(taken, free)
}

func TestAConfiguredBindThatIsHeldIsRefusedRatherThanWalked(t *testing.T) {
	taken := aHeldLoopbackAddress(t)

	held := standing(t)
	held.configure("bind: " + taken + "\n")
	held.bindsWhatIsConfigured()

	free := windowAddressesNothingHolds(t)
	printed, failures, code := held.reading(oneDecisionDocument, "serve", "--document", "-")
	if code == 0 {
		t.Fatalf("an operator stated the held bind %s and a page was served at %q", taken, printed)
	}
	if !strings.Contains(failures, taken) {
		t.Errorf("refusing a configured bind that is held does not name it: %s", failures)
	}
	held.boundNothing(taken, free)
}

func TestAStatedReachedAtAddressPinsThePortToo(t *testing.T) {
	addresses := holdingTheWindowDownTo(t, 1)

	held := standing(t)
	held.configure("reached_at: " + reachedAt + "\n")
	held.bindsWhatIsConfigured()

	free := windowAddressesNothingHolds(t)
	printed, failures, code := held.reading(oneDecisionDocument, "serve", "--document", "-")
	if code == 0 {
		t.Fatalf("an operator stated where the surface is reached at, the port at the top of the window is held, and a page was served at %q", printed)
	}
	if !strings.Contains(failures, addresses[0]) {
		t.Errorf("refusing a held port an operator committed to does not name it: %s", failures)
	}
	held.boundNothing(addresses[0], free)
}

func TestTheDefaultBindIsAFixedLoopbackPort(t *testing.T) {
	requireTheFixedPortIsFree(t)

	held := standing(t)
	held.bindsWhatIsConfigured()

	held.serve(oneDecisionDocument)
	bound := held.record().Address

	if bound != surface.DefaultAddress {
		t.Fatalf("with nothing stated the surface bound %q, and the address it binds when nothing states one is %q", bound, surface.DefaultAddress)
	}
	host, _, err := net.SplitHostPort(bound)
	if err != nil {
		t.Fatalf("the surface bound %q: %v", bound, err)
	}
	if stated := net.ParseIP(host); stated == nil || !stated.IsLoopback() {
		t.Errorf("the surface bound the host %q, and the address it binds when nothing states one is a loopback one", host)
	}

	if _, failures, code := held.run("stop"); code != 0 {
		t.Fatalf("agora stop exited %d: %s", code, failures)
	}
	untilThePortIsFree(t, bound)

	held.serve(oneDecisionDocument)
	if again := held.record().Address; again != bound {
		t.Errorf("a surface started again after the first stopped bound %q, and the first bound %q", again, bound)
	}
}

const reachedAt = "https://agora.example:8443/"

func TestServePrintsTheAddressStatedOverTheAddressBound(t *testing.T) {
	held := standing(t)
	held.configure("reached_at: " + reachedAt + "\n")

	printed := held.serve(oneDecisionDocument)

	open := held.record()
	if want := thePageAt(reachedAt, held.opened); printed != want {
		t.Errorf("agora serve printed %q, and an operator stated the surface is reached at %q", printed, want)
	}
	if strings.Contains(printed, open.Address) {
		t.Errorf("agora serve printed the address it bound, %q", open.Address)
	}
}

func TestASurfaceThatWasRevisedPrintsTheAddressStatedToo(t *testing.T) {
	// scenario: serve prints the address stated over the address bound
	held := standing(t)
	held.configure("reached_at: " + reachedAt + "\n")

	started := held.serve(oneDecisionDocument)
	revised := held.serveOnto(held.opened, oneDecisionDocument)

	if started != revised {
		t.Errorf("starting a surface printed %q and revising the page it opened printed %q", started, revised)
	}
	if want := thePageAt(reachedAt, held.opened); revised != want {
		t.Errorf("revising a live surface printed %q, and an operator stated the surface is reached at %q", revised, want)
	}
}

func TestWithNoAddressToReachAtStatedBothExitsPrintTheAddressBound(t *testing.T) {
	// scenario: serve returns before the person has read anything
	held := standing(t)

	started := held.serve(oneDecisionDocument)
	bound := held.record().Address
	revised := held.serveOnto(held.opened, oneDecisionDocument)

	want := thePageAt("http://"+bound, held.opened)
	if started != want {
		t.Errorf("starting a surface printed %q and bound %q", started, bound)
	}
	if revised != want {
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
	if _, err := os.Stat(filepath.Join(held.surfaceHome(), "surface.json")); err == nil {
		t.Error("a refused address left a session behind, and it starts no server")
	}
}

func TestAnAddressNothingAnswersOnIsStillServed(t *testing.T) {
	held := standing(t)

	unserved := "http://" + aFreeLoopbackAddress(t) + "/"
	held.configure("reached_at: " + unserved + "\n")

	printed := held.serve(oneDecisionDocument)
	if want := thePageAt(unserved, held.opened); printed != want {
		t.Errorf("agora serve printed %q, and an operator stated the surface is reached at %q", printed, want)
	}

	answered, err := http.Get(thePageAt("http://"+held.record().Address, held.opened))
	if err != nil {
		t.Fatalf("reaching the address the surface bound: %v", err)
	}
	defer func() { _ = answered.Body.Close() }()
	if answered.StatusCode != http.StatusOK {
		t.Errorf("the page the surface bound answered with %d, and the address it states is not one it serves", answered.StatusCode)
	}
}

func TestTheAddressTheSurfaceIsReachedAtIsNotAHostItAdmits(t *testing.T) {
	held := standing(t)
	held.configure("reached_at: " + reachedAt + "\n")

	printed := held.serve(oneDecisionDocument)
	if want := thePageAt(reachedAt, held.opened); printed != want {
		t.Fatalf("agora serve printed %q, and an operator stated the surface is reached at %q", printed, want)
	}

	request, err := http.NewRequest(http.MethodGet, "http://"+held.record().Address+"/", nil)
	if err != nil {
		t.Fatalf("addressing the surface: %v", err)
	}
	request.Host = "agora.example:8443"

	answered, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("reaching the surface: %v", err)
	}
	defer func() { _ = answered.Body.Close() }()

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
	if _, err := os.Stat(filepath.Join(held.surfaceHome(), "surface.json")); err == nil {
		t.Error("a refused configured bind left a session behind, and it starts no server")
	}
}

func TestThePortAnOperatorIsToldToForwardIsTheOneTheSurfaceBinds(t *testing.T) {
	section := documentedForwarding(t)

	_, port, err := net.SplitHostPort(surface.DefaultAddress)
	if err != nil {
		t.Fatalf("reading the port the surface binds from %q: %v", surface.DefaultAddress, err)
	}

	if !strings.Contains(section, surface.DefaultAddress) {
		t.Errorf("what an operator is told to forward does not name %s, which is what the surface binds:\n%s", surface.DefaultAddress, section)
	}
	for _, named := range regexp.MustCompile(`\b\d{4,5}\b`).FindAllString(section, -1) {
		if named != port {
			t.Errorf("what an operator is told to forward names the port %s, and the surface binds %s", named, port)
		}
	}
}

func TestTheAddressAPersonOpensNamesTheToolTheyInstalled(t *testing.T) {
	section := documentedForwarding(t)

	for _, line := range strings.Split(section, "\n") {
		if !strings.Contains(line, "reached_at") && !strings.Contains(line, "localhost/") {
			continue
		}
		if strings.Contains(line, "agora") {
			t.Errorf("the address a person is told to open is %q, and somebody who installed socrates has no reason to know what serves it", line)
		}
	}
}

func documentedForwarding(t *testing.T) string {
	t.Helper()

	const heading = "## Reading the page from another machine"

	held, err := os.ReadFile(filepath.Join(source.Root(t), "README.md"))
	if err != nil {
		t.Fatalf("reading the README: %v", err)
	}

	at := strings.Index(string(held), heading)
	if at < 0 {
		t.Fatalf("the README does not say how a person on another machine reads the page")
	}
	rest := string(held)[at+len(heading):]
	if end := strings.Index(rest, "\n## "); end >= 0 {
		return rest[:end]
	}
	return rest
}

func TestAPageServedWithNothingStatedSaysWhereItCanBeOpenedFrom(t *testing.T) {
	held := standing(t)

	printed, failures, code := held.reading(oneDecisionDocument, "serve", "--document", "-", "--address", aPortTheOperatingSystemChooses)
	if code != 0 {
		t.Fatalf("agora serve exited %d: %s", code, failures)
	}

	for _, said := range []string{"this machine", surface.ReachedAtField} {
		if !strings.Contains(failures, said) {
			t.Errorf("what a person reads is %q, without saying %q, and a browser elsewhere gets an empty reply that names no cause", failures, said)
		}
	}

	address := strings.TrimSpace(printed)
	answered, err := http.Get(address)
	if err != nil {
		t.Fatalf("reaching the address agora serve printed: %v", err)
	}
	defer func() { _ = answered.Body.Close() }()
	if answered.StatusCode != http.StatusOK {
		t.Errorf("the page answered with %d, and what was said is a statement rather than a refusal", answered.StatusCode)
	}
}

func TestWhatAPersonIsToldToForwardIsThePortTheWindowWalkedDownTo(t *testing.T) {
	// scenario: A page served with nothing stated says where it can be opened from
	addresses := holdingTheWindowDownTo(t, 1)

	held := standing(t)
	held.bindsWhatIsConfigured()

	printed, failures, code := held.reading(oneDecisionDocument, "serve", "--document", "-")
	if code != 0 {
		t.Fatalf("agora serve exited %d: %s", code, failures)
	}
	if bound := held.record().Address; bound != addresses[1] {
		t.Fatalf("the port at the top of the window is held and the surface bound %q, and the next port down is %q", bound, addresses[1])
	}

	if !strings.Contains(failures, addresses[1]) {
		t.Errorf("what a person is told to forward is %q, and the surface bound %s", failures, addresses[1])
	}
	if strings.Contains(failures, addresses[0]) {
		t.Errorf("what a person is told to forward names %s, and the surface walked past it:\n%s", addresses[0], failures)
	}
	if got := strings.TrimSpace(printed); !strings.HasPrefix(got, "http://"+addresses[1]+surface.SessionPrefix) {
		t.Errorf("agora serve printed %q, and the surface bound %s", got, addresses[1])
	}
}

func TestWhatACallerReadsIsTheAddressAndNothingElse(t *testing.T) {
	held := standing(t)

	printed, failures, code := held.reading(oneDecisionDocument, "serve", "--document", "-", "--address", aPortTheOperatingSystemChooses)
	if code != 0 {
		t.Fatalf("agora serve exited %d: %s", code, failures)
	}

	lines := strings.Split(strings.TrimSpace(printed), "\n")
	if len(lines) != 1 {
		t.Fatalf("the stream a caller parses holds %d lines, and a caller reads it to learn where the page is:\n%s", len(lines), printed)
	}
	if theSessionIn(t, lines[0]) == "" {
		t.Errorf("the one line a caller reads is %q, and it is the address of the page", lines[0])
	}
}

func TestAStatedAddressIsAnOperatorWhoHasAlreadyAnswered(t *testing.T) {
	held := standing(t)
	held.configure("reached_at: http://socrates.example:8443/\n")

	_, failures, code := held.reading(oneDecisionDocument, "serve", "--document", "-", "--address", aPortTheOperatingSystemChooses)
	if code != 0 {
		t.Fatalf("agora serve exited %d: %s", code, failures)
	}
	if strings.Contains(failures, "this machine") || strings.Contains(failures, surface.ReachedAtField) {
		t.Errorf("what a person reads is %q, and an operator who stated where the surface is reached at has answered that already", failures)
	}
}

func TestWhatAnOperatorStatesAboutTheSurfaceIsReadFromTheOneFile(t *testing.T) {
	held := standing(t)

	bound := aFreeLoopbackAddress(t)
	reachedAt := "https://agora.example:8443/"
	held.configure("bind: " + bound + "\nreached_at: " + reachedAt + "\n")
	held.bindsWhatIsConfigured()

	printed := strings.TrimSpace(held.serve(oneDecisionDocument))

	if open := held.record(); open.Address != bound {
		t.Errorf("the one file states the bind %q and the surface bound %q", bound, open.Address)
	}
	if !strings.HasPrefix(printed, reachedAt+strings.TrimPrefix(surface.SessionPrefix, "/")) {
		t.Errorf("agora serve printed %q, and the one file states the surface is reached at %q", printed, reachedAt)
	}
	if _, err := os.Stat(filepath.Join(held.surfaceHome(), "config.yaml")); err == nil {
		t.Errorf("the surface read a configuration file of its own under %s", held.surfaceHome())
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

	open := held.record()
	if open.Address != passed {
		t.Errorf("a caller passed %q, an operator stated %q, and the surface bound %q", passed, configured, open.Address)
	}
	if !strings.HasPrefix(strings.TrimSpace(printed), "http://"+passed+surface.SessionPrefix) {
		t.Errorf("agora serve printed %q, and the caller passed %q", strings.TrimSpace(printed), passed)
	}
}

func TestAConfiguredBindIsWhatTheSurfaceBinds(t *testing.T) {
	held := standing(t)

	wanted := aFreeLoopbackAddress(t)
	held.configure("bind: " + wanted + "\n")
	held.bindsWhatIsConfigured()

	printed := held.serve(oneDecisionDocument)

	if open := held.record(); open.Address != wanted {
		t.Errorf("an operator stated the bind %q and the surface bound %q", wanted, open.Address)
	}
	if want := thePageAt("http://"+wanted, held.opened); printed != want {
		t.Errorf("agora serve printed %q, and an operator stated the bind %q", printed, wanted)
	}
}

func inFrontOf(t *testing.T, bound string) string {
	t.Helper()

	upstream, err := url.Parse("http://" + bound)
	if err != nil {
		t.Fatalf("addressing the surface the proxy stands in front of: %v", err)
	}

	forwarding := &httputil.ReverseProxy{
		FlushInterval: -1,
		Rewrite: func(request *httputil.ProxyRequest) {
			request.SetURL(upstream)
			request.Out.Host = bound
		},
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
	held.bindsWhatIsConfigured()

	printed := held.serve(oneDecisionDocument)
	if want := thePageAt("http://"+through, held.opened); printed != want {
		t.Fatalf("agora serve printed %q, and the surface is reached at %q", printed, want)
	}
	if held.record().Address != bound {
		t.Fatalf("the surface bound %q, and the proxy stands in front of %q", held.record().Address, bound)
	}

	page := opening(t, printed)
	page.evaluate(chose("c-1", 0), nil)
	page.write(within("c-1", ".words-field"), inTheirOwnWords)
	page.click(within("c-1", ".commit"))
	page.click(".finish")
	page.waitFor(`document.querySelector('.send').disabled === true`)

	answered, failures, code := held.poll()
	if code != 0 {
		t.Fatalf("agora poll exited %d: %s", code, failures)
	}
	if !strings.Contains(answered, "patch c-1") {
		t.Errorf("what agora poll printed does not carry the command the option carries:\n%s", answered)
	}
	if !strings.Contains(answered, inTheirOwnWords) {
		t.Errorf("what agora poll printed does not carry what the person wrote:\n%s", answered)
	}

	nothingLeftTheMachine(t, printed, page.requests.all())
}
