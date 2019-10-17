FROM golang:1.13.1-buster as builder

RUN mkdir -p $GOPATH/src/github.com/sergiorua/kube-network-flow
ADD . $GOPATH/src/github.com/sergiorua/kube-network-flow
WORKDIR $GOPATH/src/github.com/sergiorua/kube-network-flow

RUN GO111MODULE=on go build -o kube-network-flow main.go

FROM debian:buster
COPY --from=builder /go/src/github.com/sergiorua/kube-network-flow/kube-network-flow /
WORKDIR /
CMD ["/kube-network-flow"]