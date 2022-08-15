package g2d

import (
	"fmt"

	"github.com/go4orward/gigl/common"
)

type MaterialAlphabetTexture struct {
	font_family   string     // font family name  (USE FIXED-WIDTH FONT LIKE "Courier New")
	font_size     int        // font size (12, 16, 21, 25, etc)
	font_rgb      [3]float32 // font color
	font_outlined bool       // flag to use outlined font
	texture       any        // texture
	texture_wh    [2]int     // texture size
	alphabet_cwh  [2]float32 // character width & height of ALPHABET_STRING
	err           error      //
}

func NewMaterialAlphabetTexture(fontfamily string, fontsize int, color string, outlined bool) *MaterialAlphabetTexture {
	rgb := common.RGBFromHexString(color)
	mab := &MaterialAlphabetTexture{font_family: fontfamily, font_size: fontsize, font_rgb: rgb, font_outlined: outlined}
	// NOTE THAT ACTUAL PREPARATION OF THIS MATERIAL IS IMPLEMENTED FOR EACH ENVIRONMENT,
	//   IT IS AUTOMATICALLY CALLED BY Renderer USING GLRenderingContext.LoadMaterial(material).
	return mab
}

// ----------------------------------------------------------------------------
// Material Interface Functions
// ----------------------------------------------------------------------------

func (self *MaterialAlphabetTexture) MaterialSummary() string {
	return fmt.Sprintf("MaterialAlphabetTexture %dx%d \n", self.texture_wh[0], self.texture_wh[1])
}

// ----------------------------------------------------------------------------
// MaterialTexture Interface Functions
// ----------------------------------------------------------------------------

func (self *MaterialAlphabetTexture) GetTexture() any {
	return self.texture
}

func (self *MaterialAlphabetTexture) SetTexture(texture any) {
	self.texture = texture
}

func (self *MaterialAlphabetTexture) GetTextureWH() [2]int {
	return self.texture_wh
}

func (self *MaterialAlphabetTexture) SetTextureWH(wh [2]int) {
	self.texture_wh = wh
}

func (self *MaterialAlphabetTexture) GetTextureRGB() [3]float32 {
	return self.font_rgb
}

func (self *MaterialAlphabetTexture) SetTextureRGB(color any) {
}

func (self *MaterialAlphabetTexture) IsTextureReady() bool {
	return (self.texture != nil && self.texture_wh[0] > 0 && self.texture_wh[1] > 0)
}

func (self *MaterialAlphabetTexture) IsTextureLoading() bool {
	return false
}

// ----------------------------------------------------------------------------
// MaterialAlphabetTexture Interface Functions
// ----------------------------------------------------------------------------

const _ALPHABET_STRING = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~?°"

func (self *MaterialAlphabetTexture) GetFontFamily() string {
	return self.font_family
}

func (self *MaterialAlphabetTexture) GetFontSize() int {
	return self.font_size
}

func (self *MaterialAlphabetTexture) GetFontColor() string {
	return common.HexStringFromRGB(self.font_rgb)
}

func (self *MaterialAlphabetTexture) GetFontOutlined() bool {
	return self.font_outlined
}

func (self *MaterialAlphabetTexture) GetAlphabetString() string {
	return _ALPHABET_STRING
}

func (self *MaterialAlphabetTexture) GetAlaphabetLength() int {
	return len([]rune(_ALPHABET_STRING)) // total rune length of _ALPHABET_STRING
}

func (self *MaterialAlphabetTexture) SetAlphabetWH(wh [2]float32) {
	self.alphabet_cwh = wh
}

func (self *MaterialAlphabetTexture) GetAlaphabetCharacterWH(scale float32) [2]float32 {
	return [2]float32{self.alphabet_cwh[0] * scale, self.alphabet_cwh[1] * scale}
}

func (self *MaterialAlphabetTexture) GetAlaphabetCharacterIndex(c rune) int {
	if c >= 32 && c < 127 { // 32: <SPACE>, 127: <DEL>
		return int(c) - 32
	} else {
		alphabet_rune := []rune(_ALPHABET_STRING)
		for i := 127 - 32; i < len(alphabet_rune); i++ {
			if c == alphabet_rune[i] {
				return i
			}
		}
		return (127 - 32) // '?'
	}
}
