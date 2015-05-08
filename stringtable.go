package yasha

import (
	"math"

	"github.com/dotabuff/yasha/bitstream"
	"github.com/dotabuff/yasha/dota"
)

type StringTables struct {
	byId   map[int]*StringTable
	byName map[string]*StringTable
}

func NewStringTables() *StringTables {
	return &StringTables{
		byId:   map[int]*StringTable{},
		byName: map[string]*StringTable{},
	}
}

func (s *StringTables) Insert(name string, id int, table *StringTable) {
	s.byName[name] = table
	s.byId[id] = table
}

func (s *StringTables) ById(id int) (table *StringTable) {
	return s.byId[id]
}

func (s *StringTables) ByName(name string) (table *StringTable) {
	return s.byName[name]
}

type StringTable struct {
	Name              string
	MaxEntries        int
	UserDataFixedSize bool
	UserDataSize      int
	UserDataSizeBits  int
	Flags             int
	Items             map[int]*StringTableItem
}

func NewStringTable(cst *dota.CSVCMsg_CreateStringTable) *StringTable {
	st := &StringTable{
		Name:              cst.GetName(),
		MaxEntries:        int(cst.GetMaxEntries()),
		UserDataFixedSize: cst.GetUserDataFixedSize(),
		UserDataSize:      int(cst.GetUserDataSize()),
		UserDataSizeBits:  int(cst.GetUserDataSizeBits()),
		Flags:             int(cst.GetFlags()),
		Items:             map[int]*StringTableItem{},
	}

	st.Update(int(cst.GetNumEntries()), cst.GetStringData())

	return st
}

type StringTableItem struct {
	Key   string
	Value []byte
}

func (p *Parser) handleCreateStringTable(cst *dota.CSVCMsg_CreateStringTable) {
	tableId := p.stringTableId
	p.stringTableId++
	p.stringTables.Insert(cst.GetName(), tableId, NewStringTable(cst))
}

func (p *Parser) handleUpdateStringTable(ust *dota.CSVCMsg_UpdateStringTable) {
	table := p.stringTables.ById(int(ust.GetTableId()))
	table.Update(int(ust.GetNumChangedEntries()), ust.GetStringData())
}

const (
	KeyHistorySize = 32
	MaxKeySize     = 1024
	MaxValueSize   = 16384
)

func (s *StringTable) Update(entries int, data []byte) {
	b := bitstream.NewBitStream(data)

	full := b.ReadBool()
	index := -1
	keys := []string{}
	items := map[int]*StringTableItem{}
	indexSize := int(math.Ceil(math.Log2(float64(s.MaxEntries))))

	for len(items) < entries {
		var key string
		var value []byte

		increment := b.ReadBool()
		if increment {
			index++
		} else {
			index = int(b.Read(indexSize))
		}

		hasName := b.ReadBool()
		if hasName {
			if full && b.ReadBool() {
				panic("stringtable key missing")
			}

			substring := b.ReadBool()
			if substring {
				sIndex := int(b.Read(5))
				sLength := int(b.Read(5))

				if sIndex >= KeyHistorySize || sLength >= MaxKeySize {
					panic("malformed substring")
				}

				if len(keys) <= sIndex {
					key = b.ReadString(MaxKeySize)
				} else {
					key = keys[sIndex][0:sLength] + b.ReadString(MaxKeySize-sLength)
				}
			} else {
				key = b.ReadString(MaxKeySize)
			}

			if len(keys) >= KeyHistorySize {
				copy(keys[0:], keys[1:])
				keys[len(keys)-1] = ""
				keys = keys[:len(keys)-1]
			}

			keys = append(keys, key)
		}

		hasValue := b.ReadBool()
		var length int
		if hasValue {
			var valSize int
			if s.UserDataFixedSize {
				length = s.UserDataSize
				valSize = s.UserDataSizeBits
			} else {
				length = int(b.Read(14))
				valSize = length * 8
			}

			if length > MaxValueSize {
				panic("stringtable value overflow")
			}

			value = b.ReadBits(valSize)
		}

		// pp(s.Name, key, value)

		items[index] = &StringTableItem{key, value}
	}

	for index, item := range items {
		s.Items[index] = item
	}
}
