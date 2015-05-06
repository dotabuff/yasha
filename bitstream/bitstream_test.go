package bitstream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	buf := []byte("stuff like this")
	b := NewBitStream(buf)

	assert.True(t, b.Good())
	assert.EqualValues(t, b.Position(), 0)
	assert.EqualValues(t, b.End(), len(buf)*8)

	for _, n := range []int{
		1, 1, 0, 0, 1, 1, 1, 0,
		0, 0, 1, 0, 1, 1, 1, 0,
		1, 0, 1, 0, 1, 1, 1, 0,
		0, 1, 1, 0, 0, 1, 1, 0,
		0, 1, 1, 0, 0, 1, 1, 0,
		0, 0, 0, 0, 0, 1, 0, 0,
		0, 0, 1, 1, 0, 1, 1, 0,
		1, 0, 0, 1, 0, 1, 1, 0,
		1, 1, 0, 1, 0, 1, 1, 0,
		1, 0, 1, 0, 0, 1, 1, 0,
		0, 0, 0, 0, 0, 1, 0, 0,
		0, 0, 1, 0, 1, 1, 1, 0,
		0, 0, 0, 1, 0, 1, 1, 0,
		1, 0, 0, 1, 0, 1, 1, 0,
		1, 1, 0, 0, 1, 1, 1, 0,
	} {
		assert.EqualValues(t, b.Read(1), n)
	}

	assert.EqualValues(t, b.Position(), 120)
	assert.False(t, b.Good())
}

func TestReadUInt(t *testing.T) {
	buf := []byte("stuff like this")
	b := NewBitStream(buf)

	for _, n := range buf {
		assert.EqualValues(t, b.ReadUInt(8), n)
	}
}

func TestReadString(t *testing.T) {
	buf := []byte("stuff like this")
	b := NewBitStream(buf)

	assert.EqualValues(t, b.ReadString(len(buf)), "stuff like this")
}
