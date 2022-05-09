build:
	docker-compose build botfreelancer

run:
	docker-compose up botfreelancer

test:
	go test -v ./...

migrate:
	migrate -path ./schema -database postgres://postgres:qwerty@0.0.0.0:5436/postgres?sslmode=disable up