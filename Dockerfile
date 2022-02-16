FROM golang:1.16-alpine AS builder

RUN apk --update add make
WORKDIR /build
ADD . .
RUN make build-docker && \
    cd bin && \
    mv *_docker enigma

FROM alpine:3.15.0

WORKDIR /enigma
COPY examples/config.yml /etc/enigma/config.yml
COPY --from=builder /build/bin .
CMD ["./enigma"]