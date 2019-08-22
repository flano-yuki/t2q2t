package cmd

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net"
	"os"

	"github.com/flano-yuki/t2q2t/config"
	"github.com/flano-yuki/t2q2t/lib"
	quic "github.com/lucas-clemente/quic-go"
	"github.com/spf13/cobra"
)

var q2tCmd = &cobra.Command{
	Use:   "q2t",
	Short: "Listen by quic, and forward to tcp",
	Long: `Listen by quic, and forward to tcp
  t2q2t q2t <Listen Addr> <forward Addr>  

  go run ./t2q2t.go q2t 0.0.0.0:2022 127.0.0.1:22`,
	Args: cobra.MinimumNArgs(2),
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

	tlsConfig := config.GenerateServerTLSConfig()
	quicConfig := config.GenerateServerQUICConfig()
	listener, err := quic.ListenAddr(addr, tlsConfig, quicConfig)
	toTcpAddr, err := net.ResolveTCPAddr("tcp4", to)
	if err != nil {
		os.Exit(0)
	}
	for {
		sess, err := listener.Accept(context.Background())
		if err != nil {
			os.Exit(0)
		}
		go q2tHandleConn(sess, toTcpAddr)
	}

}

func q2tHandleConn(sess quic.Session, toTcpAddr *net.TCPAddr) error {
	for {
		stream, err := sess.AcceptStream(context.Background())
		if err != nil {
			return err
		}
		go q2tHandleStream(stream, toTcpAddr)
	}
}

func q2tHandleStream(stream quic.Stream, toTcpAddr *net.TCPAddr) error {
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
