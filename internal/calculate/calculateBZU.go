package calc

// 30 25 45
// 30 30 40
// 45 15 50

func Lost(c int) (int, int, int) {
	B := int(float64(c) * 0.075)
	Z := int(float64(c) * 0.03)
	U := int(float64(c) * 0.1125)
	return B, Z, U
}
func Set(c int) (int, int, int) {
	B := int(float64(c) * 0.075)
	Z := int(float64(c) * 0.033)
	U := int(float64(c) * 0.1)
	return B, Z, U
}
func Get(c int) (int, int, int) {
	B := int(float64(c) * 0.125)
	Z := int(float64(c) * 0.0167)
	U := int(float64(c) * 0.125)
	return B, Z, U
}
