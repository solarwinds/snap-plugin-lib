/*
File contains implementation of a validation tree, which is used to efficiently validate and filter metrics
added during collection process.

Validation tree is created based on parsed namespace elements defined in namespace.go.

Let's consider that user have defined following metrics:
	/pl/gr1/sub1/m1
	/pl/gr1/sub2/m1
	/pl/gr1/sub3/m1
	/pl/gr1/sub3/m2
	/pl/gr2/[dyn1]/m1

Generated tree looks like:

                  n1(pl)
             /               \
         n2(gr1)            n3(gr2)
    /       |       \           \
n4(sub1) n5(sub2) n6(sub3)     n7([dyn1])
    |       |        |    \          |
  n8(m1)  n9(m1)  n10(m1) n11(m2)  n12(m2)

If user wanted to check if some metric is valid,
	ie. /pl/g1/sub3/m1
following steps are taken:
	1) Extract selector "/pl/gr1/sub3/m1" into array of elements: [pl, gr1, sub2, m1]
	2) Start traversing the tree from head (node n1), compare n1 matches to first element of selector (pl)
	3) Check if node n1 (representing "pl" element) contains sub-element representing "gr1" -> yes, it's n2
	4) Check if node n2 (representing "gr1" element) contains sub-element representing "sub3" -> yes, it's n6
	5) Check if node n6 (representing "sub3" element) contains sub-element representing "m1" -> yes, it's n10
	6) Check if node n10 is a leaf.

Steps 3-5 should be executed in O(1) since node contains map reference to its sub-nodes, so we don't need to
compare with "sub1" and "sub2" at that level.

Filtering tree looks very similar, although can contain more sophisticated nodes (like regular expression).
 */

package metrictree

import (
	"errors"
	"sort"
	"strings"
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
	mixedElementsLevel
	leafLevel
)

type TreeValidator struct {
	strategy       int            // used to distinguish between definition and filtering tree
	definitionTree *TreeValidator // used in filtering tree (reference to definition tree)

	head *Node
}

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

func NewMetricFilter(definitionTree *TreeValidator) *TreeValidator {
	return &TreeValidator{
		strategy:       metricFilteringStrategy,
		definitionTree: definitionTree,
	}
}

func (tv *TreeValidator) AddRule(ns string) error {
	parsedNs, err := ParseNamespace(ns, tv.strategy == metricFilteringStrategy)
	if err != nil {
		return err
	}

	return tv.add(parsedNs)
}

func (tv *TreeValidator) IsPartiallyValid(ns string) bool {
	isValid, _ := tv.isValid(ns, false)
	return isValid
}

func (tv *TreeValidator) IsValid(ns string) (bool, []string) {
	isValid, trace := tv.isValid(ns, true)
	return isValid, trace
}

func (tv *TreeValidator) HasRules() bool {
	return tv.head != nil
}

// second value indicated name of dynamic group for the element
// ["", "group", ""] indicated that 2nd parameter should be treated as dynamic
func (tv *TreeValidator) isValid(ns string, fullMatch bool) (bool, []string) {
	nsSep := strings.Split(ns, "/")[1:]
	groupIndicator := make([]string, len(nsSep))

	if tv.head == nil {
		return true, groupIndicator // special case - there is no rule defined so everything is valid
	}

	isValid := false
	var stackFromValidBranch []*Node

	tv.traverse(tv.head, nil, func(n *Node, stack []*Node) (bool, bool) {
		idx := len(stack) // len of stack indicated to which string element should we match
		if idx >= len(nsSep) {
			return false, true
		}
		if !n.currentElement.Match(nsSep[idx]) {
			return false, true
		}
		if _, ok := n.currentElement.(*staticRecursiveAnyElement); ok {
			isValid = true
			stackFromValidBranch = stack
			return false, false
		}

		if len(nsSep)-1 == idx {
			switch {
			case fullMatch && n.nodeType == leafLevel:
				isValid = true
				stackFromValidBranch = stack
				return false, false
			case !fullMatch:
				isValid = true
				return false, false
			default:
				return false, true
			}
		}

		return true, true
	})

	for idx, node := range stackFromValidBranch {
		if groupNode, ok := node.currentElement.(*dynamicAnyElement); ok {
			groupIndicator[idx] = groupNode.group
		}
	}

	return isValid, groupIndicator
}

func (tv *TreeValidator) ListRules() []string {
	if tv.head == nil {
		return []string{}
	}

	nsList := []string{}
	tv.traverse(tv.head, nil, func(n *Node, stack []*Node) (bool, bool) {
		if n.nodeType != leafLevel {
			return true, true
		}

		stack = append(stack, n)

		nsElems := []string{}
		for _, stackEl := range stack {
			nsElems = append(nsElems, stackEl.currentElement.String())
		}

		nsString := "/" + strings.Join(nsElems, "/")
		nsList = append(nsList, nsString)

		return false, true
	})

	sort.Strings(nsList)
	return nsList
}

// Define function executed when each node is reached during traverse.
// First return value indicates if traversing should be continued (go to next level) after processing this node (false ends traversing of current branch)
// Second return value indicates if traversing should be continued on other branches (false ends traversing of tree)
type traverseFn func(*Node, []*Node) (bool, bool)

// Traversing tree and executing function on each node.
func (tv *TreeValidator) traverse(n *Node, stack []*Node, fn traverseFn) bool {

	procBranch, procTree := fn(n, stack)
	if !procTree {
		return false
	}
	if !procBranch {
		return true
	}

	stack = append(stack, n)

	for _, subNode := range n.concreteSubNodes { // todo: optimalize O(n) into O(1)
		cont := tv.traverse(subNode, stack, fn)
		if !cont {
			return false
		}
	}

	for _, subNode := range n.regexSubNodes {
		cont := tv.traverse(subNode, stack, fn)
		if !cont {
			return false
		}
	}

	return true
}

func (tv *TreeValidator) add(parsedNs *Namespace) error {
	switch tv.strategy {
	case metricDefinitionStrategy:
		if !parsedNs.IsUsableForDefinition() {
			return errors.New("can't add rule - some namespace elements are not allowed in definition")
		}
	case metricFilteringStrategy:
		defPresent := tv.definitionTree.HasRules()
		if !parsedNs.IsUsableForFiltering(defPresent) {
			return errors.New("can't add rule - some namespace elements are not allowed in filtering when metric definition wasn't provided")
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
	return tv.attachBranchToNode(nodeToUpdate, nodesToAttach)
}

func (tv *TreeValidator) attachBranchToNode(node *Node, attachedNodes *Node) error {
	isNextNodeStatic := !attachedNodes.currentElement.IsDynamic()

	if node.nodeType == onlyStaticElementsLevel && !isNextNodeStatic {
		return errors.New("only static elements may added at current level")
	}

	if node.nodeType == onlyDynamicElementsLevel && !isNextNodeStatic {
		return errors.New("there can be only one dynamic element at current level")
	}

	if !attachedNodes.currentElement.HasRegexp() {
		node.concreteSubNodes[attachedNodes.currentElement.String()] = attachedNodes
	} else {
		node.regexSubNodes = append(node.regexSubNodes, attachedNodes)
	}
	return nil
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

	currNode := &Node{
		currentElement:   ns.elements[0],
		concreteSubNodes: map[string]*Node{},
		regexSubNodes:    []*Node{},
	}
	nextNode := tv.createNodes(&Namespace{elements: ns.elements[1:]})

	if tv.strategy == metricFilteringStrategy {
		currNode.nodeType = mixedElementsLevel
	}
	if tv.strategy == metricDefinitionStrategy {
		if !nextNode.currentElement.IsDynamic() {
			currNode.nodeType = onlyStaticElementsLevel
		} else {
			currNode.nodeType = onlyDynamicElementsLevel
		}
	}

	if nextNode.currentElement.HasRegexp() {
		currNode.regexSubNodes = []*Node{nextNode}
	} else {
		currNode.concreteSubNodes = map[string]*Node{nextNode.currentElement.String(): nextNode}
	}

	return currNode
}
