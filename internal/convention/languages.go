package convention

import (
	"regexp"
	"slices"
)

const ClientBundle = "surface.js"

var languageRegistered = regexp.MustCompile(`registerLanguage\("([^"]+)"`)

func LanguagesTheBundleRegisters(client string) []string {
	var named []string

	for _, found := range languageRegistered.FindAllStringSubmatch(client, -1) {
		named = append(named, found[1])
	}
	slices.Sort(named)
	return slices.Compact(named)
}

func LanguagesOnOneSideOnly(stated, registered []string) (statedOnly, registeredOnly []string) {
	for _, named := range stated {
		if !slices.Contains(registered, named) {
			statedOnly = append(statedOnly, named)
		}
	}
	for _, named := range registered {
		if !slices.Contains(stated, named) {
			registeredOnly = append(registeredOnly, named)
		}
	}
	return statedOnly, registeredOnly
}
