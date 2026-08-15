#ifndef DB_H
#define DB_H

#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>

typedef struct {
	char key[32];
} db_user;

typedef struct {
	size_t user_id;
	size_t ts;
	size_t lines;
} db_log;

typedef struct {
	size_t   user_id;
	char     method[8];
	char     path[64];
	uint8_t  ip[16];
	uint8_t  is_ip6;
	uint32_t status;
	uint64_t ts;
	uint64_t dur;
} db_request;

int db_init();
int db_fini();
int db_user_add(const db_user* user, size_t* id);
int db_user_get_by_id(const db_user** user, size_t id);
int db_user_get_by_key(const db_user** user, size_t* id, const char key[32]);
int db_request_add(const db_request* req, size_t* id);
int db_request_get_by_id(const db_request** req, size_t id);

#endif
