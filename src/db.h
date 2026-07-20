#ifndef DB_H
#define DB_H

struct db;

int  db_init(struct db** db);
void db_fini(struct db* db);

int  db_upload_new(struct db* db);
int  db_upload_finish(struct db* db);

int  db_requests_start(struct db* db);
int  db_requests_fill(struct db* db, int id, int ip_len, const char* ip_buf, int method_len, const char* method_buf, int uri_len, const char* uri_buf);
int  db_requests_end(struct db* db, int id);
int  db_requests(struct db* db);

#endif
