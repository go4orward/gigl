package c3d

import "math"

func AddAB(a [3]float32, b [3]float32) [3]float32 {
	return [3]float32{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
}

func SubAB(a [3]float32, b [3]float32) [3]float32 {
	return [3]float32{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
}

func AverageAB(a [3]float32, b [3]float32) [3]float32 {
	return [3]float32{(a[0] + b[0]) / 2, (a[1] + b[1]) / 2, (a[2] + b[2]) / 2}
}

func SumAll(all [][3]float32) [3]float32 {
	sum := [3]float32{0, 0, 0}
	for _, v := range all {
		sum[0] += v[0]
		sum[1] += v[1]
		sum[2] += v[2]
	}
	return sum
}

func AverageAll(all [][3]float32) [3]float32 {
	sum := SumAll(all)
	count := float32(len(all))
	return [3]float32{sum[0] / count, sum[1] / count, sum[2] / count}
}

func DotAB(a [3]float32, b [3]float32) float32 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func CrossAB(a [3]float32, b [3]float32) [3]float32 {
	return [3]float32{a[1]*b[2] - a[2]*b[1], a[2]*b[0] - a[0]*b[2], a[0]*b[1] - a[1]*b[0]}
}

func Length(v [3]float32) float32 {
	return float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])))
}

func Normalize(v [3]float32) [3]float32 {
	len := Length(v)
	return [3]float32{v[0] / len, v[1] / len, v[2] / len}
}

func GetFaceNormal(v0 [3]float32, v1 [3]float32, v2 [3]float32) bool {
	// v01 := SubAB(v1, v0)
	// v02 := SubAB(v2, v0)
	// return CrossAB(v01, v02) > 0
	return false
}

func IsPointInside(p [3]float32, v0 [3]float32, v1 [3]float32, v2 [3]float32) bool {
	p0, p1, p2 := SubAB(v0, p), SubAB(v1, p), SubAB(v2, p)
	c01, c12, c20 := CrossAB(p0, p1), CrossAB(p1, p2), CrossAB(p2, p0)
	d012, d120, d201 := DotAB(c01, c12), DotAB(c12, c20), DotAB(c20, c01) // test each pair with same direction
	if d012 > 0 && d120 > 0 && d201 > 0 {                                 // all in the same direction
		return true // point is strictly inside the triangle
	} else {
		// 	if (Point.isZero(c01) && d120 < 0) return false;    // point is on side 01, but it's outside
		// 	if (Point.isZero(c12) && d201 < 0) return false;    // point is on side 12, but it's outside
		// 	if (Point.isZero(c20) && d012 < 0) return false;    // point is on side 20, but it's outside
		return false // point may on the border
	}
}
