package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/dotabuff/yasha/send_tables"
)

var p = spew.Dump

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
func (br *BitReader) BytesLeft() int   { return len(br.buffer) - (br.currentBit / 8) }

type Vector3 struct {
	X, Y, Z float64
}

func (v Vector3) String() string {
	return fmt.Sprintf("{{ x: %f, y: %f, z: %f }}", v.X, v.Y, v.Z)
}

type Vector2 struct {
	X, Y float64
}

func (v Vector2) String() string {
	return fmt.Sprintf("{{ x: %f, y: %f }}", v.X, v.Y)
}

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

func (br *BitReader) ReadVarInt() (result uint) {
	var b uint
	count := 0

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

func (br *BitReader) ReadBitsAsBytes(n int) []byte {
	result := make([]byte, (n+7)/8)
	i := 0
	for n > 7 {
		n -= 8
		result[i] = byte(br.ReadUBits(8))
		i++
	}
	if n != 0 {
		result[i] = byte(br.ReadUBits(n))
	}

	return result
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

func (br *BitReader) ReadStringN(n int) string {
	buf := []byte{}
	for n > 0 {
		c := br.ReadByte()
		if c == 0 {
			break
		}
		buf = append(buf, c)
		n--
	}
	return string(buf)
}

func (br *BitReader) ReadFloat(prop *send_tables.SendProp) float64 {
	if result, ok := br.ReadSpecialFloat(prop); ok {
		return float64(result)
	}
	dividend := float64(br.ReadUBits(prop.NumBits))
	divisor := (1 << uint(prop.NumBits)) - 1
	f := dividend / float64(divisor)
	r := prop.HighValue - prop.LowValue
	return f*r + prop.LowValue
}

func (br *BitReader) ReadLengthPrefixedString() string {
	stringLength := uint(br.ReadUBits(9))

	if stringLength > 0 {
		return string(br.ReadBytes(int(stringLength)))
	}
	return ""
}

func (br *BitReader) ReadVector(prop *send_tables.SendProp) *Vector3 {
	var x, y, z float64
	x = br.ReadFloat(prop)
	y = br.ReadFloat(prop)

	if prop.Flags&send_tables.SPROP_NORMAL == 0 {
		z = br.ReadFloat(prop)
	} else {
		f := float64(x*x + y*y)
		if 1.0 >= f {
			z = 0
		} else {
			z = float64(math.Sqrt(1.0 - f))
		}
		if signbit := br.ReadBoolean(); signbit {
			z = -z
		}
	}

	return &Vector3{
		X: x, Y: y, Z: z,
	}
}

func (br *BitReader) ReadVectorXY(prop *send_tables.SendProp) *Vector2 {
	return &Vector2{
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

func (br *BitReader) nReadVarUnsignedInt64() uint64 {
	var readCount, tmpBuf uint32
	var value uint64

	for {
		if readCount == VarInt64Max {
			return value
		}

		tmpBuf = uint32(br.ReadBits(8))
		value |= uint64(tmpBuf&0x7F) << (7 * readCount)
		readCount += 1
		if tmpBuf&0x80 != 0x80 {
			break
		}
	}
	return value
}

func (br *BitReader) nReadVarSignedInt64() int64 {
	value := br.nReadVarUnsignedInt64()
	return int64(value>>1) ^ -int64(value&1)
}

func (br *BitReader) ReadInt64(prop *send_tables.SendProp) int64 {
	if prop.Flags&send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT != 0 {
		if prop.Flags&send_tables.SPROP_UNSIGNED != 0 {
			return int64(br.nReadVarUnsignedInt64())
		} else {
			return br.nReadVarSignedInt64()
		}
	} else {
		negate := false
		sbits := prop.NumBits - 32

		if prop.Flags&send_tables.SPROP_UNSIGNED == 0 {
			sbits -= 1
			negate = br.ReadBoolean()
		}

		a := int64(br.ReadUBits(32))
		b := int64(br.ReadUBits(sbits))
		val := (b << 32) | a
		if negate {
			return -val
		} else {
			return val
		}
	}
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

// FIXME:
// Rememver that our PE decoder is wrong, it needs to construct key names recursively.
// https://github.com/spheenik/clarity/tree/master/src/main/java/clarity/decoder/SendTableFlattener.java
// https://gist.githubusercontent.com/onethirtyfive/07899a78622dc18679c3/raw/19d411910016170e4c4ee2782fd4a987e9ce2afc/gistfile1.txt
func (br *BitReader) ReadPropertiesValues(mapping []*send_tables.SendProp, multiples map[string]int, indices []int) map[string]interface{} {
	values := map[string]interface{}{}

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

			var value interface{}

			switch prop.Type {
			case send_tables.DPT_Int:
				value = br.ReadInt(prop)
			case send_tables.DPT_Int64:
				value = br.ReadInt64(prop)
			case send_tables.DPT_Float:
				value = br.ReadFloat(prop)
			case send_tables.DPT_Vector:
				value = br.ReadVector(prop)
			case send_tables.DPT_VectorXY:
				value = br.ReadVectorXY(prop)
			case send_tables.DPT_String:
				value = br.ReadLengthPrefixedString()
			default:
				panic("unknown type")
			}

			p(prop.Type, key, value)
			values[key] = value
		}
	}

	return values
}
