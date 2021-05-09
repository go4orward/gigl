package globe

import (
	"github.com/go4orward/gigl"
	"github.com/go4orward/gigl/common"
	"github.com/go4orward/gigl/g3d"
	"github.com/go4orward/gigl/g3d/c3d"
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

func (self *WorldRenderer) Clear(globe *Globe) {
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
		self.axes = g3d.NewSceneObject_3DAxes(self.rc, length)
	}
	self.renderer.RenderSceneObject(self.axes, wcamera.gcam.GetProjMatrix(), wcamera.gcam.GetViewMatrix())
}

// ----------------------------------------------------------------------------
// Rendering the World
// ----------------------------------------------------------------------------

func (self *WorldRenderer) RenderWorld(globe *Globe, wcamera *WorldCamera) {
	if globe.IsReadyToRender() {
		// Render the Globe
		new_viewmodel := wcamera.gcam.GetViewMatrix().MultiplyToTheRight(&globe.modelmatrix)
		self.renderer.RenderSceneObject(globe.GSphere, wcamera.gcam.GetProjMatrix(), new_viewmodel)
		// Render the GlowRing (in CAMERA space)
		distance := c3d.Length(wcamera.gcam.GetCenter())
		translation := common.NewMatrix4().SetTranslation(0, 0, -distance)
		self.renderer.RenderSceneObject(globe.GlowRing, wcamera.gcam.GetProjMatrix(), translation)
	}
}
