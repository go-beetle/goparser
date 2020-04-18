package main

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"goparser/html"
	"goparser/parser"
	"strings"
	"sync"
)

func HTMLParse(s string) (parser.Tree, bool) {
	r, err := html.Parse(strings.NewReader(s))
	if err != nil {
		return parser.Tree{}, false
	}
	for r.Type != html.ElementNode {
		switch r.Type {
		case html.DocumentNode:
			r = r.FirstChild
		case html.DoctypeNode:
			r = r.NextSibling
		case html.CommentNode:
			r = r.NextSibling
		}
	}
	root_node := parser.Node{Pointer: r, NodeType: 0}
	root_node.ResetNodeType()
	root_node.MakeAttributeMap()
	return parser.Tree{root_node}, true
}

func main() {
	host1 := "https://xkcd.com/2/"
	host2 := "https://xkcd.com/3/"
	resp1, _ := soup.Get(host1)
	doc1, _ := HTMLParse(resp1)
	resp2, _ := soup.Get(host2)
	doc2, _ := HTMLParse(resp2)
	//doc1.root.printTree(0)
	temp := parser.Template{doc1,
		make(map[string]parser.TextMap, 0),
		make(map[string]parser.AttrMap, 0)}
	temp.AddMaps("test",
		parser.TextMap{make(map[uint32]string), sync.Mutex{}},
		parser.AttrMap{make(map[uint32]html.StringMap), sync.Mutex{}})
	data, ok := temp.TextDiff(&doc2)
	fmt.Println(data, ok)
	data2, ok := temp.AttrDiff(&doc2)
	fmt.Println(data2, ok)
	fmt.Println(temp)

}
