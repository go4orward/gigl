package g3d

import (
	"fmt"
	"strings"

	"github.com/go4orward/gigl"
	"github.com/go4orward/gigl/common"
	cst "github.com/go4orward/gigl/common/constants"
	"github.com/go4orward/gigl/g2d"
)

// ----------------------------------------------------------------------------
// OverlayLabel
// ----------------------------------------------------------------------------

type OverlayLabel struct {
	text    string       // text of the label
	xyz     [3]float32   // origin of the label (in WORLD space)
	chwh    [2]float32   // character width & height
	color   string       // color of the label (like "#ff0000")
	offset  [2]float32   // offset from origin (in pixels in CAMERA space)
	offref  string       // offset reference type, like "L_TOP", "R_BTM", "CENTER", etc
	angle   float32      // rotation angle
	bkgtype string       // background type, like "box:#ffff00:#000000", "under:#000000", etc
	txtobj  *SceneObject // SceneObject for rendering the label text
	bkgobj  *SceneObject // SceneObject for rendering the background
}

func (self *OverlayLabel) String(chwh [2]float32) string {
	return fmt.Sprintf("OvLabel{xy:%v text:'%s'}", self.xyz, self.text)
}

func (self *OverlayLabel) SetCharacterWH(chwh [2]float32) *OverlayLabel {
	self.chwh = chwh // character width & height
	return self
}

func (self *OverlayLabel) SetPose(rotation float32, offset_reference string, offset [2]float32) *OverlayLabel {
	text_width := self.chwh[0] * float32(len([]rune(self.text)))
	text_height := self.chwh[1]
	self.offref = offset_reference
	switch self.offref {
	case "L_TOP":
		self.offset = [2]float32{offset[0], offset[1] - text_height/2}
	case "L_MID":
		self.offset = [2]float32{offset[0], offset[1]}
	case "L_BTM":
		self.offset = [2]float32{offset[0], offset[1] + text_height/2}
	case "M_TOP":
		self.offset = [2]float32{offset[0] - text_width/2, offset[1] - text_height/2}
	case "CENTER":
		self.offset = [2]float32{offset[0] - text_width/2, offset[1]}
	case "M_BTM":
		self.offset = [2]float32{offset[0] - text_width/2, offset[1] + text_height/2}
	case "R_TOP":
		self.offset = [2]float32{offset[0] - text_width, offset[1] - text_height/2}
	case "R_MID":
		self.offset = [2]float32{offset[0] - text_width, offset[1]}
	case "R_BTM":
		self.offset = [2]float32{offset[0] - text_width, offset[1] + text_height/2}
	default: // same as "L_MID"
		self.offset = [2]float32{offset[0], offset[1]}
	}
	return self
}

func (self *OverlayLabel) SetBackground(bkgtype string) *OverlayLabel {
	self.bkgtype = bkgtype
	return self
}

// ----------------------------------------------------------------------------
// OverlayLabelLayer
// ----------------------------------------------------------------------------

type OverlayLabelLayer struct {
	rc     gigl.GLRenderingContext //
	Labels []*OverlayLabel         //

	// variables shared by all the labels
	alphabet_geometry *g2d.Geometry                // geometry with single vertex at origin
	alphabet_texture  *g2d.MaterialAlphabetTexture // alphabet texture
	alphabet_shader   gigl.GLShader                // label text shader
}

func NewOverlayLabelLayer(rc gigl.GLRenderingContext, fontsize int, outlined bool) *OverlayLabelLayer {
	self := OverlayLabelLayer{rc: rc} // let 'fontsize' of ALPHABET texture to be 20, by default
	self.Labels = make([]*OverlayLabel, 0)
	self.alphabet_geometry = g2d.NewGeometryOrigin() // trivial geometry with single vertex at origin
	self.alphabet_texture = g2d.NewMaterialAlphabetTexture("Courier New", fontsize, "#ffffff", outlined)
	self.alphabet_shader = nil
	return &self
}

func (self *OverlayLabelLayer) Render(proj *common.Matrix4, view *common.Matrix4) {
	// 'Overlay' interface function, called by Renderer
	renderer := NewRenderer(self.rc)
	for _, label := range self.Labels {
		if label.bkgobj != nil {
			renderer.RenderSceneObject(label.bkgobj, proj, view)
		}
		if label.txtobj != nil {
			renderer.RenderSceneObject(label.txtobj, proj, view)
		}
	}
}

// ----------------------------------------------------------------------------
// Managing Labels
// ----------------------------------------------------------------------------

func (self *OverlayLabelLayer) AddLabel(labels ...*OverlayLabel) *OverlayLabelLayer {
	for i := 0; i < len(labels); i++ {
		label := labels[i]
		if label.txtobj == nil && label.text != "" {
			label.txtobj = self.build_labeltext_scene_object(label)
		}
		if label.bkgobj == nil && label.bkgtype != "" {
			label.bkgobj = self.build_labelbkg_scene_object(label)
		}
		self.Labels = append(self.Labels, label)
	}
	return self
}

func (self *OverlayLabelLayer) FindLabel(label_text string) *OverlayLabel {
	for _, label := range self.Labels {
		if label.text == label_text {
			return label
		}
	}
	return nil
}

func (self *OverlayLabelLayer) CreateLabel(label_text string, xyz [3]float32, color string) *OverlayLabel {
	chwh := self.alphabet_texture.GetAlaphabetCharacterWH(1.0)
	label := &OverlayLabel{text: label_text, xyz: xyz, chwh: chwh, color: color}
	return label
}

func (self *OverlayLabelLayer) AddTextLabel(label_text string, xyz [3]float32, color string, offref string) *OverlayLabel {
	// Convenience function to quickly add a Label,
	//   which simplifies all the following steps:
	//   label := layer.CreateLabel();  label.SetPose();  layer.AddLabel(label)
	chwh := self.alphabet_texture.GetAlaphabetCharacterWH(1.0)
	label := &OverlayLabel{text: label_text, xyz: xyz, chwh: chwh, color: color}
	label.SetPose(0, offref, [2]float32{0, 0})
	self.AddLabel(label)
	return label
}

func (self *OverlayLabelLayer) AddLabelsForTest() *OverlayLabelLayer {
	self.AddTextLabel("AhjgyZ", [3]float32{0, .5, .5}, "#ff0000", "")
	label2 := self.CreateLabel("Hello!", [3]float32{0, 1, 1}, "#0000ff")
	label2.SetPose(0, "L_BTM", [2]float32{30, 30}).SetBackground("under:#000000")
	return self.AddLabel(label2)
}

// ----------------------------------------------------------------------------
// Building Label's SceneObjects
// ----------------------------------------------------------------------------

func (self *OverlayLabelLayer) build_labeltext_scene_object(label *OverlayLabel) *SceneObject {
	if self.alphabet_shader == nil {
		var vertex_shader_code = `
		precision mediump float;
		uniform   mat4  proj;		// Projection matrix
		uniform   mat4  vwmd;		// View*Model matrix
		uniform   vec2  asp;		// aspect ratio, w : h
		uniform   vec3  orgn;		// origin of the label (WORLD XYZ coordinates)
		uniform   vec3  offr;		// offset (CAMERA XY in pixels) & rotation angle
		uniform   vec3  whlen;		// character width & height, and label length
		attribute vec2  gvxy;		// geometry's vertex XY position (CAMERA XY in pixel)
		attribute vec2  cpose;		// character index & code
		varying float v_code; 		// character code (index of the character in the alphabet texture)
		void main() {
			vec4 origin = proj * vwmd * vec4(orgn, 1.0);
			if (origin.w != 0.0) { origin = origin / origin.w; }
			vec2 ch_off = vec2(offr.x + whlen[0]/2.0, offr.y) + vec2(cpose[0] * whlen[0], 0.0);
			vec2 offset = vec2((ch_off.x + gvxy.x) * 2.0 / asp[0], (ch_off.y + gvxy.y) * 2.0 / asp[1]);
			gl_Position = vec4(origin.x + offset.x, origin.y + offset.y, origin.z, 1.0);
			gl_PointSize = whlen[1];	// character height
			v_code  = cpose[1];
		}`
		var fragment_shader_code = `
		precision mediump float;
		uniform sampler2D text;		// alphabet texture (ASCII characters from SPACE to DEL)
		uniform   vec3 	whlen;		// character width & height, and label length
		uniform   vec4  color;		// color of the label
		varying   float v_code;     // character code (index of the character in the alphabet texture)
		void main() {
			vec2 uv = gl_PointCoord;
			if (uv[0] < 0.0 || uv[0] > 1.0) discard;
			if (uv[1] < 0.0 || uv[1] > 1.0) discard;
			float u = uv[0] * (whlen[1]/whlen[0]) - 0.5, v = uv[1];
			if (u < 0.0 || u > 1.0 || v < 0.0 || v > 1.0) discard;
			uv = vec2((u + v_code)/whlen[2], v);	// position in the texture (relative to alphabet_length)
			gl_FragColor = texture2D(text, uv) * color;
		}`
		self.alphabet_shader, _ = self.rc.CreateShader(vertex_shader_code, fragment_shader_code)
	}
	shader := self.alphabet_shader.Copy()
	offr := []float32{float32(label.offset[0]), float32(label.offset[1]), 0}
	whlen := []float32{label.chwh[0], label.chwh[1], float32(self.alphabet_texture.GetAlaphabetLength())}
	lrgba := common.RGBAFromHexString(label.color)                                            // label color RGBA
	shader.SetBindingForUniform(cst.Mat4, "proj", "renderer.proj")                            // Projection matrix
	shader.SetBindingForUniform(cst.Mat4, "vwmd", "renderer.vwmd")                            // View*Model matrix
	shader.SetBindingForUniform(cst.Vec2, "asp", "renderer.aspect")                           // AspectRatio
	shader.SetBindingForUniform(cst.Vec3, "orgn", label.xyz[:])                               // label origin
	shader.SetBindingForUniform(cst.Vec3, "offr", offr)                                       // label offset & rotation_angle
	shader.SetBindingForUniform(cst.Vec3, "whlen", whlen)                                     // ch_width, ch_height, alphabet_length
	shader.SetBindingForUniform(cst.Vec4, "color", lrgba[:])                                  // label color
	shader.SetBindingForUniform(cst.Sampler2D, "text", "material.texture")                    // texture sampler (unit:0)
	shader.SetBindingForAttribute(cst.Vec2, "gvxy", "geometry.coords")                        // point coordinates
	shader.SetBindingForAttribute(cst.Vec2, "cpose", "instance.pose:2:0")                     // character pose (:<stride>:<offset>)
	shader.CheckBindings()                                                                    // check validity of the shader
	scnobj := NewSceneObject(self.alphabet_geometry, self.alphabet_texture, shader, nil, nil) // shader for drawing POINTS (for each character)
	// calculate the pose of each character (rune)
	label_runes := []rune(label.text)
	scnobj.SetInstanceBuffer(len(label_runes), 2, nil)
	for i := 0; i < len(label_runes); i++ { // save character index & code for each rune
		scnobj.SetInstancePoseValues(i, 0, float32(i), float32(self.alphabet_texture.GetAlaphabetCharacterIndex(label_runes[i])))
	}
	scnobj.UseDepth = true
	scnobj.UseBlend = true
	return scnobj
}

func (self *OverlayLabelLayer) build_labelbkg_scene_object(label *OverlayLabel) *SceneObject {
	tlen := label.chwh[0] * float32(len([]rune(label.text)))
	ltop := [2]float32{label.offset[0] - 4, label.offset[1] + label.chwh[1]/2}
	lbtm := [2]float32{label.offset[0] - 4, label.offset[1] - label.chwh[1]/2}
	rtop := [2]float32{label.offset[0] + tlen + 4, label.offset[1] + label.chwh[1]/2}
	rbtm := [2]float32{label.offset[0] + tlen + 4, label.offset[1] - label.chwh[1]/2}
	bkgtype_split := strings.Split(label.bkgtype, ":")
	if len(bkgtype_split) < 2 {
		common.Logger.Error("Failed to build_labelbkg_scene_object() : invalid background type '%s'\n", label.bkgtype)
		return nil
	}
	bkgtype0 := bkgtype_split[0]
	geometry := g2d.NewGeometry()
	material := g2d.NewMaterialColors(bkgtype_split[1]) // Material from color string
	switch bkgtype0 {
	case "box": // "box:#ffff00:#000000", "box:<FillColor>:<BorderColor>"
		geometry.SetVertices([][2]float32{lbtm, rbtm, rtop, ltop})
		geometry.SetEdges([][]uint32{{0, 1, 2, 3, 0}})
		geometry.SetFaces([][]uint32{{0, 1, 2, 3}})
		geometry.BuildDataBuffers(true, true, true)
		if len(bkgtype_split) >= 3 { // EDGE color (border color)
			material.SetColorForDrawMode(2, bkgtype_split[2])
		}
	case "under": // "under:#000000", "under:<UnderlineColor>"
		geometry.SetVertices([][2]float32{{0, 0}, lbtm, rbtm})
		switch label.offref {
		case "L_TOP", "L_MID", "L_BTM":
			geometry.SetEdges([][]uint32{{0, 1, 2}})
		case "R_TOP", "R_MID", "R_BTM":
			geometry.SetEdges([][]uint32{{1, 2, 0}})
		case "M_TOP", "CENTER", "M_BTM":
			geometry.SetEdges([][]uint32{{1, 2}})
		default:
			geometry.SetEdges([][]uint32{{1, 2}})
		}
		geometry.BuildDataBuffers(true, true, false)
	default:
		common.Logger.Error("Failed to build_labelbkg_scene_object() : invalid background type '%s'\n", label.bkgtype)
		return nil
	}
	var vertex_shader_code = `
		precision mediump float;
		uniform   mat4  proj;		// Projection matrix
		uniform   mat4  vwmd;		// View*Model matrix
		uniform   vec3  orgn;		// origin of the label (WORLD XYZ coordinates)
		uniform   vec2  asp;		// aspect ratio, w : h
		attribute vec2  gvxy;		// geometry's vertex XY position (CAMERA XY in pixel)
		void main() {
			vec4 origin = proj * vwmd * vec4(orgn, 1.0);
			if (origin.w != 0.0) { origin = origin / origin.w; }
			vec2 offset = vec2(gvxy.x * 2.0 / asp[0], gvxy.y * 2.0 / asp[1]);
			gl_Position = vec4(origin.xy + offset.xy, origin.z, 1.0);
		}`
	var fragment_shader_code = `
		precision mediump float;
		uniform vec4 color;			// color RGBA
		void main() {
			gl_FragColor = color;
		}`
	shader, _ := self.rc.CreateShader(vertex_shader_code, fragment_shader_code)
	shader.SetBindingForUniform(cst.Mat4, "proj", "renderer.proj")     // Projection matrix
	shader.SetBindingForUniform(cst.Mat4, "vwmd", "renderer.vwmd")     // View*Model matrix
	shader.SetBindingForUniform(cst.Vec2, "asp", "renderer.aspect")    // AspectRatio
	shader.SetBindingForUniform(cst.Vec3, "orgn", label.xyz[:])        // label origin
	shader.SetBindingForUniform(cst.Vec4, "color", "material.color")   // label color
	shader.SetBindingForAttribute(cst.Vec2, "gvxy", "geometry.coords") // point coordinates from 2D geometry
	shader.CheckBindings()                                             // check validity of the shader
	scnobj := NewSceneObject(geometry, material, nil, shader, shader)  // shader for drawing EDGEs & FACEs
	scnobj.UseDepth = true
	return scnobj
}
