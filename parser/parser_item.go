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
