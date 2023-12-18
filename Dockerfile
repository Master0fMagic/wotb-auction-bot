FROM golang:1.21-alpine

RUN apk add --no-cache git make gcc g++ musl-dev linux-headers gmp-dev mpfr-dev

WORKDIR /go/src/app

COPY . .

RUN make build

WORKDIR /bin

# Run the bot when the container starts
CMD ["./wotb-auction-bot"]
