package world

import "math"

// Globe is centered at (0,0) with radius 1.0
const InRadian = (math.Pi / 180.0)
const InDegree = (180.0 / math.Pi)

// ------------------------------------------------------------------------
// Longitude/Latitude  =>  X/Y/Z
// ------------------------------------------------------------------------

func GetXYZFromLonLat(lon_in_degree float32, lat_in_degree float32, radius float32) [3]float32 {
	// Get XYZ world coordinates from longitude(λ)/latitude(φ) in degree
	return GetXYZFromLL(lon_in_degree*InRadian, lat_in_degree*InRadian, radius)
}

func GetXYZFromLL(lon_in_radian float32, lat_in_radian float32, radius float32) [3]float32 {
	// Get XYZ world coordinates from longitude(λ)/latitude(φ) in radian
	lon := float64(lon_in_radian) // λ(lambda; longitude)
	lat := float64(lat_in_radian) // φ(phi;    latitude )
	return [3]float32{
		radius * float32(math.Cos(lon)*math.Cos(lat)), // dist * cosλ * cosφ;
		radius * float32(math.Sin(lon)*math.Cos(lat)), // dist * sinλ * cosφ;
		radius * float32(math.Sin(lat))}               // dist * sinφ;
}

// ------------------------------------------------------------------------
// X/Y/Z  =>  Longitude/Latitude/R
// ------------------------------------------------------------------------

// TODO : Compare the implementation of 'getLLFromXYZ()' with the below:
// Ref: https://en.wikipedia.org/wiki/Vector_fields_in_cylindrical_and_spherical_coordinates
//   radius  = Math.sqrt(x*x + y*y + z*z);
//        λ  = arctan( y / x )              0 <=    λ    <= 2π
//   (π - φ) = arccos( z / radius )         0 <= (π - φ) <=  π

func GetLLFromXYZ(x float32, y float32, z float32) [3]float32 {
	// Get longitude(λ)/latitude(φ) in radian + radius from XYZ world coordinates
	radius := math.Sqrt(float64(x*x + y*y + z*z))
	if radius == 0 {
		return [3]float32{0, 0, 0}
	}
	lon := math.Asin(float64(z) / radius)
	lat := math.Atan2(float64(y), float64(x))
	// let cosφ = Math.cos(φ);
	// let cosλ = (x / radius) / cosφ;
	// let sinλ = (y / radius) / cosφ;
	// let λ = Math.asin( sinλ );    // -PI/2 ~ +PI/2
	// λ = (cosλ >= 0 ? λ : (sinλ > 0 ? (+Math.PI - λ) : (-Math.PI - λ)));
	return [3]float32{float32(lon), float32(lat), float32(radius)}
}

func GetLonLatFromXYZ(x float32, y float32, z float32) [3]float32 {
	// Get longitude(λ)/latitude(φ) in degree + radius from XYZ world coordinates
	llr := GetLLFromXYZ(x, y, z)
	return [3]float32{llr[0] * InDegree, llr[1] * InDegree, llr[2]}
}
