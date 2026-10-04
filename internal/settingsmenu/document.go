package settingsmenu

import (
	"bytes"
	"regexp"
	"strings"

	"go.yaml.in/yaml/v3"
)

const (
	blankMarker        = "#socrates:blank-line"
	keySeparator       = "."
	indentStep         = 2
	longestImplicitKey = 1024
)

var (
	blankLine   = regexp.MustCompile(`(?m)^[ \t]*\r?\n`)
	blankMark   = regexp.MustCompile(`(?m)^[ \t]*` + regexp.QuoteMeta(blankMarker) + `[ \t]*\r?\n`)
	emptyOnly   = regexp.MustCompile(`(?m)^\{\}\r?\n`)
	explicitKey = regexp.MustCompile(`(?m)^([ \t]*)\? ([^\r\n|>][^\r\n]*)\r?\n([ \t]*):`)
)

type document struct {
	root yaml.Node
}

func documentOf(content []byte) (*document, error) {
	held := blankLine.ReplaceAllString(string(content), blankMarker+"\n")

	var root yaml.Node
	if err := yaml.Unmarshal([]byte(held), &root); err != nil {
		return nil, err
	}
	return &document{root: root}, nil
}

func keyPath(key string) []string {
	return strings.Split(key, keySeparator)
}

func (d *document) set(key []string, value string) {
	holder := d.holderOf(key, true)
	leaf := key[len(key)-1]

	if stated := valueUnder(holder, leaf); stated != nil {
		stated.Kind = yaml.ScalarNode
		stated.Tag = ""
		stated.Style = 0
		stated.Value = value
		stated.Content = nil
		return
	}
	holder.Content = append(holder.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: leaf},
		&yaml.Node{Kind: yaml.ScalarNode, Value: value})
}

func (d *document) remove(key []string) {
	pruneFrom(d.top(), key)
}

func (d *document) states(key []string) bool {
	holder := d.holderOf(key, false)
	return holder != nil && valueUnder(holder, key[len(key)-1]) != nil
}

func (d *document) render() ([]byte, error) {
	if d.root.Kind == 0 {
		return nil, nil
	}

	var out bytes.Buffer
	encoder := yaml.NewEncoder(&out)
	encoder.SetIndent(indentStep)
	if err := encoder.Encode(&d.root); err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}

	written := blankMark.ReplaceAllString(out.String(), newline)
	written = onOneLine(written)
	if d.namesNothing() {
		written = emptyOnly.ReplaceAllString(written, "")
	}
	if strings.TrimSpace(written) == "" {
		return nil, nil
	}
	return []byte(written), nil
}

func onOneLine(written string) string {
	return explicitKey.ReplaceAllStringFunc(written, func(held string) string {
		parts := explicitKey.FindStringSubmatch(held)
		indent, key, valueIndent := parts[1], parts[2], parts[3]
		if indent != valueIndent || len(key) > longestImplicitKey {
			return held
		}
		return indent + key + ":"
	})
}

func (d *document) namesNothing() bool {
	return len(d.root.Content) == 0 ||
		d.root.Content[0].Kind != yaml.MappingNode ||
		len(d.root.Content[0].Content) == 0
}

func (d *document) top() *yaml.Node {
	if d.root.Kind == 0 {
		d.root = yaml.Node{Kind: yaml.DocumentNode}
	}
	if len(d.root.Content) == 0 {
		d.root.Content = []*yaml.Node{{Kind: yaml.MappingNode}}
	}
	if d.root.Content[0].Kind != yaml.MappingNode {
		d.root.Content[0] = &yaml.Node{Kind: yaml.MappingNode}
	}
	return d.root.Content[0]
}

func (d *document) holderOf(key []string, building bool) *yaml.Node {
	holder := d.top()

	for _, segment := range key[:len(key)-1] {
		nested := valueUnder(holder, segment)
		if nested == nil || nested.Kind != yaml.MappingNode {
			if !building {
				return nil
			}
			nested = mappingUnder(holder, segment)
		}
		holder = nested
	}
	return holder
}

func valueUnder(holder *yaml.Node, name string) *yaml.Node {
	if holder == nil {
		return nil
	}
	for position := 0; position+1 < len(holder.Content); position += 2 {
		if holder.Content[position].Value == name {
			return holder.Content[position+1]
		}
	}
	return nil
}

func mappingUnder(holder *yaml.Node, name string) *yaml.Node {
	if stated := valueUnder(holder, name); stated != nil {
		stated.Kind = yaml.MappingNode
		stated.Tag = ""
		stated.Style = 0
		stated.Value = ""
		return stated
	}

	nested := &yaml.Node{Kind: yaml.MappingNode}
	holder.Content = append(holder.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: name}, nested)
	return nested
}

func pruneFrom(holder *yaml.Node, segments []string) bool {
	for position := 0; position+1 < len(holder.Content); position += 2 {
		if holder.Content[position].Value != segments[0] {
			continue
		}
		if len(segments) > 1 {
			nested := holder.Content[position+1]
			if nested.Kind != yaml.MappingNode || !pruneFrom(nested, segments[1:]) {
				return false
			}
			if len(nested.Content) > 0 {
				return true
			}
		}
		holder.Content = append(holder.Content[:position], holder.Content[position+2:]...)
		return true
	}
	return false
}
