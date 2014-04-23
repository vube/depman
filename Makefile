test:
	`cd ${GOPATH}/src/ && go clean ./...`
	go install ./...
	depman
	depman help
	depman --help

	depman --verbose
	depman --debug

	depman --clear-cache
	depman install --skip-cache

	go vet ./...
	golint .
	go test -i ./...
	go test -p 1 ./...

