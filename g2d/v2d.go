package g2d

import "math"

type V2d [2]float32

// ----------------------------------------------------------------------------
// Constructor
// ----------------------------------------------------------------------------

func NewV2d(x float32, y float32) *V2d {
	return &V2d{x, y}
}

func NewV2dBySub(a [2]float32, b [2]float32) *V2d {
	return &V2d{a[0] - b[0], a[1] - b[1]}
}

func NewV2dBySum(vs ...[2]float32) *V2d {
	sum := V2d{0, 0}
	for _, v := range vs {
		sum[0] += v[0]
		sum[1] += v[1]
	}
	return &sum
}

func NewV2dByAvg(vs ...[2]float32) *V2d {
	sum := V2d{0, 0}
	for _, v := range vs {
		sum[0] += v[0]
		sum[1] += v[1]
	}
	length := float32(len(vs))
	if length > 0 {
		sum[0] /= length
		sum[1] /= length
	}
	return &sum
}

// ----------------------------------------------------------------------------
//
// ----------------------------------------------------------------------------

func (v *V2d) Clone() *V2d {
	return &V2d{v[0], v[1]}
}

func (v *V2d) Length() float32 {
	return float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1])))
}

func (v *V2d) Normalize() *V2d {
	length := v.Length()
	if length > 0 {
		v[0] /= length
		v[1] /= length
	}
	return v
}

func (v *V2d) Scale(sx float32, sy float32) *V2d {
	v[0] *= sx
	v[1] *= sy
	return v
}

func (v *V2d) Shift(dx float32, dy float32) *V2d {
	v[0] += dx
	v[1] += dy
	return v
}

func (v *V2d) Add(v2 *V2d) *V2d {
	v[0] += v2[0]
	v[1] += v2[1]
	return v
}

func (a *V2d) Dot(b *V2d) float32 {
	return (a[0]*b[0] + a[1]*b[1])
}

func (a *V2d) Cross(b *V2d) float32 {
	return (a[0]*b[1] - a[1]*b[0]) // in 2D, (ax,ay,0) x (bx,by,0) = (0,0,ax*by-ay*bx)
}

// ----------------------------------------------------------------------------
// Inclusion/Occlusion Test
// ----------------------------------------------------------------------------

func IsCCW(v0 V2d, v1 V2d, v2 V2d) bool {
	v01 := NewV2dBySub(v1, v0)
	v02 := NewV2dBySub(v2, v0)
	return v01.Cross(v02) > 0
}

func IsPointInside(p V2d, v0 V2d, v1 V2d, v2 V2d) bool {
	p0, p1, p2 := NewV2dBySub(v0, p), NewV2dBySub(v1, p), NewV2dBySub(v2, p)
	c01, c12, c20 := p0.Cross(p1), p1.Cross(p2), p2.Cross(p0)
	return (c01 > 0 && c12 > 0 && c20 > 0) || (c01 < 0 && c12 < 0 && c20 < 0)
}
