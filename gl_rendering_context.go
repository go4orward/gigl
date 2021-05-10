package gigl

type GLRenderingContext interface {
	GetWH() [2]int
	GetConstants() *GLConstants
	GetEnvVariable(vname string, dtype string) interface{}

	// Material & Shader
	CreateMaterial(source string, options ...interface{}) (GLMaterial, error)
	CreateShader(vertex_shader string, fragment_shader string) (GLShader, error)

	// Data Buffer
	CreateDataBuffer(binding_target interface{}, data_slice interface{}) interface{}
	GLBindBuffer(binding_target interface{}, buffer interface{})

	// Texture
	GLCreateTexture() interface{}
	GLActiveTexture(texture_unit int)
	GLBindTexture(target interface{}, texture interface{})
	GLTexImage2DFromPixelBuffer(target interface{}, level int, internalformat interface{}, width int, height int, border int, format interface{}, dtype interface{}, pixels []uint8)
	GLTexImage2DFromImgObject(target interface{}, level int, internalformat interface{}, format interface{}, dtype interface{}, canvas interface{})
	GLGenerateMipmap(target interface{})
	GLTexParameteri(target interface{}, pname interface{}, pvalue interface{})

	// Uniforms
	GLUniform1i(location interface{}, v0 int)
	GLUniform1f(location interface{}, v0 float32)
	GLUniform2f(location interface{}, v0 float32, v1 float32)
	GLUniform3f(location interface{}, v0 float32, v1 float32, v2 float32)
	GLUniform4f(location interface{}, v0 float32, v1 float32, v2 float32, v3 float32)
	GLUniformMatrix3fv(location interface{}, transpose bool, values []float32)
	GLUniformMatrix4fv(location interface{}, transpose bool, values []float32)

	// Attributes
	GLVertexAttribPointer(location interface{}, size int, dtype interface{}, normalized bool, stride_in_byte int, offset_in_byte int)
	GLEnableVertexAttribArray(location interface{})
	GLVertexAttribDivisor(location interface{}, divisor int)

	// Preparing to Render
	GLClearColor(r float32, g float32, b float32, a float32)
	GLClear(mask interface{})
	GLEnable(cap interface{})
	GLDisable(cap interface{})
	GLDepthFunc(ftn interface{})
	GLBlendFunc(sfactor interface{}, dfactor interface{})
	GLUseProgram(program interface{})

	// Rendering
	GLDrawArrays(mode interface{}, first int, count int)
	GLDrawArraysInstanced(mode interface{}, first int, count int, pose_count int)
	GLDrawElements(mode interface{}, count int, dtype interface{}, offset int)
	GLDrawElementsInstanced(mode interface{}, element_count int, dtype interface{}, offset int, pose_count int)

	// WebGL Extensions
	SetupExtension(extname string)
	IsExtensionReady(extname string) bool
}
