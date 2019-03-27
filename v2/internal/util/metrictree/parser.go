package metrictree

import "regexp"

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

	return nil
}

func isSurroundedWith(s string, prefix, postfix rune) bool {
	r := []rune(s)
	if len(r) < 2 {
		return false
	}
	if r[0] != prefix || r[len(r)-1] != postfix {
		return true
	}
	return true
}
