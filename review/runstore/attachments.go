package runstore

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/review/dialogue"
)

const (
	attachmentsDirectory = "attachments"
	sniffedBytes         = 512
	untypedFile          = "application/octet-stream"
	plainText            = "text/plain"
)

var outsideAFileName = regexp.MustCompile(`[^A-Za-z0-9._-]`)

func (r Round) Attach(sources []string) ([]dialogue.Attachment, error) {
	if len(sources) == 0 {
		return nil, nil
	}
	for _, source := range sources {
		if err := mustBeAFile(source); err != nil {
			return nil, err
		}
	}

	parent := filepath.Join(r.Directory, attachmentsDirectory)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return nil, fmt.Errorf(creating, parent, err)
	}
	directory, err := os.MkdirTemp(parent, "")
	if err != nil {
		return nil, fmt.Errorf(creating, parent, err)
	}
	directory, err = filepath.Abs(directory)
	if err != nil {
		_ = os.RemoveAll(directory)
		return nil, fmt.Errorf(resolving, directory, err)
	}

	attached := make([]dialogue.Attachment, 0, len(sources))
	for position, source := range sources {
		name := filepath.Base(source)
		stored := filepath.Join(directory, strconv.Itoa(position+1)+"-"+outsideAFileName.ReplaceAllString(name, "_"))
		kind, err := copyFile(source, stored)
		if err != nil {
			_ = os.RemoveAll(directory)
			return nil, err
		}
		attached = append(attached, dialogue.Attachment{Name: name, Type: kind, Path: stored})
	}
	return attached, nil
}

func mustBeAFile(source string) error {
	info, err := os.Stat(source)
	if errors.Is(err, fs.ErrNotExist) {
		return refusal.New(noSuchAttachment, source)
	}
	if err != nil {
		return fmt.Errorf(reading, source, err)
	}
	if !info.Mode().IsRegular() {
		return refusal.New(attachmentIsNotAFile, source)
	}
	return nil
}

func copyFile(source, destination string) (string, error) {
	from, err := os.Open(source)
	if err != nil {
		return "", fmt.Errorf(copying, source, err)
	}
	defer func() { _ = from.Close() }()

	to, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return "", fmt.Errorf(copying, source, err)
	}

	head := make([]byte, sniffedBytes)
	read, err := io.ReadFull(from, head)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		_ = to.Close()
		return "", fmt.Errorf(copying, source, err)
	}
	head = head[:read]

	if _, err := to.Write(head); err != nil {
		_ = to.Close()
		return "", fmt.Errorf(copying, source, err)
	}
	if _, err := io.Copy(to, from); err != nil {
		_ = to.Close()
		return "", fmt.Errorf(copying, source, err)
	}
	if err := to.Close(); err != nil {
		return "", fmt.Errorf(copying, source, err)
	}
	return typeOf(source, head), nil
}

func typeOf(source string, head []byte) string {
	if named := mime.TypeByExtension(filepath.Ext(source)); named != "" {
		return withoutParameters(named)
	}
	if len(head) == 0 {
		return untypedFile
	}
	if utf8.Valid(head) && !bytes.Contains(head, []byte{0}) {
		return plainText
	}
	return untypedFile
}

func withoutParameters(kind string) string {
	media, _, found := strings.Cut(kind, ";")
	if !found {
		return kind
	}
	return strings.TrimSpace(media)
}
