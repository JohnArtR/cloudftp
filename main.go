package main

import "github.com/JohnArtR/cloudftp/server"
import "github.com/JohnArtR/cloudftp/client"
import "github.com/JohnArtR/cloudftp/paradise"
import "flag"
import "log"

var (
	gracefulChild = flag.Bool("graceful", false, "listen on fd open 3 (internal use only)")
	stressTest    = flag.Bool("stressTest", false, "start a client making connections")
)

func main() {
	log.Println(" [INFO] Starting BigFTP server")
	flag.Parse()
	if *stressTest {
		go client.StressTest()
	}
	go server.Monitor()
	fm := paradise.NewDefaultFileSystem()
	am := paradise.NewDefaultAuthSystem()
	server.Start(fm, am, *gracefulChild)
}
