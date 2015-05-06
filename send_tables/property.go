package send_tables

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/dotabuff/yasha/bitstream"
	"github.com/dotabuff/yasha/dota"
)

type SendPropFlag int
type SendPropType int

const (
	SPROP_UNSIGNED SendPropFlag = 1 << 0
	SPROP_COORD    SendPropFlag = 1 << iota
	SPROP_NOSCALE
	SPROP_ROUNDDOWN
	SPROP_ROUNDUP
	SPROP_NORMAL
	SPROP_EXCLUDE
	SPROP_XYZE
	SPROP_INSIDEARRAY
	SPROP_PROXY_ALWAYS_YES
	SPROP_IS_A_VECTOR_ELEM
	SPROP_COLLAPSIBLE
	SPROP_COORD_MP
	SPROP_COORD_MP_LOWPRECISION
	SPROP_COORD_MP_INTEGRAL
	SPROP_CELL_COORD
	SPROP_CELL_COORD_LOWPRECISION
	SPROP_CELL_COORD_INTEGRAL
	SPROP_CHANGES_OFTEN
	SPROP_ENCODED_AGAINST_TICKCOUNT
)

const (
	DPT_Int SendPropFlag = iota
	DPT_Float
	DPT_Vector
	DPT_VectorXY
	DPT_String
	DPT_Array
	DPT_DataTable
	DPT_Int64
	DPT_NUMSendPropTypes
)

const (
	PROPERTY_MAX_STRING_LENGTH = 0x200
	PROPERTY_MAX_ELEMENTS      = 100
)

type SendProp struct {
	Name      string
	NetName   string
	Flags     SendPropFlag
	Priority  int
	ClassName string
	Elements  int
	LowValue  float32
	HighValue float32
	Bits      int
	ArrayType *SendProp
}

func NewSendProp(prop *dota.CSVCMsg_SendTableSendpropT, netName string) *SendProp {
	return &SendProp{
		Name:      prop.GetVarName(),
		NetName:   netName,
		Flags:     SendPropFlag(prop.GetFlags()),
		Priority:  int(prop.GetPriority()),
		ClassName: prop.GetDtName(),
		Elements:  int(prop.GetNumElements()),
		LowValue:  prop.GetLowValue(),
		HighValue: prop.GetHighValue(),
		Bits:      int(prop.GetNumBits()),
	}
}

func (s *SendProp) readInt(b *bitstream.BitStream) uint {
	flags := s.Flags

	if flags&SPROP_ENCODED_AGAINST_TICKCOUNT != 0 {
		if flags&SPROP_UNSIGNED != 0 {
			return b.ReadVarUInt32()
		} else {
			return b.ReadVarSInt32()
		}
	} else if flags&SPROP_UNSIGNED != 0 {
		return b.ReadUInt(s.Bits)
	} else {
		return b.ReadSInt(s.Bits)
	}
}

func (s *SendProp) readFloat(b *bitstream.BitStream) float32 {
	flags := s.Flags

	if flags&SPROP_COORD != 0 {
		return b.ReadCoord()
	}

	if flags&SPROP_COORD_MP != 0 {
		integral := int(flags & SPROP_COORD_MP_INTEGRAL)
		lowPrecision := int(flags & SPROP_COORD_MP_LOWPRECISION)
		return b.ReadCoordMp(integral, lowPrecision)
	}

	if flags&SPROP_NOSCALE != 0 {
		var f float32
		buf := bytes.NewReader(b.ReadBits(32 * 8))
		err := binary.Read(buf, binary.LittleEndian, &f) // FIXME: let's hope it's LE
		if err != nil {
			panic(err)
		}
		return f
	}

	if flags&SPROP_NORMAL != 0 {
		return b.ReadNormal()
	}

	if flags&SPROP_CELL_COORD != 0 || flags&SPROP_CELL_COORD_INTEGRAL != 0 || flags&SPROP_CELL_COORD_LOWPRECISION != 0 {
		lowPrecision := flags&SPROP_CELL_COORD_LOWPRECISION != 0
		integral := flags&SPROP_COORD_MP_INTEGRAL != 0

		return b.ReadCellCoord(s.Bits, integral, lowPrecision)
	}

	dividend := float32(b.Read(s.Bits))
	divisor := float32(b.Read(int(1<<uint(s.Bits) - 1)))

	low, high := float32(s.LowValue), float32(s.HighValue)
	return (dividend/divisor)*high - low + low
}

type Vector3 struct {
	X, Y, Z float32
}

func (s *SendProp) ReadVector3(b *bitstream.BitStream) *Vector3 {
	vec := &Vector3{}
	vec.X = s.readFloat(b)
	vec.Y = s.readFloat(b)

	if s.Flags&SPROP_NORMAL != 0 {
		sign := b.Read(1)
		f := vec.X*vec.X + vec.Y*vec.Y

		if f >= 0 {
			vec.Z = float32(math.Sqrt(1 - float64(f)))
		}

		if sign != 0 {
			vec.Z *= -1
		}
	} else {
		vec.Z = s.readFloat(b)
	}

	return vec
}

type Vector2 struct {
	X, Y float32
}

func (s *SendProp) ReadVector2(b *bitstream.BitStream) *Vector2 {
	vec := &Vector2{}
	vec.X = s.readFloat(b)
	vec.Y = s.readFloat(b)
	return vec
}

func (s *SendProp) ReadString(b *bitstream.BitStream) string {
	length := b.Read(9)

	if length > PROPERTY_MAX_STRING_LENGTH {
		panic("string too long")
	}

	return string(b.ReadBits(int(8 * length)))
}

func (s *SendProp) ReadInt64(b *bitstream.BitStream) uint64 {
	flags := s.Flags

	if flags&SPROP_ENCODED_AGAINST_TICKCOUNT != 0 {
		if flags&SPROP_UNSIGNED != 0 {
			return b.ReadVarUInt64()
		} else {
			return b.ReadVarSInt64()
		}
	}

	negate := false
	sbits := s.Bits - 32
	if flags&SPROP_UNSIGNED == 0 {
		sbits--
		negate = b.Read(1) == 1
	}

	x := b.Read(32)
	y := b.Read(sbits)
	val := uint64(x)<<32 | uint64(y)

	if negate {
		return -val
	}

	return val
}

func (s *SendProp) ReadArray(b *bitstream.BitStream) []int {
	elements := s.Elements
	bits := int(math.Floor(math.Log2(float64(elements))))
	count := b.Read(bits)
	if count > PROPERTY_MAX_ELEMENTS {
		panic("too many elements in array")
	}

	aType := s.ArrayType
	props := []int{}
	for i := 0; i < count; i++ {
		props = append(props, NewProperty(b, aType))
	}

	return props
}
