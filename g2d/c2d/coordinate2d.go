package c2d

import "math"

func AddAB(a [2]float32, b [2]float32) [2]float32 {
	return [2]float32{a[0] + b[0], a[1] + b[1]}
}

func SubAB(a [2]float32, b [2]float32) [2]float32 {
	return [2]float32{a[0] - b[0], a[1] - b[1]}
}

func CrossAB(a [2]float32, b [2]float32) float32 {
	return a[0]*b[1] - a[1]*b[0] // in 2D, (ax,ay,0) x (bx,by,0) = (0,0,ax*by-ay*bx)
}

func Length(v [2]float32) float32 {
	return float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1])))
}

func Normalize(v [2]float32) [2]float32 {
	len := Length(v)
	return [2]float32{v[0] / len, v[1] / len}
}

func IsCCW(v0 [2]float32, v1 [2]float32, v2 [2]float32) bool {
	v01 := SubAB(v1, v0)
	v02 := SubAB(v2, v0)
	return CrossAB(v01, v02) > 0
}

func IsPointInside(p [2]float32, v0 [2]float32, v1 [2]float32, v2 [2]float32) bool {
	p0, p1, p2 := SubAB(v0, p), SubAB(v1, p), SubAB(v2, p)
	c01, c12, c13 := CrossAB(p0, p1), CrossAB(p1, p2), CrossAB(p2, p0)
	return (c01 > 0 && c12 > 0 && c13 > 0) || (c01 < 0 && c12 < 0 && c13 < 0)
}

// ----------------------------------------------------------------------------
// Bounding Box
// ----------------------------------------------------------------------------

func BBoxInit() [2][2]float32 {
	return [2][2]float32{{+3.4e38, +3.4e38}, {-3.4e38, -3.4e38}}
}

func BBoxIsSet(b [2][2]float32) bool {
	return (b[0][0] <= b[1][0] && b[0][1] <= b[1][1])
}

func BBoxAddPoint(bbox *[2][2]float32, v [2]float32) {
	if bbox[0][0] > v[0] {
		bbox[0][0] = v[0]
	}
	if bbox[0][1] > v[1] {
		bbox[0][1] = v[1]
	}
	if bbox[1][0] < v[0] {
		bbox[1][0] = v[0]
	}
	if bbox[1][1] < v[1] {
		bbox[1][1] = v[1]
	}
}

func BBoxMerge(b1 [2][2]float32, b2 [2][2]float32) [2][2]float32 {
	if b1[0][0] > b2[0][0] {
		b1[0][0] = b2[0][0]
	}
	if b1[0][1] > b2[0][1] {
		b1[0][1] = b2[0][1]
	}
	if b1[1][0] < b2[1][0] {
		b1[1][0] = b2[1][0]
	}
	if b1[1][1] < b2[1][1] {
		b1[1][1] = b2[1][1]
	}
	return b1
}

func BBoxSize(b [2][2]float32) [2]float32 {
	return [2]float32{(-b[0][0] + b[1][0]), (-b[0][1] + b[1][1])}
}

func BBoxCenter(b [2][2]float32) [2]float32 {
	return [2]float32{(b[0][0] + b[1][0]) / 2, (b[0][1] + b[1][1]) / 2}
}

func BBoxInside(b [2][2]float32, v [2]float32) bool {
	return b[0][0] <= v[0] && v[0] <= b[1][0] && b[0][1] <= v[1] && v[1] <= b[1][1]
}
