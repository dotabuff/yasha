package packet_entities

type UpdateType int

const (
	Create UpdateType = iota
	Delete
	Leave
	Preserve
)
