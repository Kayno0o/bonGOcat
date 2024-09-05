.PHONY: build, run
build:
	go build -o dist/bongo
	chmod +x dist/bongo

run:
	./dist/bongo
