# t2q2t

t2q2t is a port forward tool that converts between TCP and QUIC. (TCP/QUIC Proxy).

You can benefit from QUIC with any application protocol using t2q2t.

usecase overview are below:

![t2q2t usecase](/img/overview.png)

- Client: listens on `0.0.0.0:2022` with TCP, and forwards to `192.168.0.1:22`(server ip) with QUIC
- Server: listen on `0.0.0.0:2022` with QUIC, and forward to `127.0.0.1:22` with TCP

```
# t2q2t <convert mode> <listen addr> <connect addr>

# on client side (TCP <-> QUIC)
$ ./t2q2t t2q 0.0.0.0:2022 192.168.0.1:2022

# on server side (QUIC <-> TCO)
$ ./t2q2t q2t 0.0.0.0:2022 127.0.0.1:22
```

You can use SSH with QUIC transport by running t2q2t on client and server.

## build

```
$ git clone https://github.com/flano-yuki/t2q2t.git
$ cd ./t2q2t
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build t2q2t.go
$ go build ./t2q2t.go
``` 

You can also use the following Dockerfile to try it out in your Docker environmemnt.

```
# In the project root directory
$ docker build -t flano-yuki/t2q2t .
$ docker run -it -p 2022:2022 -p 2022:2022/udp --rm flano-yuki/t2q2t
tcp/quic port forward tool
  t2q2t <subcommand> <Listen Addr> <forward Addr>

  go run ./t2q2t.go t2q 0.0.0.0:2022 127.0.0.1:2022
  go run ./t2q2t.go q2t 0.0.0.0:2022 127.0.0.1:22

Usage:
  t2q2t [command]

Available Commands:
  help        Help about any command
  q2t         Listen by quic, and forward to tcp
  t2q         Listen by tcp, and forward to quic
  version     Print the version number of t2q2t

Flags:
  -h, --help      help for t2q2t
  -v, --verbose   verbose output

Use "t2q2t [command] --help" for more information about a command.
```


## convert mode (sub command)

- `t2q`: Listen by TCP, and connect to QUIC
- `q2t`: Listen by QUIC, and connect to TCP


## Note
This is a PoC.

t2q2t uses "t2q2t" as the ALPN identifier. It is not intended to communicate with other QUIC implementations.

TODO
- Improve error handling
- Improve connection management
- (Output connection statistics)
- Many other improvements

