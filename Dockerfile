FROM golang:1.20 as build

WORKDIR /build
COPY cmd /build/cmd
#COPY internal /build/internal
COPY pkg /build/pkg
COPY go.mod go.sum main.go /build/

RUN GOOS=linux go build -ldflags "-s -w" -o /build/book-it .

FROM debian:bookworm-slim

RUN apt update && apt install -y ca-certificates

WORKDIR /app
COPY --from=build /build/book-it /app/book-it

LABEL authors="bryce@thuilot.io"
LABEL repository="github.com/bthuilot/book-it"
ENTRYPOINT ["/app/book-it"]