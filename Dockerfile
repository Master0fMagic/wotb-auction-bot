FROM golang:1.21-alpine

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN go build -o /bin/wotb-auction-bot .

WORKDIR /bin

# Run the bot when the container starts
CMD ["./wotb-auction-bot"]
