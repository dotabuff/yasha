package send_tables

import (
	"io/ioutil"
	"sort"

	"github.com/davecgh/go-spew/spew"
	dota "github.com/elobuff/d2rp/dota"
)

type DPTType int

const (
	DPT_Int DPTType = iota
	DPT_Float
	DPT_Vector
	DPT_VectorXY
	DPT_String
	DPT_Array
	DPT_DataTable
	DPT_Int64
	DPT_NUMSendPropTypes
)

type Flag int

const (
	SPROP_UNSIGNED                  Flag = 1 << 0
	SPROP_COORD                     Flag = 1 << 1
	SPROP_NOSCALE                   Flag = 1 << 2
	SPROP_ROUNDDOWN                 Flag = 1 << 3
	SPROP_ROUNDUP                   Flag = 1 << 4
	SPROP_NORMAL                    Flag = 1 << 5
	SPROP_EXCLUDE                   Flag = 1 << 6
	SPROP_XYZE                      Flag = 1 << 7
	SPROP_INSIDEARRAY               Flag = 1 << 8
	SPROP_PROXY_ALWAYS_YES          Flag = 1 << 9
	SPROP_IS_A_VECTOR_ELEM          Flag = 1 << 10
	SPROP_COLLAPSIBLE               Flag = 1 << 11
	SPROP_COORD_MP                  Flag = 1 << 12
	SPROP_COORD_MP_LOWPRECISION     Flag = 1 << 13
	SPROP_COORD_MP_INTEGRAL         Flag = 1 << 14
	SPROP_CELL_COORD                Flag = 1 << 15
	SPROP_CELL_COORD_LOWPRECISION   Flag = 1 << 16
	SPROP_CELL_COORD_INTEGRAL       Flag = 1 << 17
	SPROP_CHANGES_OFTEN             Flag = 1 << 18
	SPROP_ENCODED_AGAINST_TICKCOUNT Flag = 1 << 19
)

type Helper struct {
	sendTables       map[string]*dota.CSVCMsg_SendTable
	flatSendTable    []dota.CSVCMsg_SendTableSendpropT
	excludedSendProp []dota.CSVCMsg_SendTableSendpropT
}

func NewHelper(sendTables map[string]*dota.CSVCMsg_SendTable) *Helper {
	return &Helper{
		sendTables:       sendTables,
		flatSendTable:    []dota.CSVCMsg_SendTableSendpropT{},
		excludedSendProp: []dota.CSVCMsg_SendTableSendpropT{},
	}
}

func (sth *Helper) LoadSendTable(sendTableName string) []dota.CSVCMsg_SendTableSendpropT {
	sth.flatSendTable = []dota.CSVCMsg_SendTableSendpropT{}
	sth.excludedSendProp = []dota.CSVCMsg_SendTableSendpropT{}
	sth.excludedSendProp = sth.sendTableGetPropsExcluded(sendTableName)
	sth.sendTableBuildHierarchy(sendTableName)
	sth.sendTableSortByPriority()
	return sth.flatSendTable
}

func (sth *Helper) DumpSendTable(sendTableName string, filePath string) {
	result := sth.LoadSendTable(sendTableName)
	s := spew.Sdump(result)
	ioutil.WriteFile(filePath, []byte(s), 0644)
}

func (sth *Helper) flagString(flag Flag) (result string) {
	panic("no easy way to do that in go.")
}

func (sth *Helper) sendTableGetPropsExcluded(sendTableName string) []dota.CSVCMsg_SendTableSendpropT {
	result := []dota.CSVCMsg_SendTableSendpropT{}
	sendTable := sth.sendTables[sendTableName]
	for _, prop := range sendTable.GetProps() {
		flags := prop.GetFlags()
		if (flags & int32(SPROP_EXCLUDE)) != 0 {
			s := dota.CSVCMsg_SendTableSendpropT{
				Flags:    &flags,
				DtName:   prop.DtName,
				NumBits:  prop.NumBits,
				Priority: prop.Priority,
				Type:     prop.Type,
				VarName:  prop.VarName,
			}
			result = append(result, s)
		}
		for _, prop := range sendTable.GetProps() {
			if prop.GetType() == 6 {
				result = append(result, sth.sendTableGetPropsExcluded(prop.GetDtName())...)
			}
		}
	}

	return result
}

func (sth *Helper) sendTableBuildHierarchy(sendTableName string) {
	result := []dota.CSVCMsg_SendTableSendpropT{}
	sth.sendTableBuildHierarchyIterateProps(sendTableName, result)
	for _, res := range result {
		sth.flatSendTable = append(sth.flatSendTable, res)
	}
}
func (sth *Helper) sendTableBuildHierarchyIterateProps(sendTableName string, result []dota.CSVCMsg_SendTableSendpropT) {
	pTable := sth.sendTables[sendTableName]
	for _, pProp := range pTable.GetProps() {
		pFlags := pProp.GetFlags()
		pType := pProp.GetType()
		if pFlags&int32(SPROP_EXCLUDE) != 0 ||
			pType == int32(DPT_Array) ||
			sth.hasExcludedSendProp(sendTableName, pProp.GetVarName()) {
			continue
		}
		if pType == int32(DPT_DataTable) {
			if pFlags&int32(SPROP_COLLAPSIBLE) != 0 {
				sth.sendTableBuildHierarchyIterateProps(pProp.GetDtName(), result)
			} else {
				sth.sendTableBuildHierarchy(pProp.GetDtName())
			}
		} else {
			result = append(result, dota.CSVCMsg_SendTableSendpropT{
				DtName:    &sendTableName,
				Flags:     pProp.Flags,
				NumBits:   pProp.NumBits,
				Priority:  pProp.Priority,
				Type:      pProp.Type,
				VarName:   pProp.VarName,
				LowValue:  pProp.LowValue,
				HighValue: pProp.HighValue,
			})
		}
	}
}

func (sth *Helper) hasExcludedSendProp(sendTableName string, pVarName string) bool {
	for _, p := range sth.excludedSendProp {
		if *p.DtName == sendTableName && *p.VarName == pVarName {
			return true
		}
	}
	return false
}

func (sth *Helper) sendTableSortByPriority() {
	prioritySet := map[int]bool{}
	for _, prop := range sth.flatSendTable {
		prioritySet[int(*prop.Priority)] = true
	}
	priorities := []int{}
	for k, _ := range prioritySet {
		priorities = append(priorities, k)
	}
	sort.Ints(priorities)

	start := 0
	for _, priority := range priorities {
		i := 0
		for {
			for i = start; i < len(sth.flatSendTable); i++ {
				p := sth.flatSendTable[i]
				if (priority == int(p.GetPriority())) ||
					((Flag(p.GetFlags())&SPROP_CHANGES_OFTEN) == SPROP_CHANGES_OFTEN) &&
						priority == 64 {
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
