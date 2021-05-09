package g3d

import (
	"fmt"
	"math"

	"github.com/go4orward/gigl/common"
	"github.com/go4orward/gigl/g3d/c3d"
)

type CameraProjection interface {
	IsPerspective() bool
	IsOrthographic() bool
	GetParameters() (wh [2]int, fov float32, zoom float32, nearfar [2]float32)
	GetMatrix() *common.Matrix4
	SetAspectRatio(width int, height int)
	SetZoom(zoom float32)

	MultiplyVector3([3]float32) [3]float32 // only for test
}

type Camera struct {
	// camera internal parameters
	projection CameraProjection // projection matrix (transformation from CAMERA to CLIP space)
	// camera pose
	viewmatrix common.Matrix4 // view matrix Mcw (transformation from WORLD to CAMERA space)
	center     [3]float32     // camera position in world space
	// Ref: http://www.songho.ca/opengl/gl_projectionmatrix.html
}

func NewPerspectiveCamera(wh [2]int, fov_in_degree float32, zoom float32) *Camera {
	// 'fov' 15 degree (default) will cover width 2 at distance 10 in full screen)
	camera := Camera{}
	camera.projection = NewPerspectiveProjection(wh, fov_in_degree, zoom, [2]float32{1.0, 100.0})
	camera.SetPose([3]float32{0, 0, 10}, [3]float32{0, 0, 0}, [3]float32{0, 1, 0})
	return &camera
}

func NewOrthographicCamera(wh [2]int, fov_in_clipwidth float32, zoom float32) *Camera {
	// 'fov' 2.6 (default) will cover width 2 at any distance in full screen)
	camera := Camera{}
	camera.projection = NewOrthographicProjection(wh, fov_in_clipwidth, zoom, [2]float32{1.0, 100.0})
	camera.SetPose([3]float32{0, 0, 10}, [3]float32{0, 0, 0}, [3]float32{0, 1, 0})
	return &camera
}

func (self *Camera) GetCenter() [3]float32 {
	return self.center
}

func (self *Camera) GetProjMatrix() *common.Matrix4 {
	return self.projection.GetMatrix()
}

func (self *Camera) GetViewMatrix() *common.Matrix4 {
	return &self.viewmatrix
}

func (self *Camera) ShowInfo() {
	if self.projection.IsPerspective() {
		wh, fov, zoom, nearfar := self.projection.GetParameters()
		fmt.Printf("Perspective Camera  centered at [%5.2f %5.2f %5.2f]\n", self.center[0], self.center[1], self.center[2])
		fmt.Printf("  Parameters : AspectRatio=[%d:%d]  fov=%.0f°  zoom=%.2f  nearfar=[%.2f %.2f]\n", wh[0], wh[1], fov, zoom, nearfar[0], nearfar[1])
	} else {
		wh, fov, zoom, nearfar := self.projection.GetParameters()
		fmt.Printf("Orthographic Camera  centered at [%5.2f %5.2f %5.2f]\n", self.center[0], self.center[1], self.center[2])
		fmt.Printf("  Parameters : AspectRatio=[%d:%d]  fov=%.1f  zoom=%.2f  nearfar=[%.2f %.2f]\n", wh[0], wh[1], fov, zoom, nearfar[0], nearfar[1])
	}
	p := self.projection.GetMatrix().GetElements() // Note that Matrix4 is column-major (just like WebGL)
	v := self.viewmatrix.GetElements()
	fmt.Printf("  [ %5.2f %5.2f %5.2f %7.2f ] [ %5.2f %5.2f %5.2f %7.2f ]\n", p[0], p[4], p[8], p[12], v[0], v[4], v[8], v[12])
	fmt.Printf("  [ %5.2f %5.2f %5.2f %7.2f ] [ %5.2f %5.2f %5.2f %7.2f ]\n", p[1], p[5], p[9], p[13], v[1], v[5], v[9], v[13])
	fmt.Printf("  [ %5.2f %5.2f %5.2f %7.2f ] [ %5.2f %5.2f %5.2f %7.2f ]\n", p[2], p[6], p[10], p[14], v[2], v[6], v[10], v[14])
	fmt.Printf("  [ %5.2f %5.2f %5.2f %7.2f ] [ %5.2f %5.2f %5.2f %7.2f ]\n", p[3], p[7], p[11], p[15], v[3], v[7], v[11], v[15])
}

// ----------------------------------------------------------------------------
// Camera Internal Parameters
// ----------------------------------------------------------------------------

func (self *Camera) SetAspectRatio(width int, height int) *Camera {
	// This function can be called to handle 'window.resize' event
	self.projection.SetAspectRatio(width, height)
	return self
}

func (self *Camera) SetZoom(zoom float32) *Camera {
	// This function can be called to handle 'wheel' event [ 0.01 ~ 1.0(default) ~ 100.0 ]
	zoom = float32(math.Max(0.001, math.Min(float64(zoom), 1000.0)))
	self.projection.SetZoom(zoom)
	return self
}

// ----------------------------------------------------------------------------
// Camera Pose
// ----------------------------------------------------------------------------

func (self *Camera) SetPose(from [3]float32, lookat [3]float32, up [3]float32) *Camera {
	camY := c3d.Normalize(up)
	camZ := c3d.Normalize(c3d.SubAB(from, lookat))
	camX := c3d.Normalize(c3d.CrossAB(camY, camZ))     // Normalize(), because 'up' may not be orthogonal
	camY = c3d.CrossAB(camZ, camX)                     // camY ('up' vector) is updated to make it orthogonal
	self.SetPoseWithCameraAxes(camX, camY, camZ, from) // update the viewMatrix
	return self
}

func (self *Camera) SetPoseWithCameraAxes(camX [3]float32, camY [3]float32, camZ [3]float32, center [3]float32) *Camera {
	// We are given:
	//   Rwc : rotation    from camera to world (camera pose   in world coordinates) = [ camX, camY, camZ ]
	//   Twc : translation from camera to world (camera origin in world coordinates)
	// Now we get:
	//   Tcw : translation from world to camera (world origin in camera coordinates)
	Twc := center
	Tcw := [3]float32{ // Tcw = - (Rwc)T * Twc
		-(camX[0]*Twc[0] + camX[1]*Twc[1] + camX[2]*Twc[2]),
		-(camY[0]*Twc[0] + camY[1]*Twc[1] + camY[2]*Twc[2]),
		-(camZ[0]*Twc[0] + camZ[1]*Twc[1] + camZ[2]*Twc[2])}
	// And we get the viewMatrix (transformation from world to camera)
	//   [ Rcw  Tcw ] = [ Rwc.transpose  - Rwc.transpose * Twc ]
	//   [  0    1  ]   [    0                1                ]
	self.viewmatrix.Set(
		camX[0], camX[1], camX[2], Tcw[0],
		camY[0], camY[1], camY[2], Tcw[1],
		camZ[0], camZ[1], camZ[2], Tcw[2],
		0, 0, 0, 1)
	self.center = [3]float32{Twc[0], Twc[1], Twc[2]}
	return self
}

func (self *Camera) SetPoseWithMcw(Mcw *common.Matrix4) *Camera {
	// We are given:
	//   Mcw : transformation from world to camera
	// And we get  Twc = - (Rcw)T * Tcw
	me := Mcw.GetElements()
	Tcw := [3]float32{me[12], me[13], me[14]}
	x := -(me[0]*Tcw[0] + me[1]*Tcw[1] + me[2]*Tcw[2])
	y := -(me[4]*Tcw[0] + me[5]*Tcw[1] + me[6]*Tcw[2])
	z := -(me[8]*Tcw[0] + me[9]*Tcw[1] + me[10]*Tcw[2])
	self.viewmatrix.SetCopy(Mcw)
	self.center = [3]float32{x, y, z}
	return self
}

func (self *Camera) Translate(tx float32, ty float32, tz float32) *Camera {
	translation := common.NewMatrix4().SetTranslation(-tx, -ty, -tz)
	self.viewmatrix.SetMultiplyMatrices(translation, &self.viewmatrix)
	self.center = [3]float32{self.center[0] + tx, self.center[1] + ty, self.center[2] + tz}
	return self
}

func (self *Camera) RotateByPitch(angle_in_degree float32) *Camera {
	// Rotate around CAMERA's +X axis
	rotation := common.NewMatrix4().SetRotationByAxis([3]float32{1, 0, 0}, -angle_in_degree)
	self.viewmatrix.SetMultiplyMatrices(rotation, &self.viewmatrix)
	return self
}

func (self *Camera) RotateByRoll(angle_in_degree float32) *Camera {
	// Rotate around CAMERA's -Z axis
	rotation := common.NewMatrix4().SetRotationByAxis([3]float32{0, 0, 1}, +angle_in_degree)
	self.viewmatrix.SetMultiplyMatrices(rotation, &self.viewmatrix)
	return self
}

func (self *Camera) RotateByYaw(angle_in_degree float32) *Camera {
	// Rotate around CAMERA's -Y axis
	rotation := common.NewMatrix4().SetRotationByAxis([3]float32{0, 1, 0}, +angle_in_degree)
	self.viewmatrix.SetMultiplyMatrices(rotation, &self.viewmatrix)
	return self
}

func (self *Camera) RotateAroundAxis(axis [3]float32, angle_in_degree float32) *Camera {
	rotation := common.NewMatrix4()
	rotation.SetRotationByAxis(axis, angle_in_degree)
	self.viewmatrix.SetMultiplyMatrices(rotation, &self.viewmatrix)
	return self
}

func (self *Camera) RotateAroundPoint(distance float32, h_angle float32, v_angle float32) *Camera {
	// Rotate camera around the point (0, 0, -distance) in CAMERA space
	trn0 := common.NewMatrix4().SetTranslation(0, 0, distance)
	rotY := common.NewMatrix4().SetRotationByAxis([3]float32{0, 1, 0}, +h_angle)
	rotX := common.NewMatrix4().SetRotationByAxis([3]float32{1, 0, 0}, +v_angle)
	trn1 := common.NewMatrix4().SetTranslation(0, 0, -distance)
	self.viewmatrix.SetMultiplyMatrices(trn1, rotX, rotY, trn0, &self.viewmatrix)
	return self
}

// ----------------------------------------------------------------------------
// Testing
// ----------------------------------------------------------------------------

func (self *Camera) TestDataBuffer(dbuffer []float32, stride int) {
	vertices := [][3]float32{}
	for i := 0; i < len(dbuffer); i += stride {
		v := dbuffer[i : i+3]
		v_world := [3]float32{v[0], v[1], v[2]}
		v_camera := self.viewmatrix.MultiplyVector3(v_world)
		v_clip := self.projection.MultiplyVector3(v_camera)
		vertices = append(vertices, v_world, v_camera, v_clip)
	}
	for j := 0; j < 3; j++ {
		if j == 0 {
			fmt.Printf("CameraTest: ")
		} else {
			fmt.Printf("            ")
		}
		for i := 0; i < len(vertices); i++ {
			if i%3 == 2 {
				fmt.Printf("%5.2f  ", vertices[i][j])
			} else {
				fmt.Printf("%5.2f ", vertices[i][j])
			}
		}
		fmt.Printf("\n")
	}
}
