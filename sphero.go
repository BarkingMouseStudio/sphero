package sphero

import (
	"bytes"
	"encoding/binary"
	"fmt"
	serial "github.com/tarm/goserial"
	"io"
	"os"
	"syscall"
)

type Config struct {
	Bluetooth serial.Config
}

type Sphero struct {
	// Current Sphero configuration
	conf *Config

	// Connection to the underlying Bluetooth serial port
	conn io.ReadWriteCloser

	// Sequence number to associate with responses
	seq uint8

	// Map of response channels to sequence numbers
	res map[uint8]chan<- interface{}

	// Async response channel
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
		res:   make(map[uint8]chan<- interface{}),
		async: async,
	}

	go s.listen()

	return s, nil
}

type Response struct {
	// Start of Packet #1 - Always FFh
	sop1,

	// Start of Packet #2 - Set to FFh when this is an acknowledgement,
	// FEh when this is an asynchronous message
	sop2,

	// Message Response - This is generated by the message decoder of the
	// virtual device (refer to the appropriate appendix for a list of values)
	mrsp,

	// Sequence Number - Echoed to the client when this is a direct message
	// response (set to 00h when SOP2 = FEh)
	seq,

	// Data Length - The number of bytes following through the end of the packet
	dlen,

	// Data - Optional data in response to the Command or based on "streaming"
	// data settings
	data,

	// Checksum - Packet checksum (as computed above)
	chk []byte
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
		sum := 0
		for b := range chkSlice {
			sum += int(b)
		}
		compChk := (sum % 256) ^ 0xff

		chkBuf := bytes.NewBuffer([]byte{buf[dataEnd]})
		var chk uint8
		binary.Read(chkBuf, binary.BigEndian, &chk)

		if compChk != int(chk) {
			err = fmt.Errorf("Invalid check, expected %v, got %v", chk, compChk)
			return
		}

		r := &Response{[]byte{sop1}, []byte{sop2}, []byte{mrsp}, []byte{seq}, []byte{dlen}, data, []byte{chk}}
		if res, ok := s.res[seq]; ok {
			res <- r
		}
		return int(dataEnd) + 1, nil
	case SOP2_ASYNC:
		if len(buf) < 7 {
			fmt.Println("Async buffer too short, waiting for more", buf)
		}
	default:
		err = fmt.Errorf("Unexpected SOP2, should be %v or %v but got %v", SOP2_ANSWER, SOP2_ASYNC, sop2)
	}
	return
}

func (s *Sphero) listen() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	var data []byte
	var buf []byte
	var err error
	var n int

	for {
		// Read data from the Sphero
		data = make([]byte, 256)
		if n, err = s.Read(data); err != nil {
			fmt.Println("Failed reading, panicking")
			panic(err)
		}

		// We didn't receive any data
		if n == 0 {
			fmt.Println("Read no data")
			continue
		}

		data = data[:n]
		fmt.Println("Read data", data)

		// Append the new data to our buffer
		buf = append(buf, data...)

		fmt.Println("Buffer added", buf)

		// Attempt to parse the buf
		if n, err = s.parse(buf); err != nil {
			panic(err)
		}

		// We successfully parsed data, trim our buffer
		if n > 0 {
			buf = buf[n:]
		}
	}
}

// Implement io.ReadWriteCloser

// Implement io.Closer
func (s *Sphero) Close() error {
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

func (s *Sphero) Send() error {
	return nil
}

func (s *Sphero) Ping(res chan<- interface{}) error {
	s.seq++
	s.res[s.seq] = res

	var data []byte
	var buf bytes.Buffer
	buf.Write([]byte{SOP1})                                  // SOP1
	buf.Write([]byte{SOP2_ANSWER})                           // SOP2
	buf.Write([]byte{DID_CORE})                              // DID
	buf.Write([]byte{CMD_PING})                              // CID
	binary.Write(&buf, binary.BigEndian, uint8(0x52))        // SEQ
	binary.Write(&buf, binary.BigEndian, uint8(len(data)+1)) // DLEN

	// Calculate the chk
	chkBytes := buf.Bytes()[2:buf.Len()]
	sum := 0
	for _, b := range chkBytes {
		sum += int(uint8(b))
	}
	chk := (sum % 256) ^ 0xff
	binary.Write(&buf, binary.BigEndian, uint8(chk)) // DLEN

	_, err := s.Write(buf.Bytes())
	return err
}
