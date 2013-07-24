package packet_entities

const (
	serialNumBits = 10
)

type PacketEntity struct {
	Tick      int
	Index     int
	SerialNum int
	ClassId   int
	Name      string
	Type      UpdateType
	Values    map[string]interface{}
}

func (pe *PacketEntity) Handle() int {
	return pe.Index | (pe.SerialNum << (serialNumBits + 1))
}

func (pe *PacketEntity) Clone() *PacketEntity {
	values := map[string]interface{}{}
	for key, value := range pe.Values {
		values[key] = value
	}
	return &PacketEntity{
		Tick:      pe.Tick,
		Index:     pe.Index,
		SerialNum: pe.SerialNum,
		ClassId:   pe.ClassId,
		Name:      pe.Name,
		Type:      pe.Type,
		Values:    values,
	}
}
