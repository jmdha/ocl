#include <assert.h>
#include <jhttp/jhttp.h>

#include "web.h"
#include "db.h"

struct jhttp jhttp;

struct route {
	char* method;
	char* path;
	int (*fn)(struct jhttp_response* res, const struct jhttp_request* req);
};

int get_index(struct jhttp_response* res, const struct jhttp_request* req) {
	const char str[] =
		"<!DOCTYPE html>"
		"<html>"
		"	<body><p>Hello!</p></body>"
		"</html>";
	snprintf(res->body, JHTTP_RESPONSE_MAX, "%s", str);
	res->status = 200;
	return 0;
}

const struct route routes[] = {
	{ "GET",  "/", get_index },
};

int handler(struct jhttp_response* res, const struct jhttp_request* req) {
	for (size_t i = 0; i < sizeof(routes) / sizeof(struct route); i++)
		if (strcmp(req->method, routes[i].method) == 0 && 
		    strcmp(req->path,   routes[i].path)   == 0)
			return routes[i].fn(res, req);
	res->status = 404;
	return 0;
}


void web_init(struct db* db) {
	if (jhttp_init(&jhttp, 8000, handler) != 0) {
		fprintf(stderr, "failed to initialize jhttp on port 8000\n");
		exit(1);
	}
}

void web_fini() {
	jhttp_fini(&jhttp);
}

void web_step() {
	jhttp_poll(&jhttp);
}
