#ifndef DB_H
#define DB_H

#include <stdint.h>
#include <stddef.h>

int    db_init();
void   db_fini();
size_t db_size();
size_t db_size_max();
int    db_requests_add(const char* method, const char* uri);
int    db_users_get_id();
int    db_logs_get_id(uint32_t user_id, const char* file_name);

#endif
