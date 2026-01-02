# Makefile
.PHONY: build run docker-build docker-run test clean

build:
	go build -o bin/email-service cmd/server/main.go

run:
	go run cmd/server/main.go

docker-build:
	docker build -f docker/Dockerfile -t email-microservice .

docker-run:
	docker run -p 8080:8080 \
		-e SMTP_HOST=${SMTP_HOST} \
		-e SMTP_PORT=${SMTP_PORT} \
		-e SMTP_USERNAME=${SMTP_USERNAME} \
		-e SMTP_PASSWORD=${SMTP_PASSWORD} \
		-e FROM_EMAIL=${FROM_EMAIL} \
		email-microservice

test:
	go test -v ./...

clean:
	rm -f bin/email-service
