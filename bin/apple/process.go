package apple

import (
	"AppleMQ/queue"
	"AppleMQ/treaty"
	"bufio"
	"log"
	"net"
)

var globalQueue = queue.NewQueue()

func process(c net.Conn) {

	// Close the connection after processing
	defer func(c net.Conn) {
		err := c.Close()
		if err != nil {
			return
		}
	}(c)
	// Read the first launch identification information
	s, err := treaty.Decode(bufio.NewReader(c))
	if err != nil {
		log.Printf("read from conn failed, err:%v\n", err)
		return
	}
	var sign bool
	if string(s) == "send" {
		sign = true
	}
	for sign {
		s, err := treaty.Decode(bufio.NewReader(c))
		if err != nil {
			log.Printf("read from conn failed, err:%v\n", err)
			break
		}
		go dealMessage(s)
	}
	for !sign {
		m := globalQueue.Take()
		m, _ = treaty.Encode(string(m))
		_, err := c.Write(m)
		if err != nil {
			break
		}
		globalQueue.DeleteHead()
	}
}
