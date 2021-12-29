
FROM golang:1.16.12-alpine3.14 AS builder
ENV GO111MODULE=on 
WORKDIR /app
COPY . ./
RUN go mod download

## Our project will now successfully build with the necessary go libraries included.
RUN go build -o scaleflixapi .
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/scaleflixapi .
EXPOSE 3000
CMD ["/app/scaleflixapi"]