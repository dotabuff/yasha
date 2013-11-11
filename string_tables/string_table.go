package string_tables

import (
	dota "github.com/elobuff/d2rp/dota"
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
}
