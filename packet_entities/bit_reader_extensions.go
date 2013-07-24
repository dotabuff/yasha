package packet_entities

import (
	"strconv"

	"github.com/elobuff/d2rp/core/send_tables"
	"github.com/elobuff/d2rp/core/utils"
	dota "github.com/elobuff/d2rp/dota"
)

func flag(prop dota.CSVCMsg_SendTableSendpropT) send_tables.Flag {
	return send_tables.Flag(prop.GetFlags())
}

func ReadNextEntityIndex(br *utils.BitReader, oldEntity int) int {
	ret := br.ReadUBits(4)
	more1 := br.ReadBoolean()
	more2 := br.ReadBoolean()
	if more1 {
		ret += (br.ReadUBits(4) << 4)
	}
	if more2 {
		ret += (br.ReadUBits(8) << 4)
	}
	return oldEntity + 1 + int(ret)
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

func ReadPropertiesIndex(br *utils.BitReader) []int {
	props := []int{}
	prop := -1
	for {
		if br.ReadBoolean() {
			prop += 1
			props = append(props, prop)
		} else {
			value := br.ReadVarInt()
			if value == 16383 {
				break
			}
			prop += (int(value) + 1)
			props = append(props, prop)
		}
	}
	return props
}

func ReadPropertiesValues(br *utils.BitReader, mapping []dota.CSVCMsg_SendTableSendpropT, multiples map[string]int, indices []int) map[string]interface{} {
	values := map[string]interface{}{}
	for _, index := range indices {
		prop := mapping[index]
		multiple := multiples[prop.GetDtName()+"."+prop.GetVarName()] > 1
		elements := 1
		if (flag(prop) & send_tables.SPROP_INSIDEARRAY) != 0 {
			elements = int(br.ReadUBits(6))
		}
		for k := 0; k < elements; k++ {
			key := prop.GetDtName() + "." + prop.GetVarName()
			if multiple {
				key += ("-" + strconv.Itoa(index))
			}
			if elements > 1 {
				key += ("-" + strconv.Itoa(k))
			}
			switch send_tables.DPTType(prop.GetType()) {
			case send_tables.DPT_Int:
				if (flag(prop) & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					values[key] = br.ReadVarInt()
				} else {
					values[key] = br.ReadInt(prop)
				}
			case send_tables.DPT_Float:
				if (flag(prop) & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = br.ReadFloat(prop)
				}
			case send_tables.DPT_Vector:
				if (flag(prop) & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = br.ReadVector(prop)
				}
			case send_tables.DPT_VectorXY:
				if (flag(prop) & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = br.ReadVectorXY(prop)
				}
			case send_tables.DPT_String:
				if (flag(prop) & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = br.ReadLengthPrefixedString()
				}
			case send_tables.DPT_Int64:
				if (flag(prop) & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = br.ReadInt64(prop)
				}
			}
		}
	}
	return values
}
