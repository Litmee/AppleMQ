package treaty

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

// Encode message encoding
func Encode(m string) ([]byte, error) {

	//  1. Read the length of the message, convert it to int32 type (4 bytes)
	l := int32(len(m))

	// 2. define an empty bytes buffer
	b := new(bytes.Buffer)

	// 3. Write the message header, and write l to b in a little-endian sequence
	err := binary.Write(b, binary.LittleEndian, l)
	if err != nil {
		return nil, err
	}

	// 4. write message entity
	err = binary.Write(b, binary.LittleEndian, []byte(m))
	if err != nil {
		return nil, err
	}

	// 5. Return the packaged message
	return b.Bytes(), nil
}

// Decode message decoding
func Decode(r *bufio.Reader) ([]byte, error) {

	// Identification parameters used to solve TCP packetization problems
	var sign int32 = 0

	// 1. Read the first 4 bytes of data, and get the content length of the message, reading in Peek mode will not clear the cache
	lByte, _ := r.Peek(4)

	// 2. Defines a bytes buffer with lByte bits of content
	buffer := bytes.NewBuffer(lByte)

	var l int32

	// 3. Read the contents of buffer into the l variable
	err := binary.Read(buffer, binary.LittleEndian, &l)
	if err != nil {
		return nil, err
	}

	s := make([]byte, l+4)

	// 4. Return the current number of readable bytes in the buffer through the Buffered method,
	// which was previously read using Peek, so the data content here should be greater than l+4
	if int32(r.Buffered()) < l+4 {
		// return nil, err
		sign = int32(r.Buffered())
		s = make([]byte, sign)
	}

	// 5. read message entity
	_, err = r.Read(s)
	if err != nil {
		return nil, err
	}

	// sign != 0 indicates that there is subcontracting and needs to be spliced
	if sign != 0 {
		_, _ = r.Peek(1)
		newS := make([]byte, (l+4)-sign)
		_, err = r.Read(newS)
		if err != nil {
			return nil, err
		}
		s = append(s, newS...)
	}

	// 6. Returns the message string with the length flag removed
	return s[4:], nil
}
