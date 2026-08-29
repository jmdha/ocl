#ifndef UTILS_H
#define UTILS_H

#include <stddef.h>
#include <stdint.h>

int gen_ran(char *out, size_t size);
uint64_t current_time_ns(void);

size_t file_size(const char* path);

#endif
