FROM golang:latest as builder
WORKDIR /go/src/enigma
ADD . /go/src/enigma
RUN make build-docker

FROM scratch
WORKDIR /enigma
COPY --from=builder /go/src/enigma/bin .
CMD ["./enigma"]