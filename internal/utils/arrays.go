package utils

func ToByte32(slice []byte) [32]byte {
	arr := [32]byte{}
	copy(arr[:], slice)

	return arr
}

func ToByte64(slice []byte) [64]byte {
	arr := [64]byte{}
	copy(arr[:], slice)

	return arr
}
