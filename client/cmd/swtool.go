package main

import (
	"github.com/csnewman/blackice-ii-serialware/client"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Println("swtool upload <device> <file>")
		log.Println("swtool clear <device>")
		log.Println("swtool ping <device>")
		return
	}

	c := client.Open(os.Args[2])
	defer c.Close()

	switch os.Args[1] {
	case "ping":
		status, err := c.Ping()
		if err != nil {
			log.Fatal(err)
		}

		if status {
			log.Println("Ping OK")
		} else {
			log.Println("Ping Fail")
		}

	case "clear":
		status, err := c.Clear()
		if err != nil {
			log.Fatal(err)
		}

		if status {
			log.Println("Clear OK")
		} else {
			log.Println("Clear Fail")
		}

	case "upload":
		if len(os.Args) != 4 {
			log.Println("swtool upload <device> <file>")
			return
		}

		dat, err := os.ReadFile(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}

		var bar *progressbar.ProgressBar
		if !client.DEBUG {
			bar = progressbar.DefaultBytes(int64(len(dat)), "uploading")
		}

		err = c.Upload(dat, func(pos int) {
			if client.DEBUG {
				log.Printf("[%d/%d] %f", pos, len(dat), float32(pos)/float32(len(dat)*100))
			} else {
				bar.Set(pos)
			}
		})

		if !client.DEBUG {
			bar.Clear()
		}

		if err != nil {
			log.Fatal(err)
		}

		log.Println("Flash OK")
	default:
		log.Fatal("Unknown command")
	}
}
