package apple

import (
	"AppleMQ/treaty"
	"log"
	"net"
	"sync"
	"time"
)

// Cluster address set
var clusterArr []string

// A collection of cluster connection structures returned by a successful connection
var clusterMQArr []*clusterMQ

// The set of surviving cluster connections
var clusterMQAliveArr []*clusterMQ

var FailLock sync.Mutex

// Synchronization failure message collection
var failureMessageCollection = make(map[string][][]byte)

type clusterMQ struct {
	id       int32
	location int32
	addr     string
	c        net.Conn
	trying   bool
	state    int8
	// state: 1 means the connection is alive, 0 means the connection is dead
}

func connection(i int, addr string) {
	for {
		d, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			log.Println("ResolveTCPAddr err: ", err)
			break
		}

		conn, err := net.DialTCP("tcp", nil, d)
		if err != nil {
			time.Sleep(time.Second * 1)
			continue
		}
		clusterMQArr[i].c = conn
		clusterMQArr[i].state = 1
		clusterMQArr[i].trying = false
		s, _ := treaty.Encode("send")
		_, err = conn.Write(s)
		if err != nil {
			break
		}
		// Sync messages
		FailLock.Lock()
		// Handle possible previously unsent messages
		f, ok := failureMessageCollection[addr]
		if ok {
			var newF [][]byte
			for j, v := range f {
				_, err = conn.Write(v)
				if err != nil {
					break
				}
				newF = f[j+1:]
			}
			if len(newF) > 0 {
				failureMessageCollection[addr] = newF
			} else {
				delete(failureMessageCollection, addr)
			}
		}
		FailLock.Unlock()
		break
	}
}

// Monitor the status of cluster connections
func monitor() {
	for {
		var noLive []int32
		for _, v := range clusterMQArr {
			// The connection is alive but not added to the alive set
			if v.state == 1 && v.location == -1 {
				clusterMQAliveArr = append(clusterMQAliveArr, v)
				v.location = int32(len(clusterMQAliveArr) - 1)
				continue
			}
			if v.state == 0 && v.location != -1 {
				noLive = append(noLive, v.location)
				v.location = -1
				v.trying = true
				go connection(int(v.id), v.addr)
				continue
			}
			if v.state == 0 && v.location == -1 && v.trying == false {
				v.trying = true
				go connection(int(v.id), v.addr)
			}
		}
		if len(noLive) > 0 {
			var newLive []*clusterMQ
			for _, v := range clusterMQAliveArr {
				if v.state != 0 {
					newLive = append(newLive, v)
				}
			}
			clusterMQAliveArr = newLive
		}
		time.Sleep(time.Second * 2)
	}
}

// Connect the machines participating in the cluster
func clusterConnection() {
	clusterMQArr = make([]*clusterMQ, len(clusterArr))
	for i, v := range clusterArr {
		clusterMQArr[i] = &clusterMQ{id: int32(i), location: -1, addr: v, c: nil, trying: true, state: 0}
		go connection(i, v)
	}
	// Start watch guard
	go monitor()
}
