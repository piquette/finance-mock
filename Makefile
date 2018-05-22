build:
	go build -v

test:
	go test -v ./...

vet:
	go vet ./...

dev: build
	./finance-mock --verbose
