FROM golang:1.16-alpine

WORKDIR /app

COPY . .
WORKDIR ./cmd

RUN go mod download && go mod verify
RUN go build -o main .

CMD ["./main"]
