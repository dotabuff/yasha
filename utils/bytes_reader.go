package utils

type BytesReader struct {
	data     []byte
	position int
}

func NewBytesReader(data []byte) *BytesReader {
	return &BytesReader{data: data, position: 0}
}

func (br *BytesReader) ReadVarInt32() (result int32) {
	var b int32 = 0x80
	var count uint

	for (b & 0x80) == 0x80 {
		if count == 5 {
			break
		} else if br.position >= len(br.data) {
			break
		}
		b = int32(br.data[br.position])
		br.position++
		result |= (b & 0x7F) << (7 * count)
		count++
	}

	return
}

func (br *BytesReader) ReadInt32() (result int32) {
	bytes := br.Read(4)
	return int32((bytes[0] << 24) + (bytes[1] << 16) + (bytes[2] << 8) + bytes[3])
}

func (br *BytesReader) Read(length int) []byte {
	res := br.data[br.position:(br.position + length)]
	br.position += length
	return res
}

func (br BytesReader) CanRead() bool {
	return br.position < len(br.data)
}
func (br *BytesReader) Skip(length int) {
	br.position += length
}
