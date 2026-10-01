package settingsmenu

import (
	"bytes"
	"errors"
	"io"

	"go.yaml.in/yaml/v3"
)

type fakeFile struct {
	Mode    string `yaml:"mode"`
	Ceiling int    `yaml:"ceiling"`
	Team    map[string]struct {
		Role string `yaml:"role"`
		Duty string `yaml:"duty"`
	} `yaml:"team"`
	Sealed *struct {
		First  string `yaml:"first"`
		Second string `yaml:"second"`
	} `yaml:"sealed"`
}

var errWrongAsAWhole = errors.New("sealed cannot state both of its fields")

func parseFake(path string, content []byte) error {
	var stated fakeFile
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	decoder.KnownFields(true)
	if err := decoder.Decode(&stated); err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	if stated.Sealed != nil && stated.Sealed.First != "" && stated.Sealed.Second != "" {
		return errWrongAsAWhole
	}
	return nil
}

func fakeDefinition(path string) Definition {
	return Definition{
		Path:  path,
		Parse: parseFake,
		Settings: func(string) ([]Setting, error) {
			return []Setting{{Key: "mode", Value: "plain", Source: "default"}}, nil
		},
		Wording: Wording{
			CannotBeRead:       "cannot read %s: %s",
			CannotBeWritten:    "cannot write %s: %s",
			CannotBeUnderstood: "%s states something unreadable.",
			ReadFromNowOn:      "Read from now on.",
		},
	}
}
