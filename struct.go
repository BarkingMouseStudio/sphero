package sphero

import (
	"errors"
	serial "github.com/Freeflow/goserial"
	"time"
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
	mrsp byte
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

type SimpleResponse uint8 // MRSP

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
