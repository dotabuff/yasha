package send_tables

import (
	"sort"

	"github.com/dotabuff/yasha/dota"
)

type Helper struct {
	sendTables       map[string]*dota.CSVCMsg_SendTable
	flatSendTable    []*SendProp
	excludedSendProp []*SendProp
}

func NewHelper() *Helper {
	return &Helper{
		sendTables:       map[string]*dota.CSVCMsg_SendTable{},
		flatSendTable:    []*SendProp{},
		excludedSendProp: []*SendProp{},
	}
}

func (sth *Helper) SetSendTable(name string, table *dota.CSVCMsg_SendTable) {
	sth.sendTables[name] = table
}

func (sth *Helper) LoadSendTable(netName string) []*SendProp {
	sth.flatSendTable = []*SendProp{}
	sth.excludedSendProp = []*SendProp{}
	sth.excludedSendProp = sth.getPropsExcluded(netName)
	sth.buildHierarchy(netName)
	sth.sortByPriority()
	return sth.flatSendTable
}

func (sth *Helper) getPropsExcluded(netName string) []*SendProp {
	result := []*SendProp{}
	sendTable := sth.sendTables[netName]
	for _, prop := range sendTable.GetProps() {
		if (prop.GetFlags() & int32(SPROP_EXCLUDE)) != 0 {
			result = append(result, NewSendProp(prop, netName))
		}
	}

	for _, prop := range sendTable.GetProps() {
		if prop.GetType() == 6 {
			result = append(result, sth.getPropsExcluded(prop.GetDtName())...)
		}
	}

	return result
}

func (sth *Helper) buildHierarchy(netName string) {
	result := []*SendProp{}
	sth.buildHierarchyIterateProps(netName, &result)
	sth.flatSendTable = append(sth.flatSendTable, result...)
}

func (sth *Helper) buildHierarchyIterateProps(netName string, result *[]*SendProp) {
	pTable := sth.sendTables[netName]
	for _, pProp := range pTable.GetProps() {
		pFlags := pProp.GetFlags()
		pType := pProp.GetType()
		if pFlags&int32(SPROP_EXCLUDE) != 0 ||
			pType == int32(DPT_Array) ||
			sth.hasExcludedSendProp(netName, pProp.GetVarName()) {
			continue
		}
		if pType == int32(DPT_DataTable) {
			if pFlags&int32(SPROP_COLLAPSIBLE) != 0 {
				sth.buildHierarchyIterateProps(pProp.GetDtName(), result)
			} else {
				sth.buildHierarchy(pProp.GetDtName())
			}
		} else {
			*result = append(*result, NewSendProp(pProp, netName))
		}
	}
}

func (sth *Helper) hasExcludedSendProp(netName string, propName string) bool {
	for _, p := range sth.excludedSendProp {
		if p.NetName == netName && p.Name == propName {
			return true
		}
	}
	return false
}

func (sth *Helper) sortByPriority() {
	priorities := []int{}
	has64 := false
	for _, prop := range sth.flatSendTable {
		priority := int(prop.Priority)
		unique := true
		for _, existing := range priorities {
			if priority == existing {
				unique = false
				break
			}
		}
		if unique {
			has64 = has64 || priority == 64
			priorities = append(priorities, priority)
		}
	}
	if !has64 {
		priorities = append(priorities, 64)
	}

	sort.Ints(priorities)

	start := 0
	for _, priority := range priorities {
		i := 0
		for {
			for i = start; i < len(sth.flatSendTable); i++ {
				p := sth.flatSendTable[i]
				if priority == p.Priority || ((p.Flags&SPROP_CHANGES_OFTEN) == SPROP_CHANGES_OFTEN) && priority == 64 {
					if i != start {
						sth.flatSendTable[i] = sth.flatSendTable[start]
						sth.flatSendTable[start] = p
					}
					start++
					break
				}
			}
			if i == len(sth.flatSendTable) {
				break
			}
		}
	}
}
