package parser

import (
	"fmt"
	"reflect"

	// "code.google.com/p/gogoprotobuf/proto"
	"github.com/elobuff/d2rp/core/utils"
	dota "github.com/elobuff/d2rp/dota"
)

const (
	headerLength = 12
	headerMagic  = "PBUFDEM"
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

type ParserBase struct {
	reader *utils.BytesReader
}

func NewParserBase(datas []byte) *ParserBase {
	if len(datas) < headerLength {
		panic("File too small.")
	}

	magic := ReadStringZ(datas, 0)
	if magic != headerMagic {
		panic("demofilestamp doesn't match, was: " + magic)
	}

	totalLength := len(datas) - headerLength
	if totalLength < 1 {
		panic("couldn't open file")
	}

	buffer := datas[headerLength:totalLength]
	return &ParserBase{
		reader: &utils.BytesReader{Data: buffer},
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
	return BumDum[event]
}

func (p *ParserBase) AsBaseType(event string) ParserBaseEventMap {
	return Maps[event]
}
