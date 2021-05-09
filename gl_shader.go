package gigl

type GLShader interface {
	String() string
	Copy() GLShader
	GetShaderProgram() interface{}
	GetUniformBindings() map[string]map[string]interface{}
	GetAttributeBindings() map[string]map[string]interface{}
	ShowInfo()

	// Bindings
	SetBindingForUniform(name string, dtype string, option interface{})
	SetBindingForAttribute(name string, dtype string, autobinding string)
	CheckBindings()
}
