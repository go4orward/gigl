package common

type GLGeometry interface {
	IsDataBufferReady() bool
	IsWebGLBufferReady() bool
	BuildWebGLBuffers(rc GLRenderingContext, for_points bool, for_lines bool, for_faces bool)
	GetWebGLBuffer(draw_mode int) (interface{}, int, [4]int)
	ShowInfo()
}
