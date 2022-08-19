package g3d

import (
	"github.com/go4orward/gigl"
	cst "github.com/go4orward/gigl/common/constants"
)

func NewShader_3DAxes(rc gigl.GLRenderingContext) gigl.GLShader {
	var vertex_shader_code = `
		precision mediump float;
		uniform mat4 pvm;			// Projection * View * Model matrix
		attribute vec3 xyz;			// XYZ coordinates
		varying vec3 v_xyz;			// (varying) XYZ coordinates
		void main() {
			gl_Position = pvm * vec4(xyz.x, xyz.y, xyz.z, 1.0);
			v_xyz = xyz;
		}`
	var fragment_shader_code = `
		precision mediump float;
		varying vec3 v_xyz;			// (varying) XYZ coordinates
		void main() {
			if      (v_xyz.x != 0.0) gl_FragColor = vec4(1.0, 0.1, 0.1, 1.0);
			else if (v_xyz.y != 0.0) gl_FragColor = vec4(0.1, 1.0, 0.1, 1.0);
			else                     gl_FragColor = vec4(0.6, 0.6, 1.0, 1.0);
		}`
	shader, _ := rc.CreateShader(vertex_shader_code, fragment_shader_code)
	shader.SetBindingForUniform(cst.Mat4, "pvm", "renderer.pvm")      // (Proj * View * Models) matrix
	shader.SetBindingForAttribute(cst.Vec3, "xyz", "geometry.coords") // vertex coordinates
	shader.CheckBindings()                                            // check validity of the shader
	return shader
}

func NewShader_ColorOnly(rc gigl.GLRenderingContext) gigl.GLShader {
	// Shader for (XYZ + NORMAL) Geometry & (COLOR) Material & (DIRECTIONAL) Lighting
	var vertex_shader_code = `
		precision mediump float;
		uniform mat4 pvm;			// Projection * View * Model matrix
		attribute vec3 xyz;			// XYZ coordinates
		void main() {
			gl_Position = pvm * vec4(xyz, 1.0);
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform vec3 color;			// single color
		void main() { 
			gl_FragColor = vec4(color.rgb, 1.0);
		}`
	shader, _ := rc.CreateShader(vertex_shader_code, fragment_shader_code)
	shader.SetBindingForUniform(cst.Mat4, "pvm", "renderer.pvm")      // (Proj * View * Models) matrix
	shader.SetBindingForUniform(cst.Vec3, "color", "material.color")  // material color
	shader.SetBindingForAttribute(cst.Vec3, "xyz", "geometry.coords") // point XYZ coordinates
	shader.CheckBindings()                                            // check validity of the shader
	return shader
}

func NewShader_NormalColor(rc gigl.GLRenderingContext) gigl.GLShader {
	// Shader for (XYZ + NORMAL) Geometry & (COLOR) Material & (DIRECTIONAL) Lighting
	var vertex_shader_code = `
		precision mediump float;
		uniform mat4 proj;			// Projection matrix
		uniform mat4 vwmd;			// ModelView matrix
		uniform mat3 light;			// directional light ([0]:direction, [1]:color, [2]:ambient) COLUMN-MAJOR!
		attribute vec3 xyz;			// XYZ coordinates
		attribute vec3 nor;			// normal vector
		varying vec3 v_light;   	// (varying) lighting intensity for the point
		void main() {
			gl_Position = proj * vwmd * vec4(xyz.x, xyz.y, xyz.z, 1.0);
			float s = sqrt( vwmd[0][0]*vwmd[0][0] + vwmd[0][1]*vwmd[0][1] + vwmd[0][2]*vwmd[0][2]);  // scaling 
			mat3  mvRot = mat3( vwmd[0][0]/s, vwmd[0][1]/s, vwmd[0][2]/s, vwmd[1][0]/s, vwmd[1][1]/s, vwmd[1][2]/s, vwmd[2][0]/s, vwmd[2][1]/s, vwmd[2][2]/s );
			vec3  normal    = mvRot * nor;               		// normal vector in camera space
			float intensity = max(dot(normal, light[0]), 0.0);	// light_intensity = dot(face_normal,light_direction)
			v_light = intensity * light[1] + light[2];        	// intensity * light_color + ambient_color
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform vec4 color;			// material color
		varying vec3 v_light;		// (varying) lighting intensity
		void main() { 
			gl_FragColor = vec4(color.rgb * v_light, color.a);
		}`
	shader, _ := rc.CreateShader(vertex_shader_code, fragment_shader_code)
	shader.SetBindingForUniform(cst.Mat4, "proj", "renderer.proj")    // (Projection) matrix
	shader.SetBindingForUniform(cst.Mat4, "vwmd", "renderer.vwmd")    // (View * Models) matrix
	shader.SetBindingForUniform(cst.Vec4, "color", "material.color")  // material color
	shader.SetBindingForUniform(cst.Mat3, "light", "lighting.dlight") // directional lighting
	shader.SetBindingForAttribute(cst.Vec3, "xyz", "geometry.coords") // point XYZ coordinates
	shader.SetBindingForAttribute(cst.Vec3, "nor", "geometry.normal") // point normal vectors
	shader.CheckBindings()                                            // check validity of the shader
	return shader
}

func NewShader_TextureOnly(rc gigl.GLRenderingContext) gigl.GLShader {
	// Shader for (XYZ + UV + NORMAL) Geometry & (TEXTURE) Material & (DIRECTIONAL) Lighting
	var vertex_shader_code = `
		precision mediump float;
		uniform mat4 proj;			// Projection matrix
		uniform mat4 vwmd;			// ModelView matrix
		attribute vec3 xyz;			// XYZ coordinates
		attribute vec2 tuv;			// texture coordinates
		varying vec2 v_tuv;			// (varying) texture coordinates
		void main() {
			gl_Position = proj * vwmd * vec4(xyz.x, xyz.y, xyz.z, 1.0);
			v_tuv = tuv;
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform sampler2D text;		// texture sampler (unit)
		varying vec2 v_tuv;			// (varying) texture coordinates
		void main() { 
			gl_FragColor = texture2D(text, v_tuv);
		}`
	shader, _ := rc.CreateShader(vertex_shader_code, fragment_shader_code)
	shader.SetBindingForUniform(cst.Mat4, "proj", "renderer.proj")         // (Projection) matrix
	shader.SetBindingForUniform(cst.Mat4, "vwmd", "renderer.vwmd")         // (View * Models) matrix
	shader.SetBindingForUniform(cst.Sampler2D, "text", "material.texture") // texture sampler (unit:0)
	shader.SetBindingForAttribute(cst.Vec3, "xyz", "geometry.coords")      // point XYZ coordinates
	shader.SetBindingForAttribute(cst.Vec2, "tuv", "geometry.textuv")      // point UV coordinates (texture)
	shader.CheckBindings()                                                 // check validity of the shader
	return shader
}

func NewShader_NormalTexture(rc gigl.GLRenderingContext) gigl.GLShader {
	// Shader for (XYZ + UV + NORMAL) Geometry & (TEXTURE) Material & (DIRECTIONAL) Lighting
	var vertex_shader_code = `
		precision mediump float;
		uniform mat4 proj;			// Projection matrix
		uniform mat4 vwmd;			// ModelView matrix
		uniform mat3 light;			// directional light ([0]:direction, [1]:color, [2]:ambient) COLUMN-MAJOR!
		attribute vec3 xyz;			// XYZ coordinates
		attribute vec2 tuv;			// texture coordinates
		attribute vec3 nor;			// normal vector
		varying vec2 v_tuv;			// (varying) texture coordinates
		varying vec3 v_light;		// (varying) lighting intensity for the point
		void main() {
			gl_Position = proj * vwmd * vec4(xyz.x, xyz.y, xyz.z, 1.0);
			float s = sqrt( vwmd[0][0]*vwmd[0][0] + vwmd[0][1]*vwmd[0][1] + vwmd[0][2]*vwmd[0][2]);  // scaling 
			mat3  mvRot = mat3( vwmd[0][0]/s, vwmd[0][1]/s, vwmd[0][2]/s, vwmd[1][0]/s, vwmd[1][1]/s, vwmd[1][2]/s, vwmd[2][0]/s, vwmd[2][1]/s, vwmd[2][2]/s );
			vec3  normal    = mvRot * nor;               		// normal vector in camera space
			float intensity = max(dot(normal, light[0]), 0.0);	// light_intensity = dot(face_normal,light_direction)
			v_light = intensity * light[1] + light[2];        	// intensity * light_color + ambient_color
			v_tuv = tuv;
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform sampler2D text;		// texture sampler (unit)
		varying vec2 v_tuv;			// (varying) texture coordinates
		varying vec3 v_light;		// (varying) lighting intensity
		void main() { 
			vec4 color = texture2D(text, v_tuv);
			gl_FragColor = vec4(color.rgb * v_light, color.a);
		}`
	shader, _ := rc.CreateShader(vertex_shader_code, fragment_shader_code)
	shader.SetBindingForUniform(cst.Mat4, "proj", "renderer.proj")         // (Projection) matrix
	shader.SetBindingForUniform(cst.Mat4, "vwmd", "renderer.vwmd")         // (View * Models) matrix
	shader.SetBindingForUniform(cst.Mat3, "light", "lighting.dlight")      // directional lighting
	shader.SetBindingForUniform(cst.Sampler2D, "text", "material.texture") // texture sampler (unit:0)
	shader.SetBindingForAttribute(cst.Vec3, "xyz", "geometry.coords")      // point XYZ coordinates
	shader.SetBindingForAttribute(cst.Vec2, "tuv", "geometry.textuv")      // point UV coordinates (texture)
	shader.SetBindingForAttribute(cst.Vec3, "nor", "geometry.normal")      // point normal vector
	shader.CheckBindings()                                                 // check validity of the shader
	return shader
}

func NewShader_InstancePoseColor(rc gigl.GLRenderingContext) gigl.GLShader {
	// Shader for (XYZ + NORMAL) Geometry & (COLOR) Material & (DIRECTIONAL) Lighting
	var vertex_shader_code = `
		precision mediump float;
		uniform mat4 proj;			// Projection matrix
		uniform mat4 vwmd;			// ModelView matrix
		attribute vec3 xyz;			// XYZ coordinates
		attribute vec3 nor;			// normal vector
		attribute vec3 ixyz;		// instance pose : XYZ translation
		attribute vec3 icolor;		// instance pose : color
		uniform mat3 light;			// [0]: direction, [1]: color, [2]: ambient_color   (column-major)
		varying vec3 v_color;    	// (varying) instance color
		varying vec3 v_light;    	// (varying) lighting intensity
		void main() {
			gl_Position = proj * vwmd * vec4(xyz.x + ixyz[0], xyz.y + ixyz[1], xyz.z + ixyz[2], 1.0);
			float s = sqrt( vwmd[0][0]*vwmd[0][0] + vwmd[0][1]*vwmd[0][1] + vwmd[0][2]*vwmd[0][2]);  // scaling 
			mat3  mvRot = mat3( vwmd[0][0]/s, vwmd[0][1]/s, vwmd[0][2]/s, vwmd[1][0]/s, vwmd[1][1]/s, vwmd[1][2]/s, vwmd[2][0]/s, vwmd[2][1]/s, vwmd[2][2]/s );
			vec3  normal    = mvRot * nor;               		// normal vector in camera space
			float intensity = max(dot(normal, light[0]), 0.0);	// light_intensity = dot(face_normal,light_direction)
			v_light = intensity * light[1] + light[2];        	// intensity * light_color + ambient_color
			v_color = icolor;
		}`
	var fragment_shader_code = `
		precision mediump float;
		varying vec3 v_color;		// (varying) instance color
		varying vec3 v_light;		// (varying) lighting intensity
		void main() { 
			gl_FragColor = vec4(v_color * v_light, 1.0);
		}`
	shader, _ := rc.CreateShader(vertex_shader_code, fragment_shader_code)
	shader.SetBindingForUniform(cst.Mat4, "proj", "renderer.proj")         // (Projection) matrix
	shader.SetBindingForUniform(cst.Mat4, "vwmd", "renderer.vwmd")         // (View * Models) matrix
	shader.SetBindingForUniform(cst.Mat3, "light", "lighting.dlight")      // directional lighting
	shader.SetBindingForAttribute(cst.Vec3, "xyz", "geometry.coords")      // point XYZ coordinates
	shader.SetBindingForAttribute(cst.Vec3, "nor", "geometry.normal")      // point normal vectors
	shader.SetBindingForAttribute(cst.Vec3, "ixyz", "instance.pose:6:0")   // instance position
	shader.SetBindingForAttribute(cst.Vec3, "icolor", "instance.pose:6:3") // instance color
	shader.CheckBindings()                                                 // check validity of the shader
	return shader
}
