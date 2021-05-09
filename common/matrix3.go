package common

import "math"

type Matrix3 struct {
	elements [9]float32 // COLUMN-MAJOR (just like WebGL)
}

func NewMatrix3() *Matrix3 {
	mtx := Matrix3{elements: [9]float32{1, 0, 0, 0, 1, 0, 0, 0, 1}} // identity matrix
	return &mtx
}

func (self *Matrix3) GetElements() *[9]float32 {
	return &self.elements // reference
}

// ----------------------------------------------------------------------------
// Setting element values
// ----------------------------------------------------------------------------

func (self *Matrix3) Set(v00 float32, v01 float32, v02 float32, v10 float32, v11 float32, v12 float32, v20 float32, v21 float32, v22 float32) *Matrix3 {
	self.elements = [9]float32{ // COLUMN-MAJOR (just like WebGL)
		v00, v10, v20,
		v01, v11, v21,
		v02, v12, v22}
	return self
}

func (self *Matrix3) SetIdentity() *Matrix3 {
	self.Set(1, 0, 0, 0, 1, 0, 0, 0, 1)
	return self
}

func (self *Matrix3) SetCopy(m *Matrix3) *Matrix3 {
	self.elements = m.elements // copy values
	return self
}

func (self *Matrix3) SetTranspose() *Matrix3 {
	e := &self.elements     // reference
	e[1], e[3] = e[3], e[1] // [0], [1], [2]
	e[2], e[6] = e[6], e[2] // [3], [4], [5]
	e[5], e[7] = e[7], e[5] // [6], [7], [8]
	return self
}

func (self *Matrix3) SetTranslation(tx float32, ty float32) *Matrix3 {
	self.Set(
		1.0, 0.0, tx,
		0.0, 1.0, ty,
		0.0, 0.0, 1.0)
	return self
}

func (self *Matrix3) SetScaling(sx float32, sy float32) *Matrix3 {
	self.Set(
		sx, 0.0, 0,
		0.0, sy, 0,
		0.0, 0.0, 1.0)
	return self
}

func (self *Matrix3) SetRotation(angle_in_degree float32) *Matrix3 {
	// Based on http://www.gamedev.net/reference/articles/article1199.asp
	cos := float32(math.Cos(float64(angle_in_degree) * (math.Pi / 180.0)))
	sin := float32(math.Sin(float64(angle_in_degree) * (math.Pi / 180.0)))
	self.Set(
		cos, -sin, 0.0,
		+sin, cos, 0.0,
		0.0, 0.0, 1.0)
	return self
}

func (self *Matrix3) SetMultiplyMatrices(matrices ...*Matrix3) *Matrix3 {
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

func (self *Matrix3) Copy() *Matrix3 {
	return &Matrix3{elements: self.elements}
}

func (self *Matrix3) Transpose() *Matrix3 {
	o := &self.elements // reference
	return &Matrix3{elements: [9]float32{o[0], o[1], o[2], o[3], o[4], o[5], o[6], o[7], o[8]}}
}

func (self *Matrix3) MultiplyToTheLeft(matrix *Matrix3) *Matrix3 {
	o := &self.elements   // reference
	m := &matrix.elements // reference
	return &Matrix3{elements: [9]float32{
		o[0]*m[0] + o[1]*m[3] + o[2]*m[6],
		o[0]*m[1] + o[1]*m[4] + o[2]*m[7],
		o[0]*m[2] + o[1]*m[5] + o[2]*m[8],
		o[3]*m[0] + o[4]*m[3] + o[5]*m[6],
		o[3]*m[1] + o[4]*m[4] + o[5]*m[7],
		o[3]*m[2] + o[4]*m[5] + o[5]*m[8],
		o[6]*m[0] + o[7]*m[3] + o[8]*m[6],
		o[6]*m[1] + o[7]*m[4] + o[8]*m[7],
		o[6]*m[2] + o[7]*m[5] + o[8]*m[8]}}
}

func (self *Matrix3) MultiplyToTheRight(matrix *Matrix3) *Matrix3 {
	o := &self.elements   // reference
	m := &matrix.elements // reference
	return &Matrix3{elements: [9]float32{
		m[0]*o[0] + m[1]*o[3] + m[2]*o[6],
		m[0]*o[1] + m[1]*o[4] + m[2]*o[7],
		m[0]*o[2] + m[1]*o[5] + m[2]*o[8],
		m[3]*o[0] + m[4]*o[3] + m[5]*o[6],
		m[3]*o[1] + m[4]*o[4] + m[5]*o[7],
		m[3]*o[2] + m[4]*o[5] + m[5]*o[8],
		m[6]*o[0] + m[7]*o[3] + m[8]*o[6],
		m[6]*o[1] + m[7]*o[4] + m[8]*o[7],
		m[6]*o[2] + m[7]*o[5] + m[8]*o[8]}}
}

// ----------------------------------------------------------------------------
// Handling Vector
// ----------------------------------------------------------------------------

func (self *Matrix3) MultiplyVector2(v [2]float32) [2]float32 {
	e := &self.elements // reference
	return [2]float32{
		e[0]*v[0] + e[3]*v[1] + e[6], // COLUMN-MAJOR
		e[1]*v[0] + e[4]*v[1] + e[7]}
}
