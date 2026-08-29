#ifndef ROUTES_H
#define ROUTES_H

#include <stddef.h>

int post_users(const char* headers[], size_t count, char* buf, size_t size);
int post_login(const char* headers[], size_t count, char* buf, size_t size);
int get_upload(const char* headers[], size_t count, char* buf, size_t size);
int post_upload(const char* headers[], size_t count, char* buf, size_t size);

#endif
