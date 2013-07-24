package parser

import (
	"reflect"
	"sort"

	"code.google.com/p/gogoprotobuf/proto"
	"github.com/elobuff/d2rp/core/utils"
	dota "github.com/elobuff/d2rp/dota"
)

type Parser struct {
	*ParserBase
	Sequence float64
	Items    map[float64]*ParserItem
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
	p.Items = map[float64]*ParserItem{}
	compressed := false
	command := p.ReadEDemoCommands(&compressed)
	if command == dota.EDemoCommands_DEM_Error {
		panic(command)
	}
	for p.reader.CanRead() {
		tick := int(p.reader.ReadVarInt32())
		length := int(p.reader.ReadVarInt32())
		pbEvent := p.AsParserBaseEvent(command)
		if pbEvent.EventType == BaseError {
			p.reader.Skip(length)
		} else {
			item := &ParserItem{
				Sequence:  p.Sequence,
				From:      pbEvent.EventType,
				EventType: pbEvent.EventType,
				ItemType:  p.AsBaseType(pbEvent.Name).ItemType,
				Tick:      tick,
			}
			p.Sequence++
			if compressed {
				item.Data = SnappyUncompress(p.reader.Read(length))
			} else {
				item.Data = p.reader.Read(length)
			}
			if pbEvent.EventType == DEM_SignonPacket {
				full := &dota.CDemoPacket{}
				ProtoUnmarshal(item.Data, full)
				p.AnalyzePacket(DEM_SignonPacket, tick, full.Data)
			} else if pbEvent.EventType == DEM_Packet {
				full := &dota.CDemoPacket{}
				ProtoUnmarshal(item.Data, full)
				p.AnalyzePacket(DEM_Packet, tick, full.Data)
			} else if pbEvent.EventType == DEM_FullPacket {
				full := &dota.CDemoFullPacket{}
				ProtoUnmarshal(item.Data, full)
				item.From = DEM_FullPacket
				item.EventType = DEM_StringTables
				item.ItemType = p.AsBaseType("DEM_StringTables").ItemType
				item.Data = nil
				item.Value = full.StringTable
				p.Items[item.Sequence] = item
				p.AnalyzePacket(DEM_FullPacket, tick, full.GetPacket().Data)
			} else if pbEvent.EventType == DEM_SendTables {
				full := &dota.CDemoSendTables{}
				ProtoUnmarshal(item.Data, full)
				p.AnalyzePacket(DEM_SendTables, tick, full.Data)
			} else {
				p.Items[item.Sequence] = item
			}
		}
		command = p.ReadEDemoCommands(&compressed)
		if command == dota.EDemoCommands_DEM_Error {
			panic(command)
		}
	}
}
func (p *Parser) AnalyzePacket(fromEvent ParserBaseEvent, tick int, data []byte) {
	reader := utils.NewBytesReader(data)
	for reader.CanRead() {
		iType := int(reader.ReadVarInt32())
		length := int(reader.ReadVarInt32())
		pbEvent := p.AsParserBaseEventNETSVC(iType)
		if pbEvent.EventType == BaseError {
			reader.Skip(length)
		} else {
			item := &ParserItem{
				Sequence:  p.Sequence,
				From:      fromEvent,
				EventType: pbEvent.EventType,
				ItemType:  p.AsBaseType(pbEvent.Name).ItemType,
				Tick:      tick,
				Data:      reader.Read(length),
			}
			p.Sequence++
			if pbEvent.EventType == SVC_UserMessage {
				message := &dota.CSVCMsg_UserMessage{}
				ProtoUnmarshal(item.Data, message)
				umEvent := p.AsParserBaseEventBUMDUM(int(message.GetMsgType()))
				if umEvent.EventType != BaseError {
					item.EventType = umEvent.EventType
					item.ItemType = p.AsBaseType(umEvent.Name).ItemType
					item.Data = message.GetMsgData()
					p.Items[item.Sequence] = item
				}
			} else {
				p.Items[item.Sequence] = item
			}
		}
	}
}

func (p *Parser) Parse(events ...ParserBaseEvent) (items ParserBaseItems) {
	emap := map[ParserBaseEvent]bool{}
	for _, event := range events {
		emap[event] = true
	}
	for _, item := range p.Items {
		if _, found := emap[item.EventType]; !found {
			continue
		}
		if item.Data != nil {
			v := reflect.New(item.ItemType)
			ProtoUnmarshal(item.Data, v.Addr().Interface().(proto.Message))
			item.Value = v.Interface()
			item.Data = nil
		}
		items = append(items, &ParserBaseItem{
			Sequence: item.Sequence,
			From:     item.From,
			ItemType: item.ItemType,
			Tick:     item.Tick,
			Value:    item.Value,
		})
	}
	sort.Sort(items)
	return items
}

/*
public IEnumerable<ParserBaseItem> Parse(params ParserBaseEvent[] events) {
  Parallel.ForEach(this.Items,
    o => {
      if (o.Datas != null && events.Contains(o.EventType)) {
        o.Value = o.Datas.Protobuf(o.ItemType);
        o.Datas = null;
      }
      });
  return this.Items.Where(o => events.Contains(o.EventType))
    .OrderBy(o => o.Sequence)
    .Select(o => new ParserBaseItem {
        Sequence = o.Sequence,
        From = o.From,
        ItemType = o.ItemType,
        Tick = o.Tick,
        Value = o.Value
        });
}
public IEnumerable<ParserBaseItem> Parse() {
  Parallel.ForEach(
    this.Items, o => {
      if (o.Datas != null) {
      o.Value = o.Datas.Protobuf(o.ItemType);
      o.Datas = null;
    }
  });
  return this.Items.OrderBy(
    o => o.Sequence
  ).Select(
    o => new ParserBaseItem {
      Sequence = o.Sequence,
      From = o.From,
      ItemType = o.ItemType,
      Tick = o.Tick,
      Value = o.Value
    }
  );
}
*/
