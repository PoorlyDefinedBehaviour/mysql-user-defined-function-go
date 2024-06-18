
build:
	bash -c "\
		export CGO_CFLAGS=$(shell mysql_config --include); \
		export CGO_ENABLED=1; \
		go build -buildmode=c-shared -o temporal.so main.go; \
		docker cp ./temporal.so mysql:/usr/lib64/mysql/plugin/temporal.so"

build-server-plugin:
	gcc -shared -o temporal.so $(shell mysql_config --include) -fPIC main.c -o temporal-server.so
	docker cp ./temporal-server.so mysql:/usr/lib64/mysql/plugin/temporal-server.so

	
install: build
	docker exec mysql mysql --user=root --password=password dev -e "DROP FUNCTION IF EXISTS temporal; CREATE FUNCTION IF NOT EXISTS temporal RETURNS STRING SONAME 'temporal.so';";

logs:
	docker exec mysql cat /var/log/mysqld.log

mysql:
	docker exec -it mysql mysql  --user=root --password=password dev --binary-as-hex=0