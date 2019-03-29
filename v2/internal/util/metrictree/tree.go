package metrictree

import (
	"errors"
)

const (
	_ = iota
	metricDefinitionStrategy
	metricFilteringStrategy
)

const (
	invalidElementLevel = iota
	onlyStaticElementsLevel
	onlyDynamicElementsLevel
	leafLevel
)

type TreeValidator struct {
	strategy int
	head     *Node
}

// todo: Node may be a namespaceElement probably
type Node struct {
	currentElement namespaceElement
	nodeType       int

	concreteSubNodes map[string]*Node
	regexSubNodes    []*Node
}

func NewMetricDefinition() *TreeValidator {
	return &TreeValidator{
		strategy: metricDefinitionStrategy,
	}
}

func NewMetricFilter() *TreeValidator {
	return &TreeValidator{
		strategy: metricFilteringStrategy,
	}
}

func (tv *TreeValidator) AddRule(ns string) error {
	parsedNs, err := ParseNamespace(ns)
	if err != nil {
		return err
	}

	return tv.add(parsedNs)
}

func (tv *TreeValidator) IsValid(ns string) bool {
	return true
}

func (tv *TreeValidator) add(parsedNs *Namespace) error {
	switch tv.strategy {
	case metricDefinitionStrategy:
		if !parsedNs.isUsableForDefinition() {
			return errors.New("can't add rule")
		}
	case metricFilteringStrategy:
		if !parsedNs.isUsableForSelection() {
			return errors.New("can't add rule")
		}
	default:
		panic("invalid strategy")
	}

	return tv.updateTree(parsedNs)
}

// this function looks where to put new namespace elements and if tree conditions are met, update the tree
func (tv *TreeValidator) updateTree(parsedNs *Namespace) error {
	// special case - tree doesn't contain anything
	if tv.head == nil {
		tv.head = tv.createNodes(parsedNs)
		return nil
	}

	return errors.New("not implemented")
}

// will create the entire branch (each level is a namespace element). Return the first namespace element that are not yet part of the tree.
func (tv *TreeValidator) createNodes(ns *Namespace) *Node {
	if len(ns.elements) == 0 {
		return nil
	}
	if len(ns.elements) == 1 {
		return &Node{
			currentElement:   ns.elements[0],
			concreteSubNodes: nil,
			regexSubNodes:    nil,
			nodeType:         leafLevel,
		}
	}

	currNode := &Node{currentElement: ns.elements[0]}
	nextNode := tv.createNodes(&Namespace{elements: ns.elements[1:]})

	if !nextNode.currentElement.IsDynamic() {
		currNode.nodeType = onlyStaticElementsLevel
	} else {
		currNode.nodeType = onlyDynamicElementsLevel
	}

	if nextNode.currentElement.HasRegexp() {
		currNode.regexSubNodes = []*Node{nextNode}
	} else {
		currNode.concreteSubNodes = map[string]*Node{nextNode.currentElement.String(): nextNode}
	}

	return currNode
}
