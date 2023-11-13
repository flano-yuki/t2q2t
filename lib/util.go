package util

import (
	"fmt"
	quic "github.com/quic-go/quic-go"
	"net"
)

const BUFFER_SIZE = 0xFFFF

func T2qRelay(fromConn *net.TCPConn, toStream quic.Stream) error {
	buff := make([]byte, BUFFER_SIZE)
	for {
		n, err := fromConn.Read(buff)
		if err != nil {
			return err
		}
		b := buff[:n]
		output("t2q(receive): " + string(b))
		_, err = toStream.Write(b)
		if err != nil {
			return err
		}
		output("t2q(send): " + string(b))
	}

	return nil
}

func Q2tRelay(fromStream quic.Stream, toConn *net.TCPConn) error {
	buff := make([]byte, BUFFER_SIZE)
	for {
		n, err := fromStream.Read(buff)
		if err != nil {
			return err
		}
		b := buff[:n]
		output("q2t(receive): " + string(b))
		_, err = toConn.Write(b)
		if err != nil {
			return err
		}
		output("t2q(send): " + string(b))
	}

	return nil
}

func output(str string) {
	//TODO
	debug := false
	if debug {
		fmt.Printf(str + "\n")
	}
}
