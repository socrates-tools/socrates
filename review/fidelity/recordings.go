package fidelity

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	TestdataDirectory  = "testdata"
	RecordingExtension = ".json"

	failuresFolder = "failures"
)

func RecordingsFor(vendor string) string {
	return filepath.Join(TestdataDirectory, vendor)
}

type Recording struct {
	Name   string
	Path   string
	Shapes []Shape
}

func Recordings(directory string) ([]Recording, error) {
	paths, err := filepath.Glob(filepath.Join(directory, "*"+RecordingExtension))
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf(noRecordings, directory)
	}
	slices.Sort(paths)

	recordings := make([]Recording, 0, len(paths))
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf(readingARecording, path, err)
		}

		var shapes []Shape
		if err := json.Unmarshal(content, &shapes); err != nil {
			return nil, fmt.Errorf(readingARecording, path, err)
		}
		recordings = append(recordings, Recording{
			Name:   strings.TrimSuffix(filepath.Base(path), RecordingExtension),
			Path:   path,
			Shapes: shapes,
		})
	}
	return recordings, nil
}

func Keep(root, vendor, name, stream string, answered error) (string, error) {
	shapes, err := ShapesOf(stream)
	if err != nil {
		return "", err
	}
	return record(filepath.Join(root, folderFor(vendor, answered)), name, shapes)
}

func folderFor(vendor string, answered error) string {
	if answered != nil {
		return filepath.Join(failuresFolder, vendor)
	}
	return vendor
}

func record(directory, name string, shapes []Shape) (string, error) {
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return "", fmt.Errorf(writingARecording, directory, err)
	}

	encoded, err := json.MarshalIndent(shapes, "", "  ")
	if err != nil {
		return "", fmt.Errorf(writingARecording, name, err)
	}

	path := filepath.Join(directory, name+RecordingExtension)
	if err := os.WriteFile(path, append(encoded, '\n'), 0o644); err != nil {
		return "", fmt.Errorf(writingARecording, path, err)
	}
	return path, nil
}
