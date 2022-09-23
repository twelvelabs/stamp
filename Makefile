.PHONY: coverage
coverage:
	make test
	go tool cover -html=coverage.tmp

.PHONY: clean
clean:
	rm -Rf ./bin
	rm coverage.tmp

.PHONY: generate
generate:
	go generate ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -cover -coverprofile=coverage.tmp ./...

.PHONY: bin/stamp
bin/stamp:
	go build -trimpath -o ./bin/stamp ./cmd/stamp

.PHONY: build
build: bin/stamp

prefix  := /usr/local
bindir  := ${prefix}/bin

.PHONY: install
install: bin/stamp
	install -d ${bindir}
	install -m755 bin/stamp ${bindir}/
