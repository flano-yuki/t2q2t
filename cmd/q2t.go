package cmd

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"golang.org/x/sync/errgroup"
	"math/big"
	"net"
	"os"

	"github.com/flano-yuki/t2q2t/lib"
	quic "github.com/lucas-clemente/quic-go"
	"github.com/spf13/cobra"
)

var q2tCmd = &cobra.Command{
	Use:   "q2t",
	Short: "Listen by quic, and forward to tcp",
	Long:  `Listen by quic, and forward to tcp
  t2q2t q2t <Listen Addr> <forward Addr>  

  go run ./t2q2t.go q2t 0.0.0.0:2022 127.0.0.1:22`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("q2t called")
		listen := args[0]
		to := args[1]

		runq2t(listen, to)
	},
}

func init() {
	rootCmd.AddCommand(q2tCmd)

}

func runq2t(listen, to string) {
	addr := listen
	fmt.Printf("Listen QUIC on: %s \n", addr)
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	toTcpAddr, err := net.ResolveTCPAddr("tcp4", to)
	if err != nil {
		os.Exit(0)
	}
	for {
		sess, err := listener.Accept(context.Background())
		if err != nil {
			os.Exit(0)
		}
		stream, err := sess.AcceptStream(context.Background())
		if err != nil {
			os.Exit(0)
		}
		go q2tHandleConn(stream, toTcpAddr)
	}

}

func q2tHandleConn(stream quic.Stream, toTcpAddr *net.TCPAddr) error {
	fmt.Printf("Connect TCP to: %s \n", toTcpAddr.String())
	conn, err := net.DialTCP("tcp", nil, toTcpAddr)
	if err != nil {
		return err
	}

	eg := errgroup.Group{}
	eg.Go(func() error { return util.T2qRelay(conn, stream) })
	eg.Go(func() error { return util.Q2tRelay(stream, conn) })

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func generateTLSConfig() *tls.Config {
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
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"t2q2t"},
	}
}
