package stylesheet

import "strings"

var allowedAtRules = named("@media", "@supports")

var allowedProperties = named(
	"display", "visibility", "position", "top", "right", "bottom", "left", "inset",
	"inset-block", "inset-inline", "z-index", "float", "clear", "box-sizing", "isolation",
	"width", "min-width", "max-width", "height", "min-height", "max-height", "aspect-ratio",
	"margin", "margin-top", "margin-right", "margin-bottom", "margin-left",
	"margin-block", "margin-inline",
	"padding", "padding-top", "padding-right", "padding-bottom", "padding-left",
	"padding-block", "padding-inline",
	"gap", "row-gap", "column-gap",
	"flex", "flex-direction", "flex-wrap", "flex-flow", "flex-grow", "flex-shrink", "flex-basis",
	"justify-content", "justify-items", "justify-self",
	"align-content", "align-items", "align-self", "order",
	"place-content", "place-items", "place-self",
	"grid", "grid-template", "grid-template-columns", "grid-template-rows", "grid-template-areas",
	"grid-area", "grid-column", "grid-column-start", "grid-column-end",
	"grid-row", "grid-row-start", "grid-row-end",
	"grid-auto-flow", "grid-auto-columns", "grid-auto-rows",
	"color", "opacity", "background", "background-color", "background-image",
	"background-position", "background-size", "background-repeat", "background-clip", "background-origin",
	"border", "border-width", "border-style", "border-color", "border-radius",
	"border-top", "border-right", "border-bottom", "border-left",
	"border-top-width", "border-right-width", "border-bottom-width", "border-left-width",
	"border-top-style", "border-right-style", "border-bottom-style", "border-left-style",
	"border-top-color", "border-right-color", "border-bottom-color", "border-left-color",
	"border-top-left-radius", "border-top-right-radius",
	"border-bottom-left-radius", "border-bottom-right-radius",
	"border-collapse", "border-spacing", "outline", "outline-width", "outline-style",
	"outline-color", "outline-offset", "box-shadow",
	"font", "font-family", "font-size", "font-weight", "font-style", "font-variant",
	"font-stretch", "font-variant-numeric", "font-feature-settings",
	"line-height", "letter-spacing", "word-spacing", "text-align", "text-align-last",
	"text-decoration", "text-decoration-line", "text-decoration-color",
	"text-decoration-style", "text-decoration-thickness", "text-underline-offset",
	"text-transform", "text-indent", "text-overflow", "text-shadow", "text-wrap",
	"white-space", "word-break", "overflow-wrap", "hyphens", "vertical-align",
	"direction", "unicode-bidi", "tab-size", "writing-mode", "content", "quotes",
	"list-style", "list-style-type", "list-style-position",
	"overflow", "overflow-x", "overflow-y", "table-layout", "caption-side", "empty-cells",
	"transform", "transform-origin", "rotate", "scale", "translate", "perspective",
	"fill", "fill-opacity", "fill-rule", "stroke", "stroke-width", "stroke-opacity",
	"stroke-linecap", "stroke-linejoin", "stroke-dasharray", "stroke-dashoffset",
	"clip-path", "clip-rule", "mask", "paint-order", "shape-rendering", "vector-effect",
	"stop-color", "stop-opacity", "text-anchor", "dominant-baseline", "marker-start",
	"marker-mid", "marker-end", "color-scheme", "accent-color", "caret-color", "cursor",
)

var allowedFunctions = named(
	"var", "calc", "min", "max", "clamp",
	"rgb", "rgba", "hsl", "hsla", "hwb", "lab", "lch", "oklab", "oklch",
	"color", "color-mix", "light-dark", "linear-gradient", "radial-gradient",
	"conic-gradient", "repeating-linear-gradient", "repeating-radial-gradient",
	"repeating-conic-gradient", "translate", "translateX", "translateY", "translate3d",
	"rotate", "rotateX", "rotateY", "rotateZ", "rotate3d", "scale", "scaleX", "scaleY",
	"scale3d", "skew", "skewX", "skewY", "matrix", "matrix3d", "perspective",
	"cubic-bezier", "steps", "minmax", "repeat", "fit-content", "polygon", "inset",
	"circle", "ellipse", "path", "format", "counter", "counters", "cross-fade",
	"selector", "supports", "not", "is", "where", "nth-child", "nth-last-child",
	"nth-of-type", "nth-last-of-type", "lang", "dir",
)

var refusedPseudo = named("has", "host", "host-context", "part", "slotted")

var theCallersRoot = named("html", "body")

const rootPseudo = "root"

func named(listed ...string) map[string]bool {
	held := make(map[string]bool, len(listed))
	for _, one := range listed {
		held[strings.ToLower(one)] = true
	}
	return held
}
