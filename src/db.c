#include <sqlite3.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "db.h"

sqlite3* DB;

static const char* SQL_INIT = 
"create table if not exists requests (   "
"	id integer primary key,          "
"       guid text not null,              "
"       name text not null               "
");                                      "
"create table if not exists uploads (    "
"	id integer primary key,          "
"       status text not null,            "
"       path text,                       "
"       size integer,                    "
"       error text,                      "
"       check (status in (               "
"              'uploading',              "
"              'done',                   "
"              'error'                   "
"             )                          "
"       )                                "
");                                      "
"create table if not exists characters ( "
"	id integer primary key,          "
"       guid text not null unique,       "
"       name text not null               "
");                                      "
"create table if not exists spells (     "
"	id integer primary key,          "
"       name text not null               "
");                                      ";

void db_init(const char* conn) {
	int   err;
	char* msg;
	if ((err = sqlite3_open(conn, &DB)) != SQLITE_OK) {
		fprintf(stderr, "Cannot open database: %s\n", sqlite3_errmsg(DB));
    		sqlite3_close(DB);
    		exit(EXIT_FAILURE);
	}
	if ((err = sqlite3_exec(DB, SQL_INIT, NULL, NULL, &msg)) != SQLITE_OK) {
		fprintf(stderr, "SQL error: %s\n", msg);
		sqlite3_free(msg);
		sqlite3_close(DB);
		exit(EXIT_FAILURE);
	}
}

int db_add_upload() {
	static const char*   sql =
	"insert into uploads (status) "
	"values ('uploading')         "
	"returning id;                ";
	static sqlite3_stmt* stmt = NULL;
	int                  err;
	int                  id;

	if (!stmt && ((err = sqlite3_prepare_v2(DB, sql, -1, &stmt, NULL)) != SQLITE_OK)) {
		fprintf(stderr, "Prepare failed: %s\n", sqlite3_errmsg(DB));
        	sqlite3_close(DB);
		exit(EXIT_FAILURE);
	}

	if (sqlite3_step(stmt) == SQLITE_ROW) {
		id = sqlite3_column_int(stmt, 0);
		sqlite3_reset(stmt);
		return id;
	} else {
		fprintf(stderr, "not row\n");
		sqlite3_reset(stmt);
		return -1;
	}
}

int db_add_character(const char* guid_ptr, size_t guid_len, const char* name_ptr, size_t name_len) {
	static const char*   sql = 
	"insert into characters (guid, name) "
	"values (?, ?)                       "
	"on conflict (guid) do nothing       "
	"returning id;                       ";
	static sqlite3_stmt* stmt = NULL;
	int                  err;
	int                  id;

	if (!stmt && ((err = sqlite3_prepare_v2(DB, sql, -1, &stmt, NULL)) != SQLITE_OK)) {
		fprintf(stderr, "Prepare failed: %s\n", sqlite3_errmsg(DB));
        	sqlite3_close(DB);
		exit(EXIT_FAILURE);
	}

	sqlite3_bind_text(stmt, 1, guid_ptr, guid_len, SQLITE_STATIC);
	sqlite3_bind_text(stmt, 2, name_ptr, name_len, SQLITE_STATIC);
	if (sqlite3_step(stmt) == SQLITE_ROW) {
		id = sqlite3_column_int(stmt, 0);
		sqlite3_reset(stmt);
		return id;
	} else {
		sqlite3_reset(stmt);
		return -1;
	}
}

int db_set_upload_path(int id, const char* path_ptr, size_t path_len) {
	static const char*   sql =
	"update uploads set path = ? "
	"where id = ?;               ";
	static sqlite3_stmt* stmt = NULL;
	int                  err;

	if (!stmt && ((err = sqlite3_prepare_v2(DB, sql, -1, &stmt, NULL)) != SQLITE_OK)) {
		fprintf(stderr, "Prepare failed: %s\n", sqlite3_errmsg(DB));
        	sqlite3_close(DB);
		exit(EXIT_FAILURE);
	}

	sqlite3_bind_text(stmt, 1, path_ptr, path_len, SQLITE_STATIC);
	sqlite3_bind_int(stmt, 2, id);
	if (sqlite3_step(stmt) == SQLITE_DONE) {
		sqlite3_reset(stmt);
		return 0;
	} else {
		sqlite3_reset(stmt);
		return -1;
	}
}

int db_set_upload_done(int id, size_t size) {
	static const char*   sql =
	"update uploads set status = ?, size = ? "
	"where id = ?;                           ";
	static sqlite3_stmt* stmt = NULL;
	int                  err;

	if (!stmt && ((err = sqlite3_prepare_v2(DB, sql, -1, &stmt, NULL)) != SQLITE_OK)) {
		fprintf(stderr, "Prepare failed: %s\n", sqlite3_errmsg(DB));
        	sqlite3_close(DB);
		exit(EXIT_FAILURE);
	}

        sqlite3_bind_text(stmt, 1, "done", strlen("done"), SQLITE_STATIC);
	sqlite3_bind_int(stmt, 2, size);
	sqlite3_bind_int(stmt, 3, id);
	if (sqlite3_step(stmt) == SQLITE_DONE) {
		sqlite3_reset(stmt);
		return 0;
	} else {
		sqlite3_reset(stmt);
		return -1;
	}
}

int db_set_upload_error(int id, const char* err_ptr, size_t err_len) {
	static const char*   sql =
	"update uploads set status = ?, error = ? "
	"where id = ?;                            ";
	static sqlite3_stmt* stmt = NULL;
	int                  err;

	if (!stmt && ((err = sqlite3_prepare_v2(DB, sql, -1, &stmt, NULL)) != SQLITE_OK)) {
		fprintf(stderr, "Prepare failed: %s\n", sqlite3_errmsg(DB));
        	sqlite3_close(DB);
		exit(EXIT_FAILURE);
	}

        sqlite3_bind_text(stmt, 1, "error", strlen("error"), SQLITE_STATIC);
        sqlite3_bind_text(stmt, 2, err_ptr, err_len, SQLITE_STATIC);
	sqlite3_bind_int(stmt, 3, id);
	if (sqlite3_step(stmt) == SQLITE_DONE) {
		sqlite3_reset(stmt);
		return 0;
	} else {
		sqlite3_reset(stmt);
		return -1;
	}
}
