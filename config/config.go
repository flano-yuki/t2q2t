package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	quic "github.com/quic-go/quic-go"
	"math/big"
	"time"
)

const ALPN = "t2q2t"

func GenerateClientQUICConfig() *quic.Config {
	return &quic.Config{
		MaxIdleTimeout: time.Duration(1) * time.Hour,
		//KeepAlivePeriod:      true,
		HandshakeIdleTimeout: time.Duration(5) * time.Second,
	}
}

func GenerateServerQUICConfig() *quic.Config {
	return &quic.Config{
		MaxIdleTimeout: time.Duration(1) * time.Hour,
		//KeepAlivePeriod:    true,
		MaxIncomingStreams: 1024,
	}
}

func GenerateClientTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{ALPN},
	}
}

func GenerateServerTLSConfig() *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{generateServerCertificate()},
		NextProtos:   []string{ALPN},
	}
}

func generateServerCertificate() tls.Certificate {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return tlsCert
}
