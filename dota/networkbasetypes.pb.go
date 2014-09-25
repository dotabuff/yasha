// Code generated by protoc-gen-gogo.
// source: networkbasetypes.proto
// DO NOT EDIT!

package dota

import proto "code.google.com/p/gogoprotobuf/proto"
import json "encoding/json"
import math "math"

// discarding unused import google_protobuf "github.com/dotabuff/yasha/dota/google/protobuf/descriptor.pb"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type SIGNONSTATE int32

const (
	SIGNONSTATE_SIGNONSTATE_NONE        SIGNONSTATE = 0
	SIGNONSTATE_SIGNONSTATE_CHALLENGE   SIGNONSTATE = 1
	SIGNONSTATE_SIGNONSTATE_CONNECTED   SIGNONSTATE = 2
	SIGNONSTATE_SIGNONSTATE_NEW         SIGNONSTATE = 3
	SIGNONSTATE_SIGNONSTATE_PRESPAWN    SIGNONSTATE = 4
	SIGNONSTATE_SIGNONSTATE_SPAWN       SIGNONSTATE = 5
	SIGNONSTATE_SIGNONSTATE_FULL        SIGNONSTATE = 6
	SIGNONSTATE_SIGNONSTATE_CHANGELEVEL SIGNONSTATE = 7
)

var SIGNONSTATE_name = map[int32]string{
	0: "SIGNONSTATE_NONE",
	1: "SIGNONSTATE_CHALLENGE",
	2: "SIGNONSTATE_CONNECTED",
	3: "SIGNONSTATE_NEW",
	4: "SIGNONSTATE_PRESPAWN",
	5: "SIGNONSTATE_SPAWN",
	6: "SIGNONSTATE_FULL",
	7: "SIGNONSTATE_CHANGELEVEL",
}
var SIGNONSTATE_value = map[string]int32{
	"SIGNONSTATE_NONE":        0,
	"SIGNONSTATE_CHALLENGE":   1,
	"SIGNONSTATE_CONNECTED":   2,
	"SIGNONSTATE_NEW":         3,
	"SIGNONSTATE_PRESPAWN":    4,
	"SIGNONSTATE_SPAWN":       5,
	"SIGNONSTATE_FULL":        6,
	"SIGNONSTATE_CHANGELEVEL": 7,
}

func (x SIGNONSTATE) Enum() *SIGNONSTATE {
	p := new(SIGNONSTATE)
	*p = x
	return p
}
func (x SIGNONSTATE) String() string {
	return proto.EnumName(SIGNONSTATE_name, int32(x))
}
func (x *SIGNONSTATE) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(SIGNONSTATE_value, data, "SIGNONSTATE")
	if err != nil {
		return err
	}
	*x = SIGNONSTATE(value)
	return nil
}

type NET_Messages int32

const (
	NET_Messages_net_NOP             NET_Messages = 0
	NET_Messages_net_Disconnect      NET_Messages = 1
	NET_Messages_net_File            NET_Messages = 2
	NET_Messages_net_SplitScreenUser NET_Messages = 3
	NET_Messages_net_Tick            NET_Messages = 4
	NET_Messages_net_StringCmd       NET_Messages = 5
	NET_Messages_net_SetConVar       NET_Messages = 6
	NET_Messages_net_SignonState     NET_Messages = 7
)

var NET_Messages_name = map[int32]string{
	0: "net_NOP",
	1: "net_Disconnect",
	2: "net_File",
	3: "net_SplitScreenUser",
	4: "net_Tick",
	5: "net_StringCmd",
	6: "net_SetConVar",
	7: "net_SignonState",
}
var NET_Messages_value = map[string]int32{
	"net_NOP":             0,
	"net_Disconnect":      1,
	"net_File":            2,
	"net_SplitScreenUser": 3,
	"net_Tick":            4,
	"net_StringCmd":       5,
	"net_SetConVar":       6,
	"net_SignonState":     7,
}

func (x NET_Messages) Enum() *NET_Messages {
	p := new(NET_Messages)
	*p = x
	return p
}
func (x NET_Messages) String() string {
	return proto.EnumName(NET_Messages_name, int32(x))
}
func (x *NET_Messages) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(NET_Messages_value, data, "NET_Messages")
	if err != nil {
		return err
	}
	*x = NET_Messages(value)
	return nil
}

type CMsgVector struct {
	X                *float32 `protobuf:"fixed32,1,opt,name=x" json:"x,omitempty"`
	Y                *float32 `protobuf:"fixed32,2,opt,name=y" json:"y,omitempty"`
	Z                *float32 `protobuf:"fixed32,3,opt,name=z" json:"z,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *CMsgVector) Reset()         { *m = CMsgVector{} }
func (m *CMsgVector) String() string { return proto.CompactTextString(m) }
func (*CMsgVector) ProtoMessage()    {}

func (m *CMsgVector) GetX() float32 {
	if m != nil && m.X != nil {
		return *m.X
	}
	return 0
}

func (m *CMsgVector) GetY() float32 {
	if m != nil && m.Y != nil {
		return *m.Y
	}
	return 0
}

func (m *CMsgVector) GetZ() float32 {
	if m != nil && m.Z != nil {
		return *m.Z
	}
	return 0
}

type CMsgVector2D struct {
	X                *float32 `protobuf:"fixed32,1,opt,name=x" json:"x,omitempty"`
	Y                *float32 `protobuf:"fixed32,2,opt,name=y" json:"y,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *CMsgVector2D) Reset()         { *m = CMsgVector2D{} }
func (m *CMsgVector2D) String() string { return proto.CompactTextString(m) }
func (*CMsgVector2D) ProtoMessage()    {}

func (m *CMsgVector2D) GetX() float32 {
	if m != nil && m.X != nil {
		return *m.X
	}
	return 0
}

func (m *CMsgVector2D) GetY() float32 {
	if m != nil && m.Y != nil {
		return *m.Y
	}
	return 0
}

type CMsgQAngle struct {
	X                *float32 `protobuf:"fixed32,1,opt,name=x" json:"x,omitempty"`
	Y                *float32 `protobuf:"fixed32,2,opt,name=y" json:"y,omitempty"`
	Z                *float32 `protobuf:"fixed32,3,opt,name=z" json:"z,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *CMsgQAngle) Reset()         { *m = CMsgQAngle{} }
func (m *CMsgQAngle) String() string { return proto.CompactTextString(m) }
func (*CMsgQAngle) ProtoMessage()    {}

func (m *CMsgQAngle) GetX() float32 {
	if m != nil && m.X != nil {
		return *m.X
	}
	return 0
}

func (m *CMsgQAngle) GetY() float32 {
	if m != nil && m.Y != nil {
		return *m.Y
	}
	return 0
}

func (m *CMsgQAngle) GetZ() float32 {
	if m != nil && m.Z != nil {
		return *m.Z
	}
	return 0
}

type CMsg_CVars struct {
	Cvars            []*CMsg_CVars_CVar `protobuf:"bytes,1,rep,name=cvars" json:"cvars,omitempty"`
	XXX_unrecognized []byte             `json:"-"`
}

func (m *CMsg_CVars) Reset()         { *m = CMsg_CVars{} }
func (m *CMsg_CVars) String() string { return proto.CompactTextString(m) }
func (*CMsg_CVars) ProtoMessage()    {}

func (m *CMsg_CVars) GetCvars() []*CMsg_CVars_CVar {
	if m != nil {
		return m.Cvars
	}
	return nil
}

type CMsg_CVars_CVar struct {
	Name             *string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Value            *string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *CMsg_CVars_CVar) Reset()         { *m = CMsg_CVars_CVar{} }
func (m *CMsg_CVars_CVar) String() string { return proto.CompactTextString(m) }
func (*CMsg_CVars_CVar) ProtoMessage()    {}

func (m *CMsg_CVars_CVar) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *CMsg_CVars_CVar) GetValue() string {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return ""
}

type CNETMsg_NOP struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *CNETMsg_NOP) Reset()         { *m = CNETMsg_NOP{} }
func (m *CNETMsg_NOP) String() string { return proto.CompactTextString(m) }
func (*CNETMsg_NOP) ProtoMessage()    {}

type CNETMsg_SplitScreenUser struct {
	Slot             *int32 `protobuf:"varint,1,opt,name=slot" json:"slot,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *CNETMsg_SplitScreenUser) Reset()         { *m = CNETMsg_SplitScreenUser{} }
func (m *CNETMsg_SplitScreenUser) String() string { return proto.CompactTextString(m) }
func (*CNETMsg_SplitScreenUser) ProtoMessage()    {}

func (m *CNETMsg_SplitScreenUser) GetSlot() int32 {
	if m != nil && m.Slot != nil {
		return *m.Slot
	}
	return 0
}

type CNETMsg_Disconnect struct {
	Reason           *ENetworkDisconnectionReason `protobuf:"varint,2,opt,name=reason,enum=dota.ENetworkDisconnectionReason,def=0" json:"reason,omitempty"`
	XXX_unrecognized []byte                       `json:"-"`
}

func (m *CNETMsg_Disconnect) Reset()         { *m = CNETMsg_Disconnect{} }
func (m *CNETMsg_Disconnect) String() string { return proto.CompactTextString(m) }
func (*CNETMsg_Disconnect) ProtoMessage()    {}

const Default_CNETMsg_Disconnect_Reason ENetworkDisconnectionReason = ENetworkDisconnectionReason_NETWORK_DISCONNECT_INVALID

func (m *CNETMsg_Disconnect) GetReason() ENetworkDisconnectionReason {
	if m != nil && m.Reason != nil {
		return *m.Reason
	}
	return Default_CNETMsg_Disconnect_Reason
}

type CNETMsg_File struct {
	TransferId       *int32  `protobuf:"varint,1,opt,name=transfer_id" json:"transfer_id,omitempty"`
	FileName         *string `protobuf:"bytes,2,opt,name=file_name" json:"file_name,omitempty"`
	IsReplayDemoFile *bool   `protobuf:"varint,3,opt,name=is_replay_demo_file" json:"is_replay_demo_file,omitempty"`
	Deny             *bool   `protobuf:"varint,4,opt,name=deny" json:"deny,omitempty"`
	IsFileRequested  *bool   `protobuf:"varint,5,opt,name=is_file_requested" json:"is_file_requested,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *CNETMsg_File) Reset()         { *m = CNETMsg_File{} }
func (m *CNETMsg_File) String() string { return proto.CompactTextString(m) }
func (*CNETMsg_File) ProtoMessage()    {}

func (m *CNETMsg_File) GetTransferId() int32 {
	if m != nil && m.TransferId != nil {
		return *m.TransferId
	}
	return 0
}

func (m *CNETMsg_File) GetFileName() string {
	if m != nil && m.FileName != nil {
		return *m.FileName
	}
	return ""
}

func (m *CNETMsg_File) GetIsReplayDemoFile() bool {
	if m != nil && m.IsReplayDemoFile != nil {
		return *m.IsReplayDemoFile
	}
	return false
}

func (m *CNETMsg_File) GetDeny() bool {
	if m != nil && m.Deny != nil {
		return *m.Deny
	}
	return false
}

func (m *CNETMsg_File) GetIsFileRequested() bool {
	if m != nil && m.IsFileRequested != nil {
		return *m.IsFileRequested
	}
	return false
}

type CNETMsg_Tick struct {
	Tick                            *uint32 `protobuf:"varint,1,opt,name=tick" json:"tick,omitempty"`
	HostComputationtime             *uint32 `protobuf:"varint,4,opt,name=host_computationtime" json:"host_computationtime,omitempty"`
	HostComputationtimeStdDeviation *uint32 `protobuf:"varint,5,opt,name=host_computationtime_std_deviation" json:"host_computationtime_std_deviation,omitempty"`
	HostFramestarttimeStdDeviation  *uint32 `protobuf:"varint,6,opt,name=host_framestarttime_std_deviation" json:"host_framestarttime_std_deviation,omitempty"`
	XXX_unrecognized                []byte  `json:"-"`
}

func (m *CNETMsg_Tick) Reset()         { *m = CNETMsg_Tick{} }
func (m *CNETMsg_Tick) String() string { return proto.CompactTextString(m) }
func (*CNETMsg_Tick) ProtoMessage()    {}

func (m *CNETMsg_Tick) GetTick() uint32 {
	if m != nil && m.Tick != nil {
		return *m.Tick
	}
	return 0
}

func (m *CNETMsg_Tick) GetHostComputationtime() uint32 {
	if m != nil && m.HostComputationtime != nil {
		return *m.HostComputationtime
	}
	return 0
}

func (m *CNETMsg_Tick) GetHostComputationtimeStdDeviation() uint32 {
	if m != nil && m.HostComputationtimeStdDeviation != nil {
		return *m.HostComputationtimeStdDeviation
	}
	return 0
}

func (m *CNETMsg_Tick) GetHostFramestarttimeStdDeviation() uint32 {
	if m != nil && m.HostFramestarttimeStdDeviation != nil {
		return *m.HostFramestarttimeStdDeviation
	}
	return 0
}

type CNETMsg_StringCmd struct {
	Command          *string `protobuf:"bytes,1,opt,name=command" json:"command,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *CNETMsg_StringCmd) Reset()         { *m = CNETMsg_StringCmd{} }
func (m *CNETMsg_StringCmd) String() string { return proto.CompactTextString(m) }
func (*CNETMsg_StringCmd) ProtoMessage()    {}

func (m *CNETMsg_StringCmd) GetCommand() string {
	if m != nil && m.Command != nil {
		return *m.Command
	}
	return ""
}

type CNETMsg_SetConVar struct {
	Convars          *CMsg_CVars `protobuf:"bytes,1,opt,name=convars" json:"convars,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *CNETMsg_SetConVar) Reset()         { *m = CNETMsg_SetConVar{} }
func (m *CNETMsg_SetConVar) String() string { return proto.CompactTextString(m) }
func (*CNETMsg_SetConVar) ProtoMessage()    {}

func (m *CNETMsg_SetConVar) GetConvars() *CMsg_CVars {
	if m != nil {
		return m.Convars
	}
	return nil
}

type CNETMsg_SignonState struct {
	SignonState       *uint32  `protobuf:"varint,1,opt,name=signon_state" json:"signon_state,omitempty"`
	SpawnCount        *uint32  `protobuf:"varint,2,opt,name=spawn_count" json:"spawn_count,omitempty"`
	NumServerPlayers  *uint32  `protobuf:"varint,3,opt,name=num_server_players" json:"num_server_players,omitempty"`
	PlayersNetworkids []string `protobuf:"bytes,4,rep,name=players_networkids" json:"players_networkids,omitempty"`
	MapName           *string  `protobuf:"bytes,5,opt,name=map_name" json:"map_name,omitempty"`
	XXX_unrecognized  []byte   `json:"-"`
}

func (m *CNETMsg_SignonState) Reset()         { *m = CNETMsg_SignonState{} }
func (m *CNETMsg_SignonState) String() string { return proto.CompactTextString(m) }
func (*CNETMsg_SignonState) ProtoMessage()    {}

func (m *CNETMsg_SignonState) GetSignonState() uint32 {
	if m != nil && m.SignonState != nil {
		return *m.SignonState
	}
	return 0
}

func (m *CNETMsg_SignonState) GetSpawnCount() uint32 {
	if m != nil && m.SpawnCount != nil {
		return *m.SpawnCount
	}
	return 0
}

func (m *CNETMsg_SignonState) GetNumServerPlayers() uint32 {
	if m != nil && m.NumServerPlayers != nil {
		return *m.NumServerPlayers
	}
	return 0
}

func (m *CNETMsg_SignonState) GetPlayersNetworkids() []string {
	if m != nil {
		return m.PlayersNetworkids
	}
	return nil
}

func (m *CNETMsg_SignonState) GetMapName() string {
	if m != nil && m.MapName != nil {
		return *m.MapName
	}
	return ""
}

type CSVCMsg_GameEvent struct {
	EventName        *string                  `protobuf:"bytes,1,opt,name=event_name" json:"event_name,omitempty"`
	Eventid          *int32                   `protobuf:"varint,2,opt,name=eventid" json:"eventid,omitempty"`
	Keys             []*CSVCMsg_GameEventKeyT `protobuf:"bytes,3,rep,name=keys" json:"keys,omitempty"`
	XXX_unrecognized []byte                   `json:"-"`
}

func (m *CSVCMsg_GameEvent) Reset()         { *m = CSVCMsg_GameEvent{} }
func (m *CSVCMsg_GameEvent) String() string { return proto.CompactTextString(m) }
func (*CSVCMsg_GameEvent) ProtoMessage()    {}

func (m *CSVCMsg_GameEvent) GetEventName() string {
	if m != nil && m.EventName != nil {
		return *m.EventName
	}
	return ""
}

func (m *CSVCMsg_GameEvent) GetEventid() int32 {
	if m != nil && m.Eventid != nil {
		return *m.Eventid
	}
	return 0
}

func (m *CSVCMsg_GameEvent) GetKeys() []*CSVCMsg_GameEventKeyT {
	if m != nil {
		return m.Keys
	}
	return nil
}

type CSVCMsg_GameEventKeyT struct {
	Type             *int32   `protobuf:"varint,1,opt,name=type" json:"type,omitempty"`
	ValString        *string  `protobuf:"bytes,2,opt,name=val_string" json:"val_string,omitempty"`
	ValFloat         *float32 `protobuf:"fixed32,3,opt,name=val_float" json:"val_float,omitempty"`
	ValLong          *int32   `protobuf:"varint,4,opt,name=val_long" json:"val_long,omitempty"`
	ValShort         *int32   `protobuf:"varint,5,opt,name=val_short" json:"val_short,omitempty"`
	ValByte          *int32   `protobuf:"varint,6,opt,name=val_byte" json:"val_byte,omitempty"`
	ValBool          *bool    `protobuf:"varint,7,opt,name=val_bool" json:"val_bool,omitempty"`
	ValUint64        *uint64  `protobuf:"varint,8,opt,name=val_uint64" json:"val_uint64,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *CSVCMsg_GameEventKeyT) Reset()         { *m = CSVCMsg_GameEventKeyT{} }
func (m *CSVCMsg_GameEventKeyT) String() string { return proto.CompactTextString(m) }
func (*CSVCMsg_GameEventKeyT) ProtoMessage()    {}

func (m *CSVCMsg_GameEventKeyT) GetType() int32 {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return 0
}

func (m *CSVCMsg_GameEventKeyT) GetValString() string {
	if m != nil && m.ValString != nil {
		return *m.ValString
	}
	return ""
}

func (m *CSVCMsg_GameEventKeyT) GetValFloat() float32 {
	if m != nil && m.ValFloat != nil {
		return *m.ValFloat
	}
	return 0
}

func (m *CSVCMsg_GameEventKeyT) GetValLong() int32 {
	if m != nil && m.ValLong != nil {
		return *m.ValLong
	}
	return 0
}

func (m *CSVCMsg_GameEventKeyT) GetValShort() int32 {
	if m != nil && m.ValShort != nil {
		return *m.ValShort
	}
	return 0
}

func (m *CSVCMsg_GameEventKeyT) GetValByte() int32 {
	if m != nil && m.ValByte != nil {
		return *m.ValByte
	}
	return 0
}

func (m *CSVCMsg_GameEventKeyT) GetValBool() bool {
	if m != nil && m.ValBool != nil {
		return *m.ValBool
	}
	return false
}

func (m *CSVCMsg_GameEventKeyT) GetValUint64() uint64 {
	if m != nil && m.ValUint64 != nil {
		return *m.ValUint64
	}
	return 0
}

type CSVCMsgList_GameEvents struct {
	Events           []*CSVCMsgList_GameEventsEventT `protobuf:"bytes,1,rep,name=events" json:"events,omitempty"`
	XXX_unrecognized []byte                          `json:"-"`
}

func (m *CSVCMsgList_GameEvents) Reset()         { *m = CSVCMsgList_GameEvents{} }
func (m *CSVCMsgList_GameEvents) String() string { return proto.CompactTextString(m) }
func (*CSVCMsgList_GameEvents) ProtoMessage()    {}

func (m *CSVCMsgList_GameEvents) GetEvents() []*CSVCMsgList_GameEventsEventT {
	if m != nil {
		return m.Events
	}
	return nil
}

type CSVCMsgList_GameEventsEventT struct {
	Tick             *int32             `protobuf:"varint,1,opt,name=tick" json:"tick,omitempty"`
	Event            *CSVCMsg_GameEvent `protobuf:"bytes,2,opt,name=event" json:"event,omitempty"`
	XXX_unrecognized []byte             `json:"-"`
}

func (m *CSVCMsgList_GameEventsEventT) Reset()         { *m = CSVCMsgList_GameEventsEventT{} }
func (m *CSVCMsgList_GameEventsEventT) String() string { return proto.CompactTextString(m) }
func (*CSVCMsgList_GameEventsEventT) ProtoMessage()    {}

func (m *CSVCMsgList_GameEventsEventT) GetTick() int32 {
	if m != nil && m.Tick != nil {
		return *m.Tick
	}
	return 0
}

func (m *CSVCMsgList_GameEventsEventT) GetEvent() *CSVCMsg_GameEvent {
	if m != nil {
		return m.Event
	}
	return nil
}

type CSVCMsg_UserMessage struct {
	MsgType          *int32 `protobuf:"varint,1,opt,name=msg_type" json:"msg_type,omitempty"`
	MsgData          []byte `protobuf:"bytes,2,opt,name=msg_data" json:"msg_data,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *CSVCMsg_UserMessage) Reset()         { *m = CSVCMsg_UserMessage{} }
func (m *CSVCMsg_UserMessage) String() string { return proto.CompactTextString(m) }
func (*CSVCMsg_UserMessage) ProtoMessage()    {}

func (m *CSVCMsg_UserMessage) GetMsgType() int32 {
	if m != nil && m.MsgType != nil {
		return *m.MsgType
	}
	return 0
}

func (m *CSVCMsg_UserMessage) GetMsgData() []byte {
	if m != nil {
		return m.MsgData
	}
	return nil
}

type CSVCMsgList_UserMessages struct {
	Usermsgs         []*CSVCMsgList_UserMessagesUsermsgT `protobuf:"bytes,1,rep,name=usermsgs" json:"usermsgs,omitempty"`
	XXX_unrecognized []byte                              `json:"-"`
}

func (m *CSVCMsgList_UserMessages) Reset()         { *m = CSVCMsgList_UserMessages{} }
func (m *CSVCMsgList_UserMessages) String() string { return proto.CompactTextString(m) }
func (*CSVCMsgList_UserMessages) ProtoMessage()    {}

func (m *CSVCMsgList_UserMessages) GetUsermsgs() []*CSVCMsgList_UserMessagesUsermsgT {
	if m != nil {
		return m.Usermsgs
	}
	return nil
}

type CSVCMsgList_UserMessagesUsermsgT struct {
	Tick             *int32               `protobuf:"varint,1,opt,name=tick" json:"tick,omitempty"`
	Msg              *CSVCMsg_UserMessage `protobuf:"bytes,2,opt,name=msg" json:"msg,omitempty"`
	XXX_unrecognized []byte               `json:"-"`
}

func (m *CSVCMsgList_UserMessagesUsermsgT) Reset()         { *m = CSVCMsgList_UserMessagesUsermsgT{} }
func (m *CSVCMsgList_UserMessagesUsermsgT) String() string { return proto.CompactTextString(m) }
func (*CSVCMsgList_UserMessagesUsermsgT) ProtoMessage()    {}

func (m *CSVCMsgList_UserMessagesUsermsgT) GetTick() int32 {
	if m != nil && m.Tick != nil {
		return *m.Tick
	}
	return 0
}

func (m *CSVCMsgList_UserMessagesUsermsgT) GetMsg() *CSVCMsg_UserMessage {
	if m != nil {
		return m.Msg
	}
	return nil
}

func init() {
	proto.RegisterEnum("dota.SIGNONSTATE", SIGNONSTATE_name, SIGNONSTATE_value)
	proto.RegisterEnum("dota.NET_Messages", NET_Messages_name, NET_Messages_value)
}
