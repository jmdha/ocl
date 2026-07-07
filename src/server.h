#ifndef SERVER_H
#define SERVER_H

struct conn;
struct request;
struct response;

void server_add(const char* path, void (*f)(struct request*, struct response*));
void server_listen(int port);

typedef enum {
	GET,
	POST,
	PATCH
} http_method;

typedef enum {
	STATUS_OK                 = 200,
	STATUS_BAD_REQUEST        = 400,
	STATUS_NOT_FOUND          = 404,
	STATUS_METHOD_NOT_ALLOWED = 405,
	STATUS_INTERNAL_SERVER_ERROR,
} http_status;

int req_read(struct request* req, char* buf, int cap);
void res_write(struct response* res, const char* buf, int len);
void res_status(struct response* res, http_status status);

#endif
