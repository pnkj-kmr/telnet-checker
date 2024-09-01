
build_linux:
	env GOOS=linux GOARCH=386 go build -o telnetchecker main.go

b:
	go build -o telnetchecker main.go

r:
	go run main.go

run:
	./telnetchecker -f samples/input.json -o samples/output.json

release:
	goreleaser release --snapshot --clean 

build:
	goreleaser build --single-target --snapshot --clean

