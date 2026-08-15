all: compress
	gcc -Wall -ggdb -fsanitize=address -lcrypto -o ocl \
	src/main.c src/db.c src/jbc.c src/web.c src/utils.c \
	src/mongoose.c

watch:
	watchexec -w src -w web -e c,h,html --restart "make && ./ocl"

compress:
		find ./web -type f \( -name '*.html' -o -name '*.css' -o -name '*.js' -o -name '*.svg' -o -name '*.json' \) -exec gzip -k -f {} \;
