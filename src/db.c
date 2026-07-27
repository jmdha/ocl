#include <lmdb.h>
#include <stddef.h>
#include <stdio.h>

#include "db.h"

MDB_env* env;
MDB_dbi  dbi_requests;

int db_init() {
	int      rc  = 0;
	MDB_txn* txn = NULL;

	if (rc == 0) rc = mdb_env_create(&env);
	if (rc == 0) rc = mdb_env_set_maxdbs(env, 64);
	if (rc == 0) rc = mdb_env_open(env, "db.lmdb", MDB_NOSUBDIR | MDB_NOSYNC, 0664);
	if (rc == 0) rc = mdb_env_set_mapsize(env, (size_t)16 * 1024 * 1024 * 1024);
	if (rc == 0) rc = mdb_txn_begin(env, NULL, 0, &txn);
	if (rc == 0) rc = mdb_dbi_open(txn, "requests", MDB_CREATE | MDB_INTEGERKEY, &dbi_requests);
	if (rc == 0) rc = mdb_txn_commit(txn);

	if (rc != 0) {
		fprintf(stderr, "failed to initialize MDB: %s\n", mdb_strerror(rc));
		mdb_txn_abort(txn);
		mdb_env_close(env);
		return 1;
	}

	return 0;
}

void db_fini() {
	mdb_dbi_close(env, dbi_requests);
	mdb_env_close(env);
}

static int next_id(MDB_txn* txn, MDB_dbi dbi, size_t* out_id) {
	MDB_cursor* cur;
	MDB_val k, v;
	int rc = mdb_cursor_open(txn, dbi, &cur);
	if (rc != 0) return rc;

	rc = mdb_cursor_get(cur, &k, &v, MDB_LAST);
	if (rc == MDB_NOTFOUND) {
		*out_id = 1;
		rc = 0;
	} else if (rc == 0) {
		*out_id = *(size_t*)k.mv_data + 1;
	}

	mdb_cursor_close(cur);
	return rc;
}

int db_requests_add(const char* method, const char* uri) {
	size_t   id     = 0;
	int      rc     = 0;
	MDB_txn* txn    = NULL;
	int      offset = 0;
	char     buf[1024];

	offset += snprintf(&buf[offset], sizeof(buf) - offset, "%s", method) + 1;
	offset += snprintf(&buf[offset], sizeof(buf) - offset, "%s", uri) + 1;

	MDB_val key = { .mv_size = sizeof(id), .mv_data = &id };
	MDB_val val = { .mv_size = offset,     .mv_data = buf };

	if (rc == 0) rc = mdb_txn_begin(env, NULL, 0, &txn);
	if (rc == 0) rc = next_id(txn, dbi_requests, &id);
	if (rc == 0) rc = mdb_put(txn, dbi_requests, &key, &val, MDB_APPEND);
	if (rc == 0) rc = mdb_txn_commit(txn);

	if (rc != 0) {
		fprintf(stderr, "failed to add request: %s\n", mdb_strerror(rc));
		if (txn) mdb_txn_abort(txn);
		return 1;
	}

	return 0;
}

int db_logs_get_id(int user_id, const char* file_name) {
	return 0;
}
