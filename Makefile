include .env

DSN=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)
PROFILING=localhost:3435 http://localhost:8000

network:
	sudo docker network create dev-net

migrateup:
	migrate -path ./db/migrations -database $(DSN) -verbose up $(STEPS)

migratedown:
	migrate -path ./db/migrations -database $(DSN) -verbose down $(STEPS)
	
redis:
	docker run --network=host --name=redis -d redis:7.4.0-alpine

server:
	GIN_MODE=release go run ./cmd/app/main.go
	
worker:
	go run ./cmd/worker/main.go

data:
	go run ./cmd/seeding/main.go

backendup:
	sudo docker compose up $(BUILD) $(D)

backenddown:
	sudo docker compose down

nginxup:
	sudo docker compose -f docker-compose.nginx.yml up $(D)

nginxdown:
	sudo docker compose -f docker-compose.nginx.yml down

dockerup:
	sudo docker compose -f docker-compose.yml -f ../puxing-fe/docker-compose.yml -f docker-compose.nginx.yml up

dockerdown:
	sudo docker compose -f docker-compose.yml -f ../puxing-fe/docker-compose.yml -f docker-compose.nginx.yml down

server_nodemon:
	nodemon --exec go run ./cmd/app/main.go --signal SIGTERM

# profiling:
#     go tool pprof -http $(PROFILING)

cacheup:
	sudo docker compose -f docker-compose.memcached.yml up	

cachedown:
	sudo docker compose -f docker-compose.memcached.yml down

.PHONY: migrateup migratedown redis server worker