package earth

import (
	"github.com/go4orward/gigl/g2d"
)

type OverlayMarkerLayer struct {
	ScnObjs []*g2d.SceneObject // 2D SceneObjects to be rendered (in CAMERA space)
}

func NewOverlayMarkerLayer() *OverlayMarkerLayer {
	mlayer := OverlayMarkerLayer{}
	mlayer.ScnObjs = make([]*g2d.SceneObject, 0)
	return &mlayer
}

// ----------------------------------------------------------------------------
// Mark 						(single instance with its own geometry)
// ----------------------------------------------------------------------------

func (self *OverlayMarkerLayer) AddMark(geometry *g2d.Geometry, color string) *OverlayMarkerLayer {
	return self
}
