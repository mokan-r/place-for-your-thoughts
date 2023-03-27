MODULE=github.com/mokan-r/place-for-your-thoughts

run: run_postgres migrate_up build run_server

build:
	go mod tidy
	go build -o server ${MODULE}/cmd/main

run_server:
	./server

run_postgres:
	docker-compose up --detach

migrate_up:
	migrate -path ./migrations -database 'postgres://mdaniell:mdaniell@localhost:5432/superheroes?sslmode=disable' up

attack:
	echo 'GET http://localhost:8888' | vegeta attack -rate 200 -duration 1s | vegeta report

clean:
	@rm -rf server