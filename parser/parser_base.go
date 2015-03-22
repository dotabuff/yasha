package parser

import (
	"io/ioutil"
	"os"

	"compress/bzip2"

	"github.com/dotabuff/yasha/dota"
	"github.com/golang/protobuf/proto"
	"github.com/siddontang/go/snappy"
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

func ReadBz2File(path string) []byte {
	fd, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	reader := bzip2.NewReader(fd)
	raw, err := ioutil.ReadAll(reader)
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

type Error string

func (e Error) Error() string { return string(e) }

type SignonPacket struct{}

func (s SignonPacket) ProtoMessage()  {}
func (s SignonPacket) Reset()         {}
func (s SignonPacket) String() string { return "" }
