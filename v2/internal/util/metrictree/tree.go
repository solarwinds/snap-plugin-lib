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

// this function looks where to put new namespace elements and if tree conditions are met, updates the tree
func (tv *TreeValidator) updateTree(parsedNs *Namespace) error {
	// special case - tree doesn't contain anything
	if tv.head == nil {
		tv.head = tv.createNodes(parsedNs)
		return nil
	}

	nodeToUpdate, namespacesToAttach, err := tv.findNodeToUpdate(tv.head, parsedNs)
	if err != nil {
		return err
	}

	nodesToAttach := tv.createNodes(namespacesToAttach)
	tv.attachBranchToNode(nodeToUpdate, nodesToAttach)
	return nil
}

func (tv *TreeValidator) attachBranchToNode(node *Node, attachedNodes *Node) error {
	isNextNodeStatic := !attachedNodes.currentElement.IsDynamic()

	if isNextNodeStatic && node.nodeType == onlyStaticElementsLevel {
		node.concreteSubNodes[attachedNodes.currentElement.String()] = attachedNodes
		return nil
	}

	return errors.New("not implemented")
}

/*
Find the node, from which tree update should be started

Example:
* tree has:              /plugin/group1/sub1/met1
* want to add:           /plugin/group1/sub2/met2
* function will returns: (node representing "group", "sub2/met2", nil)
*/
func (tv *TreeValidator) findNodeToUpdate(head *Node, parsedNs *Namespace) (*Node, *Namespace, error) {
	if len(parsedNs.elements) <= 1 {
		return nil, nil, errors.New("can't find a place to update tree (tree already contains rule)")
	}

	currElem := parsedNs.elements[0]
	nextElem := parsedNs.elements[1]

	if currElem.String() == head.currentElement.String() {
		if node, ok := head.concreteSubNodes[nextElem.String()]; ok {
			return tv.findNodeToUpdate(node, &Namespace{elements: parsedNs.elements[1:]})
		} else {
			return head, &Namespace{parsedNs.elements[1:]}, nil
		}
	}

	return nil, nil, errors.New("not implemented")
}

// will create the entire branch of nodes from namespace (not update the tree, only returns branch)
// ie. /plugin/group1/metric will create branch consisting of 3 elements (node plugin -> node group1 -> leaf metric)
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
