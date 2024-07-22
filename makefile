.PHONY: build
build:
	go build -o dist/bongo
	chmod +x dist/bongo

.PHONY: run
run:
	./dist/bongo
