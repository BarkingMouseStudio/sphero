package sphero

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

type Response struct {
	Sop1 byte
	Sop2 byte
	Mrsp byte
	Seq  uint8
	Dlen uint8
	Data []byte
	Chk  uint8
}

// Returns the appropriate error from the MRSP if any
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

func (r *Response) Color() (*Color, error) {
	c := new(Color)
	if len(r.Data) != binary.Size(c) {
		return c, fmt.Errorf("Could not parse %#x as Color", c)
	}
	buf := bytes.NewBuffer(r.Data)
	binary.Read(buf, binary.BigEndian, c)
	return c, nil
}

type AsyncResponse struct {
	Sop1   byte
	Sop2   byte
	IdCode byte
	Dlen   uint16
	Data   []byte
	Chk    uint8
}

func (r *AsyncResponse) Location() (*Location, error) {
	loc := new(Location)
	if len(r.Data) != binary.Size(loc) {
		return loc, fmt.Errorf("Could not parse %#x as Location", loc)
	}
	buf := bytes.NewBuffer(r.Data)
	binary.Read(buf, binary.BigEndian, loc)
	return loc, nil
}

func (r *AsyncResponse) PowerState() (*PowerState, error) {
	ps := new(PowerState)
	if len(r.Data) != binary.Size(ps) {
		return ps, fmt.Errorf("Could not parse %#x as PowerState", ps)
	}
	buf := bytes.NewBuffer(r.Data)
	binary.Read(buf, binary.BigEndian, ps)
	return ps, nil
}

func (r *AsyncResponse) Collision() (*Collision, error) {
	c := new(Collision)
	if len(r.Data) != binary.Size(c) {
		return c, fmt.Errorf("Could not parse %#x as Collision", c)
	}
	buf := bytes.NewBuffer(r.Data)
	binary.Read(buf, binary.BigEndian, c)
	return c, nil
}

type Quaternion struct {
	X, Y, Z, W uint16
}

type Color struct {
	R, G, B uint8
}

type PowerState struct {
	RecVer, PowerState      uint8
	BattVoltage, NumCharges uint16
	TimeSinceChg            time.Duration
}

type Location struct {
	XPos, YPos, XVel, YVel, SoG uint16
}

type Collision struct {
	X, Y, Z    int16
	Axis       uint8
	XMag, YMag int16
	Speed      uint8
	Time       time.Time
}
