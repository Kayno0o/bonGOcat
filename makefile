docker:
	docker-compose up --build

run:
	./app/bongo

prod: docker run