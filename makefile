all: 
	gcc -O3 -Wall -ggdb -std=c23 -lzstd -llmdb -Ilib -o ocl \
	src/main.c src/db.c src/web.c

watch:
	watchexec -w src -w web -e c,h,html --restart "make && ./ocl"
