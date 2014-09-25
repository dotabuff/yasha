package yasha

import "github.com/dotabuff/yasha/utils"

type UpdateType int

const (
	Create UpdateType = iota
	Delete
	Leave
	Preserve
)

type PacketEntity struct {
	Tick         int
	Index        int
	SerialNum    int
	ClassId      int
	EntityHandle int
	Name         string
	Type         UpdateType
	Values       map[string]interface{}
	Delta        map[string]interface{}
	OldDelta     map[string]interface{}
}

const serialNumBits = 11

func (pe *PacketEntity) Handle() int {
	return pe.Index | (pe.SerialNum << serialNumBits)
}

func (pe *PacketEntity) Clone() *PacketEntity {
	values := map[string]interface{}{}
	for key, value := range pe.Values {
		values[key] = value
	}
	return &PacketEntity{
		Tick:         pe.Tick,
		Index:        pe.Index,
		SerialNum:    pe.SerialNum,
		ClassId:      pe.ClassId,
		EntityHandle: pe.EntityHandle,
		Name:         pe.Name,
		Type:         pe.Type,
		Values:       values,
	}
}

func ReadUpdateType(br *utils.BitReader) UpdateType {
	result := Preserve
	if !br.ReadBoolean() {
		if br.ReadBoolean() {
			result = Create
		}
	} else {
		result = Leave
		if br.ReadBoolean() {
			result = Delete
		}
	}
	return result
}
