package client

import (
	"errors"
	"go.bug.st/serial"
	"log"
)

var ErrCRCMismatch = errors.New("serialware: crc mismatch")
var ErrClearFail = errors.New("serialware: clear fail")
var ErrUploadFail = errors.New("serialware: upload fail")
var ErrFlashFail = errors.New("serialware: flash fail")
var ErrTimeout = errors.New("serialware: timeout")

const DEBUG = false

type Conn struct {
	p serial.Port
}

func Open(path string) *Conn {
	port, err := serial.Open(path, &serial.Mode{
		BaudRate:          115200,
		DataBits:          8,
		Parity:            serial.NoParity,
		StopBits:          serial.OneStopBit,
		InitialStatusBits: nil,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Conn{
		p: port,
	}
}

func (c *Conn) Close() error {
	return c.p.Close()
}

func calcCRC(data []byte) byte {
	var crc byte

	for i := 0; i < len(data); i++ {
		crc ^= data[i]
		for j := 0; j < 8; j++ {
			if crc&1 == 1 {
				crc ^= 0x91
			}
			crc >>= 1
		}
	}

	return crc
}

type packet struct {
	T    byte
	Data []byte
}

func (c *Conn) readByte() (byte, error) {
	got := make([]byte, 1)

	_, err := c.p.Read(got)

	return got[0], err
}

func (c *Conn) readPacket() (*packet, error) {
	for {
		m, err := c.readByte()
		if err != nil {
			return nil, err
		}

		if m != 0x5c {
			continue
		}

		t, err := c.readByte()
		if err != nil {
			return nil, err
		}

		if DEBUG {
			log.Println("packet started t=", t)
		}

		l, err := c.readByte()
		if err != nil {
			return nil, err
		}

		data := make([]byte, l)
		_, err = c.p.Read(data)
		if err != nil {
			return nil, err
		}

		crc, err := c.readByte()
		if err != nil {
			return nil, err
		}

		calcCRC := calcCRC(data)

		if DEBUG {
			log.Println("crc=", crc, "calc=", calcCRC)
		}

		if crc != calcCRC {
			return nil, ErrCRCMismatch
		}

		return &packet{
			T:    t,
			Data: data,
		}, nil
	}
}

func (c *Conn) writeByte(b byte) error {
	_, err := c.p.Write([]byte{b})

	return err
}

func (c *Conn) writePacket(p packet) error {
	for {
		if DEBUG {
			log.Println("Sending packet t=", p.T)
		}

		if err := c.writeByte(0x5c); err != nil {
			return err
		}

		if err := c.writeByte(p.T); err != nil {
			return err
		}

		if err := c.writeByte(byte(len(p.Data))); err != nil {
			return err
		}

		if len(p.Data) > 0 {
			_, err := c.p.Write(p.Data)
			if err != nil {
				return err
			}
		}

		crc := calcCRC(p.Data)

		if err := c.writeByte(crc); err != nil {
			return err
		}

		if DEBUG {
			log.Println("awaiting conf conf")
		}

		resp, err := c.readByte()
		if err != nil {
			return err
		}

		if resp != 0x5e {
			if DEBUG {
				log.Println("Write failure", resp)
			}

			continue
		}

		if DEBUG {
			log.Println("Write ok", resp)
		}

		return nil
	}
}

func (c *Conn) Ping() (bool, error) {
	err := c.writePacket(packet{
		T:    5,
		Data: []byte{},
	})
	if err != nil {
		return false, err
	}

	p, err := c.readPacket()
	if err != nil {
		return false, err
	}

	return p.T == 5 && p.Data[0] == 123, nil
}

func (c *Conn) Clear() (bool, error) {
	err := c.writePacket(packet{
		T:    10,
		Data: []byte{},
	})
	if err != nil {
		return false, err
	}

	p, err := c.readPacket()
	if err != nil {
		return false, err
	}

	return p.T == 10 && p.Data[0] == 0, nil
}

func (c *Conn) SendChunk(chunk []byte) (bool, error) {
	err := c.writePacket(packet{
		T:    11,
		Data: chunk,
	})
	if err != nil {
		return false, err
	}

	p, err := c.readPacket()
	if err != nil {
		return false, err
	}

	return p.T == 11 && p.Data[0] == 0, nil
}

func (c *Conn) Complete() (bool, error) {
	err := c.writePacket(packet{
		T:    12,
		Data: []byte{},
	})
	if err != nil {
		return false, err
	}

	p, err := c.readPacket()
	if err != nil {
		return false, err
	}

	return p.T == 12 && p.Data[0] == 0, nil
}

func (c *Conn) Upload(data []byte, progress func(int)) error {
	if progress != nil {
		progress(0)
	}

	status, err := c.Clear()
	if err != nil {
		return err
	}

	if !status {
		return ErrClearFail
	}

	for i := 0; i < len(data); i += 255 {
		end := i + 255
		if end > len(data) {
			end = len(data)
		}

		seg := data[i:end]

		status, err := c.SendChunk(seg)
		if err != nil {
			return err
		}

		if !status {
			return ErrUploadFail
		}

		if progress != nil {
			progress(end)
		}
	}

	status, err = c.Complete()
	if err != nil {
		return err
	}

	if !status {
		return ErrFlashFail
	}

	return nil
}

func (c *Conn) WriteUser(data []byte, progress func(int)) error {
	if progress != nil {
		progress(0)
	}

	for i := 0; i < len(data); i += 250 {
		end := i + 250
		if end > len(data) {
			end = len(data)
		}

		seg := data[i:end]

		err := c.writePacket(packet{
			T:    21,
			Data: seg,
		})
		if err != nil {
			return err
		}

		p, err := c.readPacket()
		if err != nil {
			return err
		}

		if p.T != 21 || p.Data[0] != 2 {
			return ErrUploadFail
		}

		if progress != nil {
			progress(end)
		}
	}

	return nil
}

func (c *Conn) ReadUser(readLen int, timeout byte, progress func(int)) ([]byte, error) {
	if progress != nil {
		progress(0)
	}

	data := make([]byte, readLen)

	for i := 0; i < len(data); i += 250 {
		end := i + 250
		if end > len(data) {
			end = len(data)
		}

		err := c.writePacket(packet{
			T: 20,
			Data: []byte{
				timeout,
				byte(end - i),
			},
		})
		if err != nil {
			return nil, err
		}

		p, err := c.readPacket()
		if err != nil {
			return nil, err
		}

		if p.T != 20 || p.Data[0] != 2 {
			return nil, ErrTimeout
		}

		for j := 0; j < end-i; j++ {
			data[i+j] = p.Data[1+j]
		}

		if progress != nil {
			progress(end)
		}
	}

	return data, nil
}
