FLAGS = -Ilib

all:
	gcc $(FLAGS) src/* -o ocl

watch:
	watchexec -w . -e c,h,html --restart "make && ./ocl -p 8081"
