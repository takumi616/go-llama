#Creates binary which is included in deploy container
FROM golang:1.22.3-alpine3.19 as builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -trimpath -ldflags "-w -s" -o main

#Deploy container
FROM alpine:latest as deploy
RUN apk update
COPY --from=builder /app/main .
CMD ["./main"]