package gigl

type GLShader interface {
	// Note that the creator NewShader() function should be implemented
	//   for each environment (like NewWebGLShader()/NewOpenGLShader()).
	CreateShaderProgram(vshader_source string, fshader_source string)
	GetErr() error // returns any error during the creator function

	// setting up shader bindings
	SetBindingForUniform(name string, dtype string, option interface{})
	SetBindingForAttribute(name string, dtype string, autobinding string)
	CheckBindings()

	// using shader bindings
	GetShaderProgram() interface{}
	GetUniformBindings() map[string]map[string]interface{}
	GetAttributeBindings() map[string]map[string]interface{}

	//
	Copy() GLShader
	String() string
	ShowInfo()
}
