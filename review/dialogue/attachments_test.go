package dialogue_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/review/dialogue"
)

func attachments(count int) []dialogue.Attachment {
	held := make([]dialogue.Attachment, count)
	for position := range held {
		held[position] = dialogue.Attachment{Name: "shot.png", Type: "image/png", Path: "/runs/one/r1/attachments/a/1-shot.png"}
	}
	return held
}

func TestAResponseCarriesUpToTenAttachments(t *testing.T) {
	for _, resolution := range []dialogue.Resolution{dialogue.Decided, dialogue.Rejected, dialogue.Dropped, dialogue.Fixed} {
		t.Run("Ten attachments with "+string(resolution), func(t *testing.T) {
			err := dialogue.Validate(dialogue.Response{Finding: "f1", Resolution: resolution, Reason: "because", Attachments: attachments(10)})
			if err != nil {
				t.Errorf("ten attachments were refused: %v", err)
			}
		})
	}

	t.Run("An eleventh attachment is refused naming the ceiling", func(t *testing.T) {
		err := dialogue.Validate(dialogue.Response{Finding: "f1", Resolution: dialogue.Rejected, Reason: "because", Attachments: attachments(11)})
		if err == nil || !strings.Contains(err.Error(), "10") {
			t.Errorf("eleven attachments were answered with %v, want a refusal naming ten", err)
		}
	})

	t.Run("The ceiling is stated once", func(t *testing.T) {
		if err := dialogue.Attaching(dialogue.Rejected, dialogue.MaxAttachments+1); err == nil {
			t.Error("one past the ceiling was accepted")
		}
		if err := dialogue.Attaching(dialogue.Rejected, dialogue.MaxAttachments); err != nil {
			t.Errorf("the ceiling itself was refused: %v", err)
		}
	})
}

func TestAForwardedResponseCarriesNoAttachment(t *testing.T) {
	err := dialogue.Validate(dialogue.Response{Finding: "f1", Resolution: dialogue.Forwarded, Attachments: attachments(1)})
	if err == nil {
		t.Fatal("a forward carrying a file was accepted")
	}
	for _, named := range []string{"forward", "pull request"} {
		if !strings.Contains(err.Error(), named) {
			t.Errorf("the refusal does not name %q: %v", named, err)
		}
	}
}

func TestAResponseWithAttachmentsSurvivesTheDocument(t *testing.T) {
	response := dialogue.Response{Finding: "f1", Resolution: dialogue.Decided, Reason: "see the screenshot", Attachments: attachments(2)}

	encoded, err := dialogue.Line(response)
	if err != nil {
		t.Fatalf("encoding: %v", err)
	}
	read, err := dialogue.Parse(encoded)
	if err != nil {
		t.Fatalf("parsing: %v", err)
	}
	answer, _ := read.For("f1")
	if len(answer.Attachments) != 2 || answer.Attachments[0].Path != response.Attachments[0].Path || answer.Attachments[1].Type != "image/png" {
		t.Errorf("the attachments read back as %+v", answer.Attachments)
	}
}
