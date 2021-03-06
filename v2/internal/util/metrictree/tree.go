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

/*
 Copyright (c) 2020 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package metrictree

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

const (
	_ TreeStrategy = iota
	metricDefinitionStrategy
	metricFilteringStrategy
)

const ( // nodeType const
	_ = iota
	onlyStaticElementsLevel
	onlyDynamicElementsLevel
	mixedElementsLevel
	leafLevel
)

type TreeStrategy int

type TreeValidator struct {
	strategy       TreeStrategy   // used to distinguish between definition and filtering tree
	definitionTree *TreeValidator // used in filtering tree (reference to definition tree)

	head *Node
}

type Node struct {
	parent         *Node
	currentElement namespaceElement
	nodeType       int
	level          int // horizontal position in tree (starts from 0)

	subNodes map[string]*Node
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

	switch tv.strategy {
	case metricDefinitionStrategy:
		if !parsedNs.IsUsableForDefinition() {
			return fmt.Errorf("can't add rule (%s) - some namespace elements are not allowed in definition", ns)
		}
	case metricFilteringStrategy:
		defPresent := tv.definitionTree.HasRules()
		if !parsedNs.IsUsableForFiltering(defPresent) {
			return fmt.Errorf("can't add rule (%s) - some namespace elements are not allowed in filtering when metric definition wasn't provided", ns)
		}

		if !tv.definitionTree.IsCompatible(ns) {
			return fmt.Errorf("can't add rule (%s) - not compatible with any metric definition", ns)
		}
	default:
		panic("invalid strategy")
	}

	return tv.updateTree(parsedNs)
}

func (tv *TreeValidator) IsPartiallyValid(ns string) bool {
	isValid, _ := tv.isValid(ns, false, false)
	return isValid
}

func (tv *TreeValidator) IsValid(ns string) (bool, []string) {
	isValid, trace := tv.isValid(ns, true, false)
	return isValid, trace
}

func (tv *TreeValidator) IsCompatible(ns string) bool {
	isCompatible, _ := tv.isValid(ns, false, true)
	return isCompatible
}

func (tv *TreeValidator) HasRules() bool {
	return tv.head != nil
}

func (tv *TreeValidator) isValid(ns string, fullMatch bool, compatibilityMode bool) (bool, []string) {
	if compatibilityMode && tv.strategy != metricDefinitionStrategy {
		panic("compatibilityMode can be only used for definition tree")
	}

	nsElems, _, err := SplitNamespace(ns)
	if err != nil {
		return false, nil
	}

	nsElems = nsElems[1:]
	groupIndicator := make([]string, len(nsElems))

	// special case - no rules defined - everything is valid and there are no groups (2nd param contains empty strings)
	if tv.head == nil {
		return true, groupIndicator
	}

	toVisit := nodeStack{}
	toVisit.Push(tv.head)

	for !toVisit.Empty() {
		visitedNode, _ := toVisit.Pop()

		if visitedNode.level >= len(nsElems) {
			continue
		}

		if nsElems[visitedNode.level] != staticAnyMatcher {
			switch compatibilityMode {
			case false:
				if !visitedNode.currentElement.Match(nsElems[visitedNode.level]) {
					continue
				}
			case true:
				if !visitedNode.currentElement.Compatible(nsElems[visitedNode.level]) {
					continue
				}
			}
		}

		if visitedNode.level == len(nsElems)-1 {
			if fullMatch && visitedNode.nodeType == leafLevel {
				return true, visitedNode.groupIndicator()
			}
			if !fullMatch {
				return true, groupIndicator
			}
		}

		if _, ok := visitedNode.currentElement.(*staticRecursiveAnyElement); ok { // if ** we don't need to match anymore
			return true, groupIndicator
		}

		for _, subNode := range visitedNode.subNodes {
			toVisit.Push(subNode)
		}

	}

	return false, groupIndicator
}

func (tv *TreeValidator) ListRules() []string {
	if tv.head == nil {
		return []string{}
	}

	var nsList []string

	toVisit := nodeStack{}
	toVisit.Push(tv.head)

	for !toVisit.Empty() {
		visitedNode, _ := toVisit.Pop()

		if visitedNode.nodeType == leafLevel {
			nsList = append(nsList, visitedNode.path())
			continue
		}

		for _, subNode := range visitedNode.subNodes {
			toVisit.Push(subNode)
		}
	}

	sort.Strings(nsList)
	return nsList
}

// this function looks where to put new namespace elements and if tree conditions are met, updates the tree
func (tv *TreeValidator) updateTree(parsedNs *Namespace) error {
	// special case - tree doesn't contain anything
	if tv.head == nil {
		tv.head = tv.createNodes(parsedNs, 0)
		return nil
	}

	nodeToUpdate, namespacesToAttach, err := tv.findNodeToUpdate(tv.head, parsedNs)
	if err != nil {
		return err
	}

	nodesToAttach := tv.createNodes(namespacesToAttach, nodeToUpdate.level+1)
	return nodeToUpdate.attachNode(nodesToAttach)
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
		if node, ok := head.subNodes[nextElem.String()]; ok {
			return tv.findNodeToUpdate(node, &Namespace{elements: parsedNs.elements[1:]})
		} else {
			return head, &Namespace{parsedNs.elements[1:]}, nil
		}
	}

	return nil, nil, errors.New("not implemented")
}

// will create the entire branch of nodes from namespace (not update the tree, only returns branch)
// ie. /plugin/group1/metric will create branch consisting of 3 elements (node plugin -> node group1 -> leaf metric)
func (tv *TreeValidator) createNodes(ns *Namespace, level int) *Node {
	if len(ns.elements) == 0 {
		return nil
	}
	if len(ns.elements) == 1 {
		return &Node{
			currentElement: ns.elements[0],
			subNodes:       nil,
			nodeType:       leafLevel,
			level:          level,
		}
	}

	currNode := &Node{
		currentElement: ns.elements[0],
		subNodes:       map[string]*Node{},
		level:          level,
	}
	nextNode := tv.createNodes(&Namespace{elements: ns.elements[1:]}, level+1)
	nextNode.parent = currNode

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

	currNode.subNodes = map[string]*Node{nextNode.currentElement.String(): nextNode}

	return currNode
}

func (n *Node) attachNode(attachedNode *Node) error {
	isNextNodeStatic := !attachedNode.currentElement.IsDynamic()

	if n.nodeType == onlyStaticElementsLevel && !isNextNodeStatic {
		return errors.New("only static elements may added at current level")
	}

	if n.nodeType == onlyDynamicElementsLevel && !isNextNodeStatic {
		return errors.New("there can be only one dynamic element at current level")
	}

	n.subNodes[attachedNode.currentElement.String()] = attachedNode
	attachedNode.parent = n

	return nil
}

func (n *Node) trace() []*Node {
	nodeTrace := make([]*Node, n.level+1)
	nodeTrace[n.level] = n

	currNode := n
	for currNode.parent != nil {
		currNode = currNode.parent
		nodeTrace[currNode.level] = currNode
	}

	return nodeTrace
}

func (n *Node) path() string {
	trace := n.trace()

	nsElems := make([]string, 0, len(trace))
	for _, node := range trace {
		nsElems = append(nsElems, node.currentElement.String())
	}

	return DefaultNsSeparator + strings.Join(nsElems, DefaultNsSeparator)
}

func (n *Node) groupIndicator() []string {
	trace := n.trace()

	groupIndicator := make([]string, len(trace))

	for i, node := range trace {
		if groupNode, ok := node.currentElement.(*dynamicAnyElement); ok {
			groupIndicator[i] = groupNode.group
		}
	}

	return groupIndicator
}
