#ifndef DB_H
#define DB_H

#include <stddef.h>

typedef enum {
	DB_UPLOAD_STATUS_UPLOADING,
	DB_UPLOAD_STATUS_DONE,
	DB_UPLOAD_STATUS_ERROR
} db_upload_status;

void db_init(const char* conn);
int  db_add_upload();
int  db_add_character(const char* guid_ptr, size_t guid_len, const char* name_ptr, size_t name_len);
int  db_set_upload_path(int id, const char* path_ptr, size_t path_len);
int  db_set_upload_done(int id, size_t size);
int  db_set_upload_error(int id, const char* err_ptr, size_t err_len);

#endif
