package fidelity

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

const untypedEvent = "(untyped)"

type Shape struct {
	Type   string   `json:"type"`
	Fields []string `json:"fields"`
}

type Difference struct {
	Type  string
	Field string
	Why   string
}

func (d Difference) String() string {
	if d.Field == "" {
		return fmt.Sprintf(typeDiffers, d.Type, d.Why)
	}
	return fmt.Sprintf(fieldDiffers, d.Type, d.Field, d.Why)
}

func ShapesOf(stream string) ([]Shape, error) {
	fields := map[string]map[string]bool{}

	for line := range strings.SplitSeq(stream, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		var read map[string]any
		if err := json.Unmarshal([]byte(trimmed), &read); err != nil {
			return nil, fmt.Errorf(notAnEvent, elide(trimmed), err)
		}

		kind, named := read["type"].(string)
		if !named || kind == "" {
			kind = untypedEvent
		}
		if fields[kind] == nil {
			fields[kind] = map[string]bool{}
		}
		collect("", read, fields[kind])
	}

	shapes := make([]Shape, 0, len(fields))
	for kind, named := range fields {
		shape := Shape{Type: kind, Fields: make([]string, 0, len(named))}
		for field := range named {
			shape.Fields = append(shape.Fields, field)
		}
		slices.Sort(shape.Fields)
		shapes = append(shapes, shape)
	}
	slices.SortFunc(shapes, func(one, other Shape) int {
		return strings.Compare(one.Type, other.Type)
	})
	return shapes, nil
}

func collect(prefix string, value any, into map[string]bool) {
	switch held := value.(type) {
	case map[string]any:
		for field, inner := range held {
			path := field
			if prefix != "" {
				path = prefix + "." + field
			}
			into[path] = true
			collect(path, inner, into)
		}
	case []any:
		for _, element := range held {
			collect(prefix+"[]", element, into)
		}
	}
}

func Compare(recorded, produced []Shape) []Difference {
	var differences []Difference

	for _, simulated := range produced {
		carried, found := shapeOf(recorded, simulated.Type)
		if !found {
			differences = append(differences, Difference{Type: simulated.Type, Why: noSuchEventRecorded})
			continue
		}
		for _, field := range simulated.Fields {
			if !slices.Contains(carried.Fields, field) {
				differences = append(differences, Difference{Type: simulated.Type, Field: field, Why: notInTheRecording})
			}
		}
	}

	return differences
}

func shapeOf(shapes []Shape, kind string) (Shape, bool) {
	for _, shape := range shapes {
		if shape.Type == kind {
			return shape, true
		}
	}
	return Shape{}, false
}

func elide(line string) string {
	const limit = 120
	if len(line) <= limit {
		return line
	}
	return line[:limit] + "..."
}
