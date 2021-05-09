package g3d

import "github.com/go4orward/gigl/common"

type Overlay interface {
	Render(proj *common.Matrix4, view *common.Matrix4)
}
