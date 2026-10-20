package convention

import (
	"regexp"
	"slices"
)

const ClientBundle = "surface.js"

var (
	languageRegistered = regexp.MustCompile(`registerLanguage\("([^"]+)"`)
	aliasesRegistered  = regexp.MustCompile(`aliases:\[((?:"[^"]+",?)+)\]`)
	nameQuoted         = regexp.MustCompile(`"([^"]+)"`)
)

func LanguagesTheBundleRegisters(client string) []string {
	var named []string

	for _, found := range languageRegistered.FindAllStringSubmatch(client, -1) {
		named = append(named, found[1])
	}
	for _, aliases := range aliasesRegistered.FindAllStringSubmatch(client, -1) {
		for _, alias := range nameQuoted.FindAllStringSubmatch(aliases[1], -1) {
			named = append(named, alias[1])
		}
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
