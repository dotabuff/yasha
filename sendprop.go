package yasha

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/dotabuff/yasha/bitstream"
	"github.com/dotabuff/yasha/dota"
)

type SendPropFlag int
type SendPropType int

const (
	SPROP_UNSIGNED SendPropFlag = 1 << iota
	SPROP_COORD
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
	DPT_Int SendPropType = iota
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
	Type      SendPropType
	ArrayType *SendProp
	Value     interface{}
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
		Type:      SendPropType(prop.GetType()),
	}
}

func SendPropFromStream(bs *bitstream.BitStream, prop *SendProp) *SendProp {
	prop.update(bs)
	return prop
}

func (s *SendProp) isExcluded(tableName string, excluded []*SendProp) bool {
	name := tableName + s.Name
	for _, ex := range excluded {
		if ex.Name == name {
			return true
		}
	}
	return false
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

	fmt.Println("readfloat")

	if flags&SPROP_COORD != 0 {
		fmt.Println("readfloat->readcoord")
		return b.ReadCoord()
	}

	if flags&SPROP_COORD_MP != 0 {
		fmt.Println("readfloat->readcoordmp")
		integral := int(flags & SPROP_COORD_MP_INTEGRAL)
		lowPrecision := int(flags & SPROP_COORD_MP_LOWPRECISION)
		return b.ReadCoordMp(integral, lowPrecision)
	}

	if flags&SPROP_NOSCALE != 0 {
		fmt.Println("readfloat->noscale??")
		var f float32
		buf := bytes.NewReader(b.ReadBits(32))
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

	val := b.Read(s.Bits)
	fval := float32(val / (1<<uint(s.Bits) - 1))
	return fval*(s.HighValue-s.LowValue) + s.LowValue
}

type Vector3 struct {
	X, Y, Z float32
}

func (s *SendProp) readVector3(b *bitstream.BitStream) *Vector3 {
	vec := &Vector3{}
	fmt.Println("first float")
	vec.X = s.readFloat(b)
	fmt.Println("second float")
	vec.Y = s.readFloat(b)

	if s.Flags&SPROP_NORMAL != 0 {
		fmt.Println("normal bool")
		sign := b.ReadBool()
		f := vec.X*vec.X + vec.Y*vec.Y

		if f >= 0 {
			vec.Z = float32(math.Sqrt(1 - float64(f)))
		}

		if sign {
			vec.Z *= -1
		}
	} else {
		fmt.Println("non-normal float")
		vec.Z = s.readFloat(b)
	}

	return vec
}

type Vector2 struct {
	X, Y float32
}

func (s *SendProp) readVector2(b *bitstream.BitStream) *Vector2 {
	vec := &Vector2{}
	vec.X = s.readFloat(b)
	vec.Y = s.readFloat(b)
	return vec
}

func (s *SendProp) readString(b *bitstream.BitStream) string {
	length := b.Read(9)

	if length > PROPERTY_MAX_STRING_LENGTH {
		panic("string too long")
	}

	return string(b.ReadBits(int(8 * length)))
}

func (s *SendProp) readInt64(b *bitstream.BitStream) uint64 {
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
		negate = b.ReadBool()
	}

	x := b.Read(32)
	y := b.Read(sbits)
	val := uint64(x)<<32 | uint64(y)

	if negate {
		return -val
	}

	return val
}

func (s *SendProp) readArray(b *bitstream.BitStream) []int {
	elements := s.Elements
	bits := int(math.Floor(math.Log2(float64(elements)) + 1))
	count := b.Read(bits)
	if count > PROPERTY_MAX_ELEMENTS {
		panic("too many elements in array")
	}

	props := []int{}
	/*
		for i := uint(0); i < count; i++ {
			props = append(props, NewProperty(b, aType))
		}
	*/

	panic("ReadArray")

	return props
}

func (s *SendProp) update(bs *bitstream.BitStream) {
	fmt.Printf("sendprop %+v updating\n", s)
	switch s.Type {
	case DPT_Int:
		s.Value = s.readInt(bs)
	case DPT_Float:
		s.Value = s.readFloat(bs)
	case DPT_Vector:
		s.Value = s.readVector3(bs)
	case DPT_VectorXY:
		s.Value = s.readVector2(bs)
	case DPT_String:
		s.Value = s.readString(bs)
	case DPT_Array:
		s.Value = s.readArray(bs)
	case DPT_Int64:
		s.Value = s.readInt64(bs)
	default:
		panic(fmt.Errorf("unknown type: %d", s.Type))
	}

	fmt.Printf("sendprop %s: %+v\n", s.Name, s.Value)
}
