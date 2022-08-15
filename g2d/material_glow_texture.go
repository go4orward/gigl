package g2d

import (
	"fmt"
	"math"

	"github.com/go4orward/gigl/common"
)

type MaterialGlowTexture struct {
	glow_rgb   [3]float32 //
	pixbuf     []uint8    //
	texture    any        // texture (js.Value for WebGL, uint32 for OpenGL)
	texture_wh [2]int     // texture size
	err        error      //
}

func NewMaterialGlowTexture(color any) *MaterialGlowTexture {
	glow_rgb := [3]float32{}
	switch color.(type) {
	case [3]float32:
		glow_rgb = color.([3]float32)
	case string:
		glow_rgb = common.RGBFromHexString(color.(string))
	}
	mtex := MaterialGlowTexture{glow_rgb: glow_rgb}
	// NOTE THAT ACTUAL PREPARATION OF THIS MATERIAL IS IMPLEMENTED FOR EACH ENVIRONMENT,
	//   IT IS AUTOMATICALLY CALLED BY Renderer USING GLRenderingContext.LoadMaterial(material).
	return &mtex
}

func (self *MaterialGlowTexture) MaterialSummary() string {
	return fmt.Sprintf("MaterialGlowTexture %dx%d \n", self.texture_wh[0], self.texture_wh[1])
}

// ----------------------------------------------------------------------------
// TEXTURE
// ----------------------------------------------------------------------------

func (self *MaterialGlowTexture) GetTexture() any {
	return self.texture
}

func (self *MaterialGlowTexture) SetTexture(texture any) {
	self.texture = texture
}

func (self *MaterialGlowTexture) GetTextureWH() [2]int {
	return self.texture_wh
}

func (self *MaterialGlowTexture) SetTextureWH(wh [2]int) {
	self.texture_wh = wh
}

func (self *MaterialGlowTexture) GetTextureRGB() [3]float32 {
	return self.glow_rgb
}

func (self *MaterialGlowTexture) SetTextureRGB(color any) {
}

func (self *MaterialGlowTexture) IsTextureReady() bool {
	return (self.texture != nil && self.texture_wh[0] > 0 && self.texture_wh[1] > 0)
}

func (self *MaterialGlowTexture) IsTextureLoading() bool {
	return false
}

// ----------------------------------------------------------------------------
// Loading Texture Image
// ----------------------------------------------------------------------------

func (self *MaterialGlowTexture) LoadTexture(callback_texture_loaded func(pixbuf []uint8, wh [2]int)) {
	set_pixbuf_with_rgba := func(pbuffer []uint8, idx int, R uint8, G uint8, B uint8, A uint8) {
		pbuffer[idx+0] = R
		pbuffer[idx+1] = G
		pbuffer[idx+2] = B
		pbuffer[idx+3] = A
	}
	// build pixel buffer for glow effect
	const width, height = 34, 2 // it has to be non-power-of-two texture with gl.NEAREST
	pixbuf := make([]uint8, (width*height)*4)
	for u := 1; u < width-1; u++ { // first (i==0) and last (i==width-1) pixel is ZERO
		ratio := (float32(u-1) / float32(width-2))
		if true { // diminishing glow for the first row (v == 0)  [ 1.0 ~ 0.5 ~ 0.0 ]
			intensity := 1.0 - ratio
			ii := intensity * intensity * 255
			r, g, b := uint8(ii*self.glow_rgb[0]), uint8(ii*self.glow_rgb[1]), uint8(ii*self.glow_rgb[2])
			set_pixbuf_with_rgba(pixbuf, (u)*4, r, g, b, uint8(ii))
		}
		if true { // glow on both side for the second row (v == 1)  [ 0.0 ~ 1.0 ~ 0.0 ]
			intensity := 1.0 - float32(math.Abs(float64(ratio*2-1)))
			ii := intensity * intensity * 255
			r, g, b := uint8(ii*self.glow_rgb[0]), uint8(ii*self.glow_rgb[1]), uint8(ii*self.glow_rgb[2])
			set_pixbuf_with_rgba(pixbuf, (width+u)*4, r, g, b, uint8(ii))
		}
	}
	self.pixbuf = pixbuf
	self.texture_wh = [2]int{width, height} // CLAMP_TO_EDGE & NEAREST(not LINEAR) for NON-POWER-OF-2 textures
	if callback_texture_loaded != nil {
		callback_texture_loaded(self.pixbuf, self.texture_wh)
	}
}
