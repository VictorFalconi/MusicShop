# Network ( in docker)
network:
	docker network create my_network
# Golang
golang:
	docker run --name backend_golang -p 6868:6868 -d golang:latest
	#docker run -d --net my_network --name backend_golang -p 6868:6868 -d golang:latest

# DB
postgres:
	docker run --name my_postgres -p 41234:5432 -e POSTGRES_USER=thanhliem -e POSTGRES_PASSWORD=liem1234 -d postgres:latest
	#docker run -d --net my_network --name my_postgres -p 41234:5432 -e POSTGRES_USER=thanhliem -e POSTGRES_PASSWORD=liem1234 -d postgres:latest

createdb:
	docker exec -it my_postgres createdb --username=thanhliem musicshop

dropdb:
	docker exec -it my_postgres dropdb --username=thanhliem musicshop

exportdb:
	docker exec -it my_postgres pg_dump --username=thanhliem musicshop > db/musicshop.sql

#Shell postgres
postgres_sh:
	docker exec -it my_postgres /bin/sh

	# back up
	#pg_dump -U postgres -h localhost postgres_musicshop >> musicshop.sql

migrateinit:
	migrate create -ext sql -dir db/migration -seq init_schema

migrateup:
	migrate -path db/migration -database "postgresql://thanhliem:liem1234@localhost:41234/musicshop?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://thanhliem:liem1234@localhost:41234/musicshop?sslmode=disable" -verbose down


.PHONY: postgres createdb dropdb postgres_sh migrateinit migrateup migratedown