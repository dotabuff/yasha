package bitstream

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

var pp = spew.Dump

type BitStream struct {
	size int
	pos  int
	data []byte
}

func NewBitStream(data []byte) *BitStream {
	return &BitStream{
		data: data,
		size: len(data) * 8,
	}
}

const (
	// Number of fraction bits in a normalized float
	NORMAL_FRACTION_BITS = 11
	// Normal denominator
	NORMAL_DENOMINATOR = 1<<NORMAL_FRACTION_BITS - 1
	// Normal resolution
	NORMAL_RESOLUTION = 1.0 / NORMAL_DENOMINATOR
	// Maximum number of bytes to read for a 32 bit varint
	VARINT32_MAX = 32/8 + 1
	// Maximum number of bytes to read for a 64 bit varint
	VARINT64_MAX = 64/8 + 2
	// Number of bits to read for the integer part of a coord
	COORD_INTEGER_BITS = 14
	// Number of bits to read for the fraction part of a coord
	COORD_FRACTION_BITS = 5
	// Coord demoniator
	COORD_DENOMINATOR = 1 << COORD_FRACTION_BITS
	// Coord resolution
	COORD_RESOLUTION = 1.0 / COORD_DENOMINATOR
	// Bits to read for multiplayer optimized coordinates
	COORD_INTEGER_BITS_MP = 11
	// Fractional part of low-precision coords
	COORD_FRACTION_BITS_MP_LOWPRECISION = 3
	// Denominator for low-precision coords
	COORD_DENOMINATOR_LOWPRECISION = 1 << COORD_FRACTION_BITS_MP_LOWPRECISION
	// Resolution for low-precision coords
	COORD_RESOLUTION_LOWPRECISION = 1.0 / COORD_DENOMINATOR_LOWPRECISION
	// Number of bits to read for the fraction part of a cell coord
	CELL_COORD_FRACTION_BITS = 5
	// Number of bits to read for a low-precision cell coord's fraction
	CELL_COORD_FRACTION_BITS_LOWPRECISION = 3

	errBitStreamOverflow = "More bits requested than available (%d > %d)"
)

func (b *BitStream) Good() bool {
	return b.pos < b.size
}

func (b *BitStream) End() int {
	return b.size
}

func (b *BitStream) Position() int {
	return b.pos
}

func (b *BitStream) seekForward(n int) {
	b.pos += n

	if b.pos > b.size {
		b.pos = b.size
	}
}

func (b *BitStream) seekBackward(n int) {
	if b.pos-n > b.pos {
		b.pos = 0
	} else {
		b.pos -= n
	}
}

// Returns result of reading n bits as uint32.
// This function can read a maximum of 32 bits at once. If the amount of data requested
// exceeds the remaining size of the current chunk it wraps around to the next one.
func (b *BitStream) Read(n int) (ret uint) {
	if b.pos+n > b.size {
		panic(fmt.Errorf(errBitStreamOverflow, n, b.size-b.pos))
	}

	if n > 32 {
		panic(fmt.Errorf(errBitStreamOverflow, n, 32))
	}

	// shortcut for aligned
	if b.pos%8 == 0 && n%8 == 0 {
		var result uint
		for i := 0; i < n/8; i++ {
			result += uint(b.data[b.pos/8] << (uint(i) * 8))
			b.pos += 8
		}
		return result
	}

	bitOffset := b.pos % 8
	nBitsToRead := bitOffset + n
	nBytesToRead := nBitsToRead / 8
	if nBitsToRead%8 != 0 {
		nBytesToRead += 1
	}

	var currentValue uint64
	for i := 0; i < nBytesToRead; i++ {
		m := b.data[(b.pos/8)+i]
		currentValue += (uint64(m) << (uint64(i) * 8))
	}
	currentValue >>= uint(bitOffset)
	currentValue &= ((1 << uint64(n)) - 1)
	b.pos += n
	return uint(currentValue)
}

func (b *BitStream) ReadBool() bool {
	return b.Read(1) == 1
}

func (b *BitStream) ReadUInt(n int) (ret uint) {
	return b.Read(n)
}

func (b *BitStream) ReadSInt(n int) (ret uint) {
	ret = b.Read(n)
	var sign uint = 1 << uint(n-1)

	if ret >= sign {
		ret = ret - sign - sign
	}

	return ret
}

func (b *BitStream) ReadNormal() (ret float32) {
	sign := b.ReadBool()
	fraction := b.Read(NORMAL_FRACTION_BITS)
	ret = float32(fraction) * NORMAL_RESOLUTION

	if sign {
		return -ret
	}
	return ret
}

func (b *BitStream) SkipNormal() {
	b.seekForward(NORMAL_FRACTION_BITS + 1)
}

func (b *BitStream) ReadVarUInt32() (ret uint) {
	var readCount, tmpBuf uint

	for {
		if readCount == VARINT32_MAX {
			return ret // return when maximum number of iterations is reached
		}

		tmpBuf = b.Read(8)
		ret |= (tmpBuf & 0x7F) << (7 * readCount)
		readCount += 1
		if (tmpBuf & 0x80) == 0 {
			break
		}
	}

	return ret
}

func (b *BitStream) ReadVarSInt32() (ret uint) {
	ret = b.ReadVarUInt32()
	return (ret >> 1) ^ -(ret & 1)
}

func (b *BitStream) ReadVarUInt64() (ret uint64) {
	var readCount, tmpBuf uint

	for {
		if readCount == VARINT64_MAX {
			return ret
		}

		tmpBuf = b.Read(8)
		ret |= uint64(tmpBuf&0x7F) << (7 * readCount)
		readCount++
		if (tmpBuf & 0x80) == 0 {
			break
		}
	}

	return ret
}

func (b *BitStream) ReadVarSInt64() (ret uint64) {
	ret = b.ReadVarUInt64()
	return (ret >> 1) ^ -(ret & 1)
}

func (b *BitStream) ReadCoord() (ret float32) {
	hasInteger, hasFraction := b.ReadBool(), b.ReadBool()

	if hasInteger || hasFraction {
		signBit := b.ReadBool()

		if hasInteger {
			ret += float32(b.Read(COORD_INTEGER_BITS) + 1)
		}

		if hasFraction {
			ret += float32(b.Read(COORD_FRACTION_BITS)) * COORD_RESOLUTION
		}

		if signBit {
			return -ret
		}
		return ret
	}

	panic("Returning default coord")
}

var (
	readCoordMpMult = [4]float32{
		+1 / (1 << COORD_FRACTION_BITS),
		-1 / (1 << COORD_FRACTION_BITS),
		+1 / (1 << COORD_FRACTION_BITS_MP_LOWPRECISION),
		-1 / (1 << COORD_FRACTION_BITS_MP_LOWPRECISION),
	}

	readCoordMpBits = [8]int{
		COORD_FRACTION_BITS,
		COORD_FRACTION_BITS,
		COORD_FRACTION_BITS + COORD_INTEGER_BITS,
		COORD_FRACTION_BITS + COORD_INTEGER_BITS_MP,
		COORD_FRACTION_BITS_MP_LOWPRECISION,
		COORD_FRACTION_BITS_MP_LOWPRECISION,
		COORD_FRACTION_BITS_MP_LOWPRECISION + COORD_INTEGER_BITS,
		COORD_FRACTION_BITS_MP_LOWPRECISION + COORD_INTEGER_BITS_MP,
	}
)

func (b *BitStream) ReadCoordMp(integral, lowPrecision int) float32 {
	var flags uint
	if integral != 0 {
		flags = b.Read(3)
	} else {
		flags = b.Read(2)
	}

	const (
		isInbound uint = 1
		isIntval  uint = 2
		isSign    uint = 4
	)

	if integral != 0 {
		if flags&isIntval != 0 {
			var toRead int
			if flags&isInbound != 0 {
				toRead = COORD_INTEGER_BITS_MP + 1
			} else {
				toRead = COORD_INTEGER_BITS + 1
			}
			bits := b.Read(toRead)

			var ret float32 = float32((bits >> 1) + 1)
			if bits&1 == 0 {
				return ret
			} else {
				return -ret
			}
		}

		return 0
	}

	var multiply float32
	if flags&isSign != 0 {
		multiply = readCoordMpMult[1+lowPrecision*2]
	} else {
		multiply = readCoordMpMult[lowPrecision*2]
	}

	val := b.Read(readCoordMpBits[(flags&(isInbound|isIntval))+uint(lowPrecision)*4])

	if flags&isIntval != 0 {
		var fracMp uint = val >> COORD_INTEGER_BITS_MP
		var frac uint = val >> COORD_FRACTION_BITS

		var maskMp uint = 1<<COORD_INTEGER_BITS_MP - 1
		var mask uint = 1<<COORD_FRACTION_BITS - 1

		var selectNotMp uint = (flags & isInbound) - 1

		frac -= fracMp
		frac &= selectNotMp
		frac += fracMp

		mask -= maskMp
		mask &= selectNotMp
		mask += maskMp

		intpart := (val & mask) + 1
		intbitsLow := intpart << COORD_FRACTION_BITS_MP_LOWPRECISION
		intbits := intpart << COORD_FRACTION_BITS
		selectNotLow := uint(lowPrecision) - 1

		intbits -= intbitsLow
		intbits &= selectNotLow
		intbits += intbitsLow

		val = frac | intbits
	}

	return float32(val) * multiply
}

func (b *BitStream) ReadCellCoord(n int, integral, lowPrecision bool) (ret float32) {
	val := b.Read(n)

	if integral {
		if val&0x80 != 0 {
			ret += 4.2949673e9 // because why not?
		}

		return ret
	}

	if lowPrecision {
		return (float32(val) + 0.125) * float32(b.Read(CELL_COORD_FRACTION_BITS_LOWPRECISION))
	}

	return (float32(val) + 0.03125) * float32(b.Read(CELL_COORD_FRACTION_BITS))
}

func (b *BitStream) ReadString(n int) string {
	buffer := []byte{}

	for n > 0 {
		got := byte(b.Read(8))
		if got == 0 {
			break
		}
		buffer = append(buffer, got)
		n--
	}

	return string(buffer)
}

func (b *BitStream) ReadBits(n int) []byte {
	buffer := []byte{}
	remaining := n

	for remaining >= 8 {
		buffer = append(buffer, byte(b.Read(8)))
		remaining -= 8
	}

	if remaining > 0 {
		buffer = append(buffer, byte(b.Read(remaining)))
	}

	return buffer
}

func (b *BitStream) ReadBitsAsBytes(n int) []byte {
	buffer := make([]byte, (n+7)/8)

	i := 0

	for n > 7 {
		n -= 8
		buffer[i] = byte(b.Read(8))
		i++
	}

	if n != 0 {
		buffer[i] = byte(b.Read(n))
	}

	return buffer
}
