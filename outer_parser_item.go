package yasha

import (
	"github.com/dotabuff/yasha/dota"
	"github.com/golang/protobuf/proto"
)

type OuterParserItem struct {
	Sequence int64
	Tick     int
	Data     []byte
	From     dota.EDemoCommands
	Object   proto.Message
	// Bits     string
}

// OuterParserItems attaches the methods of Interface to []*OuterParserItem, sorting in increasing order by Sequence.
type OuterParserItems []*OuterParserItem

func (p OuterParserItems) Len() int           { return len(p) }
func (p OuterParserItems) Less(i, j int) bool { return p[i].Sequence < p[j].Sequence }
func (p OuterParserItems) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
