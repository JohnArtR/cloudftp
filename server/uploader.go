package server

import "io"
import (
	"time"
	"github.com/JohnArtR/cloudftp/db"
	"fmt"
	"log"
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
		p.writeMessage(226, fmt.Sprintf("OK, received %d bytes", len(passive.data)))
		err := p.saveToDisk(passive)
		if err != nil {
			log.Println(" [DEBUG] Error while saving file to disk")
			p.writeMessage(550, "Error while dump file to disk")
			p.closePassive(passive)
			return
		}
		fileInfo := db.File{Name:p.param, Size:len(passive.data), UploadTime:time.Now(), User:p.user.Username}
		err = db.FileRepo.SaveFile(fileInfo)
		if err != nil {
			p.writeMessage(550, "Something error")
			log.Printf(" [ERROR] Error while inserting FileInfo: %s", err)
		}
	} else {
		p.writeMessage(550, "Error with upload: "+err.Error())
		p.closePassive(passive)
		return
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
	total = int64(len(p.buffer))
	for {
		temp_buffer := make([]byte, 1048576) // reads 1MB at a time
		n, err = passive.connection.Read(temp_buffer)
		total += int64(n)

		if err != nil {
			break
		}
		passive.data = append(passive.data, temp_buffer[0:n]...)
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
