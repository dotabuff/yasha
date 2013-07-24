package string_tables

type StringTable struct {
	Tick  int
	Index int
	Name  string
	Items map[int]*StringTableItem
}

type StringTableItem struct {
	Str  string
	Data []byte
}
