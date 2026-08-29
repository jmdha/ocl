#include <stddef.h>
#include <stdint.h>
#include <sys/random.h>
#include <sys/stat.h>
#include <time.h>

#include "utils.h"

int gen_ran(char *out, size_t size) {
	static const char hex[] = "0123456789abcdef";
	
	if (out == NULL)
		return -1;

	if (size == 0)
		return -1;
	
	for (size_t i = 0; i < size - 1; i++) {
		for (;;) {
			uint8_t tmp;
			if (getrandom(&tmp, 1, 0) != 1)
				return -1;
			if (tmp >= sizeof(hex) - 1)
				continue;
			out[i] = hex[tmp];
			break;
		}
	}

	out[size - 1] = '\0';
	
	return 0;
}

uint64_t current_time_ns(void) {
    struct timespec ts;
    clock_gettime(CLOCK_REALTIME, &ts);
    return (uint64_t)ts.tv_sec * 1000000000ULL + ts.tv_nsec;
}

size_t file_size(const char* path) {
	struct stat st;
	return stat(path, &st) == 0 ? (size_t) st.st_size : 0;
}
