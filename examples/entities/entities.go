package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dotabuff/yasha"
	"github.com/dotabuff/yasha/utils"
)

const MAX_COORDINATE float64 = 16384

// Example of printing any updates to hero coordinates

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected a .dem file as argument")
	}

	for _, path := range os.Args[1:] {
		parser := yasha.ParserFromFile(path)
		parser.OnEntityPreserved = func(pe *yasha.PacketEntity) {
			if strings.HasPrefix(pe.Name, "DT_DOTA_Unit_Hero_") {
				if _, ok := pe.Delta["DT_DOTA_BaseNPC.m_vecOrigin"]; ok {
					coord := coordFromCell(pe)
					fmt.Printf("%30s | X: %5.0f Y: %5.0f\n", pe.Name[18:len(pe.Name)], coord.X, coord.Y)
				}
			}
		}
		parser.Parse()
	}
}

type Coordinate struct {
	X, Y float64
}

func coordFromCell(pe *yasha.PacketEntity) Coordinate {
	cellbits, ok := pe.Values["DT_BaseEntity.m_cellbits"].(int)
	if !ok {
		return Coordinate{X: 0, Y: 0}
	}
	cellWidth := float64(uint(1) << uint(cellbits))

	var cX, cY, vX, vY float64

	if vO2, ok := pe.Values["DT_DOTA_BaseNPC.m_vecOrigin"].(*utils.Vector2); ok {
		cX = float64(pe.Values["DT_DOTA_BaseNPC.m_cellX"].(int))
		cY = float64(pe.Values["DT_DOTA_BaseNPC.m_cellY"].(int))
		vX, vY = vO2.X, vO2.Y
	} else {
		vO3 := pe.Values["DT_BaseEntity.m_vecOrigin"].(*utils.Vector3)
		cX = float64(pe.Values["DT_BaseEntity.m_cellX"].(int))
		cY = float64(pe.Values["DT_BaseEntity.m_cellY"].(int))
		vX, vY = vO3.X, vO3.Y
	}

	x := ((cX * cellWidth) - MAX_COORDINATE) + vX
	y := ((cY * cellWidth) - MAX_COORDINATE) + vY

	return Coordinate{X: x, Y: y}
}
