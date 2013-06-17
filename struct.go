package sphero

import (
	"errors"
	serial "github.com/Freeflow/goserial"
)

var (
	NotImplemented = errors.New("This feature is not yet implemented")
)

type Config struct {
	Bluetooth serial.Config
}

type Response struct {
	sop1 byte
	sop2 byte
	mrsp uint8
	seq  uint8
	dlen uint8
	data []byte
	chk  uint8
}

type AsyncResponse struct {
	sop1    byte
	sop2    byte
	idCode  byte
	dlenMSB byte
	dlenLSB byte
	data    []byte
	chk     byte
}

type SimpleResponse uint8

type ColorResponse struct {
	R, G, B uint8
}

type LocatorResponse struct {
	XPos, YPos, XVel, YVel, SoG uint16
}

type CollisionResponse struct {
	X, Y, Z    int16
	Axis       uint8
	XMag, YMag int16
	Speed      uint8
	Time       time.Time
}
