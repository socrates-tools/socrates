package sanitise

import "strings"

var allowedElements = names(
	"p", "br", "hr", "strong", "b", "em", "i", "u", "s", "small", "sub", "sup",
	"code", "pre", "kbd", "samp", "span", "div", "section", "article",
	"h3", "h4", "h5", "h6", "blockquote", "figure", "figcaption",
	"ul", "ol", "li", "dl", "dt", "dd",
	"table", "thead", "tbody", "tfoot", "tr", "th", "td", "caption", "colgroup", "col",
	"svg", "g", "path", "rect", "circle", "ellipse", "line", "polyline", "polygon",
	"text", "tspan", "textPath", "defs", "marker", "symbol", "title", "desc",
	"linearGradient", "radialGradient", "stop", "clipPath", "mask", "pattern",
)

var droppedWithTheirContent = names(
	"script", "style", "template", "iframe", "object", "embed", "noscript",
	"form", "input", "button", "select", "option", "textarea", "label", "fieldset",
	"audio", "video", "canvas", "math", "portal", "dialog",
)

var voidElements = names("br", "hr", "col", "path", "rect", "circle", "ellipse", "line", "polyline", "polygon", "stop")

var allowedAttributes = names(
	"class", "id", "role", "aria-label", "aria-labelledby", "aria-describedby", "aria-hidden",
	"colspan", "rowspan", "scope", "datetime", "lang", "dir",
	"viewBox", "xmlns", "width", "height", "preserveAspectRatio",
	"d", "x", "y", "dx", "dy", "x1", "y1", "x2", "y2", "cx", "cy", "r", "rx", "ry", "points",
	"fill", "fill-opacity", "fill-rule", "stroke", "stroke-width", "stroke-opacity",
	"stroke-linecap", "stroke-linejoin", "stroke-dasharray", "stroke-dashoffset",
	"transform", "opacity", "offset", "stop-color", "stop-opacity",
	"text-anchor", "dominant-baseline", "alignment-baseline", "font-size", "font-family",
	"font-weight", "font-style", "letter-spacing", "marker-start", "marker-mid", "marker-end",
	"markerWidth", "markerHeight", "markerUnits", "refX", "refY", "orient", "overflow",
	"gradientUnits", "gradientTransform", "clip-path", "clip-rule", "mask", "patternUnits",
	"vector-effect", "paint-order", "shape-rendering",
)

var referencingAttributes = names(
	"fill", "stroke", "clip-path", "mask", "marker-start", "marker-mid", "marker-end",
)

func names(listed ...string) map[string]bool {
	held := make(map[string]bool, len(listed))
	for _, one := range listed {
		held[strings.ToLower(one)] = true
	}
	return held
}
