package core

import (
	"reflect"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/dotabuff/d2rp/core/string_tables"
	"github.com/dotabuff/d2rp/dota"
)

type CombatLogEntry interface {
	Type() dota.DOTA_COMBATLOG_TYPES
	Timestamp() float32
}

type combatLogParser struct {
	stsh     *string_tables.StateHelper
	distinct map[dota.DOTA_COMBATLOG_TYPES][]map[interface{}]bool
}

/*
The default is mostly:

1 SourceName
2 TargetName
3 AttackerName
4 InflictorName
5 AttackerIsillusion
6 TargetIsIllusion
7 Value
8 Health
9 Timestamp
10 TargetSourceName
11 TimestampRaw
12 AttackerIsHero
13 TargetIsHero
*/
func (c combatLogParser) parse(obj *dota.CSVCMsg_GameEvent) CombatLogEntry {
	keys := obj.GetKeys()

	var v CombatLogEntry
	t := dota.DOTA_COMBATLOG_TYPES(keys[0].GetValByte())
	switch t {
	case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_ABILITY:
		v = &CombatLogAbility{}
	case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_DAMAGE:
		v = &CombatLogDamage{}
	case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_DEATH:
		v = &CombatLogDeath{}
	case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_GAME_STATE:
		v = &CombatLogGameState{}
	case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_GOLD:
		v = &CombatLogGold{}
	case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_HEAL:
		v = &CombatLogHeal{}
	case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_ITEM:
		v = &CombatLogItem{}
	case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_LOCATION:
		v = &CombatLogLocation{}
	case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MODIFIER_ADD:
		v = &CombatLogModifierAdd{}
	case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MODIFIER_REMOVE:
		v = &CombatLogModifierRemove{}
	case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_PURCHASE:
		v = &CombatLogPurchase{}
	case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_XP:
		v = &CombatLogXP{}
	default:
		pp(t, keys)
		return nil
	}

	c.assign(v, keys)
	return v
}

type CombatLogItem struct {
	Target             string  `logIndex:"2" logTable:"CombatLogNames"`
	User               string  `logIndex:"3" logTable:"CombatLogNames"`
	Item               string  `logIndex:"4" logTable:"CombatLogNames"`
	AttackerIsIllusion bool    `logIndex:"5"`
	TargetIsIllusion   bool    `logIndex:"6"`
	Time               float32 `logIndex:"9"`
	UserIsHero         bool    `logIndex:"12"`
	TargetIsHero       bool    `logIndex:"13"`
}

func (c CombatLogItem) Type() dota.DOTA_COMBATLOG_TYPES {
	return dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_ITEM
}
func (c CombatLogItem) Timestamp() float32 {
	return c.Time
}

// TODO: only observed 2,3,4,7,13 so far, but the others would make sense too.
type CombatLogAbility struct {
	Target             string  `logIndex:"2" logTable:"CombatLogNames"`
	Attacker           string  `logIndex:"3" logTable:"CombatLogNames"`
	Ability            string  `logIndex:"4" logTable:"CombatLogNames"`
	AttackerIsIllusion bool    `logIndex:"5"`
	TargetIsIllusion   bool    `logIndex:"6"`
	IsDebuff           int     `logIndex:"7"` // seen values: 0,1,2
	Time               float32 `logIndex:"9"`
	TargetSource       string  `logIndex:"10" logTable:"CombatLogNames"`
	AttackerIsHero     bool    `logIndex:"12"`
	TargetIsHero       bool    `logIndex:"13"`
}

func (c CombatLogAbility) Type() dota.DOTA_COMBATLOG_TYPES {
	return dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_ABILITY
}
func (c CombatLogAbility) Timestamp() float32 {
	return c.Time
}

/*
 0: val_byte:0 ),
 1: val_short:3 ),
 2: val_short:27 ),
 3: val_short:3 ),
 4: val_short:0 ),
 5: val_bool:false ),
 6: val_bool:false ),
 7: val_short:70 ),
 8: val_short:429 ),
 9: val_float:229.45338 ),
 10: val_short:27 ),
 11: val_float:238.43335 ),
 12: val_bool:true ),
 13: val_bool:false )
*/
type CombatLogDamage struct {
	Source             string  `logIndex:"1" logTable:"CombatLogNames"`
	Target             string  `logIndex:"2" logTable:"CombatLogNames"`
	Attacker           string  `logIndex:"3" logTable:"CombatLogNames"`
	Cause              string  `logIndex:"4" logTable:"CombatLogNames"`
	AttackerIsIllusion bool    `logIndex:"5"`
	TargetIsIllusion   bool    `logIndex:"6"`
	Value              int     `logIndex:"7"`
	Health             int     `logIndex:"8"`
	Time               float32 `logIndex:"9"`
	TargetSource       string  `logIndex:"10" logTable:"CombatLogNames"`
	AttackerIsHero     bool    `logIndex:"12"`
	TargetIsHero       bool    `logIndex:"13"`
}

func (c CombatLogDamage) Type() dota.DOTA_COMBATLOG_TYPES {
	return dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_DAMAGE
}
func (c CombatLogDamage) Timestamp() float32 {
	return c.Time
}

type CombatLogLocation struct {
	Source             string  `logIndex:"1" logTable:"CombatLogNames"`
	Target             string  `logIndex:"2" logTable:"CombatLogNames"`
	Attacker           string  `logIndex:"3" logTable:"CombatLogNames"`
	Modifier           string  `logIndex:"4" logTable:"CombatLogNames"`
	AttackerIsIllusion bool    `logIndex:"5"`
	TargetIsIllusion   bool    `logIndex:"6"`
	Value              int     `logIndex:"7"`
	Health             int     `logIndex:"8"`
	Time               float32 `logIndex:"9"`
	TargetSource       string  `logIndex:"10" logTable:"CombatLogNames"`
	AttackerIsHero     bool    `logIndex:"12"`
	TargetIsHero       bool    `logIndex:"13"`
}

func (c CombatLogLocation) Type() dota.DOTA_COMBATLOG_TYPES {
	return dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_LOCATION
}
func (c CombatLogLocation) Timestamp() float32 {
	return c.Time
}

type CombatLogHeal struct {
	Source             string  `logIndex:"1" logTable:"CombatLogNames"`
	Target             string  `logIndex:"2" logTable:"CombatLogNames"`
	Attacker           string  `logIndex:"3" logTable:"CombatLogNames"`
	Modifier           string  `logIndex:"4" logTable:"CombatLogNames"`
	AttackerIsIllusion bool    `logIndex:"5"`
	TargetIsIllusion   bool    `logIndex:"6"`
	Value              int     `logIndex:"7"`
	Health             int     `logIndex:"8"`
	Time               float32 `logIndex:"9"`
	TargetSource       string  `logIndex:"10" logTable:"CombatLogNames"`
	AttackerIsHero     bool    `logIndex:"12"`
	TargetIsHero       bool    `logIndex:"13"`
}

func (c CombatLogHeal) Type() dota.DOTA_COMBATLOG_TYPES {
	return dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_HEAL
}
func (c CombatLogHeal) Timestamp() float32 {
	return c.Time
}

type CombatLogModifierAdd struct {
	Source             string  `logIndex:"1" logTable:"CombatLogNames"`
	Target             string  `logIndex:"2" logTable:"CombatLogNames"`
	Attacker           string  `logIndex:"3" logTable:"CombatLogNames"`
	Modifier           string  `logIndex:"4" logTable:"CombatLogNames"`
	AttackerIsIllusion bool    `logIndex:"5"`
	TargetIsIllusion   bool    `logIndex:"6"`
	IsDebuff           bool    `logIndex:"7"`
	Health             int     `logIndex:"8"`
	Time               float32 `logIndex:"9"`
	TargetSource       string  `logIndex:"10" logTable:"CombatLogNames"`
	AttackerIsHero     bool    `logIndex:"12"`
	TargetIsHero       bool    `logIndex:"13"`
}

func (c CombatLogModifierAdd) Type() dota.DOTA_COMBATLOG_TYPES {
	return dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MODIFIER_ADD
}

func (c CombatLogModifierAdd) Timestamp() float32 {
	return c.Time
}

type CombatLogModifierRemove struct {
	Target             string  `logIndex:"2" logTable:"CombatLogNames"`
	Caster             string  `logIndex:"3" logTable:"CombatLogNames"`
	Modifier           string  `logIndex:"4" logTable:"CombatLogNames"`
	AttackerIsIllusion bool    `logIndex:"5"`
	TargetIsIllusion   bool    `logIndex:"6"`
	IsDebuff           bool    `logIndex:"7"`
	Health             int     `logIndex:"8"`
	Time               float32 `logIndex:"9"`
	AttackerIsHero     bool    `logIndex:"12"`
	TargetIsHero       bool    `logIndex:"13"`
}

func (c CombatLogModifierRemove) Type() dota.DOTA_COMBATLOG_TYPES {
	return dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MODIFIER_REMOVE
}
func (c CombatLogModifierRemove) Timestamp() float32 {
	return c.Time
}

type CombatLogDeath struct {
	Source             string  `logIndex:"1" logTable:"CombatLogNames"`
	Target             string  `logIndex:"2" logTable:"CombatLogNames"`
	Attacker           string  `logIndex:"3" logTable:"CombatLogNames"`
	Cause              string  `logIndex:"4" logTable:"CombatLogNames"`
	AttackerIsIllusion bool    `logIndex:"5"`
	TargetIsIllusion   bool    `logIndex:"6"`
	Time               float32 `logIndex:"9"`
	TargetSource       string  `logIndex:"10" logTable:"CombatLogNames"`
	AttackerIsHero     bool    `logIndex:"12"`
	TargetIsHero       bool    `logIndex:"13"`
}

func (c CombatLogDeath) Type() dota.DOTA_COMBATLOG_TYPES {
	return dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_DEATH
}
func (c CombatLogDeath) Timestamp() float32 {
	return c.Time
}

type CombatLogPurchase struct {
	Time  float32 `logIndex:"9"`
	Buyer string  `logIndex:"2" logTable:"CombatLogNames"`
	Item  string  `logIndex:"7" logTable:"CombatLogNames"`
}

func (c CombatLogPurchase) Type() dota.DOTA_COMBATLOG_TYPES {
	return dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_PURCHASE
}
func (c CombatLogPurchase) Timestamp() float32 {
	return c.Time
}

type CombatLogGold struct {
	Target int     `logIndex:"2"`
	Value  int     `logIndex:"7"`
	Time   float32 `logIndex:"9"`
}

func (c CombatLogGold) Type() dota.DOTA_COMBATLOG_TYPES {
	return dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_GOLD
}
func (c CombatLogGold) Timestamp() float32 {
	return c.Time
}

type CombatLogGameState struct {
	State int     `logIndex:"7"`
	Time  float32 `logIndex:"9"`
}

func (c CombatLogGameState) Type() dota.DOTA_COMBATLOG_TYPES {
	return dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_GAME_STATE
}
func (c CombatLogGameState) Timestamp() float32 {
	return c.Time
}

type CombatLogXP struct {
	Target string  `logIndex:"2" logTable:"CombatLogNames"`
	Value  int     `logIndex:"7"`
	Time   float32 `logIndex:"9"`
}

func (c CombatLogXP) Type() dota.DOTA_COMBATLOG_TYPES {
	return dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_XP
}
func (c CombatLogXP) Timestamp() float32 {
	return c.Time
}

func atoi(a string) int {
	i, _ := strconv.Atoi(a)
	return i
}

func strtbl(table map[int]*string_tables.StringTableItem, keys []*dota.CSVCMsg_GameEventKeyT, index int) string {
	return table[int(keys[index].GetValShort())].Str
}

func (c combatLogParser) logDistinct(t dota.DOTA_COMBATLOG_TYPES, keys []*dota.CSVCMsg_GameEventKeyT) {
	if c.distinct[t] == nil {
		c.distinct[t] = make([]map[interface{}]bool, len(keys))
	}
	for i, key := range keys {
		if i == 9 || i == 11 {
			continue
		}
		if c.distinct[t][i] == nil {
			c.distinct[t][i] = map[interface{}]bool{}
		}
		switch key.GetType() {
		case 2:
			c.distinct[t][i][key.GetValFloat()] = true
		case 4:
			c.distinct[t][i][key.GetValShort()] = true
		case 5:
			c.distinct[t][i][key.GetValByte()] = true
		case 6:
			c.distinct[t][i][key.GetValBool()] = true
		}
	}
}

func (c combatLogParser) assign(v CombatLogEntry, keys []*dota.CSVCMsg_GameEventKeyT) {
	rv := reflect.ValueOf(v).Elem()
	rt := rv.Type()
	fieldIndices := make([]int, rv.NumField())
	for i, _ := range fieldIndices {
		fieldTag := rt.Field(i).Tag
		logIndex := atoi(fieldTag.Get("logIndex"))

		// this allows us to be backwards compatible, but we still need to care in
		// the code that uses the log that some of the newer fields may not be set.
		if logIndex <= 0 || logIndex >= len(keys) {
			continue
		}

		field := rv.Field(i)
		key := keys[logIndex]

		switch key.GetType() {
		case 2:
			field.SetFloat(float64(key.GetValFloat()))
		case 4:
			valShort := key.GetValShort()
			if logTable := fieldTag.Get("logTable"); logTable != "" {
				table := c.stsh.GetTableNow(logTable).Items
				field.SetString(table[int(valShort)].Str)
			} else if field.Kind() == reflect.Bool {
				field.SetBool(valShort == 1)
			} else {
				field.SetInt(int64(valShort))
			}
		case 5:
			field.SetInt(int64(key.GetValByte()))
		case 6:
			field.SetBool(key.GetValBool())
		default:
			panic("unknown GameEventKey Type" + spew.Sdump(key) + " in " + spew.Sdump(keys))
		}
	}
}
