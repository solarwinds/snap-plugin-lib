package metrictree

import (
	"regexp"
	"strings"
)

const nsSeparator = '/'
const regexBeginIndicator = '{'
const regexEndIndicator = '}'
const staticAnyMatcher = '*'
const dynamicElementBeginIndicator = '['
const dynamicElementEndIndicator = ']'
const dynamicElementEqualIndicator = '='

type Namespace struct {
	elements []namespaceElement
}

// Parsing whole selector (ie. "/plugin/[group={reg}]/group2/metric1) into smaller elements
func ParseNamespaceElement(s string) namespaceElement {
	if isSurroundedWith(s, dynamicElementBeginIndicator, dynamicElementEndIndicator) {
		dynElem := s[1 : len(s)-1]
		eqIndex := strings.Index(dynElem, string(dynamicElementEqualIndicator))

		if eqIndex != -1 {
			if len(dynElem) >= 3 && eqIndex > 0 && eqIndex < len(dynElem)-1 {
				groupName := dynElem[0:eqIndex]
				groupValue := dynElem[eqIndex+1:]

				return newDynamicSpecificElement(groupName, groupValue)
			}
		}
	}

	if isSurroundedWith(s, regexBeginIndicator, regexEndIndicator) {
		regexStr := s[1 : len(s)-1]
		r, err := regexp.Compile(regexStr)
		if err != nil {
			// todo: log error
			return nil
		}
		return newStaticRegexpElement(r)
	}

	if s == string(staticAnyMatcher) {
		return newStaticAnyElement()
	}

	if isValidIdentifier(s) {
		return newStaticSpecificElement(s)
	}

	return nil
}

func isSurroundedWith(s string, prefix, postfix rune) bool {
	r := []rune(s)
	if len(r) < 2 {
		return false
	}
	if r[0] != prefix || r[len(r)-1] != postfix {
		return false
	}
	return true
}

func isValidIdentifier(s string) bool {
	return len(s) > 0 // todo: check is contains valid characters
}
