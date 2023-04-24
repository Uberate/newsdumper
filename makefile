output:
	mkdir ./output

.PHONY: test
test:
	go test ./...

.PHONY: build
build: output
	go build -o ./output/darwin-amd64-newsdumper cmd/bin/main.go
	cp -rf systemds ./output/
	cp makefile output/makefile

.PHONY: clean
clean:
	rm -rf ./output

.PHONY: linux-clean
linux-clean:
	rm -rf /usr/local/bin/newsdumper

PATH := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

.PHONY: linux-install
linux-install:
	cp ./output/newsdumper /usr/local/bin/newsdumper
	mkdir /var/
	systemctl enable $(path)output/newsdumper.service
	systemctl enable $(path)newsdumper.timer