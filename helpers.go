package sphero

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Used to generate mask1 or mask2 for SetDataStreaming command
func ApplyMasks32(masks []uint32) uint32 {
	var mask uint32 = 0
	for _, m := range masks {
		mask |= m
	}
	return mask
}

func ParseColor(data []byte) (*ColorResponse, error) {
	c := new(ColorResponse)
	if len(data) != 3 {
		return c, fmt.Errorf("Could not format %#x as ColorResponse", c)
	}
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &c.R)
	binary.Read(buf, binary.BigEndian, &c.G)
	binary.Read(buf, binary.BigEndian, &c.B)
	return c, nil
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
