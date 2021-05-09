package g3d

import (
	"github.com/go4orward/gigl"
	"github.com/go4orward/gigl/common"
	"github.com/go4orward/gigl/g2d"
)

type OverlayMarkerLayer struct {
	rc      gigl.GLRenderingContext //
	Markers []*SceneObject          // list of OverlayMarkers to be rendered (in pixels in CAMERA space)
}

func NewOverlayMarkerLayer(rc gigl.GLRenderingContext) *OverlayMarkerLayer {
	self := OverlayMarkerLayer{rc: rc}
	self.Markers = make([]*SceneObject, 0)
	return &self
}

func (self *OverlayMarkerLayer) Render(proj *common.Matrix4, vwmd *common.Matrix4) {
	// 'Overlay' interface function, called by Renderer
	renderer := NewRenderer(self.rc)
	for _, marker := range self.Markers {
		if marker.poses != nil {
			renderer.RenderSceneObject(marker, proj, vwmd)
		} else {
			vwmd := vwmd.MultiplyToTheRight(&marker.modelmatrix)
			renderer.RenderSceneObject(marker, proj, vwmd)
		}
	}
}

// ----------------------------------------------------------------------------
// Managing Markers
// ----------------------------------------------------------------------------

func (self *OverlayMarkerLayer) AddMarker(marker ...*SceneObject) *OverlayMarkerLayer {
	for i := 0; i < len(marker); i++ {
		self.Markers = append(self.Markers, marker[i])
	}
	return self
}

func (self *OverlayMarkerLayer) AddArrowMarker(size float32, color string, outline_color string, rotation float32, xyz [3]float32) *OverlayMarkerLayer {
	// Convenience function to quickly add a Arrow marker,
	//   which is equivalent to : arrow := layer.CreateArrowMarker();  layer.AddMarker(label)
	arrow := self.CreateArrowMarker(size, color, outline_color, false)
	arrow.Rotate([3]float32{0, 0, 1}, rotation).Translate(xyz[0], xyz[1], xyz[2])
	self.Markers = append(self.Markers, arrow)
	return self
}

func (self *OverlayMarkerLayer) AddArrowHeadMarker(size float32, color string, outline_color string, rotation float32, xyz [3]float32) *OverlayMarkerLayer {
	// Convenience function to quickly add a Arrow marker,
	//   which is equivalent to : ahead := layer.CreateArrowHeadMarker();  ahead.Translate();  layer.AddMarker(ahead)
	ahead := self.CreateArrowHeadMarker(size, color, outline_color, false)
	ahead.Rotate([3]float32{0, 0, 1}, rotation).Translate(xyz[0], xyz[1], xyz[2])
	self.Markers = append(self.Markers, ahead)
	return self
}

func (self *OverlayMarkerLayer) AddSpriteMarker(imgpath string, color string, wh [2]float32, xyz [3]float32, offref string) *OverlayMarkerLayer {
	// Convenience function to quickly add a SPRITE marker,
	//   which is equivalent to : sprite := layer.CreateSpriteMarker();  sprite.Translate();  layer.AddMarker(sprite)
	sprite := self.CreateSpriteMarker(imgpath, color, wh, offref, false).Translate(xyz[0], xyz[1], xyz[2])
	return self.AddMarker(sprite)
}

func (self *OverlayMarkerLayer) AddMarkersForTest() *OverlayMarkerLayer {
	ahead1 := self.CreateArrowHeadMarker(20, "#ffaaaa", "#ff0000", false)
	ahead2 := self.CreateArrowHeadMarker(20, "#aaffaa", "#00ff00", false).Translate(0.5, 0.5, 0.5)
	ahead3 := self.CreateArrowHeadMarker(20, "#ffaaaa", "#ff0000", true)
	ahead3.SetupPoses(3, 3, []float32{0, 0, -1, 1, 1, 1, 0, 0, 1})
	sprite := self.CreateSpriteMarker("/assets/map_marker.png", "#ff0000", [2]float32{20, 20}, "M_BTM", false)
	sprite.SetupPoses(3, 2, []float32{0, 0, 1, 1, 1, 1})
	return self.AddMarker(ahead1, ahead2, ahead3, sprite)
}

// ----------------------------------------------------------------------------
// Arrow Marker
// ----------------------------------------------------------------------------

func (self *OverlayMarkerLayer) CreateArrowMarker(size float32, color string, outline_color string, use_poses bool) *SceneObject {
	geometry := g2d.NewGeometry() // 2D geometry of ARROW pointing left, with tip at (0,0)
	geometry.SetVertices([][2]float32{{0, 0}, {0.5, -0.3}, {0.5, -0.15}, {1, -0.15}, {1, 0.15}, {0.5, 0.15}, {0.5, 0.3}})
	geometry.SetFaces([][]uint32{{0, 1, 2, 3, 4, 5, 6}})
	geometry.SetEdges([][]uint32{{0, 1, 2, 3, 4, 5, 6, 0}})
	geometry.Scale(size, size).BuildDataBuffers(true, true, true) // marker size is 10 pixels
	material, _ := self.rc.CreateMaterial(color)
	material.SetColorForDrawMode(2, outline_color)
	shader := self.GetShaderForMarker(use_poses)
	marker := NewSceneObject(geometry, material, nil, shader, shader)
	return marker
}

func (self *OverlayMarkerLayer) CreateArrowHeadMarker(size float32, color string, outline_color string, use_poses bool) *SceneObject {
	geometry := g2d.NewGeometry() // 2D geometry of ARROW pointing left, with tip at (0,0)
	geometry.SetVertices([][2]float32{{0, 0}, {1, -0.6}, {1, +0.6}})
	geometry.SetFaces([][]uint32{{0, 1, 2}})
	geometry.SetEdges([][]uint32{{0, 1, 2, 0}})
	geometry.Scale(size, size).BuildDataBuffers(true, true, true) // marker size is 10 pixels
	material, _ := self.rc.CreateMaterial(color)
	material.SetColorForDrawMode(2, outline_color)
	shader := self.GetShaderForMarker(use_poses)
	marker := NewSceneObject(geometry, material, nil, shader, shader)
	return marker
}

func (self *OverlayMarkerLayer) GetShaderForMarker(use_poses bool) gigl.GLShader {
	var shader gigl.GLShader = nil
	if !use_poses { // Shader for single instance (located at (0,0))
		var vertex_shader_code = `
			precision mediump float;
			uniform   mat4 proj;		// 3D Projection matrix
			uniform   mat4 vwmd;		// 3D View * Model matrix
			uniform   vec2 asp;			// aspect ratio, w : h
			attribute vec2 pxy;			// 2D vertex XY coordinates (offset; pixels in CAMERA space)
			void main() {
				vec4 origin = proj * vwmd * vec4(0.0, 0.0, 0.0, 1.0);
				if (origin.w != 0.0) { origin = origin / origin.w; }
				vec2 offset = vec2(pxy.x * 2.0 / asp[0], pxy.y * 2.0 / asp[1]);
				gl_Position = vec4(origin.x + offset.x, origin.y + offset.y, origin.z, 1.0);
			}`
		var fragment_shader_code = `
			precision mediump float;
			uniform vec3 color;			// color
			void main() { 
				gl_FragColor = vec4(color.r, color.g, color.b, 1.0);
			}`
		shader, _ = self.rc.CreateShader(vertex_shader_code, fragment_shader_code)
		shader.SetBindingForUniform("proj", "mat4", "renderer.proj")    // Projection matrix
		shader.SetBindingForUniform("vwmd", "mat4", "renderer.vwmd")    // View*Model matrix
		shader.SetBindingForUniform("asp", "vec2", "renderer.aspect")   // AspectRatio
		shader.SetBindingForUniform("color", "vec3", "material.color")  // material color
		shader.SetBindingForAttribute("pxy", "vec2", "geometry.coords") // 2D offset coordinates (in CAMERA space)
	} else { // Shader for multiple instance poses ('iorg')
		var vertex_shader_code = `
			precision mediump float;
			uniform   mat4 proj;		// 3D Projection matrix
			uniform   mat4 vwmd;		// 3D View * Model matrix
			uniform   vec2 asp;			// aspect ratio, w : h
			attribute vec3 iorg;		// 3D world XYZ coordinates of the origin
			attribute vec2 pxy;			// 2D vertex XY coordinates (offset; pixels in CAMERA space)
			void main() {
				vec4 origin = proj * vwmd * vec4(iorg, 1.0);
				if (origin.w != 0.0) { origin = origin / origin.w; }
				vec2 offset = vec2(pxy.x * 2.0 / asp[0], pxy.y * 2.0 / asp[1]);
				gl_Position = vec4(origin.x + offset.x, origin.y + offset.y, origin.z, 1.0);
			}`
		var fragment_shader_code = `
			precision mediump float;
			uniform vec3 color;			// color
			void main() { 
				gl_FragColor = vec4(color.r, color.g, color.b, 1.0);
			}`
		shader, _ = self.rc.CreateShader(vertex_shader_code, fragment_shader_code)
		shader.SetBindingForUniform("proj", "mat4", "renderer.proj")       // 3D Projection matrix
		shader.SetBindingForUniform("vwmd", "mat4", "renderer.vwmd")       // 3D View*Model matrix
		shader.SetBindingForUniform("asp", "vec2", "renderer.aspect")      // AspectRatio
		shader.SetBindingForUniform("color", "vec3", "material.color")     // material color
		shader.SetBindingForAttribute("iorg", "vec3", "instance.pose:3:0") // instance pose (:<stride>:<offset>)
		shader.SetBindingForAttribute("pxy", "vec2", "geometry.coords")    // 2D offset coordinates (in CAMERA space)
	}
	shader.CheckBindings() // check validity of the shader
	return shader
}

// ----------------------------------------------------------------------------
// Sprite Marker
// ----------------------------------------------------------------------------

func (self *OverlayMarkerLayer) CreateSpriteMarker(imgpath string, color string, wh [2]float32, offref string, use_poses bool) *SceneObject {
	geometry := g2d.NewGeometry_Origin() // 2D geometry with only one vertex at (0,0)
	material, _ := self.rc.CreateMaterial(imgpath)
	material.SetColorForDrawMode(0, color)
	// wh := [2]float32{float32(material.GetTextureWH()[0]), float32(material.GetTextureWH()[1])}
	var offrot [3]float32
	switch offref {
	case "L_TOP":
		offrot = [3]float32{+wh[0] / 2, -wh[1] / 2, 0}
	case "M_TOP":
		offrot = [3]float32{0, -wh[1] / 2, 0}
	case "R_TOP":
		offrot = [3]float32{-wh[0] / 2, -wh[1] / 2, 0}
	case "L_CTR":
		offrot = [3]float32{+wh[0] / 2, 0, 0}
	case "M_CTR", "CENTER":
		offrot = [3]float32{0, 0, 0}
	case "R_CTR":
		offrot = [3]float32{-wh[0] / 2, 0, 0}
	case "L_BTM":
		offrot = [3]float32{+wh[0] / 2, +wh[1] / 2, 0}
	case "M_BTM":
		offrot = [3]float32{0, +wh[1] / 2, 0}
	case "R_BTM":
		offrot = [3]float32{-wh[0] / 2, +wh[1] / 2, 0}
	default:
	}
	shader := self.GetShaderForSpriteMarker(wh, offrot, use_poses)
	sprite := NewSceneObject(geometry, material, shader, nil, nil)
	sprite.UseBlend = true
	return sprite
}

func (self *OverlayMarkerLayer) GetShaderForSpriteMarker(wh [2]float32, offrot [3]float32, use_poses bool) gigl.GLShader {
	var shader gigl.GLShader = nil
	if !use_poses { // Shader for single instance (located at (0,0))
		var vertex_shader_code = `
			precision mediump float;
			uniform   mat4 proj;		// 3D Projection matrix
			uniform   mat4 vwmd;		// 3D View * Model matrix
			uniform   vec2  asp;		// aspect ratio, w : h
			uniform   vec2  wh;			// size of the sprite
			uniform   vec3  offr;		// offset of the label (CAMERA XY in pixel) & rotation_angle
			attribute vec2  pxy;		// vertex XY coordinates (pixels in CAMERA space)
			void main() {
				vec4 origin = proj * vwmd * vec4(0.0, 0.0, 0.0, 1.0);
				if (origin.w != 0.0) { origin = origin / origin.w; }
				vec2 offset = vec2((offr.x + pxy.x) * 2.0 / asp[0], (offr.y + pxy.y) * 2.0 / asp[1]);
				gl_Position = vec4(origin.x + offset.x, origin.y + offset.y, origin.z, 1.0);
				if (wh[0] > wh[1]) {
					gl_PointSize = wh[0];	// sprite size is its width
				} else {
					gl_PointSize = wh[1];	// sprite size is its height
				}
			}`
		var fragment_shader_code = `
			precision mediump float;
			uniform sampler2D texture;	// alphabet texture (ASCII characters from SPACE to DEL)
			uniform   vec4  color;		// color of the sprite
			uniform   vec2  wh;			// size of the sprite
			void main() {
				vec2 uv = gl_PointCoord;
				if (wh[0] > wh[1]) {
					uv[1] = (uv[1] - 0.5) * wh[0]/wh[1] + 0.5;
				} else {
					uv[0] = (uv[0] - 0.5) * wh[1]/wh[0] + 0.5;
				}
				if (uv[0] < 0.0 || uv[0] > 1.0) discard;
				if (uv[1] < 0.0 || uv[1] > 1.0) discard;
				gl_FragColor = texture2D(texture, uv) * color;
			}`
		shader, _ = self.rc.CreateShader(vertex_shader_code, fragment_shader_code)
		shader.SetBindingForUniform("proj", "mat4", "renderer.proj")            // 3D Projection matrix
		shader.SetBindingForUniform("vwmd", "mat4", "renderer.vwmd")            // 3D View*Model matrix
		shader.SetBindingForUniform("asp", "vec2", "renderer.aspect")           // AspectRatio
		shader.SetBindingForUniform("wh", "vec2", wh[:])                        // sprite size
		shader.SetBindingForUniform("offr", "vec3", offrot[:])                  // sprite offset & rotation angle
		shader.SetBindingForUniform("color", "vec4", "material.color")          // color to be multiplied with sprite texture
		shader.SetBindingForUniform("texture", "sampler2D", "material.texture") // texture sampler (unit:0)
		shader.SetBindingForAttribute("pxy", "vec2", "geometry.coords")         // offset coordinates (in CAMERA space)
	} else { // Shader for multiple instance poses ('ixyz')
		var vertex_shader_code = `
			precision mediump float;
			uniform   mat4 proj;		// 3D Projection matrix
			uniform   mat4 vwmd;		// 3D View * Model matrix
			uniform   vec2  asp;		// aspect ratio, w : h
			uniform   vec2  wh;			// size of the sprite
			uniform   vec3  offr;		// offset of the label (CAMERA XY in pixel) & rotation_angle
			attribute vec2  pxy;		// vertex XY coordinates (pixels in CAMERA space)
			attribute vec3  ixyz;		// sprite instance position (in WORLD XYZ)
			void main() {
				vec4 origin = proj * vwmd * vec4(ixyz, 1.0);
				if (origin.w != 0.0) { origin = origin / origin.w; }
				vec2 offset = vec2((offr.x + pxy.x) * 2.0 / asp[0], (offr.y + pxy.y) * 2.0 / asp[1]);
				gl_Position = vec4(origin.x + offset.x, origin.y + offset.y, origin.z, 1.0);
				if (wh[0] > wh[1]) {
					gl_PointSize = wh[0];	// sprite size is its width
				} else {
					gl_PointSize = wh[1];	// sprite size is its height
				}
			}`
		var fragment_shader_code = `
			precision mediump float;
			uniform sampler2D texture;	// alphabet texture (ASCII characters from SPACE to DEL)
			uniform   vec4  color;		// color of the sprite
			uniform   vec2  wh;			// size of the sprite
			void main() {
				vec2 uv = gl_PointCoord;
				if (wh[0] > wh[1]) {
					uv[1] = (uv[1] - 0.5) * wh[0]/wh[1] + wh[1]/2.0;
				} else {
					uv[0] = (uv[0] - 0.5) * wh[1]/wh[0] + wh[0]/2.0;
				}
				if (uv[0] < 0.0 || uv[0] > 1.0) discard;
				if (uv[1] < 0.0 || uv[1] > 1.0) discard;
				gl_FragColor = texture2D(texture, uv) * color;
			}`
		shader, _ = self.rc.CreateShader(vertex_shader_code, fragment_shader_code)
		shader.SetBindingForUniform("proj", "mat4", "renderer.proj")            // 3D Projection matrix
		shader.SetBindingForUniform("vwmd", "mat4", "renderer.vwmd")            // 3D View*Model matrix
		shader.SetBindingForUniform("asp", "vec2", "renderer.aspect")           // AspectRatio
		shader.SetBindingForUniform("wh", "vec2", wh[:])                        // sprite size
		shader.SetBindingForUniform("offr", "vec3", offrot[:])                  // sprite offset & rotation angle
		shader.SetBindingForUniform("color", "vec4", "material.color")          // color to be multiplied with sprite texture
		shader.SetBindingForUniform("texture", "sampler2D", "material.texture") // texture sampler (unit:0)
		shader.SetBindingForAttribute("ixyz", "vec3", "instance.pose:3:0")      // 3D sprite instance position (in WORLD XY)
		shader.SetBindingForAttribute("pxy", "vec2", "geometry.coords")         // 2D offset coordinates (in CAMERA space)
	}
	shader.CheckBindings() // check validity of the shader
	return shader
}
