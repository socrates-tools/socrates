package convention

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const ClientRefusalsGroup = "refusals"

var ClientMessages = []string{
	filepath.Join(ToolDirectory("agora"), "web", "packages", "surface", "messages.ts"),
}

func ClientWording(root string, modules []string) (Wording, error) {
	var wording Wording

	for _, module := range modules {
		held, err := os.ReadFile(filepath.Join(root, module))
		if err != nil {
			return Wording{}, fmt.Errorf("reading the client's messages module %s: %w", module, err)
		}

		tokens, err := lexTypeScript(string(held))
		if err != nil {
			return Wording{}, fmt.Errorf("reading the client's messages module %s: %w", module, err)
		}

		refusals, declared := declaredGroup(tokens, ClientRefusalsGroup)
		if !declared {
			return Wording{}, fmt.Errorf("the client's messages module %s declares no `export const %s = {...}`, so no refusal in it can be read", module, ClientRefusalsGroup)
		}

		at := func(token typeScriptToken) string {
			return fmt.Sprintf("%s:%d", module, token.Line)
		}

		for _, token := range stringsOf(tokens) {
			if strings.Contains(token.Text, EmDash) {
				wording.EmDashes = append(wording.EmDashes, at(token))
			}
		}
		for position, token := range refusals {
			if token.Kind != stringToken || isPropertyKey(refusals, position) {
				continue
			}
			if !namesANextAction(token.Text) {
				wording.Actionless = append(wording.Actionless, at(token))
			}
		}
	}

	slices.Sort(wording.EmDashes)
	slices.Sort(wording.Actionless)
	return wording, nil
}
