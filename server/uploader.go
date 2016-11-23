package server

import "io"
import (
	"time"
	"log"
	//"io/ioutil"
	//"os"
)

func (p *Paradise) saveToDisk(passive *Passive) error {
	err := FileService.Save(passive.data, p.param)
	//err := ioutil.WriteFile(Settings.StorageDirectory + "/" + p.param, passive.data, os.FileMode(0777))
	//f, err := os.Create(Settings.StorageDirectory + "/lol.jpg")
	if err != nil {
		return err
	}
	return nil
}

func (p *Paradise) HandleStore() {
	passive := p.lastPassive()
	if passive == nil {
		return
	}
	log.Printf(" [INFO] Receiving file %s", p.param)
	p.writeMessage(150, "Data transfer starting")
	if waitTimeout(&passive.waiter, time.Minute) {
		p.writeMessage(550, "Could not get passive connection.")
		p.closePassive(passive)
		return
	}
	if passive.listenFailedAt > 0 {
		p.writeMessage(550, "Could not get passive connection.")
		p.closePassive(passive)
		return
	}

	_, err := p.storeOrAppend(passive)
	if err == io.EOF {
		p.writeMessage(226, "OK, received some bytes") // TODO send total in message
		err := p.saveToDisk(passive)
		if err != nil {
			log.Println(" [DEBUG] Error while saving file to disk")
		}

	} else {
		p.writeMessage(550, "Error with upload: "+err.Error())
	}

	p.closePassive(passive)
}

func (p *Paradise) storeOrAppend(passive *Passive) (int64, error) {
	passive.data = []byte{}
	var err error
	err = p.readFirst512Bytes(passive)
	if err != nil {
		return 0, err
	}

	// TODO run p.buffer thru mime type checker
	// if mime type bad, reject upload

	passive.data = append(passive.data, p.buffer...)

	var total int64
	var n int
	var iter int = 0
	total = int64(len(p.buffer))
	for {
		iter++
		temp_buffer := make([]byte, 20971520) // reads 20MB at a time
		n, err = passive.connection.Read(temp_buffer)
		total += int64(n)

		if err != nil {
			break
		}
		passive.data = append(passive.data, temp_buffer[0:n]...)
		log.Printf(" [DEBUG] Data transfer control (%s) (%d)" +
			"\n\tbuf size: %d" +
			"\n\treaded: %d" +
			"\n\ttotal real: %d", p.param, iter, len(temp_buffer), n, len(passive.data))
		if err != nil {
			break
		}
	}
	log.Println(p.path, " Done ", total)

	return total, err
}

func (p *Paradise) readFirst512Bytes(passive *Passive) error {
	p.buffer = make([]byte, 0)
	var err error
	for {
		temp_buffer := make([]byte, 512)
		n, err := passive.connection.Read(temp_buffer)

		if err != nil {
			break
		}
		//log.Print(string(temp_buffer))
		p.buffer = append(p.buffer, temp_buffer[0:n]...)

		if len(p.buffer) >= 512 {
			break
		}
	}

	if err != nil && err != io.EOF {
		return err
	}

	// you have a buffer filled to 512, or less if file is less than 512
	return nil
}
