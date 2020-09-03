FROM golang:latest as builder
WORKDIR /go/src/enigma
RUN go get github.com/go-redis/redis
ADD . /go/src/enigma
RUN make build-docker

FROM scratch
WORKDIR /enigma
COPY --from=builder /go/src/enigma/bin .
CMD ["./enigma"]