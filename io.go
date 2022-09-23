// Copyright (c) 2022 Thierry FOURNIER (tfournier@arpalert.org)

package milter

import "bufio"
import "net"

type bufferIO struct {
	Conn net.Conn
	Reader *bufio.Reader
}

func (b *bufferIO)InitBufferIO(conn net.Conn) {
	b.Conn = conn
	b.Reader = bufio.NewReader(conn)
}

func (b *bufferIO)Close()(error) {
	return b.Conn.Close()
}

func (b *bufferIO)Write(data []byte)(error) {
	var err error
	var length int

	for {
		length, err = b.Conn.Write(data)
		if err != nil {
			return err
		}
		data = data[length:]
		if len(data) > 0 {
			continue
		}
		return nil
	}
}

func (b *bufferIO)Read(want []byte)(int, error) {
	var expected_length = len(want)
	var current_length = 0
	var err error
	var l int

	for {
		l, err = b.Reader.Read(want[current_length:])
		if l > 0 {
			current_length += l
			if current_length >= expected_length {
				return len(want), nil
			}
		}
		if err != nil {
			return 0, err
		}
	}
}

