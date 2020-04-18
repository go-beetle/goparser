package parser

import (
	"fmt"
	"goparser/html"
)

const (
	TEXT_NODE     = 0
	ATTR_NODE     = 1
	TEXT_VAR_NODE = 2
	ATTR_VAR_NODE = 3
)

type Node struct {
	Pointer  *html.Node
	NodeType uint32
	NodeID   uint32
}

type NodeChildren []Node

func (nc NodeChildren) isEmpty() bool {
	return len(nc) == 0
}

func equalLength(nc1, nc2 NodeChildren) bool {
	return len(nc1) == len(nc2)
}

func isDiffText(n1, n2 Node) bool {
	return n1.Pointer.Data != n2.Pointer.Data
}

func isDiffAttr(n1, n2 Node) bool {
	n1_attrs := n1.Pointer.Attr1
	n2_attrs := n2.Pointer.Attr1
	if len(n1_attrs) != len(n2_attrs) {
		return true
	}
	for key, val := range n1_attrs {
		if _, ok := n2_attrs[key]; !ok {
			return true
		}
		if n2_attrs[key] != val {
			return true
		}
	}
	return false
}

func nodeTypeMatch(n1, n2 Node) bool {
	return (n1.NodeType == n2.NodeType) || (n1.NodeType == TEXT_VAR_NODE && n2.NodeType == TEXT_NODE) || (n1.NodeType == ATTR_VAR_NODE && n2.NodeType == ATTR_NODE)
}

func (n *Node) isTextNode() bool {
	return n.NodeType == TEXT_NODE
}

func (n *Node) isAttrNode() bool {
	return n.NodeType == ATTR_NODE
}

func (n *Node) isVarNode() bool {
	return n.NodeType == ATTR_VAR_NODE || n.NodeType == TEXT_VAR_NODE
}

func (n *Node) setTextVar(nodeID uint32) {
	n.NodeType = TEXT_VAR_NODE
	n.Pointer.Type = TEXT_VAR_NODE
	n.NodeID = nodeID
}

func (n *Node) setAttrVar(nodeID uint32) {
	n.NodeType = ATTR_VAR_NODE
	n.Pointer.Type = ATTR_VAR_NODE
	n.NodeID = nodeID
}

func (n *Node) getAttrDiffs(n1 Node) (html.StringMap, html.StringMap) {
	attr_map, temp_map := make(html.StringMap), make(html.StringMap)
	n_attrs, n1_attrs := n.Pointer.Attr1, n1.Pointer.Attr1
	for key, val := range n_attrs {
		if n1_attrs.ContainsKey(key) {
			if !matchVal(val, n1_attrs[key]) {
				attr_map.AddPair(key, n1_attrs[key])
				temp_map.AddPair(key, val)
			}
		} else {
			attr_map.AddPair(key, "")
			temp_map.AddPair(key, val)
		}
	}
	for key, val := range n1_attrs {
		if !n_attrs.ContainsKey(key) {
			attr_map.AddPair(key, val)
			temp_map.AddPair(key, "")
		}
	}
	return attr_map, temp_map
}

func (n *Node) getChildren() NodeChildren {
	child := n.Pointer.FirstChild
	var children []Node
	for child != nil {
		children = append(children, Node{Pointer: child, NodeType: uint32(child.Type)})
		child = child.NextSibling
	}
	return children
}

func (n *Node) MakeAttributeMap() {
	n.Pointer.Attr1 = make(html.StringMap)
	for _, attr := range n.Pointer.Attr {
		n.Pointer.Attr1[attr.Key] = attr.Val
	}
	children := n.getChildren()
	for _, child := range children {
		child.MakeAttributeMap()
	}
}

func (n *Node) ResetNodeType() {
	if n.Pointer.Type == html.TextNode {
		n.Pointer.Type = TEXT_NODE
	} else {
		n.Pointer.Type = ATTR_NODE
	}
	children := n.getChildren()
	for _, child := range children {
		child.ResetNodeType()
	}
}

func (n *Node) printTree(k int) {
	for i := 0; i < k; i++ {
		fmt.Print(" ")
	}
	fmt.Println(n.Pointer)
	children := n.getChildren()
	for _, child := range children {
		child.printTree(k + 1)
	}
}
