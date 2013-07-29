package utils

import (
	"bytes"
	"encoding/binary"
	"math"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/elobuff/d2rp/core"
	"github.com/elobuff/d2rp/core/send_tables"
)

var debug = false

func dump(args ...interface{}) {
	if debug {
		spew.Dump(args...)
	}
}
func printLn(args ...interface{}) {
	if debug {
		spew.Println(args...)
	}
}

const (
	CoordIntegerBits    = 14
	CoordFractionalBits = 5
	CoordDenominator    = (1 << CoordFractionalBits)
	CoordResolution     = (1.0 / CoordDenominator)

	NormalFractionalBits = 11
	NormalDenominator    = ((1 << NormalFractionalBits) - 1)
	NormalResolution     = (1.0 / NormalDenominator)
)

type BitReader struct {
	buffer     []byte
	currentBit int
}

func NewBitReader(buffer []byte) *BitReader {
	if len(buffer) == 0 {
		panic("empty buffer?")
	}
	return &BitReader{buffer: buffer}
}

func (br *BitReader) Length() int      { return len(br.buffer) }
func (br *BitReader) CurrentBit() int  { return br.currentBit }
func (br *BitReader) CurrentByte() int { return br.currentBit / 8 }
func (br *BitReader) BitsLeft() int    { return (len(br.buffer) * 8) - br.currentBit }
func (br *BitReader) BytesLeft() int   { return len(br.buffer) - (br.currentBit * 8) }

type SeekOrigin int

const (
	Current SeekOrigin = iota
	Begin
	End
)

func (br *BitReader) SeekBits(offset int, origin SeekOrigin) {
	if origin == Current {
		br.currentBit += offset
	} else if origin == Begin {
		br.currentBit = offset
	} else if origin == End {
		br.currentBit = (len(br.buffer) * 8) - offset
	}
	if br.currentBit < 0 || br.currentBit > (len(br.buffer)*8) {
		panic("out of range")
	}
}

func (br *BitReader) ReadUBitsByteAligned(nBits int) uint {
	if nBits%8 != 0 {
		panic("Must be multple of 8")
	}

	if br.currentBit%8 != 0 {
		panic("Current bit is not byte-aligned")
	}

	var result uint
	for i := 0; i < nBits/8; i++ {
		result += uint(br.buffer[br.CurrentByte()] << (uint(i) * 8))
		br.currentBit += 8
	}
	return result
}

func (br *BitReader) ReadUBitsNotByteAligned(nBits int) uint {
	bitOffset := br.currentBit % 8
	nBitsToRead := bitOffset + nBits
	nBytesToRead := nBitsToRead / 8
	if nBitsToRead%8 != 0 {
		nBytesToRead += 1
	}

	var currentValue uint64
	for i := 0; i < nBytesToRead; i++ {
		b := br.buffer[br.CurrentByte()+i]
		currentValue += (uint64(b) << (uint64(i) * 8))
	}
	currentValue >>= uint(bitOffset)
	currentValue &= ((1 << uint64(nBits)) - 1)
	br.currentBit += nBits
	return uint(currentValue)
}

func (br *BitReader) ReadVarInt() uint {
	var b uint
	var count int
	var result uint
	for {
		if count == 5 {
			return result
		} else if br.CurrentByte() >= len(br.buffer) {
			return result
		}
		b = br.ReadUBits(8)
		result |= (b & 0x7f) << uint(7*count)
		count++
		if (b & 0x80) != 0x80 {
			break
		}
	}
	return result
}

func (br *BitReader) ReadUBits(nBits int) uint {
	if nBits <= 0 || nBits > 32 {
		panic("Value must be a positive integer between 1 and 32 inclusive.")
	}
	if (br.currentBit + nBits) > (len(br.buffer) * 8) {
		panic("Out of range")
	}
	if br.currentBit%8 == 0 && nBits%8 == 0 {
		return br.ReadUBitsByteAligned(nBits)
	}
	return br.ReadUBitsNotByteAligned(nBits)
}

func (br *BitReader) ReadBits(nBits int) int {
	result := br.ReadUBits(nBits - 1)
	if br.ReadBoolean() {
		result = -((1 << (uint(nBits) - 1)) - result)
	}
	return int(result)
}

func (br *BitReader) ReadBoolean() bool {
	if br.CurrentBit()+1 > br.Length()*8 {
		panic("Out of range")
	}
	currentByte := br.currentBit / 8
	bitOffset := br.currentBit % 8
	result := br.buffer[currentByte]&(1<<uint(bitOffset)) != 0
	br.currentBit++
	return result
}

func (br *BitReader) ReadByte() byte {
	return byte(br.ReadUBits(8))
}

func (br *BitReader) ReadSByte() int8 {
	return int8(br.ReadBits(8))
}

func (br *BitReader) ReadBytes(nBytes int) []byte {
	if nBytes <= 0 {
		panic("Must be positive integer: nBytes")
	}
	result := make([]byte, nBytes)
	for i := 0; i < nBytes; i++ {
		result[i] = br.ReadByte()
	}
	return result
}

func (br *BitReader) ReadBitFloat() float32 {
	b := bytes.NewBuffer(br.ReadBytes(4))
	var f float32
	binary.Read(b, binary.LittleEndian, &f)
	return f
}

func (br *BitReader) ReadBitNormal() float32 {
	signbit := br.ReadBoolean()
	fractval := float32(br.ReadUBits(NormalFractionalBits))
	value := fractval * NormalResolution
	if signbit {
		value = -value
	}
	return value
}

func (br *BitReader) ReadBitCellCoord(bits int, integral, lowPrecision bool) (value float32) {
	if integral {
		value = float32(br.ReadBits(bits))
	} else {
		intval := br.ReadBits(bits)
		if lowPrecision {
			fractval := float32(br.ReadBits(3))
			value = float32(intval) + (fractval * (1.0 / (1 << 3)))
		} else {
			fractval := float32(br.ReadBits(5))
			value = float32(intval) + (fractval * (1.0 / (1 << 5)))
		}
	}
	return value
}

func (br *BitReader) ReadBitCoord() (value float32) {
	intFlag := br.ReadBoolean()
	fractFlag := br.ReadBoolean()
	if intFlag || fractFlag {
		negative := br.ReadBoolean()
		if intFlag {
			value += float32(br.ReadUBits(CoordIntegerBits)) + 1
		}
		if fractFlag {
			value += float32(br.ReadUBits(CoordFractionalBits)) * CoordResolution
		}
		if negative {
			value = -value
		}
	}
	return value
}

func (br *BitReader) ReadString() string {
	bs := []byte{}
	for {
		b := br.ReadByte()
		if b == 0 {
			break
		}
		bs = append(bs, b)
	}
	return string(bs)
}

func (br *BitReader) ReadFloat(prop *send_tables.SendProp) float32 {
	if result, ok := br.ReadSpecialFloat(prop); ok {
		return result
	}
	dwInterp := uint(br.ReadUBits(prop.NumBits))
	result := float32(dwInterp) / float32((uint(1)<<uint(prop.NumBits))-1)
	result = float32(prop.LowValue+(prop.HighValue-prop.LowValue)) * result
	return result
}

func (br *BitReader) ReadLengthPrefixedString() string {
	stringLength := uint(br.ReadUBits(9))
	if stringLength > 0 {
		return string(br.ReadBytes(int(stringLength)))
	}
	return ""
}

func (br *BitReader) ReadVector(prop *send_tables.SendProp) *core.Vector {
	result := &core.Vector{
		X: br.ReadFloat(prop),
		Y: br.ReadFloat(prop),
	}
	if prop.Flags&send_tables.SPROP_NORMAL == 0 {
		result.Z = br.ReadFloat(prop)
	} else {
		signbit := br.ReadBoolean()
		v0v0v1v1 := float64(result.X*result.X + result.Y*result.Y)
		if v0v0v1v1 < 1.0 {
			result.Z = float32(math.Sqrt(1.0 - v0v0v1v1))
		} else {
			result.Z = 0.0
		}
		if signbit {
			result.Z *= -1.0
		}
	}

	return result
}

func (br *BitReader) ReadVectorXY(prop *send_tables.SendProp) *core.Vector {
	return &core.Vector{
		X: br.ReadFloat(prop),
		Y: br.ReadFloat(prop),
	}
}

func (br *BitReader) ReadInt(prop *send_tables.SendProp) int {
	if prop.Flags&send_tables.SPROP_UNSIGNED != 0 {
		return int(br.ReadUBits(prop.NumBits))
	}
	return br.ReadBits(prop.NumBits)
}

func (br *BitReader) ReadSpecialFloat(prop *send_tables.SendProp) (float32, bool) {
	if prop.Flags&send_tables.SPROP_COORD != 0 {
		return br.ReadBitCoord(), true
	} else if prop.Flags&send_tables.SPROP_COORD_MP != 0 {
		panic("wtf")
	} else if prop.Flags&send_tables.SPROP_COORD_MP_INTEGRAL != 0 {
		panic("wtf")
	} else if prop.Flags&send_tables.SPROP_COORD_MP_LOWPRECISION != 0 {
		panic("wtf")
	} else if prop.Flags&send_tables.SPROP_CELL_COORD != 0 {
		return br.ReadBitCellCoord(prop.NumBits, false, false), true
	} else if prop.Flags&send_tables.SPROP_CELL_COORD_INTEGRAL != 0 {
		return br.ReadBitCellCoord(prop.NumBits, true, false), true
	} else if prop.Flags&send_tables.SPROP_CELL_COORD_LOWPRECISION != 0 {
		return br.ReadBitCellCoord(prop.NumBits, false, true), true
	} else if prop.Flags&send_tables.SPROP_NOSCALE != 0 {
		return br.ReadBitFloat(), true
	} else if prop.Flags&send_tables.SPROP_NORMAL != 0 {
		return br.ReadBitNormal(), true
	}
	return 0, false
}

func (br *BitReader) ReadInt64(prop *send_tables.SendProp) uint64 {
	var low, high uint
	if prop.Flags&send_tables.SPROP_UNSIGNED != 0 {
		low = br.ReadUBits(32)
		high = br.ReadUBits(32)
	} else {
		br.SeekBits(1, Current)
		low = br.ReadUBits(32)
		high = br.ReadUBits(31)
	}
	res := uint64(high)
	res = (res << 32)
	res = res | uint64(low)
	return res
}

func (br *BitReader) ReadNextEntityIndex(oldEntity int) int {
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

func (br *BitReader) ReadPropertiesIndex() []int {
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
			prop += 1
			prop += int(value)
			props = append(props, prop)
		}
	}
	return props
}

func (br *BitReader) ReadPropertiesValues(mapping []*send_tables.SendProp, multiples map[string]int, indices []int) map[string]interface{} {
	values := map[string]interface{}{}
	debug = true

	for _, index := range indices {
		prop := mapping[index]
		name := prop.DtName + "." + prop.VarName
		multiple := multiples[name] > 1
		elements := 1
		if prop.Flags&send_tables.SPROP_INSIDEARRAY != 0 {
			elements = int(br.ReadUBits(6))
		}
		for k := 0; k < elements; k++ {
			key := name
			if multiple {
				key += "-" + strconv.Itoa(index)
			}
			if elements > 1 {
				key += "-" + strconv.Itoa(k)
			}
			switch prop.Type {
			case send_tables.DPT_Int:
				if (prop.Flags & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					values[key] = br.ReadVarInt()
				} else {
					values[key] = br.ReadInt(prop)
				}
			case send_tables.DPT_Float:
				if (prop.Flags & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = br.ReadFloat(prop)
				}
			case send_tables.DPT_Vector:
				if (prop.Flags & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = br.ReadVector(prop)
				}
			case send_tables.DPT_VectorXY:
				if (prop.Flags & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = br.ReadVectorXY(prop)
				}
			case send_tables.DPT_String:
				if (prop.Flags & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = br.ReadLengthPrefixedString()
				}
			case send_tables.DPT_Int64:
				if (prop.Flags & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = br.ReadInt64(prop)
				}
			default:
				panic("unknown type")
			}
		}
	}

	debug = false
	return values
}
