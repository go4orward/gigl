package gigl

type GLMaterial interface {
	MaterialSummary() string
}

type GLMaterialColors interface {
	MaterialSummary() string

	GetDrawModeColor(draw_mode int) [4]float32
}

type GLMaterialTexture interface {
	MaterialSummary() string

	GetTexture() any
	SetTexture(texture any)
	GetTextureWH() [2]int
	SetTextureWH(wh [2]int)
	GetTextureRGB() [3]float32
	SetTextureRGB(color any)
	IsTextureReady() bool
	IsTextureLoading() bool
}
