package core

import (
	"io/ioutil"
	"math"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/elobuff/d2rp/core/packet_entities"
	"github.com/elobuff/d2rp/core/parser"
	"github.com/elobuff/d2rp/core/send_tables"
	"github.com/elobuff/d2rp/core/string_tables"
	"github.com/elobuff/d2rp/core/utils"
	dota "github.com/elobuff/d2rp/dota"
)

type Parser struct {
	Parser                *parser.Parser
	AllChat               []*dota.CUserMsg_SayText2
	Baseline              map[int]map[string]interface{}
	ClassIdNumBits        int
	ClassInfos            *dota.CDemoClassInfo
	ClassInfosIdMapping   map[string]int
	ClassInfosNameMapping map[int]string
	CombatLog             []*CombatLogEntry
	ConVars               map[string]interface{}
	Entities              []*packet_entities.PacketEntity
	EntityCreated         func(*packet_entities.PacketEntity)
	EntityDeleted         func(*packet_entities.PacketEntity)
	EntityPreserved       func(*packet_entities.PacketEntity, map[string]interface{})
	FileHeader            *dota.CDemoFileHeader
	FileInfo              *dota.CDemoFileInfo
	GameEventMap          map[int32]*dota.CSVCMsg_GameEventListDescriptorT
	LastHitMinutes        int
	Mapping               map[int][]*send_tables.SendProp
	Multiples             map[int]map[string]int
	PlayerResourceIndex   int
	RawClicks             []*RawClick
	ServerInfo            *dota.CSVCMsg_ServerInfo
	StartTime             float64
	Sth                   *send_tables.Helper
	Stsh                  *string_tables.StateHelper
	VoiceData             map[int][]byte
	VoiceInit             *dota.CSVCMsg_VoiceInit
}

func ParserFromFile(path string) *Parser {
	return NewParser(parser.ReadFile(path))
}

func NewParser(data []byte) *Parser {
	return &Parser{
		Parser:                parser.NewParser(data),
		AllChat:               []*dota.CUserMsg_SayText2{},
		Baseline:              map[int]map[string]interface{}{},
		ClassInfosIdMapping:   map[string]int{},
		ClassInfosNameMapping: map[int]string{},
		CombatLog:             []*CombatLogEntry{},
		ConVars:               map[string]interface{}{},
		GameEventMap:          map[int32]*dota.CSVCMsg_GameEventListDescriptorT{},
		Mapping:               map[int][]*send_tables.SendProp{},
		Multiples:             map[int]map[string]int{},
		PlayerResourceIndex:   -1,
		RawClicks:             []*RawClick{},
		VoiceData:             map[int][]byte{},
	}
}

func (p *Parser) Parse() {
	p.Sth = send_tables.NewHelper()
	p.Stsh = string_tables.NewStateHelper()
	p.Entities = make([]*packet_entities.PacketEntity, 2048)

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
	for _, item := range items {
		switch obj := item.Object.(type) {
		case *dota.CSVCMsg_SendTable:
			p.Sth.SetSendTable(obj.GetNetTableName(), obj)
		case *dota.CSVCMsg_CreateStringTable, *dota.CSVCMsg_UpdateStringTable, *dota.CDemoStringTables:
			p.Stsh.AppendPacket(item)
		case *dota.CDemoFileHeader:
			// (demo_file_stamp:"PBUFDEM\000" network_protocol:40 server_name:"Valve Dota 2 Server #8 (srcds038)" client_name:"SourceTV Demo" map_name:"dota" game_directory:"dota" fullpackets_version:2 allow_clientside_entities:true allow_clientside_particles:true )
			p.FileHeader = obj
		case *dota.CSVCMsg_ServerInfo:
			p.ServerInfo = obj
			p.ClassIdNumBits = int(math.Log(float64(obj.GetMaxClasses()))/math.Log(2)) + 1
		case *dota.CDemoClassInfo:
			p.onCDemoClassInfo(obj)
		case *dota.CSVCMsg_GameEventList:
			for _, descriptor := range obj.GetDescriptors() {
				p.GameEventMap[descriptor.GetEventid()] = descriptor
			}
		}
	}

	for _, item := range items {
		switch obj := item.Object.(type) {
		case *dota.CSVCMsg_CreateStringTable, *dota.CSVCMsg_UpdateStringTable,
			*dota.CDemoStringTables, *dota.CDemoFileHeader, *dota.CSVCMsg_ServerInfo,
			*dota.CDemoClassInfo, *dota.CSVCMsg_SendTable,
			*dota.CSVCMsg_GameEventList:
			// those have been handled above, please keep in sync.
		case *dota.CDemoSyncTick,
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
			*dota.CNETMsg_Tick,
			*dota.CSVCMsg_ClassInfo,
			*dota.CSVCMsg_SetView,
			*dota.CSVCMsg_Sounds,
			*dota.CSVCMsg_TempEntities,
			*dota.CUserMsg_SendAudio,
			*dota.CUserMsg_TextMsg,
			*dota.CUserMsg_VoiceMask:
			// see NOTES for why we ignore them.
		case *dota.CDemoFileInfo:
			p.onFileInfo(obj)
		case *dota.CDOTAUserMsg_ChatEvent:
			p.onChatEvent(item.Tick, obj)
		case *dota.CDOTAUserMsg_OverheadEvent:
			p.onOverheadEvent(item.Tick, obj)
		case *dota.CDOTAUserMsg_SpectatorPlayerClick:
			p.onSpectatorPlayerClick(item.Tick, obj)
		case *dota.CNETMsg_SetConVar:
			p.onSetConVar(obj)
		case *dota.CNETMsg_SignonState:
			p.onSignonState(obj)
		case *dota.CSVCMsg_GameEvent:
			p.onGameEvent(item.Tick, obj)
		case *dota.CSVCMsg_VoiceInit:
			p.VoiceInit = obj
		case *dota.CUserMsg_SayText2:
			p.onSayText2(obj)
		case *dota.CSVCMsg_PacketEntities:
			// to skip to a specific time, we have to handle more.
			if item.From == dota.EDemoCommands_DEM_Packet {
				p.ParsePacket(item.Tick, obj)
			}
		case *dota.CDemoStop:
			// the demo ends here.
		case *dota.CSVCMsg_VoiceData:
			/*
				client := int(obj.GetClient())
				d, ok := p.VoiceData[client]
				if ok {
					p.VoiceData[client] = append(d, obj.GetVoiceData()...)
				} else {
					p.VoiceData[client] = obj.GetVoiceData()
				}
			*/
		default:
			spew.Dump(obj)
		}
	}
}

func (p *Parser) onSpectatorPlayerClick(tick int, obj *dota.CDOTAUserMsg_SpectatorPlayerClick) {
	p.RawClicks = append(p.RawClicks, &RawClick{
		Tick:   tick,
		Entity: int(obj.GetEntindex()),
		Type:   int(obj.GetOrderType()),
		Target: int(obj.GetTargetIndex()),
	})
}

func (p *Parser) onSetConVar(obj *dota.CNETMsg_SetConVar) {
	for _, cvar := range obj.GetConvars().GetCvars() {
		p.ConVars[cvar.GetName()] = cvar.GetValue()
	}
}

func (p *Parser) onSignonState(obj *dota.CNETMsg_SignonState) {}

func (p *Parser) onSayText2(obj *dota.CUserMsg_SayText2) {
	p.AllChat = append(p.AllChat, obj)
}

// http://www.cyborgmatt.com/2013/01/dota-2-replay-parser-bruno/#cdotausermsg-overheadevent
//
// Known alerts:
// OVERHEAD_ALERT_BLOCK
// OVERHEAD_ALERT_BONUS_SPELL_DAMAGE
// OVERHEAD_ALERT_CRITICAL
// OVERHEAD_ALERT_DAMAGE
// OVERHEAD_ALERT_DENY
// OVERHEAD_ALERT_GOLD
// OVERHEAD_ALERT_HEAL
// OVERHEAD_ALERT_MANA_ADD
// OVERHEAD_ALERT_MISS
// (message_type:OVERHEAD_ALERT_GOLD value:121 target_player_entindex:1 target_entindex:1052 )
// (message_type:OVERHEAD_ALERT_MISS value:0 target_entindex:1039 )
func (p *Parser) onOverheadEvent(tick int, obj *dota.CDOTAUserMsg_OverheadEvent) {
}

// this shows up as a summary, unfortunately only at the end of the file.
// (playback_time:346.93335 playback_ticks:10408 playback_frames:4886 )
func (p *Parser) onFileInfo(obj *dota.CDemoFileInfo) {
	p.FileInfo = obj
}

// http://www.cyborgmatt.com/2013/01/dota-2-replay-parser-bruno/#chatevents
// they are displayed during game on the left side.
//
// just some:
//
// (type:CHAT_MESSAGE_AEGIS value:0 playerid_1:6 playerid_2:-1 )
// (type:CHAT_MESSAGE_CONNECT value:0 playerid_1:0 playerid_2:-1 )
// (type:CHAT_MESSAGE_DISCONNECT value:0 playerid_1:0 playerid_2:-1 )
// (type:CHAT_MESSAGE_FIRSTBLOOD value:409 playerid_1:6 playerid_2:3 )
// (type:CHAT_MESSAGE_GLYPH_USED value:0 playerid_1:2 playerid_2:-1 )
// (type:CHAT_MESSAGE_HERO_DENY value:0 playerid_1:7 playerid_2:7 )
// (type:CHAT_MESSAGE_HERO_KILL value:409 playerid_1:3 playerid_2:6 )
// (type:CHAT_MESSAGE_ITEM_PURCHASE value:208 playerid_1:6 playerid_2:-1 )
// (type:CHAT_MESSAGE_ITEM_PURCHASE value:92 playerid_1:2 playerid_2:-1 )
// (type:CHAT_MESSAGE_ITEM_PURCHASE value:92 playerid_1:7 playerid_2:-1 )
// (type:CHAT_MESSAGE_PAUSED value:0 playerid_1:4 playerid_2:-1 )
// (type:CHAT_MESSAGE_RECONNECT value:0 playerid_1:10 playerid_2:-1 )
// (type:CHAT_MESSAGE_ROSHAN_KILL value:200 playerid_1:3 playerid_2:-1 )
// (type:CHAT_MESSAGE_RUNE_BOTTLE value:0 playerid_1:7 playerid_2:-1 )
// (type:CHAT_MESSAGE_RUNE_BOTTLE value:1 playerid_1:7 playerid_2:-1 )
// (type:CHAT_MESSAGE_RUNE_BOTTLE value:2 playerid_1:7 playerid_2:-1 )
// (type:CHAT_MESSAGE_RUNE_BOTTLE value:3 playerid_1:7 playerid_2:-1 )
// (type:CHAT_MESSAGE_RUNE_BOTTLE value:4 playerid_1:7 playerid_2:-1 )
// (type:CHAT_MESSAGE_RUNE_PICKUP value:0 playerid_1:7 playerid_2:-1 )
// (type:CHAT_MESSAGE_RUNE_PICKUP value:1 playerid_1:7 playerid_2:-1 )
// (type:CHAT_MESSAGE_RUNE_PICKUP value:2 playerid_1:5 playerid_2:-1 )
// (type:CHAT_MESSAGE_RUNE_PICKUP value:3 playerid_1:2 playerid_2:-1 )
// (type:CHAT_MESSAGE_RUNE_PICKUP value:4 playerid_1:7 playerid_2:-1 )
// (type:CHAT_MESSAGE_STREAK_KILL value:374 playerid_1:0 playerid_2:1 playerid_3:1 playerid_4:8 playerid_5:3 playerid_6:-1 )
// (type:CHAT_MESSAGE_TOWER_KILL value:2 playerid_1:4 playerid_2:-1 )
// (type:CHAT_MESSAGE_TOWER_KILL value:3 playerid_1:9 playerid_2:-1 )
// (type:CHAT_MESSAGE_UNPAUSE_COUNTDOWN value:1 playerid_1:-1 playerid_2:-1 )
// (type:CHAT_MESSAGE_UNPAUSE_COUNTDOWN value:2 playerid_1:-1 playerid_2:-1 )
// (type:CHAT_MESSAGE_UNPAUSE_COUNTDOWN value:3 playerid_1:-1 playerid_2:-1 )
// (type:CHAT_MESSAGE_UNPAUSED value:0 playerid_1:3 playerid_2:-1 )
func (p *Parser) onChatEvent(tick int, obj *dota.CDOTAUserMsg_ChatEvent) {
	// FIXME: forward that to the analyzer?
	/*
		switch obj.GetType() {
		case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_DISCONNECT, dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_DISCONNECT_WAIT_FOR_RECONNECT:
			a.players[playersIndex(int(obj.GetPlayerid_1()))].disconnectTime = tick
			a.chatLog = append(a.chatLog, &ChatMessage{Tick: tick, Playerid: obj.GetPlayerid_1(), Type: "disconnect"})
		case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_RECONNECT:
			dc := a.players[playersIndex(int(obj.GetPlayerid_1()))].disconnectTime
			message := fmt.Sprintf("%f", float64(tick-dc)/30.0)
			a.chatLog = append(a.chatLog, &ChatMessage{Tick: tick, Playerid: obj.GetPlayerid_1(), Type: "reconnect", Message: message})
		case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_ABANDON:
			a.chatLog = append(a.chatLog, &ChatMessage{Tick: tick, Playerid: obj.GetPlayerid_1(), Type: "abandon"})
		case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_PAUSED:
			a.pauseTick = tick
			a.chatLog = append(a.chatLog, &ChatMessage{Tick: tick, Playerid: obj.GetPlayerid_1(), Type: "pause"})
		case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_UNPAUSED:
			message := fmt.Sprintf("%f", float32(tick-a.pauseTick)/30.0)
			a.chatLog = append(a.chatLog, &ChatMessage{Tick: tick, Playerid: obj.GetPlayerid_1(), Type: "unpause", Message: message})
		case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_SAFE_TO_LEAVE_ABANDONER_EARLY:
			a.chatLog = append(a.chatLog, &ChatMessage{Tick: tick, Playerid: obj.GetPlayerid_1(), Type: "abandon"})
			a.significant = false
		case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_TOWER_DENY:
			a.chatLog = append(a.chatLog, &ChatMessage{Tick: tick, Playerid: obj.GetPlayerid_1(), Type: "tower_deny"})
		case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_RUNE_BOTTLE:
			r := RUNE_NAMES[obj.GetValue()]
			a.chatLog = append(a.chatLog, &ChatMessage{Tick: tick, Playerid: obj.GetPlayerid_1(), Type: "rune_bottle", Message: r})
		default:
			spew.Dump(obj)
		}
	*/
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

		p.CombatLog = append(p.CombatLog, log)
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
		p.ClassInfosNameMapping[id] = name
		p.ClassInfosIdMapping[name] = id
		props := p.Sth.LoadSendTable(name)
		multiples := map[string]int{}
		for _, prop := range props {
			key := prop.DtName + "." + prop.VarName
			multiples[key] += 1
		}
		p.Multiples[id] = multiples
		p.Mapping[id] = props
	}

	p.Stsh.ClassInfosNameMapping = p.ClassInfosNameMapping
	p.Stsh.Mapping = p.Mapping
	p.Stsh.Multiples = p.Multiples
}

func (p *Parser) ParsePacket(tick int, pe *dota.CSVCMsg_PacketEntities) {
	br := utils.NewBitReader(pe.GetEntityData())
	currentIndex := -1
	for i := 0; i < int(pe.GetUpdatedEntries()); i++ {
		currentIndex = br.ReadNextEntityIndex(currentIndex)
		uType := packet_entities.ReadUpdateType(br)
		switch uType {
		case packet_entities.Create:
			p.EntityCreate(br, currentIndex, tick)
		case packet_entities.Preserve:
			p.EntityPreserve(br, currentIndex, tick)
		case packet_entities.Delete:
			p.EntityDelete(br, currentIndex, tick)
		case packet_entities.Leave:
			// TODO: figure out what this is good for.
		default:
			panic("new meta!")
		}
	}
}

func (p *Parser) EntityCreate(br *utils.BitReader, currentIndex, tick int) {
	pe := &packet_entities.PacketEntity{
		Tick:      tick,
		ClassId:   int(br.ReadUBits(p.ClassIdNumBits)),
		SerialNum: int(br.ReadUBits(10)),
		Index:     currentIndex,
		Type:      packet_entities.Create,
		Values:    map[string]interface{}{},
	}
	pe.Name = p.ClassInfosNameMapping[pe.ClassId]

	indices := br.ReadPropertiesIndex()
	pMapping := p.Mapping[pe.ClassId]
	pMultiples := p.Multiples[pe.ClassId]
	values := br.ReadPropertiesValues(pMapping, pMultiples, indices)

	baseline, foundBaseline := p.Baseline[pe.ClassId]
	if foundBaseline {
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

	p.Entities[pe.Index] = pe
	if p.PlayerResourceIndex == -1 && pe.Name == "DT_DOTA_PlayerResource" {
		p.PlayerResourceIndex = pe.Index
	}

	if p.EntityCreated != nil {
		p.EntityCreated(pe)
	}
}

func (p *Parser) EntityPreserve(br *utils.BitReader, currentIndex, tick int) {
	pe := p.Entities[currentIndex]
	pe.Tick = tick
	pe.Type = packet_entities.Preserve
	indices := br.ReadPropertiesIndex()
	classId := p.ClassInfosIdMapping[pe.Name]
	values := br.ReadPropertiesValues(p.Mapping[classId], p.Multiples[classId], indices)

	for key, value := range values {
		pe.Values[key] = value
	}

	if p.EntityPreserved != nil {
		p.EntityPreserved(pe, values)
	}
}

func (p *Parser) EntityDelete(br *utils.BitReader, currentIndex, tick int) {
	pe := p.Entities[currentIndex]

	if p.EntityDeleted != nil {
		p.EntityDeleted(pe)
	}
}

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

// just for debugging, i haven't found a way to actually play these files yet,
// they might need some sort of header.
func (p *Parser) WriteVoiceData(dir string) {
	for client, voice := range p.VoiceData {
		oga := dir + "/voice_" + strconv.Itoa(client) + ".oga"
		spew.Println("writing", oga, "...")
		ioutil.WriteFile(oga, voice, 0644)
	}
}

// just for debugging
func (p *Parser) WriteStringTables(dir string) {
	for _, table := range p.Stsh.GetStateAtTick(int(p.FileInfo.GetPlaybackTicks())) {
		if len(table.Items) > 0 {
			err := ioutil.WriteFile(dir+"/"+table.Name+".txt", []byte(spew.Sdump(table.Items)), 0644)
			if err != nil {
				panic(err)
			}
		}
	}
}
