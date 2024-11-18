FROM golang:1.23-alpine AS GO_BUILD
ADD . /work
WORKDIR /work
RUN go build -o drawio-go main.go

FROM alpine AS RUNTIME

WORKDIR /app
EXPOSE 8080

VOLUME /webroot
ENTRYPOINT /app/entrypoint.sh
HEALTHCHECK CMD wget -q -O /dev/null http://localhost:8080/health || exit 1

RUN apk add unzip wget

COPY --from=GO_BUILD /work/drawio-go /app/
ADD entrypoint.sh /app/
