package yasha

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/dotabuff/yasha/dota"
)

type DtHierarchy struct {
	Name string
	Prop *SendProp
}

type FlatSendTable struct {
	Name       string
	Properties []*DtHierarchy
}

func NewFlatSendTable(name string, properties []*DtHierarchy) *FlatSendTable {
	return &FlatSendTable{
		Name:       name,
		Properties: properties,
	}
}

type SendTable struct {
	Id          int
	Name        string
	Decodable   bool
	counter     int
	props       []*SendProp // insertion order
	propsById   map[int]*SendProp
	propsByName map[string]*SendProp
}

func NewSendTable(name string, decodable bool, id int) *SendTable {
	return &SendTable{
		Id:          id,
		Name:        name,
		Decodable:   decodable,
		counter:     0,
		props:       []*SendProp{},
		propsById:   map[int]*SendProp{},
		propsByName: map[string]*SendProp{},
	}
}

func (s *SendTable) Insert(prop *SendProp) {
	counter := s.counter
	s.counter++
	s.props = append(s.props, prop)
	s.propsById[counter] = prop
	s.propsByName[prop.Name] = prop
}

func (p *Parser) flattenSendTables() {
	/*
		FIXME: this miiiight just be important
		for _, table := range p.sendTables.tablesById {
			var last *SendProp
			for _, prop := range table.props {
				pp(prop)
				if prop.Type == DPT_Array {
					if last != nil {
						prop.ArrayType = last
					} else {
						panic("invalid array prop")
					}
				}
				last = prop
			}
		}
	*/

	for tableId, table := range p.sendTables.tablesById {
		excludes := p.buildExcludeList(table)
		props := []*DtHierarchy{}
		p.buildHierarchy(table, excludes, &props, table.Name)

		// set of all possible priorities
		priorities := make(map[int]bool, 64)
		for _, prop := range props {
			priorities[prop.Prop.Priority] = true
		}

		offset := 0
		for priority, _ := range priorities {
			cursor := offset

			for cursor < len(props) {
				prop := props[cursor].Prop
				if prop.Priority == priority || (SPROP_CHANGES_OFTEN&prop.Flags != 0 && priority == 64) {
					props[cursor], props[offset] = props[offset], props[cursor]
					offset++
				}
				cursor++
			}
		}

		pp(props)
		p.flatSendTables[tableId] = NewFlatSendTable(table.Name, props)
	}
}

func (p *Parser) handleSendTable(st *dota.CSVCMsg_SendTable) {
	sendTableId := p.sendTableId
	p.sendTableId++
	sendTableId++

	table := NewSendTable(st.GetNetTableName(), st.GetNeedsDecoder(), sendTableId)
	for _, prop := range st.GetProps() {
		table.Insert(NewSendProp(prop, table.Name))
	}
	p.sendTables.Insert(table)
}

type SendTables struct {
	tablesByName map[string]*SendTable
	tablesById   map[int]*SendTable
}

func NewSendTables() *SendTables {
	return &SendTables{
		tablesByName: map[string]*SendTable{},
		tablesById:   map[int]*SendTable{},
	}
}

func (s *SendTables) Insert(table *SendTable) {
	s.tablesByName[table.Name] = table
	s.tablesById[table.Id] = table
}

func (s *SendTables) ByName(name string) *SendTable {
	return s.tablesByName[name]
}

func (s *SendTables) ById(id int) *SendTable {
	return s.tablesById[id]
}

func (p *Parser) buildExcludeList(sendTable *SendTable) []*SendProp {
	result := []*SendProp{}

	for _, prop := range sendTable.propsById {
		if prop.Flags&SPROP_EXCLUDE != 0 {
			result = append(result, prop)
		}
	}

	for _, prop := range sendTable.propsById {
		if prop.Type == DPT_DataTable {
			inner := p.sendTables.ByName(prop.ClassName)
			result = append(result, p.buildExcludeList(inner)...)
		}
	}

	return result
}

func (p *Parser) buildHierarchy(table *SendTable, excludes []*SendProp, props *[]*DtHierarchy, base string) {
	dtProp := []*DtHierarchy{}
	p.gatherProperties(table, &dtProp, excludes, props, base)
	for _, prop := range dtProp {
		*props = append(*props, prop)
	}
}

func (p *Parser) gatherProperties(table *SendTable, dtProp *[]*DtHierarchy, excludes []*SendProp, props *[]*DtHierarchy, base string) {
	for _, prop := range table.propsById {
		if (SPROP_EXCLUDE|SPROP_INSIDEARRAY)&prop.Flags != 0 {
			continue
		} else if prop.isExcluded(table.Name, excludes) {
			continue
		}

		if prop.Type == DPT_DataTable {
			innerTable := p.sendTables.ByName(prop.ClassName)
			if innerTable == nil {
				panic(spew.Errorf("couldn't find inner sendtable: %p", prop.ClassName))
			}

			if SPROP_COLLAPSIBLE&prop.Flags != 0 {
				p.gatherProperties(innerTable, dtProp, excludes, props, base)
			} else {
				p.buildHierarchy(innerTable, excludes, props, base+"."+prop.ClassName+"."+prop.Name)
			}
		} else {
			*dtProp = append(*dtProp, &DtHierarchy{base + "." + prop.Name, prop})
		}
	}
}
