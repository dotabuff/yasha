package utils

import "fmt"

type Vector3 struct {
	X, Y, Z float64
}

func (v Vector3) String() string {
	return fmt.Sprintf("{{ x: %f, y: %f, z: %f }}", v.X, v.Y, v.Z)
}

type Vector2 struct {
	X, Y float64
}

func (v Vector2) String() string {
	return fmt.Sprintf("{{ x: %f, y: %f }}", v.X, v.Y)
}
