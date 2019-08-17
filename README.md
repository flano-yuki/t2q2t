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
$ cd ./t2q2t.git
$ go build ./t2q2t.go

``` 

## convert mode (sub command)

- `t2q`: Listen by TCP, and connect to QUIC
- `q2t`: Listen by QUIC, and connect to TCP


## Note
This is a PoC.

t2q2t uses "t2q2t" as the ALPN identifier. It is not intended to communicate with other QUIC implementations.

TODO
- Versioning
- Use QUIC stream multiplex
- Improve error handling
- Many other improvements

