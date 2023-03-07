FROM golang:1.18 as builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY main.go ./
RUN go build -o kube-network-flow main.go

FROM debian:buster
COPY --from=builder /app/kube-network-flow /
WORKDIR /
CMD ["/kube-network-flow"]