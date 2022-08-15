package g3d

import "math"

type V3d [3]float32

// ----------------------------------------------------------------------------
// Constructor
// ----------------------------------------------------------------------------

func NewV3d(x float32, y float32, z float32) *V3d {
	return &V3d{x, y, z}
}

func NewV3dBySub(a [3]float32, b [3]float32) *V3d {
	return &V3d{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
}

func NewV3dBySum(vs ...[3]float32) *V3d {
	sum := V3d{0, 0, 0}
	for _, v := range vs {
		sum[0] += v[0]
		sum[1] += v[1]
		sum[2] += v[2]
	}
	return &sum
}

func NewV3dByAvg(vs ...[3]float32) *V3d {
	sum := V3d{0, 0, 0}
	for _, v := range vs {
		sum[0] += v[0]
		sum[1] += v[1]
		sum[2] += v[2]
	}
	length := float32(len(vs))
	if length > 0 {
		sum[0] /= length
		sum[1] /= length
		sum[2] /= length
	}
	return &sum
}

func NewV3dByCross(v0 V3d, v1 V3d) *V3d {
	return v0.Cross(&v1)
}

func NewV3dByFaceNormal(v0 V3d, v1 V3d, v2 V3d) *V3d {
	v01 := NewV3dBySub(v1, v0).Normalize()
	v02 := NewV3dBySub(v2, v0).Normalize()
	return v01.Cross(v02).Normalize()
}

// ----------------------------------------------------------------------------
//
// ----------------------------------------------------------------------------

func (v *V3d) Clone() *V3d {
	return &V3d{v[0], v[1], v[2]}
}

func (v *V3d) Length() float32 {
	return float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])))
}

func (v *V3d) Normalize() *V3d {
	length := v.Length()
	if length > 0 {
		v[0] /= length
		v[1] /= length
		v[2] /= length
	}
	return v
}

func (v *V3d) Scale(sx float32, sy float32, sz float32) *V3d {
	v[0] *= sx
	v[1] *= sy
	v[2] *= sz
	return v
}

func (v *V3d) Shift(dx float32, dy float32, dz float32) *V3d {
	v[0] += dx
	v[1] += dy
	v[2] += dz
	return v
}

func (v *V3d) Add(v2 *V3d) *V3d {
	v[0] += v2[0]
	v[1] += v2[1]
	v[2] += v2[2]
	return v
}

func (v *V3d) Dot(v2 *V3d) float32 {
	return (v[0]*v2[0] + v[1]*v2[1] + v[2]*v2[2])
}

func (v *V3d) Cross(v2 *V3d) *V3d {
	return &V3d{v[1]*v2[2] - v[2]*v2[1], v[2]*v2[0] - v[0]*v2[2], v[0]*v2[1] - v[1]*v2[0]}
}

// ----------------------------------------------------------------------------
// Inclusion/Occlusion Test
// ----------------------------------------------------------------------------

func IsPointInside(p V3d, v0 V3d, v1 V3d, v2 V3d) bool {
	p0, p1, p2 := NewV3dBySub(v0, p), NewV3dBySub(v1, p), NewV3dBySub(v2, p)
	c01, c12, c20 := p0.Cross(p1), p1.Cross(p2), p2.Cross(p0)
	d012, d120, d201 := c01.Dot(c12), c12.Dot(c20), c20.Dot(c01) // test each pair with same direction
	if d012 > 0 && d120 > 0 && d201 > 0 {                        // all in the same direction
		return true // point is strictly inside the triangle
	} else {
		// 	if (Point.isZero(c01) && d120 < 0) return false;    // point is on side 01, but it's outside
		// 	if (Point.isZero(c12) && d201 < 0) return false;    // point is on side 12, but it's outside
		// 	if (Point.isZero(c20) && d012 < 0) return false;    // point is on side 20, but it's outside
		return false // point may on the border
	}
}
