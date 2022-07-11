package apple

import (
	"AppleMQ/queue"
	"AppleMQ/treaty"
	"bufio"
	"log"
	"net"
	"time"
)

// Global queue parameters
var globalQueue = queue.NewQueue()

func processStandalone(c net.Conn) {

	// Close the connection after processing
	defer func(c net.Conn) {
		err := c.Close()
		if err != nil {
			return
		}
	}(c)
	reader := bufio.NewReader(c)
	// Read the first launch identification information
	s, err := treaty.Decode(reader)
	if err != nil {
		log.Printf("read from conn failed, err:%v\n", err)
		return
	}
	var sign bool
	if string(s) == "send" {
		sign = true
	}
	go func() {
		time.Sleep(time.Second * 30)
		log.Println(globalQueue.Size())
		log.Println(len(failureMessageCollection["127.0.0.1:9083"]))
	}()
	for sign {
		s, err = treaty.Decode(reader)
		if err != nil {
			log.Printf("read from conn failed, err:%v\n", err)
			break
		}
		go dealMessageStandalone(s)
	}
	for !sign {
		m := globalQueue.Take()
		m, _ = treaty.Encode(string(m))
		_, err := c.Write(m)
		if err != nil {
			break
		}
	}
}

func processCluster(c net.Conn) {

	// Close the connection after processing
	defer func(c net.Conn) {
		err := c.Close()
		if err != nil {
			return
		}
	}(c)
	reader := bufio.NewReader(c)
	// Read the first launch identification information
	s, err := treaty.Decode(reader)
	if err != nil {
		log.Printf("read from conn failed, err:%v\n", err)
		return
	}
	var sign bool
	// Determine the first frame identification data
	if string(s) == "send" {
		sign = true
	}
	for sign {
		s, err = treaty.Decode(reader)
		if err != nil {
			log.Printf("read from conn failed, err:%v\n", err)
			log.Println(globalQueue.Size())
			log.Println(len(failureMessageCollection["127.0.0.1:9083"]))
			break
		}
		log.Println("MQ received the news: ", string(s))
		go dealMessageCluster(s)
	}
	for !sign {
		m := globalQueue.Take()
		m, _ = treaty.Encode(string(m))
		_, err := c.Write(m)
		if err != nil {
			break
		}
	}
}
