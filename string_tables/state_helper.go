package string_tables

import (
	"sort"

	"github.com/davecgh/go-spew/spew"
	"github.com/elobuff/d2rp/core/parser"
	dota "github.com/elobuff/d2rp/dota"
)

type Cache struct {
	Bits        int
	IsFixedSize bool
	MaxEntries  int
	Name        string
}
type CreateStringTables []*dota.CSVCMsg_CreateStringTable

type StateHelper struct {
	cache         map[int]*Cache
	packets       parser.ParserBaseItems
	fullPackets   parser.ParserBaseItems
	lastIndexUsed int
}

func NewStateHelper(items parser.ParserBaseItems) *StateHelper {
	packets := parser.ParserBaseItems{}
	fullPackets := parser.ParserBaseItems{}

	for _, item := range items {
		switch item.Object.(type) {
		case *dota.CSVCMsg_CreateStringTable, *dota.CSVCMsg_UpdateStringTable:
			packets = append(packets, item)
		case *dota.CDemoStringTables, *dota.CDemoFullPacket:
			fullPackets = append(fullPackets, item)
		}
	}
	helper := &StateHelper{packets: packets, fullPackets: fullPackets}
	helper.populateCache()
	return helper
}

func (helper *StateHelper) GetStateAtTick(tick int) map[int]*StringTable {
	spew.Println("<GetStateAtTick tick=", tick, ">")
	defer spew.Println("</GetStateAtTick>")
	helper.lastIndexUsed = -1
	result := map[int]*StringTable{}
	packets := parser.ParserBaseItems{}
	for _, item := range helper.packets {
		if item.Tick <= tick {
			packets = append(packets, item)
		}
	}
	sort.Sort(packets)

	for _, item := range packets {
		switch t := item.Object.(type) {
		case *dota.CSVCMsg_CreateStringTable:
			helper.lastIndexUsed += 1
			st := &StringTable{
				Index: helper.lastIndexUsed,
				Name:  t.GetName(),
				Tick:  item.Tick,
			}
			st.Items = Parse(
				t.GetStringData(),
				t.GetNumEntries(),
				t.GetMaxEntries(),
				t.GetUserDataFixedSize(),
				t.GetUserDataSizeBits(),
			)
			result[helper.lastIndexUsed] = st
		case *dota.CSVCMsg_UpdateStringTable:
			stc := helper.cache[int(t.GetTableId())]
			ustr := Parse(
				t.GetStringData(),
				t.GetNumChangedEntries(),
				int32(stc.MaxEntries),
				stc.IsFixedSize,
				int32(stc.Bits),
			)
			for key, value := range ustr {
				resItems := result[int(t.GetTableId())].Items
				if innerItem, exists := resItems[key]; exists {
					innerItem.Str = value.Str
					innerItem.Data = value.Data
				} else {
					result[int(t.GetTableId())].Items[key] = value
				}
			}
		default:
			panic("this shouldn't be happening")
		}
	}

	return result
}

func (helper *StateHelper) populateCache() {
	helper.lastIndexUsed = -1
	helper.cache = map[int]*Cache{}

	sortItems := parser.ParserBaseItems{}
	for _, item := range helper.packets {
		sortItems = append(sortItems, item)
	}
	sort.Sort(sortItems)

	items := []*dota.CSVCMsg_CreateStringTable{}
	for _, item := range sortItems {
		switch cst := item.Object.(type) {
		case *dota.CSVCMsg_CreateStringTable:
			items = append(items, cst)
		}
	}

	for _, item := range items {
		helper.lastIndexUsed += 1
		c := &Cache{
			Bits:        int(item.GetUserDataSizeBits()),
			IsFixedSize: item.GetUserDataFixedSize(),
			MaxEntries:  int(item.GetMaxEntries()),
			Name:        item.GetName(),
		}
		helper.cache[helper.lastIndexUsed] = c
	}
}
