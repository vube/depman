default: install run test doc

install:
	`cd ${GOPATH}/src/ && go clean ./...`
	go install ./...

run:
	depman
	depman help
	depman --help

	depman --verbose
	depman --debug

	depman --clear-cache
	depman install --skip-cache

test:
	go vet ./...
	golint .
	go test -i ./...
	go test -p 1 ./...

doc:
	godocdown . > Readme.md

html: doc
	markdown Readme.md > Readme.html
