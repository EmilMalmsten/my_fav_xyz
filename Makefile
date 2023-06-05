include .env

migrateup:
	migrate -path sql/schema -database ${DB_URL} -verbose up

migratedown:
	migrate -path sql/schema -database ${DB_URL} -verbose down
