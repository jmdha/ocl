#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <lmdb.h>
#include <time.h>
#include <openssl/rand.h>

#include "db.h"
#include "utils.h"

#define NS_USER "id:user"
#define NS_LOG  "id:log"

MDB_env* env;
MDB_dbi  dbi_meta;
MDB_dbi  dbi_user;
MDB_dbi  dbi_log;

static int next_id(MDB_txn *txn, MDB_dbi dbi_meta, const char *namespace, uint64_t *out_id) {
	MDB_val key = { .mv_size = strlen(namespace), .mv_data = (void *)namespace };
	MDB_val val = { 0 };
	uint64_t id;
	int rc;
	
	rc = mdb_get(txn, dbi_meta, &key, &val);
	if (rc == MDB_NOTFOUND) {
		id = 1;
	} else if (rc == 0) {
		if (val.mv_size != sizeof(uint64_t)) return -1; // corruption guard
		memcpy(&id, val.mv_data, sizeof(id));
		id += 1;
	} else {
		return rc; // real LMDB error
	}
	
	val.mv_size = sizeof(id);
	val.mv_data = &id;
	rc = mdb_put(txn, dbi_meta, &key, &val, 0);
	if (rc != 0) return rc;
	
	*out_id = id;
	return 0;
}

int db_init() {
	int      rc  = 0;
	MDB_txn* txn = NULL;

	if (rc == 0) rc = mdb_env_create(&env);
	if (rc == 0) rc = mdb_env_set_maxdbs(env, 32);
	if (rc == 0) rc = mdb_env_set_mapsize(env, (size_t)1 * 1024 * 1024 * 1024);
	if (rc == 0) rc = mdb_env_open(env, "db.lmdb", MDB_NOSUBDIR | MDB_NOSYNC, 0664);

	if (rc == 0) rc = mdb_txn_begin(env, NULL, 0, &txn);
	if (rc == 0) rc = mdb_dbi_open(txn, "meta",   MDB_CREATE, &dbi_meta);
	if (rc == 0) rc = mdb_dbi_open(txn, "user",   MDB_CREATE, &dbi_user);
	if (rc == 0) rc = mdb_dbi_open(txn, "log",    MDB_CREATE | MDB_INTEGERKEY, &dbi_log);
	if (rc == 0) rc = mdb_txn_commit(txn);

	if (rc != 0) {
		fprintf(stderr, "failed to initialize MDB: %s\n", mdb_strerror(rc));
		mdb_txn_abort(txn);
		mdb_env_close(env);
		return 1;
	}

	return 0;
}

int db_fini() {
	mdb_env_close(env);
	return 0;
}

int db_user_new(db_user* user, char apikey[64]) {
	int      rc   = 0;
	MDB_val  key  = { 0 };
	MDB_val  val  = { 0 };
	MDB_txn* txn  = NULL;
	uint8_t  hash[32];

	key.mv_size = sizeof(hash);
	key.mv_data = hash;
	val.mv_size = sizeof(db_user);
	val.mv_data = user;

	if (rc == 0) rc = generate_api_key(apikey);
	if (rc == 0) rc = hash_api_key(hash, apikey);
	if (rc == 0) rc = mdb_txn_begin(env, NULL, 0, &txn);
	if (rc == 0) rc = next_id(txn, dbi_meta, NS_USER, &user->id);
	if (rc == 0) rc = mdb_put(txn, dbi_user, &key, &val, 0);
	if (rc == 0) rc = mdb_txn_commit(txn);
	if (rc != 0 && txn) mdb_txn_abort(txn);

	return rc;
}

int db_user_get_by_id(db_user* user, uint64_t id) {
	int         rc     = 0;
	int         found  = 0;
	MDB_txn*    txn    = NULL;
	MDB_cursor* cursor = NULL;
	MDB_val     key    = { 0 };
	MDB_val     val    = { 0 };
	
	if (rc == 0) rc = mdb_txn_begin(env, NULL, MDB_RDONLY, &txn);
	if (rc == 0) rc = mdb_cursor_open(txn, dbi_user, &cursor);
	if (rc == 0) rc = mdb_cursor_get(cursor, &key, &val, MDB_FIRST);
	
	while (rc == 0 && !found) {
		if (val.mv_size != sizeof(db_user)) { rc = -1; break; }
		
		if (((db_user*)val.mv_data)->id == id) {
			memcpy(user, val.mv_data, sizeof(db_user));
			found = 1;
			break;
		}
		
		rc = mdb_cursor_get(cursor, &key, &val, MDB_NEXT);
	}
	
	if (rc == MDB_NOTFOUND) rc = 0;
	if (rc == 0 && !found)  rc = MDB_NOTFOUND;
	
	if (cursor) mdb_cursor_close(cursor);
	if (txn)    mdb_txn_abort(txn);
	
	return rc;
}

int db_user_get_by_key(db_user* user, const char apikey[64]) {
	int      rc   = 0;
	MDB_val  key  = { 0 };
	MDB_val  val  = { 0 };
	MDB_txn* txn  = NULL;
	uint8_t  hash[32];

	key.mv_size = sizeof(hash);
	key.mv_data = hash;

	if (rc == 0) rc = hash_api_key(hash, apikey);
	if (rc == 0) rc = mdb_txn_begin(env, NULL, MDB_RDONLY, &txn);
	if (rc == 0) rc = mdb_get(txn, dbi_user, &key, &val);
	if (rc == 0) rc = val.mv_size != sizeof(db_user);
	if (rc == 0) memcpy(user, val.mv_data, sizeof(db_user));
	if (txn)     mdb_txn_abort(txn);

	return rc;
}
