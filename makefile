all:
	gcc -Wall -ggdb -fsanitize=address -lcrypto -o ocl \
	src/main.c src/db.c src/jbc.c src/http.c src/utils.c \
	src/mongoose.c

watch:
	watchexec -w src -w web -e c,h,html --restart "make && ./ocl localhost:8081"
