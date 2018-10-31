package gen

import "strings"

var (
	nameWordsSeparator         = []string{".", " ", "-", "_"}
	goNameCapitalizedException = []string{
		"ID",
		"JSON",
		"HTTP",
	}
)

func toGoFieldName(s string) string {
	result := new(strings.Builder)
eachPart:
	for _, part := range splitByOneOf(s, nameWordsSeparator...) {
		partUpperCase := strings.ToUpper(part)
		for _, exception := range goNameCapitalizedException {
			if partUpperCase == exception {
				result.WriteString(partUpperCase)
				continue eachPart
			}
		}
		result.WriteString(strings.Title(part))
	}
	return result.String()
}

func splitByOneOf(s string, separators ...string) []string {
	lastSeparatorIndex := len(separators) - 1
	result := make([]string, 0)
	var splitter func(substrings []string, separatorIndex int)
	splitter = func(substrings []string, separatorIndex int) {
		for _, substring := range substrings {
			parts := strings.Split(substring, separators[separatorIndex])
			if separatorIndex < lastSeparatorIndex {
				splitter(parts, separatorIndex+1)
			} else {
				result = append(result, parts...)
			}
		}
	}
	splitter([]string{s}, 0)
	return result
}
