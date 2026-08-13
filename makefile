all: 
	gcc -Wall -O3 -ggdb -Ilib -o ocl \
	-lssl -lcrypto \
	-D LOG_INFO \
	src/main.c src/db.c src/log.c src/utils.c

watch:
	watchexec -w src -w web -e c,h,html --restart "make && ./ocl"
