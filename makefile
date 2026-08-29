all:
	gcc -std=c99 -pedantic -D_DEFAULT_SOURCE -Wall -o ocl \
	src/main.c src/db.c src/jbc.c src/http.c src/routes.c src/utils.c

watch:
	watchexec -w src -w web -e c,h,html --restart "make && ./ocl 8081"
