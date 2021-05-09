package g2d

import (
	"math"
)

var geometry_origin *Geometry // Geometry with only one vertex at (0,0)

func NewGeometry_Origin() *Geometry {
	if geometry_origin == nil { // A singlton is shared for all the uses
		geometry_origin = NewGeometry().SetVertices([][2]float32{{0, 0}})
		geometry_origin.BuildDataBuffers(true, false, false)
	}
	return geometry_origin
}

func NewGeometry_Rectangle(size float32) *Geometry {
	hs := size / 2
	geometry := NewGeometry() // create an empty geometry
	geometry.SetVertices([][2]float32{{-hs, -hs}, {hs, -hs}, {hs, hs}, {-hs, hs}})
	geometry.SetFaces([][]uint32{{0, 1, 2, 3}})
	geometry.SetEdges([][]uint32{{0, 1, 2, 3}})
	return geometry
}

func NewGeometry_Triangle(size float32) *Geometry {
	geometry := NewGeometry_Polygon(3, size, -30) // 3 vertices and 1 triangular face
	geometry.SetEdges([][]uint32{{0, 1, 2, 0}})   // 1 edge connecting all the vertices
	return geometry
}

func NewGeometry_Polygon(n int, radius float32, starting_angle_in_degree float32) *Geometry {
	geometry := NewGeometry() // create an empty geometry
	radian := float64(starting_angle_in_degree * (math.Pi / 180.0))
	radian_step := (2 * math.Pi) / float64(n)
	face_indices := make([]uint32, n)
	for i := 0; i < n; i++ {
		geometry.AddVertex([2]float32{radius * float32(math.Cos(radian)), radius * float32(math.Sin(radian))})
		face_indices[i] = uint32(i)
		radian += radian_step
	}
	geometry.AddFace(face_indices)
	return geometry
}
