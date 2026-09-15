package fidelity

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const criteriaFolder = "criteria"

var CriteriaDirectory = filepath.Join(TestdataDirectory, criteriaFolder)

type Observation struct {
	Vendor      string            `json:"vendor"`
	Version     string            `json:"version"`
	Criterion   string            `json:"criterion"`
	Question    string            `json:"question"`
	Observed    map[string]string `json:"observed"`
	Established bool              `json:"established"`
}

func Observations(directory string) ([]Observation, error) {
	paths, err := filepath.Glob(filepath.Join(directory, "*"+RecordingExtension))
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf(noObservations, directory)
	}
	slices.Sort(paths)

	observed := make([]Observation, 0, len(paths))
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf(readingARecording, path, err)
		}

		var one Observation
		if err := json.Unmarshal(content, &one); err != nil {
			return nil, fmt.Errorf(readingARecording, path, err)
		}
		observed = append(observed, one)
	}
	return observed, nil
}

func KeepObservation(root string, made Observation) (string, error) {
	directory := filepath.Join(root, criteriaFolder)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return "", fmt.Errorf(writingARecording, directory, err)
	}

	encoded, err := json.MarshalIndent(made, "", "  ")
	if err != nil {
		return "", fmt.Errorf(writingARecording, made.Criterion, err)
	}

	path := filepath.Join(directory, ObservationName(made)+RecordingExtension)
	if err := os.WriteFile(path, append(encoded, '\n'), 0o644); err != nil {
		return "", fmt.Errorf(writingARecording, path, err)
	}
	return path, nil
}

func ObservationName(made Observation) string {
	return strings.Join([]string{made.Vendor, made.Criterion}, "-")
}
