package g2d

import (
	"fmt"

	"github.com/go4orward/gigl/common"
)

type MaterialColors struct {
	colors [4][4]float32 // color ([0]:common, [1]:vert, [2]:edge, [3]:face)
}

func NewMaterialColors(colors ...any) *MaterialColors {
	mc := MaterialColors{}
	for i := 0; i < 4; i++ {
		rgba := [4]float32{0, 0, 0, 0}
		if i < len(colors) {
			switch colors[i].(type) {
			case [3]float32:
				rgb := colors[i].([3]float32)
				rgba = [4]float32{rgb[0], rgb[1], rgb[2], 1.0}
			case [4]float32:
				rgba = colors[i].([4]float32)
			case string:
				rgba = common.RGBAFromHexString(colors[i].(string))
			}
			mc.colors[i] = rgba
		} else if i > 0 {
			mc.colors[i] = mc.colors[0]
		}
	}
	return &mc
}

func (self *MaterialColors) MaterialSummary() string {
	colors := ""
	for i := 0; i < len(self.colors); i++ {
		colors += common.HexStringFromRGBA(self.colors[i]) + " "
	}
	return fmt.Sprintf("MaterialColor %s\n", colors)
}

func (self *MaterialColors) GetDrawModeColor(draw_mode int) [4]float32 {
	return self.colors[draw_mode]
}

func (self *MaterialColors) SetColorForDrawMode(draw_mode int, color string) *MaterialColors {
	// 'draw_mode' :  0:common, 1:vertex, 2:edges, 3:faces
	if color != "" {
		rgba := common.RGBAFromHexString(color)
		switch draw_mode {
		case 1:
			self.colors[1] = rgba // vertex color
		case 2:
			self.colors[2] = rgba // edge color
		case 3:
			self.colors[3] = rgba // face color
		default:
			self.colors[0] = rgba // otherwise
			self.colors[1] = rgba
			self.colors[2] = rgba
			self.colors[3] = rgba
		}
	}
	return self
}
