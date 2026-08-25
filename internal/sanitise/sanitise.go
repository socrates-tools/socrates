package sanitise

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func Fragment(supplied string) string {
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
		reduce(&reduced, node)
	}
	return reduced.String()
}

func reduce(reduced *strings.Builder, node *html.Node) {
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
		reduceChildren(reduced, node)
		return
	}

	reduced.WriteString("<")
	reduced.WriteString(node.Data)
	for _, attribute := range node.Attr {
		if kept(attribute) {
			reduced.WriteString(" ")
			reduced.WriteString(attribute.Key)
			reduced.WriteString(`="`)
			reduced.WriteString(html.EscapeString(attribute.Val))
			reduced.WriteString(`"`)
		}
	}
	if voidElements[name] {
		reduced.WriteString("/>")
		return
	}
	reduced.WriteString(">")
	reduceChildren(reduced, node)
	reduced.WriteString("</")
	reduced.WriteString(node.Data)
	reduced.WriteString(">")
}

func reduceChildren(reduced *strings.Builder, node *html.Node) {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		reduce(reduced, child)
	}
}

func kept(attribute html.Attribute) bool {
	if attribute.Namespace != "" {
		return false
	}
	name := strings.ToLower(attribute.Key)
	if !allowedAttributes[name] {
		return false
	}
	if referencingAttributes[name] {
		return referencesThePageItself(attribute.Val)
	}
	return true
}

func referencesThePageItself(value string) bool {
	trimmed := strings.ToLower(strings.TrimSpace(value))
	if !strings.HasPrefix(trimmed, "url(") {
		return !strings.Contains(trimmed, "://")
	}
	return strings.HasPrefix(trimmed, "url(#") || strings.HasPrefix(trimmed, `url("#`) || strings.HasPrefix(trimmed, "url('#")
}
