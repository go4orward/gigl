package g3d

import (
	"fmt"
	"math"

	"github.com/go4orward/gigl/common"
)

type OrthographicProjection struct {
	wh      [2]int         // camera aspect_ratio (width & height)
	zoom    float32        // camera zoom level
	nearfar [2]float32     // near-far cliping planes (Z values in camera space)
	matrix  common.Matrix4 // projection matrix
	fov     float32        // field of view (in degree), default value is 90
	// (FieldOfView is always ZERO in OrthographicProjection)
}

func NewOrthographicProjection(wh [2]int, fov_in_clipwidth float32, zoom float32, nearfar [2]float32) *OrthographicProjection {
	// FOV 2.6 (clipwidth) will cover width 2 at any distance in full screen
	fov_in_clipwidth = float32(math.Max(0.002, math.Min(float64(fov_in_clipwidth), 200000.0)))
	projection := OrthographicProjection{wh: wh, fov: fov_in_clipwidth, zoom: zoom, nearfar: nearfar}
	projection.update_proj_matrix()
	return &projection
}

// ----------------------------------------------------------------------------
// Interface Functions
// ----------------------------------------------------------------------------

func (self *OrthographicProjection) IsPerspective() bool {
	return false
}
func (self *OrthographicProjection) IsOrthographic() bool {
	return true
}

func (self *OrthographicProjection) GetParameters() ([2]int, float32, float32, [2]float32) {
	return self.wh, self.fov, self.zoom, self.nearfar
}

func (self *OrthographicProjection) GetMatrix() *common.Matrix4 {
	return &self.matrix
}

func (self *OrthographicProjection) MultiplyVector3(v [3]float32) [3]float32 {
	return self.matrix.MultiplyVector3(v)
}

// ----------------------------------------------------------------------------
// Setting Camera Parameters
// ----------------------------------------------------------------------------

func (self *OrthographicProjection) SetAspectRatio(width int, height int) {
	self.wh = [2]int{width, height}
	self.update_proj_matrix()
}

func (self *OrthographicProjection) SetZoom(zoom float32) {
	if zoom <= 0.0 || zoom >= 1000.0 {
		fmt.Printf("Camera.SetZoom() failed : invalid zoom = %.1f\n", zoom)
	}
	self.zoom = zoom
	self.update_proj_matrix()
}

func (self *OrthographicProjection) SetNearFar(near float32, far float32) {
	self.nearfar = [2]float32{near, far}
	self.update_proj_matrix()
}

func (self *OrthographicProjection) update_proj_matrix() {
	aspect_ratio := float32(self.wh[0]) / float32(self.wh[1])
	var clip_width, clip_height float32
	if aspect_ratio > 1.0 {
		clip_width, clip_height = 2.0*aspect_ratio, float32(2.0) // CLIP space width & height (2.0)
	} else {
		clip_width, clip_height = float32(2.0), 2.0/aspect_ratio // CLIP space width (2.0) & height
	}
	// For detail, refer to http://www.songho.ca/opengl/gl_projectionmatrix.html
	// Ref: http://www.songho.ca/opengl/gl_projectionmatrix.html
	x := 2 * self.zoom / clip_width
	y := 2 * self.zoom / clip_height
	p := 1.0 / (self.nearfar[1] - self.nearfar[0])
	z := (self.nearfar[1] + self.nearfar[0]) * p
	ff := 2.0 / self.fov // fov 2.0 will cover width 2 at any distance in full screen
	self.matrix.Set(
		ff*x, 0.0, 0.0, 0.0,
		0.0, ff*y, 0.0, 0.0,
		0.0, 0.0, -2*p, -z,
		0.0, 0.0, 0.0, 1.0)
}
