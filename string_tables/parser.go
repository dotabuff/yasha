package string_tables

import (
	"math"

	"github.com/dotabuff/d2rp/core/utils"
)

func Parse(bytes []byte, numEntries, maxEntries int32, isFixedSize bool, numBits int32) map[int]*StringTableItem {
	result := map[int]*StringTableItem{}

	lastEntry := -1
	history := make([]string, 0, 33)

	br := utils.NewBitReader(bytes)
	isOption := br.ReadBoolean()
	bitsPerIndex := int(math.Log(float64(maxEntries)) / math.Log(2))

	// br.SeekBits(1, utils.Begin)

	for i := int32(0); i < numEntries; i++ {
		item := &StringTableItem{}
		entryIndex := lastEntry + 1

		if br.ReadBoolean() {
			entryIndex = int(br.ReadUBits(bitsPerIndex))
		} else {
			entryIndex++
		}

		lastEntry = entryIndex

		if br.ReadBoolean() {
			value := ""
			if isOption && br.ReadBoolean() {
				panic("this is... wrong?")
			}

			if br.ReadBoolean() {
				index := int(br.ReadUBits(5))
				bytestocopy := int(br.ReadUBits(5))
				if index > len(history) {
					value = br.ReadString()
				} else {
					value = history[index][0:bytestocopy] + br.ReadString()
				}
			} else {
				value = br.ReadString()
			}
			item.Str = value
			history = append(history, value)
		}

		if br.ReadBoolean() {
			length := 0
			if isFixedSize {
				length = int(numBits)
			} else {
				length = int(br.ReadUBits(14) * 8)
			}
			item.Data = br.ReadBytes(length)
		}
		if len(history) > 32 {
			history = history[1:]
		}
		result[entryIndex] = item
	}
	return result
}
