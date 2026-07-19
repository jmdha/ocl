#ifndef DB_H
#define DB_H

struct db;

int  db_init(struct db** db);
void db_fini(struct db* db);

int  db_upload_new(struct db* db);
int  db_upload_finish(struct db* db);

int  db_metrics_requests_add(struct db* db, int method_len, const char* method_buf, int uri_len, const char* uri_buf);
int  db_metrics_requests(struct db* db);

#endif
