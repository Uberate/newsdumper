.PHONY: build
output:
	mkdir ./output

build: output
	go build -o ./output/newsdumper cmd/main.go
	cp systemds ./output/

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