package core

import (
	"math"

	"github.com/davecgh/go-spew/spew"
	"github.com/elobuff/d2rp/core/packet_entities"
	"github.com/elobuff/d2rp/core/parser"
	"github.com/elobuff/d2rp/core/send_tables"
	"github.com/elobuff/d2rp/core/string_tables"
	"github.com/elobuff/d2rp/core/utils"
	dota "github.com/elobuff/d2rp/dota"
	"sort"
)

type Parser struct {
	Parser                *parser.Parser
	ClassIdNumBits        int
	ClassInfosIdMapping   map[string]int
	ClassInfosNameMapping map[int]string
	FileHeader            *dota.CDemoFileHeader
	FileInfo              *dota.CDemoFileInfo
	GameEventMap          map[int32]*dota.CSVCMsg_GameEventListDescriptorT
	Mapping               map[int][]*send_tables.SendProp
	Multiples             map[int]map[string]int
	ServerInfo            *dota.CSVCMsg_ServerInfo
	Sth                   *send_tables.Helper
	Stsh                  *string_tables.StateHelper
	VoiceInit             *dota.CSVCMsg_VoiceInit

	ActiveModifiers map[int]*dota.CDOTAModifierBuffTableEntry
	Entities        []*packet_entities.PacketEntity
	ByHandle        map[int]*packet_entities.PacketEntity

	OnEntityCreated        func(*packet_entities.PacketEntity)
	OnEntityDeleted        func(*packet_entities.PacketEntity)
	OnEntityPreserved      func(*packet_entities.PacketEntity)
	OnVoiceData            func(obj *dota.CSVCMsg_VoiceData)
	OnSpectatorPlayerClick func(tick int, obj *dota.CDOTAUserMsg_SpectatorPlayerClick)
	OnSetConVar            func(obj *dota.CNETMsg_SetConVar)
	OnSayText2             func(tick int, obj *dota.CUserMsg_SayText2)
	OnOverheadEvent        func(tick int, obj *dota.CDOTAUserMsg_OverheadEvent)
	OnChatEvent            func(tick int, obj *dota.CDOTAUserMsg_ChatEvent)
	OnCombatLog            func(log *CombatLogEntry)
	OnTablename            func(name string)
}

func ParserFromFile(path string) *Parser {
	return NewParser(parser.ReadFile(path))
}

func NewParser(data []byte) *Parser {
	return &Parser{Parser: parser.NewParser(data)}
}

func (p *Parser) Parse() {
	p.Sth = send_tables.NewHelper()
	p.Stsh = string_tables.NewStateHelper()
	p.Entities = make([]*packet_entities.PacketEntity, 2048)
	p.ClassInfosIdMapping = map[string]int{}
	p.ClassInfosNameMapping = map[int]string{}
	p.GameEventMap = map[int32]*dota.CSVCMsg_GameEventListDescriptorT{}
	p.Mapping = map[int][]*send_tables.SendProp{}
	p.Multiples = map[int]map[string]int{}
	p.ByHandle = map[int]*packet_entities.PacketEntity{}
	p.ActiveModifiers = map[int]*dota.CDOTAModifierBuffTableEntry{}

	// in order to successfully process data every tick, we need to maintain
	// order.  First of all the string and send tables for the tick have to be
	// done, then everything else.  But to also maintain streaming behaviour, we
	// simply stuff all data for one tick into a buffer, and iterate it twice.
	//
	// A bit of a waste of iterations because it still needs an expensive type
	// switch, but less overhead than doing a sort instead.

	currentTick := 0
	buffer := []*parser.ParserBaseItem{}

	p.Parser.Analyze(func(item *parser.ParserBaseItem) {
		// we got a new tick, process the previous items
		if item.Tick > currentTick {
			p.processTick(buffer)
			buffer = make([]*parser.ParserBaseItem, 1, len(buffer))
			buffer[0] = item
			currentTick = item.Tick
		} else {
			// tick is ongoing
			buffer = append(buffer, item)
		}
	})

	// and process all remaining in the last tick
	p.processTick(buffer)
}

func (p *Parser) processTick(items []*parser.ParserBaseItem) {
	p.Stsh.ActiveModifierDelta = string_tables.ModifierBuffs{}

	for _, item := range items {
		switch obj := item.Object.(type) {
		case *dota.CSVCMsg_SendTable:
			p.Sth.SetSendTable(obj.GetNetTableName(), obj)
		case *dota.CSVCMsg_ServerInfo:
			p.ServerInfo = obj
			p.ClassIdNumBits = int(math.Log(float64(obj.GetMaxClasses()))/math.Log(2)) + 1
		case *dota.CDemoClassInfo:
			p.onCDemoClassInfo(obj)
		}
	}

	for _, item := range items {
		switch obj := item.Object.(type) {
		case *dota.CSVCMsg_CreateStringTable, *dota.CSVCMsg_UpdateStringTable, *dota.CDemoStringTables:
			p.Stsh.AppendPacket(item)
		case *dota.CDemoFileHeader:
			p.FileHeader = obj
		case *dota.CSVCMsg_GameEventList:
			for _, descriptor := range obj.GetDescriptors() {
				p.GameEventMap[descriptor.GetEventid()] = descriptor
			}
		}
	}

	for _, item := range items {
		switch obj := item.Object.(type) {
		case *dota.CDemoClassInfo,
			*dota.CDemoFileHeader,
			*dota.CDemoStringTables,
			*dota.CSVCMsg_CreateStringTable,
			*dota.CSVCMsg_GameEventList,
			*dota.CSVCMsg_SendTable,
			*dota.CSVCMsg_ServerInfo,
			*dota.CSVCMsg_UpdateStringTable:
			// those have been handled above, please keep in sync.
		case *dota.CDemoStop,
			*dota.CDemoSyncTick,
			*dota.CDOTAUserMsg_CreateLinearProjectile,
			*dota.CDOTAUserMsg_DestroyLinearProjectile,
			*dota.CDOTAUserMsg_DodgeTrackingProjectiles,
			*dota.CDOTAUserMsg_GlobalLightColor,
			*dota.CDOTAUserMsg_GlobalLightDirection,
			*dota.CDOTAUserMsg_HalloweenDrops,
			*dota.CDOTAUserMsg_HudError,
			*dota.CDOTAUserMsg_LocationPing,
			*dota.CDOTAUserMsg_MapLine,
			*dota.CDOTAUserMsg_MinimapEvent,
			*dota.CDOTAUserMsg_NevermoreRequiem,
			*dota.CDOTAUserMsg_ParticleManager,
			*dota.CDOTAUserMsg_SendRoshanPopup,
			*dota.CDOTAUserMsg_SendStatPopup,
			*dota.CDOTAUserMsg_SharedCooldown,
			*dota.CDOTAUserMsg_UnitEvent,
			*dota.CDOTAUserMsg_WorldLine,
			*dota.CNETMsg_SignonState,
			*dota.CNETMsg_Tick,
			*dota.CSVCMsg_ClassInfo,
			*dota.CSVCMsg_SetView,
			*dota.CSVCMsg_Sounds,
			*dota.CSVCMsg_TempEntities,
			*dota.CUserMsg_SendAudio,
			*dota.CUserMsg_TextMsg,
			*dota.CUserMsg_VoiceMask:
			// see NOTES for why we ignore them.
		case *dota.CSVCMsg_PacketEntities:
			// to skip to a specific time, we have to handle more.
			if item.From == dota.EDemoCommands_DEM_Packet {
				p.ParsePacket(item.Tick, obj)
			}
		case *dota.CDemoFileInfo:
			p.FileInfo = obj
		case *dota.CSVCMsg_VoiceInit:
			p.VoiceInit = obj
		case *dota.CSVCMsg_GameEvent:
			p.onGameEvent(item.Tick, obj)
		case *dota.CDOTAUserMsg_ChatEvent:
			if p.OnChatEvent != nil {
				p.OnChatEvent(item.Tick, obj)
			}
		case *dota.CDOTAUserMsg_OverheadEvent:
			if p.OnOverheadEvent != nil {
				p.OnOverheadEvent(item.Tick, obj)
			}
		case *dota.CDOTAUserMsg_SpectatorPlayerClick:
			if p.OnSpectatorPlayerClick != nil {
				p.OnSpectatorPlayerClick(item.Tick, obj)
			}
		case *dota.CUserMsg_SayText2:
			if p.OnSayText2 != nil {
				p.OnSayText2(item.Tick, obj)
			}
		case *dota.CSVCMsg_VoiceData:
			if p.OnVoiceData != nil {
				p.OnVoiceData(obj)
			}
		case *dota.CNETMsg_SetConVar:
			if p.OnSetConVar != nil {
				p.OnSetConVar(obj)
			}
		default:
			spew.Dump(obj)
		}
	}

	// sort.Sort(p.Stsh.ActiveModifierDelta)
	if len(p.Stsh.ActiveModifierDelta) > 0 {
		sort.Sort(p.Stsh.ActiveModifierDelta)
		modNames := p.Stsh.GetTableNow("ModifierNames").Items
		for _, buff := range p.Stsh.ActiveModifierDelta {
			switch buff.GetEntryType() {
			case dota.DOTA_MODIFIER_ENTRY_TYPE_DOTA_MODIFIER_ENTRY_TYPE_ACTIVE:
				p.ActiveModifiers[int(buff.GetIndex())] = buff
				modName := modNames[int(buff.GetModifierClass())].Str

				switch modName {
				case "modifier_kill":
					// spew.Dump(modName)
					// spew.Dump(p.ByHandle[int(buff.GetCaster())].Name)
					// spew.Dump(p.ByHandle[int(buff.GetAbility())].Name)
					// spew.Dump(p.ByHandle[int(buff.GetParent())].Name)
				case "modifier_item_mekansm":
				case "modifier_enigma_black_hole_pull":
				case "modifier_enigma_black_hole_thinker":
				case "modifier_rune_regen":
					// spew.Dump(modName)
					// spew.Dump(p.ByHandle[int(buff.GetCaster())].Name)
				case "modifier_smoke_of_deceit":
					// spew.Dump(modName)
					// spew.Dump(p.ByHandle[int(buff.GetCaster())])
					// spew.Dump(p.ByHandle[int(buff.GetParent())])
				case "modifier_item_dustofappearance":
					// (entry_type:DOTA_MODIFIER_ENTRY_TYPE_ACTIVE parent:951211 index:4 serial_num:604 modifier_class:67 ability_level:1 creation_time:1061.2625 duration:12 caster:64700 ability:1995890 aura:false )
					// spew.Dump(modName)
					// spew.Dump(p.ByHandle[int(buff.GetCaster())].Name)
					// spew.Dump(p.ByHandle[int(buff.GetParent())].Name)
				case "modifier_truesight":
					// (entry_type:DOTA_MODIFIER_ENTRY_TYPE_ACTIVE parent:177062 index:2 serial_num:6100 modifier_class:912 creation_time:2245.6614 duration:0.5 caster:265746 aura:true )
					// spew.Dump(modName)
					// spew.Dump(p.ByHandle[int(buff.GetCaster())].Name)
					// spew.Dump(p.ByHandle[int(buff.GetParent())].Name)
				case "modifier_item_buff_ward":
					// (entry_type:DOTA_MODIFIER_ENTRY_TYPE_ACTIVE parent:324575 index:1 serial_num:1110 modifier_class:68 creation_time:1175.0347 duration:240 caster:324575 aura:false )
					// spew.Dump(modName)
					// spew.Dump(p.ByHandle[int(buff.GetCaster())].Name)
					// spew.Dump(p.ByHandle[int(buff.GetParent())].Name)
				case "modifier_item_ward_true_sight":
					// (entry_type:DOTA_MODIFIER_ENTRY_TYPE_ACTIVE parent:324575 index:2 serial_num:1111 modifier_class:71 creation_time:1175.0347 duration:240 caster:324575 aura:false range:800 )
					// spew.Dump(modName)
					// spew.Dump(p.ByHandle[int(buff.GetCaster())].Name)
					// spew.Dump(p.ByHandle[int(buff.GetParent())].Name)
				case "modifier_item_sentry_ward":
					// (entry_type:DOTA_MODIFIER_ENTRY_TYPE_ACTIVE parent:1146009 index:33 serial_num:1121 modifier_class:70 ability_level:1 creation_time:1177.3341 caster:1146009 ability:779306 aura:false )
					// spew.Dump(modName)
					// spew.Dump(p.ByHandle[int(buff.GetCaster())].Name)
					// spew.Dump(p.ByHandle[int(buff.GetAbility())].Name)
					// spew.Dump(p.ByHandle[int(buff.GetParent())].Name)
				case "modifier_item_gem_of_true_sight":
					// spew.Dump("add")
					// spew.Dump(modName)
					// spew.Dump(p.ByHandle[int(buff.GetCaster())].Name)
					// spew.Dump(p.ByHandle[int(buff.GetAbility())].Values["DT_BaseEntity.m_iName"])
					// spew.Dump(p.ByHandle[int(buff.GetParent())].Name)
				default:
					// spew.Dump("add " + modName)
					// spew.Dump(buff)
				}
			case dota.DOTA_MODIFIER_ENTRY_TYPE_DOTA_MODIFIER_ENTRY_TYPE_REMOVED:
				mod := p.ActiveModifiers[int(buff.GetIndex())]
				modName := modNames[int(mod.GetModifierClass())].Str
				// spew.Dump("remove " + modName)
				switch modName {
				case "modifier_item_gem_of_true_sight":
					// spew.Dump(mod)
					// spew.Dump(p.ByHandle[int(mod.GetCaster())].Name)
				}
				delete(p.ActiveModifiers, int(buff.GetIndex()))
			}
		}
	}
}

func (p *Parser) onGameEvent(tick int, obj *dota.CSVCMsg_GameEvent) {
	desc := p.GameEventMap[obj.GetEventid()]
	dName := desc.GetName()

	switch dName {
	case "hltv_versioninfo":
		// version : <*>type:5 val_byte:1
	case "hltv_message":
		// text : <*>type:1 val_string:"Please wait for broadcast to start ..."
	case "hltv_status":
		// clients : <*>type:3 val_long:523
		// slots : <*>type:3 val_long:3840
		// proxies : <*>type:4 val_short:59
		// master : <*>type:1 val_string:"146.66.152.49:28027"
	case "dota_combatlog":
		if p.OnCombatLog == nil {
			return
		}
		keys := obj.GetKeys()
		table := p.Stsh.GetTableNow("CombatLogNames").Items

		log := &CombatLogEntry{
			Type:               dota.DOTA_COMBATLOG_TYPES(keys[0].GetValByte()).String()[15:],
			AttackerIsIllusion: keys[5].GetValBool(),
			TargetIsIllusion:   keys[6].GetValBool(),
			Value:              keys[7].GetValShort(),
			Health:             keys[8].GetValShort(),
			Timestamp:          keys[9].GetValFloat(),
		}

		if k := table[int(keys[1].GetValShort())]; k != nil {
			log.SourceName = k.Str
		}
		if k := table[int(keys[2].GetValShort())]; k != nil {
			log.TargetName = k.Str
		}
		if k := table[int(keys[3].GetValShort())]; k != nil {
			log.AttackerName = k.Str
		}
		if k := table[int(keys[4].GetValShort())]; k != nil {
			log.InflictorName = k.Str
		}
		if k := table[int(keys[10].GetValShort())]; k != nil {
			log.TargetSourceName = k.Str
		}

		p.OnCombatLog(log)
	case "dota_chase_hero":
		// target1 : <*>type:4 val_short:1418
		// target2 : <*>type:4 val_short:0
		// type : <*>type:5 val_byte:0
		// priority : <*>type:4 val_short:15
		// gametime : <*>type:2 val_float:2710.3667
		// highlight : <*>type:6 val_bool:false
		// target1playerid : <*>type:5 val_byte:1
		// target2playerid : <*>type:5 val_byte:32
		// eventtype : <*>type:4 val_short:1
	case "dota_tournament_item_event":
		// event_type : <*>type:4 val_short:0  => witness first blood
		// event_type : <*>type:4 val_short:1  => witness killing spree
		// event_type : <*>type:4 val_short:3  => witness hero deny
	default:
		dKeys := desc.GetKeys()
		spew.Println(dName)
		for n, key := range obj.GetKeys() {
			spew.Println(dKeys[n].GetName(), ":", key)
		}
	}
}

func (p *Parser) onCDemoClassInfo(cdci *dota.CDemoClassInfo) {
	for _, class := range cdci.GetClasses() {
		id, name := int(class.GetClassId()), class.GetTableName()
		p.ClassInfosIdMapping[name] = id
		p.ClassInfosNameMapping[id] = name

		props := p.Sth.LoadSendTable(name)
		multiples := map[string]int{}
		for _, prop := range props {
			key := prop.DtName + "." + prop.VarName
			multiples[key] += 1
		}
		p.Multiples[id] = multiples
		p.Mapping[id] = props

		if p.OnTablename != nil {
			p.OnTablename(name)
		}
	}

	p.Stsh.ClassInfosNameMapping = p.ClassInfosNameMapping
	p.Stsh.Mapping = p.Mapping
	p.Stsh.Multiples = p.Multiples
}

type packet struct {
	br          *utils.BitReader
	index, tick int
}

func (p *Parser) ParsePacket(tick int, pe *dota.CSVCMsg_PacketEntities) {
	br := utils.NewBitReader(pe.GetEntityData())
	currentIndex := -1

	createPackets := []*packet_entities.PacketEntity{}
	preservePackets := []*packet_entities.PacketEntity{}
	deletePackets := []*packet_entities.PacketEntity{}

	for i := 0; i < int(pe.GetUpdatedEntries()); i++ {
		currentIndex = br.ReadNextEntityIndex(currentIndex)
		uType := packet_entities.ReadUpdateType(br)

		switch uType {
		case packet_entities.Create:
			createPackets = append(createPackets, p.entityCreate(br, currentIndex, tick))
		case packet_entities.Preserve:
			preservePackets = append(preservePackets, p.entityPreserve(br, currentIndex, tick))
		case packet_entities.Delete:
			deletePackets = append(deletePackets, p.entityDelete(br, currentIndex, tick))
		}
	}

	/*
		if pe.EntityHandle == 445806 {
			spew.Dump("tick:", tick)
			spew.Dump(pe)
		}

		if p.OnEntityCreated != nil {
			p.OnEntityCreated(pe)
		}
	*/

	for _, pe := range createPackets {
		p.Entities[pe.Index] = pe
		p.ByHandle[pe.Handle()] = pe
		if p.OnEntityCreated != nil {
			p.OnEntityCreated(pe)
		}
	}

	for _, pe := range preservePackets {
		p.OnEntityPreserved(pe)
	}

	for _, pe := range deletePackets {
		if p.OnEntityDeleted != nil {
			p.OnEntityDeleted(pe)
		}
		p.Entities[pe.Index] = nil
		delete(p.ByHandle, pe.Handle())
	}
}

func (p *Parser) entityCreate(br *utils.BitReader, currentIndex, tick int) *packet_entities.PacketEntity {
	pe := &packet_entities.PacketEntity{
		Tick:      tick,
		ClassId:   int(br.ReadUBits(p.ClassIdNumBits)),
		SerialNum: int(br.ReadUBits(10)),
		Index:     currentIndex,
		Type:      packet_entities.Create,
		Values:    map[string]interface{}{},
	}
	pe.EntityHandle = pe.Handle()
	pe.Name = p.ClassInfosNameMapping[pe.ClassId]

	indices := br.ReadPropertiesIndex()
	pMapping := p.Mapping[pe.ClassId]
	pMultiples := p.Multiples[pe.ClassId]
	values := br.ReadPropertiesValues(pMapping, pMultiples, indices)

	baseline, foundBaseline := p.Stsh.Baseline[pe.ClassId]
	if foundBaseline {
		for key, value := range baseline {
			pe.Values[key] = value
		}
	}
	for key, value := range values {
		pe.Values[key] = value
	}

	return pe
}

func (p *Parser) entityPreserve(br *utils.BitReader, currentIndex, tick int) *packet_entities.PacketEntity {
	pe := p.Entities[currentIndex]
	pe.Tick = tick
	pe.Type = packet_entities.Preserve
	indices := br.ReadPropertiesIndex()
	classId := p.ClassInfosIdMapping[pe.Name]
	pe.Delta = br.ReadPropertiesValues(p.Mapping[classId], p.Multiples[classId], indices)
	pe.OldDelta = map[string]interface{}{}

	for key, value := range pe.Delta {
		pe.OldDelta[key] = pe.Values[key]
		pe.Values[key] = value
	}

	return pe
}

func (p *Parser) entityDelete(br *utils.BitReader, currentIndex, tick int) *packet_entities.PacketEntity {
	return p.Entities[currentIndex]
}

/*
// logs chronologically close to this event (usually within 1-2 ticks).
func (p *Parser) CombatLogsCloseTo(now float32) (logs []*CombatLogEntry) {
	closestDelta := 0.5 // be generous and look around ±0.5 second

	for _, log := range p.CombatLog {
		if log.Type == "DEATH" && log.TargetName == "npc_dota_observer_wards" {
			logDelta := math.Abs(float64(now - log.Timestamp))
			if closestDelta > logDelta {
				logs = append(logs, log)
				closestDelta = logDelta
			}
		}
	}

	return logs
}
*/
