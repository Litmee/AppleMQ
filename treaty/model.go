package treaty

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

type message struct {
	length int32
	body   []byte
}

// Encode 消息编码
func Encode(m string) ([]byte, error) {

	//  1. 读取消息的长度, 转换成 int32 类型(占 4 个字节)
	l := int32(len(m))

	// 2. 定义一个空 bytes 缓冲区
	b := new(bytes.Buffer)

	// 3. 写入消息头, 通过小端序列的方式把 l 写入 b
	err := binary.Write(b, binary.LittleEndian, l)
	if err != nil {
		return nil, err
	}

	// 4. 写入消息实体
	err = binary.Write(b, binary.LittleEndian, []byte(m))
	if err != nil {
		return nil, err
	}

	// 5. 返回封装好的消息
	return b.Bytes(), nil
}

// Decode 消息解码
func Decode(r *bufio.Reader) (string, error) {

	// 1. 读取前 4 个字节的数据, 及获取消息的内容长度, Peek 方式的读取是不会清掉缓存的
	lByte, _ := r.Peek(4)

	// 2. 定义一个以 lByte 位内容的 bytes 缓冲区
	buffer := bytes.NewBuffer(lByte)

	var l int32

	// 3. 将 buffer 的内容读取到 l 变量
	err := binary.Read(buffer, binary.LittleEndian, &l)
	if err != nil {
		return "", err
	}

	// 4. 通过 Buffered 方法返回缓冲中现有的可读取的字节数, 前面使用 Peek 读取, 所以这里的数据内容应该大于 l+4
	if int32(r.Buffered()) < l+4 {
		return "", err
	}

	// 5. 读取消息的实体
	s := make([]byte, int(l+4))
	_, err = r.Read(s)
	if err != nil {
		return "", err
	}

	// 6. 返回去掉长度标识的消息字符串
	return string(s[4:]), nil
}
