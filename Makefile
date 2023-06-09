include .env

migrateup:
	migrate -path sql/schema -database ${DB_URL} -verbose up

migratedown:
	migrate -path sql/schema -database ${DB_URL} -verbose down

testdb_migrateup:
	migrate -path sql/schema -database ${TEST_DB_URL} -verbose up

testdb_migratedown:
	migrate -path sql/schema -database ${TEST_DB_URL} -verbose down
