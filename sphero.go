// TODO: Better response interface
package sphero

import (
	"bytes"
	"encoding/binary"
	"fmt"
	serial "github.com/Freeflow/goserial"
	"io"
	"os"
	"syscall"
	"time"
)

type Sphero struct {
	conf  *Config
	conn  io.ReadWriteCloser
	seq   uint8
	res   map[uint8]chan<- *Response
	kill  chan bool
	async chan<- *AsyncResponse
}

func NewSphero(conf *Config, async chan<- *AsyncResponse) (*Sphero, error) {
	var conn io.ReadWriteCloser
	var err error
	if conn, err = serial.OpenPort(&conf.Bluetooth); err != nil {
		return nil, err
	}

	s := &Sphero{
		conf:  conf,
		conn:  conn,
		seq:   0,
		res:   make(map[uint8]chan<- *Response),
		kill:  make(chan bool, 1),
		async: async,
	}

	go s.listen()

	return s, nil
}

func (s *Sphero) parse(buf []byte) (n int, err error) {
	sop1 := buf[0]

	if sop1 != SOP1 {
		err = fmt.Errorf("SOP1 must be FFh but got %p", sop1)
		return
	}

	sop2 := buf[1]

	switch sop2 {
	case SOP2_ANSWER:
		var r *Response
		r.sop1 = sop1
		r.sop2 = sop2
		r.mrsp = buf[2]

		packetBuf := bytes.NewBuffer(buf)

		// Skip the SOP bytes and MRSP
		packetBuf.Next(3)

		binary.Read(packetBuf, binary.BigEndian, &r.seq)
		binary.Read(packetBuf, binary.BigEndian, &r.dlen)

		dataStart := 5

		/*
			We haven't read enough data yet to finish parsing this message, we'll wait
			for more.
		*/
		if len(buf) < int(r.dlen)+dataStart {
			return
		}

		dataEnd := dataStart + (int(r.dlen) - 1)
		r.data = buf[dataStart:dataEnd]

		chkBytes := buf[2:dataEnd]
		computedChk := computeChk(chkBytes)

		// Skip over the data bytes
		packetBuf.Next(int(r.dlen) - 1)
		binary.Read(packetBuf, binary.BigEndian, &r.chk)

		n = int(dataEnd) + 1

		/*
			Verify the check matches what we expect. If it doesn't match we return an
			error and let the listener throw away the bad data.
		*/
		if computedChk != r.chk {
			err = fmt.Errorf("Invalid check: expected %#x but got %#x", r.chk, computedChk)
			return
		}

		/*
			Send the response over the channel associated with the seq number, if it
			exists.
		*/
		if res, ok := s.res[r.seq]; ok {
			res <- r
		}
	case SOP2_ASYNC:
		if len(buf) < 7 {
			return
		}

		var r *AsyncResponse
		r.sop1 = sop1
		r.sop2 = sop2
		r.idCode = buf[2]

		packetBuf := bytes.NewBuffer(buf)

		// Skip the SOP bytes and ID CODE
		packetBuf.Next(3)

		binary.Read(packetBuf, binary.BigEndian, &r.dlenMSB)
		binary.Read(packetBuf, binary.LittleEndian, &r.dlenLSB)

		dataStart := 5

		if len(buf) < int(r.dlenMSB)+dataStart {
			return
		}

		dataEnd := dataStart + (int(r.dlenMSB) - 1)
		r.data = buf[dataStart:dataEnd]

		chkBytes := buf[2:dataEnd]
		computedChk := computeChk(chkBytes)

		// Skip over the data bytes
		packetBuf.Next(int(r.dlenMSB) - 1)
		binary.Read(packetBuf, binary.BigEndian, &r.chk)

		n = int(dataEnd) + 1

		/*
			Verify the check matches what we expect. If it doesn't match we return an
			error and let the listener throw away the bad data.
		*/
		if computedChk != r.chk {
			err = fmt.Errorf("Invalid async check: expected %#x but got %#x", r.chk, computedChk)
			return
		}

		s.async <- r
	default:
		err = fmt.Errorf("Unexpected SOP2, should be %#x or %#x but got %#x", SOP2_ANSWER, SOP2_ASYNC, sop2)
		n = 1 // Chomp 1 byte and maybe we'll recover
	}
	return
}

func (s *Sphero) listen() {
	var data []byte
	var buf []byte
	var n int
	var err error

	for {
		select {
		case <-s.kill:
			return
		default:
			data = make([]byte, 256)

			/*
				Since EOF errors are expected when the Sphero indicates it doesn't
				expect to send more data (e.g. all responses have been sent for commands
				received so far and async responses are turned off).

				EBADF errors are also expected if we've initiated a `Close` while `Read`
				was blocking.
			*/
			if n, err = s.Read(data); err != nil && err != io.EOF {
				if pathErr, ok := err.(*os.PathError); ok {
					if pathErr.Err != syscall.EBADF {
						panic(pathErr)
					}
				} else {
					panic(err)
				}
			}

			/*
				Trim the 256 byte data by the number of bytes actually read and append
				it to our buffer.
			*/
			if n > 0 {
				data = data[:n]
				buf = append(buf, data...)
			}

			/*
				If our buffer is too short to form a meaningful response, we wait until
				we've read more. Answers need to be at least 6 bytes, async responses
				need to be at least 7.
			*/
			if len(buf) < 6 {
				continue
			}

			if n, err = s.parse(buf); err != nil {
				fmt.Println("Parse:", err)
			}

			// Trim our buffer by the number of bytes successfully parsed.
			if n > 0 {
				buf = buf[n:]
			}
		}
	}
}

// Implement io.ReadWriteCloser

// Implement io.Closer
func (s *Sphero) Close() error {
	s.kill <- true // Signal to kill our goroutine
	return s.conn.Close()
}

// Implement io.Writer
func (s *Sphero) Write(data []byte) (int, error) {
	return s.conn.Write(data)
}

// Implement io.Reader
func (s *Sphero) Read(data []byte) (int, error) {
	return s.conn.Read(data)
}

func (s *Sphero) Send(did, cid uint8, data []byte, res chan<- *Response) error {
	s.seq++
	s.res[s.seq] = res

	var buf bytes.Buffer
	buf.Write([]byte{SOP1})                                  // SOP1
	buf.Write([]byte{SOP2_ANSWER})                           // SOP2
	buf.Write([]byte{did})                                   // DID
	buf.Write([]byte{cid})                                   // CID
	binary.Write(&buf, binary.BigEndian, s.seq)              // SEQ
	binary.Write(&buf, binary.BigEndian, uint8(len(data)+1)) // DLEN

	if data != nil {
		buf.Write(data)
	}

	chk := computeChk(buf.Bytes()[2:buf.Len()])
	binary.Write(&buf, binary.BigEndian, chk) // DLEN

	_, err := s.Write(buf.Bytes())
	return err
}

// Core

func (s *Sphero) Ping(res chan<- *Response) error {
	return s.Send(DID_CORE, CMD_PING, nil, res)
}

func (s *Sphero) Sleep(wakeup time.Duration, macro uint8, orbBasic uint16, res chan<- *Response) error {
	var data bytes.Buffer
	binary.Write(&data, binary.BigEndian, uint16(wakeup))
	binary.Write(&data, binary.BigEndian, macro)
	binary.Write(&data, binary.BigEndian, orbBasic)
	return s.Send(DID_CORE, CMD_SLEEP, data.Bytes(), res)
}

// Sphero

func (s *Sphero) SetHeading() error {
	return NotImplemented
}

func (s *Sphero) SetStabilization() error {
	return NotImplemented
}

func (s *Sphero) SetRotationRate(rate uint8, res chan<- *Response) error {
	return NotImplemented
}

func (s *Sphero) SelfLevel() error {
	return NotImplemented
}

func (s *Sphero) SetDataStreaming(n, m int16, mask uint32, pcnt uint8, mask2 uint32, res chan<- *Response) error {
	var data bytes.Buffer
	binary.Write(&data, binary.BigEndian, n)
	binary.Write(&data, binary.BigEndian, m)
	binary.Write(&data, binary.BigEndian, mask)
	binary.Write(&data, binary.BigEndian, pcnt)
	binary.Write(&data, binary.BigEndian, mask2)
	return s.Send(DID_SPHERO, CMD_SET_DATA_STREAMING, data.Bytes(), res)
}

func (s *Sphero) ConfigureCollisionDetection(method, xThreshold, xSpeed, yThreshold, ySpeed, deadTime uint8, res chan<- *Response) error {
	return NotImplemented
}

func (s *Sphero) ConfigureLocator(flags uint8, x, y, yawTare uint16, res chan<- *Response) error {
	return NotImplemented
}

func (s *Sphero) ReadLocator(res chan<- *Response) error {
	return NotImplemented
}

func (s *Sphero) SetRGBLEDOutput(red, green, blue uint8, res chan<- *Response) error {
	var data bytes.Buffer
	binary.Write(&data, binary.BigEndian, red)
	binary.Write(&data, binary.BigEndian, green)
	binary.Write(&data, binary.BigEndian, blue)

	// User flag - 0x01 would set "user LED color"
	data.Write([]byte{0x00})

	return s.Send(DID_SPHERO, CMD_SET_RGB_LED, data.Bytes(), res)
}

func (s *Sphero) SetBackLEDOutput(brightness uint8, res chan<- *Response) error {
	var data bytes.Buffer
	binary.Write(&data, binary.BigEndian, brightness)
	return s.Send(DID_SPHERO, CMD_SET_BACK_LED, data.Bytes(), res)
}

func (s *Sphero) GetRGBLED(res chan<- *Response) error {
	return s.Send(DID_SPHERO, CMD_GET_RGB_LED, nil, res)
}

func (s *Sphero) Roll() error {
	return NotImplemented
}

func (s *Sphero) SetRawMotorValues() error {
	return NotImplemented
}
