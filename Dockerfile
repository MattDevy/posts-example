FROM golang:1.16-alpine
RUN apk add --no-progress build-base

COPY . /app

WORKDIR /app

RUN go build -o server ./cmd/server

CMD ["./server"]