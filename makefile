all:
	gcc -Wall -ggdb -lsqlite3 -lcjson -Ilib -o ocl \
	src/main.c src/server.c src/http.c src/db.c src/utils.c src/worker.c
