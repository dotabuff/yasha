package parser

import (
	"code.google.com/p/gogoprotobuf/proto"
	dota "github.com/dotabuff/d2rp/dota"
)

type Error string

func (e Error) Error() string { return string(e) }

type BaseEvent struct {
	Name  string
	Value int
}

type SignonPacket struct{}

func (s SignonPacket) ProtoMessage()  {}
func (s SignonPacket) Reset()         {}
func (s SignonPacket) String() string { return "" }

func (p *ParserBase) AsBaseEvent(commandName string) (proto.Message, error) {
	switch commandName {
	case "DEM_Stop":
		return &dota.CDemoStop{}, nil
	case "DEM_FileHeader":
		return &dota.CDemoFileHeader{}, nil
	case "DEM_FileInfo":
		return &dota.CDemoFileInfo{}, nil
	case "DEM_SyncTick":
		return &dota.CDemoSyncTick{}, nil
	case "DEM_SendTables":
		return &dota.CDemoSendTables{}, nil
	case "DEM_ClassInfo":
		return &dota.CDemoClassInfo{}, nil
	case "DEM_StringTables":
		return &dota.CDemoStringTables{}, nil
	case "DEM_Packet":
		return &dota.CDemoPacket{}, nil
	case "DEM_SignonPacket":
		return &SignonPacket{}, nil
	case "DEM_ConsoleCmd":
		return &dota.CDemoConsoleCmd{}, nil
	case "DEM_CustomData":
		return &dota.CDemoCustomData{}, nil
	case "DEM_CustomDataCallbacks":
		return &dota.CDemoCustomDataCallbacks{}, nil
	case "DEM_UserCmd":
		return &dota.CDemoUserCmd{}, nil
	case "DEM_FullPacket":
		return &dota.CDemoFullPacket{}, nil
	case "net_NOP":
		return &dota.CNETMsg_NOP{}, nil
	case "net_Disconnect":
		return &dota.CNETMsg_Disconnect{}, nil
	case "net_File":
		return &dota.CNETMsg_File{}, nil
	case "net_SplitScreenUser":
		return &dota.CNETMsg_SplitScreenUser{}, nil
	case "net_Tick":
		return &dota.CNETMsg_Tick{}, nil
	case "net_StringCmd":
		return &dota.CNETMsg_StringCmd{}, nil
	case "net_SetConVar":
		return &dota.CNETMsg_SetConVar{}, nil
	case "net_SignonState":
		return &dota.CNETMsg_SignonState{}, nil
	case "svc_ServerInfo":
		return &dota.CSVCMsg_ServerInfo{}, nil
	case "svc_SendTable":
		return &dota.CSVCMsg_SendTable{}, nil
	case "svc_ClassInfo":
		return &dota.CSVCMsg_ClassInfo{}, nil
	case "svc_SetPause":
		return &dota.CSVCMsg_SetPause{}, nil
	case "svc_CreateStringTable":
		return &dota.CSVCMsg_CreateStringTable{}, nil
	case "svc_UpdateStringTable":
		return &dota.CSVCMsg_UpdateStringTable{}, nil
	case "svc_VoiceInit":
		return &dota.CSVCMsg_VoiceInit{}, nil
	case "svc_VoiceData":
		return &dota.CSVCMsg_VoiceData{}, nil
	case "svc_Print":
		return &dota.CSVCMsg_Print{}, nil
	case "svc_Sounds":
		return &dota.CSVCMsg_Sounds{}, nil
	case "svc_SetView":
		return &dota.CSVCMsg_SetView{}, nil
	case "svc_FixAngle":
		return &dota.CSVCMsg_FixAngle{}, nil
	case "svc_CrosshairAngle":
		return &dota.CSVCMsg_CrosshairAngle{}, nil
	case "svc_BSPDecal":
		return &dota.CSVCMsg_BSPDecal{}, nil
	case "svc_SplitScreen":
		return &dota.CSVCMsg_SplitScreen{}, nil
	case "svc_UserMessage":
		return &dota.CSVCMsg_UserMessage{}, nil
	case "svc_GameEvent":
		return &dota.CSVCMsg_GameEvent{}, nil
	case "svc_PacketEntities":
		return &dota.CSVCMsg_PacketEntities{}, nil
	case "svc_TempEntities":
		return &dota.CSVCMsg_TempEntities{}, nil
	case "svc_Prefetch":
		return &dota.CSVCMsg_Prefetch{}, nil
	case "svc_Menu":
		return &dota.CSVCMsg_Menu{}, nil
	case "svc_GameEventList":
		return &dota.CSVCMsg_GameEventList{}, nil
	case "svc_GetCvarValue":
		return &dota.CSVCMsg_GetCvarValue{}, nil
	case "svc_PacketReliable":
		return &dota.CSVCMsg_PacketReliable{}, nil
	case "UM_AchievementEvent":
		return &dota.CUserMsg_AchievementEvent{}, nil
	case "UM_CloseCaption":
		return &dota.CUserMsg_CloseCaption{}, nil
	case "UM_CurrentTimescale":
		return &dota.CUserMsg_CurrentTimescale{}, nil
	case "UM_DesiredTimescale":
		return &dota.CUserMsg_DesiredTimescale{}, nil
	case "UM_Fade":
		return &dota.CUserMsg_Fade{}, nil
	case "UM_GameTitle":
		return &dota.CUserMsg_GameTitle{}, nil
	case "UM_Geiger":
		return &dota.CUserMsg_Geiger{}, nil
	case "UM_HintText":
		return &dota.CUserMsg_HintText{}, nil
	case "UM_HudMsg":
		return &dota.CUserMsg_HudMsg{}, nil
	case "UM_HudText":
		return &dota.CUserMsg_HudText{}, nil
	case "UM_KeyHintText":
		return &dota.CUserMsg_KeyHintText{}, nil
	case "UM_MessageText":
		return &dota.CUserMsg_MessageText{}, nil
	case "UM_RequestState":
		return &dota.CUserMsg_RequestState{}, nil
	case "UM_ResetHUD":
		return &dota.CUserMsg_ResetHUD{}, nil
	case "UM_Rumble":
		return &dota.CUserMsg_Rumble{}, nil
	case "UM_SayText":
		return &dota.CUserMsg_SayText{}, nil
	case "UM_SayText2":
		return &dota.CUserMsg_SayText2{}, nil
	case "UM_SayTextChannel":
		return &dota.CUserMsg_SayTextChannel{}, nil
	case "UM_Shake":
		return &dota.CUserMsg_Shake{}, nil
	case "UM_ShakeDir":
		return &dota.CUserMsg_ShakeDir{}, nil
	case "UM_StatsCrawlMsg":
		return &dota.CUserMsg_StatsCrawlMsg{}, nil
	case "UM_StatsSkipState":
		return &dota.CUserMsg_StatsSkipState{}, nil
	case "UM_TextMsg":
		return &dota.CUserMsg_TextMsg{}, nil
	case "UM_Tilt":
		return &dota.CUserMsg_Tilt{}, nil
	case "UM_Train":
		return &dota.CUserMsg_Train{}, nil
	case "UM_VGUIMenu":
		return &dota.CUserMsg_VGUIMenu{}, nil
	case "UM_VoiceMask":
		return &dota.CUserMsg_VoiceMask{}, nil
	case "UM_VoiceSubtitle":
		return &dota.CUserMsg_VoiceSubtitle{}, nil
	case "UM_SendAudio":
		return &dota.CUserMsg_SendAudio{}, nil
	case "DOTA_UM_AIDebugLine":
		return &dota.CDOTAUserMsg_AIDebugLine{}, nil
	case "DOTA_UM_ChatEvent":
		return &dota.CDOTAUserMsg_ChatEvent{}, nil
	case "DOTA_UM_CombatHeroPositions":
		return &dota.CDOTAUserMsg_CombatHeroPositions{}, nil
	case "DOTA_UM_CombatLogData":
		return &dota.CDOTAUserMsg_CombatLogData{}, nil
	case "DOTA_UM_CombatLogShowDeath":
		return &dota.CDOTAUserMsg_CombatLogShowDeath{}, nil
	case "DOTA_UM_CreateLinearProjectile":
		return &dota.CDOTAUserMsg_CreateLinearProjectile{}, nil
	case "DOTA_UM_DestroyLinearProjectile":
		return &dota.CDOTAUserMsg_DestroyLinearProjectile{}, nil
	case "DOTA_UM_DodgeTrackingProjectiles":
		return &dota.CDOTAUserMsg_DodgeTrackingProjectiles{}, nil
	case "DOTA_UM_GlobalLightColor":
		return &dota.CDOTAUserMsg_GlobalLightColor{}, nil
	case "DOTA_UM_GlobalLightDirection":
		return &dota.CDOTAUserMsg_GlobalLightDirection{}, nil
	case "DOTA_UM_InvalidCommand":
		return &dota.CDOTAUserMsg_InvalidCommand{}, nil
	case "DOTA_UM_LocationPing":
		return &dota.CDOTAUserMsg_LocationPing{}, nil
	case "DOTA_UM_MapLine":
		return &dota.CDOTAUserMsg_MapLine{}, nil
	case "DOTA_UM_MiniKillCamInfo":
		return &dota.CDOTAUserMsg_MiniKillCamInfo{}, nil
	case "DOTA_UM_MinimapDebugPoint":
		return &dota.CDOTAUserMsg_MinimapDebugPoint{}, nil
	case "DOTA_UM_MinimapEvent":
		return &dota.CDOTAUserMsg_MinimapEvent{}, nil
	case "DOTA_UM_NevermoreRequiem":
		return &dota.CDOTAUserMsg_NevermoreRequiem{}, nil
	case "DOTA_UM_OverheadEvent":
		return &dota.CDOTAUserMsg_OverheadEvent{}, nil
	case "DOTA_UM_SetNextAutobuyItem":
		return &dota.CDOTAUserMsg_SetNextAutobuyItem{}, nil
	case "DOTA_UM_SharedCooldown":
		return &dota.CDOTAUserMsg_SharedCooldown{}, nil
	case "DOTA_UM_SpectatorPlayerClick":
		return &dota.CDOTAUserMsg_SpectatorPlayerClick{}, nil
	case "DOTA_UM_TutorialTipInfo":
		return &dota.CDOTAUserMsg_TutorialTipInfo{}, nil
	case "DOTA_UM_UnitEvent":
		return &dota.CDOTAUserMsg_UnitEvent{}, nil
	case "DOTA_UM_ParticleManager":
		return &dota.CDOTAUserMsg_ParticleManager{}, nil
	case "DOTA_UM_BotChat":
		return &dota.CDOTAUserMsg_BotChat{}, nil
	case "DOTA_UM_HudError":
		return &dota.CDOTAUserMsg_HudError{}, nil
	case "DOTA_UM_ItemPurchased":
		return &dota.CDOTAUserMsg_ItemPurchased{}, nil
	case "DOTA_UM_Ping":
		return &dota.CDOTAUserMsg_Ping{}, nil
	case "DOTA_UM_ItemFound":
		return &dota.CDOTAUserMsg_ItemFound{}, nil
	case "DOTA_UM_SwapVerify":
		return &dota.CDOTAUserMsg_SwapVerify{}, nil
	case "DOTA_UM_WorldLine":
		return &dota.CDOTAUserMsg_WorldLine{}, nil
	case "DOTA_UM_TournamentDrop":
		return &dota.CDOTAUserMsg_TournamentDrop{}, nil
	case "DOTA_UM_ItemAlert":
		return &dota.CDOTAUserMsg_ItemAlert{}, nil
	case "DOTA_UM_HalloweenDrops":
		return &dota.CDOTAUserMsg_HalloweenDrops{}, nil
	case "DOTA_UM_ChatWheel":
		return &dota.CDOTAUserMsg_ChatWheel{}, nil
	case "DOTA_UM_ReceivedXmasGift":
		return &dota.CDOTAUserMsg_ReceivedXmasGift{}, nil
	case "DOTA_UM_UpdateSharedContent":
		return &dota.CDOTAUserMsg_UpdateSharedContent{}, nil
	case "DOTA_UM_TutorialRequestExp":
		return &dota.CDOTAUserMsg_TutorialRequestExp{}, nil
	case "DOTA_UM_TutorialPingMinimap":
		return &dota.CDOTAUserMsg_TutorialPingMinimap{}, nil
	case "DOTA_UM_ShowSurvey":
		return &dota.CDOTAUserMsg_ShowSurvey{}, nil
	case "DOTA_UM_TutorialFade":
		return &dota.CDOTAUserMsg_TutorialFade{}, nil
	case "DOTA_UM_AddQuestLogEntry":
		return &dota.CDOTAUserMsg_AddQuestLogEntry{}, nil
	case "DOTA_UM_SendStatPopup":
		return &dota.CDOTAUserMsg_SendStatPopup{}, nil
	case "DOTA_UM_TutorialFinish":
		return &dota.CDOTAUserMsg_TutorialFinish{}, nil
	case "DOTA_UM_SendRoshanPopup":
		return &dota.CDOTAUserMsg_SendRoshanPopup{}, nil
	case "DOTA_UM_SendGenericToolTip":
		return &dota.CDOTAUserMsg_SendGenericToolTip{}, nil
	}
	return nil, Error("Type not found: " + commandName)
}

func (p *ParserBase) AsBaseEventNETSVC(value int) (proto.Message, error) {
	switch value {
	case 1:
		return &dota.CNETMsg_Disconnect{}, nil
	case 2:
		return &dota.CNETMsg_File{}, nil
	case 3:
		return &dota.CNETMsg_SplitScreenUser{}, nil
	case 4:
		return &dota.CNETMsg_Tick{}, nil
	case 5:
		return &dota.CNETMsg_StringCmd{}, nil
	case 6:
		return &dota.CNETMsg_SetConVar{}, nil
	case 7:
		return &dota.CNETMsg_SignonState{}, nil
	case 8:
		return &dota.CSVCMsg_ServerInfo{}, nil
	case 9:
		return &dota.CSVCMsg_SendTable{}, nil
	case 10:
		return &dota.CSVCMsg_ClassInfo{}, nil
	case 11:
		return &dota.CSVCMsg_SetPause{}, nil
	case 12:
		return &dota.CSVCMsg_CreateStringTable{}, nil
	case 13:
		return &dota.CSVCMsg_UpdateStringTable{}, nil
	case 14:
		return &dota.CSVCMsg_VoiceInit{}, nil
	case 15:
		return &dota.CSVCMsg_VoiceData{}, nil
	case 16:
		return &dota.CSVCMsg_Print{}, nil
	case 17:
		return &dota.CSVCMsg_Sounds{}, nil
	case 18:
		return &dota.CSVCMsg_SetView{}, nil
	case 19:
		return &dota.CSVCMsg_FixAngle{}, nil
	case 20:
		return &dota.CSVCMsg_CrosshairAngle{}, nil
	case 21:
		return &dota.CSVCMsg_BSPDecal{}, nil
	case 22:
		return &dota.CSVCMsg_SplitScreen{}, nil
	case 23:
		return &dota.CSVCMsg_UserMessage{}, nil
	case 25:
		return &dota.CSVCMsg_GameEvent{}, nil
	case 26:
		return &dota.CSVCMsg_PacketEntities{}, nil
	case 27:
		return &dota.CSVCMsg_TempEntities{}, nil
	case 28:
		return &dota.CSVCMsg_Prefetch{}, nil
	case 29:
		return &dota.CSVCMsg_Menu{}, nil
	case 30:
		return &dota.CSVCMsg_GameEventList{}, nil
	case 31:
		return &dota.CSVCMsg_GetCvarValue{}, nil
	case 32:
		return &dota.CSVCMsg_PacketReliable{}, nil
	}
	return nil, Error("not found")
}

func (p *ParserBase) AsBaseEventBUMDUM(value int) (proto.Message, error) {
	switch value {
	case 1:
		return &dota.CUserMsg_AchievementEvent{}, nil
	case 2:
		return &dota.CUserMsg_CloseCaption{}, nil
	case 4:
		return &dota.CUserMsg_CurrentTimescale{}, nil
	case 5:
		return &dota.CUserMsg_DesiredTimescale{}, nil
	case 6:
		return &dota.CUserMsg_Fade{}, nil
	case 7:
		return &dota.CUserMsg_GameTitle{}, nil
	case 8:
		return &dota.CUserMsg_Geiger{}, nil
	case 9:
		return &dota.CUserMsg_HintText{}, nil
	case 10:
		return &dota.CUserMsg_HudMsg{}, nil
	case 11:
		return &dota.CUserMsg_HudText{}, nil
	case 12:
		return &dota.CUserMsg_KeyHintText{}, nil
	case 13:
		return &dota.CUserMsg_MessageText{}, nil
	case 14:
		return &dota.CUserMsg_RequestState{}, nil
	case 15:
		return &dota.CUserMsg_ResetHUD{}, nil
	case 16:
		return &dota.CUserMsg_Rumble{}, nil
	case 17:
		return &dota.CUserMsg_SayText{}, nil
	case 18:
		return &dota.CUserMsg_SayText2{}, nil
	case 19:
		return &dota.CUserMsg_SayTextChannel{}, nil
	case 20:
		return &dota.CUserMsg_Shake{}, nil
	case 21:
		return &dota.CUserMsg_ShakeDir{}, nil
	case 22:
		return &dota.CUserMsg_StatsCrawlMsg{}, nil
	case 23:
		return &dota.CUserMsg_StatsSkipState{}, nil
	case 24:
		return &dota.CUserMsg_TextMsg{}, nil
	case 25:
		return &dota.CUserMsg_Tilt{}, nil
	case 26:
		return &dota.CUserMsg_Train{}, nil
	case 27:
		return &dota.CUserMsg_VGUIMenu{}, nil
	case 28:
		return &dota.CUserMsg_VoiceMask{}, nil
	case 29:
		return &dota.CUserMsg_VoiceSubtitle{}, nil
	case 30:
		return &dota.CUserMsg_SendAudio{}, nil
	case 65:
		return &dota.CDOTAUserMsg_AIDebugLine{}, nil
	case 66:
		return &dota.CDOTAUserMsg_ChatEvent{}, nil
	case 67:
		return &dota.CDOTAUserMsg_CombatHeroPositions{}, nil
	case 68:
		return &dota.CDOTAUserMsg_CombatLogData{}, nil
	case 70:
		return &dota.CDOTAUserMsg_CombatLogShowDeath{}, nil
	case 71:
		return &dota.CDOTAUserMsg_CreateLinearProjectile{}, nil
	case 72:
		return &dota.CDOTAUserMsg_DestroyLinearProjectile{}, nil
	case 73:
		return &dota.CDOTAUserMsg_DodgeTrackingProjectiles{}, nil
	case 74:
		return &dota.CDOTAUserMsg_GlobalLightColor{}, nil
	case 75:
		return &dota.CDOTAUserMsg_GlobalLightDirection{}, nil
	case 76:
		return &dota.CDOTAUserMsg_InvalidCommand{}, nil
	case 77:
		return &dota.CDOTAUserMsg_LocationPing{}, nil
	case 78:
		return &dota.CDOTAUserMsg_MapLine{}, nil
	case 79:
		return &dota.CDOTAUserMsg_MiniKillCamInfo{}, nil
	case 80:
		return &dota.CDOTAUserMsg_MinimapDebugPoint{}, nil
	case 81:
		return &dota.CDOTAUserMsg_MinimapEvent{}, nil
	case 82:
		return &dota.CDOTAUserMsg_NevermoreRequiem{}, nil
	case 83:
		return &dota.CDOTAUserMsg_OverheadEvent{}, nil
	case 84:
		return &dota.CDOTAUserMsg_SetNextAutobuyItem{}, nil
	case 85:
		return &dota.CDOTAUserMsg_SharedCooldown{}, nil
	case 86:
		return &dota.CDOTAUserMsg_SpectatorPlayerClick{}, nil
	case 87:
		return &dota.CDOTAUserMsg_TutorialTipInfo{}, nil
	case 88:
		return &dota.CDOTAUserMsg_UnitEvent{}, nil
	case 89:
		return &dota.CDOTAUserMsg_ParticleManager{}, nil
	case 90:
		return &dota.CDOTAUserMsg_BotChat{}, nil
	case 91:
		return &dota.CDOTAUserMsg_HudError{}, nil
	case 92:
		return &dota.CDOTAUserMsg_ItemPurchased{}, nil
	case 93:
		return &dota.CDOTAUserMsg_Ping{}, nil
	case 94:
		return &dota.CDOTAUserMsg_ItemFound{}, nil
	case 96:
		return &dota.CDOTAUserMsg_SwapVerify{}, nil
	case 97:
		return &dota.CDOTAUserMsg_WorldLine{}, nil
	case 98:
		return &dota.CDOTAUserMsg_TournamentDrop{}, nil
	case 99:
		return &dota.CDOTAUserMsg_ItemAlert{}, nil
	case 100:
		return &dota.CDOTAUserMsg_HalloweenDrops{}, nil
	case 101:
		return &dota.CDOTAUserMsg_ChatWheel{}, nil
	case 102:
		return &dota.CDOTAUserMsg_ReceivedXmasGift{}, nil
	case 103:
		return &dota.CDOTAUserMsg_UpdateSharedContent{}, nil
	case 104:
		return &dota.CDOTAUserMsg_TutorialRequestExp{}, nil
	case 105:
		return &dota.CDOTAUserMsg_TutorialPingMinimap{}, nil
	case 107:
		return &dota.CDOTAUserMsg_ShowSurvey{}, nil
	case 108:
		return &dota.CDOTAUserMsg_TutorialFade{}, nil
	case 109:
		return &dota.CDOTAUserMsg_AddQuestLogEntry{}, nil
	case 110:
		return &dota.CDOTAUserMsg_SendStatPopup{}, nil
	case 111:
		return &dota.CDOTAUserMsg_TutorialFinish{}, nil
	case 112:
		return &dota.CDOTAUserMsg_SendRoshanPopup{}, nil
	case 113:
		return &dota.CDOTAUserMsg_SendGenericToolTip{}, nil
	}

	return nil, Error("Unknown BUMDUM")
}

/*
var NetSvc []ParserBaseEventMap
var BumDum []ParserBaseEventMap
var Maps = map[string]ParserBaseEventMap{}

func init() {
	NetSvc = []ParserBaseEventMap{
		{MapType: NET, Value: 0, Name: "net_NOP", EventType: NET_NOP, ItemType: reflect.TypeOf(dota.CNETMsg_NOP{})},
		{MapType: NET, Value: 1, Name: "net_Disconnect", EventType: NET_Disconnect, ItemType: reflect.TypeOf(dota.CNETMsg_Disconnect{})},
		{MapType: NET, Value: 2, Name: "net_File", EventType: NET_File, ItemType: reflect.TypeOf(dota.CNETMsg_File{})},
		{MapType: NET, Value: 3, Name: "net_SplitScreenUser", EventType: NET_SplitScreenUser, ItemType: reflect.TypeOf(dota.CNETMsg_SplitScreenUser{})},
		{MapType: NET, Value: 4, Name: "net_Tick", EventType: NET_Tick, ItemType: reflect.TypeOf(dota.CNETMsg_Tick{})},
		{MapType: NET, Value: 5, Name: "net_StringCmd", EventType: NET_StringCmd, ItemType: reflect.TypeOf(dota.CNETMsg_StringCmd{})},
		{MapType: NET, Value: 6, Name: "net_SetConVar", EventType: NET_SetConVar, ItemType: reflect.TypeOf(dota.CNETMsg_SetConVar{})},
		{MapType: NET, Value: 7, Name: "net_SignonState", EventType: NET_SignonState, ItemType: reflect.TypeOf(dota.CNETMsg_SignonState{})},
		{MapType: SVC, Value: 8, Name: "svc_ServerInfo", EventType: SVC_ServerInfo, ItemType: reflect.TypeOf(dota.CSVCMsg_ServerInfo{})},
		{MapType: SVC, Value: 9, Name: "svc_SendTable", EventType: SVC_SendTable, ItemType: reflect.TypeOf(dota.CSVCMsg_SendTable{})},
		{MapType: SVC, Value: 10, Name: "svc_ClassInfo", EventType: SVC_ClassInfo, ItemType: reflect.TypeOf(dota.CSVCMsg_ClassInfo{})},
		{MapType: SVC, Value: 11, Name: "svc_SetPause", EventType: SVC_SetPause, ItemType: reflect.TypeOf(dota.CSVCMsg_SetPause{})},
		{MapType: SVC, Value: 12, Name: "svc_CreateStringTable", EventType: SVC_CreateStringTable, ItemType: reflect.TypeOf(dota.CSVCMsg_CreateStringTable{})},
		{MapType: SVC, Value: 13, Name: "svc_UpdateStringTable", EventType: SVC_UpdateStringTable, ItemType: reflect.TypeOf(dota.CSVCMsg_UpdateStringTable{})},
		{MapType: SVC, Value: 14, Name: "svc_VoiceInit", EventType: SVC_VoiceInit, ItemType: reflect.TypeOf(dota.CSVCMsg_VoiceInit{})},
		{MapType: SVC, Value: 15, Name: "svc_VoiceData", EventType: SVC_VoiceData, ItemType: reflect.TypeOf(dota.CSVCMsg_VoiceData{})},
		{MapType: SVC, Value: 16, Name: "svc_Print", EventType: SVC_Print, ItemType: reflect.TypeOf(dota.CSVCMsg_Print{})},
		{MapType: SVC, Value: 17, Name: "svc_Sounds", EventType: SVC_Sounds, ItemType: reflect.TypeOf(dota.CSVCMsg_Sounds{})},
		{MapType: SVC, Value: 18, Name: "svc_SetView", EventType: SVC_SetView, ItemType: reflect.TypeOf(dota.CSVCMsg_SetView{})},
		{MapType: SVC, Value: 19, Name: "svc_FixAngle", EventType: SVC_FixAngle, ItemType: reflect.TypeOf(dota.CSVCMsg_FixAngle{})},
		{MapType: SVC, Value: 20, Name: "svc_CrosshairAngle", EventType: SVC_CrosshairAngle, ItemType: reflect.TypeOf(dota.CSVCMsg_CrosshairAngle{})},
		{MapType: SVC, Value: 21, Name: "svc_BSPDecal", EventType: SVC_BSPDecal, ItemType: reflect.TypeOf(dota.CSVCMsg_BSPDecal{})},
		{MapType: SVC, Value: 22, Name: "svc_SplitScreen", EventType: SVC_SplitScreen, ItemType: reflect.TypeOf(dota.CSVCMsg_SplitScreen{})},
		{MapType: SVC, Value: 23, Name: "svc_UserMessage", EventType: SVC_UserMessage, ItemType: reflect.TypeOf(dota.CSVCMsg_UserMessage{})},
		{MapType: SVC, Value: 25, Name: "svc_GameEvent", EventType: SVC_GameEvent, ItemType: reflect.TypeOf(dota.CSVCMsg_GameEvent{})},
		{MapType: SVC, Value: 26, Name: "svc_PacketEntities", EventType: SVC_PacketEntities, ItemType: reflect.TypeOf(dota.CSVCMsg_PacketEntities{})},
		{MapType: SVC, Value: 27, Name: "svc_TempEntities", EventType: SVC_TempEntities, ItemType: reflect.TypeOf(dota.CSVCMsg_TempEntities{})},
		{MapType: SVC, Value: 28, Name: "svc_Prefetch", EventType: SVC_Prefetch, ItemType: reflect.TypeOf(dota.CSVCMsg_Prefetch{})},
		{MapType: SVC, Value: 29, Name: "svc_Menu", EventType: SVC_Menu, ItemType: reflect.TypeOf(dota.CSVCMsg_Menu{})},
		{MapType: SVC, Value: 30, Name: "svc_GameEventList", EventType: SVC_GameEventList, ItemType: reflect.TypeOf(dota.CSVCMsg_GameEventList{})},
		{MapType: SVC, Value: 31, Name: "svc_GetCvarValue", EventType: SVC_GetCvarValue, ItemType: reflect.TypeOf(dota.CSVCMsg_GetCvarValue{})},
		{MapType: SVC, Value: 32, Name: "svc_PacketReliable", EventType: SVC_PacketReliable, ItemType: reflect.TypeOf(dota.CSVCMsg_PacketReliable{})},
	}
	}
	Maps = map[string]ParserBaseEventMap{
		"DEM_Stop":                         {MapType: DEM, Value: 0, Name: "DEM_Stop", EventType: DEM_Stop, ItemType: reflect.TypeOf(dota.CDemoStop{})},
		"DEM_FileHeader":                   {MapType: DEM, Value: 1, Name: "DEM_FileHeader", EventType: DEM_FileHeader, ItemType: reflect.TypeOf(dota.CDemoFileHeader{})},
		"DEM_FileInfo":                     {MapType: DEM, Value: 2, Name: "DEM_FileInfo", EventType: DEM_FileInfo, ItemType: reflect.TypeOf(dota.CDemoFileInfo{})},
		"DEM_SyncTick":                     {MapType: DEM, Value: 3, Name: "DEM_SyncTick", EventType: DEM_SyncTick, ItemType: reflect.TypeOf(dota.CDemoSyncTick{})},
		"DEM_SendTables":                   {MapType: DEM, Value: 4, Name: "DEM_SendTables", EventType: DEM_SendTables, ItemType: reflect.TypeOf(dota.CDemoSendTables{})},
		"DEM_ClassInfo":                    {MapType: DEM, Value: 5, Name: "DEM_ClassInfo", EventType: DEM_ClassInfo, ItemType: reflect.TypeOf(dota.CDemoClassInfo{})},
		"DEM_StringTables":                 {MapType: DEM, Value: 6, Name: "DEM_StringTables", EventType: DEM_StringTables, ItemType: reflect.TypeOf(dota.CDemoStringTables{})},
		"DEM_Packet":                       {MapType: DEM, Value: 7, Name: "DEM_Packet", EventType: DEM_Packet, ItemType: reflect.TypeOf(dota.CDemoPacket{})},
		"DEM_SignonPacket":                 {MapType: DEM, Name: "DEM_SignonPacket", EventType: DEM_SignonPacket},
		"DEM_ConsoleCmd":                   {MapType: DEM, Value: 9, Name: "DEM_ConsoleCmd", EventType: DEM_ConsoleCmd, ItemType: reflect.TypeOf(dota.CDemoConsoleCmd{})},
		"DEM_CustomData":                   {MapType: DEM, Value: 10, Name: "DEM_CustomData", EventType: DEM_CustomData, ItemType: reflect.TypeOf(dota.CDemoCustomData{})},
		"DEM_CustomDataCallbacks":          {MapType: DEM, Value: 11, Name: "DEM_CustomDataCallbacks", EventType: DEM_CustomDataCallbacks, ItemType: reflect.TypeOf(dota.CDemoCustomDataCallbacks{})},
		"DEM_UserCmd":                      {MapType: DEM, Value: 12, Name: "DEM_UserCmd", EventType: DEM_UserCmd, ItemType: reflect.TypeOf(dota.CDemoUserCmd{})},
		"DEM_FullPacket":                   {MapType: DEM, Value: 13, Name: "DEM_FullPacket", EventType: DEM_FullPacket, ItemType: reflect.TypeOf(dota.CDemoFullPacket{})},
		"net_NOP":                          {MapType: NET, Value: 0, Name: "net_NOP", EventType: NET_NOP, ItemType: reflect.TypeOf(dota.CNETMsg_NOP{})},
		"net_Disconnect":                   {MapType: NET, Value: 1, Name: "net_Disconnect", EventType: NET_Disconnect, ItemType: reflect.TypeOf(dota.CNETMsg_Disconnect{})},
		"net_File":                         {MapType: NET, Value: 2, Name: "net_File", EventType: NET_File, ItemType: reflect.TypeOf(dota.CNETMsg_File{})},
		"net_SplitScreenUser":              {MapType: NET, Value: 3, Name: "net_SplitScreenUser", EventType: NET_SplitScreenUser, ItemType: reflect.TypeOf(dota.CNETMsg_SplitScreenUser{})},
		"net_Tick":                         {MapType: NET, Value: 4, Name: "net_Tick", EventType: NET_Tick, ItemType: reflect.TypeOf(dota.CNETMsg_Tick{})},
		"net_StringCmd":                    {MapType: NET, Value: 5, Name: "net_StringCmd", EventType: NET_StringCmd, ItemType: reflect.TypeOf(dota.CNETMsg_StringCmd{})},
		"net_SetConVar":                    {MapType: NET, Value: 6, Name: "net_SetConVar", EventType: NET_SetConVar, ItemType: reflect.TypeOf(dota.CNETMsg_SetConVar{})},
		"net_SignonState":                  {MapType: NET, Value: 7, Name: "net_SignonState", EventType: NET_SignonState, ItemType: reflect.TypeOf(dota.CNETMsg_SignonState{})},
		"svc_ServerInfo":                   {MapType: SVC, Value: 8, Name: "svc_ServerInfo", EventType: SVC_ServerInfo, ItemType: reflect.TypeOf(dota.CSVCMsg_ServerInfo{})},
		"svc_SendTable":                    {MapType: SVC, Value: 9, Name: "svc_SendTable", EventType: SVC_SendTable, ItemType: reflect.TypeOf(dota.CSVCMsg_SendTable{})},
		"svc_ClassInfo":                    {MapType: SVC, Value: 10, Name: "svc_ClassInfo", EventType: SVC_ClassInfo, ItemType: reflect.TypeOf(dota.CSVCMsg_ClassInfo{})},
		"svc_SetPause":                     {MapType: SVC, Value: 11, Name: "svc_SetPause", EventType: SVC_SetPause, ItemType: reflect.TypeOf(dota.CSVCMsg_SetPause{})},
		"svc_CreateStringTable":            {MapType: SVC, Value: 12, Name: "svc_CreateStringTable", EventType: SVC_CreateStringTable, ItemType: reflect.TypeOf(dota.CSVCMsg_CreateStringTable{})},
		"svc_UpdateStringTable":            {MapType: SVC, Value: 13, Name: "svc_UpdateStringTable", EventType: SVC_UpdateStringTable, ItemType: reflect.TypeOf(dota.CSVCMsg_UpdateStringTable{})},
		"svc_VoiceInit":                    {MapType: SVC, Value: 14, Name: "svc_VoiceInit", EventType: SVC_VoiceInit, ItemType: reflect.TypeOf(dota.CSVCMsg_VoiceInit{})},
		"svc_VoiceData":                    {MapType: SVC, Value: 15, Name: "svc_VoiceData", EventType: SVC_VoiceData, ItemType: reflect.TypeOf(dota.CSVCMsg_VoiceData{})},
		"svc_Print":                        {MapType: SVC, Value: 16, Name: "svc_Print", EventType: SVC_Print, ItemType: reflect.TypeOf(dota.CSVCMsg_Print{})},
		"svc_Sounds":                       {MapType: SVC, Value: 17, Name: "svc_Sounds", EventType: SVC_Sounds, ItemType: reflect.TypeOf(dota.CSVCMsg_Sounds{})},
		"svc_SetView":                      {MapType: SVC, Value: 18, Name: "svc_SetView", EventType: SVC_SetView, ItemType: reflect.TypeOf(dota.CSVCMsg_SetView{})},
		"svc_FixAngle":                     {MapType: SVC, Value: 19, Name: "svc_FixAngle", EventType: SVC_FixAngle, ItemType: reflect.TypeOf(dota.CSVCMsg_FixAngle{})},
		"svc_CrosshairAngle":               {MapType: SVC, Value: 20, Name: "svc_CrosshairAngle", EventType: SVC_CrosshairAngle, ItemType: reflect.TypeOf(dota.CSVCMsg_CrosshairAngle{})},
		"svc_BSPDecal":                     {MapType: SVC, Value: 21, Name: "svc_BSPDecal", EventType: SVC_BSPDecal, ItemType: reflect.TypeOf(dota.CSVCMsg_BSPDecal{})},
		"svc_SplitScreen":                  {MapType: SVC, Value: 22, Name: "svc_SplitScreen", EventType: SVC_SplitScreen, ItemType: reflect.TypeOf(dota.CSVCMsg_SplitScreen{})},
		"svc_UserMessage":                  {MapType: SVC, Value: 23, Name: "svc_UserMessage", EventType: SVC_UserMessage, ItemType: reflect.TypeOf(dota.CSVCMsg_UserMessage{})},
		"svc_GameEvent":                    {MapType: SVC, Value: 25, Name: "svc_GameEvent", EventType: SVC_GameEvent, ItemType: reflect.TypeOf(dota.CSVCMsg_GameEvent{})},
		"svc_PacketEntities":               {MapType: SVC, Value: 26, Name: "svc_PacketEntities", EventType: SVC_PacketEntities, ItemType: reflect.TypeOf(dota.CSVCMsg_PacketEntities{})},
		"svc_TempEntities":                 {MapType: SVC, Value: 27, Name: "svc_TempEntities", EventType: SVC_TempEntities, ItemType: reflect.TypeOf(dota.CSVCMsg_TempEntities{})},
		"svc_Prefetch":                     {MapType: SVC, Value: 28, Name: "svc_Prefetch", EventType: SVC_Prefetch, ItemType: reflect.TypeOf(dota.CSVCMsg_Prefetch{})},
		"svc_Menu":                         {MapType: SVC, Value: 29, Name: "svc_Menu", EventType: SVC_Menu, ItemType: reflect.TypeOf(dota.CSVCMsg_Menu{})},
		"svc_GameEventList":                {MapType: SVC, Value: 30, Name: "svc_GameEventList", EventType: SVC_GameEventList, ItemType: reflect.TypeOf(dota.CSVCMsg_GameEventList{})},
		"svc_GetCvarValue":                 {MapType: SVC, Value: 31, Name: "svc_GetCvarValue", EventType: SVC_GetCvarValue, ItemType: reflect.TypeOf(dota.CSVCMsg_GetCvarValue{})},
		"svc_PacketReliable":               {MapType: SVC, Value: 32, Name: "svc_PacketReliable", EventType: SVC_PacketReliable, ItemType: reflect.TypeOf(dota.CSVCMsg_PacketReliable{})},
		"UM_AchievementEvent":              {MapType: BUM, Value: 1, Name: "UM_AchievementEvent", EventType: UM_AchievementEvent, ItemType: reflect.TypeOf(dota.CUserMsg_AchievementEvent{})},
		"UM_CloseCaption":                  {MapType: BUM, Value: 2, Name: "UM_CloseCaption", EventType: UM_CloseCaption, ItemType: reflect.TypeOf(dota.CUserMsg_CloseCaption{})},
		"UM_CurrentTimescale":              {MapType: BUM, Value: 4, Name: "UM_CurrentTimescale", EventType: UM_CurrentTimescale, ItemType: reflect.TypeOf(dota.CUserMsg_CurrentTimescale{})},
		"UM_DesiredTimescale":              {MapType: BUM, Value: 5, Name: "UM_DesiredTimescale", EventType: UM_DesiredTimescale, ItemType: reflect.TypeOf(dota.CUserMsg_DesiredTimescale{})},
		"UM_Fade":                          {MapType: BUM, Value: 6, Name: "UM_Fade", EventType: UM_Fade, ItemType: reflect.TypeOf(dota.CUserMsg_Fade{})},
		"UM_GameTitle":                     {MapType: BUM, Value: 7, Name: "UM_GameTitle", EventType: UM_GameTitle, ItemType: reflect.TypeOf(dota.CUserMsg_GameTitle{})},
		"UM_Geiger":                        {MapType: BUM, Value: 8, Name: "UM_Geiger", EventType: UM_Geiger, ItemType: reflect.TypeOf(dota.CUserMsg_Geiger{})},
		"UM_HintText":                      {MapType: BUM, Value: 9, Name: "UM_HintText", EventType: UM_HintText, ItemType: reflect.TypeOf(dota.CUserMsg_HintText{})},
		"UM_HudMsg":                        {MapType: BUM, Value: 10, Name: "UM_HudMsg", EventType: UM_HudMsg, ItemType: reflect.TypeOf(dota.CUserMsg_HudMsg{})},
		"UM_HudText":                       {MapType: BUM, Value: 11, Name: "UM_HudText", EventType: UM_HudText, ItemType: reflect.TypeOf(dota.CUserMsg_HudText{})},
		"UM_KeyHintText":                   {MapType: BUM, Value: 12, Name: "UM_KeyHintText", EventType: UM_KeyHintText, ItemType: reflect.TypeOf(dota.CUserMsg_KeyHintText{})},
		"UM_MessageText":                   {MapType: BUM, Value: 13, Name: "UM_MessageText", EventType: UM_MessageText, ItemType: reflect.TypeOf(dota.CUserMsg_MessageText{})},
		"UM_RequestState":                  {MapType: BUM, Value: 14, Name: "UM_RequestState", EventType: UM_RequestState, ItemType: reflect.TypeOf(dota.CUserMsg_RequestState{})},
		"UM_ResetHUD":                      {MapType: BUM, Value: 15, Name: "UM_ResetHUD", EventType: UM_ResetHUD, ItemType: reflect.TypeOf(dota.CUserMsg_ResetHUD{})},
		"UM_Rumble":                        {MapType: BUM, Value: 16, Name: "UM_Rumble", EventType: UM_Rumble, ItemType: reflect.TypeOf(dota.CUserMsg_Rumble{})},
		"UM_SayText":                       {MapType: BUM, Value: 17, Name: "UM_SayText", EventType: UM_SayText, ItemType: reflect.TypeOf(dota.CUserMsg_SayText{})},
		"UM_SayText2":                      {MapType: BUM, Value: 18, Name: "UM_SayText2", EventType: UM_SayText2, ItemType: reflect.TypeOf(dota.CUserMsg_SayText2{})},
		"UM_SayTextChannel":                {MapType: BUM, Value: 19, Name: "UM_SayTextChannel", EventType: UM_SayTextChannel, ItemType: reflect.TypeOf(dota.CUserMsg_SayTextChannel{})},
		"UM_Shake":                         {MapType: BUM, Value: 20, Name: "UM_Shake", EventType: UM_Shake, ItemType: reflect.TypeOf(dota.CUserMsg_Shake{})},
		"UM_ShakeDir":                      {MapType: BUM, Value: 21, Name: "UM_ShakeDir", EventType: UM_ShakeDir, ItemType: reflect.TypeOf(dota.CUserMsg_ShakeDir{})},
		"UM_StatsCrawlMsg":                 {MapType: BUM, Value: 22, Name: "UM_StatsCrawlMsg", EventType: UM_StatsCrawlMsg, ItemType: reflect.TypeOf(dota.CUserMsg_StatsCrawlMsg{})},
		"UM_StatsSkipState":                {MapType: BUM, Value: 23, Name: "UM_StatsSkipState", EventType: UM_StatsSkipState, ItemType: reflect.TypeOf(dota.CUserMsg_StatsSkipState{})},
		"UM_TextMsg":                       {MapType: BUM, Value: 24, Name: "UM_TextMsg", EventType: UM_TextMsg, ItemType: reflect.TypeOf(dota.CUserMsg_TextMsg{})},
		"UM_Tilt":                          {MapType: BUM, Value: 25, Name: "UM_Tilt", EventType: UM_Tilt, ItemType: reflect.TypeOf(dota.CUserMsg_Tilt{})},
		"UM_Train":                         {MapType: BUM, Value: 26, Name: "UM_Train", EventType: UM_Train, ItemType: reflect.TypeOf(dota.CUserMsg_Train{})},
		"UM_VGUIMenu":                      {MapType: BUM, Value: 27, Name: "UM_VGUIMenu", EventType: UM_VGUIMenu, ItemType: reflect.TypeOf(dota.CUserMsg_VGUIMenu{})},
		"UM_VoiceMask":                     {MapType: BUM, Value: 28, Name: "UM_VoiceMask", EventType: UM_VoiceMask, ItemType: reflect.TypeOf(dota.CUserMsg_VoiceMask{})},
		"UM_VoiceSubtitle":                 {MapType: BUM, Value: 29, Name: "UM_VoiceSubtitle", EventType: UM_VoiceSubtitle, ItemType: reflect.TypeOf(dota.CUserMsg_VoiceSubtitle{})},
		"UM_SendAudio":                     {MapType: BUM, Value: 30, Name: "UM_SendAudio", EventType: UM_SendAudio, ItemType: reflect.TypeOf(dota.CUserMsg_SendAudio{})},
		"DOTA_UM_AIDebugLine":              {MapType: DUM, Value: 65, Name: "DOTA_UM_AIDebugLine", EventType: DOTA_UM_AIDebugLine, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_AIDebugLine{})},
		"DOTA_UM_ChatEvent":                {MapType: DUM, Value: 66, Name: "DOTA_UM_ChatEvent", EventType: DOTA_UM_ChatEvent, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_ChatEvent{})},
		"DOTA_UM_CombatHeroPositions":      {MapType: DUM, Value: 67, Name: "DOTA_UM_CombatHeroPositions", EventType: DOTA_UM_CombatHeroPositions, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_CombatHeroPositions{})},
		"DOTA_UM_CombatLogData":            {MapType: DUM, Value: 68, Name: "DOTA_UM_CombatLogData", EventType: DOTA_UM_CombatLogData, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_CombatLogData{})},
		"DOTA_UM_CombatLogShowDeath":       {MapType: DUM, Value: 70, Name: "DOTA_UM_CombatLogShowDeath", EventType: DOTA_UM_CombatLogShowDeath, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_CombatLogShowDeath{})},
		"DOTA_UM_CreateLinearProjectile":   {MapType: DUM, Value: 71, Name: "DOTA_UM_CreateLinearProjectile", EventType: DOTA_UM_CreateLinearProjectile, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_CreateLinearProjectile{})},
		"DOTA_UM_DestroyLinearProjectile":  {MapType: DUM, Value: 72, Name: "DOTA_UM_DestroyLinearProjectile", EventType: DOTA_UM_DestroyLinearProjectile, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_DestroyLinearProjectile{})},
		"DOTA_UM_DodgeTrackingProjectiles": {MapType: DUM, Value: 73, Name: "DOTA_UM_DodgeTrackingProjectiles", EventType: DOTA_UM_DodgeTrackingProjectiles, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_DodgeTrackingProjectiles{})},
		"DOTA_UM_GlobalLightColor":         {MapType: DUM, Value: 74, Name: "DOTA_UM_GlobalLightColor", EventType: DOTA_UM_GlobalLightColor, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_GlobalLightColor{})},
		"DOTA_UM_GlobalLightDirection":     {MapType: DUM, Value: 75, Name: "DOTA_UM_GlobalLightDirection", EventType: DOTA_UM_GlobalLightDirection, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_GlobalLightDirection{})},
		"DOTA_UM_InvalidCommand":           {MapType: DUM, Value: 76, Name: "DOTA_UM_InvalidCommand", EventType: DOTA_UM_InvalidCommand, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_InvalidCommand{})},
		"DOTA_UM_LocationPing":             {MapType: DUM, Value: 77, Name: "DOTA_UM_LocationPing", EventType: DOTA_UM_LocationPing, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_LocationPing{})},
		"DOTA_UM_MapLine":                  {MapType: DUM, Value: 78, Name: "DOTA_UM_MapLine", EventType: DOTA_UM_MapLine, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_MapLine{})},
		"DOTA_UM_MiniKillCamInfo":          {MapType: DUM, Value: 79, Name: "DOTA_UM_MiniKillCamInfo", EventType: DOTA_UM_MiniKillCamInfo, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_MiniKillCamInfo{})},
		"DOTA_UM_MinimapDebugPoint":        {MapType: DUM, Value: 80, Name: "DOTA_UM_MinimapDebugPoint", EventType: DOTA_UM_MinimapDebugPoint, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_MinimapDebugPoint{})},
		"DOTA_UM_MinimapEvent":             {MapType: DUM, Value: 81, Name: "DOTA_UM_MinimapEvent", EventType: DOTA_UM_MinimapEvent, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_MinimapEvent{})},
		"DOTA_UM_NevermoreRequiem":         {MapType: DUM, Value: 82, Name: "DOTA_UM_NevermoreRequiem", EventType: DOTA_UM_NevermoreRequiem, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_NevermoreRequiem{})},
		"DOTA_UM_OverheadEvent":            {MapType: DUM, Value: 83, Name: "DOTA_UM_OverheadEvent", EventType: DOTA_UM_OverheadEvent, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_OverheadEvent{})},
		"DOTA_UM_SetNextAutobuyItem":       {MapType: DUM, Value: 84, Name: "DOTA_UM_SetNextAutobuyItem", EventType: DOTA_UM_SetNextAutobuyItem, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_SetNextAutobuyItem{})},
		"DOTA_UM_SharedCooldown":           {MapType: DUM, Value: 85, Name: "DOTA_UM_SharedCooldown", EventType: DOTA_UM_SharedCooldown, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_SharedCooldown{})},
		"DOTA_UM_SpectatorPlayerClick":     {MapType: DUM, Value: 86, Name: "DOTA_UM_SpectatorPlayerClick", EventType: DOTA_UM_SpectatorPlayerClick, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_SpectatorPlayerClick{})},
		"DOTA_UM_TutorialTipInfo":          {MapType: DUM, Value: 87, Name: "DOTA_UM_TutorialTipInfo", EventType: DOTA_UM_TutorialTipInfo, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_TutorialTipInfo{})},
		"DOTA_UM_UnitEvent":                {MapType: DUM, Value: 88, Name: "DOTA_UM_UnitEvent", EventType: DOTA_UM_UnitEvent, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_UnitEvent{})},
		"DOTA_UM_ParticleManager":          {MapType: DUM, Value: 89, Name: "DOTA_UM_ParticleManager", EventType: DOTA_UM_ParticleManager, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_ParticleManager{})},
		"DOTA_UM_BotChat":                  {MapType: DUM, Value: 90, Name: "DOTA_UM_BotChat", EventType: DOTA_UM_BotChat, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_BotChat{})},
		"DOTA_UM_HudError":                 {MapType: DUM, Value: 91, Name: "DOTA_UM_HudError", EventType: DOTA_UM_HudError, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_HudError{})},
		"DOTA_UM_ItemPurchased":            {MapType: DUM, Value: 92, Name: "DOTA_UM_ItemPurchased", EventType: DOTA_UM_ItemPurchased, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_ItemPurchased{})},
		"DOTA_UM_Ping":                     {MapType: DUM, Value: 93, Name: "DOTA_UM_Ping", EventType: DOTA_UM_Ping, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_Ping{})},
		"DOTA_UM_ItemFound":                {MapType: DUM, Value: 94, Name: "DOTA_UM_ItemFound", EventType: DOTA_UM_ItemFound, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_ItemFound{})},
		"DOTA_UM_SwapVerify":               {MapType: DUM, Value: 96, Name: "DOTA_UM_SwapVerify", EventType: DOTA_UM_SwapVerify, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_SwapVerify{})},
		"DOTA_UM_WorldLine":                {MapType: DUM, Value: 97, Name: "DOTA_UM_WorldLine", EventType: DOTA_UM_WorldLine, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_WorldLine{})},
		"DOTA_UM_TournamentDrop":           {MapType: DUM, Value: 98, Name: "DOTA_UM_TournamentDrop", EventType: DOTA_UM_TournamentDrop, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_TournamentDrop{})},
		"DOTA_UM_ItemAlert":                {MapType: DUM, Value: 99, Name: "DOTA_UM_ItemAlert", EventType: DOTA_UM_ItemAlert, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_ItemAlert{})},
		"DOTA_UM_HalloweenDrops":           {MapType: DUM, Value: 100, Name: "DOTA_UM_HalloweenDrops", EventType: DOTA_UM_HalloweenDrops, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_HalloweenDrops{})},
		"DOTA_UM_ChatWheel":                {MapType: DUM, Value: 101, Name: "DOTA_UM_ChatWheel", EventType: DOTA_UM_ChatWheel, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_ChatWheel{})},
		"DOTA_UM_ReceivedXmasGift":         {MapType: DUM, Value: 102, Name: "DOTA_UM_ReceivedXmasGift", EventType: DOTA_UM_ReceivedXmasGift, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_ReceivedXmasGift{})},
		"DOTA_UM_UpdateSharedContent":      {MapType: DUM, Value: 103, Name: "DOTA_UM_UpdateSharedContent", EventType: DOTA_UM_UpdateSharedContent, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_UpdateSharedContent{})},
		"DOTA_UM_TutorialRequestExp":       {MapType: DUM, Value: 104, Name: "DOTA_UM_TutorialRequestExp", EventType: DOTA_UM_TutorialRequestExp, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_TutorialRequestExp{})},
		"DOTA_UM_TutorialPingMinimap":      {MapType: DUM, Value: 105, Name: "DOTA_UM_TutorialPingMinimap", EventType: DOTA_UM_TutorialPingMinimap, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_TutorialPingMinimap{})},
		"DOTA_UM_ShowSurvey":               {MapType: DUM, Value: 107, Name: "DOTA_UM_ShowSurvey", EventType: DOTA_UM_ShowSurvey, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_ShowSurvey{})},
		"DOTA_UM_TutorialFade":             {MapType: DUM, Value: 108, Name: "DOTA_UM_TutorialFade", EventType: DOTA_UM_TutorialFade, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_TutorialFade{})},
		"DOTA_UM_AddQuestLogEntry":         {MapType: DUM, Value: 109, Name: "DOTA_UM_AddQuestLogEntry", EventType: DOTA_UM_AddQuestLogEntry, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_AddQuestLogEntry{})},
		"DOTA_UM_SendStatPopup":            {MapType: DUM, Value: 110, Name: "DOTA_UM_SendStatPopup", EventType: DOTA_UM_SendStatPopup, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_SendStatPopup{})},
		"DOTA_UM_TutorialFinish":           {MapType: DUM, Value: 111, Name: "DOTA_UM_TutorialFinish", EventType: DOTA_UM_TutorialFinish, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_TutorialFinish{})},
		"DOTA_UM_SendRoshanPopup":          {MapType: DUM, Value: 112, Name: "DOTA_UM_SendRoshanPopup", EventType: DOTA_UM_SendRoshanPopup, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_SendRoshanPopup{})},
		"DOTA_UM_SendGenericToolTip":       {MapType: DUM, Value: 113, Name: "DOTA_UM_SendGenericToolTip", EventType: DOTA_UM_SendGenericToolTip, ItemType: reflect.TypeOf(dota.CDOTAUserMsg_SendGenericToolTip{})},
	}
}
*/
