FLAGS = -Ilib

all:
	gcc $(FLAGS) -o ocl \
	$$(find src -type f ! -name '*test*')

test:
	gcc $(FLAGS) -o testrunner \
	$$(find src -type f ! -name '*main*')

watch:
	watchexec -w . -e c,h,html --restart "make test && ./testrunner && make && ./ocl -p 8081"
