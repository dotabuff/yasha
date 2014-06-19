package string_tables

import (
	"math"

	"github.com/davecgh/go-spew/spew"
	"github.com/dotabuff/d2rp/core/utils"
)

const (
	MaxNameLength  = 0x400
	KeyHistorySize = 32
)

type CSTObject interface {
	GetStringData() []byte
	GetNumEntries() int32
	GetMaxEntries() int32
	GetUserDataFixedSize() bool
	GetUserDataSizeBits() int32
}

type USTObject interface {
	GetStringData() []byte
	GetNumChangedEntries() int32
}

func ParseCST(obj CSTObject) map[int]*StringTableItem {
	return Parse(
		obj.GetStringData(),
		int(obj.GetNumEntries()),
		int(obj.GetMaxEntries()),
		int(obj.GetUserDataSizeBits()),
		obj.GetUserDataFixedSize(),
	)
}

func ParseUST(obj USTObject, meta *CacheItem) map[int]*StringTableItem {
	return Parse(
		obj.GetStringData(),
		int(obj.GetNumChangedEntries()),
		meta.MaxEntries,
		meta.Bits,
		meta.IsFixedSize,
	)
}

func Parse(data []byte, numEntries, maxEntries, dataSizeBits int, dataFixedSize bool) map[int]*StringTableItem {
	br := utils.NewBitReader(data)

	bitsPerIndex := int(math.Log(float64(maxEntries)) / math.Log(2))
	keyHistory := make([]string, 0, KeyHistorySize)
	result := map[int]*StringTableItem{}
	mysteryFlag := br.ReadBoolean()
	index := -1
	nameBuf := ""

	for len(result) < numEntries {
		if br.ReadBoolean() {
			index++
		} else {
			index = int(br.ReadUBits(bitsPerIndex))
		}
		nameBuf = ""
		if br.ReadBoolean() {
			if mysteryFlag && br.ReadBoolean() {
				panic("mysteryFlag assertion failed!")
			}
			if br.ReadBoolean() {
				basis := br.ReadUBits(5)
				length := br.ReadUBits(5)
				if int(basis) > len(keyHistory) {
					spew.Dump("Ignoring invalid history index...", keyHistory, basis, length)
					nameBuf += br.ReadStringN(MaxNameLength)
				} else {
					nameBuf += keyHistory[basis][0:length] + br.ReadStringN(int(MaxNameLength-length))
				}
			} else {
				nameBuf += br.ReadStringN(MaxNameLength)
			}
			if len(keyHistory) >= KeyHistorySize {
				copy(keyHistory[0:], keyHistory[1:])
				keyHistory[len(keyHistory)-1] = "" // or the zero value of T
				keyHistory = keyHistory[:len(keyHistory)-1]
			}
			keyHistory = append(keyHistory, nameBuf)
		}
		value := []byte{}
		if br.ReadBoolean() {
			bitLength := 0
			if dataFixedSize {
				bitLength = dataSizeBits
			} else {
				bitLength = int(br.ReadUBits(14) * 8)
			}
			value = append(value, br.ReadBitsAsBytes(bitLength)...)
		}
		result[index] = &StringTableItem{Str: nameBuf, Data: value}
	}

	return result
}

/*
		item := &StringTableItem{}
		entryIndex := index + 1

		if !br.ReadBoolean() {
			entryIndex = int(br.ReadUBits(bitsPerIndex))
		}

		index = entryIndex

		if br.ReadBoolean() {
			value := ""
			substringcheck := br.ReadBoolean()
			if mysteryFlag && substringcheck {
				panic("substringcheck and mysteryFlag are true! Whatever shall we do?")
			} else if substringcheck {
				index := int(br.ReadUBits(5))
				bytestocopy := int(br.ReadUBits(5))
				value = keyHistory[index][0:bytestocopy] + br.ReadString()
			} else {
				value = br.ReadString()
			}
			item.Str = value
			keyHistory = append(keyHistory, value)
		}

		if br.ReadBoolean() {
			if isFixedSize {
				item.Data = []byte{byte(br.ReadBits(int(numBits)))}
			} else {
				length := int(br.ReadUBits(14))
				item.Data = br.ReadBytes(length)
			}
		}
		if len(keyHistory) > 32 {
			keyHistory = keyHistory[1:]
		}
		result[entryIndex] = item
	}
	return result
}
*/
