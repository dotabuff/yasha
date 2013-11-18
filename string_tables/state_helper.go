package string_tables

import (
	"io/ioutil"
	"os"
	"strconv"

	"code.google.com/p/gogoprotobuf/proto"
	"github.com/davecgh/go-spew/spew"
	"github.com/elobuff/d2rp/core/parser"
	"github.com/elobuff/d2rp/core/send_tables"
	"github.com/elobuff/d2rp/core/utils"
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

	ClassInfosNameMapping map[int]string
	Mapping               map[int][]*send_tables.SendProp
	Multiples             map[int]map[string]int
	Baseline              map[int]map[string]interface{}
	pendingBaseline       []*StringTableItem
}

func NewStateHelper() *StateHelper {
	return &StateHelper{
		packets:         parser.ParserBaseItems{},
		metaTables:      map[int]*CacheItem{},
		baseTables:      map[int]*StringTable{},
		evolution:       map[int][]*StringTable{},
		current:         map[int]*StringTable{},
		Baseline:        map[int]map[string]interface{}{},
		pendingBaseline: []*StringTableItem{},
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

	switch table.Name {
	case "ActiveModifiers":
		parseActiveModifiers(table.Items)
	case "instancebaseline":
		helper.updateInstanceBaseline(table.Items)
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
	switch current.Name {
	case "ActiveModifiers":
		parseActiveModifiers(update)
	case "instancebaseline":
		helper.updateInstanceBaseline(update)
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

func (helper *StateHelper) updateInstanceBaseline(update map[int]*StringTableItem) {
	for _, item := range helper.pendingBaseline {
		helper.updateInstanceBaselineItem(item)
	}
	for _, item := range update {
		helper.updateInstanceBaselineItem(item)
	}
}

func (helper *StateHelper) updateInstanceBaselineItem(item *StringTableItem) {
	classId, err := strconv.Atoi(item.Str)
	if err != nil {
		panic(err)
	}

	className := helper.ClassInfosNameMapping[classId]
	if className == "DT_DOTAPlayer" {
		return
	}

	mapping := helper.Mapping[classId]
	multiples := helper.Multiples[classId]
	if len(mapping) < 1 || len(multiples) < 1 {
		helper.pendingBaseline = append(helper.pendingBaseline, item)
		return
	}

	baseline, found := helper.Baseline[classId]
	if !found {
		baseline = map[string]interface{}{}
		helper.Baseline[classId] = baseline
	}

	br := utils.NewBitReader(item.Data)
	indices := br.ReadPropertiesIndex()
	baseValues := br.ReadPropertiesValues(mapping, multiples, indices)
	for key, value := range baseValues {
		baseline[key] = value
	}
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
			if table.Tick >= tick {
				return state
			}
			state[table.Index] = table
		}
	}

	return state
}

func (helper *StateHelper) GetTableAtTick(tick int, tableName string) (result *StringTable) {
	for _, evo := range helper.evolution {
		for _, table := range evo {
			if table.Name == tableName {
				if table.Tick >= tick {
					return result
				}
				result = table
			}
		}
	}

	return result
}

func (helper *StateHelper) GetTableNow(tableName string) (result *StringTable) {
	for _, table := range helper.current {
		if table.Name == tableName {
			return table
		}
	}

	return
}
