package g3d

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go4orward/gigl"
	"github.com/go4orward/gigl/common"
)

type Renderer struct {
	rc   gigl.GLRenderingContext
	axes *SceneObject
}

func NewRenderer(rc gigl.GLRenderingContext) *Renderer {
	renderer := Renderer{rc: rc, axes: nil}
	return &renderer
}

// ----------------------------------------------------------------------------
// Clear
// ----------------------------------------------------------------------------

func (self *Renderer) Clear(scene *Scene) {
	c := self.rc.GetConstants()
	rgb := scene.GetBkgColor()
	self.rc.GLClearColor(rgb[0], rgb[1], rgb[2], 1.0) // set clearing color
	self.rc.GLClear(c.COLOR_BUFFER_BIT)               // clear the canvas
	self.rc.GLClear(c.DEPTH_BUFFER_BIT)               // clear the canvas
}

// ----------------------------------------------------------------------------
// Rendering Axes
// ----------------------------------------------------------------------------

func (self *Renderer) RenderAxes(camera *Camera, length float32) {
	// Render three axes (X:RED, Y:GREEN, Z:BLUE) for visual reference
	if self.axes == nil {
		self.axes = NewSceneObject_3DAxes(self.rc, length)
	}
	self.RenderSceneObject(self.axes, camera.projection.GetMatrix(), &camera.viewmatrix)
	// camera.TestDataBuffer(self.axes.geometry.data_buffer_vpoints, self.axes.geometry.vpoint_info[0])
}

// ----------------------------------------------------------------------------
// Rendering Scene
// ----------------------------------------------------------------------------

func (self *Renderer) RenderScene(scene *Scene, camera *Camera) {
	// Render all the SceneObjects in the Scene
	for _, sobj := range scene.objects {
		new_viewmodel := camera.viewmatrix.MultiplyToTheRight(&sobj.modelmatrix)
		self.RenderSceneObject(sobj, camera.projection.GetMatrix(), new_viewmodel)
	}
	// Render all the OverlayLayers
	for _, overlay := range scene.overlays {
		overlay.Render(camera.projection.GetMatrix(), &camera.viewmatrix)
	}
}

// ----------------------------------------------------------------------------
// Rendering SceneObject
// ----------------------------------------------------------------------------

func (self *Renderer) RenderSceneObject(scnobj *SceneObject, proj *common.Matrix4, vwmd *common.Matrix4) error {
	rc, c := self.rc, self.rc.GetConstants()
	// Set DepthTest & Blending options
	if scnobj.UseDepth {
		rc.GLEnable(c.DEPTH_TEST) // Enable depth test
		rc.GLDepthFunc(c.LEQUAL)  // Near things obscure far things
	} else {
		rc.GLDisable(c.DEPTH_TEST) // Disable depth test
	}
	if scnobj.UseBlend {
		rc.GLEnable(c.BLEND)                         // for pre-multiplied alpha
		rc.GLBlendFunc(c.ONE, c.ONE_MINUS_SRC_ALPHA) // for pre-multiplied alpha
		// rc.BlendFunc(c.SRC_ALPHA, c.ONE_MINUS_SRC_ALPHA) // for non pre-multiplied alpha
	} else {
		rc.GLDisable(c.BLEND) // Disable blending
	}
	// If necessary, then build WebGLBuffers for the SceneObject's Geometry
	if scnobj.Geometry.IsDataBufferReady() == false {
		return errors.New("Failed to RenderSceneObject() : empty geometry data buffer")
	}
	if scnobj.Geometry.IsWebGLBufferReady() == false {
		scnobj.Geometry.BuildWebGLBuffers(self.rc, true, true, true)
	}
	if scnobj.poses != nil && scnobj.poses.IsRcBufferReady() == false {
		scnobj.poses.BuildRcBuffer(self.rc)
		if !self.rc.IsExtensionReady("ANGLE") {
			self.rc.SetupExtension("ANGLE")
		}
	}
	// R3: Render the object with FACE shader
	if scnobj.FShader != nil {
		err := self.render_scene_object_with_shader(scnobj, proj, vwmd, 3, scnobj.FShader)
		if err != nil {
			return err
		}
	}
	// R2: Render the object with EDGE shader
	if scnobj.EShader != nil {
		err := self.render_scene_object_with_shader(scnobj, proj, vwmd, 2, scnobj.EShader)
		if err != nil {
			return err
		}
	}
	// R1: Render the object with VERTEX shader
	if scnobj.VShader != nil {
		err := self.render_scene_object_with_shader(scnobj, proj, vwmd, 1, scnobj.VShader)
		if err != nil {
			return err
		}
	}
	// Render all the children
	for _, child := range scnobj.children {
		new_viewmodel := vwmd.MultiplyToTheRight(&child.modelmatrix)
		self.RenderSceneObject(child, proj, new_viewmodel)
	}
	return nil
}

func (self *Renderer) render_scene_object_with_shader(scnobj *SceneObject, proj *common.Matrix4, vwmd *common.Matrix4, draw_mode int, shader gigl.GLShader) error {
	rc, c := self.rc, self.rc.GetConstants()
	// 1. Decide which Shader to use
	if shader == nil {
		return errors.New("Failed to RenderSceneObject() : shader not found")
	}
	rc.GLUseProgram(shader.GetShaderProgram())
	// 2. bind the uniforms of the shader program
	for uname, umap := range shader.GetUniformBindings() {
		if err := self.bind_uniform(uname, umap, draw_mode, scnobj.Material, proj, vwmd); err != nil {
			if err.Error() != "Texture is not ready" {
				fmt.Println(err.Error())
			}
			return err
		}
	}
	// 3. bind the attributes of the shader program
	for aname, amap := range shader.GetAttributeBindings() {
		if err := self.bind_attribute(aname, amap, draw_mode, scnobj.Geometry, scnobj.poses); err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	// 4. draw  (Note that ARRAY_BUFFER was binded already in the attribut-binding step)
	switch draw_mode {
	case 3: // draw TRIANGLES (FACES)
		buffer, count, _ := scnobj.Geometry.GetWebGLBuffer(draw_mode)
		if count > 0 {
			rc.GLBindBuffer(c.ELEMENT_ARRAY_BUFFER, buffer)
			if scnobj.poses == nil {
				// fmt.Printf("draw FACES with drawElements()\n")
				rc.GLDrawElements(c.TRIANGLES, count, c.UNSIGNED_INT, 0) // (mode, count, type, offset)
			} else {
				// fmt.Printf("draw FACES with drawElementsInstancedANGLE()\n")
				rc.GLDrawElementsInstanced(c.TRIANGLES, count, c.UNSIGNED_INT, 0, scnobj.poses.Count)
			}
		}
	case 2: // draw LINES (EDGES)
		buffer, count, _ := scnobj.Geometry.GetWebGLBuffer(draw_mode)
		if count > 0 {
			rc.GLBindBuffer(c.ELEMENT_ARRAY_BUFFER, buffer)
			if scnobj.poses == nil {
				rc.GLDrawElements(c.LINES, count, c.UNSIGNED_INT, 0) // (mode, count, type, offset)
			} else {
				rc.GLDrawElementsInstanced(c.LINES, count, c.UNSIGNED_INT, 0, scnobj.poses.Count)
			}
		}
	case 1: // draw POINTS (VERTICES)
		_, count, pinfo := scnobj.Geometry.GetWebGLBuffer(draw_mode)
		if count > 0 {
			vert_count := count / pinfo[0] // number of vertices
			if scnobj.poses == nil {
				rc.GLDrawArrays(c.POINTS, 0, vert_count) // (mode, first, count)
			} else {
				rc.GLDrawArraysInstanced(c.POINTS, 0, vert_count, scnobj.poses.Count)
			}
		}
	default:
		err := fmt.Errorf("Unknown mode to draw : %d\n", draw_mode)
		fmt.Printf(err.Error())
		return err
	}
	return nil
}

func (self *Renderer) bind_uniform(uname string, umap map[string]interface{},
	draw_mode int, material gigl.GLMaterial, proj *common.Matrix4, vwmd *common.Matrix4) error {
	rc, c := self.rc, self.rc.GetConstants()
	if umap["location"] == nil {
		err := errors.New("Failed to bind uniform : call 'shader.CheckBinding()' before rendering")
		return err
	}
	location, dtype := umap["location"], umap["dtype"].(string)
	if umap["autobinding"] != nil {
		autobinding := umap["autobinding"].(string)
		// fmt.Printf("Uniform (%s) : autobinding= '%s'\n", dtype, autobinding)
		autobinding_split := strings.Split(autobinding, ":")
		autobinding0 := autobinding_split[0]
		switch autobinding0 {
		case "renderer.aspect": // vec2
			wh := rc.GetWH()
			rc.GLUniform2f(location, float32(wh[0]), float32(wh[1]))
			return nil
		case "renderer.proj": // mat4
			e := (*proj.GetElements())[:]
			rc.GLUniformMatrix4fv(location, false, e) // gl.uniformMatrix4fv(location, transpose, values_array)
			return nil
		case "renderer.vwmd": // mat4
			e := (*vwmd.GetElements())[:]
			rc.GLUniformMatrix4fv(location, false, e) // gl.uniformMatrix4fv(location, transpose, values_array)
			return nil
		case "renderer.pvm": // mat4
			pvm := proj.MultiplyToTheRight(vwmd)      // (Proj * View * Models) matrix
			e := (*pvm.GetElements())[:]              //
			rc.GLUniformMatrix4fv(location, false, e) // gl.uniformMatrix4fv(location, transpose, values_array)
			return nil
		case "material.color":
			c := [4]float32{0, 1, 1, 1}
			if material != nil {
				c = material.GetDrawModeColor(draw_mode) // get color from material (for the DrawMode)
			}
			switch dtype {
			case "vec3":
				rc.GLUniform3f(location, c[0], c[1], c[2])
				return nil
			case "vec4":
				rc.GLUniform4f(location, c[0], c[1], c[2], c[3])
				return nil
			}
		case "material.texture":
			if material == nil || !material.IsTextureReady() || material.IsTextureLoading() {
				return errors.New("Texture is not ready")
			}
			txt_unit := 0
			if len(autobinding_split) >= 2 {
				txt_unit, _ = strconv.Atoi(autobinding_split[1])
			}
			rc.GLActiveTexture(txt_unit)                          // activate texture unit N
			rc.GLBindTexture(c.TEXTURE_2D, material.GetTexture()) // bind the texture
			rc.GLUniform1i(location, txt_unit)                    // give shader the unit number
			return nil
		case "lighting.dlight": // mat3
			dlight := common.NewMatrix3().Set(0, 1, 0, 0, 1, 0, 1, 1, 0) // directional light (in camera space)
			e := (*dlight.GetElements())[:]                              // (direction[3] + intensity[3] + ambient[3])
			rc.GLUniformMatrix3fv(location, false, e)                    // gl.uniformMatrix4fv(location, transpose, values_array)
			return nil
		}
		return fmt.Errorf("Failed to bind uniform '%s' (%s) with %v", uname, dtype, umap)
	} else if umap["value"] != nil {
		v := umap["value"].([]float32)
		switch dtype {
		case "int":
			rc.GLUniform1i(location, int(v[0]))
			return nil
		case "float":
			rc.GLUniform1f(location, v[0])
			return nil
		case "vec2":
			rc.GLUniform2f(location, v[0], v[1])
			return nil
		case "vec3":
			rc.GLUniform3f(location, v[0], v[1], v[2])
			return nil
		case "vec4":
			rc.GLUniform4f(location, v[0], v[1], v[2], v[3])
			return nil
		}
		return fmt.Errorf("Failed to bind uniform '%s' (%s) with %v", uname, dtype, v)
	} else {
		return fmt.Errorf("Failed to bind uniform '%s' (%s)", uname, dtype)
	}
}

func (self *Renderer) bind_attribute(aname string, amap map[string]interface{},
	draw_mode int, geometry gigl.GLGeometry, poses *SceneObjectPoses) error {
	rc, c := self.rc, self.rc.GetConstants()
	if amap["location"] == nil {
		err := errors.New("Failed to bind attribute : call 'shader.CheckBinding()' before rendering")
		return err
	}
	location, dtype := amap["location"], amap["dtype"].(string)
	autobinding := amap["autobinding"].(string)
	// fmt.Printf("Attribute (%s) : autobinding= '%s'\n", dtype, autobinding)
	autobinding_split := strings.Split(autobinding, ":")
	autobinding0 := autobinding_split[0]
	switch autobinding0 {
	case "geometry.coords": // 3 * float32 in 12 bytes (3 float32)
		buffer, _, pinfo := geometry.GetWebGLBuffer(1)
		rc.GLBindBuffer(c.ARRAY_BUFFER, buffer)
		rc.GLVertexAttribPointer(location, 3, c.FLOAT, false, pinfo[0]*4, pinfo[1]*4)
		rc.GLEnableVertexAttribArray(location)
		if rc.IsExtensionReady("ANGLE") {
			// context.ext_angle.vertexAttribDivisorANGLE(attribute_loc, divisor);
			rc.GLVertexAttribDivisor(location, 0) // divisor == 0
		}
		return nil
	case "geometry.textuv": // 2 * uint16 in 4 bytes (1 float32)
		buffer, _, pinfo := geometry.GetWebGLBuffer(1)
		rc.GLBindBuffer(c.ARRAY_BUFFER, buffer)
		rc.GLVertexAttribPointer(location, 2, c.UNSIGNED_SHORT, true, pinfo[0]*4, pinfo[2]*4)
		rc.GLEnableVertexAttribArray(location)
		if pinfo[1] == pinfo[2] {
			fmt.Printf("Renderer Warning : Texture UV coordinates not found (pinfo=%v)\n", pinfo)
		}
		if rc.IsExtensionReady("ANGLE") {
			// context.ext_angle.vertexAttribDivisorANGLE(attribute_loc, divisor);
			rc.GLVertexAttribDivisor(location, 0) // divisor == 0
		}
		return nil
	case "geometry.normal": // 3 * byte in 4 bytes (1 float32)
		buffer, _, pinfo := geometry.GetWebGLBuffer(1)
		count := get_count_from_type(dtype)
		rc.GLBindBuffer(c.ARRAY_BUFFER, buffer)
		rc.GLVertexAttribPointer(location, count, c.BYTE, true, pinfo[0]*4, pinfo[3]*4)
		rc.GLEnableVertexAttribArray(location)
		if pinfo[1] == pinfo[3] {
			fmt.Printf("Renderer Warning : Normal vectors not found (pinfo=%v)\n", pinfo)
		}
		if rc.IsExtensionReady("ANGLE") {
			// context.ext_angle.vertexAttribDivisorANGLE(attribute_loc, divisor);
			rc.GLVertexAttribDivisor(location, 0) // divisor == 0
		}
		return nil
	case "instance.pose":
		if poses != nil && len(autobinding_split) == 3 { // it's like "instance.pose:<stride>:<offset>"
			count := get_count_from_type(dtype)
			stride, _ := strconv.Atoi(autobinding_split[1])
			offset, _ := strconv.Atoi(autobinding_split[2])
			rcbuffer, _, _ := poses.GetRcBuffer()
			rc.GLBindBuffer(c.ARRAY_BUFFER, rcbuffer)
			rc.GLVertexAttribPointer(location, count, c.FLOAT, false, stride*4, offset*4)
			rc.GLEnableVertexAttribArray(location)
			// context.ext_angle.vertexAttribDivisorANGLE(attribute_loc, divisor);
			rc.GLVertexAttribDivisor(location, 1) // divisor == 1
			return nil
		}
	default:
		buffer, stride_i, offset_i := amap["buffer"], amap["stride"], amap["offset"]
		if buffer != nil && stride_i != nil && offset_i != nil {
			count, stride, offset := get_count_from_type(dtype), stride_i.(int), offset_i.(int)
			rc.GLBindBuffer(c.ARRAY_BUFFER, buffer)
			rc.GLVertexAttribPointer(location, count, c.FLOAT, false, stride*4, offset*4)
			rc.GLEnableVertexAttribArray(location)
			if rc.IsExtensionReady("ANGLE") {
				// context.ext_angle.vertexAttribDivisorANGLE(attribute_loc, divisor);
				rc.GLVertexAttribDivisor(location, 0) // divisor == 0
			}
		}
	}
	return fmt.Errorf("Failed to bind attribute '%s' (%s) with %v", aname, dtype, amap)
}

func get_count_from_type(dtype string) int {
	switch dtype {
	case "float":
		return 1
	case "vec2":
		return 2
	case "vec3":
		return 3
	case "vec4":
		return 4
	default:
		return 0
	}
}
