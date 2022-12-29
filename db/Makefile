postgres:
	docker run --name my_postgres -p 41234:5432 -e POSTGRES_USER=thanhliem -e POSTGRES_PASSWORD=liem1234 -d postgres:latest

createdb:
	docker exec -it my_postgres createdb --username=thanhliem music_shop

dropdb:
	docker exec -it my_postgres dropdb --username=thanhliem music_shop

exportdb:
	docker exec -it my_postgres pg_dump --username=thanhliem music_shop > db/music_shop.sql

#Shell postgres
postgres_sh:
	docker exec -it my_postgres /bin/sh

	# back up
	#pg_dump -U postgres -h localhost postgres_musicshop >> musicshop.sql

migrateinit:
	migrate create -ext sql -dir db/migration -seq init_schema

migrateup:
	migrate -path db/migration -database "postgresql://thanhliem:liem1234@localhost:41234/music_shop?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://thanhliem:liem1234@localhost:41234/music_shop?sslmode=disable" -verbose down


.PHONY: postgres createdb dropdb postgres_sh migrateinit migrateup migratedown