export APP := dnsbuster 
VERSION=$(shell git describe --tags --always --dirty)

all: compile

test: 
	go vet ./...
	go test -v ./...

compile:
	go build -buildvcs=false -o `pwd`/$(APP)
	strip -s `pwd`/$(APP)

clean:
	rm -rf $(APP)  &2> /dev/null
	go clean

.PHONY: all compile clean
