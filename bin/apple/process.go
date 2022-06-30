package apple

import (
	"AppleMQ/treaty"
	"bufio"
	"log"
	"net"
)

func process(c net.Conn) {

	// Close the connection after processing
	defer func(c net.Conn) {
		err := c.Close()
		if err != nil {
			return
		}
	}(c)

	for {
		s, err := treaty.Decode(bufio.NewReader(c))
		if err != nil {
			log.Printf("read from conn failed, err:%v\n", err)
			break
		}
		log.Println("接收到的数据:", s)
	}
}
