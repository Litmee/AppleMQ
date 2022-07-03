package apple

import (
	"AppleMQ/message"
	"AppleMQ/treaty"
	"encoding/json"
	"log"
)

// Synchronized cluster messages
func separate(m []byte) {
	b, _ := treaty.Encode(string(m))
	for _, v := range clusterMQArr {
		_, err := v.c.Write(b)
		if err != nil {
			v.state = 0
			log.Println("Message synchronization failed")
			FailLock.Lock()
			s, ok := failureMessageCollection[v.addr]
			if ok {
				s = append(s, b)
				continue
			}
			s = make([][]byte, 1)
			s[0] = b
			failureMessageCollection[v.addr] = s
			FailLock.Unlock()
		}
	}
}

func dealMessage(m []byte) {
	globalQueue.Push(m)
	var info message.AppleMessage
	err := json.Unmarshal(m, &info)
	if err != nil {
		log.Println("message json serialization failed, err: ", err)
		return
	}
	if info.Sign == 0 {
		info.Sign = 1
		marshal, _ := json.Marshal(info)
		// Distribute messages to other cluster nodes
		go separate(marshal)
	}
}
