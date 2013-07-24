package utils

import (
	"bytes"
	"encoding/binary"
)

const (
	CoordIntegerBits            = 14
	CoordFractionalBits         = 5
	CoordDenominator            = (1 << CoordFractionalBits)
	CoordResolution     float64 = (1.0 / CoordDenominator)

	NormalFractionalBits         = 11
	NormalDenominator            = ((1 << NormalFractionalBits) - 1)
	NormalResolution     float64 = (1.0 / NormalDenominator)
)

type BitReader struct {
	Buffer     []byte
	currentBit int
}

func (br *BitReader) Length() int      { return len(br.Buffer) }
func (br *BitReader) CurrentBit() int  { return br.currentBit }
func (br *BitReader) CurrentByte() int { return br.currentBit / 8 }
func (br *BitReader) BitsLeft() int    { return (len(br.Buffer) * 8) - br.currentBit }
func (br *BitReader) BytesLeft() int   { return len(br.Buffer) - br.CurrentByte() }

// origin -1: current
// origin  0: begin
// origin  1: end
func (br *BitReader) SeekBits(offset int, origin int) {
	if origin == -1 {
		br.currentBit += offset
	} else if origin == 0 {
		br.currentBit = offset
	} else if origin == 1 {
		br.currentBit = (br.Length() * 8) - offset
	}
	if br.currentBit < 0 || br.currentBit > (br.Length()*8) {
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
		result += uint(br.Buffer[br.CurrentByte()] << (uint(i) * 8))
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
		b := br.Buffer[br.CurrentByte()+1]
		currentValue += (uint64(b) << uint(i*8))
	}
	currentValue >>= uint(bitOffset)
	currentValue &= ((1 << uint(nBits)) - 1)
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
		} else if br.CurrentByte() >= br.Length() {
			return result
		}
		b = br.ReadUBits(8)
		result |= (b & 0x7f) << uint(7*count)
		count++
		if (b & 0x80) == 0x80 {
			break
		}
	}
	return result
}
func (br *BitReader) ReadUBits(nBits int) uint {
	if nBits <= 0 || nBits > 32 {
		panic("Value must be a positive integer between 1 and 32 inclusive.")
	}
	if br.CurrentBit()+nBits > br.Length()*8 {
		panic("Out of range")
	}
	if br.CurrentBit()%8 == 0 && nBits%8 == 0 {
		return br.ReadUBitsByteAligned(nBits)
	}
	return br.ReadUBitsNotByteAligned(nBits)
}
func (br *BitReader) ReadBits(nBits int) int {
	result := br.ReadUBits(nBits - 1)
	if br.ReadBoolean() {
		result = -((1 << uint(nBits-1)) - result)
	}
	return int(result)
}
func (br *BitReader) ReadBoolean() bool {
	if br.CurrentBit()+1 > br.Length()*8 {
		panic("Out of range")
	}
	currentByte := br.CurrentBit() / 8
	bitOffset := br.CurrentBit() % 8
	result := br.Buffer[currentByte]&(1<<uint(bitOffset)) != 0
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
func (br *BitReader) ReadBitNormal() float64 {
	signbit := br.ReadBoolean()
	fractval := float64(br.ReadUBits(NormalFractionalBits))
	value := fractval * NormalResolution
	if signbit {
		value = -value
	}
	return value
}
func (br *BitReader) ReadBitCellCoord(bits int, integral, lowPrecision bool) float64 {
	intval := 0
	fractval := 0
	value := 0.0
	if integral {
		value = float64(br.ReadBits(bits))
	} else {
		intval = br.ReadBits(bits)
		if lowPrecision {
			fractval = br.ReadBits(3)
			value = float64(intval) + (float64(fractval) * (1.0 / (1 << 3)))
		} else {
			fractval = br.ReadBits(5)
			value = float64(intval) + (float64(fractval) * (1.0 / (1 << 5)))
		}
	}
	return value
}
func (br *BitReader) ReadBitCoord() float64 {
	intFlag := br.ReadBoolean()
	fractFlag := br.ReadBoolean()
	value := 0.0
	if intFlag || fractFlag {
		negative := br.ReadBoolean()
		if intFlag {
			value += float64(br.ReadUBits(CoordIntegerBits)) + 1
		}
		if fractFlag {
			value += float64(br.ReadUBits(CoordFractionalBits)) * CoordResolution
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
