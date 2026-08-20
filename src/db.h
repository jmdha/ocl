#ifndef DB_H
#define DB_H

#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>

typedef struct {
	uint8_t key[32];
} db_user;

typedef struct {
	size_t user_id;
	size_t ts;
	size_t lines;
} db_log;

int db_init();
int db_fini();
int db_user_create(db_user* user, size_t* id);
int db_user_add(const db_user* user, size_t* id);
int db_user_get_by_id(const db_user** user, size_t id);
int db_user_get_by_key(const db_user** user, size_t* id, const char key[32]);

#endif
