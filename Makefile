PROGRAM=stopwatch
BIN=bin/stopwatch
VERSION=`bash version.sh`
SOURCEDIR=src/github.com/z0rr0/stopwatch


all: lint

install:
	go install -ldflags "$(VERSION)" github.com/z0rr0/stopwatch

lint: install
	go vet github.com/z0rr0/stopwatch
	golint github.com/z0rr0/stopwatch

run: lint
	$(GOPATH)/$(BIN)

test: lint
	# go tool cover -html=coverage.out
	# go tool trace trace.test trace.out
	go test -race -v -cover -coverprofile=ratest_coverage.out -trace trace.out github.com/z0rr0/stopwatch

clean:
	rm -f $(GOPATH)/$(BIN)
