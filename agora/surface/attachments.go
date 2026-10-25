package surface

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

const (
	messageFolder      = "message"
	longestStoredName  = 100
	longestExtension   = 16
	storedNameTemplate = "%d-%s"
)

func attachmentsNamed(names []string, uploads Uploads, claimed map[string]bool) ([]Attachment, error) {
	var attachments []Attachment
	for _, name := range names {
		staged, carried := uploads.Parts[name]
		if !carried {
			return nil, fmt.Errorf(answerNamesNoPart, name)
		}
		if claimed[name] {
			return nil, fmt.Errorf(answerAttachesAPartTwice, name)
		}
		claimed[name] = true
		attachments = append(attachments, Attachment{Name: staged.Name, Type: staged.Type, Bytes: staged.Bytes, Path: staged.Path})
	}
	return attachments, nil
}

func refuseUnnamedParts(uploads Uploads, claimed map[string]bool) error {
	for _, name := range slices.Sorted(maps.Keys(uploads.Parts)) {
		if !claimed[name] {
			return fmt.Errorf(partIsNamedByNoAnswer, name)
		}
	}
	return nil
}

func lay(answered *Answered, folder string) error {
	taken := map[string]bool{messageFolder: true, answerFile: true}
	for at := range answered.Decisions {
		decided := &answered.Decisions[at]
		owner := distinctFolder(taken, decided.ID)
		if err := layOwner(folder, owner, decided.Attachments); err != nil {
			return err
		}
	}
	return layOwner(folder, messageFolder, answered.Attachments)
}

func distinctFolder(taken map[string]bool, identifier string) string {
	owner := storedName(identifier)
	for at := 2; taken[owner]; at++ {
		owner = storedName(identifier) + "-" + strconv.Itoa(at)
	}
	taken[owner] = true
	return owner
}

func layOwner(folder, owner string, attachments []Attachment) error {
	if len(attachments) == 0 {
		return nil
	}
	within := filepath.Join(folder, owner)
	if err := os.MkdirAll(within, 0o700); err != nil {
		return fmt.Errorf(layingTheAttachments, err)
	}
	for at := range attachments {
		target := filepath.Join(within, fmt.Sprintf(storedNameTemplate, at+1, storedName(attachments[at].Name)))
		if err := os.Rename(attachments[at].Path, target); err != nil {
			return fmt.Errorf(layingTheAttachments, err)
		}
		attachments[at].Path = target
	}
	return nil
}

func storedName(stated string) string {
	base := stated[strings.LastIndexAny(stated, `/\`)+1:]

	var kept strings.Builder
	for _, character := range base {
		if unicode.IsLetter(character) || unicode.IsDigit(character) || strings.ContainsRune("._-", character) {
			kept.WriteRune(character)
			continue
		}
		kept.WriteRune('_')
	}

	name := strings.TrimLeft(kept.String(), ".")
	if name == "" {
		return unnamedAttachedFile
	}
	return shortened(name)
}

func shortened(name string) string {
	characters := []rune(name)
	if len(characters) <= longestStoredName {
		return name
	}
	extension := []rune(filepath.Ext(name))
	if len(extension) > longestExtension {
		extension = nil
	}
	return string(characters[:longestStoredName-len(extension)]) + string(extension)
}

func relocated(answered Answered, from, to string) Answered {
	moved := answered
	moved.Attachments = relocatedAttachments(answered.Attachments, from, to)
	if answered.Decisions != nil {
		moved.Decisions = make([]Decided, len(answered.Decisions))
	}
	for at, decided := range answered.Decisions {
		decided.Attachments = relocatedAttachments(decided.Attachments, from, to)
		moved.Decisions[at] = decided
	}
	return moved
}

func relocatedAttachments(attachments []Attachment, from, to string) []Attachment {
	if len(attachments) == 0 {
		return nil
	}
	moved := slices.Clone(attachments)
	for at := range moved {
		within, err := filepath.Rel(from, moved[at].Path)
		if err == nil {
			moved[at].Path = filepath.Join(to, within)
		}
	}
	return moved
}
