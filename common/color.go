package common

import "fmt"

func ParseHexColor(s string) [4]float32 {
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
