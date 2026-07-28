all: 
	gcc -Wall -ggdb -std=c99 -lzstd -llmdb -Ilib -o ocl \
	src/main.c src/db.c src/web.c src/utils.c

watch:
	watchexec -w src -w web -e c,h,html --restart "make && ./ocl"
