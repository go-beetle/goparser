package parser

import (
	"goparser/html"
	"sync"
)

type TextMap struct {
	Data map[uint32]string
	Mu   sync.Mutex
}

type AttrMap struct {
	Data map[uint32]html.StringMap
	Mu   sync.Mutex
}

func (d *TextMap) getNextNodeID() uint32 {
	return uint32(len(d.Data) + 1)
}

func (d *AttrMap) getNextNodeID() uint32 {
	return uint32(len(d.Data) + 1)
}

func matchVal(s1, s2 string) bool {
	return s1 == s2
}
