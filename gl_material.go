package gigl

type GLMaterial interface {
	ShowInfo()

	// COLOR
	SetColorForDrawMode(draw_mode int, color string) GLMaterial
	SetDrawModeColor(draw_mode int, color [4]float32) GLMaterial
	GetDrawModeColor(draw_mode int) [4]float32

	// TEXTURE
	GetTexture() interface{}
	GetTextureWH() [2]int
	IsTextureReady() bool
	IsTextureLoading() bool

	// // Glow texture
	InitializeGlowTexture(color string)
	// // Alphabet texture
	InitializeAlphabetTexture(color string, fontsize int, outlined bool)
	GetAlaphabetLength() int
	GetAlaphabetCharacterWH(scale float32) [2]float32
	GetAlaphabetCharacterIndex(c rune) int
}
