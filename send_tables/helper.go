package send_tables

import (
	"sort"

	"github.com/davecgh/go-spew/spew"
	dota "github.com/elobuff/d2rp/dota"
)

func foo() { spew.Dump("hi") }

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

type SendProp struct {
	DtName    string
	VarName   string
	Type      DPTType
	Flags     Flag
	Priority  int
	NumBits   int
	LowValue  float64
	HighValue float64
}

type Helper struct {
	sendTables       map[string]*dota.CSVCMsg_SendTable
	flatSendTable    []*SendProp
	excludedSendProp []*SendProp
}

func NewHelper(sendTables map[string]*dota.CSVCMsg_SendTable) *Helper {
	return &Helper{
		sendTables:       sendTables,
		flatSendTable:    []*SendProp{},
		excludedSendProp: []*SendProp{},
	}
}

func (sth *Helper) LoadSendTable(sendTableName string) []*SendProp {
	sth.flatSendTable = []*SendProp{}
	sth.excludedSendProp = []*SendProp{}
	sth.excludedSendProp = sth.getPropsExcluded(sendTableName)
	sth.buildHierarchy(sendTableName)
	sth.sortByPriority()
	return sth.flatSendTable
}

func (sth *Helper) getPropsExcluded(sendTableName string) []*SendProp {
	result := []*SendProp{}
	sendTable := sth.sendTables[sendTableName]
	for _, prop := range sendTable.GetProps() {
		flags := prop.GetFlags()
		if (flags & int32(SPROP_EXCLUDE)) != 0 {
			result = append(result, &SendProp{
				Flags:    Flag(flags),
				DtName:   prop.GetDtName(),
				NumBits:  int(prop.GetNumBits()),
				Priority: int(prop.GetPriority()),
				Type:     DPTType(prop.GetType()),
				VarName:  prop.GetVarName(),
			})
		}
	}

	for _, prop := range sendTable.GetProps() {
		if prop.GetType() == 6 {
			result = append(result, sth.getPropsExcluded(prop.GetDtName())...)
		}
	}

	return result
}

func (sth *Helper) buildHierarchy(sendTableName string) {
	result := []*SendProp{}
	sth.buildHierarchyIterateProps(sendTableName, &result)
	sth.flatSendTable = append(sth.flatSendTable, result...)
}
func (sth *Helper) buildHierarchyIterateProps(sendTableName string, result *[]*SendProp) {
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
				sth.buildHierarchyIterateProps(pProp.GetDtName(), result)
			} else {
				sth.buildHierarchy(pProp.GetDtName())
			}
		} else {
			*result = append(*result, &SendProp{
				DtName:    sendTableName,
				Flags:     Flag(pProp.GetFlags()),
				NumBits:   int(pProp.GetNumBits()),
				Priority:  int(pProp.GetPriority()),
				Type:      DPTType(pProp.GetType()),
				VarName:   pProp.GetVarName(),
				LowValue:  float64(pProp.GetLowValue()),
				HighValue: float64(pProp.GetHighValue()),
			})
		}
	}
}

func (sth *Helper) hasExcludedSendProp(sendTableName string, pVarName string) bool {
	for _, p := range sth.excludedSendProp {
		if p.DtName == sendTableName && p.VarName == pVarName {
			return true
		}
	}
	return false
}

func (sth *Helper) sortByPriority() {
	priorities := []int{}
	has64 := false
	for _, prop := range sth.flatSendTable {
		priority := prop.Priority
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
