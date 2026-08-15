#include <string.h>
#include <stdio.h>
#include <stdlib.h>

#include "db.h"
#include "jbc.h"

jbc* db_users;
jbc* db_logs;

int db_init() {
	int rc = 0;

	if (rc == 0) rc = jbc_init(&db_users, "db/users.jbc", sizeof(db_user), 1 * 1024 * 1024 * 1024);
	if (rc == 0) rc = jbc_init(&db_logs,  "db/logs.jbc",  sizeof(db_log),  1 * 1024 * 1024 * 1024);

	return rc;
}

int db_fini() {
	int rc = 0;

	if (rc == 0) jbc_fini(db_users);
	if (rc == 0) jbc_fini(db_logs);

	return rc;
}

int db_user_add(const db_user* user, size_t* id) {
	return jbc_add(db_users, id, user);
}

int db_user_get_by_id(const db_user** user, size_t id) {
	return jbc_ref(db_users, id, (const void**) user);
}

int db_user_get_by_key(const db_user** user, size_t* id, const char key[32]) {
	for (size_t i = 0; i < jbc_len(db_users); i++) {
		if (jbc_ref(db_users, i, (const void**) user) != 0)
			return 1;
		if (memcmp((*user)->key, key, sizeof((*user)->key) - 1) == 0) {
			*id = i;
			return 0;
		}
	}
	return 1;
}
