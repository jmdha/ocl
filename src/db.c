#include <lmdb.h>
#include <netinet/in.h>
#include <stddef.h>
#include <stdio.h>
#include <string.h>

#include "db.h"

MDB_env* env;
MDB_dbi  dbi_requests;
MDB_dbi  dbi_users;
MDB_dbi  dbi_logs;
MDB_dbi  dbi_log_ids;

int db_init() {
	int      rc  = 0;
	MDB_txn* txn = NULL;

	if (rc == 0) rc = mdb_env_create(&env);
	if (rc == 0) rc = mdb_env_set_maxdbs(env, 32);
	if (rc == 0) rc = mdb_env_set_mapsize(env, (size_t)16 * 1024 * 1024 * 1024);
	if (rc == 0) rc = mdb_env_open(env, "db.lmdb", MDB_NOSUBDIR | MDB_NOSYNC, 0664);

	if (rc == 0) rc = mdb_txn_begin(env, NULL, 0, &txn);
	if (rc == 0) rc = mdb_dbi_open(txn, "requests", MDB_CREATE | MDB_INTEGERKEY, &dbi_requests);
	if (rc == 0) rc = mdb_dbi_open(txn, "users",    MDB_CREATE | MDB_INTEGERKEY, &dbi_users);
	if (rc == 0) rc = mdb_dbi_open(txn, "logs",     MDB_CREATE | MDB_INTEGERKEY, &dbi_logs);
	if (rc == 0) rc = mdb_dbi_open(txn, "log_ids",  MDB_CREATE,                  &dbi_log_ids);
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
	mdb_dbi_close(env, dbi_users);
	mdb_dbi_close(env, dbi_logs);
	mdb_dbi_close(env, dbi_log_ids);
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

size_t db_size() {
	MDB_envinfo info;
	MDB_stat    stat;
	if (mdb_env_info(env, &info) != MDB_SUCCESS) return 0;
	if (mdb_env_stat(env, &stat) != MDB_SUCCESS) return 0;
	return (size_t)(info.me_last_pgno + 1) * stat.ms_psize;
}

size_t db_size_max() {
	MDB_envinfo info;
	mdb_env_info(env, &info);
	return info.me_mapsize;
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

int db_logs_get_id(uint32_t user_id, const char* file_name) {
	

	return 0;
}
