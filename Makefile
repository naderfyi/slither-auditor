build:
	go build -o bin/solidityAuditor

run: build
	./bin/solidityAuditor