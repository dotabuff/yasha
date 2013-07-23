package core

import (
	"fmt"
)

type Vector struct {
	X, Y, Z float64
}

func (v Vector) String() string {
	return fmt.Sprintf("{{ x: %f, y: %f: z: %f }}", v.X, v.Y, v.Z)
}
func (v Vector) StringXY() string {
	return fmt.Sprintf("{{ x: %f, y: %f }}", v.X, v.Y)
}
