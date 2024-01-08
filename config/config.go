package config

import (
	"crypto/tls"
	quic "github.com/quic-go/quic-go"
	"time"
)

const ALPN = "t2q2t"

func GenerateClientQUICConfig() *quic.Config {
	return &quic.Config{
		MaxIdleTimeout:       time.Duration(1) * time.Hour,
		HandshakeIdleTimeout: time.Duration(5) * time.Second,
	}
}

func GenerateServerQUICConfig() *quic.Config {
	return &quic.Config{
		MaxIdleTimeout:       time.Duration(1) * time.Hour,
		MaxIncomingStreams:   1024,
	}
}

func GenerateClientTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{ALPN},
	}
}

func GenerateServerTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	cert, _ := tls.LoadX509KeyPair(certFile, keyFile)
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{ALPN},
	}
}
