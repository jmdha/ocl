#ifndef DB_H
#define DB_H

#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "jbc.h"

typedef struct {
	// TODO: Store hash instead of key directly
	char key[32];
} db_user;

typedef struct {
	size_t user_id;
	char filename[64];
} db_log;

int db_init(void);
int db_fini(void);
int db_user_add(const db_user* user, size_t* id);
int db_user_get(const db_user** user, size_t id);
int db_log_add(const db_log* log, size_t* id);
int db_log_get(const db_log** log, size_t id);

int db_user_get_id(size_t* id, const char* key);
int db_log_get_id(size_t* id, size_t user_id, const char* filename);

#endif
