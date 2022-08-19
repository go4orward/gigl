package gigl

import (
	"fmt"
	"strings"

	"github.com/go4orward/gigl/common"
	cst "github.com/go4orward/gigl/common/constants"
)

type GLShader interface {
	// Note that the creator NewShader() function should be implemented
	//   for each environment (like NewWebGLShader()/NewOpenGLShader()).
	CreateShaderProgram(vshader_source string, fshader_source string)
	GetShaderProgram() any
	IsReady() bool
	GetErr() error // returns any error during the creator function

	// setting up shader bindings
	SetBindingForUniform(btype cst.BindType, name string, target any)
	SetBindingForAttribute(btype cst.BindType, name string, target any)
	CheckBindings()
	GetUniformBindings() map[string]BindTarget
	GetAttributeBindings() map[string]BindTarget

	//
	Copy() GLShader
	String() string
	Summary() string
}

// ----------------------------------------------------------------------------
// Setting up Shader Bindings
// ----------------------------------------------------------------------------

type BindTarget struct {
	Type   cst.BindType // data type      of the variable
	Loc    any          // location       of the variable
	Target any          // target (value) of the variable
}

func (self *BindTarget) String() string {
	locinfo := "Loc:N"
	if self.Loc != nil {
		locinfo = "Loc:Y"
	}
	switch self.Target.(type) {
	case string:
		return fmt.Sprintf("{%2d %s %q}", self.Type, locinfo, self.Target.(string))
	case []float32:
		return fmt.Sprintf("{%2d %s %v}", self.Type, locinfo, self.Target.([]float32))
	case []int:
		return fmt.Sprintf("{%2d %s %v}", self.Type, locinfo, self.Target.([]int))
	default:
		return fmt.Sprintf("{%2d %s}", self.Type, locinfo)
	}
}

type GLShaderBinder struct {
	Uniforms   map[string]BindTarget // shader uniforms to bind
	Attributes map[string]BindTarget // shader attributes to bind
}

func (self *GLShaderBinder) InitBindings() {
	self.Uniforms = map[string]BindTarget{}
	self.Attributes = map[string]BindTarget{}
}

func (self *GLShaderBinder) SetBindingForUniform(btype cst.BindType, name string, target any) {
	// Set uniform binding with its name, data_type, and AUTO/MANUAL option.
	switch target.(type) {
	case string: // let Renderer bind the uniform variable automatically
		starget := target.(string)
		starget_split := strings.Split(starget, ":")
		starget0 := starget_split[0] // "material.texture:0" (with texture UNIT value)
		switch starget0 {
		case "lighting.dlight": // [mat3](3D) directional light information with (direction[3], color[3], ambient[3])
		case "material.color": //  [vec3] uniform color taken from Material
		case "material.texture": // [sampler2D] texture sampler(unit), like "material.texture:0"
		case "renderer.aspect": // AspectRatio of camera, Width : Height
		case "renderer.pvm": //  [mat3](2D) or [mat4](3D) (Proj * View * Model) matrix
		case "renderer.proj": // [mat3](2D) or [mat4](3D) (Projection) matrix
		case "renderer.vwmd": // [mat3](2D) or [mat4](3D) (View * Model) matrix
		default:
			common.Logger.Warn("Failed to SetBindingForUniform('%s') : unknown target '%s'\n", name, starget)
			return
		}
		self.Uniforms[name] = BindTarget{Type: btype, Target: starget}
	case []float32: // let Renderer set the uniform manually, with the given values
		self.Uniforms[name] = BindTarget{Type: btype, Target: target.([]float32)}
	default:
		common.Logger.Warn("Failed to SetBindingForUniform('%s') : invalid option %v\n", name, target)
	}
}

func (self *GLShaderBinder) SetBindingForAttribute(btype cst.BindType, name string, target any) {
	// Set attribute binding with its name, data_type, with AUTO_BINDING option.
	switch target.(type) {
	case string: // let Renderer bind the uniform variable automatically
		starget := target.(string)
		starget_split := strings.Split(starget, ":")
		starget0 := starget_split[0]
		switch starget0 {
		case "geometry.coords": // point coordinates
		case "geometry.textuv": // texture UV coordinates
		case "geometry.normal": // (3D only) normal vector
		case "instance.pose", "instance.color": // instance pose or color, like "instance.pose:<stride>:<offset>"
			if len(starget_split) != 3 {
				common.Logger.Warn("Failed to SetBindingForAttribute('%s') : try 'instance.pose:<stride>:<offset>'\n", name)
				return
			}
		default:
			common.Logger.Warn("Failed to SetBindingForAttribute('%s') : invalid target '%s'\n", name, starget)
			return
		}
		self.Attributes[name] = BindTarget{Type: btype, Target: starget}
	default:
		common.Logger.Warn("Failed to SetBindingForAttribute('%s') : invalid option %v\n", name, target)
		return
	}
}

func (self *GLShaderBinder) GetUniformBindings() map[string]BindTarget {
	return self.Uniforms
}

func (self *GLShaderBinder) GetAttributeBindings() map[string]BindTarget {
	return self.Attributes
}
