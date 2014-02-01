package parser

import (
	"sort"

	"code.google.com/p/gogoprotobuf/proto"
	"github.com/davecgh/go-spew/spew"
	"github.com/dotabuff/d2rp/core/utils"
	dota "github.com/dotabuff/d2rp/dota"
)

func foo() { spew.Dump("hi") }

type Parser struct {
	*ParserBase
	Sequence int64
	Items    map[int64]*ParserItem
}

func ParserFromFile(path string) *Parser {
	return NewParser(ReadFile(path))
}

func NewParser(data []byte) *Parser {
	parser := &Parser{ParserBase: NewParserBase(data)}
	parser.Analyze()
	return parser
}

func (p *Parser) Analyze() {
	p.Sequence = 1
	p.Items = map[int64]*ParserItem{}
	command, compressed := p.ReadEDemoCommands()
	if command == dota.EDemoCommands_DEM_Error {
		panic(command)
	}
	for p.reader.CanRead() {
		tick := int(p.reader.ReadVarInt32())
		length := int(p.reader.ReadVarInt32())
		obj, err := p.AsBaseEvent(command.String())
		if err == nil {
			item := &ParserItem{
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
				p.AnalyzePacket(dota.EDemoCommands_DEM_SignonPacket, tick, full.GetData())
			case *dota.CDemoPacket:
				ProtoUnmarshal(item.Data, o)
				p.AnalyzePacket(dota.EDemoCommands_DEM_Packet, tick, o.GetData())
			case *dota.CDemoFullPacket:
				ProtoUnmarshal(item.Data, o)
				item.From = dota.EDemoCommands_DEM_FullPacket
				item.Data = nil
				item.Object = o.GetStringTable()
				p.Items[item.Sequence] = item
				p.AnalyzePacket(dota.EDemoCommands_DEM_FullPacket, tick, o.GetPacket().GetData())
			case *dota.CDemoSendTables:
				ProtoUnmarshal(item.Data, o)
				p.AnalyzePacket(dota.EDemoCommands_DEM_SendTables, tick, o.GetData())
			default:
				p.Items[item.Sequence] = item
			}
		} else {
			p.reader.Skip(length)
		}
		command, compressed = p.ReadEDemoCommands()
		if command == dota.EDemoCommands_DEM_Error {
			panic(command)
		}
	}
}
func (p *Parser) AnalyzePacket(fromEvent dota.EDemoCommands, tick int, data []byte) {
	reader := utils.NewBytesReader(data)
	for reader.CanRead() {
		iType := int(reader.ReadVarInt32())
		length := int(reader.ReadVarInt32())
		obj, err := p.AsBaseEventNETSVC(iType)
		if err != nil {
			spew.Println(err)
			reader.Skip(length)
		} else {
			item := &ParserItem{
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
					p.Items[item.Sequence] = item
				}
			default:
				p.Items[item.Sequence] = item
			}
		}
	}
}

func (p *Parser) Parse(check func(proto.Message) bool) (items ParserBaseItems) {
	for _, item := range p.Items {
		if item.Data == nil {
			continue
		}
		if !check(item.Object) {
			continue
		}
		items = append(items, parseOne(item))
	}
	sort.Sort(items)
	return items
}

func parseOne(item *ParserItem) *ParserBaseItem {
	err := ProtoUnmarshal(item.Data, item.Object)
	if err != nil {
		spew.Println("parseOne()")
		spew.Dump(item)
		panic(err)
		return &ParserBaseItem{}
	}
	item.Data = nil
	return &ParserBaseItem{
		Sequence: item.Sequence,
		Tick:     item.Tick,
		From:     item.From,
		Object:   item.Object,
	}
}
