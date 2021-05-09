package common

import "math"

type Matrix4 struct {
	elements [16]float32 // COLUMN-MAJOR (just like WebGL)
}

func NewMatrix4() *Matrix4 {
	matrix := Matrix4{elements: [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}} // identity matrix
	return &matrix
}

func (self *Matrix4) GetElements() *[16]float32 {
	return &self.elements // reference
}

// ----------------------------------------------------------------------------
// Setting element values
// ----------------------------------------------------------------------------

func (self *Matrix4) Set(
	v00 float32, v01 float32, v02 float32, v03 float32,
	v10 float32, v11 float32, v12 float32, v13 float32,
	v20 float32, v21 float32, v22 float32, v23 float32,
	v30 float32, v31 float32, v32 float32, v33 float32) *Matrix4 {
	self.elements = [16]float32{ // COLUMN-MAJOR (just like WebGL)
		v00, v10, v20, v30,
		v01, v11, v21, v31,
		v02, v12, v22, v32,
		v03, v13, v23, v33}
	return self
}

func (self *Matrix4) SetIdentity() *Matrix4 {
	self.Set(1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1)
	return self
}

func (self *Matrix4) SetCopy(m *Matrix4) *Matrix4 {
	self.elements = m.elements
	return self
}

func (self *Matrix4) SetTranspose() *Matrix4 {
	e := &self.elements // reference
	e[1], e[4] = e[4], e[1]
	e[2], e[8] = e[8], e[2]   // [0], [4], [ 8], [12]
	e[3], e[12] = e[12], e[3] // [1], [5], [ 9], [13]
	e[6], e[9] = e[9], e[6]   // [2], [6], [10], [14]
	e[7], e[13] = e[13], e[7] // [3], [7], [11], [15]
	e[11], e[14] = e[14], e[11]
	return self
}

func (self *Matrix4) SetTranslation(tx float32, ty float32, tz float32) *Matrix4 {
	self.Set(
		1.0, 0.0, 0.0, tx,
		0.0, 1.0, 0.0, ty,
		0.0, 0.0, 1.0, tz,
		0.0, 0.0, 0.0, 1.0)
	return self
}

func (self *Matrix4) SetScaling(sx float32, sy float32, sz float32) *Matrix4 {
	self.Set(
		sx, 0.0, 0.0, 0,
		0.0, sy, 0.0, 0,
		0.0, 0.0, sz, 0,
		0.0, 0.0, 0.0, 1.0)
	return self
}

func (self *Matrix4) normalize_vector(v [3]float32) [3]float32 {
	len := float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])))
	return [3]float32{v[0] / len, v[1] / len, v[2] / len}
}

func (self *Matrix4) SetRotationByAxis(axis [3]float32, angle_in_degree float32) *Matrix4 {
	axis = self.normalize_vector(axis)
	// Based on http://www.gamedev.net/reference/articles/article1199.asp
	c := float32(math.Cos(float64(angle_in_degree) * (math.Pi / 180.0)))
	s := float32(math.Sin(float64(angle_in_degree) * (math.Pi / 180.0)))
	t := 1 - c
	x, y, z := axis[0], axis[1], axis[2]
	tx, ty := t*x, t*y
	self.Set(
		tx*x+c, tx*y-s*z, tx*z+s*y, 0,
		tx*y+s*z, ty*y+c, ty*z-s*x, 0,
		tx*z-s*y, ty*z+s*x, t*z*z+c, 0,
		0, 0, 0, 1)
	return self
}

func (self *Matrix4) SetMultiplyMatrices(matrices ...*Matrix4) *Matrix4 {
	if len(matrices) > 0 {
		m := matrices[0] // multiply all the matrices first,
		for i := 1; i < len(matrices); i++ {
			m = m.MultiplyToTheRight(matrices[i])
		}
		self.SetCopy(m) // and then copy (overwriting old values)
	}
	return self
}

// ----------------------------------------------------------------------------
// Creating new matrix
// ----------------------------------------------------------------------------

func (self *Matrix4) Copy() *Matrix4 {
	return &Matrix4{elements: self.elements}
}

func (self *Matrix4) Transpose() *Matrix4 {
	o := &self.elements // reference
	return &Matrix4{elements: [16]float32{
		o[0], o[1], o[2], o[3],
		o[4], o[5], o[6], o[7],
		o[8], o[9], o[10], o[11],
		o[12], o[13], o[14], o[15]}}
}

func (self *Matrix4) MultiplyToTheLeft(matrix *Matrix4) *Matrix4 {
	o := &self.elements   // reference       (M*O)T = OT * MT
	m := &matrix.elements // reference
	return &Matrix4{elements: [16]float32{
		o[0]*m[0] + o[1]*m[4] + o[2]*m[8] + o[3]*m[12],  // [0], [4], [8], [12]
		o[0]*m[1] + o[1]*m[5] + o[2]*m[9] + o[3]*m[13],  // [1], [5], [9], [13]
		o[0]*m[2] + o[1]*m[6] + o[2]*m[10] + o[3]*m[14], // [2], [6], [10], [14]
		o[0]*m[3] + o[1]*m[7] + o[2]*m[11] + o[3]*m[15], // [3], [7], [11], [15]
		o[4]*m[0] + o[5]*m[4] + o[6]*m[8] + o[7]*m[12],  // 2nd row
		o[4]*m[1] + o[5]*m[5] + o[6]*m[9] + o[7]*m[13],
		o[4]*m[2] + o[5]*m[6] + o[6]*m[10] + o[7]*m[14],
		o[4]*m[3] + o[5]*m[7] + o[6]*m[11] + o[7]*m[15],
		o[8]*m[0] + o[9]*m[4] + o[10]*m[8] + o[11]*m[12], // 3rd
		o[8]*m[1] + o[9]*m[5] + o[10]*m[9] + o[11]*m[13],
		o[8]*m[2] + o[9]*m[6] + o[10]*m[10] + o[11]*m[14],
		o[8]*m[3] + o[9]*m[7] + o[10]*m[11] + o[11]*m[15],
		o[12]*m[0] + o[13]*m[4] + o[14]*m[8] + o[15]*m[12], // 4th
		o[12]*m[1] + o[13]*m[5] + o[14]*m[9] + o[15]*m[13],
		o[12]*m[2] + o[13]*m[6] + o[14]*m[10] + o[15]*m[14],
		o[12]*m[3] + o[13]*m[7] + o[14]*m[11] + o[15]*m[15]}}
}

func (self *Matrix4) MultiplyToTheRight(matrix *Matrix4) *Matrix4 {
	o := &self.elements   // reference        (O*M)T = MT * OT
	m := &matrix.elements // reference
	return &Matrix4{elements: [16]float32{
		m[0]*o[0] + m[1]*o[4] + m[2]*o[8] + m[3]*o[12],  // [0], [4], [8], [12]
		m[0]*o[1] + m[1]*o[5] + m[2]*o[9] + m[3]*o[13],  // [1], [5], [9], [13]
		m[0]*o[2] + m[1]*o[6] + m[2]*o[10] + m[3]*o[14], // [2], [6], [10], [14]
		m[0]*o[3] + m[1]*o[7] + m[2]*o[11] + m[3]*o[15], // [3], [7], [11], [15]
		m[4]*o[0] + m[5]*o[4] + m[6]*o[8] + m[7]*o[12],  // 2nd row
		m[4]*o[1] + m[5]*o[5] + m[6]*o[9] + m[7]*o[13],
		m[4]*o[2] + m[5]*o[6] + m[6]*o[10] + m[7]*o[14],
		m[4]*o[3] + m[5]*o[7] + m[6]*o[11] + m[7]*o[15],
		m[8]*o[0] + m[9]*o[4] + m[10]*o[8] + m[11]*o[12], // 3rd
		m[8]*o[1] + m[9]*o[5] + m[10]*o[9] + m[11]*o[13],
		m[8]*o[2] + m[9]*o[6] + m[10]*o[10] + m[11]*o[14],
		m[8]*o[3] + m[9]*o[7] + m[10]*o[11] + m[11]*o[15],
		m[12]*o[0] + m[13]*o[4] + m[14]*o[8] + m[15]*o[12], // 4th
		m[12]*o[1] + m[13]*o[5] + m[14]*o[9] + m[15]*o[13],
		m[12]*o[2] + m[13]*o[6] + m[14]*o[10] + m[15]*o[14],
		m[12]*o[3] + m[13]*o[7] + m[14]*o[11] + m[15]*o[15]}}
}

// ----------------------------------------------------------------------------
// Handling Vector
// ----------------------------------------------------------------------------

func (self *Matrix4) MultiplyVector3(v [3]float32) [3]float32 {
	e := &self.elements // reference
	return [3]float32{
		e[0]*v[0] + e[4]*v[1] + e[8]*v[2] + e[12], // COLUMN-MAJOR
		e[1]*v[0] + e[5]*v[1] + e[9]*v[2] + e[13],
		e[2]*v[0] + e[6]*v[1] + e[10]*v[2] + e[14]}
}
