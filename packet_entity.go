package yasha

import (
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/dotabuff/yasha/bitstream"
	"github.com/dotabuff/yasha/dota"
)

type UpdateType int

const (
	Create UpdateType = iota
	Delete
	Leave
	Preserve
)
const serialNumBits = 11

type EntityDelta struct {
	ID     int
	Fields []int
}

type EntityState int

const (
	EntityStateDefault EntityState = iota
	EntityStateCreated
	EntityStateOverwritten
	EntityStateUpdated
	EntityStateDeleted
)

type PacketEntity struct {
	ID        int
	Type      UpdateType
	ClassId   int
	Class     *ClassInfo
	Name      string
	SerialNum int

	properties map[int]*SendProp
	flat       *FlatSendTable

	Tick         int
	Index        int
	EntityHandle int

	Values   map[string]interface{}
	Delta    map[string]interface{}
	OldDelta map[string]interface{}
}

func NewPacketEntity(id int, class *ClassInfo, flatSendTable *FlatSendTable, serialNum int) *PacketEntity {
	return &PacketEntity{
		ID:        id,
		Class:     class,
		flat:      flatSendTable,
		SerialNum: serialNum,
	}
}

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

func (pe *PacketEntity) update(id int, eClass *ClassInfo, flatSendTable *FlatSendTable) {
	panic("not implemented")
}
func (pe *PacketEntity) setState(state EntityState) {
	panic("not implemented")
}

func readUpdateType(bs *bitstream.BitStream) UpdateType {
	if !bs.ReadBool() {
		if bs.ReadBool() {
			return Create
		} else {
			return Preserve
		}
	} else if bs.ReadBool() {
		return Delete
	}
	return Leave
}

func readHeader(previousId int, bs *bitstream.BitStream) (id int, updateType UpdateType) {
	nId := bs.Read(6)

	switch nId & 0x30 {
	case 16:
		nId = (nId & 15) | bs.Read(4)<<4
	case 32:
		nId = (nId & 15) | bs.Read(8)<<4
	case 48:
		nId = (nId & 15) | bs.Read(28)<<4
	}

	updateType = readUpdateType(bs)

	return previousId + int(nId+1), updateType
}

const (
	instancebaseline = "instancebaseline"
	dotaMaxEntities  = 0x3FFF // 16383
)

func (p *Parser) handleEntity(tick int, pes *dota.CSVCMsg_PacketEntities) {
	bs := bitstream.NewBitStream(pes.GetEntityData())
	id := -1

	baseline := p.stringTables.ByName(instancebaseline)
	if baseline == nil {
		panic("couldn't find instancebaselinu")
	}

	var updateType UpdateType
	for i := int32(0); i < pes.GetUpdatedEntries(); i++ {
		id, updateType = readHeader(id, bs)

		if id > dotaMaxEntities {
			panic(spew.Errorf("entity id too large: %d > %d", id, dotaMaxEntities))
		}

		ent := p.Entities[id]

		switch updateType {
		case Create:
			pp("create", id)

			classId := int(bs.Read(p.classIdNumBits))
			serialNum := int(bs.Read(10))
			eClass := p.classInfos.byId[classId]
			if eClass == nil {
				panic("no fitting entity class found")
			}

			flatSendTable := p.flatSendTables[classId]
			if flatSendTable == nil {
				panic("no suitable flat send table found")
			}

			// pp(classId, serialNum, eClass, flatSendTable)
			if ent == nil {
				ent = NewPacketEntity(id, eClass, flatSendTable, serialNum)
				p.Entities[id] = ent
			} else {
				ent.update(id, eClass, flatSendTable)
				ent.setState(EntityStateOverwritten)
			}

			c := strconv.Itoa(classId)
			var base *StringTableItem
			for _, item := range baseline.Items {
				if item.Key == c {
					base = item
				}
			}

			if base == nil {
				panic("no matching instancebaseline found")
			}

			baseStream := bitstream.NewBitStream(base.Value)
			ent.updateFromBitStream(baseStream, nil)
			delta := &EntityDelta{}
			ent.updateFromBitStream(bs, delta)
		case Preserve:
		case Delete:
		case Leave:
			// ignore
		}
	}
}

func (e *PacketEntity) updateFromBitStream(bs *bitstream.BitStream, delta *EntityDelta) {
	fields := readFieldIds(bs)
	if e.properties == nil {
		e.properties = map[int]*SendProp{}
	}

	for _, field := range fields {
		// if field >= len(e.properties) {
		// 	panic(spew.Sprintf("unknown sendprop: %d > %d", field, len(e.properties)))
		// }

		property := e.properties[field]
		if property == nil {
			// pp("field:", field)
			// pp("property:", e.flat.Properties)
			// pp(e)
			e.properties[field] = SendPropFromStream(bs, e.flat.Properties[field].Prop)
			e.properties[field].Name = e.flat.Properties[field].Prop.Name
		} else {
			property.update(bs)
		}
	}
}

func readFieldIds(bs *bitstream.BitStream) []int {
	fields := []int{}
	field := -1
	for {
		if bs.ReadBool() {
			field += 1
			fields = append(fields, field)
		} else {
			value := bs.ReadVarUInt32()
			if value == 16383 {
				break
			}
			field += 1
			field += int(value)
			fields = append(fields, field)
		}
	}
	return fields
}
