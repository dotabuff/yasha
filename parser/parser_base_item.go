package parser

import (
	"code.google.com/p/gogoprotobuf/proto"
	"reflect"
)

type ParserBaseItem struct {
	Sequence float64
	From     ParserBaseEvent
	ItemType reflect.Type
	Tick     int
	Value    proto.Message
}

// ParserBaseItems attaches the methods of Interface to []*ParserBaseItem, sorting in increasing order by Sequence.
type ParserBaseItems []*ParserBaseItem

func (p ParserBaseItems) Len() int           { return len(p) }
func (p ParserBaseItems) Less(i, j int) bool { return p[i].Sequence < p[j].Sequence }
func (p ParserBaseItems) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
