package constants

type BindType uint8

const (
	None      BindType = 0
	Vec1      BindType = 1
	Vec2      BindType = 2
	Vec3      BindType = 3
	Vec4      BindType = 4
	Mat3      BindType = 9
	Mat4      BindType = 16
	Sampler2D BindType = 99
)
