package sphero

// Used to generate mask or mask2 in SetDataStreaming command
func applyMasks32(masks []uint32) uint32 {
	var mask uint32 = 0
	for _, m := range masks {
		mask |= m
	}
	return mask
}

// Computes the modulo 256 sum of the bytes, bit inverted
// (1's complement). The value is used as a verification
// on commands and responses.
func computeChk(data []byte) uint8 {
	sum := 0
	for _, b := range data {
		sum += int(uint8(b))
	}
	chk := (sum % 256) ^ 0xff
	return uint8(chk)
}
