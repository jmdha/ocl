#include <sqlite3.h>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>

#include "db.h"

const char schema[] = {
	#embed "../schema.sql"
};

const char* sql_requests_start = "insert into requests (start) values (?)";
const char* sql_requests_fill  = "update requests set ip = ?, method = ?, uri = ? where id = ?";
const char* sql_requests_end   = "update requests set end = ? where id = ?;";

const char* sql_requests_count = "select count(*) from requests;";

struct db {
	sqlite3* db;
	sqlite3_stmt* requests_start;
	sqlite3_stmt* requests_fill;
	sqlite3_stmt* requests_end;
	sqlite3_stmt* requests_count;
};

int db_init(struct db** db) {
	*db = calloc(1, sizeof(struct db));

	if (*db == NULL) {
		fprintf(stderr, "out of memory\n");
		return 1;
	}

	if (sqlite3_open("abc.db", &(*db)->db) != SQLITE_OK) {
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

	if (sqlite3_prepare_v2((*db)->db, sql_requests_start, -1, &(*db)->requests_start, NULL) != SQLITE_OK) {
		fprintf(stderr, "sql err: %s\n", sqlite3_errmsg((*db)->db));
		db_fini(*db);
		*db = NULL;
		return 1;
	}

	if (sqlite3_prepare_v2((*db)->db, sql_requests_fill, -1, &(*db)->requests_fill, NULL) != SQLITE_OK) {
		fprintf(stderr, "sql err: %s\n", sqlite3_errmsg((*db)->db));
		db_fini(*db);
		*db = NULL;
		return 1;
	}

	if (sqlite3_prepare_v2((*db)->db, sql_requests_end, -1, &(*db)->requests_end, NULL) != SQLITE_OK) {
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

	if (db->requests_start)   sqlite3_finalize(db->requests_start);
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

int db_requests_start(struct db* db) {
	sqlite3_bind_int(db->requests_start, 1, time(NULL));

	if (sqlite3_step(db->requests_start) != SQLITE_DONE) {
		fprintf(stderr, "sql stmt failed with: %s\n", sqlite3_errmsg(db->db));
		sqlite3_reset(db->requests_start);
		sqlite3_clear_bindings(db->requests_start);
		return -1;
	}

	int id = sqlite3_last_insert_rowid(db->db);
	sqlite3_reset(db->requests_start);
	sqlite3_clear_bindings(db->requests_start);
	return id;
}

int db_requests_fill(struct db* db, int id, int ip_len, const char* ip_buf, int method_len, const char* method_buf, int uri_len, const char* uri_buf) {
	sqlite3_bind_text(db->requests_fill, 1, ip_buf, ip_len, SQLITE_STATIC);
	sqlite3_bind_text(db->requests_fill, 2, method_buf, method_len, SQLITE_STATIC);
	sqlite3_bind_text(db->requests_fill, 3, uri_buf, uri_len, SQLITE_STATIC);
	sqlite3_bind_int(db->requests_fill,  4, id);

	if (sqlite3_step(db->requests_fill) != SQLITE_DONE) {
		fprintf(stderr, "sql stmt failed with: %s\n", sqlite3_errmsg(db->db));
		sqlite3_reset(db->requests_fill);
		sqlite3_clear_bindings(db->requests_fill);
		return -1;
	}

	sqlite3_reset(db->requests_fill);
	sqlite3_clear_bindings(db->requests_fill);
	return 0;
}

int db_requests_end(struct db* db, int id) {
	sqlite3_bind_int(db->requests_end, 1, time(NULL));
	sqlite3_bind_int(db->requests_end, 2, id);

	if (sqlite3_step(db->requests_end) != SQLITE_DONE) {
		fprintf(stderr, "sql stmt failed with: %s\n", sqlite3_errmsg(db->db));
		sqlite3_reset(db->requests_end);
		sqlite3_clear_bindings(db->requests_end);
		return 1;
	}

	sqlite3_reset(db->requests_end);
	sqlite3_clear_bindings(db->requests_end);
	return 0;
}

int db_requests(struct db* db) {
	int out = -1;

	if (sqlite3_step(db->requests_count) == SQLITE_ROW)
		out = sqlite3_column_int(db->requests_count, 0);
	else
		fprintf(stderr, "sql stmt failed with: %s\n", sqlite3_errmsg(db->db));

	sqlite3_reset(db->requests_count);
	sqlite3_clear_bindings(db->requests_count);
	return out;
}
