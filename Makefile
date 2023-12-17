GOLINT := golangci-lint

all: dep lint build

dep:
	go mod download

check-lint:
	@which $(GOLINT) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GO_PATH)/bin v1.54.2

lint: dep check-lint
	 $(GOLINT) run --timeout 1h -c .golangci.yml

build: dep
	CGO_ENABLED=1  go build -v -o build/bin/wotb-auction-bot  ./
	@echo "Done building."

build-docker:
	docker build -t wotb-auction-bot .

