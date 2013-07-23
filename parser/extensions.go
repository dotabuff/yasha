package parser

import (
	"bytes"
)

func ReadStringZ(datas []byte, offset int) string {
	idx := bytes.IndexByte(datas[offset:], '\000')
	if idx < 0 {
		return ""
	}
	return string(datas[offset:idx])
}
