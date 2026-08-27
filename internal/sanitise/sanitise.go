package sanitise

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func Fragment(supplied, prefix string) string {
	if strings.TrimSpace(supplied) == "" {
		return ""
	}

	holding := &html.Node{Type: html.ElementNode, Data: "div", DataAtom: atom.Div}
	parsed, err := html.ParseFragment(strings.NewReader(supplied), holding)
	if err != nil {
		return ""
	}

	var reduced strings.Builder
	for _, node := range parsed {
		reduce(&reduced, node, prefix)
	}
	return reduced.String()
}

func reduce(reduced *strings.Builder, node *html.Node, prefix string) {
	switch node.Type {
	case html.TextNode:
		reduced.WriteString(html.EscapeString(node.Data))
		return
	case html.ElementNode:
	default:
		return
	}

	name := strings.ToLower(node.Data)
	if droppedWithTheirContent[name] {
		return
	}
	if !allowedElements[name] {
		reduceChildren(reduced, node, prefix)
		return
	}

	reduced.WriteString("<")
	reduced.WriteString(node.Data)
	for _, attribute := range node.Attr {
		if value, held := kept(attribute, prefix); held {
			reduced.WriteString(" ")
			reduced.WriteString(attribute.Key)
			reduced.WriteString(`="`)
			reduced.WriteString(html.EscapeString(value))
			reduced.WriteString(`"`)
		}
	}
	if voidElements[name] {
		reduced.WriteString("/>")
		return
	}
	reduced.WriteString(">")
	reduceChildren(reduced, node, prefix)
	reduced.WriteString("</")
	reduced.WriteString(node.Data)
	reduced.WriteString(">")
}

func reduceChildren(reduced *strings.Builder, node *html.Node, prefix string) {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		reduce(reduced, child, prefix)
	}
}

func kept(attribute html.Attribute, prefix string) (string, bool) {
	if attribute.Namespace != "" {
		return "", false
	}
	name := strings.ToLower(attribute.Key)
	if !allowedAttributes[name] {
		return "", false
	}
	if namingAttributes[name] {
		return prefixedNames(attribute.Val, prefix), true
	}
	if referencingAttributes[name] {
		return prefixedReference(attribute.Val, prefix)
	}
	return attribute.Val, true
}

func prefixedNames(value, prefix string) string {
	var prefixed []string
	for named := range strings.FieldsSeq(value) {
		prefixed = append(prefixed, prefix+named)
	}
	return strings.Join(prefixed, " ")
}

func prefixedReference(value, prefix string) (string, bool) {
	trimmed := strings.TrimSpace(value)
	if !strings.HasPrefix(strings.ToLower(trimmed), reference) {
		if strings.Contains(strings.ToLower(trimmed), "://") {
			return "", false
		}
		return value, true
	}

	inside := trimmed[len(reference):]
	if !strings.HasSuffix(inside, ")") {
		return "", false
	}
	inside = strings.TrimSpace(strings.TrimSuffix(inside, ")"))
	for _, quote := range []string{`"`, "'"} {
		if len(inside) >= 2 && strings.HasPrefix(inside, quote) && strings.HasSuffix(inside, quote) {
			inside = inside[1 : len(inside)-1]
		}
	}

	named, fragmentary := strings.CutPrefix(inside, "#")
	if !fragmentary || named == "" || strings.ContainsAny(named, `"'() `) {
		return "", false
	}
	return reference + "#" + prefix + named + ")", true
}
