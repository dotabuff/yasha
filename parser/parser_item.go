package parser

import (
	"code.google.com/p/gogoprotobuf/proto"
	dota "github.com/dotabuff/d2rp/dota"
)

type ParserItem struct {
	Sequence int64
	Tick     int
	Data     []byte
	From     dota.EDemoCommands
	Object   proto.Message
	// Bits     string
}

// ParserItems attaches the methods of Interface to []*ParserItem, sorting in increasing order by Sequence.
type ParserItems []*ParserItem

func (p ParserItems) Len() int           { return len(p) }
func (p ParserItems) Less(i, j int) bool { return p[i].Sequence < p[j].Sequence }
func (p ParserItems) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
