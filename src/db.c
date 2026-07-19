#include <sqlite3.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "db.h"

const char* schema =
"create table if not exists requests ( "
"	id integer primary key,        "
"       method string not null,        "
"       uri string not null            "
");                                    ";

const char* sql_requests_add = "insert into requests (method, uri) values (?, ?);";

const char* sql_requests_count = "select count(*) from requests;";

struct db {
	sqlite3* db;
	sqlite3_stmt* requests_add;
	sqlite3_stmt* requests_count;
};

int db_init(struct db** db) {
	*db = calloc(1, sizeof(struct db));

	if (*db == NULL) {
		fprintf(stderr, "out of memory\n");
		return 1;
	}

	if (sqlite3_open(":memory:", &(*db)->db) != SQLITE_OK) {
		fprintf(stderr, "cannot open database: %s\n", sqlite3_errmsg((*db)->db));
		db_fini(*db);
		*db = NULL;
		return 1;
	}

	if (sqlite3_exec((*db)->db, schema, NULL, NULL, NULL) != SQLITE_OK) {
		fprintf(stderr, "SQL err: %s\n", sqlite3_errmsg((*db)->db));
		db_fini(*db);
		*db = NULL;
		return 1;
	}

	if (sqlite3_prepare_v2((*db)->db, sql_requests_add, -1, &(*db)->requests_add, NULL) != SQLITE_OK) {
		fprintf(stderr, "sql err: %s\n", sqlite3_errmsg((*db)->db));
		db_fini(*db);
		*db = NULL;
		return 1;
	}

	if (sqlite3_prepare_v2((*db)->db, sql_requests_count, -1, &(*db)->requests_count, NULL) != SQLITE_OK) {
		fprintf(stderr, "sql err: %s\n", sqlite3_errmsg((*db)->db));
		db_fini(*db);
		*db = NULL;
		return 1;
	}

	return 0;
}

void db_fini(struct db* db) {
	if (!db) return;

	if (db->requests_add)   sqlite3_finalize(db->requests_add);
	if (db->requests_count) sqlite3_finalize(db->requests_count);
	if (db->db)             sqlite3_close(db->db);
	free(db);
}

int db_upload_new(struct db* db) {
	return 0;
}

int db_upload_finish(struct db* db) {
	return 0;
}

int db_metrics_requests_add(struct db* db, int method_len, const char* method_buf, int uri_len, const char* uri_buf) {
	sqlite3_bind_text(db->requests_add, 1, method_buf, method_len, SQLITE_STATIC);
	sqlite3_bind_text(db->requests_add, 2, uri_buf, uri_len, SQLITE_STATIC);

	if (sqlite3_step(db->requests_add) != SQLITE_DONE) {
		fprintf(stderr, "sql stmt failed with: %s\n", sqlite3_errmsg(db->db));
		sqlite3_reset(db->requests_add);
		sqlite3_clear_bindings(db->requests_add);
		return 1;
	}

	sqlite3_reset(db->requests_add);
	sqlite3_clear_bindings(db->requests_add);
	return 0;
}

int db_metrics_requests(struct db* db) {
	int out = -1;

	if (sqlite3_step(db->requests_count) == SQLITE_ROW)
		out = sqlite3_column_int(db->requests_count, 0);
	else
		fprintf(stderr, "sql stmt failed with: %s\n", sqlite3_errmsg(db->db));

	sqlite3_reset(db->requests_count);
	sqlite3_clear_bindings(db->requests_count);
	return out;
}
