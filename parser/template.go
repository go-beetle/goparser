package parser

import (
	"goparser/html"
	"sync"
)

type Template struct {
	Tree     Tree
	TextMaps map[string]TextMap
	AttrMaps map[string]AttrMap
}

func (t *Template) AddMaps(host string, text TextMap, attr AttrMap) bool {
	_, ok := t.AttrMaps[host]
	if ok {
		return false
	}
	t.AttrMaps[host] = attr
	_, ok = t.TextMaps[host]
	if ok {
		return false
	}
	t.TextMaps[host] = text
	return true
}

func (t *Template) updateTemplateAttr(dataID uint32, attr_map map[string]string) {
	for host, _ := range t.AttrMaps {
		t.AttrMaps[host].Data[dataID] = attr_map
	}
}

func (t *Template) updateTemplateText(dataID uint32, text string) {
	for host, _ := range t.TextMaps {
		t.TextMaps[host].Data[dataID] = text
	}
}

func (t *Template) AttrDiff(new_tree *Tree) (AttrMap, bool) {
	data := AttrMap{make(map[uint32]html.StringMap), sync.Mutex{}}
	ok := t.attrDiffHelper(&t.Tree.Root, &new_tree.Root, &data)
	if !ok {
		return AttrMap{}, false
	}
	return data, true
}

func (t *Template) attrDiffHelper(t1, t2 *Node, data *AttrMap) bool {
	t1_children := t1.getChildren()
	t2_children := t2.getChildren()
	if !equalLength(t1_children, t2_children) {
		return false
	}
	if t1_children.isEmpty() {
		return true
	}
	for i, child := range t1_children {
		if !nodeTypeMatch(child, t2_children[i]) {
			return false
		}
		if child.isVarNode() {
			diffAttr, _ := child.getAttrDiffs(t2_children[i])
			data.Data[child.NodeID] = diffAttr
		} else if child.isAttrNode() {
			if isDiffAttr(child, t2_children[i]) {
				diffAttr, tempAttr := child.getAttrDiffs(t2_children[i])
				nodeID := data.getNextNodeID()
				child.setAttrVar(nodeID)
				data.Data[nodeID] = diffAttr
				t.updateTemplateAttr(nodeID, tempAttr)
			}
		}
		ok := t.attrDiffHelper(&child, &t2_children[i], data)
		if !ok {
			return false
		}
	}
	return true
}

func (t *Template) TextDiff(new_tree *Tree) (TextMap, bool) {
	data := TextMap{make(map[uint32]string), sync.Mutex{}}
	ok := t.textDiffHelper(&t.Tree.Root, &new_tree.Root, &data)
	if !ok {
		return TextMap{}, false
	}
	return data, true
}

func (t *Template) textDiffHelper(t1, t2 *Node, data *TextMap) bool {
	t1_children := t1.getChildren()
	t2_children := t2.getChildren()
	if !equalLength(t1_children, t2_children) {
		return false
	}
	if t1_children.isEmpty() {
		return true
	}
	for i, child := range t1_children {
		if !nodeTypeMatch(child, t2_children[i]) {
			return false
		}
		if child.isVarNode() {
			data.Data[child.NodeID] = t2_children[i].Pointer.Data
		} else if child.isTextNode() {
			if isDiffText(child, t2_children[i]) {
				nodeID := data.getNextNodeID()
				child.setTextVar(nodeID)
				data.Data[nodeID] = t2_children[i].Pointer.Data
				t.updateTemplateText(nodeID, child.Pointer.Data)
			}
		}
		ok := t.textDiffHelper(&child, &t2_children[i], data)
		if !ok {
			return false
		}
	}
	return true
}
