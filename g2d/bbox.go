package g2d

import "math"

type BBox [2][2]float32

// ----------------------------------------------------------------------------
//
// ----------------------------------------------------------------------------

func NewBBox(minx float32, miny float32, maxx float32, maxy float32) *BBox {
	return &BBox{{minx, miny}, {maxx, maxy}}
}

func NewBBoxEmpty() *BBox {
	return &BBox{{+math.MaxFloat32, +math.MaxFloat32}, {-math.MaxFloat32, -math.MaxFloat32}}
}

// ----------------------------------------------------------------------------
// Basic Member Functions
// ----------------------------------------------------------------------------

func (b *BBox) IsEmpty() bool {
	return (b[0][0] > b[1][0] || b[0][1] > b[1][1])
}

func (b *BBox) Width() float32 {
	return b[1][0] - b[0][0]
}

func (b *BBox) Height() float32 {
	return b[1][1] - b[0][1]
}

func (b *BBox) Shape() [2]float32 {
	return [2]float32{(b[1][0] - b[0][0]), (b[1][1] - b[0][1])}
}

func (b *BBox) Center() [2]float32 {
	return [2]float32{(b[0][0] + b[1][0]) / 2, (b[0][1] + b[1][1]) / 2}
}

func (b *BBox) AddPoint(p *[2]float32) {
	if b[0][0] > p[0] {
		b[0][0] = p[0]
	}
	if b[0][1] > p[1] {
		b[0][1] = p[1]
	}
	if b[1][0] < p[0] {
		b[1][0] = p[0]
	}
	if b[1][1] < p[1] {
		b[1][1] = p[1]
	}
}

func (b *BBox) Merge(b2 *BBox) *BBox {
	if b[0][0] > b2[0][0] {
		b[0][0] = b2[0][0]
	}
	if b[0][1] > b2[0][1] {
		b[0][1] = b2[0][1]
	}
	if b[1][0] < b2[1][0] {
		b[1][0] = b2[1][0]
	}
	if b[1][1] < b2[1][1] {
		b[1][1] = b2[1][1]
	}
	return b
}

// ----------------------------------------------------------------------------
// Inclusion/Occlusion Test
// ----------------------------------------------------------------------------

func (b *BBox) IsIncludingPoint(v *V2d) bool {
	return b[0][0] <= v[0] && v[0] <= b[1][0] && b[0][1] <= v[1] && v[1] <= b[1][1]
}
