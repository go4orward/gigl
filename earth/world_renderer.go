package earth

import (
	"github.com/go4orward/gigl"
	"github.com/go4orward/gigl/common"
	"github.com/go4orward/gigl/g2d"
	"github.com/go4orward/gigl/g3d"
)

type WorldRenderer struct {
	rc       gigl.GLRenderingContext // WebGL context
	renderer *g3d.Renderer           // Renderer for rendering 3D SceneObjects
	axes     *g3d.SceneObject        // XYZ axes for visual reference (only if required)
}

func NewWorldRenderer(rc gigl.GLRenderingContext) *WorldRenderer {
	renderer := WorldRenderer{rc: rc, renderer: g3d.NewRenderer(rc), axes: nil}
	return &renderer
}

// ----------------------------------------------------------------------------
// Clear
// ----------------------------------------------------------------------------

func (self *WorldRenderer) Clear(globe *WorldGlobe) {
	rc, c := self.rc, self.rc.GetConstants()
	rgb := globe.GetBkgColor()
	rc.GLClearColor(rgb[0], rgb[1], rgb[2], 1.0) // set clearing color
	rc.GLClear(c.COLOR_BUFFER_BIT)               // clear the canvas
	rc.GLClear(c.DEPTH_BUFFER_BIT)               // clear the canvas
}

// ----------------------------------------------------------------------------
// Rendering Axes
// ----------------------------------------------------------------------------

func (self *WorldRenderer) RenderAxes(wcamera *WorldCamera, length float32) {
	// Render three axes (X:RED, Y:GREEN, Z:BLUE) for visual reference
	if self.axes == nil {
		// create geometry with three lines for X,Y,Z  axes, with origin at (0,0,0)
		geometry := g3d.NewGeometry() // create an empty geometry
		geometry.SetVertices([][3]float32{{0, 0, 0}, {length, 0, 0}, {0, length, 0}, {0, 0, length}})
		geometry.SetEdges([][]uint32{{0, 1}, {0, 2}, {0, 3}})           // add three edges
		geometry.BuildDataBuffers(true, true, false)                    // build data buffers for vertices and edges
		shader := g3d.NewShader_3DAxes(self.rc)                         // create shader, and set its bindings
		self.axes = g3d.NewSceneObject(geometry, nil, nil, shader, nil) // set up the scene object (draw EDGES only)
	}
	self.renderer.RenderSceneObject(self.axes, wcamera.camera.GetProjMatrix(), wcamera.camera.GetViewMatrix())
}

// ----------------------------------------------------------------------------
// Rendering the World
// ----------------------------------------------------------------------------

func (self *WorldRenderer) RenderWorld(globe *WorldGlobe, wcamera *WorldCamera) {
	if globe.GSphere == nil || globe.GSphere.Material == nil {
		return
	} else if mtex := globe.GSphere.Material.(*g2d.MaterialTexture); mtex.IsLoading() {
		// It may take a long time to load the texture image (ONLY IF it's done asynchronously).
		// We choose not to draw the globe, until the world image texture is loaded and ready.
		return
	}
	// Render the Globe
	new_viewmodel := wcamera.camera.GetViewMatrix().MultiplyToTheRight(&globe.modelmatrix)
	self.renderer.RenderSceneObject(globe.GSphere, wcamera.camera.GetProjMatrix(), new_viewmodel)
	// Render the GlowRing (in CAMERA space)
	camera_center := g3d.V3d(wcamera.camera.GetCenter())
	translation := common.NewMatrix4().SetTranslation(0, 0, -camera_center.Length())
	self.renderer.RenderSceneObject(globe.GlowRing, wcamera.camera.GetProjMatrix(), translation)
}
