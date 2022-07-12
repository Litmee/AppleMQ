package apple

import (
	"log"
	"net"
)

func Run() {

	// Read the apple.ini configuration file
	readConf()

	// Start the service
	listen, err := net.Listen("tcp", "127.0.0.1:"+options["port"])
	log.Println("start port: ", options["port"])
	if err != nil {
		panic(err)
	}

	// Check and process the value of a configuration item
	ok := checkAndDealOption()

	for {
		// listen for connections
		conn, err := listen.Accept()
		if err != nil {
			log.Println("accept failed, err:%V\n", err)
			continue
		}
		// Start a separate goroutine to handle the connection
		if ok {
			go processCluster(conn)
		} else {
			go processStandalone(conn)
		}
	}
}
