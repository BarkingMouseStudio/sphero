// TODO: Better read stability
// TODO: SetDataStreaming
// TODO: Better response interface
package sphero

import (
	"bytes"
	"encoding/binary"
	"fmt"
	serial "github.com/Freeflow/goserial"
	"io"
	"time"
)

type Sphero struct {
	conf  *Config
	conn  io.ReadWriteCloser
	seq   uint8
	res   map[uint8]chan<- *Response
	kill  chan bool
	async chan<- interface{}
}

func NewSphero(conf *Config, async chan<- interface{}) (*Sphero, error) {
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
	if len(buf) < 2 {
		return
	}

	sop1 := buf[0]

	if sop1 != SOP1 {
		err = fmt.Errorf("SOP1 must be FFh but got %p", sop1)
		return
	}

	sop2 := buf[1]

	switch sop2 {
	case SOP2_ANSWER:
		if len(buf) < 6 {
			fmt.Println("Answer buffer too short, waiting for more", buf)
			return
		}

		// TODO: All of these buffer reads are ugly and terrible
		dlenBuf := bytes.NewBuffer([]byte{buf[4]})
		var dlen uint8
		binary.Read(dlenBuf, binary.BigEndian, &dlen)

		if len(buf) < int(dlen)+5 {
			fmt.Println("Buffer shorter than expected length")
			return
		}

		mrspBuf := bytes.NewBuffer([]byte{buf[2]})
		var mrsp uint8
		binary.Read(mrspBuf, binary.BigEndian, &mrsp)

		seqBuf := bytes.NewBuffer([]byte{buf[3]})
		var seq uint8
		binary.Read(seqBuf, binary.BigEndian, &seq)

		dataEnd := 5 + (dlen - 1)
		data := buf[5:dataEnd]

		// Calculate the chk
		chkSlice := buf[2:dataEnd]
		compChk := ComputeChk(chkSlice)

		chkBuf := bytes.NewBuffer([]byte{buf[dataEnd]})
		var chk uint8
		binary.Read(chkBuf, binary.BigEndian, &chk)

		if compChk != chk {
			err = fmt.Errorf("Invalid check, expected %#x, got %#x", chk, compChk, buf[2:dataEnd])
			return
		}

		r := &Response{
			sop1: sop1,
			sop2: sop2,
			mrsp: mrsp,
			seq:  seq,
			dlen: dlen,
			data: data,
			chk:  chk,
		}
		if res, ok := s.res[seq]; ok {
			res <- r
		}
		n = int(dataEnd) + 1
	case SOP2_ASYNC:
		if len(buf) < 7 {
			fmt.Println("Async buffer too short, waiting for more", buf)
		}
		n = 1
		/* ID_CODE = buffer[2]
		   dlen = buffer.readUInt16BE(3)
		   if buffer.length < dlen + 5 {
		     return
		   }
		   startOfData = 5
		   dataEnd = startOfData + (dlen - 1)
		   endOfPacket = endOfData + 1
		   DATA = buffer.slice(startOfData, endOfData)
		   // msg = SOP2: SOP2, ID_CODE: ID_CODE, DATA: DATA
		   r := &AsyncResponse{[]byte{sop1}, []byte{sop2}, []byte{id}, []byte{data}}
		   s.async <- r
		   n = int(dataEnd) + 1
		 } */
	default:
		err = fmt.Errorf("Unexpected SOP2, should be %#x or %#x but got %#x", SOP2_ANSWER, SOP2_ASYNC, sop2)
		n = 1 // Chomp 1 byte and maybe we'll recover
	}
	return
}

func (s *Sphero) listen() {
	var data []byte
	var buf []byte
	var err error
	var n int

	for {
		select {
		case <-s.kill:
			fmt.Println("Killing goroutine")
			return
		default:
			data = make([]byte, 256)
			if n, err = s.Read(data); err != nil {
				fmt.Println("Read:", err)
			}
			if n > 0 {
				data = data[:n]
				buf = append(buf, data...)
			}
			if len(buf) > 1 {
				if n, err = s.parse(buf); err != nil {
					fmt.Println("Parse:", err)
				}
				if n > 0 {
					buf = buf[n:]
				}
			}
		}
	}
}

// Implement io.ReadWriteCloser

// Implement io.Closer
func (s *Sphero) Close() error {
	s.kill <- true // Kill our goroutine
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

	chk := ComputeChk(buf.Bytes()[2:buf.Len()])
	binary.Write(&buf, binary.BigEndian, chk) // DLEN

	fmt.Printf("Writing: %#x\n", buf.Bytes())

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
	data.Write([]byte{0x00}) // User flag - this would set the "user LED color" if 0x01
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
