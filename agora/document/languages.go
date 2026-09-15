package document

import (
	"embed"
	"encoding/json"
	"slices"
)

const languagesFile = "languages.json"

//go:embed languages.json
var registered embed.FS

var languages = languagesRegistered()

func Languages() []string {
	return slices.Clone(languages)
}

func languagesRegistered() []string {
	held, err := registered.ReadFile(languagesFile)
	if err != nil {
		panic(err)
	}

	var named []string
	if err := json.Unmarshal(held, &named); err != nil {
		panic(err)
	}
	slices.Sort(named)
	return named
}
