package gigl

// ----------------------------------------------------------------------------
// GLRenderingContext
// ----------------------------------------------------------------------------

type GLRenderingContext interface {
	GetWH() [2]int
	GetConstants() *GLConstants
	GetEnvVariable(vname string, dtype string) interface{}

	// Material
	LoadMaterial(material GLMaterial) error
	SetupMaterial(material GLMaterial) error

	// Shader
	CreateShader(vertex_shader string, fragment_shader string) (GLShader, error)
	// DataBuffer
	CreateDataBufferVAO() *VAO
	CreateVtxDataBuffer(data_slice []float32) interface{}
	CreateIdxDataBuffer(data_slice []uint32) interface{}

	// Binding DataBuffer
	GLBindBuffer(binding_target uint32, buffer interface{})

	// Binding Texture
	GLActiveTexture(texture_unit int)
	GLBindTexture(target uint32, texture interface{})

	// Binding Uniforms
	GLUniform1i(location interface{}, v0 int)
	GLUniform1f(location interface{}, v0 float32)
	GLUniform2f(location interface{}, v0 float32, v1 float32)
	GLUniform3f(location interface{}, v0 float32, v1 float32, v2 float32)
	GLUniform4f(location interface{}, v0 float32, v1 float32, v2 float32, v3 float32)
	GLUniformMatrix3fv(location interface{}, transpose bool, values []float32)
	GLUniformMatrix4fv(location interface{}, transpose bool, values []float32)

	// Binding Attributes
	GLVertexAttribPointer(location interface{}, size int, dtype uint32, normalized bool, stride_in_byte int, offset_in_byte int)
	GLEnableVertexAttribArray(location interface{})
	GLVertexAttribDivisor(location interface{}, divisor int)

	// Preparing to Render
	GLClearColor(r float32, g float32, b float32, a float32)
	GLClear(mask uint32)
	GLEnable(cap uint32)
	GLDisable(cap uint32)
	GLDepthFunc(ftn uint32)
	GLBlendFunc(sfactor uint32, dfactor uint32)
	GLUseProgram(program interface{})

	// Rendering
	GLDrawArrays(mode uint32, first int, count int)
	GLDrawArraysInstanced(mode uint32, first int, count int, pose_count int)
	GLDrawElements(mode uint32, count int, dtype uint32, offset int)
	GLDrawElementsInstanced(mode uint32, element_count int, dtype uint32, offset int, pose_count int)

	// WebGL Extensions
	SetupExtension(extname string)
	IsExtensionReady(extname string) bool
}
