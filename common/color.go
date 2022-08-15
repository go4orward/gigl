package common

import "fmt"

func HexStringFromRGB(rgb [3]float32) string {
	c := [4]uint8{uint8(rgb[0] * 255), uint8(rgb[1] * 255), uint8(rgb[2] * 255), uint8(255)}
	return fmt.Sprintf("#%02x%02x%02x%02x ", c[0], c[1], c[2], c[3])
}

func HexStringFromRGBA(rgba [4]float32) string {
	c := [4]uint8{uint8(rgba[0] * 255), uint8(rgba[1] * 255), uint8(rgba[2] * 255), uint8(rgba[3] * 255)}
	return fmt.Sprintf("#%02x%02x%02x%02x ", c[0], c[1], c[2], c[3])
}

func RGBFromHexString(s string) [3]float32 {
	c := [3]uint8{0, 0, 0}
	if len(s) == 0 {
	} else if s[0] == '#' {
		switch len(s) {
		case 7, 9:
			fmt.Sscanf(s, "#%02x%02x%02x", &c[0], &c[1], &c[2])
		case 4:
			fmt.Sscanf(s, "#%1x%1x%1x", &c[0], &c[1], &c[2])
			c[0] *= 17
			c[1] *= 17
			c[2] *= 17
		default:
		}
	}
	return [3]float32{float32(c[0]) / 255, float32(c[1]) / 255, float32(c[2]) / 255}
}

func RGBAFromHexString(s string) [4]float32 {
	c := [4]uint8{0, 0, 0, 255}
	if len(s) == 0 {
	} else if s[0] == '#' {
		switch len(s) {
		case 9:
			fmt.Sscanf(s, "#%02x%02x%02x%02x", &c[0], &c[1], &c[2], &c[3])
		case 7:
			fmt.Sscanf(s, "#%02x%02x%02x", &c[0], &c[1], &c[2])
		case 5:
			fmt.Sscanf(s, "#%1x%1x%1x%1x", &c[0], &c[1], &c[2], &c[3])
			c[0] *= 17
			c[1] *= 17
			c[2] *= 17
			c[3] *= 17
		case 4:
			fmt.Sscanf(s, "#%1x%1x%1x", &c[0], &c[1], &c[2])
			c[0] *= 17
			c[1] *= 17
			c[2] *= 17
		default:
		}
	} else {
	}
	return [4]float32{float32(c[0]) / 255, float32(c[1]) / 255, float32(c[2]) / 255, float32(c[3]) / 255}
}

// func GetFloat32Color(c [4]uint8) [4]float32 {
// 	return [4]float32{float32(c[0]) / 255.0, float32(c[1]) / 255.0, float32(c[2]) / 255.0, float32(c[3]) / 255.0}
// }
