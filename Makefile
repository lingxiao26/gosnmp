.PHONY: build

ALL: build

build:
	CGO_ENABLED=0 go build -o gosnmp .