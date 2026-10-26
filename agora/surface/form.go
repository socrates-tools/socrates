package surface

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/socrates-tools/socrates/agora/render"
)

const (
	mebibyte = 1 << 20

	FileCeiling   = 20 * mebibyte
	AnswerCeiling = 100 * mebibyte
	TrayCeiling   = render.FilesPerTray

	answerPart          = "answer"
	incomingPattern     = "incoming-*"
	stagedPrefix        = "part-"
	unknownAttachedType = "application/octet-stream"
	unnamedAttachedFile = "file"
)

var errFileIsOverTheCeiling = errors.New("a file is over the ceiling")

type Staged struct {
	Name  string
	Type  string
	Bytes int64
	Path  string
}

type Uploads struct {
	Folder string
	Parts  map[string]Staged
}

type formReading struct {
	writer  http.ResponseWriter
	reader  *multipart.Reader
	session string
	sent    Sent
	stated  bool
	uploads Uploads
}

func readForm(writer http.ResponseWriter, request *http.Request, session string) (Sent, Uploads, bool) {
	request.Body = http.MaxBytesReader(writer, request.Body, AnswerCeiling)
	reader, err := request.MultipartReader()
	if err != nil {
		http.Error(writer, fmt.Sprintf(answerIsNotFormData, err), http.StatusBadRequest)
		return Sent{}, Uploads{}, false
	}

	reading := &formReading{writer: writer, reader: reader, session: session, uploads: Uploads{Parts: map[string]Staged{}}}
	if !reading.readEveryPart() {
		reading.discard()
		return Sent{}, Uploads{}, false
	}
	return reading.sent, reading.uploads, true
}

func (f *formReading) discard() {
	if f.uploads.Folder != "" {
		_ = os.RemoveAll(f.uploads.Folder)
	}
}

func (f *formReading) refuse(failure error, status int) bool {
	var overTheCeiling *http.MaxBytesError
	if errors.As(failure, &overTheCeiling) {
		http.Error(f.writer, fmt.Sprintf(answerIsOverTheCeiling, AnswerCeiling/mebibyte), http.StatusRequestEntityTooLarge)
		return false
	}
	if errors.Is(failure, errFileIsOverTheCeiling) {
		http.Error(f.writer, fmt.Sprintf(fileIsOverTheCeiling, FileCeiling/mebibyte), http.StatusRequestEntityTooLarge)
		return false
	}
	http.Error(f.writer, failure.Error(), status)
	return false
}

func (f *formReading) readEveryPart() bool {
	for {
		part, err := f.reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return f.refuse(err, http.StatusBadRequest)
		}
		if !f.readPart(part) {
			return false
		}
	}
	if !f.stated {
		return f.refuse(fmt.Errorf(answerHasNoAnswerPart, answerPart), http.StatusBadRequest)
	}
	return true
}

func (f *formReading) readPart(part *multipart.Part) bool {
	defer func() { _ = part.Close() }()

	name := part.FormName()
	if name == answerPart {
		return f.readAnswerPart(part)
	}
	return f.stage(name, part)
}

func (f *formReading) readAnswerPart(part *multipart.Part) bool {
	if f.stated {
		return f.refuse(fmt.Errorf(answerPartIsRepeated, answerPart), http.StatusBadRequest)
	}
	raw, err := io.ReadAll(io.LimitReader(part, BodyCeiling+1))
	if err != nil {
		return f.refuse(err, http.StatusBadRequest)
	}
	if len(raw) > BodyCeiling {
		return f.refuse(fmt.Errorf(bodyIsOverTheCeiling, BodyCeiling), http.StatusRequestEntityTooLarge)
	}

	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&f.sent); err != nil {
		return f.refuse(err, http.StatusBadRequest)
	}
	f.stated = true
	return true
}

func (f *formReading) stage(name string, part *multipart.Part) bool {
	if _, repeated := f.uploads.Parts[name]; repeated {
		return f.refuse(fmt.Errorf(answerPartIsRepeated, name), http.StatusBadRequest)
	}
	if f.uploads.Folder == "" {
		folder, err := f.freshFolder()
		if err != nil {
			return f.refuse(err, http.StatusInternalServerError)
		}
		f.uploads.Folder = folder
	}

	path := filepath.Join(f.uploads.Folder, stagedPrefix+strconv.Itoa(len(f.uploads.Parts)+1))
	held, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return f.refuse(fmt.Errorf(stagingTheAttachment, err), http.StatusInternalServerError)
	}
	written, copying := io.Copy(held, io.LimitReader(part, FileCeiling+1))
	closing := held.Close()
	if copying != nil {
		return f.refuse(copying, http.StatusBadRequest)
	}
	if closing != nil {
		return f.refuse(fmt.Errorf(stagingTheAttachment, closing), http.StatusInternalServerError)
	}
	if written > FileCeiling {
		return f.refuse(errFileIsOverTheCeiling, http.StatusRequestEntityTooLarge)
	}

	f.uploads.Parts[name] = Staged{Name: statedName(part), Type: statedType(part), Bytes: written, Path: path}
	return true
}

func (f *formReading) freshFolder() (string, error) {
	within, err := SessionPath(f.session)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(within, 0o700); err != nil {
		return "", fmt.Errorf(stagingTheAttachment, err)
	}
	folder, err := os.MkdirTemp(within, incomingPattern)
	if err != nil {
		return "", fmt.Errorf(stagingTheAttachment, err)
	}
	return folder, nil
}

func statedName(part *multipart.Part) string {
	if named := part.FileName(); named != "" {
		return named
	}
	return unnamedAttachedFile
}

func statedType(part *multipart.Part) string {
	if stated := part.Header.Get("Content-Type"); stated != "" {
		return stated
	}
	return unknownAttachedType
}
