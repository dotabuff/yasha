package string_tables

import (
	dota "github.com/dotabuff/d2rp/dota"
)

type StringTable struct {
	Tick  int
	Index int
	Name  string
	Items map[int]*StringTableItem
}

type StringTableItem struct {
	Str          string
	Data         []byte
	ModifierBuff *dota.CDOTAModifierBuffTableEntry
	Userinfo     *Userinfo
}

const (
	MAX_PLAYER_NAME_LENGTH = 32
	MAX_CUSTOM_FILES       = 4  // max 4 files
	SIGNED_GUID_LEN        = 32 // Hashed CD Key (32 hex alphabetic chars + 0 terminator )
)

type rawUserinfo struct {
	Xuid        uint64
	Name        [MAX_PLAYER_NAME_LENGTH]byte
	UserID      int32
	Guid        [SIGNED_GUID_LEN + 1]byte
	FriendsID   uint32
	FriendsName [MAX_PLAYER_NAME_LENGTH]byte
	/*
		Fakeplayer  int32
		Ishltv      int32
		#if defined( REPLAY_ENABLED )
			true if player is the Replay proxy
			bool			isreplay;
		#endif
			// custom files CRC for this player
			CRC32_t			customFiles[MAX_CUSTOM_FILES];
			// this counter increases each time the server downloaded a new file
			unsigned char	filesDownloaded;
	*/
}

type Userinfo struct {
	XUID        uint64 // network xuid
	Name        string // scoreboard information
	UserID      int    // local server user ID, unique while server is running
	GUID        string // global unique player identifer
	FriendsID   uint   // friends identification number
	FriendsName string // friends name
	SteamID     uint64
}
