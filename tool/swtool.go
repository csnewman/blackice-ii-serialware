package main

import (
	"github.com/schollz/progressbar/v3"
	"go.bug.st/serial"
	"log"
	"os"
	"time"
)

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

func (c *Conn) Close() {
	if err := c.p.Close(); err != nil {
		log.Panic(err)
	}
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

type Packet struct {
	T    byte
	Data []byte
}

func (c *Conn) ReadByte() byte {
	got := make([]byte, 1)

	_, err := c.p.Read(got)
	if err != nil {
		log.Fatal(err)
	}

	return got[0]
}

func (c *Conn) ReadPacket() Packet {
	for {
		m := c.ReadByte()
		if m != 0x5c {
			continue
		}

		t := c.ReadByte()

		if DEBUG {
			log.Println("Packet started t=", t)
		}

		l := c.ReadByte()
		data := make([]byte, l)
		_, err := c.p.Read(data)
		if err != nil {
			log.Fatal(err)
		}

		crc := c.ReadByte()
		calcCRC := calcCRC(data)

		if DEBUG {
			log.Println("crc=", crc, "calc=", calcCRC)
		}

		if crc != calcCRC {
			log.Fatal("crc mismatch")
		}

		return Packet{
			T:    t,
			Data: data,
		}
	}
}

func (c *Conn) WriteByte(b byte) {
	_, err := c.p.Write([]byte{b})
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Conn) WritePacket(p Packet) {
	for {
		if DEBUG {
			log.Println("Sending packet t=", p.T)
		}

		c.WriteByte(0x5c)
		c.WriteByte(p.T)
		c.WriteByte(byte(len(p.Data)))

		if len(p.Data) > 0 {
			_, err := c.p.Write(p.Data)
			if err != nil {
				log.Fatal(err)
			}
		}

		crc := calcCRC(p.Data)
		c.WriteByte(crc)

		if DEBUG {
			log.Println("awaiting conf conf")
		}

		resp := c.ReadByte()
		if resp != 0x5e {
			if DEBUG {
				log.Println("Write failure", resp)
			}

			continue
		}

		if DEBUG {
			log.Println("Write ok", resp)
		}

		return
	}
}

func main() {
	if len(os.Args) < 3 {
		log.Println("swtool upload <device> <file>")
		log.Println("swtool clear <device>")
		log.Println("swtool shell <device>")
		log.Println("swtool ping <device>")
		return
	}

	c := Open(os.Args[2])
	defer c.Close()

	switch os.Args[1] {
	case "ping":
		c.WritePacket(Packet{
			T:    5,
			Data: []byte{},
		})

		p := c.ReadPacket()

		if p.T != 5 || p.Data[0] != 123 {
			log.Println("Ping Fail")
		} else {
			log.Println("Ping OK")
		}

	case "clear":
		c.WritePacket(Packet{
			T:    10,
			Data: []byte{},
		})

		p := c.ReadPacket()

		if p.T != 10 || p.Data[0] != 0 {
			log.Println("Clear Fail")
		} else {
			log.Println("Clear OK")
		}

	case "upload":
		if len(os.Args) != 4 {
			log.Println("swtool upload <device> <file>")
			return
		}

		c.WritePacket(Packet{
			T:    10,
			Data: []byte{},
		})

		p := c.ReadPacket()

		if p.T != 10 || p.Data[0] != 0 {
			log.Fatal("Clear Fail")
		} else {
			log.Println("Clear OK")
		}

		dat, err := os.ReadFile(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}

		var bar *progressbar.ProgressBar
		if !DEBUG {
			bar = progressbar.DefaultBytes(int64(len(dat)), "uploading")
		}

		for i := 0; i < len(dat); i += 255 {
			if DEBUG {
				log.Printf("[%d/%d] %f", i, len(dat), float32(i)/float32(len(dat)))
			}

			end := i + 255
			if end > len(dat)-1 {
				end = len(dat) - 1
			}

			seg := dat[i:end]

			c.WritePacket(Packet{
				T:    11,
				Data: seg,
			})

			p := c.ReadPacket()

			if p.T != 11 || p.Data[0] != 0 {
				log.Fatal("Chunk Fail")
			}

			if !DEBUG {
				bar.Add(len(seg))
			}
		}

		if !DEBUG {
			bar.Clear()
		}

		c.WritePacket(Packet{
			T:    12,
			Data: []byte{},
		})

		p = c.ReadPacket()

		if p.T != 12 || p.Data[0] != 0 {
			log.Fatal("Flash Fail")
		} else {
			log.Println("Flash OK")
		}

	case "shell":
		for {
			log.Println("TODO")
			time.Sleep(time.Second)
		}

	default:
		log.Fatal("Unknown command")
	}
}
