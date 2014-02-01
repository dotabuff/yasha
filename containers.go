package core

type AbilityTracker struct {
	HeroHandle int
	Level      int
	Tick       int
	Name       string
}

type Abilities []*AbilityTracker

func (p Abilities) Len() int           { return len(p) }
func (p Abilities) Less(i, j int) bool { return p[i].Tick < p[j].Tick }
func (p Abilities) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type LastHitTracker struct {
	HeroHandle int
	Tick       int
	LastHit    int
}

type LastHits []*LastHitTracker

func (p LastHits) Len() int           { return len(p) }
func (p LastHits) Less(i, j int) bool { return p[i].Tick < p[j].Tick }
func (p LastHits) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Dotabuff struct {
	Players []*Player
}

type Player struct {
	Abilities               Abilities
	LastHits                LastHits
	HeroHandle              int
	AFK                     bool
	Assists                 uint
	BattleBonusActive       bool
	BattleBonusRate         int
	BroadcasterChannel      int
	BroadcasterChannelSlot  int
	BroadcasterLanguage     int
	BuybackCooldownTime     float32
	ConnectionState         int
	Deaths                  uint
	DenyCount               uint
	FakeClient              bool
	FullyJoinedServer       bool
	HasRandomed             bool
	HasRepicked             bool
	IsBroadcaster           bool
	Kills                   uint
	LastBuybackTime         int
	LastHitCount            uint
	LastHitMultikill        uint
	LastHitStreak           uint
	Level                   int
	MetaExperienceAwarded   int
	MetaExperienceBonusRate int
	MetaExperience          int
	MetaLevel               int
	PlayerNames             string
	PlayerSteamIDs          uint64
	PlayerTeams             int
	PossibleHeroSelection   int
	ReliableGold            uint
	RespawnSeconds          int
	SelectedHeroID          int
	SelectedHero            string
	Streak                  uint
	SuggestedHeroes         int
	TimedRewardCrates       int
	TimedRewardDrops        int
	TotalEarnedGold         uint
	TotalEarnedXP           uint
	UnitShareMasks          int
	UnreliableGold          uint
	VoiceChatBanned         bool
}

type CombatLogEntry struct {
	Type               string
	SourceName         string
	TargetName         string
	AttackerName       string
	InflictorName      string
	AttackerIsIllusion bool
	TargetIsIllusion   bool
	Value              int32
	Health             int32
	Timestamp          float32
	TargetSourceName   string
}
