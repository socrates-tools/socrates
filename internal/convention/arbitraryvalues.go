package convention

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var ClassSourcesStatement = filepath.Join(ToolDirectory("agora"), "web", "packages", "ui", "theme.css")

var sourceStatement = regexp.MustCompile(`(?m)^@source\s+(not\s+)?(?:'([^']*)'|"([^"]*)")\s*;`)

var classCandidate = regexp.MustCompile("[^\\s\"'`]+")

var arbitraryValueWritten = []*regexp.Regexp{
	regexp.MustCompile(`\w-\[`),
	regexp.MustCompile(`\w-[\w-]*/\[`),
	regexp.MustCompile(`@\[`),
	regexp.MustCompile(`(?:^|:)!?\[(?:&|@|[A-Za-z-][\w-]*:)[^\]]*\]`),
	regexp.MustCompile(`\w-\((?:[a-z-]+:)?--`),
}

type classSourceStatement struct {
	pattern string
	not     bool
}

func ClassSourcesStated(root, statement string) ([]string, error) {
	held, err := os.ReadFile(filepath.Join(root, statement))
	if err != nil {
		return nil, err
	}

	var stated []classSourceStatement
	for _, match := range sourceStatement.FindAllStringSubmatch(string(held), -1) {
		pattern := match[2] + match[3]
		stated = append(stated, classSourceStatement{
			pattern: path.Join(filepath.ToSlash(filepath.Dir(statement)), pattern),
			not:     match[1] != "",
		})
	}
	if !slices.ContainsFunc(stated, func(each classSourceStatement) bool { return !each.not }) {
		return nil, fmt.Errorf("%s states no @source, so there is no class to read", filepath.ToSlash(statement))
	}

	var included, excluded []*regexp.Regexp
	var starts []string
	for _, each := range stated {
		expression := regexp.MustCompile("^" + globExpression(each.pattern) + "(?:/.*)?$")
		if each.not {
			excluded = append(excluded, expression)
			continue
		}
		included = append(included, expression)
		starts = append(starts, globStart(each.pattern))
	}

	var read []string
	for _, start := range slices.Compact(slices.Sorted(slices.Values(starts))) {
		err := filepath.WalkDir(filepath.Join(root, start), func(found string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			relative, relErr := filepath.Rel(root, found)
			if relErr != nil {
				return relErr
			}
			relative = filepath.ToSlash(relative)
			if slices.ContainsFunc(excluded, func(expression *regexp.Regexp) bool { return expression.MatchString(relative) }) {
				if entry.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if entry.IsDir() {
				return nil
			}
			if slices.ContainsFunc(included, func(expression *regexp.Regexp) bool { return expression.MatchString(relative) }) {
				read = append(read, relative)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	slices.Sort(read)
	return slices.Compact(read), nil
}

func ArbitraryValuesWrittenInClasses(root, statement string) ([]string, error) {
	sources, err := ClassSourcesStated(root, statement)
	if err != nil {
		return nil, err
	}

	var found []string
	for _, source := range sources {
		held, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(source)))
		if err != nil {
			return nil, err
		}
		for number, line := range strings.Split(string(held), "\n") {
			for _, candidate := range classCandidate.FindAllString(line, -1) {
				if arbitraryValueIn(candidate) {
					found = append(found, fmt.Sprintf("%s:%d writes the arbitrary value %q, and a class takes a token or a step of the scale, so add one to theme.css instead", source, number+1, candidate))
				}
			}
		}
	}
	return found, nil
}

func arbitraryValueIn(candidate string) bool {
	return slices.ContainsFunc(arbitraryValueWritten, func(written *regexp.Regexp) bool { return written.MatchString(candidate) })
}

func globStart(pattern string) string {
	segments := strings.Split(pattern, "/")
	var plain []string
	for _, segment := range segments {
		if strings.ContainsAny(segment, "*?[{") {
			break
		}
		plain = append(plain, segment)
	}
	return filepath.FromSlash(path.Join(plain...))
}

func globExpression(pattern string) string {
	var expression strings.Builder
	for position := 0; position < len(pattern); position++ {
		switch {
		case strings.HasPrefix(pattern[position:], "**/"):
			expression.WriteString("(?:.*/)?")
			position += 2
		case strings.HasPrefix(pattern[position:], "**"):
			expression.WriteString(".*")
			position++
		case pattern[position] == '*':
			expression.WriteString("[^/]*")
		case pattern[position] == '?':
			expression.WriteString("[^/]")
		default:
			expression.WriteString(regexp.QuoteMeta(string(pattern[position])))
		}
	}
	return expression.String()
}
