package parser

import (
	"reflect"
)

type ParserItem struct {
	Sequence  float64
	From      ParserBaseEvent
	EventType ParserBaseEvent
	ItemType  reflect.Type
	Tick      int
	Data      []byte
	Value     interface{}
	Bits      string
}

// ParserItems attaches the methods of Interface to []*ParserItem, sorting in increasing order by Sequence.
type ParserItems []*ParserItem

func (p ParserItems) Len() int           { return len(p) }
func (p ParserItems) Less(i, j int) bool { return p[i].Sequence < p[j].Sequence }
func (p ParserItems) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
