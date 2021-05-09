package g2d

import (
	"github.com/go4orward/gigl/common"
)

type Overlay interface {
	Render(pvm *common.Matrix3)
}
