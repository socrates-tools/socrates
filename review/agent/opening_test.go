package agent

import (
	"slices"
	"strings"
	"testing"
)

func TestTheOpeningEventIsReadWhereverItArrivesInTheStream(t *testing.T) {
	stream := strings.Join([]string{
		`{"type":"system","subtype":"hook_started","hook_id":"h-1","hook_name":"SessionStart:startup","hook_event":"SessionStart","session_id":"s-1"}`,
		`{"type":"system","subtype":"hook_response","hook_id":"h-1","hook_name":"SessionStart:startup","hook_event":"SessionStart","output":"","exit_code":0,"outcome":"success","session_id":"s-1"}`,
		`{"type":"system","subtype":"init","session_id":"s-1","model":"claude-opus","tools":["Read","Grep","mcp__serena__find_symbol"],"mcp_servers":[{"name":"serena","status":"connected"}]}`,
		`{"type":"assistant","session_id":"s-1","message":{"model":"claude-opus","usage":{},"content":[{"type":"text","text":"read it"}]}}`,
		`{"type":"result","subtype":"success","session_id":"s-1","is_error":false,"result":"read it"}`,
	}, "\n") + "\n"

	read, err := parseStream(strings.NewReader(stream))
	if err != nil {
		t.Fatalf("a stream whose opening event follows a hook's events could not be read: %v", err)
	}

	if !slices.Equal(read.tools, []string{"Read", "Grep", "mcp__serena__find_symbol"}) {
		t.Errorf("the session's tools read as %v, so an opening event that is not the first line was not read", read.tools)
	}
	if !slices.Equal(read.toolSources, []string{"serena"}) {
		t.Errorf("the session's sources read as %v, so an opening event that is not the first line was not read", read.toolSources)
	}
}
