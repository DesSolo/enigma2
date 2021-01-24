FROM golang:1.15.7-alpine AS builder

RUN apk --update add make
WORKDIR /build
ADD . .
RUN make build-docker

FROM alpine

WORKDIR /enigma
COPY --from=builder /build/bin .
RUN mv *_docker enigma
CMD ["./enigma"]