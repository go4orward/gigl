package g3d

import (
	"fmt"
	"math"

	"github.com/go4orward/gigl/common"
)

type PerspectiveProjection struct {
	wh      [2]int         // camera aspect_ratio (width & height)
	zoom    float32        // camera zoom level
	nearfar [2]float32     // near-far cliping planes (Z values in camera space)
	matrix  common.Matrix4 // projection matrix
	fov     float32        // field of view (in degree), default value is 90
}

func NewPerspectiveProjection(wh [2]int, fov_in_degree float32, zoom float32, nearfar [2]float32) *PerspectiveProjection {
	// FOV 15 (degree) will cover width 2 at distance 10 in full screen
	fov_in_degree = float32(math.Max(10, math.Min(float64(fov_in_degree), 170)))
	projection := PerspectiveProjection{wh: wh, fov: fov_in_degree, zoom: zoom, nearfar: nearfar}
	projection.update_proj_matrix()
	return &projection
}

// ----------------------------------------------------------------------------
// Interface Functions
// ----------------------------------------------------------------------------

func (self *PerspectiveProjection) IsPerspective() bool {
	return true
}
func (self *PerspectiveProjection) IsOrthographic() bool {
	return false
}

func (self *PerspectiveProjection) GetParameters() ([2]int, float32, float32, [2]float32) {
	return self.wh, self.fov, self.zoom, self.nearfar
}

func (self *PerspectiveProjection) GetMatrix() *common.Matrix4 {
	return &self.matrix
}

func (self *PerspectiveProjection) MultiplyVector3(v [3]float32) [3]float32 {
	return self.matrix.MultiplyVector3(v)
}

// ----------------------------------------------------------------------------
// Setting Camera Parameters
// ----------------------------------------------------------------------------

func (self *PerspectiveProjection) SetAspectRatio(width int, height int) {
	self.wh = [2]int{width, height}
	self.update_proj_matrix()
}

func (self *PerspectiveProjection) SetZoom(zoom float32) {
	if zoom <= 0.0 || zoom >= 1000.0 {
		fmt.Printf("Camera.SetZoom() failed : invalid zoom = %.1f\n", zoom)
	}
	self.zoom = zoom
	self.update_proj_matrix()
}

func (self *PerspectiveProjection) update_proj_matrix() {
	aspect_ratio := float32(self.wh[0]) / float32(self.wh[1])
	var clip_width, clip_height float32
	if aspect_ratio > 1.0 {
		clip_width, clip_height = 2.0*aspect_ratio, float32(2.0) // CLIP space width & height (2.0)
	} else {
		clip_width, clip_height = float32(2.0), 2.0/aspect_ratio // CLIP space width (2.0) & height
	}
	// For detail, refer to http://www.songho.ca/opengl/gl_projectionmatrix.html
	// this.projectionMatrix.makePerspective( left, left + width, top, top - height, nearfar[0], nearfar[1] );
	x := 2 * self.zoom * self.nearfar[0] / clip_width
	y := 2 * self.zoom * self.nearfar[0] / clip_height
	c := -(self.nearfar[1] + self.nearfar[0]) / (self.nearfar[1] - self.nearfar[0])
	d := -2 * self.nearfar[1] * self.nearfar[0] / (self.nearfar[1] - self.nearfar[0])
	// factor for 'field of view' (fov 15 degree will cover width 2 at distance 10 in full screen)
	// ff := float32(math.Tan(float64(self.fov/2)*InRadian)) * 70
	ff := 1.0 / float32(math.Tan(float64(self.fov/2)*InRadian))
	self.matrix.Set(
		ff*x, 0.0, 0.0, 0.0,
		0.0, ff*y, 0.0, 0.0,
		0.0, 0.0, c, d,
		0.0, 0.0, -1, 0.0,
	)
}
