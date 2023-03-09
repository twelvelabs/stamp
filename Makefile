.PHONY: coverage
coverage:
	make test
	go tool cover -html=coverage.out

.PHONY: clean
clean:
	rm -Rf ./bin
	rm coverage.out

.PHONY: generate
generate:
	go generate ./...

.PHONY: lint
lint:
	stylist check

.PHONY: test
test:
	go mod tidy
	go test -cover -coverprofile=coverage.out ./...
	@cat coverage.out | grep -v "_mock.go" | grep -v "_enum.go" > coverage.out.new
	@mv coverage.out.new coverage.out

.PHONY: dist/stamp
dist/stamp:
	go mod tidy
	go build -trimpath -o ./dist/stamp ./cmd/stamp

.PHONY: build
build: dist/stamp

dst_dir  := /usr/local/bin

.PHONY: install
install: dist/stamp
	install -d ${dst_dir}
	install -m755 bin/stamp ${dst_dir}/
