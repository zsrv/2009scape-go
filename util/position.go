package util

import (
	"fmt"
	"math"
)

var BuildArea = []int{104, 120, 136, 168}

type Position struct {
	X     int
	Z     int
	Plane int

	// build area
	BAIndex int
	BASizeX int
	BASizeZ int
}

func NewPosition(x int, z int, plane int) *Position {
	return &Position{
		X:     x,
		Z:     z,
		Plane: plane,

		BAIndex: 0,
		BASizeX: BuildArea[0], // 0 = BAIndex
		BASizeZ: BuildArea[0], // 0 = BAIndex
	}
}

func (p *Position) UpdateBuildArea(index int) {
	p.BAIndex = index
	p.BASizeX = BuildArea[p.BAIndex]
	p.BASizeZ = BuildArea[p.BAIndex]
}

func (p *Position) Equals(other *Position) bool {
	return p.X == other.X && p.Z == other.Z && p.Plane == other.Plane
}

func (p *Position) Copy() *Position {
	return NewPosition(p.X, p.Z, p.Plane)
}

func (p *Position) Clone(other *Position) {
	p.X = other.X
	p.Z = other.Z
	p.Plane = other.Plane
}

func (p *Position) ToString() string {
	return fmt.Sprintf("(%v, %v, %v)", p.X, p.Z, p.Plane)
}

func (p *Position) Near(other *Position, distance int) bool {
	return int(math.Abs(float64(p.X-other.X))) <= distance && int(math.Abs(float64(p.Z-other.Z))) <= distance
}

func (p *Position) DistanceTo(other *Position) float64 {
	return math.Sqrt(math.Pow(float64(p.X-other.X), 2) + math.Pow(float64(p.Z-other.Z), 2))
}

// range: 0-200
func (p *Position) MapSquareX() int {
	return p.X >> 6
}

func (p *Position) MapSquareZ() int {
	return p.Z >> 6
}

// range: 0-1600
func (p *Position) ZoneX() int {
	return p.X >> 3
}

func (p *Position) ZoneZ() int {
	return p.Z >> 3
}

// local to the mapsquare
func (p *Position) MapLocalX() int {
	return p.X & 63
}

func (p *Position) MapLocalZ() int {
	return p.Z & 63
}

// local to the build area
func (p *Position) BAStartX() int {
	return (p.ZoneX() - (p.BASizeX >> 4)) << 3
}

func (p *Position) BAEndX() int {
	return (p.ZoneX() + (p.BASizeX >> 4)) << 3
}

func (p *Position) BAStartZ() int {
	return (p.ZoneZ() - (p.BASizeZ >> 4)) << 3
}

func (p *Position) BAEndZ() int {
	return (p.ZoneZ() + (p.BASizeZ >> 4)) << 3
}

func (p *Position) BALocalX() int {
	return p.X - p.BAStartX()
}

func (p *Position) BALocalZ() int {
	return p.Z - p.BAStartZ()
}

// GPI
func (p *Position) HighRes() int {
	return p.Z | p.X<<14 | p.Plane<<28
}

func (p *Position) LowRes() int {
	return p.MapSquareZ() | p.MapSquareX()<<8 | p.Plane
}
