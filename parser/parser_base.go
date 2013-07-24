package parser

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"code.google.com/p/gogoprotobuf/proto"
	"code.google.com/p/snappy-go/snappy"
	"github.com/elobuff/d2rp/core/utils"
	dota "github.com/elobuff/d2rp/dota"
)

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
type ParserBaseEventMap struct {
	Name      string
	EventType ParserBaseEvent
	ItemType  reflect.Type
	MapType   ParserBaseEventMapType
	Value     int
}

func SnappyUncompress(compressed []byte) []byte {
	dst := make([]byte, 0, len(compressed))
	out, err := snappy.Decode(dst, compressed)
	if err != nil {
		panic(err)
	}
	return out
}

func ProtoUnmarshal(data []byte, obj proto.Message) {
	err := proto.Unmarshal(data, obj)
	if err != nil {
		panic(err)
	}
}

func ReadFile(path string) []byte {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return raw
}

type ParserBase struct {
	reader *utils.BytesReader
}

func NewParserBase(data []byte) *ParserBase {
	if len(data) < headerLength {
		panic("File too small.")
	}

	magic := ReadStringZ(data, 0)
	if magic != headerMagic {
		panic("demofilestamp doesn't match, was: " + magic)
	}

	totalLength := len(data) - headerLength
	if totalLength < 1 {
		panic("couldn't open file")
	}

	buffer := data[headerLength:(headerLength + totalLength)]
	return &ParserBase{
		reader: utils.NewBytesReader(buffer),
	}
}

func (p *ParserBase) ReadEDemoCommands(compressed *bool) dota.EDemoCommands {
	command := dota.EDemoCommands(p.reader.ReadVarInt32())
	*compressed = (command & dota.EDemoCommands_DEM_IsCompressed) == dota.EDemoCommands_DEM_IsCompressed
	command = command & ^dota.EDemoCommands_DEM_IsCompressed
	return command
}

func (p *ParserBase) AsParserBaseEvent(event fmt.Stringer) ParserBaseEventMap {
	return Maps[event.String()]
}
func (p *ParserBase) AsParserBaseEventNETSVC(event int) ParserBaseEventMap {
	return NetSvc[event]
}

func (p *ParserBase) AsParserBaseEventBUMDUM(event int) ParserBaseEventMap {
	if event > len(BumDum) {
		return ParserBaseEventMap{EventType: BaseError}
	}
	return BumDum[event]
}

func (p *ParserBase) AsBaseType(event string) ParserBaseEventMap {
	return Maps[event]
}
