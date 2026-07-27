#ifndef DB_H
#define DB_H

int  db_init();
void db_fini();

int  db_requests_add(const char* method, const char* uri);

int  db_logs_get_id(int user_id, const char* file_name);

#endif
