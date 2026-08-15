#include <stddef.h>
#include <stdint.h>
#include <sys/random.h>
#include <time.h>

#include "utils.h"

int gen_ran(char *out, size_t size) {
	static const char hex[] = "0123456789abcdef";
	
	if (out == NULL || size == 0)
		return -1;
	
	size_t chars = size - 1;
	size_t bytes = (chars + 1) / 2;
	
	uint8_t buf[bytes];
	
	if (getrandom(buf, bytes, 0) != (ssize_t)bytes)
		return -1;
	
	for (size_t i = 0; i < chars; i++) {
		uint8_t nibble;
		
		if ((i & 1) == 0)
			nibble = buf[i / 2] >> 4;
		else
			nibble = buf[i / 2] & 0x0f;
		
		out[i] = hex[nibble];
	}
	
	out[chars] = '\0';
	return 0;
}

uint64_t current_time_ns(void) {
    struct timespec ts;
    clock_gettime(CLOCK_REALTIME, &ts);
    return (uint64_t)ts.tv_sec * 1000000000ULL + ts.tv_nsec;
}
