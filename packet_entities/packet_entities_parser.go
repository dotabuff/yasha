package packet_entities

import (
	"math"
	"sort"
	"strconv"

	"github.com/elobuff/d2rp/core/parser"
	"github.com/elobuff/d2rp/core/send_tables"
	"github.com/elobuff/d2rp/core/utils"
	dota "github.com/elobuff/d2rp/dota"
)

type PacketEntitiesParser struct {
	baseline              map[int]map[string]interface{}
	classIdNumBits        int
	classInfosIdMapping   map[string]int
	classInfosNameMapping map[int]string
	createClassHandlers   map[string]Handler
	deleteClassHandlers   map[string]Handler
	mapping               map[int][]dota.CSVCMsg_SendTableSendpropT
	multiples             map[int]map[string]int
	packets               parser.ParserBaseItems
	preserveClassHandlers map[string]HandlerPreserve
	sendTablesHelper      *send_tables.Helper
	entities              []*PacketEntity
}

func (p PacketEntitiesParser) Entities() []*PacketEntity {
	return p.entities
}

func NewParser(items parser.ParserBaseItems) {
	p := PacketEntitiesParser{
		entities: make([]*PacketEntity, 0, 2048),
	}

	var serverInfo dota.CSVCMsg_ServerInfo
	packets := parser.ParserBaseItems{}
	for _, item := range items {
		switch value := item.Value.(type) {
		case dota.CSVCMsg_ServerInfo:
			serverInfo = value
		case dota.CSVCMsg_PacketEntities:
			if item.From == parser.DEM_Packet {
				packets = append(packets, item)
			}
		}
	}
	sort.Sort(packets)
	p.packets = packets

	p.classIdNumBits = int(math.Log(float64(serverInfo.GetMaxClasses()))/math.Log(2)) + 1

	var classInfos *dota.CDemoClassInfo
	var instanceBaseline *dota.CDemoStringTablesTableT
	sendTables := map[string]dota.CSVCMsg_SendTable{}

	for _, item := range items {
		switch value := item.Value.(type) {
		case dota.CDemoClassInfo:
			if classInfos == nil {
				classInfos = &value
			}
		case dota.CSVCMsg_SendTable:
			sendTables[value.GetNetTableName()] = value
		case dota.CDemoStringTables:
			for _, table := range value.GetTables() {
				if table.GetTableName() == "instancebaseline" {
					instanceBaseline = table
				}
			}
		}
	}

	for _, info := range classInfos.GetClasses() {
		id, name := int(info.GetClassId()), info.GetTableName()
		p.classInfosNameMapping[id] = name
		p.classInfosIdMapping[name] = id
		props := p.sendTablesHelper.LoadSendTable(name)
		p.mapping[id] = props

		m := map[string]int{}
		for _, prop := range props {
			m[prop.GetDtName()+"."+prop.GetVarName()] += 1
		}
		p.multiples[id] = m
	}

	for _, item := range instanceBaseline.GetItems() {
		classId, err := strconv.Atoi(item.GetStr())
		if err != nil {
			panic(err)
		}
		br := utils.NewBitReader(item.GetData())
		indices := br.ReadPropertiesIndex()
		p.baseline[classId] = br.ReadPropertiesValues(
			p.mapping[classId],
			p.multiples[classId],
			indices,
		)
	}
}

func (p *PacketEntitiesParser) AddCreateHandler(className string, callback Callback) {
	p.createClassHandlers[className] = Handler{ClassName: className, Callback: callback}
}

func (p *PacketEntitiesParser) AddDeleteHandler(className string, callback Callback) {
	p.deleteClassHandlers[className] = Handler{ClassName: className, Callback: callback}
}

func (p *PacketEntitiesParser) AddPreserveHandler(className string, callback PreserveCallback) {
	p.preserveClassHandlers[className] = HandlerPreserve{ClassName: className, Callback: callback}
}

func (p *PacketEntitiesParser) Parse() {
	for _, packet := range p.packets {
		p.ParsePacket(packet)
	}
}

func (p *PacketEntitiesParser) EntityCreate(br *utils.BitReader, currentIndex, tick int) {
	pe := &PacketEntity{
		Tick:      tick,
		ClassId:   int(br.ReadUBits(p.classIdNumBits)),
		SerialNum: int(br.ReadUBits(10)),
		Index:     currentIndex,
		Type:      Create,
		Values:    map[string]interface{}{},
	}
	pe.Name = p.classInfosNameMapping[pe.ClassId]
	indices := br.ReadPropertiesIndex()
	values := br.ReadPropertiesValues(
		p.mapping[pe.ClassId],
		p.multiples[pe.ClassId],
		indices,
	)
	if baseline, ok := p.baseline[pe.ClassId]; ok {
		for key, baseValue := range baseline {
			if subValue, ok := values[key]; ok {
				pe.Values[key] = subValue
			} else {
				pe.Values[key] = baseValue
			}
		}
	} else {
		for key, value := range values {
			pe.Values[key] = value
		}
	}

	p.entities[pe.Index] = pe
	if handler, ok := p.createClassHandlers[pe.Name]; ok {
		handler.Callback(pe)
	}
}

func (p *PacketEntitiesParser) EntityDelete(br *utils.BitReader, currentIndex, tick int) {
	pe := p.entities[currentIndex].Clone()
	pe.Tick = tick
	pe.Type = Delete
	p.entities[currentIndex] = nil
	if handler, ok := p.deleteClassHandlers[pe.Name]; ok {
		handler.Callback(pe)
	}
}

func (p *PacketEntitiesParser) EntityPreserve(br *utils.BitReader, currentIndex, tick int) {
	pe := p.entities[currentIndex]
	pe.Tick = tick
	pe.Type = Preserve
	indices := br.ReadPropertiesIndex()
	values := br.ReadPropertiesValues(
		p.mapping[p.classInfosIdMapping[pe.Name]],
		p.multiples[p.classInfosIdMapping[pe.Name]],
		indices,
	)
	for key, value := range values {
		pe.Values[key] = value
	}

	if handler, ok := p.preserveClassHandlers[pe.Name]; ok {
		handler.Callback(pe, values)
	}
}

func (p *PacketEntitiesParser) ParsePacket(packet *parser.ParserBaseItem) {
	pe := (packet.Value).(dota.CSVCMsg_PacketEntities)
	br := utils.NewBitReader(pe.GetEntityData())
	currentIndex := -1
	for i := 0; i < int(pe.GetUpdatedEntries()); i++ {
		currentIndex = br.ReadNextEntityIndex(currentIndex)
		switch ReadUpdateType(br) {
		case Preserve:
			p.EntityPreserve(br, currentIndex, packet.Tick)
		case Create:
			p.EntityCreate(br, currentIndex, packet.Tick)
		case Delete:
			p.EntityDelete(br, currentIndex, packet.Tick)
		}
	}
}

type Callback func(*PacketEntity)

type Handler struct {
	Callback  Callback
	ClassName string
}

type PreserveCallback func(*PacketEntity, map[string]interface{})

type HandlerPreserve struct {
	Callback  PreserveCallback
	ClassName string
}
