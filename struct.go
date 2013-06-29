package sphero

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Represents a command response.
type Response struct {
	Sop1 byte
	Sop2 byte
	Mrsp byte
	Seq  uint8
	Dlen uint8
	Data []byte
	Chk  uint8
}

// Returns the appropriate error from the message response (MRSP) field, if any.
func (r *Response) Error() (err error) {
	switch r.Mrsp {
	case ORBOTIX_RSP_CODE_OK:
		err = nil
	case ORBOTIX_RSP_CODE_EGEN:
		err = GeneralError
	case ORBOTIX_RSP_CODE_ECHKSUM:
		err = ChecksumError
	case ORBOTIX_RSP_CODE_EFRAG:
		err = CommandFragmentError
	case ORBOTIX_RSP_CODE_EBAD_CMD:
		err = UnknownCommandError
	case ORBOTIX_RSP_CODE_EUNSUPP:
		err = UnsupportedCommandError
	case ORBOTIX_RSP_CODE_EBAD_MSG:
		err = BadMessageFormatError
	case ORBOTIX_RSP_CODE_EPARAM:
		err = InvalidParametersError
	case ORBOTIX_RSP_CODE_EEXEC:
		err = FailedExecuteCommandError
	case ORBOTIX_RSP_CODE_EBAD_DID:
		err = UnknownDeviceError
	case ORBOTIX_RSP_CODE_POWER_NOGOOD:
		err = PowerTooLowError
	case ORBOTIX_RSP_CODE_PAGE_ILLEGAL:
		err = IllegalPageError
	case ORBOTIX_RSP_CODE_FLASH_FAIL:
		err = FlashFailError
	case ORBOTIX_RSP_CODE_MA_CORRUPT:
		err = ApplicationCorruptError
	case ORBOTIX_RSP_CODE_MSG_TIMEOUT:
		err = MessageTimeoutError
	default:
		err = UnknownError
	}
	return
}

/*
	Parses the data portion of the response into a PowerState struct.
	This method may be inaccurate or fail if the data doesn't represent a color.
*/
func (r *Response) PowerState() (*PowerState, error) {
	ps := new(PowerState)
	if len(r.Data) != binary.Size(ps) {
		return ps, fmt.Errorf("Could not parse %#x as PowerState", ps)
	}
	buf := bytes.NewBuffer(r.Data)
	binary.Read(buf, binary.BigEndian, ps)
	return ps, nil
}

/*
	Parses the data portion of the response into a Color struct.
	This method may be inaccurate or fail if the data doesn't represent a color.
*/
func (r *Response) Color() (*Color, error) {
	c := new(Color)
	if len(r.Data) != binary.Size(c) {
		return c, fmt.Errorf("Could not parse %#x as Color", c)
	}
	buf := bytes.NewBuffer(r.Data)
	binary.Read(buf, binary.BigEndian, c)
	return c, nil
}

/*
	Represents an async response from one of the async data commands:
	 	- SetDataStreaming
	 	- ConfigureLocator
	 	- ConfigureCollisionDetection
*/
type AsyncResponse struct {
	Sop1   byte
	Sop2   byte
	IdCode byte
	Dlen   uint16
	Data   []byte
	Chk    uint8
}

/*
	Unpacks sensor data from an async response.
	This method may be inaccurate or fail if the data doesn't represent sensor data.
*/
func (r *AsyncResponse) Sensors(d interface{}) error {
	buf := bytes.NewBuffer(r.Data)
	return binary.Read(buf, binary.BigEndian, d)
}

/*
	Parses the data portion of the async response into a Location struct.
	This method may be inaccurate or fail if the data doesn't represent a
	location.
*/
func (r *AsyncResponse) Location() (*Location, error) {
	loc := new(Location)
	if len(r.Data) != binary.Size(loc) {
		return loc, fmt.Errorf("Could not parse %#x as Location", loc)
	}
	buf := bytes.NewBuffer(r.Data)
	binary.Read(buf, binary.BigEndian, loc)
	return loc, nil
}

/*
	Parses the data portion of the async response into a Collision struct.
	This method may be inaccurate or fail if the data doesn't represent a
	collision.
*/
func (r *AsyncResponse) Collision() (*Collision, error) {
	c := new(Collision)
	if len(r.Data) != binary.Size(c) {
		return c, fmt.Errorf("Could not parse %#x as Collision", r.Data)
	}
	buf := bytes.NewBuffer(r.Data)
	binary.Read(buf, binary.BigEndian, c)
	return c, nil
}

// Simple Color struct. See GetRGBLED.
type Color struct {
	R, G, B uint8
}

// Represents the power state of the Sphero. See SetPowerNotification.
type PowerState struct {
	RecVer, PowerState                    uint8
	BattVoltage, NumCharges, TimeSinceChg uint16
}

// Represents collision data from the Locator service. See ConfigureLocator.
type Location struct {
	XPos, YPos, XVel, YVel, SoG uint16
}

// Represents collision data from the Collision service. See ConfigureCollisionDetection.
type Collision struct {
	X, Y, Z    int16
	Axis       uint8
	XMag, YMag int16
	Speed      uint8
	TimeStamp  int32
}
