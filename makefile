all: 
	gcc -Wall -ggdb -lsqlite3 -Ilib -o ocl \
	src/main.c src/db.c src/web.c \
	lib/mongoose/mongoose.c

watch:
	watchexec -w src -w web -e c,h,html --restart "make && ./ocl"
