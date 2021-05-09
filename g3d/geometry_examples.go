package g3d

import (
	"math"
)

const InRadian = (math.Pi / 180.0)
const InDegree = (180.0 / math.Pi)

func NewGeometry_Polygon(n int, radius float32, starting_angle_in_degree float32) *Geometry {
	geometry := NewGeometry() // create an empty geometry
	radian := float64(starting_angle_in_degree * (math.Pi / 180.0))
	radian_step := (2 * math.Pi) / float64(n)
	face_indices := make([]uint32, n)
	for i := 0; i < n; i++ {
		geometry.AddVertex([3]float32{radius * float32(math.Cos(radian)), radius * float32(math.Sin(radian)), 0})
		face_indices[i] = uint32(i)
		radian += radian_step
	}
	geometry.AddFace(face_indices)
	return geometry
}

func NewGeometry_Cube(xsize float32, ysize float32, zsize float32) *Geometry {
	geometry := NewGeometry()
	x, y, z := xsize/2, ysize/2, zsize/2
	geometry.SetVertices([][3]float32{
		{-x, -y, -z}, {+x, -y, -z}, {+x, +y, -z}, {-x, +y, -z},
		{-x, -y, +z}, {+x, -y, +z}, {+x, +y, +z}, {-x, +y, +z}})
	geometry.SetFaces([][]uint32{
		{0, 3, 2, 1}, {0, 1, 5, 4}, {1, 2, 6, 5}, {2, 3, 7, 6}, {3, 0, 4, 7}, {4, 5, 6, 7}})
	return geometry
}

func NewGeometry_CubeWithTexture(xsize float32, ysize float32, zsize float32) *Geometry {
	geometry := NewGeometry_Cube(xsize, ysize, zsize)
	geometry.SetTextureUVs([][]float32{
		{0, 1, 1, 1, 1, 0, 0, 0}, {0, 1, 1, 1, 1, 0, 0, 0}, {0, 1, 1, 1, 1, 0, 0, 0}, {0, 1, 1, 1, 1, 0, 0, 0}, {0, 1, 1, 1, 1, 0, 0, 0}, {0, 1, 1, 1, 1, 0, 0, 0}})
	return geometry
}

func NewGeometry_Sphere(radius float32, wsegs int, hsegs int) *Geometry {
	//   Sphere with the minimum number of vertices (to be used with face normal vectors)
	geometry := NewGeometry()
	wnum, hnum := wsegs+1, hsegs+1
	wstep := math.Pi * 2.0 / float32(wsegs)
	hstep := math.Pi / float32(hsegs)
	for i := 0; i < wnum; i++ {
		lon := float64(wstep * float32(i)) // longitude (λ)
		for j := 1; j <= hnum; j++ {
			lat := float64(-math.Pi/2.0 + hstep*float32(j)) // latitude (φ)
			x := radius * float32(math.Cos(lon)*math.Cos(lat))
			y := radius * float32(math.Sin(lon)*math.Cos(lat))
			z := radius * float32(math.Sin(lat))
			geometry.AddVertex([3]float32{x, y, z})
		}
	}
	npole := geometry.AddVertex([3]float32{0, 0, +radius}) // north pole
	spole := geometry.AddVertex([3]float32{0, 0, -radius}) // south pole
	for i := 0; i < wnum; i++ {                            // quadratic faces on the side
		for j := 0; j < hnum-1; j++ {
			start := uint32((i+0)*hnum + j)
			wnext := uint32((i+1)%wnum*hnum + j)
			geometry.AddFace([]uint32{start, wnext, wnext + 1, start + 1})
		}
	}
	for i := 0; i < wnum; i++ { // faces around south pole
		start := uint32((i + 0) * hnum)
		wnext := uint32((i + 1) % wnum * hnum)
		geometry.AddFace([]uint32{spole, wnext, start})
	}
	for i := 0; i < wnum; i++ { // faces around north pole
		start := uint32((i+0)*hnum + (hnum - 1))
		wnext := uint32((i+1)%wnum*hnum + (hnum - 1))
		geometry.AddFace([]uint32{npole, start, wnext})
	}
	return geometry
}

func NewGeometry_Cylinder(nsides int, radius float32, height float32, stt_angle float32, solid bool) *Geometry {
	geometry := NewGeometry()
	rad_step := math.Pi * 2.0 / float64(nsides)
	for i := 0; i < nsides; i++ {
		rad := rad_step*float64(i) + float64(stt_angle)*InRadian
		cosR, sinR := float32(math.Cos(rad))*radius, float32(math.Sin(rad))*radius
		geometry.AddVertex([3]float32{cosR, sinR, 0})
		geometry.AddVertex([3]float32{cosR, sinR, height})
		i2, next2 := uint32(i*2), uint32((i+1)%nsides)*2
		geometry.AddFace([]uint32{i2 + 1, i2, next2, next2 + 1}) // face on the side
	}
	if solid {
		btm_face_indices := make([]uint32, nsides)
		top_face_indices := make([]uint32, nsides)
		for i := 0; i < nsides; i++ {
			btm_face_indices[nsides-1-i] = uint32(i * 2) // reversed
			top_face_indices[i] = uint32(i*2 + 1)        //
		}
		geometry.AddFace(btm_face_indices) // face on the bottom (reversed)
		geometry.AddFace(top_face_indices) // face on the top
	}
	return geometry
}

func NewGeometry_Pyramid(nsides int, radius float32, height float32, stt_angle float32, solid bool) *Geometry {
	geometry := NewGeometry()
	rad_step := math.Pi * 2.0 / float64(nsides)
	apex := uint32(nsides)
	for i := 0; i < nsides; i++ {
		rad := rad_step*float64(i) + float64(stt_angle)*InRadian
		cosR, sinR := float32(math.Cos(rad))*radius, float32(math.Sin(rad))*radius
		geometry.AddVertex([3]float32{cosR, sinR, 0.0})
		next := uint32((i + 1) % nsides)
		geometry.AddFace([]uint32{apex, uint32(i), next}) // face on the side
	}
	geometry.AddVertex([3]float32{0, 0, height}) // apex
	if solid {
		btm_face_indices := make([]uint32, nsides)
		for i := 0; i < nsides; i++ {
			btm_face_indices[nsides-1-i] = uint32(i) // reversed
		}
		geometry.AddFace(btm_face_indices) // face on the bottom
	}
	return geometry
}

func NewGeometry_SolidFromFaceAndHeight(face [][3]float32, height float32) *Geometry {
	geometry := NewGeometry()
	flen := len(face)
	btm_list, top_list := make([]uint32, flen), make([]uint32, flen)
	for i := 0; i < flen; i++ {
		b := geometry.AddVertex([3]float32{face[i][0], face[i][1], face[i][2]})
		t := geometry.AddVertex([3]float32{face[i][0], face[i][1], face[i][2] + height})
		b2, t2 := uint32((int(b)+2)%(2*flen)), uint32((int(t)+2)%(2*flen))
		geometry.AddFace([]uint32{b, b2, t2, t}) // face on the side
		top_list[i] = t
		btm_list[flen-1-i] = b
	}
	geometry.AddFace(btm_list) // bottom face
	geometry.AddFace(top_list) // top face
	return geometry
}

func NewGeometry_SolidFromCentersAndRadii(centers [][3]float32, radii []float32, nsegments int) *Geometry {
	geometry := NewGeometry()
	nsides, astep := uint32(nsegments), math.Pi*2/float64(nsegments)
	geometry.AddVertex(centers[0]) // bottom center
	for i := 1; i < len(centers); i++ {
		center := centers[i]
		radius := float64(radii[i])
		curr := uint32(len(geometry.verts)) // geometry.countVertices();
		if i == 1 {
			for j := uint32(0); j < nsides; j++ {
				a := astep * float64(j)
				cosA, sinA := float32(math.Cos(a)*radius), float32(math.Sin(a)*radius)
				geometry.AddVertex([3]float32{center[0] + cosA, center[1] + sinA, center[2]})
				geometry.AddFace([]uint32{0, curr + (j+1)%nsides, curr + j})
			}
		} else if i < len(centers)-1 {
			prev := curr - nsides
			for j := uint32(0); j < nsides; j++ {
				a := astep * float64(j)
				cosA, sinA := float32(math.Cos(a)*radius), float32(math.Sin(a)*radius)
				geometry.AddVertex([3]float32{center[0] + cosA, center[1] + sinA, center[2]})
				geometry.AddFace([]uint32{prev + j, prev + (j+1)%nsides, curr + (j+1)%nsides, curr + j})
			}
		} else {
			prev := curr - nsides
			geometry.AddVertex(centers[i]) // top center
			for j := uint32(0); j < nsides; j++ {
				geometry.AddFace([]uint32{curr, prev + j, prev + (j+1)%nsides})
			}
		}
	}
	return geometry
}

func NewGeometry_EmptyExample() *Geometry {
	geometry := NewGeometry()
	geometry.SetVertices([][3]float32{{0, 0, 0}, {1, 0, 0}, {0, 1, 0}, {0, 0, 1}})
	geometry.SetFaces([][]uint32{{0, 2, 1}, {0, 1, 3}, {1, 2, 3}, {2, 0, 3}})
	return geometry
}
