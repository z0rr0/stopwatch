PROGRAM=stopwatch
BIN=bin/stopwatch
VERSION=`bash version.sh`
SOURCEDIR=src/github.com/z0rr0/stopwatch


all: test

install:
	go install -ldflags "$(VERSION)" github.com/z0rr0/stopwatch

lint: install
	go vet github.com/z0rr0/stopwatch
	golint github.com/z0rr0/stopwatch

run: lint
	$(GOPATH)/$(BIN)

test: lint
	# go tool cover -html=ratest_coverage.out
	# go tool trace ratest.test trace.out
	go test -race -v -cover -coverprofile=ratest_coverage.out -trace trace.out github.com/z0rr0/stopwatch

clean:
	rm -f $(GOPATH)/$(BIN)
