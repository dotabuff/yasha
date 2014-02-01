package parser

import (
	"io/ioutil"

	"code.google.com/p/gogoprotobuf/proto"
	"code.google.com/p/snappy-go/snappy"
	"github.com/dotabuff/d2rp/dota"
)

func SnappyUncompress(compressed []byte) []byte {
	dst := make([]byte, 0, len(compressed))
	out, err := snappy.Decode(dst, compressed)
	if err != nil {
		panic(err)
	}
	return out
}

func ProtoUnmarshal(data []byte, obj proto.Message) error {
	return proto.Unmarshal(data, obj)
}

func ReadFile(path string) []byte {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return raw
}

const (
	headerLength = 12
	headerMagic  = "PBUFDEM"
)

const (
	DEM ParserBaseEventMapType = iota
	NET
	SVC
	BUM
	DUM
)

type ParserBaseEvent int
type ParserBaseEventMapType int
type ItemType int

type ParserBaseItem struct {
	Sequence int64
	Tick     int
	From     dota.EDemoCommands
	Object   proto.Message
}

// ParserBaseItems attaches the methods of Interface to []*ParserBaseItem, sorting in increasing order by Sequence.
type ParserBaseItems []*ParserBaseItem

func (p ParserBaseItems) Len() int           { return len(p) }
func (p ParserBaseItems) Less(i, j int) bool { return p[i].Sequence < p[j].Sequence }
func (p ParserBaseItems) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
