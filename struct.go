package sphero

import (
	"time"
)

// TODO: Implement better response interface

type SimpleResponse uint8 // MRSP

type Color struct {
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
