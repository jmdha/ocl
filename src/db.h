#ifndef DB_H
#define DB_H

#include <stdint.h>
#include <stdbool.h>

typedef struct {
	uint64_t id;
} db_user;

typedef struct {
	uint64_t user_id;
	uint64_t ts;
	uint32_t lines;
} db_log;

int db_init();
int db_fini();

int db_user_new(db_user* user, char key[64]);
int db_user_get_by_id(db_user* user, uint64_t id);
int db_user_get_by_key(db_user* user, const char key[64]);

#endif
