#include <stdio.h>
#include <utest.h/utest.h>

UTEST_STATE();

int main(int argc, const char** argv) {
	if (utest_main(0, argv) != 0)
		return -1;
	printf("Hello World!\n");
}
