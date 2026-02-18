.PHONY: build build-all install clean

build:
	go build -o kubectl-rollout .

build-all:
	@mkdir -p dist
	GOOS=linux   GOARCH=amd64 go build -o dist/kubectl-rollout-linux-amd64   .
	GOOS=linux   GOARCH=arm64 go build -o dist/kubectl-rollout-linux-arm64   .
	GOOS=darwin  GOARCH=amd64 go build -o dist/kubectl-rollout-darwin-amd64  .
	GOOS=darwin  GOARCH=arm64 go build -o dist/kubectl-rollout-darwin-arm64  .

install: build
	@mkdir -p $(HOME)/.krew/bin
	cp kubectl-rollout $(HOME)/.krew/bin/kubectl-rollout

clean:
	rm -f kubectl-rollout
	rm -rf dist/
