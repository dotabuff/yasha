package string_tables

import (
	"io/ioutil"
	"os"
	"sort"

	"code.google.com/p/gogoprotobuf/proto"
	"github.com/davecgh/go-spew/spew"
	"github.com/elobuff/d2rp/core/parser"
	dota "github.com/elobuff/d2rp/dota"
)

func p(v ...interface{}) { spew.Dump(v...) }

type CacheItem struct {
	Bits        int
	IsFixedSize bool
	MaxEntries  int
	Name        string
}

type StateHelper struct {
	packets parser.ParserBaseItems

	lastCreateIndex int
	// contains meta information for parsing the table
	metaTables map[int]*CacheItem
	// the first CSTs
	baseTables map[int]*StringTable
	// every UST we get, we calculate the ST and put it in here.
	evolution map[int][]*StringTable
	// every UST we get, we calculate the ST and put it in here.
	current map[int]*StringTable
}

func NewStateHelper() *StateHelper {
	return &StateHelper{
		packets:    parser.ParserBaseItems{},
		metaTables: map[int]*CacheItem{},
		baseTables: map[int]*StringTable{},
		evolution:  map[int][]*StringTable{},
		current:    map[int]*StringTable{},
	}
}

func writeStringTables(directory string, tick int, t string) {
	err := os.MkdirAll(directory, os.ModePerm|os.ModeDir)
	if err != nil {
		panic(err)
	}
	path := spew.Sprintf("%s/tick_%010d.txt", directory, tick)
	err = ioutil.WriteFile(path, []byte(t), 0644)
	if err != nil {
		panic(err)
	}
}

func (helper *StateHelper) AppendPacket(packet *parser.ParserBaseItem) {
	switch obj := packet.Object.(type) {
	case *dota.CDemoStringTables:
		// looks like we don't need them, they are just like CST
		// helper.OnCDST(packet.Tick, obj)
	case *dota.CSVCMsg_CreateStringTable:
		helper.packets = append(helper.packets, packet)
		helper.OnCST(packet.Tick, obj)
	case *dota.CSVCMsg_UpdateStringTable:
		helper.packets = append(helper.packets, packet)
		helper.OnUST(packet.Tick, obj)
	default:
		panic("Cannot handle this type")
	}
}

func (helper *StateHelper) OnCDST(tick int, obj *dota.CDemoStringTables) {
	/*
		for _, table := range obj.GetTables() {
			writeStringTables("CDemoStringTables/"+table.GetTableName(), tick, spew.Sdump(table))
		}
	*/
}

func (helper *StateHelper) OnCST(tick int, obj *dota.CSVCMsg_CreateStringTable) {
	if tick != 0 {
		// tested against a ton of replays, hasn't happened yet...
		panic("creating string table after first tick?")
	}

	helper.metaTables[helper.lastCreateIndex] = &CacheItem{
		Bits:        int(obj.GetUserDataSizeBits()),
		IsFixedSize: obj.GetUserDataFixedSize(),
		MaxEntries:  int(obj.GetMaxEntries()),
		Name:        obj.GetName(),
	}

	table := StringTable{
		Index: helper.lastCreateIndex,
		Name:  obj.GetName(),
		Tick:  tick,
		Items: Parse(
			obj.GetStringData(),
			obj.GetNumEntries(),
			obj.GetMaxEntries(),
			obj.GetUserDataFixedSize(),
			obj.GetUserDataSizeBits(),
		),
	}

	if table.Name == "ActiveModifiers" {
		parseActiveModifiers(table.Items)
	}

	// writeStringTables("CreateStringTable/"+table.Name, tick, spew.Sdump(table))

	helper.baseTables[helper.lastCreateIndex] = &table
	helper.current[helper.lastCreateIndex] = &table
	helper.evolution[helper.lastCreateIndex] = append(helper.evolution[helper.lastCreateIndex], &table)

	helper.lastCreateIndex++
}

// NOTE:
// We ignore the "userinfo" table decoding process since it's a PITA and has no useful info anyway.
// In case we ever need it, a struct describing the binary is at:
// https://github.com/mitsuhiko/dota2-demoinfo2/blob/4ca45a87c631787eab140d313a3f21210b543741/demofile.h#L48
func (helper *StateHelper) OnUST(tick int, obj *dota.CSVCMsg_UpdateStringTable) {
	tableId := int(obj.GetTableId())

	meta := helper.metaTables[tableId]
	update := Parse(
		obj.GetStringData(),
		obj.GetNumChangedEntries(),
		int32(meta.MaxEntries),
		meta.IsFixedSize,
		int32(meta.Bits),
	)

	current := helper.current[tableId]
	if current.Name == "ActiveModifiers" {
		parseActiveModifiers(update)
	}
	// writeStringTables("UpdateStringTable/"+current.Name, tick, spew.Sdump(update))

	for key, value := range update {
		current.Items[key] = value
	}

	stCopy := &StringTable{
		Index: current.Index,
		Items: map[int]*StringTableItem{},
		Name:  current.Name,
		Tick:  current.Tick,
	}

	for key, value := range current.Items {
		stCopy.Items[key] = value
	}

	helper.evolution[tableId] = append(helper.evolution[tableId], stCopy)
}

func parseActiveModifiers(entries map[int]*StringTableItem) {
	for _, e := range entries {
		if len(e.Data) > 0 {
			o := &dota.CDOTAModifierBuffTableEntry{}
			err := proto.Unmarshal(e.Data, o)
			if err != nil {
				panic(err)
			}
			e.Data = e.Data[:0]
			e.ModifierBuff = o
		} else {
			// spew.Dump(e)
		}
	}
}

func (helper *StateHelper) GetStateAtTick(tick int) map[int]*StringTable {
	state := map[int]*StringTable{}

	for _, evo := range helper.evolution {
		for _, table := range evo {
			if table.Tick > tick {
				return state
			}
			state[table.Index] = table
		}
	}

	return state
}

/*
	index := -1
	packets := parser.ParserBaseItems{}
	for _, item := range helper.packets {
		if item.Tick <= tick {
			packets = append(packets, item)
		}
	}
	sort.Sort(packets)

	result := map[int]*StringTable{}

	for _, item := range packets {
		switch t := item.Object.(type) {
		case *dota.CSVCMsg_CreateStringTable:
			index++
			st := &StringTable{
				Index: index,
				Name:  t.GetName(),
				Tick:  item.Tick,
				Items: Parse(
					t.GetStringData(),
					t.GetNumEntries(),
					t.GetMaxEntries(),
					t.GetUserDataFixedSize(),
					t.GetUserDataSizeBits(),
				),
			}
			result[index] = st
		case *dota.CSVCMsg_UpdateStringTable:
			stc := helper.metaTables[int(t.GetTableId())]
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
*/

// Not used anywhere atm
func (helper *StateHelper) RepopulateCache() {
	sortItems := parser.ParserBaseItems{}
	for _, sortItem := range helper.packets {
		sortItems = append(sortItems, sortItem)
	}
	sort.Sort(sortItems)

	index := 0
	helper.metaTables = map[int]*CacheItem{}

	for _, sortItem := range sortItems {
		switch item := sortItem.Object.(type) {
		case *dota.CSVCMsg_CreateStringTable:
			helper.metaTables[index] = &CacheItem{
				Bits:        int(item.GetUserDataSizeBits()),
				IsFixedSize: item.GetUserDataFixedSize(),
				MaxEntries:  int(item.GetMaxEntries()),
				Name:        item.GetName(),
			}
			index++
		}
	}
}
