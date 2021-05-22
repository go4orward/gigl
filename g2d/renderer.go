package g2d

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go4orward/gigl"
	"github.com/go4orward/gigl/common"
)

type Renderer struct {
	rc   gigl.GLRenderingContext //
	axes *SceneObject            //
}

func NewRenderer(rc gigl.GLRenderingContext) *Renderer {
	renderer := Renderer{rc: rc, axes: nil}
	return &renderer
}

// ----------------------------------------------------------------------------
// Clear
// ----------------------------------------------------------------------------

func (self *Renderer) Clear(scene *Scene) {
	rc, c := self.rc, self.rc.GetConstants()
	rgb := scene.GetBkgColor()
	rc.GLClearColor(rgb[0], rgb[1], rgb[2], 1.0) // Set clearing color
	rc.GLClear(c.COLOR_BUFFER_BIT)               // clear the canvas
	rc.GLClear(c.DEPTH_BUFFER_BIT)               // clear the canvas
}

func (self *Renderer) RenderAxes(camera *Camera, length float32) {
	// Render three axes (X:RED, Y:GREEN, Z:BLUE) for visual reference
	if self.axes == nil {
		// create geometry with two lines for X,Y axes, with origin at (0,0)
		geometry := NewGeometry()                                            // create an empty geometry
		geometry.SetVertices([][2]float32{{0, 0}, {length, 0}, {0, length}}) // add three vertices
		geometry.SetEdges([][]uint32{{0, 1}, {0, 2}})                        // add two edges
		geometry.BuildDataBuffers(true, true, false)                         // build data buffers for vertices and edges
		shader := NewShader_2DAxes(self.rc)                                  // create shader, and set its bindings
		self.axes = NewSceneObject(geometry, nil, nil, shader, nil)          // set up the scene object (draw LINES)
	}
	rc, c := self.rc, self.rc.GetConstants()
	rc.GLBindBuffer(c.ARRAY_BUFFER, nil)
	rc.GLBindBuffer(c.ELEMENT_ARRAY_BUFFER, nil)
	self.RenderSceneObject(self.axes, &camera.pjvwmatrix) // (Proj * View) matrix
}

// ----------------------------------------------------------------------------
// Rendering Scene
// ----------------------------------------------------------------------------

func (self *Renderer) RenderScene(scene *Scene, camera *Camera) {
	// Render all the scene objects
	for _, scnobj := range scene.objects {
		pvm_matrix := camera.pjvwmatrix.MultiplyToTheRight(&scnobj.modelmatrix)
		self.RenderSceneObject(scnobj, pvm_matrix) // (Proj * View * Model) matrix
	}
	// Render all the OverlayLayers
	for _, overlay := range scene.overlays {
		overlay.Render(&camera.pjvwmatrix)
	}
}

// ----------------------------------------------------------------------------
// Rendering SceneObject
// ----------------------------------------------------------------------------

func (self *Renderer) RenderSceneObject(scnobj *SceneObject, pvm *common.Matrix3) error {
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
	geom := scnobj.Geometry
	if geom.IsDataBufferReady() == false {
		return errors.New("Failed to RenderSceneObject() : empty geometry data buffer")
	}
	if scnobj.vao == nil {
		scnobj.vao = rc.CreateDataBufferVAO()
	}
	if scnobj.vao.VertBuffer == nil && scnobj.vao.FvtxBuffer == nil {
		// create data buffers & buffer information for RenderingContext, and save them in VAO
		scnobj.vao.VertBuffer = rc.CreateVtxDataBuffer(geom.GetVtxBuffer(0))
		scnobj.vao.VertBufferInfo = geom.GetVtxBufferInfo(0)
		if geom.GetIdxBuffer(2) != nil {
			scnobj.vao.EdgeBuffer = rc.CreateIdxDataBuffer(geom.GetIdxBuffer(2))
			scnobj.vao.EdgeBufferCount = geom.GetIdxBufferCount(2)
		}
		if geom.GetIdxBuffer(3) != nil {
			if geom.GetVtxBuffer(1) != nil {
				scnobj.vao.FvtxBuffer = rc.CreateVtxDataBuffer(geom.GetVtxBuffer(1))
				scnobj.vao.FvtxBufferInfo = geom.GetVtxBufferInfo(1)
			}
			scnobj.vao.FaceBuffer = rc.CreateIdxDataBuffer(geom.GetIdxBuffer(3))
			scnobj.vao.FaceBufferCount = geom.GetIdxBufferCount(3)
		}
	}
	if scnobj.vao.InstanceBuffer == nil && scnobj.instance_buffer != nil {
		scnobj.vao.InstanceBuffer = rc.CreateVtxDataBuffer(scnobj.instance_buffer)
		if !self.rc.IsExtensionReady("ANGLE") {
			self.rc.SetupExtension("ANGLE")
		}
	}
	// R3: Render the object with FACE shader
	if scnobj.FShader != nil {
		err := self.render_scene_object_with_shader(scnobj, pvm, 3, scnobj.FShader)
		if err != nil {
			return err
		}
	}
	// R2: Render the object with EDGE shader
	if scnobj.EShader != nil {
		err := self.render_scene_object_with_shader(scnobj, pvm, 2, scnobj.EShader)
		if err != nil {
			return err
		}
	}
	// R1: Render the object with VERTEX shader
	if scnobj.VShader != nil {
		err := self.render_scene_object_with_shader(scnobj, pvm, 1, scnobj.VShader)
		if err != nil {
			return err
		}
	}
	// Render all the children
	for _, child := range scnobj.children {
		new_pvm := pvm.MultiplyToTheRight(&child.modelmatrix)
		self.RenderSceneObject(child, new_pvm)
	}
	return nil
}

func (self *Renderer) render_scene_object_with_shader(scnobj *SceneObject, pvm *common.Matrix3, draw_mode int, shader gigl.GLShader) error {
	rc, c := self.rc, self.rc.GetConstants()
	// 1. Decide which Shader to use
	if shader == nil {
		return errors.New("Failed to render SceneObject() : shader not found")
	} else if shader.GetErr() != nil {
		return errors.New("Failed to render SceneObject() : shader has error")
	}
	rc.GLUseProgram(shader.GetShaderProgram())
	// 2. bind the uniforms of the shader program
	for uname, umap := range shader.GetUniformBindings() {
		if err := self.bind_uniform(uname, umap, draw_mode, scnobj, pvm); err != nil {
			if err.Error() != "Texture is not ready" {
				fmt.Println(err.Error())
			}
			return err
		}
	}
	// 3. bind the attributes of the shader program
	for aname, amap := range shader.GetAttributeBindings() {
		if err := self.bind_attribute(aname, amap, draw_mode, scnobj); err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	// 4. draw  (Note that ARRAY_BUFFER was binded already in the attribut-binding step)
	switch draw_mode {
	case 3: // draw TRIANGLES (FACES)
		buffer, count := scnobj.vao.GetIdxBuffer(3)
		if count > 0 {
			rc.GLBindBuffer(c.ELEMENT_ARRAY_BUFFER, buffer)
			if scnobj.instance_count == 0 {
				rc.GLDrawElements(c.TRIANGLES, count, c.UNSIGNED_INT, 0) // (mode, count, type, offset)
			} else {
				rc.GLDrawElementsInstanced(c.TRIANGLES, count, c.UNSIGNED_INT, 0, scnobj.instance_count)
			}
		}
	case 2: // draw LINES (EDGES)
		buffer, count := scnobj.vao.GetIdxBuffer(2)
		if count > 0 {
			rc.GLBindBuffer(c.ELEMENT_ARRAY_BUFFER, buffer)
			if scnobj.instance_count == 0 {
				rc.GLDrawElements(c.LINES, count, c.UNSIGNED_INT, 0) // (mode, count, type, offset)
			} else {
				rc.GLDrawElementsInstanced(c.LINES, count, c.UNSIGNED_INT, 0, scnobj.instance_count)
			}
		}
	case 1: // draw POINTS (VERTICES)
		_, binfo := scnobj.vao.GetVtxBuffer(0, 0) // nverts, stride, size, offset
		nverts := binfo[0]
		if nverts > 0 {
			if scnobj.instance_count == 0 {
				rc.GLDrawArrays(c.POINTS, 0, nverts) // (mode, first, count)
			} else {
				rc.GLDrawArraysInstanced(c.POINTS, 0, nverts, scnobj.instance_count)
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
	draw_mode int, scnobj *SceneObject, pvm *common.Matrix3) error {
	rc, c := self.rc, self.rc.GetConstants()
	if umap["location"] == nil {
		err := fmt.Errorf("Failed to bind uniform '%s' : call 'shader.CheckBinding()' before rendering", uname)
		return err
	}
	location, dtype := umap["location"], umap["dtype"].(string)
	if umap["autobinding"] != nil {
		autobinding := umap["autobinding"].(string)
		autobinding_split := strings.Split(autobinding, ":")
		autobinding0 := autobinding_split[0]
		switch autobinding0 {
		case "material.color":
			c := [4]float32{0, 1, 1, 1}
			if scnobj.Material != nil {
				c = scnobj.Material.GetDrawModeColor(draw_mode) // get color from material (for the DrawMode)
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
			if scnobj.Material == nil || !scnobj.Material.IsTextureReady() || scnobj.Material.IsTextureLoading() {
				return errors.New("Texture is not ready")
			}
			txt_unit := 0
			if len(autobinding_split) >= 2 {
				txt_unit, _ = strconv.Atoi(autobinding_split[1])
			}
			rc.GLActiveTexture(txt_unit)                                 // activate texture unit N
			rc.GLBindTexture(c.TEXTURE_2D, scnobj.Material.GetTexture()) // bind the texture
			rc.GLUniform1i(location, txt_unit)                           // give shader the unit number
			return nil
		case "renderer.aspect": // vec2
			wh := rc.GetWH()
			rc.GLUniform2f(location, float32(wh[0]), float32(wh[1]))
			return nil
		case "renderer.pvm": // mat3
			elements := pvm.GetElements()                       // ModelView matrix
			rc.GLUniformMatrix3fv(location, false, elements[:]) // gl.uniformMatrix3fv(location, transpose, values_array)
			return nil
		}
		return fmt.Errorf("Failed to bind uniform '%s' (%s) with %v", uname, dtype, autobinding)
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

func (self *Renderer) bind_attribute(aname string, amap map[string]interface{}, draw_mode int, scnobj *SceneObject) error {
	rc, c := self.rc, self.rc.GetConstants()
	if amap["location"] == nil {
		err := fmt.Errorf("Failed to bind attribute '%s' : call 'shader.CheckBinding()' before rendering", aname)
		return err
	}
	location, dtype := amap["location"], amap["dtype"].(string)
	autobinding := amap["autobinding"].(string)
	// fmt.Printf("Attribute (%s) : autobinding= '%s'\n", dtype, autobinding)
	autobinding_split := strings.Split(autobinding, ":")
	autobinding0 := autobinding_split[0]
	switch autobinding0 {
	case "geometry.coords": // 2 * float32 in 8 bytes (2 float32)
		buffer, binfo := scnobj.vao.GetVtxBuffer(0, 0) // [4]int{ nverts, stride, size, offset }
		rc.GLBindBuffer(c.ARRAY_BUFFER, buffer)
		rc.GLVertexAttribPointer(location, binfo[2], c.FLOAT, false, binfo[1]*4, binfo[3]*4)
		rc.GLEnableVertexAttribArray(location)
		if rc.IsExtensionReady("ANGLE") {
			// context.ext_angle.vertexAttribDivisorANGLE(attribute_loc, divisor);
			rc.GLVertexAttribDivisor(location, 0) // divisor == 0
		}
		return nil
	case "geometry.textuv": // 2 * uint16 in 4 bytes (1 float32)
		buffer, binfo := scnobj.vao.GetVtxBuffer(0, 1) // [4]int{ nverts, stride, size, offset }
		rc.GLBindBuffer(c.ARRAY_BUFFER, buffer)
		rc.GLVertexAttribPointer(location, binfo[2], c.UNSIGNED_SHORT, true, binfo[1]*4, binfo[3]*4)
		rc.GLEnableVertexAttribArray(location)
		if rc.IsExtensionReady("ANGLE") {
			// context.ext_angle.vertexAttribDivisorANGLE(attribute_loc, divisor);
			rc.GLVertexAttribDivisor(location, 0) // divisor == 0
		}
		return nil
	case "instance.pose":
		if scnobj.vao.InstanceBuffer != nil && len(autobinding_split) == 3 { // it's like "instance.pose:<stride>:<offset>"
			count := get_count_from_type(dtype)
			stride, _ := strconv.Atoi(autobinding_split[1])
			offset, _ := strconv.Atoi(autobinding_split[2])
			rc.GLBindBuffer(c.ARRAY_BUFFER, scnobj.vao.InstanceBuffer)
			rc.GLVertexAttribPointer(location, count, c.FLOAT, false, stride*4, offset*4)
			rc.GLEnableVertexAttribArray(location)
			// context.ext_angle.vertexAttribDivisorANGLE(attribute_loc, divisor);
			rc.GLVertexAttribDivisor(location, 1) // divisor == 1
			return nil
		}
	case "instance.color":
		if scnobj.vao.InstanceBuffer != nil && len(autobinding_split) == 3 { // it's like "instance.color:<stride>:<offset>"
			count := get_count_from_type(dtype)
			stride, _ := strconv.Atoi(autobinding_split[1])
			offset, _ := strconv.Atoi(autobinding_split[2])
			rc.GLBindBuffer(c.ARRAY_BUFFER, scnobj.vao.InstanceBuffer)
			rc.GLVertexAttribPointer(location, count, c.BYTE, true, stride*4, offset*4)
			rc.GLEnableVertexAttribArray(location)
			// context.ext_angle.vertexAttribDivisorANGLE(attribute_loc, divisor);
			rc.GLVertexAttribDivisor(location, 1) // divisor == 1
			return nil
		}
	default:
		buffer, stride_i, offset_i := amap["buffer"], amap["stride"], amap["offset"]
		if buffer != nil && stride_i != nil && offset_i != nil {
			size, stride, offset := get_count_from_type(dtype), stride_i.(int), offset_i.(int)
			rc.GLBindBuffer(c.ARRAY_BUFFER, buffer)
			rc.GLVertexAttribPointer(location, size, c.FLOAT, false, stride*4, offset*4)
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
