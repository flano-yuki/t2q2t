FROM golang:1.13.0-alpine3.10 as builder
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN apk --no-cache add git
WORKDIR /go/src/github.com/flano-yuki/t2q2t
COPY go.mod go.sum ./
RUN GO111MODULE=on go mod download
COPY . .
RUN go build t2q2t.go

FROM alpine:3.10 as executor
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /go/src/github.com/flano-yuki/t2q2t/t2q2t /t2q2t
RUN addgroup go \
  && adduser -D -G go go \
  && chown -R go:go /t2q2t
EXPOSE 2022
EXPOSE 2022/udp
ENTRYPOINT ["/t2q2t"]
