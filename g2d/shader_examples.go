package g2d

import (
	"github.com/go4orward/gigl"
	cst "github.com/go4orward/gigl/common/constants"
)

func NewShaderFor2DAxes(rc gigl.GLRenderingContext) gigl.GLShader {
	// Shader for three axes - X(RED) & Y(GREEN) - for visual reference,
	//   with auto-binded (Proj * View * Model) matrix and XY coordinates
	var vertex_shader_code = `
		precision mediump float;
		uniform   mat3 pvm;			// Projection * View * Model matrix
		attribute vec2 axy;			// XY coordinates
		varying   vec2 v_xy;		// (varying) XY coordinates
		void main() {
			vec3 new_pos = pvm * vec3(axy, 1.0);
			gl_Position = vec4(new_pos.x, new_pos.y, 0.0, 1.0);
			v_xy = axy;
		}`
	var fragment_shader_code = `
		precision mediump float;
		varying vec2 v_xy;			// (varying) XY coordinates
		void main() {
			if (v_xy.x != 0.0) gl_FragColor = vec4(1.0, 0.1, 0.1, 1.0);
			else               gl_FragColor = vec4(0.1, 1.0, 0.1, 1.0);
		}`
	shader, _ := rc.CreateShader(vertex_shader_code, fragment_shader_code)
	shader.SetBindingForUniform(cst.Mat3, "pvm", "renderer.pvm")      // Proj*View*Model matrix
	shader.SetBindingForAttribute(cst.Vec2, "axy", "geometry.coords") // point coordinates
	shader.CheckBindings()                                            // check validity of the shader
	return shader
}

func NewShaderForMaterialColors(rc gigl.GLRenderingContext) gigl.GLShader {
	// Shader with auto-binded color and (Proj * View * Model) matrix
	var vertex_shader_code = `
		precision mediump float;
		uniform   mat3 pvm;			// Projection * View * Model matrix
		attribute vec2 xy;			// XY coordinates
		void main() {
			vec3 new_pos = pvm * vec3(xy.x, xy.y, 1.0);
			gl_Position = vec4(new_pos.x, new_pos.y, 0.0, 1.0);
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform vec4 color;			// color RGBA
		void main() { 
			gl_FragColor = color;
		}`
	shader, _ := rc.CreateShader(vertex_shader_code, fragment_shader_code)
	shader.SetBindingForUniform(cst.Mat3, "pvm", "renderer.pvm")     // Proj*View*Model matrix
	shader.SetBindingForUniform(cst.Vec4, "color", "material.color") // material color
	shader.SetBindingForAttribute(cst.Vec2, "xy", "geometry.coords") // point coordinates
	shader.CheckBindings()                                           // check validity of the shader
	return shader
}

func NewShaderForMaterialTexture(rc gigl.GLRenderingContext) gigl.GLShader {
	// Shader with auto-binded color and (Proj * View * Model) matrix
	var vertex_shader_code = `
		precision mediump float;
		uniform   mat3 pvm;			// Projection * View * Model matrix
		attribute vec2 xy;			// XY coordinates
		attribute vec2 uv;			// UV coordinates
		varying vec2 v_uv;			// (varying) UV coordinates
		void main() {
			vec3 new_pos = pvm * vec3(xy.x, xy.y, 1.0);
			gl_Position = vec4(new_pos.x, new_pos.y, 0.0, 1.0);
			v_uv = uv;
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform sampler2D text;		// texture sampler (unit)
		varying vec2 v_uv;			// (varying) UV coordinates
		void main() { 
			gl_FragColor = texture2D(text, v_uv);
		}`
	shader, _ := rc.CreateShader(vertex_shader_code, fragment_shader_code)
	shader.SetBindingForUniform(cst.Mat3, "pvm", "renderer.pvm")           // Proj*View*Model matrix
	shader.SetBindingForUniform(cst.Sampler2D, "text", "material.texture") // texture sampler (unit:0)
	shader.SetBindingForAttribute(cst.Vec2, "xy", "geometry.coords")       // point coordinates
	shader.SetBindingForAttribute(cst.Vec2, "uv", "geometry.textuv")       // texture UV coordinates
	shader.CheckBindings()                                                 // check validity of the shader
	return shader
}

func NewShaderForInstancePoseColor(rc gigl.GLRenderingContext) gigl.GLShader {
	// Shader with instance pose, for rendering multiple instances of a same geometry
	var vertex_shader_code = `
		precision mediump float;
		uniform   mat3 pvm;			// Projection * View * Model matrix
		attribute vec2 xy;			// XY coordinates
		attribute vec2 ixy;			// instance position : XY translation
		attribute vec3 icolor;		// instance color    : RGB
		varying   vec3 v_color;		// (varying) color
		void main() {
			vec3 new_pos = pvm * vec3(xy.x + ixy.x, xy.y + ixy.y, 1.0);
			gl_Position = vec4(new_pos.x, new_pos.y, 0.0, 1.0);
			v_color = icolor;
		}`
	var fragment_shader_code = `
		precision mediump float;
		varying   vec3 v_color;		// (varying) color
		void main() { 
			gl_FragColor = vec4(v_color, 1.0);
		}`
	shader, _ := rc.CreateShader(vertex_shader_code, fragment_shader_code)
	shader.SetBindingForUniform(cst.Mat3, "pvm", "renderer.pvm")           // Proj*View*Model matrix
	shader.SetBindingForAttribute(cst.Vec2, "xy", "geometry.coords")       // point coordinates
	shader.SetBindingForAttribute(cst.Vec2, "ixy", "instance.pose:5:0")    // instance pose (XY coordinates)
	shader.SetBindingForAttribute(cst.Vec3, "icolor", "instance.pose:5:2") // instance color (RGB packed in single float32)
	shader.CheckBindings()                                                 // check validity of the shader
	return shader
}
