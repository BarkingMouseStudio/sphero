package sphero

// Used to generate mask1 or mask2 for SetDataStreaming command
func ApplyMasks32(masks []uint32) uint32 {
	var mask uint32 = 0
	for _, m := range masks {
		mask |= m
	}
	return mask
}

// Computes the modulo 256 sum of the bytes, bit inverted (1's complement)
func ComputeChk(data []byte) uint8 {
	sum := 0
	for _, b := range data {
		sum += int(uint8(b))
	}
	chk := (sum % 256) ^ 0xff
	return uint8(chk)
}
