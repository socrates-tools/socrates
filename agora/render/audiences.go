package render

type audienceStated struct {
	tone  string
	icon  string
	title string
	hint  string
}

var audiencesSocratesReviewWrites = map[string]audienceStated{
	"needs a person": {tone: "call", icon: "user", title: "Your call", hint: "Only a person can decide this. The reviewer will not answer it."},
	"can be patched": {tone: "neutral", icon: "wrench", title: "Author can patch", hint: "The session that wrote the change can fix this without you."},
	"informational":  {tone: "neutral", icon: "eye", title: "For information", hint: "Nothing to decide."},
}
