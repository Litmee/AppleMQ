package message

type AppleMessage struct {
	// Sign: 0 means from the producer, 1 means from the cluster machine synchronization
	Sign int
	Body string
	Tag  string
}
