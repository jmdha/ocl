#ifndef SERVER_H
#define SERVER_H

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
} http_status;

void res_write(struct response* res, const char* buf, int len);
void res_status(struct response* res, http_status status);

#endif
