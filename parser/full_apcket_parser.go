package parser

import (
	"io/ioutil"

	"github.com/elobuff/d2rp/core/utils"

	"code.google.com/p/gogoprotobuf/proto"
	//	"github.com/elobuff/d2rp/core/utils"
	"code.google.com/p/snappy-go/snappy"
	dota "github.com/elobuff/d2rp/dota"
)

func SnappyUncompress(compressed []byte) []byte {
	dst := make([]byte, 0, len(compressed))
	out, err := snappy.Decode(dst, compressed)
	if err != nil {
		panic(err)
	}
	return out
}

type FullPacketParser struct {
	*ParserBase
	Sequence float64
	Items    map[float64]*ParserItem
}

func LoadFromFile(path string) *FullPacketParser {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return NewFullPacketParser(raw)
}

func NewFullPacketParser(data []byte) *FullPacketParser {
	fppb := NewParserBase(data)
	fpp := &FullPacketParser{ParserBase: fppb}
	fpp.Analyze()
	return fpp
}

func (fpp *FullPacketParser) Analyze() {
	fpp.Sequence = 1
	fpp.Items = map[float64]*ParserItem{}
	compressed := false
	command := fpp.ReadEDemoCommands(&compressed)
	if command == dota.EDemoCommands_DEM_Error {
		panic(command)
	}
	for fpp.reader.CanRead() {
		tick := int(fpp.reader.ReadVarInt32())
		length := int(fpp.reader.ReadVarInt32())
		pbEvent := fpp.AsParserBaseEvent(command)
		if pbEvent.EventType == BaseError {
			fpp.reader.Skip(length)
		} else {
			item := &ParserItem{
				Sequence:  fpp.Sequence,
				From:      pbEvent.EventType,
				EventType: pbEvent.EventType,
				ItemType:  fpp.AsBaseType(pbEvent.Name).ItemType,
				Tick:      tick,
			}
			fpp.Sequence++
			if compressed {
				item.Data = SnappyUncompress(fpp.reader.Read(length))
			} else {
				item.Data = fpp.reader.Read(length)
			}
			if pbEvent.EventType == DEM_FullPacket {
				full := dota.CDemoFullPacket{}
				err := proto.Unmarshal(item.Data, &full)
				if err != nil {
					panic(err)
				}
				item.From = DEM_FullPacket
				item.EventType = DEM_StringTables
				item.Data = nil
				item.Value = full.GetStringTable()
				fpp.Items[item.Sequence] = item
				fpp.AnalyzePacket(DEM_FullPacket, tick, full.GetPacket().Data)
			} else if pbEvent.EventType == DEM_SendTables && item.Tick == 0 {
				full := dota.CDemoSendTables{}
				err := proto.Unmarshal(item.Data, &full)
				if err != nil {
					panic(err)
				}
				fpp.AnalyzePacket(DEM_SendTables, tick, full.GetData())
			} else {
				fpp.Items[item.Sequence] = item
			}
		}
		command = fpp.ReadEDemoCommands(&compressed)
		if command == dota.EDemoCommands_DEM_Error {
			panic(command)
		}
	}
}

func (fpp *FullPacketParser) AnalyzePacket(fromEvent ParserBaseEvent, tick int, data []byte) {
	reader := utils.BytesReader{Data: data}
	for reader.CanRead() {
		iType := int(reader.ReadVarInt32())
		pbEvent := fpp.AsParserBaseEventNETSVC(iType)
		length := int(reader.ReadVarInt32())
		if pbEvent.EventType != BaseError {
			item := &ParserItem{
				Sequence:  fpp.Sequence,
				From:      fromEvent,
				EventType: pbEvent.EventType,
				ItemType:  fpp.AsBaseType(pbEvent.Name).ItemType,
				Tick:      tick,
				Data:      reader.Read(length),
			}
			fpp.Sequence++
			fpp.Items[fpp.Sequence] = item
		} else {
			reader.Skip(int(length))
		}
	}
}
