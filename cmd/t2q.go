package cmd

import (
	"context"
	"fmt"
	"github.com/oniyan/t2q2t/config"
	"github.com/oniyan/t2q2t/lib"
	quic "github.com/quic-go/quic-go"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"net"
)

var t2qCmd = &cobra.Command{
	Use:   "t2q",
	Short: "Listen by tcp, and forward to quic",
	Long: `Listen by tcp, and forward to quic
  t2q2t t2q <Listen Addr> <forward Addr>  

  go run ./t2q2t.go t2q 0.0.0.0:2022 127.0.0.1:2022:`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("t2q called")
		listen := args[0]
		to := args[1]

		err := runt2q(listen, to)
		if err != nil {
			fmt.Printf("[Error] %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(t2qCmd)
}

func runt2q(listen, to string) error {
	listenTcpAddr, err := net.ResolveTCPAddr("tcp4", listen)
	toTcpAddr, err := net.ResolveTCPAddr("tcp4", to)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listen TCP on: %s \n", listenTcpAddr.String())

	lt, err := net.ListenTCP("tcp", listenTcpAddr)

	tlsConf := config.GenerateClientTLSConfig()
	quicConf := config.GenerateClientQUICConfig()
	var sess quic.Connection = nil
	for {
		conn, err := lt.AcceptTCP()
		if err != nil {
			return err
		}

		// TODO
		// and, if connection has closed
		if sess == nil {
			sess, err = quic.DialAddr(context.Background(), toTcpAddr.String(), tlsConf, quicConf)
			if err != nil {
				return err
			}
			fmt.Printf("Connect QUIC to: %s \n", toTcpAddr.String())
		}

		// TODO error handling
		go t2qHandleConn(conn, sess)
	}
	return nil
}

func t2qHandleConn(conn *net.TCPConn, sess quic.Connection) error {
	var stream quic.Stream
	stream, err := sess.OpenStreamSync(context.Background())
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
