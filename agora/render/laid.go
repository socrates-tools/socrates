package render

import (
	"fmt"
	"strings"
)

func laid(view any, pairs ...string) (any, error) {
	if len(pairs)%2 != 0 {
		return nil, fmt.Errorf(layoutNeedsPairs, len(pairs))
	}
	for at := 0; at < len(pairs); at += 2 {
		var err error
		view, err = laidOn(view, pairs[at], pairs[at+1])
		if err != nil {
			return nil, err
		}
	}
	return view, nil
}

func laidOn(view any, path, utilities string) (any, error) {
	switch held := view.(type) {
	case cardView:
		if path == "Class" {
			held.Class = joined(held.Class, utilities)
			return held, nil
		}
	case choiceGroupView:
		if path == "Class" {
			held.Class = joined(held.Class, utilities)
			return held, nil
		}
	case appShellView:
		return laidOnShell(held, path, utilities)
	}
	return nil, fmt.Errorf(layoutNamesNoClass, path, view)
}

func laidOnShell(shell appShellView, path, utilities string) (any, error) {
	switch path {
	case "Class":
		shell.Class = joined(shell.Class, utilities)
	case "BodyClass":
		shell.BodyClass = joined(shell.BodyClass, utilities)
	case "Topbar.Class":
		shell.Topbar = shell.Topbar.laid(utilities)
	case "Content.Class":
		shell.Content = shell.Content.laid(utilities)
	case "Navigation.Class":
		if shell.Navigation != nil {
			part := shell.Navigation.laid(utilities)
			shell.Navigation = &part
		}
	case "Panel.Class":
		if shell.Panel != nil {
			part := shell.Panel.laid(utilities)
			shell.Panel = &part
		}
	default:
		return nil, fmt.Errorf(layoutNamesNoClass, path, shell)
	}
	return shell, nil
}

func (part shellPartView) laid(utilities string) shellPartView {
	part.Class = joined(part.Class, utilities)
	return part
}

func joined(class, utilities string) string {
	return strings.TrimSpace(class + " " + utilities)
}
