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

var (
	masks = [64]uint64{
		0x0, 0x1, 0x3, 0x7,
		0xf, 0x1f, 0x3f, 0x7f,
		0xff, 0x1ff, 0x3ff, 0x7ff,
		0xfff, 0x1fff, 0x3fff, 0x7fff,
		0xffff, 0x1ffff, 0x3ffff, 0x7ffff,
		0xfffff, 0x1fffff, 0x3fffff, 0x7fffff,
		0xffffff, 0x1ffffff, 0x3ffffff, 0x7ffffff,
		0xfffffff, 0x1fffffff, 0x3fffffff, 0x7fffffff,
		0xffffffff, 0x1ffffffff, 0x3ffffffff, 0x7ffffffff,
		0xfffffffff, 0x1fffffffff, 0x3fffffffff, 0x7fffffffff,
		0xffffffffff, 0x1ffffffffff, 0x3ffffffffff, 0x7ffffffffff,
		0xfffffffffff, 0x1fffffffffff, 0x3fffffffffff, 0x7fffffffffff,
		0xffffffffffff, 0x1ffffffffffff, 0x3ffffffffffff, 0x7ffffffffffff,
		0xfffffffffffff, 0x1fffffffffffff, 0x3fffffffffffff, 0x7fffffffffffff,
	}

	shift = [64]byte{
		0x0, 0x1, 0x2, 0x3,
		0x4, 0x5, 0x6, 0x7,
		0x8, 0x9, 0xa, 0xb,
		0xc, 0xd, 0xe, 0xf,
		0x10, 0x11, 0x12, 0x13,
		0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b,
		0x1c, 0x1d, 0x1e, 0x1f,
		0x20, 0x21, 0x22, 0x23,
		0x24, 0x25, 0x26, 0x27,
		0x28, 0x29, 0x2a, 0x2b,
		0x2c, 0x2d, 0x2e, 0x2f,
		0x30, 0x31, 0x32, 0x33,
		0x34, 0x35, 0x36, 0x37,
		0x38, 0x39, 0x3a, 0x3b,
		0x3c, 0x3d, 0x3e, 0x3f,
	}
)

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
	if n > b.size-b.pos {
		panic(fmt.Errorf(errBitStreamOverflow, n, b.size-b.pos))
	}

	if n > 32 {
		panic(fmt.Errorf(errBitStreamOverflow, n, 32))
	}

	start := b.pos / 8
	end := (b.pos + n - 1) / 8
	s := b.pos % 8
	if start == end {
		ret = uint(uint64(b.data[start]>>shift[s]) & masks[n])
	} else { // wrap around
		ret = uint(uint64(b.data[start]>>shift[s]) | uint64(b.data[end]<<(8-shift[s]))&masks[n])
	}

	b.pos += n
	return ret
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
	sign := b.Read(1)
	fraction := b.Read(NORMAL_FRACTION_BITS)
	ret = float32(fraction) * NORMAL_RESOLUTION

	if sign == 1 {
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
		if (tmpBuf & 0x80) != 0 {
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
		readCount += 1
		if (tmpBuf & 0x80) != 0 {
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
	integer, fraction := b.Read(1), b.Read(1)

	if integer != 0 || fraction != 0 {
		signBit := b.Read(1)

		if integer != 0 {
			integer = b.Read(COORD_INTEGER_BITS)
			integer += 1
		}

		if fraction != 0 {
			fraction = b.Read(COORD_FRACTION_BITS)
		}

		ret = float32(integer) + float32(fraction)*COORD_RESOLUTION

		if signBit != 0 {
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
	buffer := make([]byte, n)

	for i := 0; i < n; i++ {
		buffer[i] = byte(b.Read(8))
		if buffer[i] == 0 {
			break
		}
	}

	return string(buffer)
}

func (b *BitStream) ReadBits(n int) []byte {
	buffer := make([]byte, n/8)
	remaining := n
	i := 0

	for remaining >= 8 {
		i++
		buffer[i] = byte(b.Read(8))
		remaining -= 8
	}

	if remaining > 8 {
		i++
		buffer[i] = byte(b.Read(remaining))
	}

	return buffer
}
