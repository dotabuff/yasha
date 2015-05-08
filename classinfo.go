package yasha

import "github.com/dotabuff/yasha/dota"

type ClassInfo struct {
	ID          int
	TableName   string
	NetworkName string
}

func NewClassInfo(ci *dota.CDemoClassInfoClassT) *ClassInfo {
	return &ClassInfo{
		ID:          int(ci.GetClassId()),
		TableName:   ci.GetTableName(),
		NetworkName: ci.GetNetworkName(),
	}
}

type ClassInfos struct {
	infos  []*ClassInfo
	byId   map[int]*ClassInfo
	byName map[string]*ClassInfo
}

func NewClassInfos() *ClassInfos {
	return &ClassInfos{
		infos:  []*ClassInfo{},
		byId:   map[int]*ClassInfo{},
		byName: map[string]*ClassInfo{},
	}
}

func (c *ClassInfos) Insert(info *ClassInfo) {
	c.infos = append(c.infos, info)
	c.byId[info.ID] = info
	c.byName[info.NetworkName] = info
}
