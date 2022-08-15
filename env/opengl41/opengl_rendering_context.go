package opengl41

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go4orward/gigl"
)

type OpenGLRenderingContext struct {
	constants gigl.GLConstants // OpenGL constant values
	wh        [2]int           // canvas width & height
}

func NewOpenGLRenderingContext(width int, height int) *OpenGLRenderingContext {
	self := OpenGLRenderingContext{}
	self.wh = [2]int{width, height}
	// get WebGL constants
	self.constants.ARRAY_BUFFER = gl.ARRAY_BUFFER
	self.constants.BLEND = gl.BLEND
	self.constants.BYTE = gl.BYTE
	self.constants.CLAMP_TO_EDGE = gl.CLAMP_TO_EDGE
	self.constants.COLOR_BUFFER_BIT = gl.COLOR_BUFFER_BIT
	self.constants.COMPILE_STATUS = gl.COMPILE_STATUS
	self.constants.DEPTH_BUFFER_BIT = gl.DEPTH_BUFFER_BIT
	self.constants.DEPTH_TEST = gl.DEPTH_TEST
	self.constants.ELEMENT_ARRAY_BUFFER = gl.ELEMENT_ARRAY_BUFFER
	self.constants.FLOAT = gl.FLOAT
	self.constants.FRAGMENT_SHADER = gl.FRAGMENT_SHADER
	self.constants.LEQUAL = gl.LEQUAL
	self.constants.LESS = gl.LESS
	self.constants.LINEAR = gl.LINEAR
	self.constants.LINES = gl.LINES
	self.constants.LINK_STATUS = gl.LINK_STATUS
	self.constants.NEAREST = gl.NEAREST
	self.constants.ONE = gl.ONE
	self.constants.ONE_MINUS_SRC_ALPHA = gl.ONE_MINUS_SRC_ALPHA
	self.constants.POINTS = gl.POINTS
	self.constants.RGBA = gl.RGBA
	self.constants.SRC_ALPHA = gl.SRC_ALPHA
	self.constants.STATIC_DRAW = gl.STATIC_DRAW
	self.constants.TEXTURE_2D = gl.TEXTURE_2D
	self.constants.TEXTURE0 = gl.TEXTURE0
	self.constants.TEXTURE1 = gl.TEXTURE1
	self.constants.TEXTURE_MIN_FILTER = gl.TEXTURE_MIN_FILTER
	self.constants.TEXTURE_WRAP_S = gl.TEXTURE_WRAP_S
	self.constants.TEXTURE_WRAP_T = gl.TEXTURE_WRAP_T
	self.constants.TRIANGLES = gl.TRIANGLES
	self.constants.UNSIGNED_BYTE = gl.UNSIGNED_BYTE
	self.constants.UNSIGNED_INT = gl.UNSIGNED_INT
	self.constants.UNSIGNED_SHORT = gl.UNSIGNED_SHORT
	self.constants.VERTEX_SHADER = gl.VERTEX_SHADER
	return &self
}

func (self *OpenGLRenderingContext) GetWH() [2]int {
	return self.wh
}

func (self *OpenGLRenderingContext) GetConstants() *gigl.GLConstants {
	return &self.constants
}

func (self *OpenGLRenderingContext) GetEnvVariable(vname string, dtype string) interface{} {
	switch dtype {
	case "int":
		return 0
	case "bool":
		return false
	default:
		return ""
	}
}

// ----------------------------------------------------------------------------
// Material & Shader
// ----------------------------------------------------------------------------

func (self *OpenGLRenderingContext) LoadMaterial(material gigl.GLMaterial) error {
	return load_material(self, material)
}

func (self *OpenGLRenderingContext) CreateShader(vertex_shader string, fragment_shader string) (gigl.GLShader, error) {
	return create_shader(self, vertex_shader, fragment_shader)
}

// ----------------------------------------------------------------------------
// OpenGL Data Buffer
// ----------------------------------------------------------------------------

func (self *OpenGLRenderingContext) CreateDataBufferVAO() *gigl.VAO {
	var vao uint32 // NON-ZERO values
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	// fmt.Printf("VAO (%T): %v\n", vao, vao)
	return &gigl.VAO{}
}

func (self *OpenGLRenderingContext) CreateVtxDataBuffer(data_slice []float32) interface{} {
	if data_slice == nil {
		return nil
	}
	// var vao uint32
	// gl.GenVertexArrays(1, &vao)
	// gl.BindVertexArray(vao)
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(data_slice)*4, gl.Ptr(data_slice), gl.STATIC_DRAW)
	// fmt.Printf("  CreateDataBuffer() : type:%T  len:%d  =>  vbo:%v\n", data_slice, len(data_buffer), vbo)
	return vbo
}

func (self *OpenGLRenderingContext) CreateIdxDataBuffer(data_slice []uint32) interface{} {
	if data_slice == nil {
		return nil
	}
	// var vao uint32
	// gl.GenVertexArrays(1, &vao)
	// gl.BindVertexArray(vao)
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(data_slice)*4, gl.Ptr(data_slice), gl.STATIC_DRAW)
	// fmt.Printf("  CreateDataBuffer() : type:%T  len:%d  =>  vbo:%v\n", data_slice, len(data_buffer), vbo)
	return vbo
}

func (self *OpenGLRenderingContext) GLBindBuffer(target uint32, buffer interface{}) {
	// 'bind_target' : c.ARRAY_BUFFER or c.ELEMENT_ARRAY_BUFFER
	if buffer == nil {
		gl.BindBuffer(target, 0)
	} else {
		gl.BindBuffer(target, buffer.(uint32))
	}
}

// ----------------------------------------------------------------------------
// Binding Texture
// ----------------------------------------------------------------------------

func (self *OpenGLRenderingContext) GLActiveTexture(texture_unit int) {
	// js_texture_unit := js.ValueOf(self.constants.TEXTURE0.(js.Value).Int() + texture_unit)
	// gl.ActiveTexture(texture_unit)
}

func (self *OpenGLRenderingContext) GLBindTexture(target uint32, texture interface{}) {
	// 'binding_target' : TEXTURE_2D
	// self.context.Call("bindTexture", target.(js.Value), texture.(js.Value))
}

// ----------------------------------------------------------------------------
// Binding Uniforms
// ----------------------------------------------------------------------------

func (self *OpenGLRenderingContext) GLUniform1i(location interface{}, v0 int) {
	gl.Uniform1i(location.(int32), int32(v0))
}

func (self *OpenGLRenderingContext) GLUniform1f(location interface{}, v0 float32) {
	gl.Uniform1f(location.(int32), v0)
}

func (self *OpenGLRenderingContext) GLUniform2f(location interface{}, v0 float32, v1 float32) {
	gl.Uniform2f(location.(int32), v0, v1)
}

func (self *OpenGLRenderingContext) GLUniform3f(location interface{}, v0 float32, v1 float32, v2 float32) {
	gl.Uniform3f(location.(int32), v0, v1, v2)
}

func (self *OpenGLRenderingContext) GLUniform4f(location interface{}, v0 float32, v1 float32, v2 float32, v3 float32) {
	gl.Uniform4f(location.(int32), v0, v1, v2, v3)
}

func (self *OpenGLRenderingContext) GLUniformMatrix3fv(location interface{}, transpose bool, values []float32) {
	gl.UniformMatrix3fv(location.(int32), 1, transpose, &values[0])
}

func (self *OpenGLRenderingContext) GLUniformMatrix4fv(location interface{}, transpose bool, values []float32) {
	// js_typed_array := self.ConvertGoSliceToJsTypedArray(values) // converted to JavaScript 'Float32Array'
	gl.UniformMatrix4fv(location.(int32), 1, transpose, &values[0])
}

// ----------------------------------------------------------------------------
// Binding Attributes
// ----------------------------------------------------------------------------

func (self *OpenGLRenderingContext) GLVertexAttribPointer(location interface{}, size int, dtype uint32, normalized bool, stride_in_byte int, offset_in_byte int) {
	gl.VertexAttribPointer(uint32(location.(int32)), int32(size), dtype, normalized, int32(stride_in_byte), gl.PtrOffset(offset_in_byte))
	// gl.VertexAttribPointerWithOffset(location.(int32), int32(size), dtype, normalized, int32(stride_in_byte), offset_in_byte)
}

func (self *OpenGLRenderingContext) GLEnableVertexAttribArray(location interface{}) {
	gl.EnableVertexAttribArray(uint32(location.(int32)))
}

func (self *OpenGLRenderingContext) GLVertexAttribDivisor(location interface{}, divisor int) {
	// extension for geometry instancing
	// if !self.ext_angle.IsNull() {
	// 	self.ext_angle.Call("vertexAttribDivisorANGLE", location.(int32), divisor)
	// }
}

// ----------------------------------------------------------------------------
// Preparing to Render
// ----------------------------------------------------------------------------

func (self *OpenGLRenderingContext) GLClearColor(r float32, g float32, b float32, a float32) {
	gl.ClearColor(r, g, b, a)
}

func (self *OpenGLRenderingContext) GLClear(mask uint32) {
	gl.Clear(mask)
}

func (self *OpenGLRenderingContext) GLEnable(cap uint32) {
	gl.Enable(cap)
}

func (self *OpenGLRenderingContext) GLDisable(cap uint32) {
	gl.Disable(cap)
}

func (self *OpenGLRenderingContext) GLDepthFunc(ftn uint32) {
	gl.DepthFunc(ftn)
}

func (self *OpenGLRenderingContext) GLBlendFunc(sfactor uint32, dfactor uint32) {
	gl.BlendFunc(sfactor, dfactor)
}

func (self *OpenGLRenderingContext) GLUseProgram(shader_program interface{}) {
	gl.UseProgram(shader_program.(uint32))
}

// ----------------------------------------------------------------------------
// Rendering
// ----------------------------------------------------------------------------

func (self *OpenGLRenderingContext) GLDrawArrays(mode uint32, first int, count int) {
	// 'mode' : POINTS
	gl.DrawArrays(mode, int32(first), int32(count))
}

func (self *OpenGLRenderingContext) GLDrawArraysInstanced(mode uint32, first int, count int, pose_count int) {
	// 'mode' : POINTS
	// if !self.ext_angle.IsNull() {
	// 	self.ext_angle.Call("drawArraysInstancedANGLE", mode.(js.Value), first, count, pose_count)
	// }
}

func (self *OpenGLRenderingContext) GLDrawElements(mode uint32, count int, dtype uint32, offset int) {
	// 'mode'  : LINES, TRIANGLES
	// 'dtype' : UNSIGNED_INT
	gl.DrawElements(mode, int32(count), dtype, gl.PtrOffset(offset))
}

func (self *OpenGLRenderingContext) GLDrawElementsInstanced(mode uint32, element_count int, dtype uint32, offset int, pose_count int) {
	// 'mode'  : LINES, TRIANGLES
	// 'dtype' : UNSIGNED_INT
	// if !self.ext_angle.IsNull() {
	// 	self.ext_angle.Call("drawElementsInstancedANGLE", mode.(js.Value), element_count, dtype.(js.Value), offset, pose_count)
	// }
}

// ----------------------------------------------------------------------------
// OpenGL Extensions
// ----------------------------------------------------------------------------

func (self *OpenGLRenderingContext) SetupExtension(extname string) {
	// switch extname {
	// case "UINT32": // extension for UINT32 index, to drawElements() with large number of vertices
	// 	self.ext_uint = self.context.Call("getExtension", "OES_element_index_uint")
	// case "ANGLE": // extension for geometry instancing
	// 	self.ext_angle = self.context.Call("getExtension", "ANGLE_instanced_arrays")
	// }
}

func (self *OpenGLRenderingContext) IsExtensionReady(extname string) bool {
	// switch extname {
	// case "UINT32": // extension for UINT32 index, to drawElements() with large number of vertices
	// 	return !self.ext_uint.IsNull() && !self.ext_uint.IsUndefined()
	// case "ANGLE": // extension for geometry instancing
	// 	return !self.ext_angle.IsNull() && !self.ext_angle.IsUndefined()
	// }
	return false
}

// ----------------------------------------------------------------------------
// private functions
// ----------------------------------------------------------------------------
