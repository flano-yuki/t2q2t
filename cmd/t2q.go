package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/flano-yuki/t2q2t/lib"
	quic "github.com/lucas-clemente/quic-go"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"net"
	"os"
)

// t2qCmd represents the t2q command
var t2qCmd = &cobra.Command{
	Use:   "t2q",
	Short: "Listen by tcp, and forward to quic",
	Long:  `Listen by tcp, and forward to quic
  t2q2t t2q <Listen Addr> <forward Addr>  

  go run ./t2q2t.go t2q 0.0.0.0:2022 127.0.0.1:2022:`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("t2q called")
		listen := args[0]
		to := args[1]

		runt2q(listen, to)
	},
}

func init() {
	rootCmd.AddCommand(t2qCmd)
}

func runt2q(listen, to string) {
	listenTcpAddr, err := net.ResolveTCPAddr("tcp4", listen)
	toTcpAddr, err := net.ResolveTCPAddr("tcp4", to)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listen TCP on: %s \n", listenTcpAddr.String())

	lt, err := net.ListenTCP("tcp", listenTcpAddr)
	for {
		conn, err := lt.AcceptTCP()
		if err != nil {
			os.Exit(0)
		}
		go t2qHandleConn(conn, toTcpAddr)
	}
}

func t2qHandleConn(conn *net.TCPConn, toTcpAddr *net.TCPAddr) error {
	fmt.Printf("Connect QUIC to: %s \n", toTcpAddr.String())

	var stream quic.Stream

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"t2q2t"},
	}
	sess, err := quic.DialAddr(toTcpAddr.String(), tlsConf, nil)
	if err != nil {
		return err
	}

	stream, err = sess.OpenStreamSync(context.Background())

	if err != nil {
		panic(err)
	}

	eg := errgroup.Group{}
	eg.Go(func() error { return util.T2qRelay(conn, stream) })
	eg.Go(func() error { return util.Q2tRelay(stream, conn) })

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
