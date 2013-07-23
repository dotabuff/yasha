package utils

type BytesReader struct {
	Data     []byte
	Position int
}

func (br *BytesReader) ReadVarInt32() (result uint32) {
	var b uint32
	var count uint

	for {
		if count == 5 {
			return
		} else if br.Position >= len(br.Data) {
			return
		}
		b = uint32(br.Data[br.Position])
		br.Position++
		result |= (b & 0x7F) << (7 * count)
		count++
		if (b & 0x80) != 0x80 {
			break
		}
	}

	return
}

func (br *BytesReader) ReadInt32() (result int32) {
	bytes := br.Read(4)
	return int32((bytes[0] << 24) + (bytes[1] << 16) + (bytes[2] << 8) + bytes[3])
}

func (br *BytesReader) Read(length int) []byte {
	res := br.Data[br.Position:length]
	br.Position += length
	return res
}

func (br BytesReader) CanRead() bool {
	return br.Position < len(br.Data)
}
func (br *BytesReader) Skip(length int) {
	br.Position += length
}
