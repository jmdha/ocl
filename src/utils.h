#ifndef UTILS_H
#define UTILS_H

#include <stdint.h>

int gen_ran(uint8_t *out, size_t size);
uint64_t current_time_ns(void);

#endif
