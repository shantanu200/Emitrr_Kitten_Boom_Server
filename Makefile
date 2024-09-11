build: 
	cd ./ && env GOOS=linux CGO_ENABLED=0 go build -o bin/server main.go

run: build
	cd ./bin && ./server
