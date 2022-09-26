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

// this function send full buffer and returns nil, otherwise it
// returns error and the buffer may be not sent.
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

// this function read the required length in buffer "want" and return
// nil. Otherwise it returns error and the buffer content is not defined
func (b *bufferIO)Read(want []byte)(error) {
	var expected_length = len(want)
	var current_length = 0
	var err error
	var l int

	for {
		l, err = b.Reader.Read(want[current_length:])
		if l > 0 {
			current_length += l
			if current_length >= expected_length {
				return nil
			}
		}
		if err != nil {
			return err
		}
	}
}

// Decode response with point of view of the client. Note if
// the received packet is PROGRESS, the fonction silently eat
// data and wait for next packet.
func (b *bufferIO)ReceivePacket()([]byte, error) {
	var length uint
	var msg []byte
	var err error

	// Read data as long as we have full message
	for {

		// decode message length and check avalaible length. If no
		// sufficient data, try again read network
		msg = make([]byte, 4)
		err = b.Read(msg)
		if err != nil {
			return nil, err
		}
		length, err = DecodeLength(msg)
		if err != nil {
			return nil, err
		}

		// Now read the message which measure the expected length
		msg = make([]byte, int(length))
		err = b.Read(msg)
		if err != nil {
			return nil, err
		}

		// Special case, if the message is SMFIR_PROGRESS, ignore it
		// and read again. This message is designed to retrigger timeout
		// during long process
		if length == 1 && toMsgType(msg[0]) == SMFIR_PROGRESS {
			continue
		}

		return msg, nil
	}
}

