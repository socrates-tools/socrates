package surface

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/socrates-tools/socrates/agora/document"
)

const ClientTypesFile = "agora/web/packages/surface/document.generated.ts"

func ClientTypes() ([]byte, error) {
	writing := &typeScript{enumerations: map[reflect.Type][]string{
		reflect.TypeFor[document.Mode](): modeValues(),
	}}
	writing.queue = append(writing.queue, reflect.TypeFor[Served]())

	var written strings.Builder
	for at := 0; at < len(writing.queue); at++ {
		declared, err := writing.declare(writing.queue[at])
		if err != nil {
			return nil, err
		}
		if at > 0 {
			written.WriteString("\n")
		}
		written.WriteString(declared)
	}
	return []byte(written.String()), nil
}

func modeValues() []string {
	var values []string
	for _, mode := range document.Modes() {
		values = append(values, string(mode))
	}
	return values
}

type typeScript struct {
	queue        []reflect.Type
	enumerations map[reflect.Type][]string
}

func (w *typeScript) declare(declared reflect.Type) (string, error) {
	if values, enumerated := w.enumerations[declared]; enumerated {
		quoted := make([]string, 0, len(values))
		for _, value := range values {
			quoted = append(quoted, "'"+value+"'")
		}
		return fmt.Sprintf(typeScriptUnion, declared.Name(), strings.Join(quoted, typeScriptUnionJoint)), nil
	}

	var declaration strings.Builder
	fmt.Fprintf(&declaration, typeScriptInterfaceOpening, declared.Name())
	for at := range declared.NumField() {
		field := declared.Field(at)
		name, optional, written := jsonKeyOf(field)
		if !written {
			continue
		}
		kind, err := w.kindOf(field.Type)
		if err != nil {
			return "", fmt.Errorf(typeScriptHasNoKindFor, declared.Name(), field.Name, err)
		}
		marker := ""
		if optional {
			marker = "?"
		}
		fmt.Fprintf(&declaration, typeScriptProperty, name, marker, kind)
	}
	declaration.WriteString("}\n")
	return declaration.String(), nil
}

func jsonKeyOf(field reflect.StructField) (name string, optional bool, written bool) {
	tag := field.Tag.Get("json")
	if tag == "" || tag == "-" {
		return "", false, false
	}
	name, options, _ := strings.Cut(tag, ",")
	return name, strings.Contains(options, "omitempty"), true
}

func (w *typeScript) kindOf(held reflect.Type) (string, error) {
	if _, enumerated := w.enumerations[held]; enumerated {
		return w.named(held), nil
	}
	switch held.Kind() {
	case reflect.String:
		return "string", nil
	case reflect.Bool:
		return "boolean", nil
	case reflect.Int:
		return "number", nil
	case reflect.Slice:
		element, err := w.kindOf(held.Elem())
		if err != nil {
			return "", err
		}
		return element + "[]", nil
	case reflect.Map:
		if held.Key().Kind() != reflect.String {
			return "", fmt.Errorf(typeScriptKeysAreNotStrings, held)
		}
		element, err := w.kindOf(held.Elem())
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(typeScriptRecord, element), nil
	case reflect.Struct:
		return w.named(held), nil
	default:
		return "", fmt.Errorf(typeScriptHasNoKindForGo, held)
	}
}

func (w *typeScript) named(held reflect.Type) string {
	for _, queued := range w.queue {
		if queued == held {
			return held.Name()
		}
	}
	w.queue = append(w.queue, held)
	return held.Name()
}
