package globe

import (
	"github.com/go4orward/gigl/g2d"
	"github.com/go4orward/gigl/g3d"
)

type GlobeMarkLayer struct {
	ScnObjs []*g2d.SceneObject // 2D SceneObjects to be rendered (in CAMERA space)
}

func NewGlobeMarkLayer() *GlobeMarkLayer {
	mlayer := GlobeMarkLayer{}
	mlayer.ScnObjs = make([]*g2d.SceneObject, 0)
	return &mlayer
}

// ----------------------------------------------------------------------------
// Mark 						(single instance with its own geometry)
// ----------------------------------------------------------------------------

func (self *GlobeMarkLayer) AddMark(geometry *g2d.Geometry, color string) *GlobeMarkLayer {
	return self
}

// ----------------------------------------------------------------------------
// Mark with Instance Poses 	(multiple instances sharing the same geometry)
// ----------------------------------------------------------------------------

func (self *GlobeMarkLayer) AddMarkWithPoses(geometry *g2d.Geometry, poses *g3d.SceneObjectPoses) *GlobeMarkLayer {
	return self
}
