package g3d

import (
	"fmt"
	"math"

	"github.com/go4orward/gigl/common"
	"github.com/go4orward/gigl/g3d/c3d"
)

// ----------------------------------------------------------------------------
// Camera Initialization Parameters
// ----------------------------------------------------------------------------

type CamInternalParams struct {
	WH      [2]int     // aspect ratio					default: [100 100]
	Fov     float32    // field of view (in degree) 	default:       15°  (2.0 CLIP_WIDTH for orthogonal)
	Zoom    float32    // zoom level					default:      1.0
	NearFar [2]float32 // distances to near/far plane	default:   [1 100]
	// 'fov' (field of view') : 15 degree (default) will cover width 2 at distance 10 in full screen)
}

type CamExternalPose struct {
	From [3]float32 // camera position in WORLD space	default: [0  0 10]
	At   [3]float32 // look-at target  in WORLD space	default: [0  0  0]
	Up   [3]float32 // camera up direction				default: [0  1  0]
}

// ----------------------------------------------------------------------------
// Camera
// ----------------------------------------------------------------------------

type Camera struct {
	perspective bool              // perspective or orthogonal
	ip          CamInternalParams //
	projmatrix  common.Matrix4    // projection matrix
	viewmatrix  common.Matrix4    // view matrix Mcw (transformation from WORLD to CAMERA space)
	center      [3]float32        // camera position in world space
}

func NewCamera(perspective bool, ip *CamInternalParams, ep *CamExternalPose) *Camera {
	self := Camera{}
	self.perspective = perspective
	self.ip = *ip
	self.UpdateProjectionMatrix()
	self.SetPose(ep.From, ep.At, ep.Up)
	self.CheckInternalParameters() // check internal parameters, and fix them if necessary
	return &self
}

func (self *Camera) IsPerspective() bool {
	return self.perspective
}

func (self *Camera) GetCenter() [3]float32 {
	return self.center
}

func (self *Camera) GetProjMatrix() *common.Matrix4 {
	return &self.projmatrix
}

func (self *Camera) GetViewMatrix() *common.Matrix4 {
	return &self.viewmatrix
}

func (self *Camera) ShowInfo() {
	wh, fov, zoom, nearfar := self.ip.WH, self.ip.Fov, self.ip.Zoom, self.ip.NearFar
	if self.perspective {
		fmt.Printf("Perspective Camera  centered at [%5.2f %5.2f %5.2f]\n", self.center[0], self.center[1], self.center[2])
		fmt.Printf("  Parameters : AspectRatio=[%d:%d]  fov=%.0f°  zoom=%.2f  nearfar=[%.2f %.2f]\n", wh[0], wh[1], fov, zoom, nearfar[0], nearfar[1])
	} else {
		fmt.Printf("Orthographic Camera  centered at [%5.2f %5.2f %5.2f]\n", self.center[0], self.center[1], self.center[2])
		fmt.Printf("  Parameters : AspectRatio=[%d:%d]  fov=%.1f  zoom=%.2f  nearfar=[%.2f %.2f]\n", wh[0], wh[1], fov, zoom, nearfar[0], nearfar[1])
	}
	p := self.projmatrix.GetElements() // Note that Matrix4 is column-major (just like WebGL)
	v := self.viewmatrix.GetElements()
	fmt.Printf("  [ %5.2f %5.2f %5.2f %7.2f ] [ %5.2f %5.2f %5.2f %7.2f ]\n", p[0], p[4], p[8], p[12], v[0], v[4], v[8], v[12])
	fmt.Printf("  [ %5.2f %5.2f %5.2f %7.2f ] [ %5.2f %5.2f %5.2f %7.2f ]\n", p[1], p[5], p[9], p[13], v[1], v[5], v[9], v[13])
	fmt.Printf("  [ %5.2f %5.2f %5.2f %7.2f ] [ %5.2f %5.2f %5.2f %7.2f ]\n", p[2], p[6], p[10], p[14], v[2], v[6], v[10], v[14])
	fmt.Printf("  [ %5.2f %5.2f %5.2f %7.2f ] [ %5.2f %5.2f %5.2f %7.2f ]\n", p[3], p[7], p[11], p[15], v[3], v[7], v[11], v[15])
}

func (self *Camera) CheckInternalParameters() *Camera {
	param_fixed := false
	if true { // check if near/far values are appropriate, w.r.t. camera's from/at position
		camdist := c3d.Length(self.center)
		if camdist < self.ip.NearFar[0] {
			self.ip.NearFar[0] = camdist / 2.0
			param_fixed = true
		}
		if camdist > self.ip.NearFar[1] {
			self.ip.NearFar[1] = camdist * 2.0
			param_fixed = true
		}
	}
	if param_fixed {
		fmt.Printf("Camera parameters fixed to wh:%v fov:%v zoom:%v nearfar:%v\n", self.ip.WH, self.ip.Fov, self.ip.Zoom, self.ip.NearFar)
		self.UpdateProjectionMatrix()
	}
	return self
}

// ----------------------------------------------------------------------------
// Camera Internal Parameters
// ----------------------------------------------------------------------------

func (self *Camera) SetAspectRatio(width int, height int) *Camera {
	// This function can be called to handle 'window.resize' event
	self.ip.WH = [2]int{width, height}
	self.UpdateProjectionMatrix()
	return self
}

func (self *Camera) SetZoom(zoom float32) *Camera {
	// This function can be called to handle 'wheel' event [ 0.01 ~ 1.0(default) ~ 100.0 ]
	self.ip.Zoom = float32(math.Max(0.001, math.Min(float64(zoom), 1000.0)))
	self.UpdateProjectionMatrix()
	return self
}

func (self *Camera) UpdateProjectionMatrix() {
	aspect_ratio := float32(self.ip.WH[0]) / float32(self.ip.WH[1])
	var clip_width, clip_height float32
	if aspect_ratio > 1.0 {
		clip_width, clip_height = 2.0*aspect_ratio, float32(2.0) // CLIP space width & height (2.0)
	} else {
		clip_width, clip_height = float32(2.0), 2.0/aspect_ratio // CLIP space width (2.0) & height
	}
	near, far := self.ip.NearFar[0], self.ip.NearFar[1]
	if self.perspective {
		// Ref: http://www.songho.ca/opengl/gl_projectionmatrix.html
		// this.projectionMatrix.makePerspective( left, left + width, top, top - height, nearfar[0], nearfar[1] );
		x := 2 * self.ip.Zoom * near / clip_width
		y := 2 * self.ip.Zoom * near / clip_height
		c := -(far + near) / (far - near)
		d := -2 * far * near / (far - near)
		// factor for 'field of view' (fov 15 degree will cover width 2 at distance 10 in full screen)
		// ff := float32(math.Tan(float64(self.fov/2)*InRadian)) * 70
		ff := 1.0 / float32(math.Tan(float64(self.ip.Fov/2)*InRadian))
		self.projmatrix.Set(
			ff*x, 0.0, 0.0, 0.0,
			0.0, ff*y, 0.0, 0.0,
			0.0, 0.0, c, d,
			0.0, 0.0, -1, 0.0,
		)
	} else {
		// Ref: http://www.songho.ca/opengl/gl_projectionmatrix.html
		x := 2 * self.ip.Zoom / clip_width
		y := 2 * self.ip.Zoom / clip_height
		p := 1.0 / (far - near)
		z := (far + near) * p
		ff := 2.0 / self.ip.Fov // fov 2.0 will cover width 2 at any distance in full screen
		self.projmatrix.Set(
			ff*x, 0.0, 0.0, 0.0,
			0.0, ff*y, 0.0, 0.0,
			0.0, 0.0, -2*p, -z,
			0.0, 0.0, 0.0, 1.0)
	}
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
		v_clip := self.projmatrix.MultiplyVector3(v_camera)
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
