all:
	gcc -Wall -ggdb -lsqlite3 -lcjson -Ilib -o ocl \
	src/main.c src/server.c src/http.c src/db.c src/utils.c src/worker.c src/log.c

watch:
	watchexec -w src -w web -e c,h,html --restart "make && ./ocl"
