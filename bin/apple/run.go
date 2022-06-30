package apple

import (
	"log"
	"net"
)

func Run() {

	// Read the apple.ini configuration file
	readConf()

	// Check and process the value of a configuration item
	checkAndDealOption()

	dealModeCluster()

	// 1. Start the service
	listen, err := net.Listen("tcp", "127.0.0.1:9082")
	if err != nil {
		panic(err)
	}

	// 2. Polling waits for the client to establish a connection
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println("accept failed, err:%V\n", err)
			continue
		}
		// 3. Start a separate goroutine to handle the connection
		go process(conn)
	}
}
