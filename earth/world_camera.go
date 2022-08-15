package earth

import (
	"math"

	"github.com/go4orward/gigl/g3d"
)

type WorldCamera struct {
	camera *g3d.Camera // Perspective Globe Camera
	// TODO: Add support for Orthographic projection camera
	// TODO: Add more cameras for flat world map projections
}

func NewWorldCamera(wh [2]int, fov float32, zoom float32) *WorldCamera {
	cam_ip := g3d.CamInternalParams{WH: wh, Fov: fov, Zoom: zoom, NearFar: [2]float32{1, 100}}
	cam_ep := g3d.CamExternalPose{From: [3]float32{10, 0, 0}, At: [3]float32{0, 0, 0}, Up: [3]float32{0, 0, 1}}
	perspective := true
	camera := g3d.NewCamera(perspective, &cam_ip, &cam_ep)
	self := WorldCamera{camera: camera}
	return &self
}

func (self *WorldCamera) Summary() string {
	return self.camera.Summary()
}

// ----------------------------------------------------------------------------
// Camera Internal Parameters
// ----------------------------------------------------------------------------

func (self *WorldCamera) SetAspectRatio(width int, height int) *WorldCamera {
	self.camera.SetAspectRatio(width, height)
	return self
}

func (self *WorldCamera) SetZoom(zoom float32) *WorldCamera {
	self.camera.SetZoom(zoom)
	return self
}

// ----------------------------------------------------------------------------
// Camera Pose
// ----------------------------------------------------------------------------

func (self *WorldCamera) SetPoseByLonLat(lon float32, lat float32, dist float32) *WorldCamera {
	Twc := GetXYZFromLonLat(lon, lat, dist)              // Camera center in WORLD space
	coslon := float32(math.Cos(float64(lon) * InRadian)) // cos(λ)
	sinlon := float32(math.Sin(float64(lon) * InRadian)) // sin(λ)
	camX := &g3d.V3d{-sinlon, +coslon, 0}                // this prevents UP vector singularity at poles
	camZ := Twc.Clone().Normalize()                      // camera's Z axis points backward (away from view frustum)
	camY := camZ.Cross(camX)
	self.camera.SetPoseWithCameraAxes(*camX, *camY, *camZ, *Twc)
	return self
}

func (self *WorldCamera) RotateAroundGlobe(horizontal_angle float32, vertical_angle float32) *WorldCamera {
	self.camera.RotateAroundPoint(10, horizontal_angle, vertical_angle)
	self.RotateByRollToHeadUpNorth()
	return self
}

func (self *WorldCamera) RotateByRollToHeadUpNorth() *WorldCamera {
	e := self.camera.GetViewMatrix().GetElements()
	// Note that {e[8],e[9]} is the NORTH (0,0,1) projected onto XY plane of Camera axes in WORLD space.
	if e[8]*e[8]+e[9]*e[9] > 0.01 { // Now compare {e[8],e[9]} with Y direction (90°) in CAMERA space.
		roll := 90 - float32(math.Atan2(float64(e[9]), float64(e[8])))*InDegree
		self.camera.RotateByRoll(roll)
	}
	return self
}

// ----------------------------------------------------------------------------
//
// ----------------------------------------------------------------------------
