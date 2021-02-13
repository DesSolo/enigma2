FROM golang:1.15.7-alpine AS builder

RUN apk --update add make
WORKDIR /build
ADD . .
RUN make build-docker && \
    cd bin && \
    mv *_docker enigma

FROM alpine

WORKDIR /enigma
COPY --from=builder /build/bin .
CMD ["./enigma"]