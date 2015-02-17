package yasha

import (
	"math"
	"sort"

	"github.com/davecgh/go-spew/spew"

	"github.com/dotabuff/yasha/dota"
	"github.com/dotabuff/yasha/parser"
	"github.com/dotabuff/yasha/send_tables"
	"github.com/dotabuff/yasha/string_tables"
	"github.com/dotabuff/yasha/utils"
)

type Parser struct {
	Parser                *parser.Parser
	combatLogParser       *combatLogParser
	ClassIdNumBits        int
	ClassInfosIdMapping   map[string]int
	ClassInfosNameMapping map[int]string
	FileHeader            *dota.CDemoFileHeader
	GameEventMap          map[int32]*dota.CSVCMsg_GameEventListDescriptorT
	Mapping               map[int][]*send_tables.SendProp
	Multiples             map[int]map[string]int
	ServerInfo            *dota.CSVCMsg_ServerInfo
	Sth                   *send_tables.Helper
	Stsh                  *string_tables.StateHelper
	VoiceInit             *dota.CSVCMsg_VoiceInit

	ActiveModifiers map[int]*dota.CDOTAModifierBuffTableEntry
	Entities        []*PacketEntity
	ByHandle        map[int]*PacketEntity

	OnEntityCreated   func(*PacketEntity)
	OnEntityDeleted   func(*PacketEntity)
	OnEntityPreserved func(*PacketEntity)

	OnActiveModifierDelta func(map[int]*string_tables.StringTableItem, string_tables.ModifierBuffs)

	OnChatEvent                 func(tick int, obj *dota.CDOTAUserMsg_ChatEvent)
	OnOverheadEvent             func(tick int, obj *dota.CDOTAUserMsg_OverheadEvent)
	OnSayText2                  func(tick int, obj *dota.CUserMsg_SayText2)
	OnSounds                    func(tick int, obj *dota.CSVCMsg_Sounds)
	OnSpectatorPlayerClick      func(tick int, obj *dota.CDOTAUserMsg_SpectatorPlayerClick)
	OnChatWheel                 func(tick int, obj *dota.CDOTAUserMsg_ChatWheel)
	OnEnemyItemAlert            func(tick int, obj *dota.CDOTAUserMsg_EnemyItemAlert)
	OnSpectatorPlayerUnitOrders func(tick int, obj *dota.CDOTAUserMsg_SpectatorPlayerUnitOrders)
	OnPredictionResult          func(tick int, obj *dota.CDOTAUserMsg_PredictionResult)

	OnFileInfo  func(obj *dota.CDemoFileInfo)
	OnSetConVar func(obj *dota.CNETMsg_SetConVar)
	OnVoiceData func(obj *dota.CSVCMsg_VoiceData)

	OnCombatLog func(log CombatLogEntry)

	OnTablename func(name string)

	BeforeTick func(tick int)
	AfterTick  func(tick int)
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
	p.Entities = make([]*PacketEntity, 2048)
	p.ClassInfosIdMapping = map[string]int{}
	p.ClassInfosNameMapping = map[int]string{}
	p.GameEventMap = map[int32]*dota.CSVCMsg_GameEventListDescriptorT{}
	p.Mapping = map[int][]*send_tables.SendProp{}
	p.Multiples = map[int]map[string]int{}
	p.ByHandle = map[int]*PacketEntity{}
	p.combatLogParser = &combatLogParser{
		stsh:     p.Stsh,
		distinct: map[dota.DOTA_COMBATLOG_TYPES][]map[interface{}]bool{},
	}

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
			p.processTick(currentTick, buffer)
			buffer = make([]*parser.ParserBaseItem, 1, len(buffer))
			buffer[0] = item
			currentTick = item.Tick
		} else {
			// tick is ongoing
			buffer = append(buffer, item)
		}
	})

	// and process all remaining in the last tick
	currentTick++
	p.processTick(currentTick, buffer)
}

func (p *Parser) PrintDistinctCombatLogTypes() {
	for k, v := range p.combatLogParser.distinct {
		spew.Println(k)
		for kk, vv := range v {
			if len(vv) > 1 {
				spew.Println(kk)
				pp(vv)
			}
		}
	}
}

func (p *Parser) processTick(tick int, items []*parser.ParserBaseItem) {
	p.Stsh.ActiveModifierDelta = string_tables.ModifierBuffs{}

	if p.BeforeTick != nil {
		p.BeforeTick(tick)
	}

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
			*dota.CDOTAUserMsg_BoosterState,
			*dota.CDOTAUserMsg_CourierKilledAlert, // NOTE: might be useful later
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
			*dota.CDemoSyncTick,
			*dota.CNETMsg_SignonState,
			*dota.CNETMsg_Tick,
			*dota.CSVCMsg_ClassInfo,
			*dota.CSVCMsg_Print,
			*dota.CSVCMsg_SetView,
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
			if p.OnFileInfo != nil {
				p.OnFileInfo(obj)
			}
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
		case *dota.CSVCMsg_Sounds:
			if p.OnSounds != nil {
				p.OnSounds(item.Tick, obj)
			}
		case *dota.CSVCMsg_VoiceData:
			if p.OnVoiceData != nil {
				p.OnVoiceData(obj)
			}
		case *dota.CNETMsg_SetConVar:
			if p.OnSetConVar != nil {
				p.OnSetConVar(obj)
			}
		case *dota.CDOTAUserMsg_ChatWheel:
			// (chat_message:k_EDOTA_CW_All_GGWP player_id:2 param_hero_id:0 )
			if p.OnChatWheel != nil {
				p.OnChatWheel(item.Tick, obj)
			}
		case *dota.CDOTAUserMsg_EnemyItemAlert:
			// (player_id:13 target_player_id:9 itemid:751 rune_type:4294967295 )
			if p.OnEnemyItemAlert != nil {
				p.OnEnemyItemAlert(item.Tick, obj)
			}
		case *dota.CDOTAUserMsg_AbilitySteal:
			/*
				(player_id:0 ability_id:412 ability_level:4 )
				(player_id:0 ability_id:305 ability_level:4 )
				(player_id:0 ability_id:117 ability_level:4 )
				(player_id:0 ability_id:131 ability_level:4 )
				(player_id:0 ability_id:413 ability_level:2 )
				(player_id:0 ability_id:413 ability_level:2 )
				(player_id:0 ability_id:275 ability_level:4 )
				(player_id:0 ability_id:411 ability_level:4 )
				(player_id:0 ability_id:131 ability_level:4 )
				(player_id:0 ability_id:305 ability_level:4 )
				(player_id:0 ability_id:412 ability_level:4 )
				(player_id:0 ability_id:132 ability_level:4 )
			*/
		case *dota.CDemoSaveGame:
			// this is not VDF... some new fun stuff instead.
		case *dota.CDOTAUserMsg_SpectatorPlayerUnitOrders:
			// (entindex:3 order_type:8 units:403 ability_index:464 queue:false )
			if p.OnSpectatorPlayerUnitOrders != nil {
				p.OnSpectatorPlayerUnitOrders(item.Tick, obj)
			}
		case *dota.CDOTAUserMsg_PredictionResult:
			// (account_id:47276380 match_id:1232716559 correct:true predictions:<item_def:11133 num_correct:1 num_fails:0 > )
			// item_def is the id from the items_game.txt in vpk
			if p.OnPredictionResult != nil {
				p.OnPredictionResult(item.Tick, obj)
			}
		default:
			spew.Dump(obj)
		}
	}

	if p.OnActiveModifierDelta != nil {
		if len(p.Stsh.ActiveModifierDelta) > 0 {
			sort.Sort(p.Stsh.ActiveModifierDelta)
			p.OnActiveModifierDelta(p.Stsh.GetTableNow("ModifierNames").Items, p.Stsh.ActiveModifierDelta)
		}
	}

	if p.AfterTick != nil {
		p.AfterTick(tick)
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
		if p.OnCombatLog != nil {
			if log := p.combatLogParser.parse(obj); log != nil {
				p.OnCombatLog(log)
			}
		}
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

func (p *Parser) ParsePacket(tick int, pe *dota.CSVCMsg_PacketEntities) {
	br := utils.NewBitReader(pe.GetEntityData())
	currentIndex := -1

	createPackets := []*PacketEntity{}
	preservePackets := []*PacketEntity{}
	deletePackets := []*PacketEntity{}

	for i := 0; i < int(pe.GetUpdatedEntries()); i++ {
		currentIndex = br.ReadNextEntityIndex(currentIndex)
		uType := ReadUpdateType(br)

		switch uType {
		case Create:
			createPackets = append(createPackets, p.entityCreate(br, currentIndex, tick))
		case Preserve:
			preservePackets = append(preservePackets, p.entityPreserve(br, currentIndex, tick))
		case Delete:
			deletePackets = append(deletePackets, p.entityDelete(br, currentIndex, tick))
		}
	}

	for _, pe := range createPackets {
		p.Entities[pe.Index] = pe
		p.ByHandle[pe.Handle()] = pe
		if p.OnEntityCreated != nil {
			p.OnEntityCreated(pe)
		}
	}

	for _, pe := range preservePackets {
		if p.OnEntityPreserved != nil {
			p.OnEntityPreserved(pe)
		}
	}

	for _, pe := range deletePackets {
		if p.OnEntityDeleted != nil {
			p.OnEntityDeleted(pe)
		}
		// p.Entities[pe.Index] = nil
		// delete(p.ByHandle, pe.Handle())
	}
}

func (p *Parser) entityCreate(br *utils.BitReader, currentIndex, tick int) *PacketEntity {
	pe := &PacketEntity{
		Tick:      tick,
		ClassId:   int(br.ReadUBits(p.ClassIdNumBits)),
		SerialNum: int(br.ReadUBits(10)),
		Index:     currentIndex,
		Type:      Create,
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

func (p *Parser) entityPreserve(br *utils.BitReader, currentIndex, tick int) *PacketEntity {
	pe := p.Entities[currentIndex]
	pe.Tick = tick
	pe.Type = Preserve
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

func (p *Parser) entityDelete(br *utils.BitReader, currentIndex, tick int) *PacketEntity {
	return p.Entities[currentIndex]
}
