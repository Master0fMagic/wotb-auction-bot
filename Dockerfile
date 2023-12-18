FROM golang:1.21-alpine

RUN apk add --no-cache make

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN make build

WORKDIR /bin

# Run the bot when the container starts
CMD ["./wotb-auction-bot"]
