package packet_entities

import (
	"math"
	"strconv"

	"github.com/elobuff/d2rp/core"
	"github.com/elobuff/d2rp/core/send_tables"
	"github.com/elobuff/d2rp/core/utils"
	dota "github.com/elobuff/d2rp/dota"
)

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
					values[key] = ReadInt(br, prop)
				}
			case send_tables.DPT_Float:
				if (flag(prop) & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = ReadFloat(br, prop)
				}
			case send_tables.DPT_Vector:
				if (flag(prop) & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = ReadVector(br, prop)
				}
			case send_tables.DPT_VectorXY:
				if (flag(prop) & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = ReadVectorXY(br, prop)
				}
			case send_tables.DPT_String:
				if (flag(prop) & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = ReadLengthPrefixedString(br)
				}
			case send_tables.DPT_Int64:
				if (flag(prop) & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = ReadInt64(br, prop)
				}
			}
		}
	}
	return values
}

func ReadVector(br *utils.BitReader, prop dota.CSVCMsg_SendTableSendpropT) *core.Vector {
	result := &core.Vector{
		X: ReadFloat(br, prop),
		Y: ReadFloat(br, prop),
	}
	if (flag(prop) & send_tables.SPROP_NORMAL) == 0 {
		result.Z = ReadFloat(br, prop)
	} else {
		signbit := br.ReadBoolean()
		v0v0v1v1 := result.X*result.X + result.Y*result.Y
		if v0v0v1v1 < 1.0 {
			result.Z = math.Sqrt(1.0 - v0v0v1v1)
		} else {
			result.Z = 0.0
		}
		if signbit {
			result.Z *= -1.0
		}
	}

	return result
}

func ReadVectorXY(br *utils.BitReader, prop dota.CSVCMsg_SendTableSendpropT) *core.Vector {
	return &core.Vector{
		X: ReadFloat(br, prop),
		Y: ReadFloat(br, prop),
	}
}

func ReadInt(br *utils.BitReader, prop dota.CSVCMsg_SendTableSendpropT) int {
	if (flag(prop) & send_tables.SPROP_UNSIGNED) != 0 {
		return int(br.ReadUBits(int(prop.GetNumBits())))
	}
	return br.ReadBits(int(prop.GetNumBits()))
}

func ReadSpecialFloat(br *utils.BitReader, prop dota.CSVCMsg_SendTableSendpropT) (float64, bool) {
	flags := flag(prop)
	if (flags & send_tables.SPROP_COORD) != 0 {
		return br.ReadBitCoord(), true
	} else if (flags & send_tables.SPROP_COORD_MP) != 0 {
		panic("utils.BitReader.ReadSpecialFloat")
	} else if (flags & send_tables.SPROP_COORD_MP_INTEGRAL) != 0 {
		panic("utils.BitReader.ReadSpecialFloat")
	} else if (flags & send_tables.SPROP_COORD_MP_LOWPRECISION) != 0 {
		panic("utils.BitReader.ReadSpecialFloat")
	} else if (flags & send_tables.SPROP_CELL_COORD) != 0 {
		return br.ReadBitCellCoord(int(prop.GetNumBits()), false, false), true
	} else if (flags & send_tables.SPROP_CELL_COORD_INTEGRAL) != 0 {
		return br.ReadBitCellCoord(int(prop.GetNumBits()), true, false), true
	} else if (flags & send_tables.SPROP_CELL_COORD_LOWPRECISION) != 0 {
		return br.ReadBitCellCoord(int(prop.GetNumBits()), false, true), true
	} else if (flags & send_tables.SPROP_NOSCALE) != 0 {
		return float64(br.ReadBitFloat()), true
	} else if (flags & send_tables.SPROP_NORMAL) != 0 {
		return br.ReadBitNormal(), true
	}
	return 0, false
}

func ReadFloat(br *utils.BitReader, prop dota.CSVCMsg_SendTableSendpropT) float64 {
	if result, ok := ReadSpecialFloat(br, prop); ok {
		return result
	}
	dwInterp := float64(br.ReadUBits(int(prop.GetNumBits())))
	bits := 1 << uint(prop.GetNumBits())
	result := dwInterp / float64(bits-1)
	low, high := float64(prop.GetLowValue()), float64(prop.GetHighValue())
	return low + (high-low)*result
}

func ReadLengthPrefixedString(br *utils.BitReader) string {
	stringLength := int(br.ReadUBits(9))
	if stringLength > 0 {
		return string(br.ReadBytes(stringLength))
	}
	return ""
}

func ReadInt64(br *utils.BitReader, prop dota.CSVCMsg_SendTableSendpropT) uint64 {
	var low, high uint
	if (flag(prop) & send_tables.SPROP_UNSIGNED) != 0 {
		low = br.ReadUBits(32)
		high = br.ReadUBits(32)
	} else {
		br.SeekBits(1, 0)
		low = br.ReadUBits(32)
		high = br.ReadUBits(31)
	}
	res := uint64(high)
	res = (res << 32)
	res = res | uint64(low)
	return res
}

func flag(prop dota.CSVCMsg_SendTableSendpropT) send_tables.Flag {
	return send_tables.Flag(prop.GetFlags())
}
