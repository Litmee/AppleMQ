package apple

import (
	"AppleMQ/message"
	"AppleMQ/treaty"
	"encoding/json"
	"log"
)

// Synchronized cluster messages
func separate(m []byte) {

	var info message.AppleMessage
	err := json.Unmarshal(m, &info)
	if err != nil {
		log.Println("message json serialization failed, err: ", err)
		return
	}
	if info.Sign != 0 {
		return
	}

	info.Sign = 1
	marshal, _ := json.Marshal(info)

	b, _ := treaty.Encode(string(marshal))
	for _, v := range clusterMQArr {
		if v.c == nil {
			goto deal
		} else {
			_, err := v.c.Write(b)
			if err != nil {
				goto deal
			}
			continue
		}
	deal:
		v.state = 0
		// log.Println("Message synchronization failed")
		FailLock.Lock()
		s, ok := failureMessageCollection[v.addr]
		if ok {
			s = append(s, b)
			failureMessageCollection[v.addr] = s
			FailLock.Unlock()
			continue
		}
		s = make([][]byte, 1)
		s[0] = b
		failureMessageCollection[v.addr] = s
		FailLock.Unlock()
	}
}

func dealMessageStandalone(m []byte) {
	globalQueue.Push(m)
}

func dealMessageCluster(m []byte) {
	globalQueue.Push(m)
	// Distribute messages to other cluster nodes
	go separate(m)
}
