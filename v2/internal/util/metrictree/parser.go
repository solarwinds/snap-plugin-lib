package metrictree

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

const (
	nsSeparator                  = "/"
	regexBeginIndicator          = '{'
	regexEndIndicator            = '}'
	staticAnyMatcher             = '*'
	staticRecursiveAnyMatcher    = "**"
	dynamicElementBeginIndicator = '['
	dynamicElementEndIndicator   = ']'
	dynamicElementEqualIndicator = '='
)

const minNamespaceElements = 2

// Parsing whole selector (ie. "/plugin/[group={reg}]/group2/metric1) into smaller elements
func ParseNamespace(s string, isFilter bool) (*Namespace, error) {
	ns := &Namespace{}
	splitNs := strings.Split(s, nsSeparator)
	if len(splitNs)-1 < minNamespaceElements {
		return nil, fmt.Errorf("namespace doesn't contain valid numbers of elements (min. %d)", minNamespaceElements)
	}
	if splitNs[0] != "" {
		return nil, fmt.Errorf("namespace should start with '%s'", nsSeparator)
	}

	for i, nsElem := range splitNs[1:] {
		parsedEl, err := parseNamespaceElement(nsElem, isFilter)
		if err != nil {
			return nil, fmt.Errorf("can't parse namespace (%s), error at index %d: %s", s, i, err)
		}
		if _, ok := parsedEl.(*staticRecursiveAnyElement); ok && i != len(splitNs[1:])-1 {
			return nil, fmt.Errorf("recurive any-matcher (**) can be placed only as the last element")
		}
		ns.elements = append(ns.elements, parsedEl)
	}

	return ns, nil
}

// Parsing single selector (ie. [group={reg}])
func parseNamespaceElement(s string, isFilter bool) (namespaceElement, error) {
	if isSurroundedWith(s, dynamicElementBeginIndicator, dynamicElementEndIndicator) { // is it group []?
		dynElem := s[1 : len(s)-1]
		eqIndex := strings.Index(dynElem, string(dynamicElementEqualIndicator))

		if eqIndex != -1 { // is it group with value [group=id]
			if isIndexInTheMiddle(eqIndex, s) {
				groupName := dynElem[0:eqIndex]
				groupValue := dynElem[eqIndex+1:]

				if !isValidIdentifier(groupName) {
					return nil, fmt.Errorf("invalid character(s) used for group name [%s]", groupName)
				}

				if isSurroundedWith(groupValue, regexBeginIndicator, regexEndIndicator) { // is it group value as regex [group={regex}]
					regexStr := groupValue[1 : len(groupValue)-1]
					r, err := regexp.Compile(regexStr)
					if err != nil {
						return nil, fmt.Errorf("invalid regular expression (%s): %s", regexStr, err)
					}
					return newDynamicRegexpElement(groupName, r), nil
				}

				if isValidIdentifier(groupValue) {
					return newDynamicSpecificElement(groupName, groupValue), nil
				}

				return nil, fmt.Errorf("invalid character(s) used for group value [%s]", groupValue)
			}

			return nil, fmt.Errorf("invalid group with value (%s)", dynElem)
		}

		if isValidIdentifier(dynElem) {
			return newDynamicAnyElement(dynElem), nil
		}

		return nil, fmt.Errorf("invalid character(s) used for group value [%s]", dynElem)
	}

	if isSurroundedWith(s, regexBeginIndicator, regexEndIndicator) { // is it {regex}
		regexStr := s[1 : len(s)-1]
		r, err := regexp.Compile(regexStr)
		if err != nil {
			return nil, fmt.Errorf("invalid regular expression (%s): %s", regexStr, err)
		}
		return newStaticRegexpElement(r), nil
	}

	if s == string(staticRecursiveAnyMatcher) { // is it **
		return newStaticRecursiveAnyElement(), nil
	}

	if s == string(staticAnyMatcher) { // is it *
		return newStaticAnyElement(), nil
	}

	if isFilter {
		return newStaticSpecificAcceptingGroupElement(s), nil
	}

	if isValidIdentifier(s) { // is it static element ie. metric
		return newStaticSpecificElement(s), nil
	}

	return nil, fmt.Errorf("invalid character(s) used for element [%s]", s)
}

/*****************************************************************************/

func isIndexInTheMiddle(idx int, s string) bool {
	return len(s) >= 3 && idx > 0 && idx < len(s)-1
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
	if len(s) == 0 {
		return false
	}

	for _, el := range s {
		if !unicode.IsLetter(el) && !unicode.IsDigit(el) && el != '-' && el != '_' {
			return false
		}
	}

	return true
}
