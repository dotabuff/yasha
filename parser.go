package core

import (
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"strings"

	"code.google.com/p/gogoprotobuf/proto"
	"github.com/davecgh/go-spew/spew"
	"github.com/elobuff/d2rp/core/packet_entities"
	"github.com/elobuff/d2rp/core/parser"
	"github.com/elobuff/d2rp/core/send_tables"
	"github.com/elobuff/d2rp/core/string_tables"
	"github.com/elobuff/d2rp/core/utils"
	dota "github.com/elobuff/d2rp/dota"
)

func foo() { spew.Dump("hi") }

type Parser struct {
	ItemsOnGround         map[int]bool
	TickTime              map[int]float64
	StartTime             float64
	LastHitMinutes        int
	PlayerIdClientId      map[int]int
	Abilities             map[*AbilityTracker]bool
	LastHits              map[*LastHitTracker]bool
	PlayerResourceIndex   int
	Baseline              map[int]map[string]interface{}
	ClassInfos            *dota.CDemoClassInfo
	ClassIdNumBits        int
	ClassInfosIdMapping   map[string]int
	ClassInfosNameMapping map[int]string
	Mapping               map[int][]*send_tables.SendProp
	Multiples             map[int]map[string]int
	Packets               parser.ParserBaseItems
	Entities              []*packet_entities.PacketEntity
	Items                 parser.ParserBaseItems
	Stsh                  *string_tables.StateHelper
	Sth                   *send_tables.Helper
	Parser                *parser.Parser
	TextMsg               []string
	AllChat               []*dota.CUserMsg_SayText2
	VoiceData             map[int][]byte
	GameEventMap          map[int32]*dota.CSVCMsg_GameEventListDescriptorT
	FileInfo              *dota.CDemoFileInfo
	StringTables          []*dota.CDemoStringTablesTableT
	EntityPreserved       func(*packet_entities.PacketEntity)
	EntityCreated         func(*packet_entities.PacketEntity)
	EntityDeleted         func(*packet_entities.PacketEntity)
}

func ParserFromFile(path string) *Parser {
	return NewParser(parser.ReadFile(path))
}

func NewParser(data []byte) *Parser {
	return &Parser{
		Parser:                parser.NewParser(data),
		ItemsOnGround:         map[int]bool{},
		TickTime:              map[int]float64{},
		PlayerIdClientId:      map[int]int{},
		Abilities:             map[*AbilityTracker]bool{},
		LastHits:              map[*LastHitTracker]bool{},
		PlayerResourceIndex:   -1,
		Baseline:              map[int]map[string]interface{}{},
		ClassInfosIdMapping:   map[string]int{},
		ClassInfosNameMapping: map[int]string{},
		Mapping:               map[int][]*send_tables.SendProp{},
		Multiples:             map[int]map[string]int{},
		Packets:               []*parser.ParserBaseItem{},
		Items:                 []*parser.ParserBaseItem{},
		TextMsg:               []string{},
		AllChat:               []*dota.CUserMsg_SayText2{},
		VoiceData:             map[int][]byte{},
		GameEventMap:          map[int32]*dota.CSVCMsg_GameEventListDescriptorT{},
	}
}

func (p *Parser) PicksBans() (picksbans []*PickBan, isCaptainsMode bool) {
	gameInfo := p.FileInfo.GetGameInfo()
	game := gameInfo.GetDota()

	if game.GetGameMode() != 2 {
		return nil, false
	}

	pbs := make([]*PickBan, 0, 10)
	for _, pb := range game.GetPicksBans() {
		pbs = append(pbs, &PickBan{
			IsPick: pb.GetIsPick(),
			Team:   int(pb.GetTeam() - 2),
			HeroId: int(pb.GetHeroId()),
		})
	}

	return pbs, true
}

func (p *Parser) Parse() {
	p.ParseBaseline()
	p.ParsePackets()
}

func (p *Parser) ParseBaseline() {
	p.Items = p.Parser.Parse(MessageFilter)
	sort.Sort(p.Items)

	p.Entities = make([]*packet_entities.PacketEntity, 2048)

	p.processItems()
	p.makeBaseline()
}

func (p *Parser) ParsePackets() {
	for _, packet := range p.Packets {
		p.ParsePacket(packet)
	}
}

func (p *Parser) Players() map[int]*Player {
	lpe := p.Entities[p.PlayerResourceIndex]
	lpev := lpe.Values

	total := map[int]map[string]interface{}{}

	for key, value := range lpev {
		keys := strings.SplitN(key, ".", 2)
		key := keys[0]
		i, _ := strconv.Atoi(keys[1])
		if _, ok := total[i]; !ok {
			total[i] = map[string]interface{}{}
		}
		total[i][key] = value
	}

	players := map[int]*Player{}

	for _, vs := range total {
		heroHandle, ok := vs["m_hSelectedHero"].(int)
		if !ok {
			continue
		}
		player := &Player{
			Abilities:  Abilities{},
			LastHits:   LastHits{},
			HeroHandle: heroHandle,
		}
		players[heroHandle] = player

		player.BroadcasterChannelSlot = vs["m_iBroadcasterChannelSlot"].(int)
		player.BroadcasterChannel = vs["m_iBroadcasterChannel"].(int)
		player.BroadcasterLanguage = vs["m_iBroadcasterLanguage"].(int)
		player.IsBroadcaster = vs["m_bIsBroadcaster"].(int) == 1

		player.BattleBonusActive = vs["m_bBattleBonusActive"].(int) == 1
		player.BattleBonusRate = vs["m_iBattleBonusRate"].(int)
		player.MetaExperienceAwarded = vs["m_iMetaExperienceAwarded"].(int)
		player.MetaExperienceBonusRate = vs["m_iMetaExperienceBonusRate"].(int)
		player.MetaExperience = vs["m_iMetaExperience"].(int)
		player.MetaLevel = vs["m_iMetaLevel"].(int)
		player.TimedRewardCrates = vs["m_iTimedRewardCrates"].(int)
		player.TimedRewardDrops = vs["m_iTimedRewardDrops"].(int)

		player.ConnectionState = vs["m_iConnectionState"].(int)
		player.FakeClient = vs["m_bFakeClient"].(int) == 1
		player.FullyJoinedServer = vs["m_bFullyJoinedServer"].(int) == 1
		player.AFK = vs["m_bAFK"].(int) == 1

		player.Kills = vs["m_iKills"].(uint)
		player.Deaths = vs["m_iDeaths"].(uint)
		player.Assists = vs["m_iAssists"].(uint)
		player.Level = vs["m_iLevel"].(int)
		player.ReliableGold = vs["m_iReliableGold"].(uint)
		player.TotalEarnedGold = vs["m_iTotalEarnedGold"].(uint)
		player.TotalEarnedXP = vs["m_iTotalEarnedXP"].(uint)
		player.UnreliableGold = vs["m_iUnreliableGold"].(uint)

		player.BuybackCooldownTime = vs["m_flBuybackCooldownTime"].(float32)
		player.DenyCount = vs["m_iDenyCount"].(uint)
		player.HasRandomed = vs["m_bHasRandomed"].(int) == 1
		player.HasRepicked = vs["m_bHasRepicked"].(int) == 1
		player.LastBuybackTime = vs["m_iLastBuybackTime"].(int)
		player.LastHitCount = vs["m_iLastHitCount"].(uint)
		if lastHitMultikill, found := vs["m_iLastHitMultikill"]; found {
			player.LastHitMultikill = lastHitMultikill.(uint)
		}
		if lastHitStreak, found := vs["m_iLastHitStreak"]; found {
			player.LastHitStreak = lastHitStreak.(uint)
		}
		player.PlayerNames = vs["m_iszPlayerNames"].(string)
		player.PlayerSteamIDs = vs["m_iPlayerSteamIDs"].(uint64)
		player.PlayerTeams = vs["m_iPlayerTeams"].(int)
		player.PossibleHeroSelection = vs["m_nPossibleHeroSelection"].(int)
		player.RespawnSeconds = vs["m_iRespawnSeconds"].(int)
		player.SelectedHeroID = vs["m_nSelectedHeroID"].(int)
		player.SelectedHero = vs["m_iszSelectedHero"].(string)
		player.Streak = vs["m_iStreak"].(uint)
		player.SuggestedHeroes = vs["m_nSuggestedHeroes"].(int)
		player.UnitShareMasks = vs["m_UnitShareMasks"].(int)
		player.VoiceChatBanned = vs["m_bVoiceChatBanned"].(int) == 1
	}

	for ability, _ := range p.Abilities {
		player := players[ability.HeroHandle]
		if player == nil {
			// "courier_burst" the time when courier is upgraded?
			spew.Println("cannot assign ability to player:")
			spew.Dump(ability)
		} else {
			player.Abilities = append(player.Abilities, ability)
		}
	}

	for lastHit, _ := range p.LastHits {
		player := players[lastHit.HeroHandle]
		if player == nil {
			spew.Println("cannot assign lasthit to player:")
			spew.Dump(lastHit)
		} else {
			player.LastHits = append(player.LastHits, lastHit)
		}
	}

	for _, player := range players {
		sort.Sort(player.Abilities)
		sort.Sort(player.LastHits)
	}

	return players
}

func (p *Parser) ParsePacket(packet *parser.ParserBaseItem) {
	pe := (packet.Object).(*dota.CSVCMsg_PacketEntities)
	br := utils.NewBitReader(pe.GetEntityData())
	currentIndex := -1
	for i := 0; i < int(pe.GetUpdatedEntries()); i++ {
		currentIndex = br.ReadNextEntityIndex(currentIndex)
		uType := packet_entities.ReadUpdateType(br)
		switch uType {
		case packet_entities.Preserve:
			p.EntityPreserve(br, currentIndex, packet.Tick)
		case packet_entities.Create:
			p.EntityCreate(br, currentIndex, packet.Tick)
		case packet_entities.Delete:
			p.EntityDelete(br, currentIndex, packet.Tick)
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

	if pe.Name == "DT_DOTAPlayer" {
		if iPlayerID, found := values["DT_DOTAPlayer.m_iPlayerID"]; found {
			if _, found := p.PlayerIdClientId[pe.Index]; !found {
				p.PlayerIdClientId[pe.Index] = iPlayerID.(int)
			}
		}
	} else {
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
		if pe.Name == "DT_DOTA_Item_Physical" {
			p.ItemsOnGround[(pe.Values["DT_DOTA_Item_Physical.m_hItem"]).(int)] = true
		}
	}

	// spew.Println("=== Create", pe.Name, "(", currentIndex, ")", "@", pe.Tick)
	// spew.Dump(values)

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

	if len(values) == 0 {
		return
	}

	// spew.Println("=== Preserve", pe.Name, "(", currentIndex, ")", "@", pe.Tick)
	// spew.Dump(values)

	switch pe.Name {
	case "DT_DOTAGamerulesProxy":
		if gameTime, found := values["DT_DOTAGamerules.m_fGameTime"]; found {
			p.TickTime[pe.Tick] = float64(gameTime.(float32))
		}

		if startTime, found := values["DT_DOTAGamerules.m_flGameStartTime"]; found {
			p.StartTime = float64(startTime.(float32))
		}

		if p.StartTime > 0 {
			if tickTime, found := p.TickTime[pe.Tick]; found {
				realGameTime := tickTime - p.StartTime
				minutesPassed := int(realGameTime / 60)
				if minutesPassed > p.LastHitMinutes {
					lpe := p.Entities[p.PlayerResourceIndex]
					for i := 0; i < 10; i++ {
						is := ".000" + strconv.Itoa(i)
						p.LastHits[&LastHitTracker{
							HeroHandle: int(lpe.Values["m_hSelectedHero"+is].(int)),
							LastHit:    int(lpe.Values["m_iLastHitCount"+is].(uint)),
							Tick:       pe.Tick,
						}] = true
					}
					p.LastHitMinutes++
				}
			}
		}
	case "DT_DOTA_PlayerResource":
		for key, value := range values {
			pe.Values[key] = value
		}
	default:
		if strings.Contains(pe.Name, "Ability") {
			if level, found := values["DT_DOTABaseAbility.m_iLevel"]; found {
				p.Abilities[&AbilityTracker{
					HeroHandle: pe.Values["DT_BaseEntity.m_hOwnerEntity"].(int),
					Level:      level.(int),
					Name:       pe.Values["DT_BaseEntity.m_iName"].(string),
					Tick:       pe.Tick,
				}] = true
			}
		}

		for key, value := range values {
			pe.Values[key] = value
		}

		/*
			if strings.HasPrefix(pe.Name, "DT_DOTA_Unit_Hero_") {
				if pe.Tick > 10000 && pe.Tick < 10500 {
					if v, found := values["DT_DOTA_BaseNPC.m_vecOrigin"]; found {
						vec := v.(*core.Vector)
						cx := pe.Values["DT_DOTA_BaseNPC.m_cellX"].(int)
						cy := pe.Values["DT_DOTA_BaseNPC.m_cellY"].(int)
						cz := pe.Values["DT_DOTA_BaseNPC.m_cellZ"].(int)
						spew.Printf("\"%s\": {%3d, %3d, %3d, %3d, %3d, %3d},\n", pe.Name, cx, cy, cz, int(vec.X), int(vec.Y), int(vec.Z))
					}
				}
			}
		*/
	}

	if p.EntityPreserved != nil {
		p.EntityPreserved(pe)
	}
}

func (p *Parser) EntityDelete(br *utils.BitReader, currentIndex, tick int) {
	pe := p.Entities[currentIndex]
	// spew.Println("=== Delete", pe.Name, "(", currentIndex, ")", "@", pe.Tick)
	// spew.Dump(pe)
	if pe.Name == "DT_DOTA_Item_Physical" {
		delete(p.ItemsOnGround, pe.Values["DT_DOTA_Item_Physical.m_hItem"].(int))
	}
	p.Entities[currentIndex] = nil
	if p.EntityDeleted != nil {
		p.EntityDeleted(pe)
	}
}

// we ignore whatever we don't care about or understand yet, this helps speed up parsing a bit.
func MessageFilter(msg proto.Message) bool {
	switch msg.(type) {
	case *dota.CDemoFileHeader:
		// (demo_file_stamp:"PBUFDEM\000" network_protocol:40 server_name:"Valve Dota 2 Server #8 (srcds038)" client_name:"SourceTV Demo" map_name:"dota" game_directory:"dota" fullpackets_version:2 allow_clientside_entities:true allow_clientside_particles:true )
	case *dota.CNETMsg_SetConVar:
		// can't see anything of value here.
		// (convars:<cvars:<name:"sv_consistency" value:"0" > cvars:<name:"sv_skyname" value:"sky_dotasky_01" > cvars:<name:"tv_transmitall" value:"1" > cvars:<name:"steamworks_sessionid_server" value:"21537759004" > cvars:<name:"think_limit" value:"0" > cvars:<name:"tv_transmitall" value:"1" > > )
	case *dota.CNETMsg_SignonState:
		// we see signon_state:3|4|5 before the SyncTick
		// (signon_state:3 spawn_count:2 num_server_players:0 )
	case *dota.CSVCMsg_ClassInfo:
		// (create_on_client:true )
		return true
	case *dota.CSVCMsg_VoiceInit:
		// (quality:5 codec:"vaudio_celt" version:3 )
	case *dota.CSVCMsg_SetView:
		// (entity_index:1 )
	case *dota.CDemoSyncTick:
		// that's where the game timer starts?
	case *dota.CDemoStop:
		// that's where the game timer stops?
	case *dota.CDOTAUserMsg_HudError:
		// looks like some kind of notification for upcoming StringTableUpdate
		// (order_id:19 )
		// (*dota.CSVCMsg_UpdateStringTable)(table_id:19 num_changed_entries:1 string_data:"\360\301\005\000A\000\302\261\320\006C\000$\032\n" )
	case *dota.CDOTAUserMsg_SpectatorPlayerClick:
		// no idea what that is good for.
		// (entindex:9 order_type:1 target_index:0 )
		// (entindex:26 order_type:4 target_index:1125 )
	case *dota.CDOTAUserMsg_MapLine:
		// I guess those are lines drawn?
		// (player_id:2 mapline:<x:-7813 y:-5988 initial:false > )
	case *dota.CDOTAUserMsg_WorldLine:
		// now some spectator drawing on the world?
		// (player_id:22 worldline:<x:-1820 y:-1819 z:127 initial:true > )
	case *dota.CUserMsg_VoiceMask:
		// this might be used for mute status of players, haven't found a use for it.
	case *dota.CSVCMsg_Sounds:
		// seem to be ingame sound sources, not particularly interesting.
	case *dota.CUserMsg_SendAudio:
		// global sounds?
		// (stop:false name:"General.Morning")
	case *dota.CDOTAUserMsg_MinimapEvent:
		// (event_type:32 entity_handle:1039755 x:-1864 y:-1740 duration:3 )
		// (event_type:32 entity_handle:35559 x:-1283 y:1418 duration:5 )
		// (event_type:4 entity_handle:1095158 x:-878 y:-694 duration:4 )
		// (event_type:64 entity_handle:1039755 x:-1410 y:1329 duration:4 )
	case *dota.CDOTAUserMsg_SendStatPopup:
		// (player_id:16 statpopup:<style:k_EDOTA_SPT_Textline stat_strings:"I spoke to Bruno in the break, he called Na`Vi's strat last game \"Electric Bear\"" > )
		// (player_id:16 statpopup:<style:k_EDOTA_SPT_Basic stat_strings:"Chat Poll Results" stat_strings:"62% of viewers expected TobiWan to miss first blood.\rThey were wrong." stat_images:1001 > )
	case *dota.CDOTAUserMsg_SendRoshanPopup:
		// (reclaimed:true gametime:2075 )
		// (reclaimed:false gametime:2314 )
	case *dota.CDOTAUserMsg_LocationPing:
		// pings, might be useful later
		// (player_id:1 location_ping:<x:-3513 y:-5111 target:1063 direct_ping:true type:0 > )
	case *dota.CDOTAUserMsg_ParticleManager:
		// http://www.cyborgmatt.com/2013/01/dota-2-replay-parser-bruno/#cdotausermsg-particlemanager
		// this might be useful, just gotta find out for what :)
		// (type:DOTA_PARTICLE_MANAGER_EVENT_CREATE index:6636 create_particle:<particle_name_index:2146 attach_type:1 entity_handle:510877 > )
		// (type:DOTA_PARTICLE_MANAGER_EVENT_RELEASE index:6636 )
		// (type:DOTA_PARTICLE_MANAGER_EVENT_UPDATE_ENT index:6637 update_particle_ent:<control_point:0 entity_handle:1078463 attach_type:5 attachment:3 > )
	case *dota.CDOTAUserMsg_OverheadEvent:
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
	case *dota.CSVCMsg_TempEntities:
		// Temporary entities are used by the server to create short-lived or
		// one-off effects on clients. They are different to standard entities in
		// that they are 'fire and forget'; once one has been created, the server
		// has nothing more to do with it. TEs have no edict or entity index, and
		// do not count toward the entity limit.
		//
		// TEs are unreliable and get dropped if too many are created at once.
		// The maximum per update is 32 in multiplayer and 255 in single player.
	case *dota.CDOTAUserMsg_CreateLinearProjectile:
		// (origin:<x:-849.5365 y:-2277 z:127.99994 > velocity:<x:687.15155 y:-1458.2067 > entindex:743 particle_index:1998 handle:70 )
	case *dota.CDOTAUserMsg_DestroyLinearProjectile:
		// (handle:70 )
	case *dota.CDOTAUserMsg_DodgeTrackingProjectiles:
		// (entindex:743 )
	default:
		return true
	}

	return false
}

func (p *Parser) processItems() {
	var serverInfo *dota.CSVCMsg_ServerInfo
	sthItems := map[string]*dota.CSVCMsg_SendTable{}
	p.Sth = send_tables.NewHelper(sthItems)

	for _, item := range p.Items {
		switch value := item.Object.(type) {
		case *dota.CSVCMsg_PacketEntities:
			if item.From == dota.EDemoCommands_DEM_Packet {
				p.Packets = append(p.Packets, item)
			}
		case *dota.CDemoFileInfo:
			if p.FileInfo == nil {
				p.FileInfo = value
			}
		case *dota.CSVCMsg_ServerInfo:
			if serverInfo == nil {
				serverInfo = value
			}
		case *dota.CSVCMsg_GameEventList:
			for _, descriptor := range value.GetDescriptors() {
				p.GameEventMap[descriptor.GetEventid()] = descriptor
			}
		case *dota.CSVCMsg_GameEvent:
			desc := p.GameEventMap[value.GetEventid()]
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
				// type : <*>type:5 val_byte:0
				// sourcename : <*>type:4 val_short:2
				// targetname : <*>type:4 val_short:74
				// attackername : <*>type:4 val_short:2
				// inflictorname : <*>type:4 val_short:0
				// attackerillusion : <*>type:6 val_bool:false
				// targetillusion : <*>type:6 val_bool:false
				// value : <*>type:4 val_short:111
				// health : <*>type:4 val_short:31
				// timestamp : <*>type:2 val_float:2166.2576
				// targetsourcename : <*>type:4 val_short:74
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
				for n, key := range value.GetKeys() {
					spew.Println(dKeys[n].GetName(), ":", key)
				}
			}
		case *dota.CDOTAUserMsg_ChatEvent:
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
		case *dota.CSVCMsg_SendTable:
			sthItems[value.GetNetTableName()] = value
		case *dota.CUserMsg_TextMsg:
			p.TextMsg = append(p.TextMsg, value.GetParam()...)
		case *dota.CUserMsg_SayText2:
			p.AllChat = append(p.AllChat, value)
		case *dota.CDOTAUserMsg_UnitEvent:
			switch value.GetMsgType() {
			case dota.EDotaEntityMessages_DOTA_UNIT_ADD_GESTURE:
				// unit starts animation
			case dota.EDotaEntityMessages_DOTA_UNIT_FADE_GESTURE:
			case dota.EDotaEntityMessages_DOTA_UNIT_REMOVE_GESTURE:
			case dota.EDotaEntityMessages_DOTA_UNIT_REMOVE_ALL_GESTURES:
			case dota.EDotaEntityMessages_DOTA_UNIT_SPEECH:
				// heroes say stuff
			case dota.EDotaEntityMessages_DOTA_UNIT_SPEECH_CLIENTSIDE_RULES:
				// the announcer says stuff
			default:
				spew.Dump(value)
			}
		case *dota.CSVCMsg_VoiceData:
			client := int(value.GetClient())
			d, ok := p.VoiceData[client]
			if ok {
				p.VoiceData[client] = append(d, value.GetVoiceData()...)
			} else {
				p.VoiceData[client] = value.GetVoiceData()
			}
		case *dota.CDemoStringTables:
			p.StringTables = value.GetTables()
		case *dota.CDemoClassInfo, *dota.CSVCMsg_UpdateStringTable, *dota.CNETMsg_Tick, *dota.CSVCMsg_CreateStringTable, *dota.CSVCMsg_ClassInfo:
			// don't print
		default:
			spew.Dump(value)
		}
	}

	sort.Sort(p.Packets)
	p.ClassIdNumBits = int(math.Log(float64(serverInfo.GetMaxClasses()))/math.Log(2)) + 1
	p.Stsh = string_tables.NewStateHelper(p.Items)

	for _, item := range p.Items {
		switch value := item.Object.(type) {
		case *dota.CDemoClassInfo:
			p.ClassInfos = value
			break
		}
	}

	return
}

func (p *Parser) makeBaseline() {
	for _, class := range p.ClassInfos.GetClasses() {
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

	stringTables := p.Stsh.GetStateAtTick(int(p.FileInfo.GetPlaybackTicks()))
	var instanceBaseline *string_tables.StringTable
	for _, value := range stringTables {
		if value.Name == "instancebaseline" {
			instanceBaseline = value
			break
		}
	}

	if instanceBaseline == nil {
		panic("no instanceBaseline")
	}

	for _, item := range instanceBaseline.Items {
		classId, err := strconv.Atoi(item.Str)
		if err != nil {
			panic(err)
		}
		className := p.ClassInfosNameMapping[classId]
		if className != "DT_DOTAPlayer" {
			br := utils.NewBitReader(item.Data)
			indices := br.ReadPropertiesIndex()
			mapping := p.Mapping[classId]
			multiples := p.Multiples[classId]
			baseValues := br.ReadPropertiesValues(mapping, multiples, indices)
			p.Baseline[classId] = baseValues
		}
	}
}

func (p *Parser) WriteVoiceData(dir string) {
	for client, voice := range p.VoiceData {
		oga := dir + "/voice_" + strconv.Itoa(client) + ".oga"
		spew.Println("writing", oga, "...")
		ioutil.WriteFile(oga, voice, 0644)
	}
}

func (p *Parser) WriteStringTables(dir string) {
	for _, table := range p.StringTables {
		prefix := dir + "/" + table.GetTableName()
		ic := table.GetItemsClientside()
		if len(ic) > 0 {
			ioutil.WriteFile(prefix+"_items_clientside.txt", []byte(spew.Sdump(ic)), 0644)
		}
		i := table.GetItems()
		if len(i) > 0 {
			ioutil.WriteFile(prefix+"_items.txt", []byte(spew.Sdump(i)), 0644)
		}
	}
}
