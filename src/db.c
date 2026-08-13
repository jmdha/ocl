#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <openssl/rand.h>
#include <jbc/jbc.h>

#include "db.h"
#include "utils.h"

jbc db_users;
jbc db_logs;

int db_init() {
	int rc = 0;

	if (rc == 0) rc = jbc_init(&db_users, "db/users.jbc", sizeof(db_user), 1 * 1024 * 1024 * 1024);
	if (rc == 0) rc = jbc_init(&db_logs,  "db/logs.jbc",  sizeof(db_log),  1 * 1024 * 1024 * 1024);

	return rc;
}

int db_fini() {
	int rc = 0;

	if (rc == 0) jbc_fini(&db_users);
	if (rc == 0) jbc_fini(&db_logs);

	return rc;
}

size_t db_size() {
	size_t size = 0;

	size += jbc_size(&db_users);
	size += jbc_size(&db_logs);

	return size;
}

size_t db_size_max() {
	size_t size = 0;

	size += db_users.max;
	size += db_logs.max;

	return size;
}

int db_user_add(const db_user* user, size_t* id) {
	return jbc_add(&db_users, id, user);
}

int db_user_get_by_id(const db_user** user, size_t id) {
	return jbc_ref(&db_users, id, (const void**) user);
}

int db_user_get_by_key(const db_user** user, size_t* id, const uint8_t apikey_hash[32]) {
	for (size_t i = 0; i < jbc_len(&db_users); i++) {
		if (jbc_ref(&db_users, i, (const void**) user) != 0)
			return 1;
		if (memcmp((*user)->apikey_hash, apikey_hash, sizeof((*user)->apikey_hash)) == 0) {
			*id = i;
			return 0;
		}
	}
	return 1;
}

