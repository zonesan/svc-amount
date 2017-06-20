GOFILES:=$(shell find . -name '*.go' | grep -v -E '(./vendor)')

all: \
	bin/linux/svc-amount-agent \
	bin/darwin/svc-amount-agent

images: GVERSION=$(shell $(CURDIR)/git-version.sh)
images: bin/linux/svc-amount-agent
	docker build --no-cache -f Dockerfile.release -t svc-amount-agent:$(GVERSION) .
	@docker tag  svc-amount-agent:$(GVERSION) svc-amount-agent:latest

check:
	@find . -name vendor -prune -o -name '*.go' -exec gofmt -s -d {} +
	@go vet $(shell go list ./... | grep -v '/vendor/')
	@go test -v $(shell go list ./... | grep -v '/vendor/')

vendor:
	dep ensure

clean:
	rm -rf bin

bin/%: LDFLAGS=-X main.Version=$(shell $(CURDIR)/git-version.sh)
bin/%: $(GOFILES)
	mkdir -p $(dir $@)
	GOOS=$(word 1, $(subst /, ,$*)) GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $@ 

