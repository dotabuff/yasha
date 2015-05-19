package yasha

import (
	"bytes"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/dotabuff/yasha/dota"
)

type OuterParser struct {
	reader   *BytesReader
	Sequence int64
	Items    map[int64]*OuterParserItem
}

func OuterParserFromFile(path string) *OuterParser {
	if strings.HasSuffix(path, ".dem.bz2") {
		return NewOuterParser(ReadBz2File(path))
	} else if strings.HasSuffix(path, ".dem") {
		return NewOuterParser(ReadFile(path))
	} else {
		panic("expected path to .dem or .dem.bz2 instead of " + path)
	}
}

func NewOuterParser(data []byte) *OuterParser {
	if len(data) < headerLength {
		panic("File too small.")
	}

	magic := ReadStringZ(data, 0)
	if magic != headerMagic {
		panic("demofilestamp doesn't match, was: " + spew.Sdump(magic))
	}

	totalLength := len(data) - headerLength
	if totalLength < 1 {
		panic("couldn't open file")
	}

	buffer := data[headerLength:(headerLength + totalLength)]
	return &OuterParser{reader: NewBytesReader(buffer)}
}

func (p *OuterParser) readEDemoCommands() (dota.EDemoCommands, bool) {
	command := dota.EDemoCommands(p.reader.ReadVarInt32())
	compressed := (command & dota.EDemoCommands_DEM_IsCompressed) == dota.EDemoCommands_DEM_IsCompressed
	command = command & ^dota.EDemoCommands_DEM_IsCompressed
	return command, compressed
}

func (p *OuterParser) Analyze(callback func(*OuterParserBaseItem)) {
	p.Sequence = 1
	command, compressed := p.readEDemoCommands()
	if command == dota.EDemoCommands_DEM_Error {
		panic(command)
	}
	for p.reader.CanRead() {
		tick := int(p.reader.ReadVarInt32())
		length := int(p.reader.ReadVarInt32())
		obj, err := p.AsBaseEvent(command.String())
		if err == nil {
			item := &OuterParserItem{
				Sequence: p.Sequence,
				Object:   obj,
				Tick:     tick,
			}
			p.Sequence++
			if compressed {
				item.Data = SnappyUncompress(p.reader.Read(length))
			} else {
				item.Data = p.reader.Read(length)
			}
			switch o := obj.(type) {
			case *SignonPacket:
				full := &dota.CDemoPacket{}
				ProtoUnmarshal(item.Data, full)
				p.AnalyzePacket(callback, dota.EDemoCommands_DEM_SignonPacket, tick, full.GetData())
			case *dota.CDemoPacket:
				ProtoUnmarshal(item.Data, o)
				p.AnalyzePacket(callback, dota.EDemoCommands_DEM_Packet, tick, o.GetData())
			case *dota.CDemoFullPacket:
				ProtoUnmarshal(item.Data, o)
				item.From = dota.EDemoCommands_DEM_FullPacket
				item.Data = nil
				item.Object = o.GetStringTable()
				callback(parseOne(item))
				p.AnalyzePacket(callback, dota.EDemoCommands_DEM_FullPacket, tick, o.GetPacket().GetData())
			case *dota.CDemoSendTables:
				ProtoUnmarshal(item.Data, o)
				p.AnalyzePacket(callback, dota.EDemoCommands_DEM_SendTables, tick, o.GetData())
			default:
				callback(parseOne(item))
			}
		} else {
			p.reader.Skip(length)
		}
		command, compressed = p.readEDemoCommands()
		if command == dota.EDemoCommands_DEM_Error {
			panic(command)
		}
	}
}

func (p *OuterParser) AnalyzePacket(callback func(*OuterParserBaseItem), fromEvent dota.EDemoCommands, tick int, data []byte) {
	reader := NewBytesReader(data)
	for reader.CanRead() {
		iType := int(reader.ReadVarInt32())
		length := int(reader.ReadVarInt32())
		obj, err := p.AsBaseEventNETSVC(iType)
		if err != nil {
			spew.Println(err)
			reader.Skip(length)
		} else {
			item := &OuterParserItem{
				Sequence: p.Sequence,
				From:     fromEvent,
				Object:   obj,
				Tick:     tick,
				Data:     reader.Read(length),
			}
			p.Sequence++
			switch obj.(type) {
			case *dota.CSVCMsg_UserMessage:
				message := &dota.CSVCMsg_UserMessage{}
				ProtoUnmarshal(item.Data, message)
				um, err := p.AsBaseEventBUMDUM(int(message.GetMsgType()))
				if err == nil {
					item.Object = um
					item.Data = message.GetMsgData()
					callback(parseOne(item))
				}
			default:
				callback(parseOne(item))
			}
		}
	}
}

func parseOne(item *OuterParserItem) *OuterParserBaseItem {
	err := ProtoUnmarshal(item.Data, item.Object)
	if err != nil {
		spew.Println("parseOne()")
		spew.Dump(item)
		panic(err)
		return &OuterParserBaseItem{}
	}
	item.Data = nil
	return &OuterParserBaseItem{
		Sequence: item.Sequence,
		Tick:     item.Tick,
		From:     item.From,
		Object:   item.Object,
	}
}

func ReadStringZ(datas []byte, offset int) string {
	idx := bytes.IndexByte(datas[offset:], '\000')
	if idx < 0 {
		return ""
	}
	return string(datas[offset:idx])
}